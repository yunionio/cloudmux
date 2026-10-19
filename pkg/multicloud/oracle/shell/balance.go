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

	"yunion.io/x/cloudmux/pkg/multicloud/oracle"
)

func init() {
	type BalanceShowOptions struct {
	}
	shellutils.R(&BalanceShowOptions{}, "balance-show", "Show account balance from subscribed service availableAmount", func(cli *oracle.SRegion, args *BalanceShowOptions) error {
		balance, err := cli.GetClient().GetBalance()
		if err != nil {
			return err
		}
		printObject(balance)
		return nil
	})

	type SubscribedServiceListOptions struct {
	}
	shellutils.R(&SubscribedServiceListOptions{}, "subscribed-service-list", "List subscribed services", func(cli *oracle.SRegion, args *SubscribedServiceListOptions) error {
		subs, err := cli.GetClient().GetOrganizationSubscriptions()
		if err != nil {
			return err
		}
		ret := []oracle.SSubscribedService{}
		for _, sub := range subs {
			services, err := cli.GetClient().GetSubscribedServices(sub.Id)
			if err != nil {
				return err
			}
			ret = append(ret, services...)
		}
		printList(ret, 0, 0, 0, []string{})
		return nil
	})

	type CommitmentListOptions struct {
	}
	shellutils.R(&CommitmentListOptions{}, "commitment-list", "List commitments under ACTIVE subscribed services", func(cli *oracle.SRegion, args *CommitmentListOptions) error {
		commitments, err := cli.GetClient().GetCommitmentsBalance()
		if err != nil {
			return err
		}
		printList(commitments, 0, 0, 0, []string{})
		return nil
	})
}
