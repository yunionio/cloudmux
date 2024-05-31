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

	"yunion.io/x/cloudmux/pkg/multicloud/baidu"
)

func init() {
	type StorageListOptions struct {
		ZoneName string
	}
	shellutils.R(&StorageListOptions{}, "storage-list", "list storage", func(cli *baidu.SRegion, args *StorageListOptions) error {
		storages, err := cli.GetStorageTypes(args.ZoneName)
		if err != nil {
			return err
		}
		printList(storages)
		return nil
	})

}
