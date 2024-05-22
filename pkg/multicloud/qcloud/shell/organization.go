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

	"yunion.io/x/cloudmux/pkg/multicloud/qcloud"
)

func init() {
	type OrganizationListOptions struct {
	}
	shellutils.R(&OrganizationListOptions{}, "organization-member-list", "List organization members", func(cli *qcloud.SRegion, args *OrganizationListOptions) error {
		ret, err := cli.GetClient().DescribeOrganizationMembers()
		if err != nil {
			return err
		}
		printList(ret, 0, 0, 0, []string{})
		return nil
	})

	type OrganizationIdOptions struct {
		ID string
	}

	shellutils.R(&OrganizationIdOptions{}, "organization-member-show", "Show organization member", func(cli *qcloud.SRegion, args *OrganizationIdOptions) error {
		ret, err := cli.GetClient().DescribeOrganizationMember(args.ID)
		if err != nil {
			return err
		}
		printObject(ret)
		return nil
	})
}
