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
	"yunion.io/x/cloudmux/pkg/cloudprovider"
	"yunion.io/x/cloudmux/pkg/multicloud/test"
	"yunion.io/x/cloudmux/pkg/multicloud/volcengine"
	"yunion.io/x/pkg/util/shellutils"
)

func init() {
	test.TestShell()
	type VpcListOptions struct {
		Limit  int      `help:"page size"`
		Offset int      `help:"page offset"`
		VpcIds []string `help:"Vpc IDs"`
	}
	shellutils.R(&VpcListOptions{}, "vpc-list", "List vpc", func(cli *volcengine.SRegion, args *VpcListOptions) error {
		vpcs, total, err := cli.GetVpcs(args.VpcIds, args.Offset, args.Limit)
		if err != nil {
			return err
		}
		printList(vpcs, total, args.Offset, args.Limit, nil)
		return nil
	})

	type VpcCreateOptions struct {
		Name string
		Desc string
		CIDR string
	}
	shellutils.R(&VpcCreateOptions{}, "vpc-create", "Create vpc", func(cli *volcengine.SRegion, args *VpcCreateOptions) error {
		opts := cloudprovider.VpcCreateOptions{
			NAME: args.Name,
			CIDR: args.CIDR,
			Desc: args.Desc,
		}
		vpc, err := cli.CreateIVpc(&opts)
		if err != nil {
			return err
		}
		printObject(vpc)
		return nil
	})

	type VpcIdOptions struct {
		ID string
	}
	shellutils.R(&VpcIdOptions{}, "vpc-show", "show vpc", func(cli *volcengine.SRegion, args *VpcIdOptions) error {
		vpc, err := cli.GetIVpcById(args.ID)
		if err != nil {
			return err
		}
		printObject(vpc)
		return nil
	})

	shellutils.R(&VpcIdOptions{}, "vpc-delete", "delete vpc", func(cli *volcengine.SRegion, args *VpcIdOptions) error {
		return cli.DeleteVpc(args.ID)
	})
}
