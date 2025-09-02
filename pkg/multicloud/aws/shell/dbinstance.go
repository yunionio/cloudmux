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
	"yunion.io/x/log"
	"yunion.io/x/pkg/util/shellutils"

	"yunion.io/x/cloudmux/pkg/cloudprovider"
	"yunion.io/x/cloudmux/pkg/multicloud/aws"
)

func init() {
	type DBInstanceListOptions struct {
		Id string
	}
	shellutils.R(&DBInstanceListOptions{}, "dbinstance-list", "List rds intances", func(cli *aws.SRegion, args *DBInstanceListOptions) error {
		instances, err := cli.GetDBInstances(args.Id)
		if err != nil {
			return err
		}
		printList(instances, 0, 0, 0, []string{})
		return nil
	})

	type DBInstanceIdOptions struct {
		ID string
	}

	shellutils.R(&DBInstanceIdOptions{}, "dbinstance-show", "Show rds intance", func(cli *aws.SRegion, args *DBInstanceIdOptions) error {
		instance, err := cli.GetDBInstance(args.ID)
		if err != nil {
			return err
		}
		printObject(instance)
		return nil
	})

	shellutils.R(&DBInstanceIdOptions{}, "dbinstance-tags-list", "Show rds intance tags", func(cli *aws.SRegion, args *DBInstanceIdOptions) error {
		instance, err := cli.ListRdsResourceTags(args.ID)
		if err != nil {
			return err
		}
		printObject(instance)
		return nil
	})

	type SDBInstanceUpdateOptions struct {
		ID string
		cloudprovider.SDBInstanceUpdateOptions
	}
	shellutils.R(&SDBInstanceUpdateOptions{}, "dbinstance-update", "Show rds intance tags", func(cli *aws.SRegion, args *SDBInstanceUpdateOptions) error {
		err := cli.Update(args.ID, args.SDBInstanceUpdateOptions)
		if err != nil {
			log.Errorln(err)
		}
		return err
	})

	type SDBInstanceClusterListOptions struct {
		Id string
	}
	shellutils.R(&SDBInstanceClusterListOptions{}, "dbinstance-cluster-list", "List rds intance clusters", func(cli *aws.SRegion, args *SDBInstanceClusterListOptions) error {
		clusters, err := cli.GetDBInstanceClusters(args.Id)
		if err != nil {
			return err
		}
		printList(clusters, 0, 0, 0, []string{})
		return nil
	})

	type SDBInstanceGlobalClusterListOptions struct {
		Id string
	}
	shellutils.R(&SDBInstanceGlobalClusterListOptions{}, "dbinstance-global-cluster-list", "List rds intance global clusters", func(cli *aws.SRegion, args *SDBInstanceGlobalClusterListOptions) error {
		clusters, err := cli.GetDBInstanceGlobalClusters(args.Id)
		if err != nil {
			return err
		}
		printList(clusters, 0, 0, 0, []string{})
		return nil
	})
}
