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

package huawei

import (
	"context"
	"fmt"
	"net/url"
	"strings"
	"time"

	"yunion.io/x/jsonutils"
	"yunion.io/x/log"
	"yunion.io/x/pkg/errors"

	billing_api "yunion.io/x/cloudmux/pkg/apis/billing"
	api "yunion.io/x/cloudmux/pkg/apis/compute"
	"yunion.io/x/cloudmux/pkg/cloudprovider"
	"yunion.io/x/cloudmux/pkg/multicloud"
)

type Attachment struct {
	ServerID     string `json:"server_id"`
	AttachmentID string `json:"attachment_id"`
	AttachedAt   string `json:"attached_at"`
	HostName     string `json:"host_name"`
	VolumeID     string `json:"volume_id"`
	Device       string `json:"device"`
	ID           string `json:"id"`
}

type DiskMeta struct {
	ResourceSpecCode string `json:"resourceSpecCode"`
	Billing          string `json:"billing"`
	ResourceType     string `json:"resourceType"`
	AttachedMode     string `json:"attached_mode"`
	Readonly         string `json:"readonly"`
}

type VolumeImageMetadata struct {
	QuickStart             string `json:"__quick_start"`
	ContainerFormat        string `json:"container_format"`
	MinRAM                 string `json:"min_ram"`
	ImageName              string `json:"image_name"`
	ImageID                string `json:"image_id"`
	OSType                 string `json:"__os_type"`
	OSFeatureList          string `json:"__os_feature_list"`
	MinDisk                string `json:"min_disk"`
	SupportKVM             string `json:"__support_kvm"`
	VirtualEnvType         string `json:"virtual_env_type"`
	SizeGB                 string `json:"size"`
	OSVersion              string `json:"__os_version"`
	OSBit                  string `json:"__os_bit"`
	SupportKVMHi1822Hiovs  string `json:"__support_kvm_hi1822_hiovs"`
	SupportXen             string `json:"__support_xen"`
	Description            string `json:"__description"`
	Imagetype              string `json:"__imagetype"`
	DiskFormat             string `json:"disk_format"`
	ImageSourceType        string `json:"__image_source_type"`
	Checksum               string `json:"checksum"`
	Isregistered           string `json:"__isregistered"`
	HwVifMultiqueueEnabled string `json:"hw_vif_multiqueue_enabled"`
	Platform               string `json:"__platform"`
}

// https://support.huaweicloud.com/api-evs/zh-cn_topic_0124881427.html
type SDisk struct {
	storage *SStorage
	multicloud.SDisk
	HuaweiDiskTags
	details *SResourceDetail

	ID                  string              `json:"id"`
	Name                string              `json:"name"`
	Status              string              `json:"status"`
	Attachments         []Attachment        `json:"attachments"`
	Description         string              `json:"description"`
	SizeGB              int                 `json:"size"`
	Metadata            DiskMeta            `json:"metadata"`
	Encrypted           bool                `json:"encrypted"`
	Bootable            string              `json:"bootable"`
	Multiattach         bool                `json:"multiattach"`
	AvailabilityZone    string              `json:"availability_zone"`
	SourceVolid         string              `json:"source_volid"`
	SnapshotID          string              `json:"snapshot_id"`
	CreatedAt           time.Time           `json:"created_at"`
	VolumeType          string              `json:"volume_type"`
	VolumeImageMetadata VolumeImageMetadata `json:"volume_image_metadata"`
	ReplicationStatus   string              `json:"replication_status"`
	UserID              string              `json:"user_id"`
	ConsistencygroupID  string              `json:"consistencygroup_id"`
	UpdatedAt           string              `json:"updated_at"`
	EnterpriseProjectId string

	ExpiredTime time.Time
}

func (self *SDisk) GetId() string {
	return self.ID
}

func (self *SDisk) GetName() string {
	if len(self.Name) == 0 {
		return self.ID
	}

	return self.Name
}

func (self *SDisk) GetGlobalId() string {
	return self.ID
}

