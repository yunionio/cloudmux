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
	type DynamodbListOptions struct {
	}
	shellutils.R(&DynamodbListOptions{}, "dynamodb-table-list", "List kinesis stream", func(cli *aws.SRegion, args *DynamodbListOptions) error {
		tables, err := cli.ListTables()
		if err != nil {
			return err
		}
		fmt.Println(tables)
		return nil
	})

	type DynamodbNameOptions struct {
		NAME string
	}
	shellutils.R(&DynamodbNameOptions{}, "dynamodb-table-show", "Show kinesis stream", func(cli *aws.SRegion, args *DynamodbNameOptions) error {
		table, err := cli.DescribeTable(args.NAME)
		if err != nil {
			return err
		}
		printObject(table)
		return nil
	})

}
