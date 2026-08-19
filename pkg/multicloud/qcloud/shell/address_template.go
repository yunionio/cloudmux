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

	"yunion.io/x/cloudmux/pkg/cloudprovider"
	"yunion.io/x/cloudmux/pkg/multicloud/qcloud"
)

func init() {
	type AddressListOptions struct {
		Id     string `help:"Address template ID"`
		Name   string `help:"Address template name"`
		Limit  int    `help:"page size"`
		Offset int    `help:"page offset"`
	}
	shellutils.R(&AddressListOptions{}, "address-list", "List IP address templates", func(cli *qcloud.SRegion, args *AddressListOptions) error {
		address, total, err := cli.GetClient().AddressList(args.Id, args.Name, args.Offset, args.Limit)
		if err != nil {
			return err
		}
		printList(address, total, args.Offset, args.Limit, []string{})
		return nil
	})

	type AddressIdOptions struct {
		ID string `help:"Address template ID"`
	}
	shellutils.R(&AddressIdOptions{}, "address-show", "Show IP address template", func(cli *qcloud.SRegion, args *AddressIdOptions) error {
		tpl, err := cli.GetClient().GetAddressTemplate(args.ID)
		if err != nil {
			return err
		}
		printObject(tpl)
		return nil
	})

	type AddressCreateOptions struct {
		NAME    string   `help:"Address template name"`
		Address []string `help:"IP, CIDR or IP range entries"`
	}
	shellutils.R(&AddressCreateOptions{}, "address-create", "Create IP address template", func(cli *qcloud.SRegion, args *AddressCreateOptions) error {
		tpl, err := cli.GetClient().CreateAddressTemplate(&cloudprovider.IpSetCreateOptions{
			Name:      args.NAME,
			Addresses: args.Address,
		})
		if err != nil {
			return err
		}
		printObject(tpl)
		return nil
	})

	type AddressUpdateOptions struct {
		ID      string   `help:"Address template ID"`
		Name    string   `help:"Address template name"`
		Address []string `help:"IP, CIDR or IP range entries, replace current entries"`
	}
	shellutils.R(&AddressUpdateOptions{}, "address-update", "Update IP address template", func(cli *qcloud.SRegion, args *AddressUpdateOptions) error {
		err := cli.GetClient().ModifyAddressTemplate(args.ID, &cloudprovider.IpSetUpdateOptions{
			Name:      args.Name,
			Addresses: args.Address,
		})
		if err != nil {
			return err
		}
		tpl, err := cli.GetClient().GetAddressTemplate(args.ID)
		if err != nil {
			return err
		}
		printObject(tpl)
		return nil
	})

	shellutils.R(&AddressIdOptions{}, "address-delete", "Delete IP address template", func(cli *qcloud.SRegion, args *AddressIdOptions) error {
		return cli.GetClient().DeleteAddressTemplate(args.ID)
	})
}
