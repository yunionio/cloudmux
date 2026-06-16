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

	"yunion.io/x/cloudmux/pkg/multicloud/aws"
)

func init() {
	type DedicatedHostListOptions struct {
		Zone string   `help:"Zone ID"`
		Id   []string `help:"Host IDs"`
	}
	shellutils.R(&DedicatedHostListOptions{}, "dedicated-host-list", "List dedicated hosts", func(cli *aws.SRegion, args *DedicatedHostListOptions) error {
		hosts, err := cli.GetDedicatedHosts(args.Zone, args.Id...)
		if err != nil {
			return err
		}
		printList(hosts, 0, 0, 0, []string{})
		return nil
	})

	type DedicatedHostShowOptions struct {
		ID string `help:"Dedicated host ID"`
	}
	shellutils.R(&DedicatedHostShowOptions{}, "dedicated-host-show", "Show dedicated host details", func(cli *aws.SRegion, args *DedicatedHostShowOptions) error {
		host, err := cli.GetDedicatedHost(args.ID)
		if err != nil {
			return err
		}
		printObject(host)
		return nil
	})

	type DedicatedHostCreateOptions struct {
		Zone          string `help:"Availability zone"`
		InstanceType  string `help:"Dedicated host instance type"`
		Name          string `help:"Host name tag"`
		Quantity      int    `help:"Number of hosts to allocate" default:"1"`
		AutoPlacement bool   `help:"Allow automatic instance placement" default:"true"`
	}
	shellutils.R(&DedicatedHostCreateOptions{}, "dedicated-host-create", "Allocate a dedicated host", func(cli *aws.SRegion, args *DedicatedHostCreateOptions) error {
		hosts, err := cli.AllocateDedicatedHost(args.Zone, args.InstanceType, args.Name, args.Quantity, args.AutoPlacement)
		if err != nil {
			return err
		}
		printList(hosts, 0, 0, 0, []string{})
		return nil
	})

	type DedicatedHostDeleteOptions struct {
		ID string `help:"Dedicated host ID"`
	}
	shellutils.R(&DedicatedHostDeleteOptions{}, "dedicated-host-delete", "Release a dedicated host", func(cli *aws.SRegion, args *DedicatedHostDeleteOptions) error {
		return cli.ReleaseDedicatedHost(args.ID)
	})
}
