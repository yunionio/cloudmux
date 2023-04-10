// @@ -0,0 +1,46 @@
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
	huawei "yunion.io/x/cloudmux/pkg/multicloud/hcso"
)

func init() {
	type ModelartsPoolListOption struct {
		PoolId string `help:"Pool Id"`
		Status string `help:"Status"`
	}

	shellutils.R(&ModelartsPoolListOption{}, "modelarts-pool-list", "List Modelarts Pool", func(cli *huawei.SRegion, args *ModelartsPoolListOption) error {
		if len(args.Status) > 0 {
			pools, err := cli.GetIModelartsPoolsWithStatus(args.Status)
			if err != nil {
				return err
			}
			printList(pools, len(pools), 0, 0, nil)
			return nil
		}
		pools, err := cli.GetIModelartsPools()
		if err != nil {
			return err
		}
		printList(pools, len(pools), 0, 0, nil)
		return nil
	})

	shellutils.R(&ModelartsPoolListOption{}, "modelarts-pool-detail", "List pool", func(cli *huawei.SRegion, args *ModelartsPoolListOption) error {
		pool, err := cli.GetIModelartsPoolById(args.PoolId)
		if err != nil {
			return err
		}
		printObject(pool)
		return nil
	})

	shellutils.R(&cloudprovider.ModelartsPoolCreateOption{}, "modelarts-pool-create", "Create Modelarts Pool", func(cli *huawei.SRegion, args *cloudprovider.ModelartsPoolCreateOption) error {
		res, err := cli.CreateIModelartsPool(args, nil)
		if err != nil {
			return err
		}
		printObject(res)
		return nil
	})

	shellutils.R(&ModelartsPoolListOption{}, "modelarts-pool-delete", "Delete Modelarts Pool", func(cli *huawei.SRegion, args *ModelartsPoolListOption) error {
		res, err := cli.DeletePool(args.PoolId)
		if err != nil {
			return err
		}
		printObject(res)
		return nil
	})

	shellutils.R(&ModelartsPoolListOption{}, "modelarts-pool-monitor", "Modelarts Pool Monitor", func(cli *huawei.SRegion, args *ModelartsPoolListOption) error {
		res, err := cli.MonitorPool(args.PoolId)
		if err != nil {
			return err
		}
		printList(res.Metrics, len(res.Metrics), 0, 0, nil)
		return nil
	})

	shellutils.R(&ModelartsPoolListOption{}, "modelarts-pool-by-id", "Modelarts Pool By Id", func(cli *huawei.SRegion, args *ModelartsPoolListOption) error {
		res, err := cli.GetIModelartsPoolById(args.PoolId)
		if err != nil {
			return err
		}
		printObject(res)
		return nil
	})
}
