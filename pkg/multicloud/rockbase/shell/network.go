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

	"yunion.io/x/cloudmux/pkg/multicloud/rockbase"
)

func init() {
	type VSwitchListOptions struct {
		Vpc string `help:"Vpc ID"`
	}
	shellutils.R(&VSwitchListOptions{}, "subnet-list", "List subnets", func(cli *rockbase.SRegion, args *VSwitchListOptions) error {
		vswitches, e := cli.GetNetworks(args.Vpc)
		if e != nil {
			return e
		}
		printList(vswitches, 0, 0, 0, nil)
		return nil
	})

	type SubnetIdOptions struct {
		ID string
	}

	shellutils.R(&SubnetIdOptions{}, "subnet-show", "Show subnet", func(cli *rockbase.SRegion, args *SubnetIdOptions) error {
		net, err := cli.GetNetwork(args.ID)
		if err != nil {
			return err
		}
		printObject(net)
		return nil
	})

	type SubnetCreateOptions struct {
		VPC  string `help:"VPC ID"`
		NAME string
		CIDR string
		Desc string
	}

	shellutils.R(&SubnetCreateOptions{}, "subnet-create", "Create subnet", func(cli *rockbase.SRegion, args *SubnetCreateOptions) error {
		net, err := cli.CreateNetwork(args.VPC, args.NAME, args.CIDR, args.Desc)
		if err != nil {
			return err
		}
		printObject(net)
		return nil
	})

	shellutils.R(&SubnetIdOptions{}, "subnet-delete", "Delete subnet", func(cli *rockbase.SRegion, args *SubnetIdOptions) error {
		return cli.DeleteNetwork(args.ID)
	})
}
