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

package baidu

import (
	api "yunion.io/x/cloudmux/pkg/apis/compute"
	"yunion.io/x/cloudmux/pkg/cloudprovider"
	"yunion.io/x/cloudmux/pkg/multicloud"
	"yunion.io/x/jsonutils"
	"yunion.io/x/pkg/errors"
	"yunion.io/x/pkg/util/httputils"
)

type SVpcResp struct {
	Vpcs        []SVpc `json:"vpcs"`
	IsTruncated bool   `json:"isTruncated"`
}

type SVpc struct {
	multicloud.SVpc
	SBaiduTag

	region *SRegion

	iwires    []cloudprovider.ICloudWire
	secgroups []cloudprovider.ICloudSecurityGroup

	IsDefault bool   `json:"IsDefault"`
	VpcID     string `json:"vpcId"`
	Name      string `json:"name"`

	CreateTime string `json:"CreateTime"`
	Cidr       string `json:"Cidr"`
}

func (region *SRegion) GetVpcs() ([]SVpc, error) {
	return region.getVpcs()
}

func (region *SRegion) getVpcs(id ...string) ([]SVpc, error) {
	param := map[string]string{}
	if len(id) > 0 {
		param["VpcId.0"] = id[0]
	}
	resp, err := region.client.request(httputils.GET, "bcc", region.Region, "v1/vpc", nil)
	if err != nil {
		return nil, errors.Wrap(err, "list instance")
	}
	vpcs := []SVpc{}
	err = resp.Unmarshal(&vpcs, "vpcs")
	if err != nil {
		return nil, errors.Wrap(err, "unmarshal instances")
	}
	return vpcs, nil
}

func (vpc *SVpc) GetId() string {
	return vpc.VpcID
}

func (vpc *SVpc) GetName() string {
	if len(vpc.Name) > 0 {
		return vpc.Name
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
	return vpc.Cidr
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
