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
	"strconv"
	"strings"
	"time"
	billing "yunion.io/x/cloudmux/pkg/apis/billing"
	api "yunion.io/x/cloudmux/pkg/apis/compute"
	"yunion.io/x/cloudmux/pkg/cloudprovider"
	"yunion.io/x/cloudmux/pkg/multicloud"
	"yunion.io/x/pkg/errors"
)

type SRedisInstance struct {
	multicloud.SElasticcacheBase
	VolcEngineTags
	region                                                                         *SRegion
	InstanceId, InstanceName, RegionId, Status, EngineVersion, InstanceClass       string
	VpcId, SubnetId, ProjectName, ChargeType, CreateTime, ExpiredTime, VpcAuthMode string
	AutoRenew                                                                      bool
	ShardedCluster, NodeNumber, MaxConnections                                     int
	ZoneIds                                                                        []string
	Capacity                                                                       struct{ Total int }
	VisitAddrs                                                                     []struct{ AddrType, Address, VIP, Port string }
}

var _ cloudprovider.ICloudElasticcache = (*SRedisInstance)(nil)

func (region *SRegion) GetIElasticcaches() ([]cloudprovider.ICloudElasticcache, error) {
	catalog, err := region.client.inventoryRequest("cn-beijing", "redis", "DescribeRegions", map[string]interface{}{})
	if err != nil {
		return nil, err
	}
	regions, err := catalog.GetArray("Regions")
	if err != nil || len(regions) == 0 {
		return nil, errors.Errorf("invalid Redis region catalog")
	}
	supported := false
	for _, r := range regions {
		id, err := r.GetString("RegionId")
		if err != nil || id == "" {
			return nil, errors.Errorf("invalid Redis region")
		}
		supported = supported || id == region.RegionId
	}
	if !supported {
		return nil, cloudprovider.ErrNotSupported
	}
	items, err := region.inventoryPages("redis", "DescribeDBInstances", "Instances", "TotalInstancesNum", "InstanceId", map[string]interface{}{})
	if err != nil {
		return nil, err
	}
	ret := make([]cloudprovider.ICloudElasticcache, 0, len(items))
	for _, item := range items {
		id, _ := item.GetString("InstanceId")
		instance, err := region.GetRedisInstance(id)
		if err != nil {
			return nil, err
		}
		ret = append(ret, instance)
	}
	return ret, nil
}
func (region *SRegion) GetIElasticcacheById(id string) (cloudprovider.ICloudElasticcache, error) {
	return region.GetRedisInstance(id)
}
func (region *SRegion) GetRedisInstance(id string) (*SRedisInstance, error) {
	if id == "" {
		return nil, cloudprovider.ErrNotFound
	}
	data, err := region.client.inventoryRequest(region.RegionId, "redis", "DescribeDBInstanceDetail", map[string]interface{}{"InstanceId": id})
	if err != nil {
		return nil, err
	}
	instance := &SRedisInstance{region: region}
	if err = data.Unmarshal(instance); err != nil {
		return nil, err
	}
	if instance.InstanceId != id || instance.RegionId != region.RegionId {
		return nil, errors.Errorf("Redis detail identity or region mismatch")
	}
	if instance.Capacity.Total <= 0 {
		return nil, errors.Errorf("missing Redis capacity")
	}
	return instance, nil
}
func (i *SRedisInstance) GetId() string       { return i.InstanceId }
func (i *SRedisInstance) GetGlobalId() string { return i.InstanceId }
func (i *SRedisInstance) GetName() string {
	if i.InstanceName != "" {
		return i.InstanceName
	}
	return i.InstanceId
}
func (i *SRedisInstance) GetProjectId() string { return i.ProjectName }
func (i *SRedisInstance) GetStatus() string {
	switch i.Status {
	case "Running":
		return api.ELASTIC_CACHE_STATUS_RUNNING
	default:
		return "unknown"
	}
}
func (i *SRedisInstance) GetInstanceType() string { return i.InstanceClass }

// Capacity.Total is MiB, not bytes or per-replica capacity.
func (i *SRedisInstance) GetCapacityMB() int { return i.Capacity.Total }
func (i *SRedisInstance) GetArchType() string {
	if i.ShardedCluster == 1 {
		return "cluster"
	}
	if i.NodeNumber > 1 {
		return api.ELASTIC_CACHE_ARCH_TYPE_MASTER
	}
	return api.ELASTIC_CACHE_ARCH_TYPE_SINGLE
}
func (i *SRedisInstance) GetNodeType() string {
	if i.NodeNumber > 1 {
		return "double"
	}
	return "single"
}
func (i *SRedisInstance) GetEngine() string        { return "redis" }
func (i *SRedisInstance) GetEngineVersion() string { return i.EngineVersion }
func (i *SRedisInstance) GetVpcId() string         { return i.VpcId }
func (i *SRedisInstance) GetNetworkId() string     { return i.SubnetId }
func (i *SRedisInstance) GetNetworkType() string   { return "vpc" }
func (i *SRedisInstance) GetZoneId() string {
	if len(i.ZoneIds) > 0 {
		return i.region.GetGlobalId() + "/" + i.ZoneIds[0]
	}
	return ""
}
func (i *SRedisInstance) address(kind string) (string, string, int) {
	for _, a := range i.VisitAddrs {
		if a.AddrType == kind {
			p, _ := strconv.Atoi(a.Port)
			return a.Address, a.VIP, p
		}
	}
	return "", "", 0
}
func (i *SRedisInstance) GetPrivateDNS() string        { d, _, _ := i.address("Private"); return d }
func (i *SRedisInstance) GetPrivateIpAddr() string     { _, v, _ := i.address("Private"); return v }
func (i *SRedisInstance) GetPrivateConnectPort() int   { _, _, p := i.address("Private"); return p }
func (i *SRedisInstance) GetPublicDNS() string         { d, _, _ := i.address("Public"); return d }
func (i *SRedisInstance) GetPublicIpAddr() string      { _, v, _ := i.address("Public"); return v }
func (i *SRedisInstance) GetPublicConnectPort() int    { _, _, p := i.address("Public"); return p }
func (i *SRedisInstance) GetConnections() int          { return i.MaxConnections }
func (i *SRedisInstance) GetMaintainStartTime() string { return "" }
func (i *SRedisInstance) GetMaintainEndTime() string   { return "" }
func (i *SRedisInstance) GetAuthMode() string {
	// Volcengine open means password-free access, the inverse of auth_mode.
	switch strings.ToLower(i.VpcAuthMode) {
	case "open":
		return "off"
	case "close":
		return "on"
	default:
		return "unknown"
	}
}
func (i *SRedisInstance) GetCreatedAt() time.Time {
	t, _ := time.Parse(time.RFC3339Nano, i.CreateTime)
	return t
}
func (i *SRedisInstance) GetExpiredAt() time.Time {
	t, _ := time.Parse(time.RFC3339Nano, i.ExpiredTime)
	return t
}
func (i *SRedisInstance) IsAutoRenew() bool { return i.AutoRenew }
func (i *SRedisInstance) GetBillingType() string {
	switch i.ChargeType {
	case "PrePaid":
		return billing.BILLING_TYPE_PREPAID
	case "PostPaid":
		return billing.BILLING_TYPE_POSTPAID
	default:
		return ""
	}
}
func (i *SRedisInstance) Refresh() error {
	r, e := i.region.GetRedisInstance(i.InstanceId)
	if e != nil {
		return e
	}
	*i = *r
	return nil
}
func (i *SRedisInstance) GetSecurityGroupIds() ([]string, error) {
	return nil, cloudprovider.ErrNotSupported
}
