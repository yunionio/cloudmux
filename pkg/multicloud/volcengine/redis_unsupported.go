// Copyright 2026 Yunion
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package volcengine

import (
	"yunion.io/x/cloudmux/pkg/cloudprovider"
	"yunion.io/x/jsonutils"
)

func (i *SRedisInstance) GetICloudElasticcacheAccounts() ([]cloudprovider.ICloudElasticcacheAccount, error) {
	return nil, cloudprovider.ErrNotSupported
}
func (i *SRedisInstance) GetICloudElasticcacheAcls() ([]cloudprovider.ICloudElasticcacheAcl, error) {
	return nil, cloudprovider.ErrNotSupported
}
func (i *SRedisInstance) GetICloudElasticcacheBackups() ([]cloudprovider.ICloudElasticcacheBackup, error) {
	return nil, cloudprovider.ErrNotSupported
}
func (i *SRedisInstance) GetICloudElasticcacheParameters() ([]cloudprovider.ICloudElasticcacheParameter, error) {
	return nil, cloudprovider.ErrNotSupported
}
func (i *SRedisInstance) GetICloudElasticcacheAccount(accountId string) (cloudprovider.ICloudElasticcacheAccount, error) {
	return nil, cloudprovider.ErrNotSupported
}
func (i *SRedisInstance) GetICloudElasticcacheAcl(aclId string) (cloudprovider.ICloudElasticcacheAcl, error) {
	return nil, cloudprovider.ErrNotSupported
}
func (i *SRedisInstance) GetICloudElasticcacheBackup(backupId string) (cloudprovider.ICloudElasticcacheBackup, error) {
	return nil, cloudprovider.ErrNotSupported
}
func (i *SRedisInstance) Restart() error                       { return cloudprovider.ErrNotSupported }
func (i *SRedisInstance) Delete() error                        { return cloudprovider.ErrNotSupported }
func (i *SRedisInstance) ChangeInstanceSpec(spec string) error { return cloudprovider.ErrNotSupported }
func (i *SRedisInstance) SetMaintainTime(maintainStartTime, maintainEndTime string) error {
	return cloudprovider.ErrNotSupported
}
func (i *SRedisInstance) AllocatePublicConnection(port int) (string, error) {
	return "", cloudprovider.ErrNotSupported
}
func (i *SRedisInstance) ReleasePublicConnection() error { return cloudprovider.ErrNotSupported }
func (i *SRedisInstance) CreateAccount(account cloudprovider.SCloudElasticCacheAccountInput) (cloudprovider.ICloudElasticcacheAccount, error) {
	return nil, cloudprovider.ErrNotSupported
}
func (i *SRedisInstance) CreateAcl(aclName, securityIps string) (cloudprovider.ICloudElasticcacheAcl, error) {
	return nil, cloudprovider.ErrNotSupported
}
func (i *SRedisInstance) CreateBackup(desc string) (cloudprovider.ICloudElasticcacheBackup, error) {
	return nil, cloudprovider.ErrNotSupported
}
func (i *SRedisInstance) FlushInstance(input cloudprovider.SCloudElasticCacheFlushInstanceInput) error {
	return cloudprovider.ErrNotSupported
}
func (i *SRedisInstance) UpdateAuthMode(noPasswordAccess bool, password string) error {
	return cloudprovider.ErrNotSupported
}
func (i *SRedisInstance) UpdateInstanceParameters(config jsonutils.JSONObject) error {
	return cloudprovider.ErrNotSupported
}
func (i *SRedisInstance) UpdateBackupPolicy(config cloudprovider.SCloudElasticCacheBackupPolicyUpdateInput) error {
	return cloudprovider.ErrNotSupported
}
func (i *SRedisInstance) UpdateSecurityGroups(secgroupIds []string) error {
	return cloudprovider.ErrNotSupported
}
