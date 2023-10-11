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

	"yunion.io/x/cloudmux/pkg/multicloud/ctyun"
)

func init() {
	type SDiskListOptions struct {
	}
	shellutils.R(&SDiskListOptions{}, "disk-list", "List disks", func(cli *ctyun.SRegion, args *SDiskListOptions) error {
		disks, err := cli.GetDisks()
		if err != nil {
			return err
		}
		printList(disks, 0, 0, 0, nil)
		return nil
	})

	type DiskCreateOptions struct {
		ZoneId   string `help:"zone id"`
		NAME     string `help:"disk name"`
		DiskType string `help:"disk type" choice:"SATA|SSD-genric|SSD|FAST-SSD" default:"SATA"`
		SizeGb   int    `help:"disk size" default:"10"`
	}
	shellutils.R(&DiskCreateOptions{}, "disk-create", "Create disk", func(cli *ctyun.SRegion, args *DiskCreateOptions) error {
		disk, err := cli.CreateDisk(args.ZoneId, args.NAME, args.DiskType, args.SizeGb)
		if err != nil {
			return err
		}
		printObject(disk)
		return nil
	})

	type DiskResizeOptions struct {
		DISK string `help:"disk id"`
		SIZE int64  `help:"disk size GB"`
	}
	shellutils.R(&DiskResizeOptions{}, "disk-resize", "Resize disk", func(cli *ctyun.SRegion, args *DiskResizeOptions) error {
		return cli.ResizeDisk(args.DISK, args.SIZE)
	})

	type DiskIdOptions struct {
		ID string
	}

	shellutils.R(&DiskIdOptions{}, "disk-delete", "Delete disk", func(cli *ctyun.SRegion, args *DiskIdOptions) error {
		return cli.DeleteDisk(args.ID)
	})

	shellutils.R(&DiskIdOptions{}, "disk-show", "Show disk", func(cli *ctyun.SRegion, args *DiskIdOptions) error {
		disk, err := cli.GetDisk(args.ID)
		if err != nil {
			return err
		}
		printObject(disk)
		return nil
	})

}
