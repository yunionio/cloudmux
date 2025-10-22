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

	"yunion.io/x/cloudmux/pkg/multicloud/oceanbase"
)

func init() {
	type DBInstanceListOptions struct {
	}
	shellutils.R(&DBInstanceListOptions{}, "dbinstance-list", "show dbinstance", func(cli *oceanbase.SRegion, args *DBInstanceListOptions) error {
		instances, err := cli.GetDBInstances()
		if err != nil {
			return err
		}
		printList(instances)
		return nil
	})

	type DBInstanceIdOptions struct {
		ID string
	}
	shellutils.R(&DBInstanceIdOptions{}, "dbinstance-show", "show dbinstance", func(cli *oceanbase.SRegion, args *DBInstanceIdOptions) error {
		instance, err := cli.GetDBInstance(args.ID)
		if err != nil {
			return err
		}
		printObject(instance)
		return nil
	})

	shellutils.R(&DBInstanceIdOptions{}, "dbinstance-delete", "delete dbinstance", func(cli *oceanbase.SRegion, args *DBInstanceIdOptions) error {
		return cli.DeleteDBInstance(args.ID)
	})

	shellutils.R(&DBInstanceIdOptions{}, "dbinstance-start", "start dbinstance", func(cli *oceanbase.SRegion, args *DBInstanceIdOptions) error {
		return cli.StartDBInstance(args.ID)
	})

	shellutils.R(&DBInstanceIdOptions{}, "dbinstance-stop", "stop dbinstance", func(cli *oceanbase.SRegion, args *DBInstanceIdOptions) error {
		return cli.StopDBInstance(args.ID)
	})
}
