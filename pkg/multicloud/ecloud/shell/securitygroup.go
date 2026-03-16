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
	"yunion.io/x/pkg/util/secrules"
	"yunion.io/x/pkg/util/shellutils"

	"yunion.io/x/cloudmux/pkg/cloudprovider"
	"yunion.io/x/cloudmux/pkg/multicloud/ecloud"
)

func init() {
	type SecgroupListOptions struct {
	}
	shellutils.R(&SecgroupListOptions{}, "secgroup-list", "List security groups", func(cli *ecloud.SRegion, args *SecgroupListOptions) error {
		groups, err := cli.GetISecurityGroups()
		if err != nil {
			return err
		}
		printList(groups)
		return nil
	})

	type SecgroupCreateOptions struct {
		Name string `help:"Security group name"`
		Desc string `help:"Description" optional:"true"`
	}
	shellutils.R(&SecgroupCreateOptions{}, "secgroup-create", "Create security group", func(cli *ecloud.SRegion, args *SecgroupCreateOptions) error {
		opts := &cloudprovider.SecurityGroupCreateInput{
			Name: args.Name,
			Desc: args.Desc,
		}
		group, err := cli.CreateSecurityGroup(opts)
		if err != nil {
			return err
		}
		printObject(group)
		return nil
	})

	type SecgroupIdOptions struct {
		ID string `help:"Security group id"`
	}
	shellutils.R(&SecgroupIdOptions{}, "secgroup-show", "Show security group detail", func(cli *ecloud.SRegion, args *SecgroupIdOptions) error {
		group, err := cli.GetISecurityGroupById(args.ID)
		if err != nil {
			return err
		}
		printObject(group)
		return nil
	})

	shellutils.R(&SecgroupIdOptions{}, "secgroup-delete", "Delete security group", func(cli *ecloud.SRegion, args *SecgroupIdOptions) error {
		group, err := cli.GetISecurityGroupById(args.ID)
		if err != nil {
			return err
		}
		return group.Delete()
	})

	type SecgroupRuleListOptions struct {
		ID string `help:"Security group ID"`
	}
	shellutils.R(&SecgroupRuleListOptions{}, "secgroup-rule-list", "List security group rules", func(cli *ecloud.SRegion, args *SecgroupRuleListOptions) error {
		rules, err := cli.GetSecurityGroupRules(args.ID)
		if err != nil {
			return err
		}
		printList(rules)
		return nil
	})

	type SecgroupRuleCreateOptions struct {
		SecgroupId string `help:"Security group ID"`
		Direction  string `help:"Direction: in|out" choices:"in|out"`
		Protocol   string `help:"Protocol: tcp|udp|icmp|ANY" default:"ANY"`
		Ports      string `help:"Port or range, e.g. 80 or 80-443" optional:"true"`
		CIDR       string `help:"Remote CIDR, e.g. 0.0.0.0/0" optional:"true"`
		Desc       string `help:"Description" optional:"true"`
	}
	shellutils.R(&SecgroupRuleCreateOptions{}, "secgroup-rule-create", "Create security group rule", func(cli *ecloud.SRegion, args *SecgroupRuleCreateOptions) error {
		opts := &cloudprovider.SecurityGroupRuleCreateOptions{
			Direction: secrules.TSecurityRuleDirection(args.Direction),
			Protocol:  args.Protocol,
			Ports:     args.Ports,
			CIDR:      args.CIDR,
			Desc:      args.Desc,
			Action:    secrules.SecurityRuleAllow,
		}
		rule, err := cli.CreateSecurityGroupRule(args.SecgroupId, opts)
		if err != nil {
			return err
		}
		printObject(rule)
		return nil
	})

	type SecgroupRuleDeleteOptions struct {
		RuleId string `help:"Security group rule ID"`
	}
	shellutils.R(&SecgroupRuleDeleteOptions{}, "secgroup-rule-delete", "Delete security group rule", func(cli *ecloud.SRegion, args *SecgroupRuleDeleteOptions) error {
		return cli.DeleteSecurityGroupRule(args.RuleId)
	})
}
