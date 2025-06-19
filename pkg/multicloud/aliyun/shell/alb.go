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
	type AlbListOptions struct {
	}
	shellutils.R(&AlbListOptions{}, "alb-list", "List ALB instances", func(cli *aliyun.SRegion, args *AlbListOptions) error {
		albs, err := cli.GetAlbs()
		if err != nil {
			return err
		}
		printList(albs, 0, 0, 0, nil)
		return nil
	})

	type AlbShowOptions struct {
		ID string `help:"ALB instance ID"`
	}
	shellutils.R(&AlbShowOptions{}, "alb-show", "Show ALB instance", func(cli *aliyun.SRegion, args *AlbShowOptions) error {
		alb, err := cli.GetAlbDetail(args.ID)
		if err != nil {
			return err
		}
		printObject(alb)
		return nil
	})

	type AlbListenerListOptions struct {
		ALB_ID string `help:"ALB instance ID"`
	}
	shellutils.R(&AlbListenerListOptions{}, "alb-listener-list", "List ALB listeners", func(cli *aliyun.SRegion, args *AlbListenerListOptions) error {
		listeners, err := cli.GetAlbListeners(args.ALB_ID)
		if err != nil {
			return err
		}
		printList(listeners, 0, 0, 0, nil)
		return nil
	})

	type AlbListenerShowOptions struct {
		ID string `help:"ALB listener ID"`
	}
	shellutils.R(&AlbListenerShowOptions{}, "alb-listener-show", "Show ALB listener", func(cli *aliyun.SRegion, args *AlbListenerShowOptions) error {
		listener, err := cli.GetAlbListener(args.ID)
		if err != nil {
			return err
		}
		printObject(listener)
		return nil
	})

	type AlbServerGroupListOptions struct {
		ALB_ID string `help:"ALB instance ID"`
	}
	shellutils.R(&AlbServerGroupListOptions{}, "alb-servergroup-list", "List ALB server groups", func(cli *aliyun.SRegion, args *AlbServerGroupListOptions) error {
		groups, err := cli.GetAlbServerGroups(args.ALB_ID)
		if err != nil {
			return err
		}
		printList(groups, 0, 0, 0, nil)
		return nil
	})

	type AlbServerGroupShowOptions struct {
		ID string `help:"ALB server group ID"`
	}
	shellutils.R(&AlbServerGroupShowOptions{}, "alb-servergroup-show", "Show ALB server group", func(cli *aliyun.SRegion, args *AlbServerGroupShowOptions) error {
		group, err := cli.GetAlbServerGroup(args.ID)
		if err != nil {
			return err
		}
		printObject(group)
		return nil
	})

	type AlbRuleListOptions struct {
		LISTENER_ID string `help:"ALB listener ID"`
	}
	shellutils.R(&AlbRuleListOptions{}, "alb-rule-list", "List ALB rules", func(cli *aliyun.SRegion, args *AlbRuleListOptions) error {
		rules, err := cli.GetAlbRules(args.LISTENER_ID)
		if err != nil {
			return err
		}
		printList(rules, 0, 0, 0, nil)
		return nil
	})

	type AlbRuleShowOptions struct {
		ID string `help:"ALB rule ID"`
	}
	shellutils.R(&AlbRuleShowOptions{}, "alb-rule-show", "Show ALB rule", func(cli *aliyun.SRegion, args *AlbRuleShowOptions) error {
		rule, err := cli.GetAlbRule(args.ID)
		if err != nil {
			return err
		}
		printObject(rule)
		return nil
	})
}
