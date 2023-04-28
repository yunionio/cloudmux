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

	"yunion.io/x/pkg/util/secrules"
	"yunion.io/x/pkg/util/shellutils"

	"yunion.io/x/cloudmux/pkg/multicloud/aws"
)

func init() {
	type SecurityGroupListOptions struct {
		VpcId string `help:"VPC ID"`
		Id    string
		Name  string `help:"Secgroup name"`
	}
	shellutils.R(&SecurityGroupListOptions{}, "security-group-list", "List security group", func(cli *aws.SRegion, args *SecurityGroupListOptions) error {
		secgrps, err := cli.GetSecurityGroups(args.VpcId, args.Name, args.Id)
		if err != nil {
			return err
		}
		printList(secgrps, 0, 0, 0, []string{})
		return nil
	})

	type SecurityGroupIdOptions struct {
		ID string
	}
	shellutils.R(&SecurityGroupIdOptions{}, "security-group-show", "Show security group", func(cli *aws.SRegion, args *SecurityGroupIdOptions) error {
		group, err := cli.GetSecurityGroup(args.ID)
		if err != nil {
			return err
		}
		printObject(group)
		rules, err := group.GetRules()
		if err != nil {
			return err
		}
		printList(rules, 0, 0, 0, []string{})
		return nil
	})

	type SecurityGroupCreateOptions struct {
		VPC  string `help:"vpcId"`
		NAME string `help:"group name"`
		DESC string `help:"group desc"`
	}
	shellutils.R(&SecurityGroupCreateOptions{}, "security-group-create", "Create  security group", func(cli *aws.SRegion, args *SecurityGroupCreateOptions) error {
		id, err := cli.CreateSecurityGroup(args.VPC, args.NAME, args.DESC)
		if err != nil {
			return err
		}
		fmt.Println(id)
		return nil
	})

	type SecurityGroupRuleDeleteOption struct {
		SECGROUP_ID string
		RULE        string
	}

	shellutils.R(&SecurityGroupRuleDeleteOption{}, "security-group-rule-delete", "Delete  security group rule", func(cli *aws.SRegion, args *SecurityGroupRuleDeleteOption) error {
		rule := secrules.MustParseSecurityRule(args.RULE)
		return cli.RemoveSecurityGroupRule(args.SECGROUP_ID, *rule)
	})

}