func (self *SDisk) GetStatus() string {
	// https://support.huaweicloud.com/api-evs/zh-cn_topic_0051803385.html
	switch self.Status {
	case "creating", "downloading":
		return api.DISK_ALLOCATING
	case "available", "in-use":
		return api.DISK_READY
	case "error":
		return api.DISK_ALLOC_FAILED
	case "attaching":
		return api.DISK_ATTACHING
	case "detaching":
		return api.DISK_DETACHING
	case "restoring-backup":
		return api.DISK_REBUILD
	case "backing-up":
		return api.DISK_BACKUP_STARTALLOC // ?
	case "error_restoring":
		return api.DISK_BACKUP_ALLOC_FAILED
	case "uploading":
		return api.DISK_SAVING //?
	case "extending":
		return api.DISK_RESIZING
	case "error_extending":
		return api.DISK_ALLOC_FAILED // ?
	case "deleting":
		return api.DISK_DEALLOC //?
	case "error_deleting":
		return api.DISK_DEALLOC_FAILED // ?
	case "rollbacking":
		return api.DISK_REBUILD
	case "error_rollbacking":
		return api.DISK_UNKNOWN
	default:
		return api.DISK_UNKNOWN
	}
}

func (self *SDisk) Refresh() error {
	new, err := self.storage.zone.region.GetDisk(self.GetId())
	if err != nil {
		return err
	}
	return jsonutils.Update(self, new)
}

func (self *SDisk) IsEmulated() bool {
	return false
}

func (self *SDisk) getResourceDetails() *SResourceDetail {
	if self.details != nil {
		return self.details
	}

	res, err := self.storage.zone.region.GetOrderResourceDetail(self.GetId())
	if err != nil {
		log.Debugln(err)
		return nil
	}

	self.details = &res
	return self.details
}

func (self *SDisk) GetBillingType() string {
	details := self.getResourceDetails()
	if details == nil {
		return billing_api.BILLING_TYPE_POSTPAID
	} else {
		return billing_api.BILLING_TYPE_PREPAID
	}
}

func (self *SDisk) GetCreatedAt() time.Time {
	return self.CreatedAt
}

func (self *SDisk) GetExpiredAt() time.Time {
	var expiredTime time.Time
	details := self.getResourceDetails()
	if details != nil {
		expiredTime = details.ExpireTime
	}

	return expiredTime
}

func (self *SDisk) GetIStorage() (cloudprovider.ICloudStorage, error) {
	return self.storage, nil
}

func (self *SDisk) GetDiskFormat() string {
	// self.volume_type ?
	return "vhd"
}

func (self *SDisk) GetDiskSizeMB() int {
	return int(self.SizeGB * 1024)
}

func (self *SDisk) checkAutoDelete(attachments []Attachment) bool {
	autodelete := false
	for _, attach := range attachments {
		if len(attach.ServerID) > 0 {
			// todo : 忽略错误？？
			vm, err := self.storage.zone.region.GetInstanceByID(attach.ServerID)
			if err != nil {
				volumes := vm.OSExtendedVolumesVolumesAttached
				for _, vol := range volumes {
					if vol.ID == self.ID && strings.ToLower(vol.DeleteOnTermination) == "true" {
						autodelete = true
					}
				}
			}

			break
		}
	}

	return autodelete
}

func (self *SDisk) GetIsAutoDelete() bool {
	if len(self.Attachments) > 0 {
		return self.checkAutoDelete(self.Attachments)
	}

	return false
}

func (self *SDisk) GetTemplateId() string {
	return self.VolumeImageMetadata.ImageID
}

// Bootable 表示硬盘是否为启动盘。
// 启动盘 != 系统盘(必须是启动盘且挂载在root device上)
func (self *SDisk) GetDiskType() string {
	if self.Bootable == "true" {
		return api.DISK_TYPE_SYS
	} else {
		return api.DISK_TYPE_DATA
	}
}

func (self *SDisk) GetFsFormat() string {
	return ""
}

func (self *SDisk) GetIsNonPersistent() bool {
	return false
}

func (self *SDisk) GetDriver() string {
	// https://support.huaweicloud.com/api-evs/zh-cn_topic_0058762431.html
	// scsi or vbd?
	// todo: implement me
	return "scsi"
}

func (self *SDisk) GetCacheMode() string {
	return "none"
}

func (self *SDisk) GetMountpoint() string {
	if len(self.Attachments) > 0 {
		return self.Attachments[0].Device
	}

	return ""
}

func (self *SDisk) GetMountServerId() string {
	if len(self.Attachments) > 0 {
		return self.Attachments[0].ServerID
	}

	return ""
}

func (self *SDisk) GetAccessPath() string {
	return ""
}

