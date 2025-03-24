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
	"yunion.io/x/cloudmux/pkg/cloudprovider"
	"yunion.io/x/cloudmux/pkg/multicloud/ksyun"
	"yunion.io/x/pkg/util/shellutils"
)

func init() {
	type UserListOptions struct {
	}
	shellutils.R(&UserListOptions{}, "cloud-user-list", "List users", func(cli *ksyun.SRegion, args *UserListOptions) error {
		ret, err := cli.GetClient().GetUsers()
		if err != nil {
			return err
		}
		printList(ret)
		return nil
	})

	shellutils.R(&cloudprovider.SClouduserCreateConfig{}, "cloud-user-create", "Create user", func(cli *ksyun.SRegion, args *cloudprovider.SClouduserCreateConfig) error {
		ret, err := cli.GetClient().CreateUser(args)
		if err != nil {
			return err
		}
		printObject(ret)
		return nil
	})

	type UserNameOptions struct {
		NAME string
	}

	shellutils.R(&UserNameOptions{}, "cloud-user-show", "Show user", func(cli *ksyun.SRegion, args *UserNameOptions) error {
		ret, err := cli.GetClient().GetUser(args.NAME)
		if err != nil {
			return err
		}
		printObject(ret)
		return nil
	})

	shellutils.R(&UserNameOptions{}, "cloud-user-delete", "Delete user", func(cli *ksyun.SRegion, args *UserNameOptions) error {
		return cli.GetClient().DeleteUser(args.NAME)
	})

}
