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

	"yunion.io/x/cloudmux/pkg/multicloud/aliyun"
)

func init() {
	type MFAListOptions struct {
	}
	shellutils.R(&MFAListOptions{}, "mfa-device-list", "List mfa device", func(cli *aliyun.SRegion, args *MFAListOptions) error {
		ret, err := cli.GetClient().ListVirtualMFADevices()
		if err != nil {
			return err
		}
		printList(ret, 0, 0, 0, []string{})
		return nil
	})

	type MFACreateOptions struct {
		NAME string
	}
	shellutils.R(&MFACreateOptions{}, "mfa-device-create", "Create mfa device", func(cli *aliyun.SRegion, args *MFACreateOptions) error {
		ret, err := cli.GetClient().CreateVirtualMFADevice(args.NAME)
		if err != nil {
			return err
		}
		printObject(ret)
		return nil
	})

	type MFAUserOptions struct {
		USER string
	}

	shellutils.R(&MFAUserOptions{}, "mfa-device-unbind", "Unbind mfa device", func(cli *aliyun.SRegion, args *MFAUserOptions) error {
		return cli.GetClient().UnbindMFADevice(args.USER)
	})

	shellutils.R(&MFAUserOptions{}, "mfa-device-disable", "Disable mfa device", func(cli *aliyun.SRegion, args *MFAUserOptions) error {
		return cli.GetClient().DisableVirtualMFA(args.USER)
	})

	type MFABindOptions struct {
		USER  string
		SN    string
		CODE1 string
		CODE2 string
	}

	shellutils.R(&MFABindOptions{}, "mfa-device-bind", "Bind mfa device", func(cli *aliyun.SRegion, args *MFABindOptions) error {
		return cli.GetClient().BindMFADevice(args.USER, args.SN, args.CODE1, args.CODE2)
	})
}
