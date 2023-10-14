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

package ksyun

import (
	api "yunion.io/x/cloudmux/pkg/apis/compute"
	"yunion.io/x/cloudmux/pkg/cloudprovider"
	"yunion.io/x/cloudmux/pkg/multicloud"
	"yunion.io/x/jsonutils"
	"yunion.io/x/log"
	"yunion.io/x/pkg/errors"
	"yunion.io/x/pkg/util/netutils"
	"yunion.io/x/pkg/util/rbacscope"
)

type SNetworkResp struct {
	SubnetSet []SNetwork `json:"SubnetSet"`
	RequestID string     `json:"RequestId"`
	NextToken string     `json:"NextToken"`
}

type SNetwork struct {
	multicloud.SResourceBase
	SKsTag
	wire *SWire

	RouteTableID                string                        `json:"RouteTableId"`
	NetworkACLID                string                        `json:"NetworkAclId"`
	NatID                       string                        `json:"NatId"`
	CreateTime                  string                        `json:"CreateTime"`
	DhcpIPTo                    string                        `json:"DhcpIpTo"`
	DNS1                        string                        `json:"Dns1"`
	CidrBlock                   string                        `json:"CidrBlock"`
	DNS2                        string                        `json:"Dns2"`
	ProvidedIpv6CidrBlock       bool                          `json:"ProvidedIpv6CidrBlock"`
	SubnetID                    string                        `json:"SubnetId"`
	SubnetType                  string                        `json:"SubnetType"`
	SubnetName                  string                        `json:"SubnetName"`
	VpcID                       string                        `json:"VpcId"`
	GatewayIP                   string                        `json:"GatewayIp"`
	AvailabilityZoneName        string                        `json:"AvailabilityZoneName"`
	DhcpIPFrom                  string                        `json:"DhcpIpFrom"`
	Ipv6CidrBlockAssociationSet []Ipv6CidrBlockAssociationSet `json:"Ipv6CidrBlockAssociationSet"`
	AvailableIPNumber           int                           `json:"AvailableIpNumber"`
	SecondaryCidrID             string                        `json:"SecondaryCidrId"`
}

func (region *SRegion) GetNetworks(vpcId ...string) ([]SNetwork, error) {
	return region.getNetworks(vpcId...)
}

func (region *SRegion) getNetwork(networkId string) ([]SNetwork, error) {
	networks := []SNetwork{}
	param := map[string]string{}
	if len(networkId) > 0 {
		param["SubnetId.1"] = networkId
	}
	resp, err := region.client.request("vpc", region.Region, "DescribeSubnets", "2016-03-04", param)
	if err != nil {
		return nil, errors.Wrap(err, "list networks")
	}
	res := SNetworkResp{}
	err = resp.Unmarshal(&res)
	if err != nil {
		return nil, errors.Wrap(err, "unmarshal instances")
	}
	networks = append(networks, res.SubnetSet...)

	return networks, nil
}

func (region *SRegion) getNetworks(vpcId ...string) ([]SNetwork, error) {
	networks := []SNetwork{}
	param := map[string]string{}
	if len(vpcId) > 0 && len(vpcId[0]) > 0 {
		param["Filter.1.Name"] = "vpc-id"
		param["Filter.1.Value.1"] = vpcId[0]
	}
	resp, err := region.client.request("vpc", region.Region, "DescribeSubnets", "2016-03-04", param)
	if err != nil {
		return nil, errors.Wrap(err, "list networks")
	}
	res := SNetworkResp{}
	err = resp.Unmarshal(&res)
	if err != nil {
		return nil, errors.Wrap(err, "unmarshal instances")
	}
	networks = append(networks, res.SubnetSet...)
	return networks, nil
}

func (net *SNetwork) GetId() string {
	return net.SubnetID
}

func (net *SNetwork) GetName() string {
	if len(net.SubnetName) == 0 {
		return net.SubnetID
	}

	return net.SubnetName
}

func (net *SNetwork) GetGlobalId() string {
	return net.SubnetID
}

// https://support.huaweicloud.com/api-vpc/zh-cn_topic_0020090591.html
func (net *SNetwork) GetStatus() string {
	return api.NETWORK_STATUS_AVAILABLE
}

func (net *SNetwork) Refresh() error {
	log.Debugf("network refresh %s", net.GetId())
	new, err := net.wire.region.getNetworks(net.GetId())
	if err != nil {
		return err
	}
	return jsonutils.Update(net, new)
}

func (net *SNetwork) IsEmulated() bool {
	return false
}

func (net *SNetwork) GetIWire() cloudprovider.ICloudWire {
	return net.wire
}

func (net *SNetwork) GetIpStart() string {
	pref, _ := netutils.NewIPV4Prefix(net.CidrBlock)
	startIp := pref.Address.NetAddr(pref.MaskLen) // 0
	startIp = startIp.StepUp()                    // 1
	startIp = startIp.StepUp()                    // 2
	return startIp.String()
}

func (net *SNetwork) GetIpEnd() string {
	pref, _ := netutils.NewIPV4Prefix(net.CidrBlock)
	endIp := pref.Address.BroadcastAddr(pref.MaskLen) // 255
	endIp = endIp.StepDown()                          // 254
	endIp = endIp.StepDown()                          // 253
	endIp = endIp.StepDown()                          // 252
	return endIp.String()
}

func (net *SNetwork) GetIpMask() int8 {
	pref, _ := netutils.NewIPV4Prefix(net.CidrBlock)
	return pref.MaskLen
}

func (net *SNetwork) GetGateway() string {
	pref, _ := netutils.NewIPV4Prefix(net.CidrBlock)
	startIp := pref.Address.NetAddr(pref.MaskLen) // 0
	startIp = startIp.StepUp()                    // 1
	return startIp.String()
}

func (net *SNetwork) GetServerType() string {
	return api.NETWORK_TYPE_GUEST
}

func (net *SNetwork) GetIsPublic() bool {
	return true
}

func (net *SNetwork) GetPublicScope() rbacscope.TRbacScope {
	return rbacscope.ScopeDomain
}

func (net *SNetwork) Delete() error {
	return cloudprovider.ErrNotImplemented
}

func (net *SNetwork) GetAllocTimeoutSeconds() int {
	return 120 // 2 minutes
}

func (net *SNetwork) GetProjectId() string {
	return ""
}

func (net *SNetwork) GetDescription() string {
	return ""
}
