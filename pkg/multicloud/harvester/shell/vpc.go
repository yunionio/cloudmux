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
	"yunion.io/x/cloudmux/pkg/multicloud/harvester"

	"yunion.io/x/pkg/errors"
	"yunion.io/x/pkg/util/shellutils"
)

func init() {
	type VpcListOptions struct {
	}
	shellutils.R(&VpcListOptions{}, "vpc-list", "list vpc", func(cli *harvester.SRegion, args *VpcListOptions) error {
		res, err := cli.GetVpcs()
		if err != nil {
			return errors.Wrap(err, "GetVpcs")
		}
		printList(res)
		return nil
	})
}
