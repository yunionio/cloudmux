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
	type NlbListOptions struct {
	}
	shellutils.R(&NlbListOptions{}, "nlb-list", "List NLB instances", func(cli *aliyun.SRegion, args *NlbListOptions) error {
		nlbs, err := cli.GetNlbs()
		if err != nil {
			return err
		}
		printList(nlbs, 0, 0, 0, nil)
		return nil
	})

	type NlbShowOptions struct {
		ID string `help:"NLB instance ID"`
	}
	shellutils.R(&NlbShowOptions{}, "nlb-show", "Show NLB instance", func(cli *aliyun.SRegion, args *NlbShowOptions) error {
		nlb, err := cli.GetNlbDetail(args.ID)
		if err != nil {
			return err
		}
		printObject(nlb)
		return nil
	})

	type NlbListenerListOptions struct {
		NLB_ID string `help:"NLB instance ID"`
	}
	shellutils.R(&NlbListenerListOptions{}, "nlb-listener-list", "List NLB listeners", func(cli *aliyun.SRegion, args *NlbListenerListOptions) error {
		listeners, err := cli.GetNlbListeners(args.NLB_ID)
		if err != nil {
			return err
		}
		printList(listeners, 0, 0, 0, nil)
		return nil
	})

	type NlbListenerShowOptions struct {
		ID string `help:"NLB listener ID"`
	}
	shellutils.R(&NlbListenerShowOptions{}, "nlb-listener-show", "Show NLB listener", func(cli *aliyun.SRegion, args *NlbListenerShowOptions) error {
		listener, err := cli.GetNlbListener(args.ID)
		if err != nil {
			return err
		}
		printObject(listener)
		return nil
	})

	type NlbServerGroupListOptions struct {
	}
	shellutils.R(&NlbServerGroupListOptions{}, "nlb-servergroup-list", "List NLB server groups", func(cli *aliyun.SRegion, args *NlbServerGroupListOptions) error {
		groups, err := cli.GetNlbServerGroups()
		if err != nil {
			return err
		}
		printList(groups, 0, 0, 0, nil)
		return nil
	})

	type NlbServerGroupShowOptions struct {
		ID string `help:"NLB server group ID"`
	}
	shellutils.R(&NlbServerGroupShowOptions{}, "nlb-servergroup-show", "Show NLB server group", func(cli *aliyun.SRegion, args *NlbServerGroupShowOptions) error {
		group, err := cli.GetNlbServerGroup(args.ID)
		if err != nil {
			return err
		}
		printObject(group)
		return nil
	})
}
