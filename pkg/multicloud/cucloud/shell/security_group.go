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
	type SecurityGroupListOptions struct {
		Id string
	}
	shellutils.R(&SecurityGroupListOptions{}, "security-group-list", "list security groups", func(cli *cucloud.SRegion, args *SecurityGroupListOptions) error {
		securityGroups, err := cli.GetSecurityGroups(args.Id)
		if err != nil {
			return err
		}
		printList(securityGroups)
		return nil
	})

	type SeucrityGroupIdOptions struct {
		ID string
	}
	shellutils.R(&SeucrityGroupIdOptions{}, "security-group-show", "Show security group", func(cli *cucloud.SRegion, args *SeucrityGroupIdOptions) error {
		secgroup, err := cli.GetSecurityGroup(args.ID)
		if err != nil {
			return err
		}
		printObject(secgroup)
		return nil
	})
}
