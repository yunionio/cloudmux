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

	"yunion.io/x/cloudmux/pkg/multicloud/aliyun"
)

func init() {
	type GtmListOptions struct {
	}
	shellutils.R(&GtmListOptions{}, "gtm-instance-list", "List Gtm", func(cli *aliyun.SRegion, args *GtmListOptions) error {
		ret, err := cli.GetClient().DescribeDnsGtmInstances()
		if err != nil {
			return err
		}
		printList(ret, 0, 0, 0, []string{})
		return nil
	})

	type GtmPoolListOptions struct {
		ID string
	}
	shellutils.R(&GtmPoolListOptions{}, "gtm-instance-address-pool-list", "List Gtm address pool", func(cli *aliyun.SRegion, args *GtmPoolListOptions) error {
		ret, err := cli.GetClient().DescribeDnsGtmInstanceAddressPools(args.ID)
		if err != nil {
			return err
		}
		printList(ret, 0, 0, 0, []string{})
		return nil
	})

	type GtmPoolShowOptions struct {
		ID string
	}
	shellutils.R(&GtmPoolShowOptions{}, "gtm-instance-address-pool-show", "Show Gtm address pool", func(cli *aliyun.SRegion, args *GtmPoolShowOptions) error {
		ret, err := cli.GetClient().DescribeDnsGtmInstanceAddressPool(args.ID)
		if err != nil {
			return err
		}
		printObject(ret)
		return nil
	})

}
