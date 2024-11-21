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
	"yunion.io/x/pkg/errors"
	"yunion.io/x/pkg/util/shellutils"

	"yunion.io/x/cloudmux/pkg/multicloud/baidu"
)

func init() {
	type diskListOptions struct {
		ZoneName    string
		StorageType string
		InstanceId  string
	}
	shellutils.R(&diskListOptions{}, "disk-list", "list disks", func(cli *baidu.SRegion, args *diskListOptions) error {
		disks, err := cli.GetDisks(args.StorageType, args.ZoneName, args.InstanceId)
		if err != nil {
			return err
		}
		printList(disks)
		return nil
	})
	type diskShowOptions struct {
		ID string `help:"ID of disk to show"`
	}
	shellutils.R(&diskShowOptions{}, "disk-show", "list disks", func(cli *baidu.SRegion, args *diskShowOptions) error {
		disk, err := cli.GetDisk(args.ID)
		if err != nil {
			return errors.Wrap(err, "Getdisk")
		}
		printObject(disk)
		return nil
	})
}
