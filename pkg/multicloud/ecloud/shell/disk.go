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
	"context"
	"yunion.io/x/pkg/util/shellutils"

	"yunion.io/x/cloudmux/pkg/multicloud/ecloud"
)

func init() {
	type VDiskListOptions struct {
	}
	shellutils.R(&VDiskListOptions{}, "disk-list", "List disks", func(cli *ecloud.SRegion, args *VDiskListOptions) error {
		disks, err := cli.GetDisks()
		if err != nil {
			return err
		}
		printList(disks)
		return nil
	})

	type DiskShowOptions struct {
		ID string `help:"Disk ID"`
	}
	shellutils.R(&DiskShowOptions{}, "disk-show", "Show disk detail", func(cli *ecloud.SRegion, args *DiskShowOptions) error {
		disk, err := cli.GetDisk(args.ID)
		if err != nil {
			return err
		}
		printObject(disk)
		return nil
	})

	type DiskDeleteOptions struct {
		ID string `help:"Disk ID"`
	}
	shellutils.R(&DiskDeleteOptions{}, "disk-delete", "Delete disk (pre-delete)", func(cli *ecloud.SRegion, args *DiskDeleteOptions) error {
		return cli.PreDeleteVolume(args.ID)
	})

	type DiskResizeOptions struct {
		ID   string `help:"Disk ID"`
		SIZE int64  `help:"New disk size GB (must be >= current)"`
	}
	shellutils.R(&DiskResizeOptions{}, "disk-resize", "Resize disk", func(cli *ecloud.SRegion, args *DiskResizeOptions) error {
		return cli.ResizeDisk(context.Background(), args.ID, args.SIZE)
	})
}
