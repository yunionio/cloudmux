// Copyright 2026 Yunion
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package volcengine

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"yunion.io/x/jsonutils"
	"yunion.io/x/pkg/errors"

	billing "yunion.io/x/cloudmux/pkg/apis/billing"
	api "yunion.io/x/cloudmux/pkg/apis/compute"
	"yunion.io/x/cloudmux/pkg/cloudprovider"
	"yunion.io/x/cloudmux/pkg/multicloud"
)

const postgresqlAPIVersion = "2022-01-01"

var postgresqlRegionPattern = regexp.MustCompile(`^[a-z0-9]+(?:-[a-z0-9]+)+$`)

// PostgreSQL uses regional HTTPS endpoints and JSON POST requests, unlike
// several of the older Volcengine APIs handled by jsonRequest.
type sPostgresqlClient struct {
	httpClient *http.Client
	regionLock sync.Mutex
	regions    map[string]bool
}

func (client *SVolcEngineClient) getPostgresqlClient() *sPostgresqlClient {
	client.lock.Lock()
	defer client.lock.Unlock()
	if client.postgresqlClient == nil {
		transport := http.DefaultTransport.(*http.Transport).Clone()
		if client.cpcfg.ProxyFunc != nil {
			transport.Proxy = client.cpcfg.ProxyFunc
		}
		client.postgresqlClient = &sPostgresqlClient{
			httpClient: &http.Client{
				Transport:     transport,
				Timeout:       30 * time.Second,
				CheckRedirect: func(req *http.Request, via []*http.Request) error { return http.ErrUseLastResponse },
			},
		}
	}
	return client.postgresqlClient
}

