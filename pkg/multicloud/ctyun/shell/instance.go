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
	"fmt"

	"yunion.io/x/pkg/util/shellutils"

	"yunion.io/x/cloudmux/pkg/multicloud/ctyun"
)

func init() {
	type InstanceListOptions struct {
		ZoneId string
		Ids    []string `help:"ID of instance to show"`
	}
	shellutils.R(&InstanceListOptions{}, "instance-list", "List intances", func(cli *ctyun.SRegion, args *InstanceListOptions) error {
		instances, err := cli.GetInstances(args.ZoneId, args.Ids)
		if err != nil {
			return err
		}
		printList(instances, 0, 0, 0, []string{})
		return nil
	})

	type SInstanceIdOptions struct {
		ID string
	}

	shellutils.R(&SInstanceIdOptions{}, "instance-start", "Start intance", func(cli *ctyun.SRegion, args *SInstanceIdOptions) error {
		return cli.StartVM(args.ID)
	})

	shellutils.R(&SInstanceIdOptions{}, "instance-stop", "Stop intance", func(cli *ctyun.SRegion, args *SInstanceIdOptions) error {
		return cli.StopVM(args.ID)
	})

	shellutils.R(&SInstanceIdOptions{}, "instance-delete", "Delete intance", func(cli *ctyun.SRegion, args *SInstanceIdOptions) error {
		return cli.DeleteVM(args.ID)
	})

	shellutils.R(&SInstanceIdOptions{}, "instance-vnc", "Show intance vnc", func(cli *ctyun.SRegion, args *SInstanceIdOptions) error {
		url, err := cli.GetInstanceVnc(args.ID)
		if err != nil {
			return err
		}
		fmt.Println(url)
		return nil
	})

}
