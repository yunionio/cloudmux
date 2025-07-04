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

	"yunion.io/x/cloudmux/pkg/cloudprovider"
	"yunion.io/x/cloudmux/pkg/multicloud/aws"
)

func init() {
	type DiskListOptions struct {
		Instance   string `help:"Instance ID"`
		Zone       string `help:"Zone ID"`
		VolumeType string `help:"Disk category" choices:"gp2|gp3|io1|io2|st1|sc1|standard"`
		DiskIds    []string
	}
	shellutils.R(&DiskListOptions{}, "disk-list", "List disks", func(cli *aws.SRegion, args *DiskListOptions) error {
		disks, err := cli.GetDisks(args.Instance, args.Zone, args.VolumeType, args.DiskIds)
		if err != nil {
			return err
		}
		printList(disks, 0, 0, 0, []string{})
		return nil
	})

	type DiskDeleteOptions struct {
		ID string `help:"Disk ID"`
	}
	shellutils.R(&DiskDeleteOptions{}, "disk-delete", "List disks", func(cli *aws.SRegion, args *DiskDeleteOptions) error {
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
	shellutils.R(&DiskResizeOptions{}, "disk-resize", "List disks", func(cli *aws.SRegion, args *DiskResizeOptions) error {
		e := cli.ResizeDisk(args.ID, args.SIZE_GB)
		if e != nil {
			return e
		}
		return nil
	})

	type VolumeCreateOptions struct {
		Name       string
		Desc       string
		VolumeType string `choices:"gp2|gp3|io1|io2|st1|sc1|standard" default:"gp2"`
		ZoneId     string
		SizeGb     int `default:"10"`
		Iops       int
		Throughput int
		SnapshotId string
	}

	shellutils.R(&VolumeCreateOptions{}, "disk-create", "create a volume", func(cli *aws.SRegion, args *VolumeCreateOptions) error {
		opts := &cloudprovider.DiskCreateConfig{
			Name:       args.Name,
			SizeGb:     args.SizeGb,
			Iops:       args.Iops,
			Throughput: args.Throughput,
			SnapshotId: args.SnapshotId,
			Desc:       args.Desc,
		}
		volume, err := cli.CreateDisk(args.ZoneId, args.VolumeType, opts)
		if err != nil {
			return err
		}
		printObject(volume)
		return nil
	})

}
