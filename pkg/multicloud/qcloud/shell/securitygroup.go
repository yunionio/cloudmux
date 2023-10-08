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
	"yunion.io/x/cloudmux/pkg/multicloud/qcloud"
)

func init() {
	type SecurityGroupListOptions struct {
		Ids    []string `help:"Secgroup Ids"`
		Name   string   `help:"Secgroup Name"`
		Limit  int      `help:"page size"`
		Offset int      `help:"page offset"`
	}
	shellutils.R(&SecurityGroupListOptions{}, "security-group-list", "List SecurityGroup", func(cli *qcloud.SRegion, args *SecurityGroupListOptions) error {
		secgrps, total, err := cli.GetSecurityGroups(args.Ids, args.Name, args.Limit, args.Offset)
		if err != nil {
			return err
		}
		printList(secgrps, total, args.Offset, args.Limit, []string{})
		return nil
	})

	type SecurityGroupOptions struct {
		ID string `help:"SecurityGroup ID"`
	}
	shellutils.R(&SecurityGroupOptions{}, "security-group-show", "Show SecurityGroup", func(cli *qcloud.SRegion, args *SecurityGroupOptions) error {
		secgroup, err := cli.GetSecurityGroup(args.ID)
		if err != nil {
			return err
		}
		printObject(secgroup)
		return nil
	})

	shellutils.R(&SecurityGroupOptions{}, "security-group-rule-list", "List SecurityGroup rules", func(cli *qcloud.SRegion, args *SecurityGroupOptions) error {
		rules, err := cli.GetSecurityGroupRules(args.ID)
		if err != nil {
			return err
		}
		printList(rules, 0, 0, 0, nil)
		return nil
	})

	shellutils.R(&SecurityGroupOptions{}, "security-group-references", "Show references of a security group", func(cli *qcloud.SRegion, args *SecurityGroupOptions) error {
		references, err := cli.DescribeSecurityGroupReferences(args.ID)
		if err != nil {
			return err
		}
		printList(references, 0, 0, 0, nil)
		return nil
	})

	shellutils.R(&SecurityGroupOptions{}, "security-group-delete", "Delete SecurityGroup", func(cli *qcloud.SRegion, args *SecurityGroupOptions) error {
		return cli.DeleteSecurityGroup(args.ID)
	})

	shellutils.R(&cloudprovider.SecurityGroupCreateInput{}, "security-group-create", "Create SecurityGroup", func(cli *qcloud.SRegion, args *cloudprovider.SecurityGroupCreateInput) error {
		secgrp, err := cli.CreateSecurityGroup(args)
		if err != nil {
			return err
		}
		printObject(secgrp)
		return nil
	})

	type AddressShowOptions struct {
		Id     string `help:"IP address ID"`
		Name   string `help:"IP address name"`
		Limit  int    `help:"page size"`
		Offset int    `help:"page offset"`
	}
	shellutils.R(&AddressShowOptions{}, "address-list", "Show address", func(cli *qcloud.SRegion, args *AddressShowOptions) error {
		address, total, err := cli.AddressList(args.Id, args.Name, args.Offset, args.Limit)
		if err != nil {
			return err
		}
		printList(address, total, args.Offset, args.Limit, []string{})
		return nil
	})

}
