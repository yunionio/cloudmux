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
	"yunion.io/x/cloudmux/pkg/multicloud/ctyun"
)

func init() {
	type SSecurityGroupListOptions struct {
	}
	shellutils.R(&SSecurityGroupListOptions{}, "security-group-list", "List secgroups", func(cli *ctyun.SRegion, args *SSecurityGroupListOptions) error {
		secgroups, e := cli.GetSecurityGroups()
		if e != nil {
			return e
		}
		printList(secgroups, 0, 0, 0, nil)
		return nil
	})

	type SSecurityGroupIdOptions struct {
		ID string `help:"Security Group ID"`
	}
	shellutils.R(&SSecurityGroupIdOptions{}, "security-group-show", "Show secgroup", func(cli *ctyun.SRegion, args *SSecurityGroupIdOptions) error {
		group, e := cli.GetSecurityGroup(args.ID)
		if e != nil {
			return e
		}
		printObject(group)
		return nil
	})

	shellutils.R(&SSecurityGroupIdOptions{}, "security-group-delete", "Delete secgroup", func(cli *ctyun.SRegion, args *SSecurityGroupIdOptions) error {
		return cli.DeleteSecurityGroup(args.ID)
	})

	type SecurityGroupCreateOptions struct {
		VpcId string `help:"vpc id"`
		Name  string `help:"secgroup name"`
	}
	shellutils.R(&cloudprovider.SecurityGroupCreateInput{}, "security-group-create", "Create secgroup", func(cli *ctyun.SRegion, args *cloudprovider.SecurityGroupCreateInput) error {
		sec, e := cli.CreateSecurityGroup(args)
		if e != nil {
			return e
		}
		printObject(sec)
		return nil
	})

	type SecurityGroupRuleCreateOptions struct {
		GROUP string `help:"secgroup id"`
		cloudprovider.SecurityGroupRuleCreateOptions
	}
	shellutils.R(&SecurityGroupRuleCreateOptions{}, "security-group-rule-create", "Create secgroup rule", func(cli *ctyun.SRegion, args *SecurityGroupRuleCreateOptions) error {
		return cli.CreateSecurityGroupRule(args.GROUP, &args.SecurityGroupRuleCreateOptions)
	})
}
