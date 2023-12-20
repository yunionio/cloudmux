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
	"yunion.io/x/cloudmux/pkg/multicloud/huawei"
)

func init() {
	type SecurityGroupListOptions struct {
		Name string `help:"Secgroup name"`
	}
	shellutils.R(&SecurityGroupListOptions{}, "security-group-list", "List security group", func(cli *huawei.SRegion, args *SecurityGroupListOptions) error {
		secgrps, e := cli.GetSecurityGroups(args.Name)
		if e != nil {
			return e
		}
		printList(secgrps, 0, 0, 0, nil)
		return nil
	})

	type SecurityGroupShowOptions struct {
		ID string `help:"ID or name of security group"`
	}
	shellutils.R(&SecurityGroupShowOptions{}, "security-group-show", "Show details of a security group", func(cli *huawei.SRegion, args *SecurityGroupShowOptions) error {
		secgrp, err := cli.GetSecurityGroup(args.ID)
		if err != nil {
			return err
		}
		printObject(secgrp)
		return nil
	})

	shellutils.R(&SecurityGroupShowOptions{}, "security-group-rule-list", "List of a security group rules", func(cli *huawei.SRegion, args *SecurityGroupShowOptions) error {
		rules, err := cli.GetSecurityGroupRules(args.ID)
		if err != nil {
			return err
		}
		printList(rules, 0, 0, 0, nil)
		return nil
	})

	shellutils.R(&cloudprovider.SecurityGroupCreateInput{}, "security-group-create", "Create security group", func(cli *huawei.SRegion, args *cloudprovider.SecurityGroupCreateInput) error {
		result, err := cli.CreateSecurityGroup(args)
		if err != nil {
			return err
		}
		printObject(result)
		return nil
	})

	type SecurityGroupRuleIdOptions struct {
		ID string
	}

	shellutils.R(&SecurityGroupRuleIdOptions{}, "security-group-rule-delete", "Delete security group rule", func(cli *huawei.SRegion, args *SecurityGroupRuleIdOptions) error {
		return cli.DeleteSecurityGroupRule(args.ID)
	})

	type SecurityGroupRuleCreateOptions struct {
		SECGROUP_ID string
		cloudprovider.SecurityGroupRuleCreateOptions
	}

	shellutils.R(&SecurityGroupRuleCreateOptions{}, "security-group-rule-create", "Create security group rule", func(cli *huawei.SRegion, args *SecurityGroupRuleCreateOptions) error {
		rule, err := cli.CreateSecurityGroupRule(args.SECGROUP_ID, &args.SecurityGroupRuleCreateOptions)
		if err != nil {
			return err
		}
		printObject(rule)
		return nil
	})

}
