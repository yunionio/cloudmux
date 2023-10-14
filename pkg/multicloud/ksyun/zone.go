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

	"yunion.io/x/cloudmux/pkg/cloudprovider"
	"yunion.io/x/cloudmux/pkg/multicloud"
	"yunion.io/x/pkg/errors"
)

type SZone struct {
	multicloud.SResourceBase
	region *SRegion
	host   *SHost
	SKsTag
	iwires    []cloudprovider.ICloudWire
	istorages []cloudprovider.ICloudStorage

	AvailabilityZone string
	/* 支持的磁盘种类集合 */
	storageTypes []string
}

func (region *SRegion) GetZones() ([]SZone, error) {
	params := map[string]string{}
	if len(region.Region) > 0 {
		params = map[string]string{"Region": region.Region}
	}
	resp, err := region.client.request("kec", "", "DescribeAvailabilityZones", "2016-03-04", params)
	if err != nil {
		return nil, errors.Wrap(err, "request zone")
	}
	zones := []SZone{}
	return zones, resp.Unmarshal(&zones, "AvailabilityZoneSet")
}

func (zone *SZone) GetId() string {
	return zone.AvailabilityZone
}

func (zone *SZone) GetName() string {
	return zone.AvailabilityZone
}

func (zone *SZone) GetI18n() cloudprovider.SModelI18nTable {
	return nil
}

func (zone *SZone) GetGlobalId() string {
	return zone.AvailabilityZone
}

func (zone *SZone) GetStatus() string {
	return "enable"
}

func (zone *SZone) Refresh() error {
	return nil
}

func (zone *SZone) IsEmulated() bool {
	return false
}

func (zone *SZone) GetIRegion() cloudprovider.ICloudRegion {
	return zone.region
}

func (zone *SZone) GetIHostById(id string) (cloudprovider.ICloudHost, error) {
	host := zone.getHost()
	if host.GetGlobalId() == id {
		return host, nil
	}
	return nil, cloudprovider.ErrNotFound
}

func (zone *SZone) getHost() *SHost {
	if zone.host == nil {
		zone.host = &SHost{zone: zone}
	}
	return zone.host
}

func (zone *SZone) GetIStorages() ([]cloudprovider.ICloudStorage, error) {
	if zone.istorages == nil {
		err := zone.fetchStorages()
		if err != nil {
			return nil, errors.Wrap(err, "fetchStorages")
		}
	}
	return zone.istorages, nil
}

func (zone *SZone) GetIStorageById(id string) (cloudprovider.ICloudStorage, error) {
	return nil, cloudprovider.ErrNotImplemented
}

func (zone *SZone) GetIWires() ([]cloudprovider.ICloudWire, error) {
	return zone.iwires, nil
}

func (zone *SRegion) getZoneById(id string) (*SZone, error) {
	izones, err := zone.GetIZones()
	if err != nil {
		return nil, err
	}
	for i := 0; i < len(izones); i += 1 {
		zone := izones[i].(*SZone)
		if zone.GetId() == id {
			return zone, nil
		}
	}
	return nil, fmt.Errorf("no such zone %s", id)
}

func (zone *SZone) GetIHosts() ([]cloudprovider.ICloudHost, error) {
	return []cloudprovider.ICloudHost{zone.getHost()}, nil
}

func (zone *SZone) GetDescription() string {
	return ""
}

func (self *SZone) fetchStorages() error {
	categories := []string{"ESSD_PL1", "ESSD_PL2", "ESSD_PL3", "SSD3.0", "EHDD"}
	self.istorages = []cloudprovider.ICloudStorage{}
	for _, sc := range categories {
		storage := SStorage{zone: self, StorageType: sc}
		self.istorages = append(self.istorages, &storage)
	}
	return nil
}
