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
	"yunion.io/x/pkg/errors"
	"yunion.io/x/pkg/util/shellutils"

	"yunion.io/x/cloudmux/pkg/multicloud/aliyun"
)

func init() {
	type WafShowOptions struct {
	}
	shellutils.R(&WafShowOptions{}, "waf-v2-instance-show", "Show waf instance", func(cli *aliyun.SRegion, args *WafShowOptions) error {
		waf, err := cli.DescribeWafInstance()
		if err != nil {
			return err
		}
		printObject(waf)
		return nil
	})

	type WafIdOptions struct {
		ID string
	}

	shellutils.R(&WafIdOptions{}, "waf-v2-domain-list", "List waf instance domains", func(cli *aliyun.SRegion, args *WafIdOptions) error {
		domains, err := cli.DescribeWafDomains(args.ID)
		if err != nil {
			return errors.Wrapf(err, "DescribeDomainNames")
		}
		printList(domains, 0, 0, 0, nil)
		return nil
	})

	type WafDomainIdOptions struct {
		ID     string
		DOMAIN string
	}

	shellutils.R(&WafDomainIdOptions{}, "waf-v2-domain-show", "Show waf domain", func(cli *aliyun.SRegion, args *WafDomainIdOptions) error {
		domain, err := cli.DescribeDomainV2(args.ID, args.DOMAIN)
		if err != nil {
			return err
		}
		printObject(domain)
		return nil
	})
}
