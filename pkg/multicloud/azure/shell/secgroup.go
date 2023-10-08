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
	"yunion.io/x/cloudmux/pkg/multicloud/azure"
)

func init() {
	type SecurityGroupListOptions struct {
	}
	shellutils.R(&SecurityGroupListOptions{}, "security-group-list", "List security group", func(cli *azure.SRegion, args *SecurityGroupListOptions) error {
		secgrps, err := cli.ListSecgroups()
		if err != nil {
			return err
		}
		printList(secgrps, len(secgrps), 0, 0, []string{})
		return nil
	})

	type SecurityGroupOptions struct {
		ID string `help:"ID or name of security group"`
	}
	shellutils.R(&SecurityGroupOptions{}, "security-group-show", "Show details of a security group", func(cli *azure.SRegion, args *SecurityGroupOptions) error {
		secgrp, err := cli.GetSecurityGroupDetails(args.ID)
		if err != nil {
			return err
		}
		printObject(secgrp)
		rules, err := secgrp.GetRules()
		if err != nil {
			return err
		}
		printList(rules, 0, 0, 0, nil)
		return nil
	})

	shellutils.R(&SecurityGroupOptions{}, "security-group-rule-list", "List security group rules", func(cli *azure.SRegion, args *SecurityGroupOptions) error {
		if secgroup, err := cli.GetSecurityGroupDetails(args.ID); err != nil {
			return err
		} else if rules, err := secgroup.GetRules(); err != nil {
			return err
		} else {
			printList(rules, len(rules), 0, 30, []string{})
			return nil
		}
	})

	type SecurityGroupCreateOptions struct {
		NAME string `help:"Security Group name"`
	}

	shellutils.R(&SecurityGroupCreateOptions{}, "security-group-create", "Create security group", func(cli *azure.SRegion, args *SecurityGroupCreateOptions) error {
		opts := &cloudprovider.SecurityGroupCreateInput{Name: args.NAME}
		secgrp, err := cli.CreateSecurityGroup(opts)
		if err != nil {
			return err
		}
		printObject(secgrp)
		return nil
	})
}
