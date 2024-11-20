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

	"yunion.io/x/cloudmux/pkg/multicloud/azure"
)

func init() {
	type CloudpolicyListOptions struct {
		Name string
	}
	shellutils.R(&CloudpolicyListOptions{}, "cloud-policy-list", "List cloudpolicies", func(cli *azure.SRegion, args *CloudpolicyListOptions) error {
		roles, err := cli.GetClient().GetRoles(args.Name)
		if err != nil {
			return err
		}
		printList(roles, 0, 0, 0, nil)
		return nil
	})

	type CloudpolicyAssignOption struct {
		OBJECT string
		ROLE   string
	}

	shellutils.R(&CloudpolicyAssignOption{}, "cloud-policy-assign-object", "Assign cloudpolicy for object", func(cli *azure.SRegion, args *CloudpolicyAssignOption) error {
		return cli.GetClient().AssignPolicy(args.OBJECT, args.ROLE)
	})

	type AssignmentListOption struct {
		ObjectId string
	}

	type CloudpolicyAssignListOptions struct {
		ID string
	}

	shellutils.R(&CloudpolicyAssignListOptions{}, "cloud-user-policy-list", "Assign cloudpolicy for object", func(cli *azure.SRegion, args *CloudpolicyAssignListOptions) error {
		ret, err := cli.GetClient().GetPrincipalPolicy(args.ID)
		if err != nil {
			return err
		}
		printList(ret, 0, 0, 0, nil)
		return nil
	})

	shellutils.R(&CloudpolicyAssignListOptions{}, "cloud-group-policy-list", "Assign cloudpolicy for object", func(cli *azure.SRegion, args *CloudpolicyAssignListOptions) error {
		ret, err := cli.GetClient().GetPrincipalPolicy(args.ID)
		if err != nil {
			return err
		}
		printList(ret, 0, 0, 0, nil)
		return nil
	})

}
