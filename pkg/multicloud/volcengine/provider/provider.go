// Copyright 2023 Yunion
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

package volcengine

import (
	"context"
	"strings"

	"github.com/pkg/errors"

	api "yunion.io/x/cloudmux/pkg/apis/compute"
	"yunion.io/x/cloudmux/pkg/cloudprovider"
	"yunion.io/x/cloudmux/pkg/multicloud/volcengine"
	"yunion.io/x/jsonutils"
)

type SVolcengineProviderFactory struct {
	cloudprovider.SPublicCloudBaseProviderFactory
}

func (self *SVolcengineProviderFactory) GetId() string {
	return volcengine.CLOUD_PROVIDER_VOLCENGINE
}

func (self *SVolcengineProviderFactory) GetName() string {
	return volcengine.CLOUD_PROVIDER_VOLCENGINE_CN
}

func (self *SVolcengineProviderFactory) ValidateCreateCloudaccountData(ctx context.Context, input cloudprovider.SCloudaccountCredential) (cloudprovider.SCloudaccount, error) {
	output := cloudprovider.SCloudaccount{}
	if len(input.AccessKeyId) == 0 {
		return output, errors.Wrap(cloudprovider.ErrMissingParameter, "access_key_id")
	}
	if len(input.AccessKeySecret) == 0 {
		return output, errors.Wrap(cloudprovider.ErrMissingParameter, "access_key_secret")
	}
	output.Account = input.AccessKeyId
	output.Secret = input.AccessKeySecret
	return output, nil
}

func (f *SVolcengineProviderFactory) ValidateUpdateCloudaccountCredential(ctx context.Context, input cloudprovider.SCloudaccountCredential, cloudaccount string) (cloudprovider.SCloudaccount, error) {
	output := cloudprovider.SCloudaccount{}
	if len(input.AccessKeyId) == 0 {
		return output, errors.Wrap(cloudprovider.ErrMissingParameter, "access_key_id")
	}
	if len(input.AccessKeySecret) == 0 {
		return output, errors.Wrap(cloudprovider.ErrMissingParameter, "access_key_secret")
	}
	output = cloudprovider.SCloudaccount{
		Account: input.AccessKeyId,
		Secret:  input.AccessKeySecret,
	}
	return output, nil
}

func validateClientCloudenv(client *volcengine.SVolcengineClient) error {
	regions := client.GetIRegions()
	if len(regions) == 0 {
		return nil
	}
	return nil
}

func parseAccount(account string) (accessKey string, projectId string) {
	segs := strings.Split(account, "::")
	if len(segs) == 2 {
		accessKey = segs[0]
		projectId = segs[1]
	} else {
		accessKey = account
		projectId = ""
	}

	return
}

func (self *SVolcengineProviderFactory) GetProvider(cfg cloudprovider.ProviderConfig) (cloudprovider.ICloudProvider, error) {
	accessKey, accountId := parseAccount(cfg.Account)
	client, err := volcengine.NewVolcengineClient(
		volcengine.NewVolcengineClientConfig(
			accessKey,
			cfg.Secret,
		).CloudproviderConfig(cfg).AccountId(accountId),
	)
	if err != nil {
		return nil, err
	}

	err = validateClientCloudenv(client)
	if err != nil {
		return nil, errors.Wrap(err, "validateClientCloudenv")
	}

	return &SVolcengineProvider{
		SBaseProvider: cloudprovider.NewBaseProvider(self),
		client:        client,
	}, nil
}

func (self *SVolcengineProviderFactory) GetClientRC(info cloudprovider.SProviderInfo) (map[string]string, error) {
	accessKey, accountId := parseAccount(info.Account)
	return map[string]string{
		"VOLCENGINE_ACCESS_KEY": accessKey,
		"VOLCENGINE_SECRET_KEY": info.Secret,
		"VOLCENGINE_REGION":     volcengine.VOLCENGINE_DEFAULT_REGION,
		"VOLCENGINE_ACCOUNT_ID": accountId,
	}, nil
}

func init() {
	factory := SVolcengineProviderFactory{}
	cloudprovider.RegisterFactory(&factory)
}

type SVolcengineProvider struct {
	cloudprovider.SBaseProvider

	client *volcengine.SVolcengineClient
}

func (self *SVolcengineProvider) GetAccountId() string {
	return self.client.GetAccountId()
}

func (self *SVolcengineProvider) GetSysInfo() (jsonutils.JSONObject, error) {
	regions := self.client.GetIRegions()
	info := jsonutils.NewDict()
	info.Add(jsonutils.NewInt(int64(len(regions))), "region_count")
	info.Add(jsonutils.NewString(volcengine.VOLCENGINE_API_VERSION), "api_version")
	return info, nil
}

func (self *SVolcengineProvider) GetBalance() (*cloudprovider.SBalanceInfo, error) {
	// GetBalance is not currently open
	return &cloudprovider.SBalanceInfo{
		Amount:   0.0,
		Currency: "CNY",
		Status:   api.CLOUD_PROVIDER_HEALTH_NORMAL,
	}, cloudprovider.ErrNotImplemented
}

func (self *SVolcengineProvider) GetBucketCannedAcls(regionId string) []string {
	return nil
}

func (self *SVolcengineProvider) GetCapabilities() []string {
	return self.client.GetCapabilities()
}

func (self *SVolcengineProvider) GetIProjects() ([]cloudprovider.ICloudProject, error) {
	return self.client.GetIProjects()
}

func (self *SVolcengineProvider) GetIRegionById(extId string) (cloudprovider.ICloudRegion, error) {
	return self.client.GetIRegionById(extId)
}

func (self *SVolcengineProvider) GetIRegions() []cloudprovider.ICloudRegion {
	return self.client.GetIRegions()
}

func (self *SVolcengineProvider) GetObjectCannedAcls(regionId string) []string {
	return nil
}

func (self *SVolcengineProvider) GetStorageClasses(regionId string) []string {
	return nil
}

func (self *SVolcengineProvider) GetSubAccounts() ([]cloudprovider.SSubAccount, error) {
	return self.client.GetSubAccounts()
}

func (self *SVolcengineProvider) GetVersion() string {
	return volcengine.VOLCENGINE_API_VERSION
}
