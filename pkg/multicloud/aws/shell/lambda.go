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

	"yunion.io/x/cloudmux/pkg/multicloud/aws"
)

func init() {
	type LambdaListOptions struct {
	}
	shellutils.R(&LambdaListOptions{}, "lambda-function-list", "List lambda functions", func(cli *aws.SRegion, args *LambdaListOptions) error {
		functions, err := cli.ListFunctions()
		if err != nil {
			return err
		}
		printList(functions, 0, 0, 0, []string{})
		return nil
	})

	type LambdaProvisionedOptions struct {
		NAME    string
		VERSION string
	}
	shellutils.R(&LambdaProvisionedOptions{}, "lambda-function-provisioned-show", "Show lambda function provisioned", func(cli *aws.SRegion, args *LambdaProvisionedOptions) error {
		ret, err := cli.GetProvisionedConcurrencyConfig(args.NAME, args.VERSION)
		if err != nil {
			return err
		}
		printObject(ret)
		return nil
	})

}
