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
	type DiskListOptions struct {
		Instance         string `help:"Instance ID"`
		Zone             string `help:"Zone ID"`
		Category         string `help:"Disk category"`
		SnapshotpolicyId string
	}
	shellutils.R(&DiskListOptions{}, "disk-list", "List disks", func(cli *aliyun.SRegion, args *DiskListOptions) error {
		disks, err := cli.GetDisks(args.Instance, args.Zone, args.Category, nil, args.SnapshotpolicyId)
		if err != nil {
			return err
		}
		printList(disks, 0, 0, 0, []string{})
		return nil
	})

	type DiskDeleteOptions struct {
		ID string `help:"Instance ID"`
	}
	shellutils.R(&DiskDeleteOptions{}, "disk-delete", "List disks", func(cli *aliyun.SRegion, args *DiskDeleteOptions) error {
		e := cli.DeleteDisk(args.ID)
		if e != nil {
			return e
		}
		return nil
	})
}
