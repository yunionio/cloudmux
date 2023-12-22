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
	"time"

	"github.com/pkg/errors"
	"yunion.io/x/cloudmux/pkg/cloudprovider"
	"yunion.io/x/cloudmux/pkg/multicloud/volcengine"
	"yunion.io/x/pkg/util/shellutils"
)

func init() {
	type DiskListOptions struct {
		Instance string   `help:"Instance ID"`
		Zone     string   `help:"Zone ID"`
		DiskType string   `help:"Disk category"`
		DiskId   []string `help:"Disk IDs"`
	}
	shellutils.R(&DiskListOptions{}, "disk-list", "List disks", func(cli *volcengine.SRegion, args *DiskListOptions) error {
		disks, e := cli.GetDisks(args.Instance, args.Zone, args.DiskType, args.DiskId)
		if e != nil {
			return e
		}
		printList(disks, 0, 0, 0, nil)
		return nil
	})

	type DiskDeleteOptions struct {
		ID string `help:"Disk ID"`
	}
	shellutils.R(&DiskDeleteOptions{}, "disk-delete", "List disks", func(cli *volcengine.SRegion, args *DiskDeleteOptions) error {
		e := cli.DeleteDisk(args.ID)
		if e != nil {
			return e
		}
		return nil
	})

	type DiskResizeOptions struct {
		ID      string `help:"Disk ID"`
		SIZE_GB int64
	}
	shellutils.R(&DiskResizeOptions{}, "disk-resize", "Resize disks", func(cli *volcengine.SRegion, args *DiskResizeOptions) error {
		e := cli.ResizeDisk(args.ID, args.SIZE_GB)
		if e != nil {
			return e
		}
		return nil
	})

	type DiskCreateOptions struct {
		Name     string
		Desc     string
		Category string `choices:"PTSSD|ESSD_PL0|ESSD_FlexPL" default:"PTSSD"`
		ZoneId   string
		SizeGb   int `default:"10"`
	}

	shellutils.R(&DiskCreateOptions{}, "disk-create", "create a disk", func(cli *volcengine.SRegion, args *DiskCreateOptions) error {
		diskId, err := cli.CreateDisk(
			args.ZoneId,
			args.Category,
			args.Name,
			args.SizeGb,
			args.Desc,
			"",
		)
		if err != nil {
			return err
		}
		err = cloudprovider.Wait(time.Second, time.Minute, func() (bool, error) {
			disk, err := cli.GetDisk(diskId)
			if err != nil {
				if errors.Cause(err) != cloudprovider.ErrNotFound {
					return false, err
				}
				return false, nil
			}
			printObject(disk)
			return true, nil
		})
		return err
	})
}
