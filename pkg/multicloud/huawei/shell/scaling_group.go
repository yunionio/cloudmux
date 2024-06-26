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

	"yunion.io/x/cloudmux/pkg/multicloud/huawei"
)

func init() {
	type ScalingGroupListOptions struct {
	}
	shellutils.R(&ScalingGroupListOptions{}, "scaling-group-list", "List scaling groups", func(cli *huawei.SRegion, args *ScalingGroupListOptions) error {
		ret, err := cli.ListScalingGroups()
		if err != nil {
			return err
		}
		printList(ret, 0, 0, 0, []string{})
		return nil
	})

	type ScalingInstanceListOptions struct {
		GROUP string
	}
	shellutils.R(&ScalingInstanceListOptions{}, "scaling-instance-list", "List scaling instances", func(cli *huawei.SRegion, args *ScalingInstanceListOptions) error {
		ret, err := cli.ListScalingInstances(args.GROUP)
		if err != nil {
			return err
		}
		printList(ret, 0, 0, 0, []string{})
		return nil
	})

}
