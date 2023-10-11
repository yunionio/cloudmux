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
	"yunion.io/x/cloudmux/pkg/multicloud/ctyun"
)

func init() {
	type SEipListOptions struct {
		Status string `choices:"ACTIVE|DOWN|FREEZING|EXPIRED"`
	}
	shellutils.R(&SEipListOptions{}, "eip-list", "List eips", func(cli *ctyun.SRegion, args *SEipListOptions) error {
		eips, e := cli.GetEips(args.Status)
		if e != nil {
			return e
		}
		printList(eips, 0, 0, 0, nil)
		return nil
	})

	shellutils.R(&cloudprovider.SEip{}, "eip-create", "Create eip", func(cli *ctyun.SRegion, args *cloudprovider.SEip) error {
		eip, e := cli.CreateEip(args)
		if e != nil {
			return e
		}
		printObject(eip)
		return nil
	})

	type EipIdOptions struct {
		ID string
	}

	shellutils.R(&EipIdOptions{}, "eip-show", "Show eip", func(cli *ctyun.SRegion, args *EipIdOptions) error {
		eip, e := cli.GetEip(args.ID)
		if e != nil {
			return e
		}
		printObject(eip)
		return nil
	})

	shellutils.R(&EipIdOptions{}, "eip-delete", "Delete eip", func(cli *ctyun.SRegion, args *EipIdOptions) error {
		return cli.DeleteEip(args.ID)
	})

}
