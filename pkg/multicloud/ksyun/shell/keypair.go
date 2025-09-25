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
	"yunion.io/x/cloudmux/pkg/multicloud/ksyun"

	"yunion.io/x/pkg/util/shellutils"
)

func init() {
	type KeyPairListOptions struct {
	}
	shellutils.R(&KeyPairListOptions{}, "keypair-list", "List keypairs", func(cli *ksyun.SRegion, args *KeyPairListOptions) error {
		keypairs, err := cli.GetClient().GetKeypairs()
		if err != nil {
			return err
		}
		printList(keypairs)
		return nil
	})

	type KeyPairCreateOptions struct {
		Name      string
		PublicKey string
	}
	shellutils.R(&KeyPairCreateOptions{}, "keypair-create", "Create keypair", func(cli *ksyun.SRegion, args *KeyPairCreateOptions) error {
		keypair, err := cli.GetClient().CreateKeypair(args.Name, args.PublicKey)
		if err != nil {
			return err
		}
		printObject(keypair)
		return nil
	})
}
