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

	"yunion.io/x/cloudmux/pkg/multicloud/cucloud"
)

func init() {
	type InstanceListOptions struct {
		ZoneId string
		Id     string
	}
	shellutils.R(&InstanceListOptions{}, "instance-list", "list instances", func(cli *cucloud.SRegion, args *InstanceListOptions) error {
		instances, err := cli.GetInstances(args.ZoneId, args.Id)
		if err != nil {
			return err
		}
		printList(instances)
		return nil
	})

	type InstanceIdOptions struct {
		ID string `help:"Instance ID"`
	}

	shellutils.R(&InstanceIdOptions{}, "instance-stop", "stop instance", func(cli *cucloud.SRegion, args *InstanceIdOptions) error {
		return cli.StopVM(args.ID)
	})

	shellutils.R(&InstanceIdOptions{}, "instance-start", "start instance", func(cli *cucloud.SRegion, args *InstanceIdOptions) error {
		return cli.StartVM(args.ID)
	})
}
