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
	"fmt"

	"yunion.io/x/pkg/util/shellutils"

	"yunion.io/x/cloudmux/pkg/multicloud/aws"
)

func init() {
	type TrifficPolicyInstanceShowOptions struct {
		ID string
	}
	shellutils.R(&TrifficPolicyInstanceShowOptions{}, "dns-traffic-policy-instance-show", "Show traffic policy instance", func(cli *aws.SRegion, args *TrifficPolicyInstanceShowOptions) error {
		ret, err := cli.GetClient().GetDnsTrafficPolicyInstance(args.ID)
		if err != nil {
			return err
		}
		printObject(ret)
		return nil
	})

	type TrifficPolicyShowOptions struct {
		ID      string
		VERSION string
	}
	shellutils.R(&TrifficPolicyShowOptions{}, "dns-traffic-policy-show", "Show traffic policy", func(cli *aws.SRegion, args *TrifficPolicyShowOptions) error {
		ret, err := cli.GetClient().GetTrafficPolicy(args.ID, args.VERSION)
		if err != nil {
			return err
		}
		printObject(ret)
		return nil
	})

	shellutils.R(&TrifficPolicyInstanceShowOptions{}, "dns-extra-address-list", "List traffic policy address", func(cli *aws.SRegion, args *TrifficPolicyInstanceShowOptions) error {
		ret, err := cli.GetClient().GetDnsExtraAddresses(args.ID)
		if err != nil {
			return err
		}
		fmt.Println(ret)
		return nil
	})

}
