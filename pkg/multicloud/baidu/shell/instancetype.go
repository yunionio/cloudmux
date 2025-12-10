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
	"yunion.io/x/cloudmux/pkg/multicloud/baidu"

	"yunion.io/x/pkg/errors"
	"yunion.io/x/pkg/util/shellutils"
)

func init() {
	type InstanceTypeListOptions struct {
		ZoneName string
		Specs    []string
		SpecIds  []string
	}
	shellutils.R(&InstanceTypeListOptions{}, "instance-type-list", "list instance types", func(cli *baidu.SRegion, args *InstanceTypeListOptions) error {
		res, err := cli.GetInstanceTypes(args.ZoneName, args.Specs, args.SpecIds)
		if err != nil {
			return errors.Wrap(err, "GetInstanceTypes")
		}
		printList(res)
		return nil
	})

	type InstanceTypePriceOptions struct {
		ZoneName string
		SpecId   string
		Spec     string
	}
	shellutils.R(&InstanceTypePriceOptions{}, "instance-type-price", "get instance type price", func(cli *baidu.SRegion, args *InstanceTypePriceOptions) error {
		res, err := cli.GetInstancePrice(args.ZoneName, args.SpecId, args.Spec)
		if err != nil {
			return errors.Wrap(err, "GetInstancePrice")
		}
		printObject(res)
		return nil
	})
}
