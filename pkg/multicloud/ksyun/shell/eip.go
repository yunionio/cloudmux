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
	"yunion.io/x/cloudmux/pkg/cloudprovider"
	"yunion.io/x/cloudmux/pkg/multicloud/ksyun"

	"yunion.io/x/pkg/errors"
	"yunion.io/x/pkg/util/shellutils"
)

func init() {
	type EipListOptions struct {
		Id []string
	}
	shellutils.R(&EipListOptions{}, "eip-list", "list eips", func(cli *ksyun.SRegion, args *EipListOptions) error {
		res, err := cli.GetEips(args.Id)
		if err != nil {
			return errors.Wrap(err, "GetEips")
		}
		printList(res)
		return nil
	})

	shellutils.R(&EipListOptions{}, "eip-line-list", "list eip lines", func(cli *ksyun.SRegion, args *EipListOptions) error {
		res, err := cli.GetLines()
		if err != nil {
			return errors.Wrap(err, "GetLines")
		}
		printList(res)
		return nil
	})

	shellutils.R(&cloudprovider.SEip{}, "eip-create", "create eip", func(cli *ksyun.SRegion, args *cloudprovider.SEip) error {
		eip, err := cli.CreateEip(args)
		if err != nil {
			return errors.Wrap(err, "CreateEip")
		}
		printObject(eip)
		return nil
	})

	type EipIdOptions struct {
		ID string
	}

	shellutils.R(&EipIdOptions{}, "eip-delete", "delete eip", func(cli *ksyun.SRegion, args *EipIdOptions) error {
		return cli.DeallocateEIP(args.ID)
	})
}
