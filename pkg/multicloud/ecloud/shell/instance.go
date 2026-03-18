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
	"context"
	"fmt"

	"yunion.io/x/pkg/util/shellutils"

	"yunion.io/x/cloudmux/pkg/multicloud/ecloud"
)

func init() {
	type InstanceListOptions struct {
		ZoneId   string
		ServerId string
	}
	shellutils.R(&InstanceListOptions{}, "instance-list", "List intances", func(cli *ecloud.SRegion, args *InstanceListOptions) error {
		instances, err := cli.GetInstances(args.ZoneId, args.ServerId)
		if err != nil {
			return err
		}
		printList(instances)
		return nil
	})
	type InstanceShowOptions struct {
		ID string
	}
	shellutils.R(&InstanceShowOptions{}, "instance-show", "Show intances", func(cli *ecloud.SRegion, args *InstanceShowOptions) error {
		instance, err := cli.GetInstance(args.ID)
		if err != nil {
			return err
		}
		printObject(instance)
		return nil
	})
	shellutils.R(&InstanceShowOptions{}, "instance-nic-list", "List intance nics", func(cli *ecloud.SRegion, args *InstanceShowOptions) error {
		nics, err := cli.GetInstanceNics(context.Background(), args.ID)
		if err != nil {
			return err
		}
		printList(nics)
		return nil
	})
	shellutils.R(&InstanceShowOptions{}, "instance-disk-list", "List intance disks", func(cli *ecloud.SRegion, args *InstanceShowOptions) error {
		disks, err := cli.GetDataDisks(args.ID)
		if err != nil {
			return err
		}
		printList(disks)
		return nil
	})

	shellutils.R(&InstanceShowOptions{}, "instance-vnc-url", "Get instance VNC url", func(cli *ecloud.SRegion, args *InstanceShowOptions) error {
		url, err := cli.GetInstanceVNCUrl(args.ID)
		if err != nil {
			return err
		}
		fmt.Println(url)
		return nil
	})

	shellutils.R(&InstanceShowOptions{}, "instance-start", "Start instance", func(cli *ecloud.SRegion, args *InstanceShowOptions) error {
		return cli.StartInstance(context.Background(), args.ID)
	})

	shellutils.R(&InstanceShowOptions{}, "instance-stop", "Stop instance", func(cli *ecloud.SRegion, args *InstanceShowOptions) error {
		return cli.StopInstance(context.Background(), args.ID)
	})

	type InstanceDeleteOptions struct {
		ID                  string `help:"Instance ID"`
		DeletePublicNetwork bool   `help:"Delete associated public network (EIP)" default:"false"`
		DeleteDataVolumes   bool   `help:"Delete associated data volumes" default:"false"`
	}
	shellutils.R(&InstanceDeleteOptions{}, "instance-delete", "Delete instance", func(cli *ecloud.SRegion, args *InstanceDeleteOptions) error {
		return cli.DeleteInstance(context.Background(), args.ID, args.DeletePublicNetwork, args.DeleteDataVolumes)
	})
}
