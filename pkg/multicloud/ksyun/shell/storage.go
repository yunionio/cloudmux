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
	"fmt"

	"yunion.io/x/cloudmux/pkg/multicloud/ksyun"

	"yunion.io/x/pkg/errors"
	"yunion.io/x/pkg/util/shellutils"
)

func init() {
	type StorageListOptions struct {
	}
	shellutils.R(&StorageListOptions{}, "storage-list", "list storage", func(cli *ksyun.SRegion, args *StorageListOptions) error {
		res, err := cli.GetIStorages()
		if err != nil {
			return errors.Wrap(err, "GetIStorages")
		}
		printList(res)
		return nil
	})
	type StorageZoneListOptions struct {
		DiskType string
	}
	shellutils.R(&StorageZoneListOptions{}, "storage-zone-list", "list storage zone", func(cli *ksyun.SRegion, args *StorageZoneListOptions) error {
		res, err := cli.GetZonesByDiskType(args.DiskType)
		if err != nil {
			return errors.Wrap(err, "GetZonesByDiskType")
		}
		fmt.Println(res)
		return nil
	})
}
