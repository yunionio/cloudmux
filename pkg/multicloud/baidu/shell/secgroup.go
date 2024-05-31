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

	"yunion.io/x/cloudmux/pkg/cloudprovider"
	"yunion.io/x/cloudmux/pkg/multicloud/baidu"
)

func init() {
	type SecurityGroupListOptions struct {
		VpcId string
	}
	shellutils.R(&SecurityGroupListOptions{}, "security-group-list", "list security groups", func(cli *baidu.SRegion, args *SecurityGroupListOptions) error {
		secgroups, err := cli.GetSecurityGroups(args.VpcId)
		if err != nil {
			return err
		}
		printList(secgroups)
		return nil
	})

	type SecurityGroupIdOptions struct {
		ID string
	}

	shellutils.R(&SecurityGroupIdOptions{}, "security-group-show", "show security group", func(cli *baidu.SRegion, args *SecurityGroupIdOptions) error {
		secgroup, err := cli.GetSecurityGroup(args.ID)
		if err != nil {
			return err
		}
		printObject(secgroup)
		return nil
	})

	shellutils.R(&SecurityGroupIdOptions{}, "security-group-delete", "delete security group", func(cli *baidu.SRegion, args *SecurityGroupIdOptions) error {
		return cli.DeleteSecurityGroup(args.ID)
	})

	shellutils.R(&cloudprovider.SecurityGroupCreateInput{}, "security-group-create", "create security group", func(cli *baidu.SRegion, args *cloudprovider.SecurityGroupCreateInput) error {
		group, err := cli.CreateSecurityGroup(args)
		if err != nil {
			return err
		}
		printObject(group)
		return nil
	})

}
