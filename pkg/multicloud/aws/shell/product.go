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
	"strings"

	"yunion.io/x/pkg/util/shellutils"

	"yunion.io/x/cloudmux/pkg/multicloud/aws"
)

func init() {
	type ProductListOptions struct {
		ServiceCode string `default:"AmazonEC2" choices:"AmazonEC2|AmazonElastiCache"`
		NextToken   string
		Filters     []string
	}
	shellutils.R(&ProductListOptions{}, "product-list", "List product", func(cli *aws.SRegion, args *ProductListOptions) error {
		filters := []aws.ProductFilter{}
		for _, filter := range args.Filters {
			info := strings.Split(filter, "=")
			if len(info) == 2 {
				filters = append(filters, aws.ProductFilter{
					Type:  "TERM_MATCH",
					Field: info[0],
					Value: info[1],
				})
			}
		}
		products, _, err := cli.GetProducts(args.ServiceCode, filters, args.NextToken)
		if err != nil {
			return err
		}
		printList(products, 0, 0, 0, []string{})
		return nil
	})
}
