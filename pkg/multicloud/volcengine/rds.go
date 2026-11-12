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
	"strings"
	api "yunion.io/x/cloudmux/pkg/apis/compute"
	"yunion.io/x/cloudmux/pkg/cloudprovider"
	"yunion.io/x/cloudmux/pkg/multicloud"
	"yunion.io/x/pkg/errors"
)

// Reconcile all supported engines atomically. Never pass a partial engine list
// to Cloudpods: that would mark the other engines' existing records as missing.
func (r *SRegion) GetIDBInstances() ([]cloudprovider.ICloudDBInstance, error) {
	ret := []cloudprovider.ICloudDBInstance{}
	supported := false
	pg, e := r.getPostgresqlInstances()
	if e == nil {
		ret = append(ret, pg...)
		supported = true
	} else if errors.Cause(e) != cloudprovider.ErrNotSupported {
		return nil, e
	}
	for _, service := range []string{"rds-mysql", "rds-mssql"} {
		items, e := r.getReadonlyRDSInstances(service)
		if e != nil {
			if errors.Cause(e) == cloudprovider.ErrNotSupported {
				continue
			}
			return nil, e
		}
		supported = true
		ret = append(ret, items...)
	}
	if !supported {
		return nil, cloudprovider.ErrNotSupported
	}
	return ret, nil
}
func (r *SRegion) getReadonlyRDSInstances(service string) ([]cloudprovider.ICloudDBInstance, error) {
	catalog, e := r.client.inventoryRequest("cn-beijing", service, "DescribeRegions", map[string]interface{}{})
	if e != nil {
		return nil, e
	}
	regions, e := catalog.GetArray("Regions")
	if e != nil || len(regions) == 0 {
		return nil, errors.Errorf("invalid %s region catalog", service)
	}
	supported := false
	for _, region := range regions {
		id, e := region.GetString("RegionId")
		if e != nil || id == "" {
			return nil, errors.Errorf("invalid RDS region")
		}
		supported = supported || id == r.RegionId
	}
	if !supported {
		return nil, cloudprovider.ErrNotSupported
	}
	array := "Instances"
	if service == "rds-mssql" {
		array = "InstancesInfo"
	}
	rows, e := r.inventoryPages(service, "DescribeDBInstances", array, "Total", "InstanceId", nil)
	if e != nil {
		return nil, e
	}
	ret := make([]cloudprovider.ICloudDBInstance, 0, len(rows))
	for _, row := range rows {
		id, _ := row.GetString("InstanceId")
		i, e := r.getReadonlyRDSInstance(service, id)
		if e != nil {
			return nil, e
		}
		ret = append(ret, i)
	}
	return ret, nil
}

type SReadonlyRDSInstance struct {
	SPostgresqlInstance
	service          string
	masterInstanceId string
}

var _ cloudprovider.ICloudDBInstance = (*SReadonlyRDSInstance)(nil)

func (r *SRegion) getReadonlyRDSInstance(service, id string) (*SReadonlyRDSInstance, error) {
	if id == "" {
		return nil, cloudprovider.ErrNotFound
	}
	data, e := r.client.inventoryRequest(r.RegionId, service, "DescribeDBInstanceDetail", map[string]interface{}{"InstanceId": id})
	if e != nil {
		return nil, e
	}
	var detail struct {
		BasicInfo struct {
			Tags                                []multicloud.STag
			MasterInstanceId, PrimaryInstanceId string
		}
		ChargeDetail sPostgresqlCharge
		Endpoints    []struct {
			EndpointType, ReadWriteMode string
			Address, Addresses          []sPostgresqlAddress
		}
		ConnectionInfo        []sPostgresqlEndpoint
		Nodes, NodeDetailInfo []sPostgresqlNode
	}
	if e := data.Unmarshal(&detail); e != nil {
		return nil, e
	}
	var basic sPostgresqlBasicInfo
	if e := data.Unmarshal(&basic, "BasicInfo"); e != nil {
		return nil, e
	}
	if basic.InstanceId != id || basic.RegionId != r.RegionId {
		return nil, errors.Errorf("%s detail identity or region mismatch", service)
	}
	i := &SReadonlyRDSInstance{service: service, SPostgresqlInstance: SPostgresqlInstance{region: r, sPostgresqlBasicInfo: basic, ChargeDetail: detail.ChargeDetail, Nodes: detail.Nodes}}
	i.Tags = detail.BasicInfo.Tags
	i.masterInstanceId = detail.BasicInfo.MasterInstanceId
	if service == "rds-mssql" {
		i.Endpoints = detail.ConnectionInfo
		i.Nodes = detail.NodeDetailInfo
		i.masterInstanceId = detail.BasicInfo.PrimaryInstanceId
	} else {
		for _, ep := range detail.Endpoints {
			addresses := ep.Address
			if len(addresses) == 0 {
				addresses = ep.Addresses
			}
			i.Endpoints = append(i.Endpoints, sPostgresqlEndpoint{EndpointType: ep.EndpointType, ReadWriteMode: ep.ReadWriteMode, Address: addresses})
		}
	}
	return i, nil
}
func (i *SReadonlyRDSInstance) GetEngine() string {
	if i.service == "rds-mssql" {
		return api.DBINSTANCE_TYPE_SQLSERVER
	}
	return api.DBINSTANCE_TYPE_MYSQL
}
func (i *SReadonlyRDSInstance) GetEngineVersion() string {
	if i.service == "rds-mssql" {
		return strings.TrimPrefix(i.DBEngineVersion, "SQLServer_")
	}
	return strings.ReplaceAll(strings.TrimPrefix(i.DBEngineVersion, "MySQL_"), "_", ".")
}
func (i *SReadonlyRDSInstance) GetMasterInstanceId() string { return i.masterInstanceId }
func (i *SReadonlyRDSInstance) Refresh() error {
	r, e := i.region.getReadonlyRDSInstance(i.service, i.InstanceId)
	if e != nil {
		return e
	}
	*i = *r
	return nil
}
