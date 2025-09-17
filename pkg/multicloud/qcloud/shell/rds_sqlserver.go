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

	"yunion.io/x/cloudmux/pkg/multicloud/qcloud"
)

func init() {
	type SQLServerSQLProductListOptions struct {
		ZONE string
	}
	shellutils.R(&SQLServerSQLProductListOptions{}, "sqlserver-product-list", "List sql server products", func(cli *qcloud.SRegion, args *SQLServerSQLProductListOptions) error {
		products, err := cli.DescribeSqlServerProductConfig(args.ZONE)
		if err != nil {
			return errors.Wrapf(err, "DescribeProductConfig")
		}
		printList(products, 0, 0, 0, nil)
		return nil
	})

	type SSQLServerSQLSkuListOptions struct {
	}

	shellutils.R(&SSQLServerSQLSkuListOptions{}, "sqlserver-sku-list", "List sqlserver skus", func(cli *qcloud.SRegion, args *SSQLServerSQLSkuListOptions) error {
		skus, err := cli.ListSQLServerSkus()
		if err != nil {
			return errors.Wrapf(err, "ListSQLServerSkus")
		}
		printList(skus, 0, 0, 0, nil)
		return nil
	})

	type SSQLServerSQLInstanceListOptions struct {
		Id string
	}
	shellutils.R(&SSQLServerSQLInstanceListOptions{}, "sqlserver-instance-list", "List sqlserver instances", func(cli *qcloud.SRegion, args *SSQLServerSQLInstanceListOptions) error {
		instances, err := cli.GetSQLServers(args.Id)
		if err != nil {
			return errors.Wrapf(err, "GetSQLServers")
		}
		printList(instances, 0, 0, 0, nil)
		return nil
	})

	type SQLServerIdOptions struct {
		ID string
	}
	shellutils.R(&SQLServerIdOptions{}, "sqlserver-instance-show", "Show sqlserver", func(cli *qcloud.SRegion, args *SQLServerIdOptions) error {
		instance, err := cli.GetSQLServer(args.ID)
		if err != nil {
			return err
		}
		printObject(instance)
		return nil
	})

	shellutils.R(&SQLServerIdOptions{}, "sqlserver-instance-delete", "Delete sqlserver", func(cli *qcloud.SRegion, args *SQLServerIdOptions) error {
		return cli.DeleteSQLServer(args.ID)
	})

	shellutils.R(&SQLServerIdOptions{}, "sqlserver-instance-delete-in-recycle-bin", "Delete sqlserver in recycle bin", func(cli *qcloud.SRegion, args *SQLServerIdOptions) error {
		return cli.DeleteSQLServerInRecycleBin(args.ID)
	})

}
