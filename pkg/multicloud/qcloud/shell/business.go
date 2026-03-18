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

	"yunion.io/x/cloudmux/pkg/multicloud/qcloud"
)

func init() {
	type QcloudBillSummaryOptions struct {
		MONTH string `help:"Billing month, e.g. 2024-02"`
		Uin   string `help:"Root account UIN, optional"`
	}
	shellutils.R(&QcloudBillSummaryOptions{}, "bill-summary", "Query bill summary", func(cli *qcloud.SRegion, args *QcloudBillSummaryOptions) error {
		result, err := cli.GetClient().DescribeBillSummary(args.MONTH, args.Uin)
		if err != nil {
			return err
		}
		printObject(result)
		return nil
	})
}
