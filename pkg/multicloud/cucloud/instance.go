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

package cucloud

import (
	"context"

	"github.com/pkg/errors"
	"yunion.io/x/cloudmux/pkg/cloudprovider"
	"yunion.io/x/cloudmux/pkg/multicloud"
	"yunion.io/x/pkg/util/httputils"
)

type SInstance struct {
	multicloud.SInstanceBase
	SQingTag
	host *SHost

	InstanceID    string `json:"instanc_id"`
	InstanceName  string
	VcpusCurrent  int
	MemoryCurrent int
	InstanceType  string
	Status        string
}

func (region *SRegion) GetInstances() ([]SInstance, error) {
	resp, err := region.client.request(httputils.POST, "instance/v1/product/order/ecs/instances", map[string]interface{}{
		"cloudRegionCode": region.GetId(),
	})
	if err != nil {
		return nil, errors.Wrap(err, "list instance")
	}
	instances := []SInstance{}
	err = resp.Unmarshal(&instances, "InstancesSet")
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
	return ins.InstanceID
}

func (ins *SInstance) GetId() string {
	return ins.InstanceID
}

func (ins *SInstance) GetInstanceType() string {
	return ins.InstanceType
}

func (ins *SInstance) GetMachine() string {
	return ins.InstanceType
}

func (ins *SInstance) GetHostname() string {
	return ""
}

func (ins *SInstance) GetName() string {
	return ins.InstanceName
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
	return ins.Status
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
	return ins.VcpusCurrent
}

func (ins *SInstance) GetVmemSizeMB() int {
	return ins.MemoryCurrent
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