func (self *SDisk) Delete(ctx context.Context) error {
	disk, err := self.storage.zone.region.GetDisk(self.GetId())
	if err != nil {
		if errors.Cause(err) == cloudprovider.ErrNotFound {
			return nil
		}
		return err
	}
	if disk.Status != "deleting" {
		// 等待硬盘ready
		cloudprovider.WaitStatus(self, api.DISK_READY, 5*time.Second, 60*time.Second)
		err := self.storage.zone.region.DeleteDisk(self.GetId())
		if err != nil {
			return err
		}
	}

	return cloudprovider.WaitDeleted(self, 10*time.Second, 120*time.Second)
}

func (self *SDisk) CreateISnapshot(ctx context.Context, name string, desc string) (cloudprovider.ICloudSnapshot, error) {
	snapshot, err := self.storage.zone.region.CreateSnapshot(self.GetId(), name, desc)
	if err != nil {
		return nil, err
	}
	return snapshot, nil
}

func (self *SDisk) GetISnapshot(id string) (cloudprovider.ICloudSnapshot, error) {
	snapshot, err := self.storage.zone.region.GetSnapshot(id)
	if err != nil {
		return nil, err
	}
	return snapshot, nil
}

func (self *SDisk) GetISnapshots() ([]cloudprovider.ICloudSnapshot, error) {
	snapshots, err := self.storage.zone.region.GetSnapshots(self.ID, "")
	if err != nil {
		return nil, err
	}

	isnapshots := make([]cloudprovider.ICloudSnapshot, len(snapshots))
	for i := 0; i < len(snapshots); i++ {
		isnapshots[i] = &snapshots[i]
	}
	return isnapshots, nil
}

func (self *SDisk) Resize(ctx context.Context, newSizeMB int64) error {
	err := cloudprovider.WaitStatus(self, api.DISK_READY, 5*time.Second, 60*time.Second)
	if err != nil {
		return err
	}

	sizeGb := newSizeMB / 1024
	err = self.storage.zone.region.resizeDisk(self.GetId(), sizeGb)
	if err != nil {
		return err
	}

	return cloudprovider.WaitStatusWithDelay(self, api.DISK_READY, 15*time.Second, 5*time.Second, 60*time.Second)
}

func (self *SDisk) Detach() error {
	err := self.storage.zone.region.DetachDisk(self.GetMountServerId(), self.GetId())
	if err != nil {
		log.Debugf("detach server %s disk %s failed: %s", self.GetMountServerId(), self.GetId(), err)
		return err
	}

	return cloudprovider.WaitCreated(5*time.Second, 60*time.Second, func() bool {
		err := self.Refresh()
		if err != nil {
			log.Debugln(err)
			return false
		}

		if self.Status == "available" {
			return true
		}

		return false
	})
}

func (self *SDisk) Attach(device string) error {
	err := self.storage.zone.region.AttachDisk(self.GetMountServerId(), self.GetId(), device)
	if err != nil {
		log.Debugf("attach server %s disk %s failed: %s", self.GetMountServerId(), self.GetId(), err)
		return err
	}

	return cloudprovider.WaitStatusWithDelay(self, api.DISK_READY, 10*time.Second, 5*time.Second, 60*time.Second)
}

// 在线卸载磁盘 https://support.huaweicloud.com/usermanual-ecs/zh-cn_topic_0036046828.html
// 对于挂载在系统盘盘位（也就是“/dev/sda”或“/dev/vda”挂载点）上的磁盘，当前仅支持离线卸载
func (self *SDisk) Reset(ctx context.Context, snapshotId string) (string, error) {
	mountpoint := self.GetMountpoint()
	if len(mountpoint) > 0 {
		err := self.Detach()
		if err != nil {
			return "", err
		}
	}

	diskId, err := self.storage.zone.region.resetDisk(self.GetId(), snapshotId)
	if err != nil {
		return diskId, err
	}

	err = cloudprovider.WaitStatus(self, api.DISK_READY, 5*time.Second, 300*time.Second)
	if err != nil {
		return "", err
	}

	if len(mountpoint) > 0 {
		err := self.Attach(mountpoint)
		if err != nil {
			return "", err
		}
	}

	return diskId, nil
}

// 华为云不支持重置
func (self *SDisk) Rebuild(ctx context.Context) error {
	return cloudprovider.ErrNotSupported
}

func (self *SRegion) GetDisk(id string) (*SDisk, error) {
	resp, err := self.list(SERVICE_EVS, "cloudvolumes/"+id, nil)
	if err != nil {
		return nil, err
	}
	ret := &SDisk{}
	return ret, resp.Unmarshal(ret, "volume")
}

