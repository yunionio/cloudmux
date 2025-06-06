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

	"yunion.io/x/cloudmux/pkg/multicloud/aws"
)

func init() {
	type LbbgIdOptions struct {
		ID string
	}
	shellutils.R(&LbbgIdOptions{}, "lb-lbbg-show", "Show loadbalancer backendgroup", func(cli *aws.SRegion, args *LbbgIdOptions) error {
		ret, err := cli.GetElbBackendgroup(args.ID)
		if err != nil {
			return err
		}
		printObject(ret)
		return nil
	})

	shellutils.R(&LbbgIdOptions{}, "lb-lbbg-delete", "Delete loadbalancer", func(cli *aws.SRegion, args *LbbgIdOptions) error {
		return cli.DeleteElbBackendGroup(args.ID)
	})

	shellutils.R(&LbbgIdOptions{}, "lb-lbbg-backend-list", "Delete loadbalancer", func(cli *aws.SRegion, args *LbbgIdOptions) error {
		backends, err := cli.GetELbBackends(args.ID)
		if err != nil {
			return err
		}
		printList(backends, 0, 0, 0, nil)
		return nil
	})

}
