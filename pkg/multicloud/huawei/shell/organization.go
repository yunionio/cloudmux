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

	"yunion.io/x/cloudmux/pkg/multicloud/huawei"
)

func init() {
	type OrganizationShowOptions struct {
	}
	shellutils.R(&OrganizationShowOptions{}, "organization-show", "Show organization", func(cli *huawei.SRegion, args *OrganizationShowOptions) error {
		org, err := cli.GetClient().ShowOrganization()
		if err != nil {
			return err
		}
		printObject(org)
		return nil
	})

	type OrganizationAccountListOptions struct {
	}
	shellutils.R(&OrganizationAccountListOptions{}, "organization-account-list", "List organization accounts", func(cli *huawei.SRegion, args *OrganizationAccountListOptions) error {
		accounts, err := cli.GetClient().ListAccounts()
		if err != nil {
			return err
		}
		printList(accounts, 0, 0, 0, nil)
		return nil
	})

	type OrganizationEntityListOptions struct {
		Parent string `help:"Parent root or organizational unit id"`
	}
	shellutils.R(&OrganizationEntityListOptions{}, "organization-entity-list", "List organization entities", func(cli *huawei.SRegion, args *OrganizationEntityListOptions) error {
		entities, err := cli.GetClient().ListEntities(args.Parent)
		if err != nil {
			return err
		}
		printList(entities, 0, 0, 0, nil)
		return nil
	})

	type OrganizationalUnitShowOptions struct {
		ID string `help:"Organizational unit id"`
	}
	shellutils.R(&OrganizationalUnitShowOptions{}, "organization-ou-show", "Show organizational unit", func(cli *huawei.SRegion, args *OrganizationalUnitShowOptions) error {
		ou, err := cli.GetClient().ShowOrganizationalUnit(args.ID)
		if err != nil {
			return err
		}
		printObject(ou)
		return nil
	})
}
