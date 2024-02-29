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
	"yunion.io/x/jsonutils"
)

type SStorage struct {
	multicloud.SStorageBase

	zone           *SZone
	ZoneId         string
	StorageType    string
	MediumType     string
	CapacityMb     int64
	CapacityUsedMb int64
	// Enabled        bool
	// SkipSync       bool
}

var ksDiskTypes = []string{"ESSD_PL1", "ESSD_PL2", "ESSD_PL3", "SSD3.0", "EHDD"}

func (storage *SStorage) GetId() string {
	return fmt.Sprintf("%s-%s-%s", storage.zone.region.client.cpcfg.Id, storage.zone.GetId(), storage.StorageType)
}

func (storage *SStorage) GetName() string {
	return fmt.Sprintf("%s-%s-%s", storage.zone.region.client.cpcfg.Name, storage.zone.GetId(), storage.StorageType)
}

func (storage *SStorage) GetGlobalId() string {
	return fmt.Sprintf("%s-%s-%s", storage.zone.region.client.cpcfg.Id, storage.zone.GetGlobalId(), storage.StorageType)
}

func (storage *SStorage) IsEmulated() bool {
	return true
}

func (storage *SStorage) GetIZone() cloudprovider.ICloudZone {
	return storage.zone
}

func (storage *SStorage) GetIDisks() ([]cloudprovider.ICloudDisk, error) {
	disks, err := storage.zone.region.GetDisks(SDiskListInput{})
	if err != nil {
		return nil, err
	}
	ret := []cloudprovider.ICloudDisk{}
	for i := range disks {
		if fmt.Sprintf("%s-%s-%s", storage.zone.region.client.cpcfg.Id, storage.zone.GetId(), disks[i].VolumeType) != storage.GetId() {
			continue
		}
		disks[i].storage = storage
		ret = append(ret, &disks[i])
	}
	return ret, nil
}

func (storage *SStorage) GetStorageType() string {
	return storage.StorageType
}

func (storage *SStorage) GetMediumType() string {
	return storage.MediumType
}

func (storage *SStorage) GetCapacityMB() int64 {
	return storage.CapacityMb
}

func (storage *SStorage) GetCapacityUsedMB() int64 {
	return storage.CapacityUsedMb
}

func (storage *SStorage) GetStorageConf() jsonutils.JSONObject {
	return jsonutils.NewDict()
}

func (storage *SStorage) GetEnabled() bool {
	return true
}

func (storage *SStorage) CreateIDisk(conf *cloudprovider.DiskCreateConfig) (cloudprovider.ICloudDisk, error) {
	return nil, cloudprovider.ErrNotSupported
}

func (storage *SStorage) GetIDiskById(id string) (cloudprovider.ICloudDisk, error) {
	disks, err := storage.GetIDisks()
	if err != nil {
		return nil, err
	}
	for i := range disks {
		if disks[i].GetGlobalId() == id {
			return disks[i], nil
		}
	}
	return nil, cloudprovider.ErrNotFound
}

func (storage *SStorage) GetMountPoint() string {
	return ""
}

func (storage *SStorage) IsSysDiskStore() bool {
	return true
}

func (storage *SStorage) DisableSync() bool {
	return false
}

func (storage *SStorage) GetIStoragecache() cloudprovider.ICloudStoragecache {
	return nil
}

func (storage *SStorage) GetStatus() string {
	return "ready"
}
