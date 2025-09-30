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
	"yunion.io/x/cloudmux/pkg/cloudprovider"
	"yunion.io/x/cloudmux/pkg/multicloud/ksyun"

	"yunion.io/x/pkg/errors"
	"yunion.io/x/pkg/util/shellutils"
)

func init() {
	type DiskListOptions struct {
		DiskId      []string
		StorageType string
		ZoneId      string
	}
	shellutils.R(&DiskListOptions{}, "disk-list", "list disks", func(cli *ksyun.SRegion, args *DiskListOptions) error {
		res, err := cli.GetDisks(args.DiskId, args.StorageType, args.ZoneId)
		if err != nil {
			return errors.Wrap(err, "GetDisks")
		}
		printList(res)
		return nil
	})

	type InstanceDiskListOptions struct {
		INSTANCE_ID string
	}
	shellutils.R(&InstanceDiskListOptions{}, "instance-disk-list", "list disks", func(cli *ksyun.SRegion, args *InstanceDiskListOptions) error {
		res, err := cli.GetDiskByInstanceId(args.INSTANCE_ID)
		if err != nil {
			return errors.Wrap(err, "GetDiskByInstanceId")
		}
		printList(res)
		return nil
	})

	type DiskCreateOptions struct {
		NAME        string
		Desc        string
		SizeGb      int
		ZoneId      string
		StorageType string
	}
	shellutils.R(&DiskCreateOptions{}, "disk-create", "create disk", func(cli *ksyun.SRegion, args *DiskCreateOptions) error {
		disk, err := cli.CreateDisk(args.StorageType, args.ZoneId, &cloudprovider.DiskCreateConfig{
			Name:   args.NAME,
			Desc:   args.Desc,
			SizeGb: args.SizeGb,
		})
		if err != nil {
			return errors.Wrap(err, "CreateDisk")
		}
		printObject(disk)
		return nil
	})

	type DiskDeleteOptions struct {
		ID string
	}
	shellutils.R(&DiskDeleteOptions{}, "disk-delete", "delete disk", func(cli *ksyun.SRegion, args *DiskDeleteOptions) error {
		return cli.DeleteDisk(args.ID)
	})

	type DiskResizeOptions struct {
		ID   string
		SIZE int64
	}
	shellutils.R(&DiskResizeOptions{}, "disk-resize", "resize disk", func(cli *ksyun.SRegion, args *DiskResizeOptions) error {
		return cli.ResizeDisk(args.ID, args.SIZE)
	})

	type DiskResetOptions struct {
		ID       string
		SNAPSHOT string
	}
	shellutils.R(&DiskResetOptions{}, "disk-reset", "reset disk", func(cli *ksyun.SRegion, args *DiskResetOptions) error {
		return cli.ResetDisk(args.ID, args.SNAPSHOT)
	})

	type DiskShowOptions struct {
		ID string
	}
	shellutils.R(&DiskShowOptions{}, "disk-show", "show disk", func(cli *ksyun.SRegion, args *DiskShowOptions) error {
		disk, err := cli.GetDisk(args.ID)
		if err != nil {
			return errors.Wrap(err, "GetDisk")
		}
		printObject(disk)
		return nil
	})

}
