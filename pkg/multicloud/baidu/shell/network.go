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
	"yunion.io/x/cloudmux/pkg/multicloud/baidu"
)

func init() {
	type NetworkListOptions struct {
		VpcId    string
		ZoneName string
	}
	shellutils.R(&NetworkListOptions{}, "network-list", "list networks", func(cli *baidu.SRegion, args *NetworkListOptions) error {
		networks, err := cli.GetNetworks(args.VpcId, args.ZoneName)
		if err != nil {
			return err
		}
		printList(networks)
		return nil
	})

	type NetworkIdOptions struct {
		ID string
	}
	shellutils.R(&NetworkIdOptions{}, "network-show", "show network", func(cli *baidu.SRegion, args *NetworkIdOptions) error {
		network, err := cli.GetNetwork(args.ID)
		if err != nil {
			return err
		}
		printObject(network)
		return nil
	})

	shellutils.R(&NetworkIdOptions{}, "network-delete", "delete network", func(cli *baidu.SRegion, args *NetworkIdOptions) error {
		return cli.DeleteNetwork(args.ID)
	})

	type NetworkCreateOptions struct {
		ZONE string
		VPC  string
		cloudprovider.SNetworkCreateOptions
	}

	shellutils.R(&NetworkCreateOptions{}, "network-create", "create network", func(cli *baidu.SRegion, args *NetworkCreateOptions) error {
		network, err := cli.CreateNetwork(args.ZONE, args.VPC, &args.SNetworkCreateOptions)
		if err != nil {
			return err
		}
		printObject(network)
		return nil
	})

}