func (client *SVolcEngineClient) postgresqlRequest(region, action string, params interface{}) (jsonutils.JSONObject, error) {
	if !postgresqlRegionPattern.MatchString(region) {
		return nil, errors.Errorf("invalid PostgreSQL region")
	}
	switch action {
	case "DescribeRegions", "DescribeDBInstances", "DescribeDBInstanceDetail":
	default:
		return nil, errors.Wrapf(cloudprovider.ErrNotSupported, "PostgreSQL action %s", action)
	}
	query := url.Values{"Action": {action}, "Version": {postgresqlAPIVersion}}
	endpoint := fmt.Sprintf("https://rds-postgresql.%s.volcengineapi.com/?%s", region, query.Encode())
	req, err := http.NewRequest(http.MethodPost, endpoint, strings.NewReader(jsonutils.Marshal(params).String()))
	if err != nil {
		return nil, errors.Wrap(err, "NewRequest")
	}
	req.Header.Set("Content-Type", "application/json")
	cred := client.getDefaultCredential(region, "rds_postgresql")
	resp, err := client.getPostgresqlClient().httpClient.Do(cred.Sign(req))
	if err != nil {
		return nil, errors.Wrapf(err, "PostgreSQL %s", action)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(io.LimitReader(resp.Body, 8*1024*1024))
	if err != nil {
		return nil, errors.Wrap(err, "Read response")
	}
	obj, err := jsonutils.Parse(body)
	if err != nil {
		return nil, errors.Wrapf(err, "invalid PostgreSQL %s response (HTTP %d)", action, resp.StatusCode)
	}
	if obj.Contains("ResponseMetadata", "Error") {
		ve := &sVolcError{StatusCode: resp.StatusCode}
		if err := obj.Unmarshal(ve, "ResponseMetadata"); err != nil {
			return nil, errors.Wrap(err, "Decode PostgreSQL error")
		}
		return nil, ve
	}
	if resp.StatusCode != http.StatusOK {
		return nil, errors.Errorf("PostgreSQL %s returned HTTP %d", action, resp.StatusCode)
	}
	result, err := obj.Get("Result")
	if err != nil {
		return nil, errors.Wrapf(err, "missing PostgreSQL %s Result", action)
	}
	return result, nil
}

func (client *SVolcEngineClient) isPostgresqlRegion(region string) (bool, error) {
	pg := client.getPostgresqlClient()
	pg.regionLock.Lock()
	defer pg.regionLock.Unlock()
	if pg.regions == nil {
		body, err := client.postgresqlRequest(VOLCENGINE_DEFAULT_REGION, "DescribeRegions", map[string]string{})
		if err != nil {
			return false, err
		}
		var regions []struct{ RegionId string }
		if !body.Contains("Regions") {
			return false, errors.Errorf("PostgreSQL DescribeRegions response has no Regions")
		}
		if err := body.Unmarshal(&regions, "Regions"); err != nil {
			return false, errors.Wrap(err, "Decode PostgreSQL regions")
		}
		if len(regions) == 0 {
			return false, errors.Errorf("PostgreSQL DescribeRegions returned an empty region catalog")
		}
		catalog := make(map[string]bool, len(regions))
		for _, r := range regions {
			if r.RegionId == "" {
				return false, errors.Errorf("PostgreSQL region has no RegionId")
			}
			catalog[r.RegionId] = true
		}
		pg.regions = catalog
	}
	return pg.regions[region], nil
}

// GetIDBInstances currently inventories PostgreSQL only. Read errors and
// inconsistent pagination must never be returned as a successful empty set:
// the caller reconciles this set with its existing inventory.
func (region *SRegion) GetIDBInstances() ([]cloudprovider.ICloudDBInstance, error) {
	supported, err := region.client.isPostgresqlRegion(region.RegionId)
	if err != nil {
		return nil, err
	}
	if !supported {
		return nil, errors.Wrapf(cloudprovider.ErrNotSupported, "PostgreSQL in %s", region.RegionId)
	}
	ids := make([]string, 0)
	seen := make(map[string]bool)
	for page := 1; ; page++ {
		body, err := region.client.postgresqlRequest(region.RegionId, "DescribeDBInstances", map[string]int{"PageNumber": page, "PageSize": 100})
		if err != nil {
			return nil, err
		}
		if !body.Contains("Total") || !body.Contains("Instances") {
			return nil, errors.Errorf("PostgreSQL DescribeDBInstances response has no Total or Instances")
		}
		if _, err := body.GetArray("Instances"); err != nil {
			return nil, errors.Wrap(err, "PostgreSQL Instances is not an array")
		}
		var result struct {
			Total     int
			Instances []struct{ InstanceId string }
		}
		if err := body.Unmarshal(&result); err != nil {
			return nil, errors.Wrap(err, "Decode PostgreSQL instances")
		}
		if result.Total < 0 {
			return nil, errors.Errorf("PostgreSQL DescribeDBInstances returned a negative Total")
		}
		for _, instance := range result.Instances {
			if instance.InstanceId == "" || seen[instance.InstanceId] {
				return nil, errors.Errorf("PostgreSQL pagination returned an empty or repeated instance ID")
			}
			seen[instance.InstanceId] = true
			ids = append(ids, instance.InstanceId)
		}
		if len(ids) == result.Total {
			break
		}
		if len(ids) > result.Total || len(result.Instances) == 0 {
			return nil, errors.Errorf("incomplete PostgreSQL pagination: received %d of %d instances", len(ids), result.Total)
		}
	}
	instances := make([]cloudprovider.ICloudDBInstance, 0, len(ids))
	for _, id := range ids {
		instance, err := region.GetPostgresqlInstance(id)
		if err != nil {
			return nil, err
		}
		instances = append(instances, instance)
	}
	return instances, nil
}

func (region *SRegion) GetIDBInstanceById(id string) (cloudprovider.ICloudDBInstance, error) {
	return region.GetPostgresqlInstance(id)
}

func (region *SRegion) GetPostgresqlInstance(id string) (*SPostgresqlInstance, error) {
	if id == "" {
		return nil, errors.Wrap(cloudprovider.ErrNotFound, "empty PostgreSQL instance ID")
	}
	body, err := region.client.postgresqlRequest(region.RegionId, "DescribeDBInstanceDetail", map[string]string{"InstanceId": id})
	if err != nil {
		return nil, err
	}
	var result struct {
		BasicInfo    sPostgresqlBasicInfo
		ChargeDetail sPostgresqlCharge
		Endpoints    []sPostgresqlEndpoint
		Nodes        []sPostgresqlNode
	}
	if err := body.Unmarshal(&result); err != nil {
		return nil, errors.Wrap(err, "Decode PostgreSQL instance detail")
	}
	if result.BasicInfo.InstanceId != id || result.BasicInfo.RegionId != region.RegionId {
		return nil, errors.Errorf("PostgreSQL detail identity or region does not match the request")
	}
	instance := &SPostgresqlInstance{
		region:               region,
		sPostgresqlBasicInfo: result.BasicInfo,
		ChargeDetail:         result.ChargeDetail,
		Endpoints:            result.Endpoints,
		Nodes:                result.Nodes,
	}
	instance.Tags = result.BasicInfo.InstanceTag
	return instance, nil
}

type sPostgresqlBasicInfo struct {
	InstanceId, InstanceName, InstanceStatus string
	DBEngineVersion, InstanceType, NodeSpec  string
	ProjectName, RegionId, ZoneId            string
	VpcID, SubnetId, StorageType             string
	CreateTime                               string
	VCPU, Memory, StorageSpace               int
	InstanceTag                              []multicloud.STag
}

type sPostgresqlCharge struct {
	ChargeType    string
	AutoRenew     bool
	ChargeEndTime string
}

type sPostgresqlAddress struct {
	Domain, IPAddress, Port, NetworkType, SubnetId string
}

type sPostgresqlEndpoint struct {
	EndpointType, ReadWriteMode string
	Address                     []sPostgresqlAddress
}

type sPostgresqlNode struct {
	NodeType, ZoneId string
}

type SPostgresqlInstance struct {
	multicloud.SDBInstanceBase
	VolcEngineTags
	sPostgresqlBasicInfo

	region       *SRegion
	ChargeDetail sPostgresqlCharge
	Endpoints    []sPostgresqlEndpoint
	Nodes        []sPostgresqlNode
}

var _ cloudprovider.ICloudDBInstance = (*SPostgresqlInstance)(nil)

func (instance *SPostgresqlInstance) GetId() string       { return instance.InstanceId }
func (instance *SPostgresqlInstance) GetGlobalId() string { return instance.InstanceId }
func (instance *SPostgresqlInstance) GetName() string {
	if instance.InstanceName != "" {
		return instance.InstanceName
	}
	return instance.InstanceId
}
func (instance *SPostgresqlInstance) GetProjectId() string { return instance.ProjectName }
func (instance *SPostgresqlInstance) GetEngine() string    { return api.DBINSTANCE_TYPE_POSTGRESQL }
func (instance *SPostgresqlInstance) GetEngineVersion() string {
	return strings.TrimPrefix(instance.DBEngineVersion, "PostgreSQL_")
}
func (instance *SPostgresqlInstance) GetInstanceType() string { return instance.NodeSpec }
func (instance *SPostgresqlInstance) GetVcpuCount() int       { return instance.VCPU }
func (instance *SPostgresqlInstance) GetVmemSizeMB() int      { return instance.Memory * 1024 }
func (instance *SPostgresqlInstance) GetDiskSizeGB() int      { return instance.StorageSpace }
func (instance *SPostgresqlInstance) GetStorageType() string  { return instance.StorageType }
func (instance *SPostgresqlInstance) GetCategory() string     { return instance.InstanceType }
func (instance *SPostgresqlInstance) GetCreatedAt() time.Time {
	created, _ := time.Parse(time.RFC3339Nano, instance.CreateTime)
	return created
}
func (instance *SPostgresqlInstance) GetExpiredAt() time.Time {
	expired, _ := time.Parse(time.RFC3339Nano, instance.ChargeDetail.ChargeEndTime)
	return expired
}
func (instance *SPostgresqlInstance) IsAutoRenew() bool       { return instance.ChargeDetail.AutoRenew }
func (instance *SPostgresqlInstance) GetIVpcId() string       { return instance.VpcID }
func (instance *SPostgresqlInstance) GetMaintainTime() string { return "" }

func (instance *SPostgresqlInstance) GetBillingType() string {
	switch strings.ToLower(instance.ChargeDetail.ChargeType) {
	case "prepaid":
		return billing.BILLING_TYPE_PREPAID
	case "postpaid":
		return billing.BILLING_TYPE_POSTPAID
	default:
		return ""
	}
}

func (instance *SPostgresqlInstance) GetStatus() string {
	switch instance.InstanceStatus {
	case "Running":
		return api.DBINSTANCE_RUNNING
	case "Creating":
		return api.DBINSTANCE_DEPLOYING
	case "Restarting":
		return api.DBINSTANCE_REBOOTING
	case "Restoring":
		return api.DBINSTANCE_RESTORING
	case "Updating", "Upgrading":
		return api.DBINSTANCE_UPGRADING
	case "MasterChanging":
		return api.DBINSTANCE_MIGRATING
	case "Deleting", "Deleted", "Released":
		return api.DBINSTANCE_DELETING
	case "Recycled":
		return api.DBINSTANCE_ISOLATE
	case "CreateError":
		return api.DBINSTANCE_CREATE_FAILED
	default:
		return api.DBINSTANCE_UNKNOWN
	}
}

func (instance *SPostgresqlInstance) Refresh() error {
	refreshed, err := instance.region.GetPostgresqlInstance(instance.InstanceId)
	if err != nil {
		return err
	}
	*instance = *refreshed
	return nil
}

func (instance *SPostgresqlInstance) zones() []string {
	zones := []string{}
	if instance.ZoneId != "" {
		zones = append(zones, instance.ZoneId)
	}
	for _, node := range instance.Nodes {
		if node.ZoneId == "" {
			continue
		}
		found := false
		for _, zone := range zones {
			found = found || zone == node.ZoneId
		}
		if !found {
			zones = append(zones, node.ZoneId)
		}
	}
	return zones
}

func (instance *SPostgresqlInstance) zone(index int) string {
	zones := instance.zones()
	if index >= len(zones) {
		return ""
	}
	return zones[index]
}

func (instance *SPostgresqlInstance) GetZone1Id() string { return instance.zone(0) }
func (instance *SPostgresqlInstance) GetZone2Id() string { return instance.zone(1) }
func (instance *SPostgresqlInstance) GetZone3Id() string { return instance.zone(2) }

func (instance *SPostgresqlInstance) addresses() []sPostgresqlAddress {
	for _, endpoint := range instance.Endpoints {
		if endpoint.ReadWriteMode == "ReadWrite" || endpoint.EndpointType == "Primary" {
			return endpoint.Address
		}
	}
	if len(instance.Endpoints) > 0 {
		return instance.Endpoints[0].Address
	}
	return nil
}

func (instance *SPostgresqlInstance) GetPort() int {
	for _, address := range instance.addresses() {
		if address.NetworkType == "Private" {
			port, _ := strconv.Atoi(address.Port)
			return port
		}
	}
	return 0
}

func (instance *SPostgresqlInstance) connection(networkType string) string {
	for _, address := range instance.addresses() {
		if address.NetworkType == networkType {
			if address.Domain != "" {
				return address.Domain
			}
			return address.IPAddress
		}
	}
	return ""
}

func (instance *SPostgresqlInstance) GetConnectionStr() string { return instance.connection("Public") }
func (instance *SPostgresqlInstance) GetInternalConnectionStr() string {
	return instance.connection("Private")
}

func (instance *SPostgresqlInstance) GetDBNetworks() ([]cloudprovider.SDBInstanceNetwork, error) {
	networks := []cloudprovider.SDBInstanceNetwork{}
	for _, address := range instance.addresses() {
		if address.NetworkType != "Private" || net.ParseIP(address.IPAddress) == nil {
			continue
		}
		subnet := address.SubnetId
		if subnet == "" {
			subnet = instance.SubnetId
		}
		if subnet != "" {
			networks = append(networks, cloudprovider.SDBInstanceNetwork{IP: address.IPAddress, NetworkId: subnet})
		}
	}
	return networks, nil
}
