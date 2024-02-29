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

type SRegion struct {
	multicloud.SRegion
	multicloud.SNoObjectStorageRegion
	multicloud.SNoLbRegion
	client *SKsyunClient

	Region     string
	RegionName string
}

func (region *SRegion) GetId() string {
	return region.Region
}

func (region *SRegion) GetGlobalId() string {
	return fmt.Sprintf("%s/%s", api.CLOUD_PROVIDER_KSYUN, region.Region)
}

func (region *SRegion) GetProvider() string {
	return api.CLOUD_PROVIDER_KSYUN
}

func (region *SRegion) GetCloudEnv() string {
	return api.CLOUD_PROVIDER_KSYUN
}

func (region *SRegion) GetGeographicInfo() cloudprovider.SGeographicInfo {
	geo, ok := map[string]cloudprovider.SGeographicInfo{
		"cn-northwest-1": api.RegionQingYang,
		"ap-singapore-1": api.RegionSingapore,
		"cn-beijing-6":   api.RegionBeijing,
		"cn-guangzhou-1": api.RegionGuangzhou,
		"cn-shanghai-2":  api.RegionShanghai,
	}[region.Region]
	if ok {
		return geo
	}
	return cloudprovider.SGeographicInfo{}
}

func (region *SRegion) GetName() string {
	return region.RegionName
}

func (region *SRegion) GetI18n() cloudprovider.SModelI18nTable {
	table := cloudprovider.SModelI18nTable{}
	table["name"] = cloudprovider.NewSModelI18nEntry(region.GetName()).CN(region.GetName()).EN(region.Region)
	return table
}

func (region *SRegion) GetStatus() string {
	return api.CLOUD_REGION_STATUS_INSERVER
}

func (region *SRegion) GetClient() *SKsyunClient {
	return region.client
}

func (region *SRegion) CreateEIP(opts *cloudprovider.SEip) (cloudprovider.ICloudEIP, error) {
	return nil, cloudprovider.ErrNotImplemented
}

func (region *SRegion) CreateISecurityGroup(conf *cloudprovider.SecurityGroupCreateInput) (cloudprovider.ICloudSecurityGroup, error) {
	return nil, cloudprovider.ErrNotImplemented
}

func (region *SRegion) GetISecurityGroupById(secgroupId string) (cloudprovider.ICloudSecurityGroup, error) {
	return nil, cloudprovider.ErrNotImplemented
}

func (region *SRegion) CreateIVpc(opts *cloudprovider.VpcCreateOptions) (cloudprovider.ICloudVpc, error) {
	return nil, cloudprovider.ErrNotImplemented
}

func (region *SRegion) GetIVpcs() ([]cloudprovider.ICloudVpc, error) {
	vpcs, err := region.GetVpcs([]string{})
	if err != nil {
		return nil, errors.Wrap(err, "GetVpcs")
	}
	res := []cloudprovider.ICloudVpc{}
	for i := 0; i < len(vpcs); i++ {
		vpcs[i].region = region
		res = append(res, &vpcs[i])
	}
	return res, nil
}

func (region *SRegion) GetIVpcById(id string) (cloudprovider.ICloudVpc, error) {
	return nil, cloudprovider.ErrNotImplemented
}

func (region *SRegion) GetCapabilities() []string {
	return region.client.GetCapabilities()
}

func (region *SRegion) GetIEipById(eipId string) (cloudprovider.ICloudEIP, error) {
	eip, err := region.GetEipById([]string{eipId})
	if err != nil {
		return nil, errors.Wrap(err, "GetEipById")
	}
	return &eip[0], nil
}

func (region *SRegion) GetIEips() ([]cloudprovider.ICloudEIP, error) {
	eips, err := region.GetEips()
	if err != nil {
		return nil, errors.Wrap(err, "GetEips")
	}
	res := []cloudprovider.ICloudEIP{}
	for i := 0; i < len(eips); i++ {
		eips[i].region = region
		res = append(res, &eips[i])
	}
	return res, nil
}

func (region *SRegion) GetIZones() ([]cloudprovider.ICloudZone, error) {
	zones, err := region.GetZones()
	if err != nil {
		return nil, errors.Wrap(err, "GetZones")
	}
	res := []cloudprovider.ICloudZone{}
	for i := 0; i < len(zones); i++ {
		zones[i].region = region
		res = append(res, &zones[i])
	}
	return res, nil
}

func (region *SRegion) GetIZoneById(id string) (cloudprovider.ICloudZone, error) {
	return nil, cloudprovider.ErrNotImplemented
}

func (region *SRegion) GetIStorages() ([]cloudprovider.ICloudStorage, error) {
	iStores := make([]cloudprovider.ICloudStorage, 0)
	izones, err := region.GetIZones()
	if err != nil {
		return nil, err
	}
	if len(izones) > 0 {
		iZoneStores, err := izones[0].GetIStorages()
		if err != nil {
			return nil, err
		}
		iStores = append(iStores, iZoneStores...)
	}
	return iStores, nil
}

func (region *SRegion) ecsGetRequest(action string, params map[string]string) (jsonutils.JSONObject, error) {
	return region.client.request("kec", region.Region, action, "2016-03-04", params)
}

func (region *SRegion) eipGetRequest(action string, params map[string]string) (jsonutils.JSONObject, error) {
	return region.client.request("eip", region.Region, action, "2016-03-04", params)
}

func (region *SRegion) diskGetRequest(action string, params map[string]string) (jsonutils.JSONObject, error) {
	return region.client.request("ebs", region.Region, action, "2016-03-04", params)
}

func (region *SRegion) vpcGetRequest(action string, params map[string]string) (jsonutils.JSONObject, error) {
	return region.client.request("vpc", region.Region, action, "2016-03-04", params)
}
