// Copyright 2023 Yunion
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package shell

import (
	"fmt"

	"yunion.io/x/cloudmux/pkg/cloudprovider"
	"yunion.io/x/cloudmux/pkg/multicloud/volcengine"
	"yunion.io/x/pkg/util/shellutils"
)

func init() {
	type SecurityGroupListOptions struct {
		VpcId            string   `help:"VPC ID"`
		Name             string   `help:"Secgroup Name"`
		SecurityGroupIds []string `help:"SecurityGroup ids"`
		Limit            int      `help:"page size"`
		Offset           int      `help:"page offset"`
	}
	shellutils.R(&SecurityGroupListOptions{}, "security-group-list", "List security group", func(cli *volcengine.SRegion, args *SecurityGroupListOptions) error {
		secgrps, total, e := cli.GetSecurityGroups(args.VpcId, args.Name, args.SecurityGroupIds, args.Offset, args.Limit)
		if e != nil {
			return e
		}
		printList(secgrps, total, args.Offset, args.Limit, []string{})
		return nil
	})

	type SecurityGroupIdOptions struct {
		ID string `help:"ID or name of security group"`
	}
	shellutils.R(&SecurityGroupIdOptions{}, "security-group-rule-list", "Show details of a security group", func(cli *volcengine.SRegion, args *SecurityGroupIdOptions) error {
		rules, err := cli.GetSecurityGroupRules(args.ID)
		if err != nil {
			return err
		}
		printList(rules, 0, 0, 0, nil)
		return nil
	})

	shellutils.R(&SecurityGroupIdOptions{}, "security-group-delete", "Delete security group", func(cli *volcengine.SRegion, args *SecurityGroupIdOptions) error {
		return cli.DeleteSecurityGroup(args.ID)
	})

	type SecurityGroupCreateOptions struct {
		NAME            string `help:"SecurityGroup name"`
		VpcId           string `help:"VPC ID"`
		ResourceGroupId string
		Desc            string `help:"SecurityGroup description"`
	}

	shellutils.R(&cloudprovider.SecurityGroupCreateInput{}, "security-group-create", "Create details of a security group", func(cli *volcengine.SRegion, args *cloudprovider.SecurityGroupCreateInput) error {
		secgroupId, err := cli.CreateSecurityGroup(args)
		if err != nil {
			return err
		}
		fmt.Printf("secgroupId: %s", secgroupId)
		return nil
	})
}
