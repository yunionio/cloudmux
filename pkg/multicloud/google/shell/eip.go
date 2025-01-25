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

	"yunion.io/x/cloudmux/pkg/multicloud/google"
)

func init() {
	type EipListOptions struct {
		Address    string
		MaxResults int
		PageToken  string
	}
	shellutils.R(&EipListOptions{}, "eip-list", "List eips", func(cli *google.SRegion, args *EipListOptions) error {
		eips, err := cli.GetEips(args.Address, args.MaxResults, args.PageToken)
		if err != nil {
			return err
		}
		printList(eips, 0, 0, 0, nil)
		return nil
	})

	shellutils.R(&EipListOptions{}, "global-eip-list", "List global eips", func(cli *google.SRegion, args *EipListOptions) error {
		eips, err := cli.GetClient().GetGlobalRegion().GetEips(args.Address)
		if err != nil {
			return err
		}
		printList(eips, 0, 0, 0, nil)
		return nil
	})

	type EipIdOptions struct {
		ID string `help:"numeric id of eip"`
	}
	shellutils.R(&EipIdOptions{}, "eip-show", "Show eip", func(cli *google.SRegion, args *EipIdOptions) error {
		eip, err := cli.GetEip(args.ID)
		if err != nil {
			return err
		}
		printObject(eip)
		return nil
	})

	type EipIdOptions2 struct {
		LINK string `help:"self link of eip"`
	}
	shellutils.R(&EipIdOptions2{}, "eip-delete", "Delete eip", func(cli *google.SRegion, args *EipIdOptions2) error {
		return cli.Delete(args.LINK)
	})

	type EipCreateOptions struct {
		NAME string
		Desc string
	}
	shellutils.R(&EipCreateOptions{}, "eip-create", "Create eip", func(cli *google.SRegion, args *EipCreateOptions) error {
		eip, err := cli.CreateEip(args.NAME, args.Desc)
		if err != nil {
			return err
		}
		printObject(eip)
		return nil
	})

}
