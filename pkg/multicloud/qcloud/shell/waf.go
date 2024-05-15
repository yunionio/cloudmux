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
	"yunion.io/x/pkg/util/shellutils"

	"yunion.io/x/cloudmux/pkg/multicloud/qcloud"
)

func init() {
	type WafListOptions struct {
	}
	shellutils.R(&WafListOptions{}, "waf-list", "List wafs", func(cli *qcloud.SRegion, args *WafListOptions) error {
		ret, err := cli.GetWafInstances()
		if err != nil {
			return err
		}
		printList(ret, 0, 0, 0, []string{})
		return nil
	})

	type WafShowOptions struct {
		DOMAIN      string
		DOMAIN_ID   string
		INSTANCE_ID string
	}

	shellutils.R(&WafShowOptions{}, "waf-show", "Show waf", func(cli *qcloud.SRegion, args *WafShowOptions) error {
		ret, err := cli.GetWafInstance(args.DOMAIN, args.DOMAIN_ID, args.INSTANCE_ID)
		if err != nil {
			return err
		}
		printObject(ret)
		return nil
	})

}
