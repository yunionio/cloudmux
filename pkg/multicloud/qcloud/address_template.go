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

package qcloud

import (
	"fmt"
	"strings"
	"time"

	"yunion.io/x/jsonutils"
	"yunion.io/x/pkg/errors"

	api "yunion.io/x/cloudmux/pkg/apis"
	"yunion.io/x/cloudmux/pkg/cloudprovider"
	"yunion.io/x/cloudmux/pkg/multicloud"
)

type AddressTemplate struct {
	multicloud.SVirtualResourceBase
	QcloudTags
	region *SRegion

	AddressSet          []string
	AddressTemplateId   string
	AddressTemplateName string
	CreatedTime         time.Time
}

func (self *AddressTemplate) GetId() string {
	return self.AddressTemplateId
}

func (self *AddressTemplate) GetName() string {
	if len(self.AddressTemplateName) > 0 {
		return self.AddressTemplateName
	}
	return self.AddressTemplateId
}

func (self *AddressTemplate) GetGlobalId() string {
	return self.AddressTemplateId
}

func (self *AddressTemplate) GetStatus() string {
	return api.STATUS_AVAILABLE
}

func (self *AddressTemplate) GetCreatedAt() time.Time {
	return self.CreatedTime
}

func (self *AddressTemplate) Refresh() error {
	tpl, err := self.region.GetAddressTemplate(self.AddressTemplateId)
	if err != nil {
		return err
	}
	return jsonutils.Update(self, tpl)
}

func (self *AddressTemplate) GetIpSetType() string {
	for _, addr := range self.AddressSet {
		if strings.Contains(addr, ":") {
			return cloudprovider.IpSetTypeIpv6CidrList
		}
	}
	return cloudprovider.IpSetTypeIpv4CidrList
}

func (self *AddressTemplate) GetAddresses() []string {
	return self.AddressSet
}

func (self *AddressTemplate) Update(opts *cloudprovider.IpSetUpdateOptions) error {
	return self.region.ModifyAddressTemplate(self.AddressTemplateId, opts)
}

func (self *AddressTemplate) Delete() error {
	return self.region.DeleteAddressTemplate(self.AddressTemplateId)
}

func (self *SRegion) AddressList(addressId, addressName string, offset, limit int) ([]AddressTemplate, int, error) {
	params := map[string]string{}
	filter := 0
	if len(addressId) > 0 {
		params[fmt.Sprintf("Filters.%d.Name", filter)] = "address-template-id"
		params[fmt.Sprintf("Filters.%d.Values.0", filter)] = addressId
		filter++
	}
	if len(addressName) > 0 {
		params[fmt.Sprintf("Filters.%d.Name", filter)] = "address-template-name"
		params[fmt.Sprintf("Filters.%d.Values.0", filter)] = addressName
		filter++
	}
	params["Offset"] = fmt.Sprintf("%d", offset)
	if limit == 0 {
		limit = 20
	}
	params["Limit"] = fmt.Sprintf("%d", limit)
	body, err := self.vpcRequest("DescribeAddressTemplates", params)
	if err != nil {
		return nil, 0, err
	}
	addressTemplates := []AddressTemplate{}
	err = body.Unmarshal(&addressTemplates, "AddressTemplateSet")
	if err != nil {
		return nil, 0, err
	}
	total, _ := body.Float("TotalCount")
	return addressTemplates, int(total), nil
}

func (self *SRegion) GetAddressTemplate(id string) (*AddressTemplate, error) {
	templates, _, err := self.AddressList(id, "", 0, 1)
	if err != nil {
		return nil, err
	}
	for i := range templates {
		if templates[i].AddressTemplateId == id {
			templates[i].region = self
			return &templates[i], nil
		}
	}
	return nil, errors.Wrapf(cloudprovider.ErrNotFound, "%s", id)
}

func (self *SRegion) GetIIpSets() ([]cloudprovider.ICloudIpSet, error) {
	ret := []cloudprovider.ICloudIpSet{}
	for {
		part, total, err := self.AddressList("", "", len(ret), 100)
		if err != nil {
			return nil, err
		}
		for i := range part {
			part[i].region = self
			ret = append(ret, &part[i])
		}
		if len(part) == 0 || len(ret) >= total {
			break
		}
	}
	return ret, nil
}

func (self *SRegion) GetIIpSetById(id string) (cloudprovider.ICloudIpSet, error) {
	return self.GetAddressTemplate(id)
}

func (self *SRegion) CreateAddressTemplate(opts *cloudprovider.IpSetCreateOptions) (*AddressTemplate, error) {
	params := map[string]string{
		"AddressTemplateName": opts.Name,
	}
	for i := range opts.Addresses {
		params[fmt.Sprintf("Addresses.%d", i)] = opts.Addresses[i]
	}
	body, err := self.vpcRequest("CreateAddressTemplate", params)
	if err != nil {
		return nil, errors.Wrapf(err, "CreateAddressTemplate")
	}
	ret := &AddressTemplate{region: self}
	err = body.Unmarshal(ret, "AddressTemplate")
	if err != nil {
		return nil, errors.Wrapf(err, "Unmarshal")
	}
	return ret, nil
}

func (self *SRegion) CreateIIpSet(opts *cloudprovider.IpSetCreateOptions) (cloudprovider.ICloudIpSet, error) {
	return self.CreateAddressTemplate(opts)
}

func (self *SRegion) ModifyAddressTemplate(id string, opts *cloudprovider.IpSetUpdateOptions) error {
	params := map[string]string{
		"AddressTemplateId": id,
	}
	if len(opts.Name) > 0 {
		params["AddressTemplateName"] = opts.Name
	}
	for i := range opts.Addresses {
		params[fmt.Sprintf("Addresses.%d", i)] = opts.Addresses[i]
	}
	_, err := self.vpcRequest("ModifyAddressTemplateAttribute", params)
	return errors.Wrapf(err, "ModifyAddressTemplateAttribute")
}

func (self *SRegion) DeleteAddressTemplate(id string) error {
	params := map[string]string{
		"AddressTemplateId": id,
	}
	_, err := self.vpcRequest("DeleteAddressTemplate", params)
	return errors.Wrapf(err, "DeleteAddressTemplate")
}
