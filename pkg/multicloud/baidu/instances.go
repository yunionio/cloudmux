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
	"time"

	api "yunion.io/x/cloudmux/pkg/apis/compute"
	"yunion.io/x/cloudmux/pkg/cloudprovider"
	"yunion.io/x/cloudmux/pkg/multicloud"
	"yunion.io/x/pkg/errors"
)

type SInstance struct {
	multicloud.SInstanceBase
	SBaiduTag
	host *SHost

	Id                    string        `json:"id"`
	Name                  string        `json:"name"`
	RoleName              string        `json:"roleName"`
	Hostname              string        `json:"hostname"`
	InstanceType          string        `json:"instanceType"`
	Spec                  string        `json:"spec"`
	Status                string        `json:"status"`
	Desc                  string        `json:"desc"`
	CreatedFrom           string        `json:"createdFrom"`
	PaymentTiming         string        `json:"paymentTiming"`
	CreateTime            time.Time     `json:"createTime"`
	ExpireTime            time.Time     `json:"expireTime"`
	InternalIP            string        `json:"internalIp"`
	PublicIP              string        `json:"publicIp"`
	CPUCount              int           `json:"cpuCount"`
	IsomerismCard         string        `json:"isomerismCard"`
	CardCount             string        `json:"cardCount"`
	NpuVideoMemory        string        `json:"npuVideoMemory"`
	MemoryCapacityInGB    int           `json:"memoryCapacityInGB"`
	LocalDiskSizeInGB     int           `json:"localDiskSizeInGB"`
	ImageID               string        `json:"imageId"`
	PlacementPolicy       string        `json:"placementPolicy"`
	SubnetID              string        `json:"subnetId"`
	VpcID                 string        `json:"vpcId"`
	ZoneName              string        `json:"zoneName"`
	DedicatedHostID       string        `json:"dedicatedHostId"`
	DeletionProtection    int           `json:"deletionProtection"`
	DeploysetList         []interface{} `json:"deploysetList"`
	AutoRenew             bool          `json:"autoRenew"`
	Ipv6                  string        `json:"ipv6"`
	NicInfo               NicInfo       `json:"nicInfo"`
	EniNum                string        `json:"eniNum"`
	Tags                  []Tags        `json:"tags"`
	NetworkCapacityInMbps int           `json:"networkCapacityInMbps"`
}
type Ips struct {
	PrivateIP       string `json:"privateIp"`
	Eip             string `json:"eip"`
	Primary         string `json:"primary"`
	EipID           string `json:"eipId"`
	EipAllocationID string `json:"eipAllocationId"`
	EipSize         string `json:"eipSize"`
	EipStatus       string `json:"eipStatus"`
	EipGroupID      string `json:"eipGroupId"`
	EipType         string `json:"eipType"`
}
type NicInfo struct {
	EniID          string        `json:"eniId"`
	EniUUID        string        `json:"eniUuid"`
	Name           string        `json:"name"`
	Type           string        `json:"type"`
	SubnetID       string        `json:"subnetId"`
	SubnetType     string        `json:"subnetType"`
	Az             string        `json:"az"`
	Description    string        `json:"description"`
	DeviceID       string        `json:"deviceId"`
	Status         string        `json:"status"`
	MacAddress     string        `json:"macAddress"`
	VpcID          string        `json:"vpcId"`
	CreatedTime    string        `json:"createdTime"`
	EniNum         int           `json:"eniNum"`
	EriNum         int           `json:"eriNum"`
	EriInfos       []interface{} `json:"eriInfos"`
	Ips            []Ips         `json:"ips"`
	SecurityGroups []interface{} `json:"securityGroups"`
}
type Tags struct {
	TagKey   string `json:"tagKey"`
	TagValue string `json:"tagValue"`
}

func (region *SRegion) GetInstances() ([]SInstance, error) {
	resp, err := region.client.list("bcc", region.Region, "v2/instance", nil)
	if err != nil {
		return nil, errors.Wrap(err, "list instance")
	}
	instances := []SInstance{}
	err = resp.Unmarshal(&instances, "instances")
	if err != nil {
		return nil, errors.Wrap(err, "unmarshal instances")
	}
	return instances, nil
}

