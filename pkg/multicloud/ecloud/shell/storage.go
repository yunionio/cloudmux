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

	"yunion.io/x/cloudmux/pkg/multicloud/ecloud"
)

func init() {
	type StorageListOptions struct {
		ZoneCode string `help:"Zone code to filter storages" optional:"true"`
	}
	shellutils.R(&StorageListOptions{}, "storage-list", "List storages in region (all zones or specific zone)", func(cli *ecloud.SRegion, args *StorageListOptions) error {
		storages, err := cli.GetStorages(args.ZoneCode)
		if err != nil {
			return err
		}
		printList(storages)
		return nil
	})

	type PoolInfoOptions struct {
		ProductType string `help:"Product type, e.g. capebs/ssdebs/ssd" required:"true"`
	}
	shellutils.R(&PoolInfoOptions{}, "pool-info-list", "List EBS pool infos by productType", func(cli *ecloud.SRegion, args *PoolInfoOptions) error {
		pools, err := cli.GetPoolInfo(args.ProductType)
		if err != nil {
			return err
		}
		printList(pools)
		return nil
	})

	type VolumeConfigOptions struct {
	}
	shellutils.R(&VolumeConfigOptions{}, "volume-config-list", "List volume type config", func(cli *ecloud.SRegion, args *VolumeConfigOptions) error {
		cfg, err := cli.GetVolumeConfig()
		if err != nil {
			return err
		}
		printList(cfg)
		return nil
	})
}

