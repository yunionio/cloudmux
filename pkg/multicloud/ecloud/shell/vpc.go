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
	"yunion.io/x/cloudmux/pkg/multicloud/ecloud"
)

func init() {
	type VpcListOptions struct {
	}
	shellutils.R(&VpcListOptions{}, "vpc-list", "List vpcs", func(cli *ecloud.SRegion, args *VpcListOptions) error {
		vpcs, err := cli.GetVpcs()
		if err != nil {
			return err
		}
		printList(vpcs)
		return nil
	})

	type VpcCreateOptions struct {
		Name string `help:"VPC name"`
		Desc string `help:"Description"`
		CIDR string `help:"CIDR block, e.g. 192.168.0.0/16" optional:"true"`
	}
	shellutils.R(&VpcCreateOptions{}, "vpc-create", "Create vpc", func(cli *ecloud.SRegion, args *VpcCreateOptions) error {
		opts := &cloudprovider.VpcCreateOptions{
			NAME: args.Name,
			Desc: args.Desc,
			CIDR: args.CIDR,
		}
		vpc, err := cli.CreateIVpc(opts)
		if err != nil {
			return err
		}
		printObject(vpc)
		return nil
	})

	type VpcIdOptions struct {
		ID string `help:"VPC id or router id"`
	}
	shellutils.R(&VpcIdOptions{}, "vpc-show", "Show vpc detail", func(cli *ecloud.SRegion, args *VpcIdOptions) error {
		ivpc, err := cli.GetVpc(args.ID)
		if err != nil {
			return err
		}
		printObject(ivpc)
		return nil
	})

	shellutils.R(&VpcIdOptions{}, "vpc-delete", "Delete vpc", func(cli *ecloud.SRegion, args *VpcIdOptions) error {
		return cli.DeleteVpc(args.ID)
	})
}
