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

	"yunion.io/x/cloudmux/pkg/multicloud/ctyun"
)

func init() {
	type SKeypairOptions struct {
		Name string
	}
	shellutils.R(&SKeypairOptions{}, "keypair-list", "List keypair", func(cli *ctyun.SRegion, args *SKeypairOptions) error {
		keypairs, err := cli.GetKeypairs(args.Name)
		if err != nil {
			return err
		}
		printList(keypairs, 0, 0, 0, nil)
		return nil
	})

	type SKeypairImportOptions struct {
		NAME       string
		PUBLIC_KEY string
	}

	shellutils.R(&SKeypairImportOptions{}, "keypair-import", "Import keypair", func(cli *ctyun.SRegion, args *SKeypairImportOptions) error {
		keypair, err := cli.ImportKeypair(args.NAME, args.PUBLIC_KEY)
		if err != nil {
			return err
		}
		printObject(keypair)
		return nil
	})

}