func (ins *SInstance) AssignSecurityGroup(secgroupId string) error {
	return nil
}

func (ins *SInstance) AttachDisk(ctx context.Context, diskId string) error {
	return nil
}

func (ins *SInstance) ChangeConfig(ctx context.Context, config *cloudprovider.SManagedVMChangeConfig) error {
	return nil
}

func (ins *SInstance) DeleteVM(ctx context.Context) error {
	return nil
}

func (ins *SInstance) DetachDisk(ctx context.Context, diskId string) error {
	return nil
}

func (ins *SInstance) GetBios() cloudprovider.TBiosType {
	return ""
}

func (ins *SInstance) GetBootOrder() string {
	return ""
}

func (ins *SInstance) GetError() error {
	return nil
}

func (ins *SInstance) GetFullOsName() string {
	return ""
}

func (ins *SInstance) GetGlobalId() string {
	return ins.Id
}

func (ins *SInstance) GetId() string {
	return ins.Id
}

func (ins *SInstance) GetInstanceType() string {
	return ins.InstanceType
}

func (ins *SInstance) GetMachine() string {
	return ins.InstanceType
}

func (ins *SInstance) GetHostname() string {
	return ins.Hostname
}

func (ins *SInstance) GetName() string {
	return ins.Name
}

func (ins *SInstance) GetOsArch() string {
	return ""
}

func (ins *SInstance) GetOsDist() string {
	return ""
}

func (ins *SInstance) GetOsLang() string {
	return ""
}

func (ins *SInstance) GetOsType() cloudprovider.TOsType {
	return ""
}

func (ins *SInstance) GetOsVersion() string {
	return ""
}

func (ins *SInstance) GetProjectId() string {
	return ins.host.projectId
}

func (ins *SInstance) GetSecurityGroupIds() ([]string, error) {
	return nil, nil
}

func (ins *SInstance) GetStatus() string {
	switch ins.Status {
	case "Running":
		return api.VM_RUNNING
	case "Stopped":
		return api.VM_READY
	case "Stopping":
		return api.VM_STOPPING
	case "Starting":
		return api.VM_STARTING
	}
	return api.VM_UNKNOWN
}

func (ins *SInstance) GetHypervisor() string {
	return ""
}

func (ins *SInstance) GetIDisks() ([]cloudprovider.ICloudDisk, error) {
	return nil, nil
}

func (ins *SInstance) GetIEIP() (cloudprovider.ICloudEIP, error) {
	return nil, nil
}

func (ins *SInstance) GetINics() ([]cloudprovider.ICloudNic, error) {
	return nil, nil
}

func (ins *SInstance) GetVNCInfo(input *cloudprovider.ServerVncInput) (*cloudprovider.ServerVncOutput, error) {
	return nil, nil
}

func (ins *SInstance) GetVcpuCount() int {
	return ins.CPUCount
}

func (ins *SInstance) GetVmemSizeMB() int {
	return ins.MemoryCapacityInGB * 1024
}

func (ins *SInstance) GetVdi() string {
	return ""
}

func (ins *SInstance) GetVga() string {
	return ""
}

func (ins *SInstance) RebuildRoot(ctx context.Context, config *cloudprovider.SManagedVMRebuildRootConfig) (string, error) {
	return "", nil
}

func (ins *SInstance) SetSecurityGroups(secgroupIds []string) error {
	return nil
}

func (ins *SInstance) StartVM(ctx context.Context) error {
	return nil
}

func (ins *SInstance) StopVM(ctx context.Context, opts *cloudprovider.ServerStopOptions) error {
	return nil
}

func (ins *SInstance) UpdateUserData(userData string) error {
	return nil
}

func (ins *SInstance) UpdateVM(ctx context.Context, input cloudprovider.SInstanceUpdateOptions) error {
	return nil
}

func (ins *SInstance) GetIHost() cloudprovider.ICloudHost {
	return nil
}

func (ins *SInstance) DeployVM(ctx context.Context, name string, username string, password string, publicKey string, deleteKeypair bool, description string) error {
	return nil
}
