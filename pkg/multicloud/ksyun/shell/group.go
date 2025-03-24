// Copyright 2023 Yunion
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package shell

import (
	"yunion.io/x/cloudmux/pkg/multicloud/ksyun"
	"yunion.io/x/pkg/util/shellutils"
)

func init() {
	type GroupListOptions struct {
	}
	shellutils.R(&GroupListOptions{}, "cloud-group-list", "List groups", func(cli *ksyun.SRegion, args *GroupListOptions) error {
		ret, err := cli.GetClient().ListGroups()
		if err != nil {
			return err
		}
		printList(ret)
		return nil
	})

	type GroupCreateOptions struct {
		NAME string
		Desc string
	}

	shellutils.R(&GroupCreateOptions{}, "cloud-group-create", "Create group", func(cli *ksyun.SRegion, args *GroupCreateOptions) error {
		ret, err := cli.GetClient().CreateGroup(args.NAME, args.Desc)
		if err != nil {
			return err
		}
		printObject(ret)
		return nil
	})

	type GroupNameOptions struct {
		NAME string
	}

	shellutils.R(&GroupNameOptions{}, "cloud-group-show", "Show group", func(cli *ksyun.SRegion, args *GroupNameOptions) error {
		ret, err := cli.GetClient().GetGroup(args.NAME)
		if err != nil {
			return err
		}
		printObject(ret)
		return nil
	})

	shellutils.R(&GroupNameOptions{}, "cloud-group-delete", "Delete group", func(cli *ksyun.SRegion, args *GroupNameOptions) error {
		return cli.GetClient().DeleteGroup(args.NAME)
	})

	shellutils.R(&GroupNameOptions{}, "cloud-group-policy-list", "List group policies", func(cli *ksyun.SRegion, args *GroupNameOptions) error {
		ret, err := cli.GetClient().ListGroupPolicies(args.NAME)
		if err != nil {
			return err
		}
		printList(ret)
		return nil
	})

}
