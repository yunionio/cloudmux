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

package provider

import (
	"context"

	"yunion.io/x/jsonutils"
	"yunion.io/x/pkg/errors"
	"yunion.io/x/pkg/util/regutils"

	api "yunion.io/x/cloudmux/pkg/apis/compute"
	"yunion.io/x/cloudmux/pkg/cloudprovider"
	"yunion.io/x/cloudmux/pkg/multicloud/harvester"
)

type SHarvesterProviderFactory struct {
	cloudprovider.SPrivateCloudBaseProviderFactory
}

func (self *SHarvesterProviderFactory) GetId() string {
	return harvester.CLOUD_PROVIDER_HARVESTER
}

func (self *SHarvesterProviderFactory) GetName() string {
	return harvester.CLOUD_PROVIDER_HARVESTER
}

func (self *SHarvesterProviderFactory) GetSupportedBrands() []string {
	return []string{}
}

func (self *SHarvesterProviderFactory) ValidateCreateCloudaccountData(ctx context.Context, input cloudprovider.SCloudaccountCredential) (cloudprovider.SCloudaccount, error) {
	output := cloudprovider.SCloudaccount{}
	if len(input.AccessKeyId) == 0 {
		return output, errors.Wrap(cloudprovider.ErrMissingParameter, "access_key_id")
	}
	if len(input.AccessKeySecret) == 0 {
		return output, errors.Wrap(cloudprovider.ErrMissingParameter, "access_key_secret")
	}
	if len(input.Host) == 0 {
		return output, errors.Wrap(cloudprovider.ErrMissingParameter, "host")
	}
	if !regutils.MatchIPAddr(input.Host) && !regutils.MatchDomainName(input.Host) {
		return output, errors.Wrap(cloudprovider.ErrInputParameter, "host should be ip or domain name")
	}
	output.AccessUrl = input.Host
	output.Account = input.AccessKeyId
	output.Secret = input.AccessKeySecret
	return output, nil
}

func (self *SHarvesterProviderFactory) ValidateUpdateCloudaccountCredential(ctx context.Context, input cloudprovider.SCloudaccountCredential, cloudaccount string) (cloudprovider.SCloudaccount, error) {
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

func (self *SHarvesterProviderFactory) GetProvider(cfg cloudprovider.ProviderConfig) (cloudprovider.ICloudProvider, error) {
	client, err := harvester.NewHarvesterClient(
		harvester.NewHarvesterClientConfig(
			cfg.URL, cfg.Account, cfg.Secret,
		).CloudproviderConfig(cfg),
	)
	if err != nil {
		return nil, err
	}
	return &SHarvesterProvider{
		SBaseProvider: cloudprovider.NewBaseProvider(self),
		client:        client,
	}, nil
}

func (self *SHarvesterProviderFactory) GetClientRC(info cloudprovider.SProviderInfo) (map[string]string, error) {
	return map[string]string{
		"HARVESTER_HOST":              info.Url,
		"HARVESTER_ACCESS_KEY_ID":     info.Account,
		"HARVESTER_ACCESS_KEY_SECRET": info.Secret,
	}, nil
}

func init() {
	factory := SHarvesterProviderFactory{}
	cloudprovider.RegisterFactory(&factory)
}

type SHarvesterProvider struct {
	cloudprovider.SBaseProvider
	client *harvester.SHarvesterClient
}

func (self *SHarvesterProvider) GetVersion() string {
	return ""
}

func (self *SHarvesterProvider) GetSysInfo() (jsonutils.JSONObject, error) {
	return jsonutils.NewDict(), nil
}

func (self *SHarvesterProvider) GetSubAccounts() ([]cloudprovider.SSubAccount, error) {
	return self.client.GetSubAccounts()
}

func (self *SHarvesterProvider) GetAccountId() string {
	return ""
}

func (self *SHarvesterProvider) GetIRegions() ([]cloudprovider.ICloudRegion, error) {
	return self.client.GetIRegions()
}

func (self *SHarvesterProvider) GetIRegionById(extId string) (cloudprovider.ICloudRegion, error) {
	return self.client.GetIRegionById(extId)
}

func (self *SHarvesterProvider) GetBalance() (*cloudprovider.SBalanceInfo, error) {
	return &cloudprovider.SBalanceInfo{
		Amount:   0.0,
		Currency: "CNY",
		Status:   api.CLOUD_PROVIDER_HEALTH_NORMAL,
	}, cloudprovider.ErrNotSupported
}

func (self *SHarvesterProvider) GetCloudRegionExternalIdPrefix() string {
	return self.client.GetCloudRegionExternalIdPrefix()
}

func (self *SHarvesterProvider) GetIProjects() ([]cloudprovider.ICloudProject, error) {
	return self.client.GetIProjects()
}

func (self *SHarvesterProvider) GetStorageClasses(regionId string) []string {
	return nil
}

func (self *SHarvesterProvider) GetBucketCannedAcls(regionId string) []string {
	return nil
}

func (self *SHarvesterProvider) GetObjectCannedAcls(regionId string) []string {
	return nil
}

func (self *SHarvesterProvider) GetCapabilities() []string {
	return self.client.GetCapabilities()
}

func (self *SHarvesterProvider) GetMetrics(opts *cloudprovider.MetricListOptions) ([]cloudprovider.MetricValues, error) {
	return nil, cloudprovider.ErrNotImplemented
	//return self.client.GetMetrics(opts)
}
