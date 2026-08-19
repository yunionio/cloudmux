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
	"yunion.io/x/cloudmux/pkg/multicloud/aliyun"
)

func init() {
	type PrefixListListOptions struct {
	}
	shellutils.R(&PrefixListListOptions{}, "prefix-list-list", "List VPC prefix lists", func(cli *aliyun.SRegion, args *PrefixListListOptions) error {
		lists, err := cli.GetPrefixLists()
		if err != nil {
			return err
		}
		printList(lists, 0, 0, 0, []string{})
		return nil
	})

	type PrefixListIdOptions struct {
		ID string `help:"Prefix list ID"`
	}
	shellutils.R(&PrefixListIdOptions{}, "prefix-list-show", "Show VPC prefix list", func(cli *aliyun.SRegion, args *PrefixListIdOptions) error {
		pl, err := cli.GetPrefixList(args.ID)
		if err != nil {
			return err
		}
		printObject(pl)
		return nil
	})

	type PrefixListCreateOptions struct {
		NAME      string   `help:"Prefix list name"`
		Desc      string   `help:"Prefix list description"`
		IpVersion string   `help:"IP version" choices:"IPv4|IPv6" default:"IPv4"`
		Cidr      []string `help:"CIDR entries"`
	}
	shellutils.R(&PrefixListCreateOptions{}, "prefix-list-create", "Create VPC prefix list", func(cli *aliyun.SRegion, args *PrefixListCreateOptions) error {
		ipSetType := cloudprovider.IpSetTypeIpv4CidrList
		if args.IpVersion == "IPv6" {
			ipSetType = cloudprovider.IpSetTypeIpv6CidrList
		}
		pl, err := cli.CreatePrefixList(&cloudprovider.IpSetCreateOptions{
			Name:      args.NAME,
			Desc:      args.Desc,
			IpSetType: ipSetType,
			Addresses: args.Cidr,
		})
		if err != nil {
			return err
		}
		printObject(pl)
		return nil
	})

	type PrefixListUpdateOptions struct {
		ID   string   `help:"Prefix list ID"`
		Name string   `help:"Prefix list name"`
		Cidr []string `help:"CIDR entries, replace current entries"`
	}
	shellutils.R(&PrefixListUpdateOptions{}, "prefix-list-update", "Update VPC prefix list", func(cli *aliyun.SRegion, args *PrefixListUpdateOptions) error {
		err := cli.ModifyPrefixList(args.ID, &cloudprovider.IpSetUpdateOptions{
			Name:      args.Name,
			Addresses: args.Cidr,
		})
		if err != nil {
			return err
		}
		pl, err := cli.GetPrefixList(args.ID)
		if err != nil {
			return err
		}
		printObject(pl)
		return nil
	})

	shellutils.R(&PrefixListIdOptions{}, "prefix-list-delete", "Delete VPC prefix list", func(cli *aliyun.SRegion, args *PrefixListIdOptions) error {
		return cli.DeletePrefixList(args.ID)
	})
}
