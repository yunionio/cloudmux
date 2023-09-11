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
	"fmt"

	"yunion.io/x/pkg/util/shellutils"

	"yunion.io/x/cloudmux/pkg/multicloud/aws"
)

func init() {
	type KeyPairListOptions struct {
		Name   string
		Finger string
	}
	shellutils.R(&KeyPairListOptions{}, "keypair-list", "List keypairs", func(cli *aws.SRegion, args *KeyPairListOptions) error {
		keypairs, err := cli.GetKeypairs(args.Finger, args.Name)
		if err != nil {
			return err
		}
		printList(keypairs, 0, 0, 0, []string{})
		return nil
	})

	type KeyPairImportOptions struct {
		NAME   string `help:"Name of new keypair"`
		PUBKEY string `help:"Public key string"`
	}
	shellutils.R(&KeyPairImportOptions{}, "keypair-import", "Import a keypair", func(cli *aws.SRegion, args *KeyPairImportOptions) error {
		keypair, err := cli.ImportKeypair(args.NAME, args.PUBKEY)
		if err != nil {
			return err
		}
		printObject(keypair)
		return nil
	})

	type KeyPairSyncOptions struct {
		PUBKEY string `help:"Public key string"`
	}
	shellutils.R(&KeyPairSyncOptions{}, "keypair-sync", "Sync a keypair", func(cli *aws.SRegion, args *KeyPairSyncOptions) error {
		key, err := cli.SyncKeypair(args.PUBKEY)
		if err != nil {
			return err
		}
		fmt.Println(key)
		return nil
	})

}
