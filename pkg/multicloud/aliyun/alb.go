// Copyright 2019 Yunion
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package aliyun

import (
	"context"
	"fmt"
	"time"

	"yunion.io/x/cloudmux/pkg/apis/compute"
	"yunion.io/x/cloudmux/pkg/cloudprovider"
	"yunion.io/x/cloudmux/pkg/multicloud"
	"yunion.io/x/jsonutils"
	"yunion.io/x/pkg/errors"
)

type SAlb struct {
	multicloud.SLoadbalancerBase
	region *SRegion
	AliyunTags

	AddressIPVersion    string
	LoadBalancerVersion string
	ResourceGroupId     string
	LoadBalancerId      string
	DNSName             string
	ZoneMappings        []struct {
		Status                string
		ZoneId                string
		VSwitchId             string
		LoadBalancerAddresses []struct {
			IntranetAddress         string
			IntranetAddressHcStatus string
			Address                 string
			Ipv4LocalAddresses      []string
		}
	}
	BandwidthCapacity        int
	DeletionProtectionConfig struct {
		Enabled bool
	}
	SysSecurityGroupId        string
	BackToOriginRouteEnabled  bool
	LoadBalancerEdition       string
	LoadBalancerBillingConfig struct {
		PayType string
	}
	CreateTime                  time.Time
	LoadBalancerName            string
	LoadBalancerBussinessStatus string
	VpcId                       string
	RegionId                    string
	AddressAllocatedMode        string
	AddressType                 string
	LoadBalancerStatus          string
}

func (alb *SAlb) GetId() string {
	return alb.LoadBalancerId
}

func (alb *SAlb) GetName() string {
	return alb.LoadBalancerName
}

func (alb *SAlb) GetStatus() string {
	return alb.LoadBalancerStatus
}

func (alb *SAlb) GetGlobalId() string {
	return alb.LoadBalancerId
}

func (alb *SAlb) GetAddress() string {
	return alb.DNSName
}

func (alb *SAlb) GetAddressType() string {
	return alb.AddressType
}

func (alb *SAlb) GetNetworkType() string {
	return "vpc"
}

func (alb *SAlb) GetType() cloudprovider.LoadbalancerType {
	return cloudprovider.LoadbalancerTypeALB
}

func (alb *SAlb) GetNetworkIds() []string {
	if len(alb.ZoneMappings) == 0 {
		alb.Refresh()
	}
	networks := []string{}
	for _, zone := range alb.ZoneMappings {
		networks = append(networks, zone.VSwitchId)
	}
	return networks
}

func (alb *SAlb) GetVpcId() string {
	return alb.VpcId
}

func (alb *SAlb) Refresh() error {
	lb, err := alb.region.GetAlb(alb.LoadBalancerId)
	if err != nil {
		return err
	}
	return jsonutils.Update(alb, lb)
}

func (alb *SAlb) GetZoneId() string {
	if len(alb.ZoneMappings) == 0 {
		alb.Refresh()
	}
	for _, zone := range alb.ZoneMappings {
		return zone.ZoneId
	}
	return ""
}

func (alb *SAlb) GetZone1Id() string {
	if len(alb.ZoneMappings) == 0 {
		alb.Refresh()
	}
	for i, zone := range alb.ZoneMappings {
		if i != 0 {
			return zone.ZoneId
		}
	}
	return ""
}

func (alb *SAlb) GetLoadbalancerSpec() string {
	return alb.LoadBalancerEdition
}

func (alb *SAlb) GetChargeType() string {
	return compute.LB_CHARGE_TYPE_BY_TRAFFIC
}

func (alb *SAlb) GetEgressMbps() int {
	if alb.BandwidthCapacity == 0 {
		alb.Refresh()
	}
	return alb.BandwidthCapacity
}

func (alb *SAlb) GetIEIP() (cloudprovider.ICloudEIP, error) {
	return nil, nil // TODO: implement
}

func (alb *SAlb) Delete(ctx context.Context) error {
	return errors.Wrap(cloudprovider.ErrNotImplemented, "Delete")
}

func (alb *SAlb) Start() error {
	return errors.Wrap(cloudprovider.ErrNotImplemented, "Start")
}

func (alb *SAlb) Stop() error {
	return errors.Wrap(cloudprovider.ErrNotImplemented, "Stop")
}

func (alb *SAlb) GetILoadBalancerListeners() ([]cloudprovider.ICloudLoadbalancerListener, error) {
	return nil, errors.Wrap(cloudprovider.ErrNotImplemented, "GetILoadBalancerListeners")
}

func (alb *SAlb) GetILoadBalancerBackendGroups() ([]cloudprovider.ICloudLoadbalancerBackendGroup, error) {
	return nil, errors.Wrap(cloudprovider.ErrNotImplemented, "GetILoadBalancerBackendGroups")
}

func (alb *SAlb) CreateILoadBalancerBackendGroup(group *cloudprovider.SLoadbalancerBackendGroup) (cloudprovider.ICloudLoadbalancerBackendGroup, error) {
	return nil, errors.Wrap(cloudprovider.ErrNotImplemented, "CreateILoadBalancerBackendGroup")
}

func (alb *SAlb) GetILoadBalancerBackendGroupById(groupId string) (cloudprovider.ICloudLoadbalancerBackendGroup, error) {
	return nil, errors.Wrap(cloudprovider.ErrNotImplemented, "GetILoadBalancerBackendGroupById")
}

func (alb *SAlb) CreateILoadBalancerListener(ctx context.Context, listener *cloudprovider.SLoadbalancerListenerCreateOptions) (cloudprovider.ICloudLoadbalancerListener, error) {
	return nil, errors.Wrap(cloudprovider.ErrNotImplemented, "CreateILoadBalancerListener")
}

func (alb *SAlb) GetILoadBalancerListenerById(listenerId string) (cloudprovider.ICloudLoadbalancerListener, error) {
	return nil, errors.Wrap(cloudprovider.ErrNotImplemented, "GetILoadBalancerListenerById")
}

func (alb *SAlb) GetSecurityGroupIds() ([]string, error) {
	if len(alb.SysSecurityGroupId) == 0 {
		alb.Refresh()
	}
	if len(alb.SysSecurityGroupId) == 0 {
		return []string{}, nil
	}
	return []string{alb.SysSecurityGroupId}, nil
}

func (region *SRegion) GetAlbs(ids []string) ([]SAlb, error) {
	params := map[string]string{}
	for i, id := range ids {
		params[fmt.Sprintf("LoadBalancerIds.%d", i+1)] = id
	}
	albs := []SAlb{}
	for {
		resp, err := region.albRequest("ListLoadBalancers", params)
		if err != nil {
			return nil, err
		}
		part := struct {
			LoadBalancers []SAlb
			NextToken     string
		}{}
		err = resp.Unmarshal(&part)
		if err != nil {
			return nil, err
		}
		albs = append(albs, part.LoadBalancers...)
		if len(part.NextToken) == 0 {
			break
		}
		params["NextToken"] = part.NextToken
	}
	for i := 0; i < len(albs); i++ {
		albs[i].region = region
	}
	return albs, nil
}

func (region *SRegion) GetAlb(id string) (*SAlb, error) {
	resp, err := region.albRequest("GetLoadBalancerAttribute", map[string]string{
		"LoadBalancerId": id,
	})
	if err != nil {
		return nil, err
	}
	ret := &SAlb{}
	err = resp.Unmarshal(ret)
	if err != nil {
		return nil, err
	}
	ret.region = region
	return ret, nil
}
