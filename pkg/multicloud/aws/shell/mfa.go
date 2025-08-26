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

	"yunion.io/x/cloudmux/pkg/multicloud/aws"
)

func init() {
	type VirtualMFADeviceListOptions struct {
	}

	shellutils.R(&VirtualMFADeviceListOptions{}, "virtual-mfa-device-list", "List virtual mfa devices", func(cli *aws.SRegion, args *VirtualMFADeviceListOptions) error {
		devices, err := cli.GetClient().GetVirtualMFADevices()
		if err != nil {
			return err
		}
		printList(devices, 0, 0, 0, nil)
		return nil
	})

	type VirtualMFADeviceDeleteOptions struct {
		SerialNumber string
		UserName     string
	}

	shellutils.R(&VirtualMFADeviceDeleteOptions{}, "virtual-mfa-device-deactivate", "Deactivate virtual mfa device", func(cli *aws.SRegion, args *VirtualMFADeviceDeleteOptions) error {
		return cli.GetClient().DeleteVirtualMFADevice(args.SerialNumber, args.UserName)
	})

	type VirtualMFADeviceCreateOptions struct {
		NAME string
	}

	shellutils.R(&VirtualMFADeviceCreateOptions{}, "virtual-mfa-device-create", "Create virtual mfa device", func(cli *aws.SRegion, args *VirtualMFADeviceCreateOptions) error {
		device, err := cli.GetClient().CreateVirtualMFADevice(args.NAME)
		if err != nil {
			return err
		}
		printObject(device)
		return nil
	})

	type VirtualMFADeviceResyncOptions struct {
		SERIAL_NUMBER string
		USER          string
		CODE1         string
		CODE2         string
	}

	shellutils.R(&VirtualMFADeviceResyncOptions{}, "virtual-mfa-device-resync", "Resync virtual mfa device", func(cli *aws.SRegion, args *VirtualMFADeviceResyncOptions) error {
		return cli.GetClient().ResyncMFADevice(args.SERIAL_NUMBER, args.USER, args.CODE1, args.CODE2)
	})
}
