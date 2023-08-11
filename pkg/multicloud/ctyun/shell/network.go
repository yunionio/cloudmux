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
	"yunion.io/x/cloudmux/pkg/multicloud/ctyun"
)

func init() {
	type NetworkListOptions struct {
		VPC string `help:"Vpc ID"`
	}
	shellutils.R(&NetworkListOptions{}, "network-list", "List subnets", func(cli *ctyun.SRegion, args *NetworkListOptions) error {
		networks, err := cli.GetNetwroks(args.VPC)
		if err != nil {
			return err
		}
		printList(networks, 0, 0, 0, nil)
		return nil
	})

	type NetworkCreateOptions struct {
		VPC string
		cloudprovider.SNetworkCreateOptions
	}
	shellutils.R(&NetworkCreateOptions{}, "network-create", "Create subnet", func(cli *ctyun.SRegion, args *NetworkCreateOptions) error {
		vpc, e := cli.CreateNetwork(args.VPC, &args.SNetworkCreateOptions)
		if e != nil {
			return e
		}
		printObject(vpc)
		return nil
	})

	type NetworkIdOptions struct {
		ID string
	}

	shellutils.R(&NetworkIdOptions{}, "network-show", "Show network", func(cli *ctyun.SRegion, args *NetworkIdOptions) error {
		network, err := cli.GetNetwork(args.ID)
		if err != nil {
			return err
		}
		printObject(network)
		return nil
	})

	shellutils.R(&NetworkIdOptions{}, "network-delete", "Delete network", func(cli *ctyun.SRegion, args *NetworkIdOptions) error {
		return cli.DeleteNetwork(args.ID)
	})

}
