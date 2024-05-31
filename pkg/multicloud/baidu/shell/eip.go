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

	"yunion.io/x/cloudmux/pkg/cloudprovider"
	"yunion.io/x/cloudmux/pkg/multicloud/baidu"
)

func init() {
	type EipListOptions struct {
		InstanceId string
	}
	shellutils.R(&EipListOptions{}, "eip-list", "list eips", func(cli *baidu.SRegion, args *EipListOptions) error {
		eips, err := cli.GetEips(args.InstanceId)
		if err != nil {
			return err
		}
		printList(eips)
		return nil
	})
	type EipIdOptions struct {
		ID string `help:"ID of eip to show"`
	}
	shellutils.R(&EipIdOptions{}, "eip-show", "list eips", func(cli *baidu.SRegion, args *EipIdOptions) error {
		eip, err := cli.GetEip(args.ID)
		if err != nil {
			return errors.Wrap(err, "GetEip")
		}
		printObject(eip)
		return nil
	})

	shellutils.R(&EipIdOptions{}, "eip-delete", "delete eip", func(cli *baidu.SRegion, args *EipIdOptions) error {
		return cli.DeleteEip(args.ID)
	})

	shellutils.R(&cloudprovider.SEip{}, "eip-create", "create eip", func(cli *baidu.SRegion, args *cloudprovider.SEip) error {
		eip, err := cli.CreateEip(args)
		if err != nil {
			return err
		}
		printObject(eip)
		return nil
	})

}
