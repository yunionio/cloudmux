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
	"yunion.io/x/cloudmux/pkg/cloudprovider"
	"yunion.io/x/cloudmux/pkg/multicloud/ksyun"

	"yunion.io/x/pkg/errors"
	"yunion.io/x/pkg/util/shellutils"
)

func init() {
	type NetworkListOptions struct {
		VpcId     []string
		NetworkId []string
		ZoneName  string
	}
	shellutils.R(&NetworkListOptions{}, "network-list", "list networks", func(cli *ksyun.SRegion, args *NetworkListOptions) error {
		res, err := cli.GetNetworks(args.VpcId, args.NetworkId, args.ZoneName)
		if err != nil {
			return errors.Wrap(err, "GetNetworks")
		}
		printList(res)
		return nil
	})

	type NetworkCreateOptions struct {
		VpcId  string
		ZoneId string
		cloudprovider.SNetworkCreateOptions
	}

	shellutils.R(&NetworkCreateOptions{}, "network-create", "create network", func(cli *ksyun.SRegion, args *NetworkCreateOptions) error {
		network, err := cli.CreateNetwork(args.VpcId, args.ZoneId, &args.SNetworkCreateOptions)
		if err != nil {
			return err
		}
		printObject(network)
		return nil
	})

	type NetworkIdOptions struct {
		ID string
	}

	shellutils.R(&NetworkIdOptions{}, "network-delete", "delete network", func(cli *ksyun.SRegion, args *NetworkIdOptions) error {
		return cli.DeleteNetwork(args.ID)
	})
}
