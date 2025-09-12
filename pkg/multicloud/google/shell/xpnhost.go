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
	type XpnHostListOptions struct {
		MaxResults int
		PageToken  string
	}
	shellutils.R(&XpnHostListOptions{}, "xpn-host-list", "List xpn hosts", func(cli *google.SRegion, args *XpnHostListOptions) error {
		hosts, err := cli.GetClient().GetXpnHosts()
		if err != nil {
			return err
		}
		printList(hosts, 0, 0, 0, nil)
		return nil
	})

	type XpnNetworkListOptions struct {
		PROJECT string
	}
	shellutils.R(&XpnNetworkListOptions{}, "xpn-network-list", "List xpn networks", func(cli *google.SRegion, args *XpnNetworkListOptions) error {
		networks, err := cli.GetClient().GetXpnNetworks(args.PROJECT)
		if err != nil {
			return err
		}
		printList(networks, 0, 0, 0, nil)
		return nil
	})

	type XpnResourceListOptions struct {
		PROJECT string
	}
	shellutils.R(&XpnResourceListOptions{}, "xpn-resource-list", "List xpn resources", func(cli *google.SRegion, args *XpnResourceListOptions) error {
		resources, err := cli.GetClient().GetXpnResources(args.PROJECT)
		if err != nil {
			return err
		}
		printList(resources, 0, 0, 0, nil)
		return nil
	})
}
