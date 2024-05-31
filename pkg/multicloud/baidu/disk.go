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
	"context"
	"fmt"
	"strings"
	"time"

	api "yunion.io/x/cloudmux/pkg/apis/compute"
	"yunion.io/x/cloudmux/pkg/cloudprovider"
	"yunion.io/x/cloudmux/pkg/multicloud"
	"yunion.io/x/pkg/errors"
)

type Attachment struct {
	InstanceID         string `json:"InstanceId"`
	MountPoint         string `json:"MountPoint"`
	DeleteWithInstance bool   `json:"DeleteWithInstance"`
}

type HistoryAttachment struct {
	InstanceID string `json:"InstanceId"`
	AttachTime string `json:"AttachTime"`
	DetachTime string `json:"DetachTime"`
	MountPoint string `json:"MountPoint"`
}

type SDisk struct {
	storage *SStorage
	multicloud.SDisk
	SBaiduTag

	Id                 string             `json:"id"`
	CreateTime         time.Time          `json:"createTime"`
	ExpireTime         interface{}        `json:"expireTime"`
	Name               string             `json:"name"`
	DiskSizeInGB       int                `json:"diskSizeInGB"`
	Status             string             `json:"status"`
	Type               string             `json:"type"`
	StorageType        string             `json:"storageType"`
	Desc               string             `json:"desc"`
	PaymentTiming      string             `json:"paymentTiming"`
	Attachments        []Attachments      `json:"attachments"`
	RegionID           string             `json:"regionId"`
	SourceSnapshotID   string             `json:"sourceSnapshotId"`
	SnapshotNum        string             `json:"snapshotNum"`
	Tags               []Tags             `json:"tags"`
	AutoSnapshotPolicy AutoSnapshotPolicy `json:"autoSnapshotPolicy"`
	ZoneName           string             `json:"zoneName"`
	IsSystemVolume     bool               `json:"isSystemVolume"`
}

type Attachments struct {
	Id         string `json:"Id"`
	InstanceID string `json:"instanceId"`
	Device     string `json:"device"`
	Serial     string `json:"serial"`
}

type AutoSnapshotPolicy struct {
	ID             string `json:"id"`
	Name           string `json:"name"`
	TimePoints     []int  `json:"timePoints"`
	RepeatWeekdays []int  `json:"repeatWeekdays"`
	RetentionDays  int    `json:"retentionDays"`
	Status         string `json:"status"`
}

func (region *SRegion) GetDisks(storageType, zoneName string) ([]SDisk, error) {
	params := map[string]interface{}{}
	if len(zoneName) > 0 {
		params["zoneName"] = zoneName
	}
	if len(storageType) > 0 {
		params["storageType"] = storageType
	}
	disks := []SDisk{}
	for {
		resp, err := region.client.bccList(region.Region, "v2/volume", params)
		if err != nil {
			return nil, errors.Wrap(err, "list disks")
		}
		temp := []SDisk{}
		err = resp.Unmarshal(&temp, "volumes")
		if err != nil {
			return nil, errors.Wrap(err, "unmarshal disks")
		}
		disks = append(disks, temp...)
		if nextMarker, _ := resp.GetString("nextMarker"); len(nextMarker) > 0 {
			params["marker"] = nextMarker
		} else {
			break
		}
	}
	return disks, nil
}

func (region *SRegion) GetDisk(diskId string) (*SDisk, error) {
	params := map[string]interface{}{}
	disk := SDisk{}
	resp, err := region.client.bccList(region.Region, fmt.Sprintf("v2/volume/%s", diskId), params)
	if err != nil {
		return nil, errors.Wrap(err, "list disks")
	}
	err = resp.Unmarshal(&disk, "volume")
	if err != nil {
		return nil, errors.Wrap(err, "unmarshal disks")
	}
	return &disk, nil
}

func (region *SRegion) GetDiskByInstanceId(instanceId string) ([]SDisk, error) {
	disks := []SDisk{}
	params := map[string]interface{}{
		"MaxResults": "1000",
	}
	params["InstanceId"] = instanceId
	resp, err := region.client.bccList(region.Region, "DescribeInstanceVolumes", params)
	if err != nil {
		return nil, errors.Wrap(err, "DescribeInstanceVolumes")
	}
	return disks, resp.Unmarshal(&disks, "Attachments")
}

func (disk *SDisk) GetIStorage() (cloudprovider.ICloudStorage, error) {
	if disk.storage == nil {
		return nil, fmt.Errorf("disk %s(%s) missing storage", disk.Name, disk.Id)
	}
	return disk.storage, nil
}

func (disk *SDisk) GetIStorageId() string {
	if disk.storage == nil {
		return ""
	}
	return disk.storage.GetGlobalId()
}

func (disk *SDisk) GetDiskFormat() string {
	return ""
}

func (disk *SDisk) GetId() string {
	return disk.Id
}

func (disk *SDisk) GetGlobalId() string {
	return disk.Id
}

func (disk *SDisk) GetName() string {
	return disk.Name
}

func (disk *SDisk) GetStatus() string {
	// creating、available、attaching、inuse、detaching、extending、deleting、error
	switch disk.Status {
	case "Available", "InUse":
		return api.DISK_READY
	case "Detaching":
		return api.DISK_DETACHING
	case "Error":
		return api.DISK_UNKNOWN
	default:
		return strings.ToLower(disk.Status)
	}
}

func (disk *SDisk) GetDiskSizeMB() int {
	return disk.DiskSizeInGB * 1024
}

func (disk *SDisk) GetIsAutoDelete() bool {
	return false
}

func (disk *SDisk) GetTemplateId() string {
	return ""
}

func (disk *SDisk) GetDiskType() string {
	if disk.IsSystemVolume {
		return api.DISK_TYPE_SYS
	}
	return api.DISK_TYPE_DATA
}

func (disk *SDisk) GetFsFormat() string {
	return ""
}

func (disk *SDisk) GetIsNonPersistent() bool {
	return false
}

func (disk *SDisk) GetIops() int {
	return 0
}

func (disk *SDisk) GetDriver() string {
	return ""
}

func (disk *SDisk) GetCacheMode() string {
	return ""
}

func (disk *SDisk) GetMountpoint() string {
	return ""
}

func (disk *SDisk) GetAccessPath() string {
	return ""
}

func (disk *SDisk) Delete(ctx context.Context) error {
	return cloudprovider.ErrNotSupported
}

func (disk *SDisk) CreateISnapshot(ctx context.Context, name string, desc string) (cloudprovider.ICloudSnapshot, error) {
	return nil, cloudprovider.ErrNotSupported
}

func (disk *SDisk) GetISnapshots() ([]cloudprovider.ICloudSnapshot, error) {
	return nil, cloudprovider.ErrNotSupported
}

func (disk *SDisk) GetExtSnapshotPolicyIds() ([]string, error) {
	return nil, cloudprovider.ErrNotSupported
}

func (disk *SDisk) Resize(ctx context.Context, newSizeMB int64) error {
	return cloudprovider.ErrNotSupported
}

func (disk *SDisk) Reset(ctx context.Context, snapshotId string) (string, error) {
	return "", cloudprovider.ErrNotSupported
}

func (disk *SDisk) Rebuild(ctx context.Context) error {
	return cloudprovider.ErrNotSupported
}

func (disk *SDisk) SetStorage(storage SStorage) {
	disk.storage = &storage
}
