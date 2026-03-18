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
	type VNetworkListOptions struct {
		RouteId string `help:"VPC ID or router ID"`
	}
	shellutils.R(&VNetworkListOptions{}, "network-list", "List networks", func(cli *ecloud.SRegion, args *VNetworkListOptions) error {
		networks, err := cli.GetNetworks(args.RouteId, "")
		if err != nil {
			return err
		}
		printList(networks)
		return nil
	})

	type NetworkShowOptions struct {
		ID string `help:"Subnet ID"`
	}
	shellutils.R(&NetworkShowOptions{}, "network-show", "Show network detail", func(cli *ecloud.SRegion, args *NetworkShowOptions) error {
		// VPCId 当前未被 OpenAPI 使用，这里仅保留参数以兼容调用习惯
		net, err := cli.GetNetwork(args.ID)
		if err != nil {
			return err
		}
		printObject(net)
		return nil
	})

	type NetworkCreateOptions struct {
		VpcId string `help:"VPC ID or router ID"`
		Name  string `help:"Network name (5-22 chars, letter first)"`
		Cidr  string `help:"CIDR, e.g. 192.168.1.0/24" default:"192.168.0.0/24"`
	}
	shellutils.R(&NetworkCreateOptions{}, "network-create", "Create network", func(cli *ecloud.SRegion, args *NetworkCreateOptions) error {
		ivpc, err := cli.GetIVpcById(args.VpcId)
		if err != nil {
			return err
		}
		vpc := ivpc.(*ecloud.SVpc)
		iwires, err := vpc.GetIWires()
		if err != nil || len(iwires) == 0 {
			return err
		}
		wire := iwires[0].(*ecloud.SWire)
		inet, err := wire.CreateINetwork(&cloudprovider.SNetworkCreateOptions{
			Name: args.Name,
			Cidr: args.Cidr,
		})
		if err != nil {
			return err
		}
		printObject(inet)
		return nil
	})

	type NetworkDeleteOptions struct {
		NetworkId string `help:"Network ID"`
	}
	shellutils.R(&NetworkDeleteOptions{}, "network-delete", "Delete network", func(cli *ecloud.SRegion, args *NetworkDeleteOptions) error {
		return cli.DeleteNetwork(args.NetworkId)
	})
}
