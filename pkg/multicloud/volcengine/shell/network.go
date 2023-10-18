// Copyright 2023 Yunion
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
	"time"

	"yunion.io/x/cloudmux/pkg/multicloud/volcengine"
	"yunion.io/x/pkg/util/shellutils"
)

func init() {
	type NetworkListOptions struct {
		Limit  int `help:"page size"`
		Offset int `help:"page offset"`
	}
	shellutils.R(&NetworkListOptions{}, "network-list", "List networkes", func(cli *volcengine.SRegion, args *NetworkListOptions) error {
		networks, total, err := cli.GetSubnets(nil, "", "", args.Offset, args.Limit)
		if err != nil {
			return err
		}
		printList(networks, total, args.Offset, args.Limit, nil)
		return nil
	})

	type NetworkShowOptions struct {
		ID string `help:"show network details"`
	}
	shellutils.R(&NetworkShowOptions{}, "network-show", "Show network details", func(cli *volcengine.SRegion, args *NetworkShowOptions) error {
		network, e := cli.GetSubnetAttributes(args.ID)
		if e != nil {
			return e
		}
		printObject(network)
		return nil
	})

	shellutils.R(&NetworkShowOptions{}, "network-delete", "Delete subnet", func(cli *volcengine.SRegion, args *NetworkShowOptions) error {
		e := cli.DeleteSubnet(args.ID)
		if e != nil {
			return e
		}
		return nil
	})

	type NetworkCreateOption struct {
		ZoneId string
		VpcId  string
		Name   string
		Desc   string
		CIDR   string
	}

	shellutils.R(&NetworkCreateOption{}, "network-create", "create network", func(cli *volcengine.SRegion, args *NetworkCreateOption) error {
		networkId, err := cli.CreateSubnet(args.ZoneId, args.VpcId, args.Name, args.CIDR, args.Desc)
		if err != nil {
			return err
		}
		time.Sleep(time.Second * 1)
		network, _, err := cli.GetSubnets([]string{networkId}, "", "", 1, 1)
		if err != nil {
			return err
		}
		printObject(network)
		return nil
	})
}
