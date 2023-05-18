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

package bingocloud

import (
	"fmt"
	"strconv"
	"strings"

	"yunion.io/x/jsonutils"
	"yunion.io/x/pkg/errors"

	api "yunion.io/x/cloudmux/pkg/apis/compute"
	"yunion.io/x/cloudmux/pkg/cloudprovider"
	"yunion.io/x/cloudmux/pkg/multicloud"
)

type SEip struct {
	multicloud.SEipBase
	BingoTags

	region *SRegion

	AddressId    string
	AddressType  string
	Bandtype     string
	Bandwidth    int
	CanAssociate bool
	InstanceId   string
	Owner        string
	PublicIp     string
	SubnetId     string
	VpcId        string
}

func (self *SEip) GetId() string {
	return self.PublicIp
}

func (self *SEip) GetGlobalId() string {
	return self.PublicIp
}

func (self *SEip) GetName() string {
	return self.PublicIp
}

func (self *SEip) GetIpAddr() string {
	return self.PublicIp
}

func (self *SEip) GetMode() string {
	return api.EIP_MODE_STANDALONE_EIP
}

func (self *SEip) GetINetworkId() string {
	return self.SubnetId
}

func (self *SEip) GetAssociationType() string {
	if len(self.InstanceId) > 0 {
		return api.EIP_ASSOCIATE_TYPE_SERVER
	}
	return ""
}

func (self *SEip) GetAssociationExternalId() string {
	return self.InstanceId
}

func (self *SEip) GetBandwidth() int {
	return self.Bandwidth
}

func (self *SEip) GetInternetChargeType() string {
	return api.EIP_CHARGE_TYPE_BY_BANDWIDTH
}

func (self *SEip) Refresh() error {
	newEip, err := self.region.GetIEipById(self.GetId())
	if err != nil {
		return err
	}
	return jsonutils.Update(self, &newEip)
}

func (self *SEip) Delete() error {
	return cloudprovider.ErrNotImplemented
}

func (self *SEip) Associate(conf *cloudprovider.AssociateConfig) error {
	nics, err := self.region.GetInstanceNics(conf.InstanceId)
	if err != nil {
		return err
	}
	if len(nics) == 0 {
		return errors.Wrapf(cloudprovider.ErrNotFound, "region.GetInstanceNics", conf.InstanceId)
	}

	params := map[string]string{}
	params["PublicIp"] = self.PublicIp
	params["InstanceId"] = conf.InstanceId
	params["NetworkInterfaceId"] = nics[0].NetworkInterfaceId

	_, err = self.region.invoke("AssociateAddress", params)
	return err
}

func (self *SEip) Dissociate() error {
	params := map[string]string{}
	params["PublicIp"] = self.PublicIp

	_, err := self.region.invoke("DisassociateAddress", params)
	//参数错误: 只有用作NAT的弹性IP才能解绑, 错误编码: 106574
	if err != nil && strings.Contains(err.Error(), "106574") {
		return nil
	}

	return err
}

func (self *SEip) ChangeBandwidth(bw int) error {
	return cloudprovider.ErrNotImplemented
}

func (self *SEip) GetProjectId() string {
	return ""
}

func (self *SEip) GetStatus() string {
	return api.EIP_STATUS_READY
}

func (self *SRegion) GetEips(ip, instanceId, nextToken string) ([]SEip, string, error) {
	params := map[string]string{}
	if len(ip) > 0 {
		params["PublicIp.1"] = ip
	}
	if len(nextToken) > 0 {
		params["NextToken"] = nextToken
	}

	idx := 1
	if len(instanceId) > 0 {
		params[fmt.Sprintf("Filter.%d.Name", idx)] = "instance-id"
		params[fmt.Sprintf("Filter.%d.Value.1", idx)] = instanceId
		idx++
	}
	params[fmt.Sprintf("Filter.%d.Name", idx)] = "owner-id"
	params[fmt.Sprintf("Filter.%d.Value.1", idx)] = self.client.user
	idx++

	resp, err := self.invoke("DescribeAddresses", params)
	if err != nil {
		return nil, "", errors.Wrapf(err, "DescribeAddresses")
	}
	ret := struct {
		AddressesSet []SEip
		NextToken    string
	}{}
	_ = resp.Unmarshal(&ret)

	if len(ip) > 0 || len(instanceId) > 0 {
		return ret.AddressesSet, ret.NextToken, nil
	}

	var floatingRet []SEip
	for _, eip := range ret.AddressesSet {
		if strings.Contains(strings.ToLower(eip.AddressType), "virtualip") {
			continue
		}
		if eip.InstanceId == "" && eip.CanAssociate {
			floatingRet = append(floatingRet, eip)
			continue
		}
		if eip.InstanceId != "" {
			nics, err := self.GetInstanceNics(eip.InstanceId)
			if err != nil {
				return nil, "", err
			}
			isSecondaryIp := false
			for i := range nics {
				if eip.PublicIp == nics[i].PrivateIPAddress {
					isSecondaryIp = true
					break
				}
				for j := range nics[i].PrivateIPAddressesSet {
					if nics[i].PrivateIPAddressesSet[j].PrivateIPAddress == eip.PublicIp && nics[i].PrivateIPAddressesSet[j].PrivateIPAddress == nics[i].PrivateIPAddressesSet[j].Association.PublicIp {
						isSecondaryIp = true
						break
					}
				}
			}
			if !isSecondaryIp {
				floatingRet = append(floatingRet, eip)
			}
		}
	}

	return floatingRet, ret.NextToken, nil
}

func (self *SRegion) CreateEIP(eip *cloudprovider.SEip) (cloudprovider.ICloudEIP, error) {
	params := map[string]string{}
	if len(eip.NetworkExternalId) > 0 {
		params["SubnetId"] = eip.NetworkExternalId
	}
	if eip.BandwidthMbps > 0 {
		params["Bandwidth"] = strconv.Itoa(eip.BandwidthMbps)
	}
	if len(eip.IP) > 0 {
		params["DesiredAddress"] = eip.IP
	}
	resp, err := self.invoke("AllocateAddress", params)
	if err != nil {
		return nil, errors.Wrapf(err, "AllocateAddress")
	}
	var ip string
	err = resp.Unmarshal(&ip, "publicIp")
	if err != nil {
		return nil, errors.Wrapf(err, "AllocateAddress")
	}
	return self.GetIEipById(ip)
}

func (self *SRegion) GetIEips() ([]cloudprovider.ICloudEIP, error) {
	part, nextToken, err := self.GetEips("", "", "")
	if err != nil {
		return nil, err
	}
	var eips []SEip
	eips = append(eips, part...)
	for len(nextToken) > 0 {
		part, nextToken, err = self.GetEips("", "", nextToken)
		if err != nil {
			return nil, err
		}
		eips = append(eips, part...)
	}
	var ret []cloudprovider.ICloudEIP
	for i := range eips {
		eips[i].region = self
		ret = append(ret, &eips[i])
	}
	return ret, nil
}

func (self *SRegion) GetIEipById(id string) (cloudprovider.ICloudEIP, error) {
	eips, _, err := self.GetEips(id, "", "")
	if err != nil {
		return nil, err
	}
	for i := range eips {
		if eips[i].GetGlobalId() == id {
			eips[i].region = self
			return &eips[i], nil
		}
	}
	return nil, errors.Wrapf(cloudprovider.ErrNotFound, id)
}
