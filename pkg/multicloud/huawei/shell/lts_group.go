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

	"yunion.io/x/cloudmux/pkg/multicloud/huawei"
)

func init() {
	type LtsGroupListOptions struct {
	}
	shellutils.R(&LtsGroupListOptions{}, "lts-group-list", "List lts groups", func(cli *huawei.SRegion, args *LtsGroupListOptions) error {
		ret, err := cli.ListLtsGroups()
		if err != nil {
			return err
		}
		printList(ret, 0, 0, 0, []string{})
		return nil
	})

	type LtsStreamListOptions struct {
		Group string
	}
	shellutils.R(&LtsStreamListOptions{}, "lts-stream-list", "List lts streams", func(cli *huawei.SRegion, args *LtsStreamListOptions) error {
		if len(args.Group) > 0 {
			ret, err := cli.ListLtsStreamsByGroup(args.Group)
			if err != nil {
				return err
			}
			printList(ret, 0, 0, 0, []string{})
			return nil
		}
		ret, err := cli.ListLtsStreams()
		if err != nil {
			return err
		}
		printList(ret, 0, 0, 0, []string{})
		return nil
	})

}
