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
	"fmt"

	api "yunion.io/x/cloudmux/pkg/apis/compute"
	"yunion.io/x/cloudmux/pkg/cloudprovider"
	"yunion.io/x/cloudmux/pkg/multicloud"
	"yunion.io/x/jsonutils"
	"yunion.io/x/pkg/errors"
)

type SecondaryCidrSet struct {
	Type            string `json:"Type"`
	Cidr            string `json:"Cidr"`
	SecondaryCidrID string `json:"SecondaryCidrId"`
}
type Ipv6CidrBlockAssociationSet struct {
	Ipv6CidrBlock string `json:"Ipv6CidrBlock"`
}

type SVpc struct {
	multicloud.SVpc
	SKsTag

	region *SRegion
	zone   *SZone

	iwires    []cloudprovider.ICloudWire
	secgroups []cloudprovider.ICloudSecurityGroup

	IsDefault                   bool                          `json:"IsDefault"`
	SecondaryCidrSet            []SecondaryCidrSet            `json:"SecondaryCidrSet"`
	VpcID                       string                        `json:"VpcId"`
	CreateTime                  string                        `json:"CreateTime"`
	CidrBlock                   string                        `json:"CidrBlock"`
	Ipv6CidrBlockAssociationSet []Ipv6CidrBlockAssociationSet `json:"Ipv6CidrBlockAssociationSet"`
	VpcName                     string                        `json:"VpcName"`
	ProvidedIpv6CidrBlock       bool                          `json:"ProvidedIpv6CidrBlock"`
}

func (region *SRegion) GetVpcs() ([]SVpc, error) {
	return region.getVpcs()
}

func (region *SRegion) getVpcs(id ...string) ([]SVpc, error) {
	param := map[string]string{}
	if len(id) > 0 {
		param["VpcId.0"] = id[0]
	}
	resp, err := region.client.request("vpc", region.Region, "DescribeVpcs", "2016-03-04", param)
	if err != nil {
		return nil, errors.Wrap(err, "list instance")
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
	if len(vpc.VpcName) > 0 {
		return vpc.VpcName
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
	return vpc.CidrBlock
}

func (vpc *SVpc) GetIWires() ([]cloudprovider.ICloudWire, error) {
	zones, err := vpc.region.GetZones()
	if err != nil {
		return nil, errors.Wrap(err, "GetZones")
	}
	wires := []cloudprovider.ICloudWire{}
	for i := 0; i < len(zones); i++ {
		wire := SWire{
			vpc:    vpc,
			zone:   &zones[i],
			region: vpc.region,
			WireId: fmt.Sprintf("%s-%s", vpc.GetId(), zones[i].GetId()),
			Id:     fmt.Sprintf("%s-%s", vpc.GetId(), zones[i].GetId()),
		}
		wires = append(wires, &wire)
	}

	return wires, nil
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