// https://support.huaweicloud.com/api-evs/zh-cn_topic_0058762430.html
func (self *SRegion) GetDisks(zoneId string) ([]SDisk, error) {
	params := url.Values{}
	if len(zoneId) > 0 {
		params.Set("availability_zone", zoneId)
	}
	ret := []SDisk{}
	for {
		resp, err := self.list(SERVICE_EVS, "cloudvolumes/detail", params)
		if err != nil {
			return nil, err
		}
		part := struct {
			Volumes []SDisk
			Count   int
		}{}
		err = resp.Unmarshal(&part)
		if err != nil {
			return nil, errors.Wrapf(err, "Unmarshal")
		}
		ret = append(ret, part.Volumes...)
		if len(ret) >= part.Count || len(part.Volumes) == 0 {
			break
		}
		params.Set("marker", part.Volumes[len(part.Volumes)-1].ID)
	}
	return ret, nil
}

// https://support.huaweicloud.com/api-evs/zh-cn_topic_0058762427.html
func (self *SRegion) CreateDisk(zoneId string, category string, name string, sizeGb int, snapshotId string, desc string, projectId string) (string, error) {
	params := map[string]interface{}{}
	obj := map[string]interface{}{
		"name":              name,
		"availability_zone": zoneId,
		"description":       desc,
		"volume_type":       category,
		"size":              sizeGb,
	}
	if len(snapshotId) > 0 {
		obj["snapshot_id"] = snapshotId
	}
	if len(projectId) > 0 {
		obj["enterprise_project_id"] = projectId
	}
	params["volume"] = obj
	// 目前只支持创建按需资源，返回job id。 如果创建包年包月资源则返回order id
	resp, err := self.post(SERVICE_EVS, "cloudvolumes", params)
	if err != nil {
		return "", errors.Wrapf(err, "create volume")
	}
	ret := struct {
		JobId     string
		OrderId   string
		VolumeIds []string
	}{}
	err = resp.Unmarshal(&ret)
	if err != nil {
		return "", errors.Wrapf(err, "Unmarshal")
	}
	id := ret.JobId
	if len(ret.OrderId) > 0 {
		id = ret.OrderId
	}

	// 按需计费
	volumeId, err := self.GetTaskEntityID(SERVICE_EVS, id, "volume_id")
	if err != nil {
		return "", errors.Wrap(err, "GetAllSubTaskEntityIDs")
	}

	if len(volumeId) == 0 {
		return "", errors.Errorf("CreateInstance job %s result is emtpy", id)
	}
	return volumeId, nil
}

// https://support.huaweicloud.com/api-evs/zh-cn_topic_0058762428.html
// 默认删除云硬盘关联的所有快照
func (self *SRegion) DeleteDisk(diskId string) error {
	resource := fmt.Sprintf("cloudvolumes/%s", diskId)
	_, err := self.delete(SERVICE_EVS, resource)
	return err
}

/*
扩容状态为available的云硬盘时，没有约束限制。
扩容状态为in-use的云硬盘时，有以下约束：
不支持共享云硬盘，即multiattach参数值必须为false。
云硬盘所挂载的云服务器状态必须为ACTIVE、PAUSED、SUSPENDED、SHUTOFF才支持扩容
*/
func (self *SRegion) resizeDisk(diskId string, sizeGB int64) error {
	params := map[string]interface{}{
		"os-extend": map[string]interface{}{
			"new_size": sizeGB,
		},
	}
	_, err := self.post(SERVICE_EVS, fmt.Sprintf("cloudvolumes/%s/action", diskId), params)
	return err
}

/*
https://support.huaweicloud.com/api-evs/zh-cn_topic_0051408629.html
只支持快照回滚到源云硬盘，不支持快照回滚到其它指定云硬盘。
只有云硬盘状态处于“available”或“error_rollbacking”状态才允许快照回滚到源云硬盘。
*/
func (self *SRegion) resetDisk(diskId, snapshotId string) (string, error) {
	params := map[string]interface{}{
		"rollback": map[string]interface{}{
			"volume_id": diskId,
		},
	}
	resource := fmt.Sprintf("cloudsnapshots/%s/rollback", snapshotId)
	resp, err := self.post(SERVICE_EVS, resource, params)
	if err != nil {
		return "", err
	}
	return resp.GetString("rollback", "volume_id")
}

func (self *SDisk) GetProjectId() string {
	return self.EnterpriseProjectId
}
