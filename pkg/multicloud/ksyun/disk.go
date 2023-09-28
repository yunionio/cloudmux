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
	"context"
	"fmt"

	"github.com/pkg/errors"
	api "yunion.io/x/cloudmux/pkg/apis/compute"
	"yunion.io/x/cloudmux/pkg/cloudprovider"
	"yunion.io/x/cloudmux/pkg/multicloud"
)

type SDiskResp struct {
	RequestID  string  `json:"RequestId"`
	Volumes    []SDisk `json:"Volumes"`
	TotalCount int     `json:"TotalCount"`
}
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
	SKsTag

	VolumeID           string              `json:"VolumeId"`
	VolumeName         string              `json:"VolumeName"`
	VolumeDesc         string              `json:"VolumeDesc,omitempty"`
	Size               int                 `json:"Size"`
	VolumeStatus       string              `json:"VolumeStatus"`
	VolumeType         string              `json:"VolumeType"`
	VolumeCategory     string              `json:"VolumeCategory"`
	InstanceID         string              `json:"InstanceId"`
	AvailabilityZone   string              `json:"AvailabilityZone"`
	ChargeType         string              `json:"ChargeType"`
	InstanceTradeType  int                 `json:"InstanceTradeType"`
	CreateTime         string              `json:"CreateTime"`
	Attachment         []Attachment        `json:"Attachment"`
	ProjectID          int                 `json:"ProjectId"`
	ExpireTime         string              `json:"ExpireTime,omitempty"`
	HistoryAttachment  []HistoryAttachment `json:"HistoryAttachment,omitempty"`
	DeleteWithInstance bool                `json:"DeleteWithInstance"`
}

type SDiskListInput struct {
	InstanceId string
	DiskId     string
	DiskName   string
}

func (region *SRegion) GetDisks(input *SDiskListInput) ([]SDisk, error) {
	return region.getDisks(input)
}

func (region *SRegion) getDisks(input *SDiskListInput) ([]SDisk, error) {
	pageNum := 0
	disks := []SDisk{}
	param := map[string]string{
		"MaxResults": "100",
		"Marker":     fmt.Sprintf("%d", pageNum),
	}
	if input != nil {
		if len(input.DiskId) > 0 {
			param["VolumeId.1"] = input.DiskId
		}
	}

	for {
		resp, err := region.client.request("ebs", region.Region, "DescribeVolumes", "2016-03-04", param)
		if err != nil {
			return nil, errors.Wrap(err, "list instance")
		}
		part := SDiskResp{}
		err = resp.Unmarshal(&part)
		if err != nil {
			return nil, errors.Wrap(err, "unmarshal instances")
		}
		disks = append(disks, part.Volumes...)
		if len(disks) >= part.TotalCount {
			break
		}
		pageNum++
		param["Marker"] = fmt.Sprintf("%d", pageNum)
	}

	return disks, nil
}

func (disk *SDisk) GetIStorage() (cloudprovider.ICloudStorage, error) {
	if disk.storage == nil {
		return nil, fmt.Errorf("disk %s(%s) missing storage", disk.VolumeName, disk.VolumeID)
	}
	return disk.storage, nil
}

func (disk *SDisk) GetIStorageId() string {
	return disk.VolumeCategory
}

func (disk *SDisk) GetDiskFormat() string {
	return ""
}

func (disk *SDisk) GetId() string {
	return disk.VolumeID
}

func (disk *SDisk) GetGlobalId() string {
	return disk.VolumeID
}

func (disk *SDisk) GetName() string {
	return disk.VolumeName
}

func (disk *SDisk) GetStatus() string {
	// creating、available、attaching、inuse、detaching、extending、deleting、error
	switch disk.VolumeStatus {
	case "available", "inuse", "in-use":
		return api.DISK_READY
	case "detaching":
		return api.DISK_DETACHING
	case "error":
		return api.DISK_UNKNOWN
	default:
		return disk.VolumeStatus
	}
}

func (disk *SDisk) GetDiskSizeMB() int {
	return disk.Size * 1024
}

func (disk *SDisk) GetIsAutoDelete() bool {
	return disk.DeleteWithInstance
}

func (disk *SDisk) GetTemplateId() string {
	return ""
}

func (disk *SDisk) GetDiskType() string {
	if disk.VolumeCategory == "system" {
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
