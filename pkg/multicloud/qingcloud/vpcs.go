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

package qingcloud

import (
	"time"

	api "yunion.io/x/cloudmux/pkg/apis/compute"
	"yunion.io/x/cloudmux/pkg/cloudprovider"
	"yunion.io/x/cloudmux/pkg/multicloud"
	"yunion.io/x/jsonutils"
	"yunion.io/x/pkg/errors"
)

type SVpcResp struct {
	Action     string `json:"action"`
	RouterSet  []SVpc `json:"router_set"`
	HasShare   bool   `json:"has_share"`
	TotalCount int    `json:"total_count"`
	RetCode    int    `json:"ret_code"`
}
type SecurityGroups struct {
	GroupID   string `json:"group_id"`
	GroupName string `json:"group_name"`
}

type SVpc struct {
	multicloud.SVpc
	SQingTag

	region *SRegion

	iwires    []cloudprovider.ICloudWire
	secgroups []cloudprovider.ICloudSecurityGroup

	Status              string           `json:"status"`
	BaseVxnet           string           `json:"base_vxnet"`
	IsApplied           int              `json:"is_applied"`
	Features            int              `json:"features"`
	VpcNetwork          string           `json:"vpc_network"`
	ConsoleID           string           `json:"console_id"`
	CreateTime          time.Time        `json:"create_time"`
	AlarmStatus         string           `json:"alarm_status"`
	PrivateIP           string           `json:"private_ip"`
	ResourceProjectInfo []interface{}    `json:"resource_project_info"`
	Owner               string           `json:"owner"`
	PlaceGroupID        string           `json:"place_group_id"`
	SecurityGroups      []SecurityGroups `json:"security_groups"`
	L3Vni               int              `json:"l3vni"`
	SubCode             int              `json:"sub_code"`
	SecurityGroupID     string           `json:"security_group_id"`
	Source              string           `json:"source"`
	Memory              int              `json:"memory"`
	StatusTime          time.Time        `json:"status_time"`
	RouterID            string           `json:"router_id"`
	Description         interface{}      `json:"description"`
	Tags                []interface{}    `json:"tags"`
	TransitionStatus    string           `json:"transition_status"`
	IsDefault           int              `json:"is_default"`
	Controller          string           `json:"controller"`
	VpcID               string           `json:"vpc_id"`
	VpcIpv6Network      string           `json:"vpc_ipv6_network"`
	Eip                 Eip              `json:"eip"`
	Hypervisor          string           `json:"hypervisor"`
	InstanceID          string           `json:"instance_id"`
	RootUserID          string           `json:"root_user_id"`
	DNSAliases          []interface{}    `json:"dns_aliases"`
	Mode                int              `json:"mode"`
	RouterType          int              `json:"router_type"`
	RouterName          string           `json:"router_name"`
	CPU                 int              `json:"cpu"`
}

func (region *SRegion) GetVpcs() ([]SVpc, error) {
	return region.getVpcs()
}

func (region *SRegion) getVpcs(id ...string) ([]SVpc, error) {
	param := map[string]string{}
	if len(id) > 0 {
		param["VpcId.0"] = id[0]
	}
	resp, err := region.client.request("vpc", "DescribeRouters", region.Region, param)
	if err != nil {
		return nil, errors.Wrap(err, "list vpcs")
	}
	vpcs := []SVpc{}
	err = resp.Unmarshal(&vpcs, "VpcSet")
	if err != nil {
		return nil, errors.Wrap(err, "unmarshal instances")
	}
	return vpcs, nil
}

func (vpc *SVpc) GetId() string {
	return vpc.VpcID
}

func (vpc *SVpc) GetName() string {
	if len(vpc.RouterName) > 0 {
		return vpc.RouterName
	}
	return vpc.VpcID
}

func (vpc *SVpc) GetGlobalId() string {
	return vpc.VpcID
}

func (vpc *SVpc) GetStatus() string {
	return api.VPC_STATUS_AVAILABLE
}

func (vpc *SVpc) Refresh() error {
	new, err := vpc.region.getVpcs(vpc.GetId())
	if err != nil {
		return err
	}
	return jsonutils.Update(vpc, new)
}

func (vpc *SVpc) IsEmulated() bool {
	return false
}

func (vpc *SVpc) GetRegion() cloudprovider.ICloudRegion {
	return vpc.region
}

func (vpc *SVpc) GetIsDefault() bool {
	// 华为云没有default vpc.
	return false
}

func (vpc *SVpc) GetCidrBlock() string {
	return vpc.VpcNetwork
}

func (vpc *SVpc) GetIWires() ([]cloudprovider.ICloudWire, error) {
	return nil, cloudprovider.ErrNotImplemented
}

func (vpc *SVpc) GetISecurityGroups() ([]cloudprovider.ICloudSecurityGroup, error) {
	return nil, cloudprovider.ErrNotImplemented
}

func (vpc *SVpc) GetIRouteTables() ([]cloudprovider.ICloudRouteTable, error) {
	return nil, cloudprovider.ErrNotImplemented
}

func (vpc *SVpc) GetIRouteTableById(routeTableId string) (cloudprovider.ICloudRouteTable, error) {
	return nil, cloudprovider.ErrNotImplemented
}

func (vpc *SVpc) Delete() error {
	// todo: 确定删除VPC的逻辑
	return cloudprovider.ErrNotImplemented
}

func (vpc *SVpc) GetIWireById(wireId string) (cloudprovider.ICloudWire, error) {
	return nil, cloudprovider.ErrNotImplemented
}

func (vpc *SVpc) GetINatGateways() ([]cloudprovider.ICloudNatGateway, error) {
	return nil, cloudprovider.ErrNotImplemented

}

func (vpc *SVpc) GetICloudVpcPeeringConnections() ([]cloudprovider.ICloudVpcPeeringConnection, error) {
	return nil, cloudprovider.ErrNotImplemented
}

func (vpc *SVpc) GetICloudAccepterVpcPeeringConnections() ([]cloudprovider.ICloudVpcPeeringConnection, error) {
	return nil, cloudprovider.ErrNotImplemented
}

func (vpc *SVpc) GetICloudVpcPeeringConnectionById(id string) (cloudprovider.ICloudVpcPeeringConnection, error) {
	return nil, cloudprovider.ErrNotImplemented
}

func (vpc *SVpc) CreateICloudVpcPeeringConnection(opts *cloudprovider.VpcPeeringConnectionCreateOptions) (cloudprovider.ICloudVpcPeeringConnection, error) {
	return nil, cloudprovider.ErrNotImplemented
}
func (vpc *SVpc) AcceptICloudVpcPeeringConnection(id string) error {
	return cloudprovider.ErrNotImplemented
}

func (vpc *SVpc) GetAuthorityOwnerId() string {
	return ""
}

func (vpc *SRegion) DeleteVpc(vpcId string) error {
	return cloudprovider.ErrNotImplemented
}
