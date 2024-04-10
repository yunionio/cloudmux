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

package shell

import (
	"yunion.io/x/cloudmux/pkg/cloudprovider"
	"yunion.io/x/pkg/errors"
)

func init() {
	r := EmptyOptionProviderR("cloudprovider")

	r.Run("system-info", "Get cloudprovider system info", func(cli cloudprovider.ICloudProvider) (any, error) {
		return cli.GetSysInfo()
	})

	r.Run("api-version", "Get cloudprovider api version", func(cli cloudprovider.ICloudProvider) (any, error) {
		version := cli.GetVersion()
		return map[string]string{
			"api_version": version,
		}, nil
	})

	r.Run("iam-login-url", "Get IAM login URL", func(cli cloudprovider.ICloudProvider) (any, error) {
		url := cli.GetIamLoginUrl()
		return map[string]string{
			"url": url,
		}, nil
	})

	r.List("region-list", "List regions of a cloudprovider", func(cli cloudprovider.ICloudProvider) (any, error) {
		return cli.GetIRegions()
	})

	r.List("subaccount-list", "List subaccounts of a cloudprovider", func(cli cloudprovider.ICloudProvider) (any, error) {
		accounts, err := cli.GetSubAccounts()
		if err != nil {
			return nil, errors.Wrap(err, "GetSubAccounts")
		}
		return accounts, nil
	})
}
