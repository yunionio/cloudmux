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

	"yunion.io/x/cloudmux/pkg/multicloud/oracle"
)

func init() {
	type SecurityGroupListOptions struct {
		VpcId string
	}
	shellutils.R(&SecurityGroupListOptions{}, "security-group-list", "list security-groups", func(cli *oracle.SRegion, args *SecurityGroupListOptions) error {
		securitygroups, err := cli.GetSecurityGroups(args.VpcId)
		if err != nil {
			return err
		}
		printList(securitygroups, 0, 0, 0, []string{})
		return nil
	})

	type SecurityGroupIdOptions struct {
		ID string
	}

	shellutils.R(&SecurityGroupIdOptions{}, "security-group-show", "Show security-group", func(cli *oracle.SRegion, args *SecurityGroupIdOptions) error {
		securitygroup, err := cli.GetSecurityGroup(args.ID)
		if err != nil {
			return err
		}
		printObject(securitygroup)
		return nil
	})

}
