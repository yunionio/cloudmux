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

	type PriceListOptions struct {
		ServiceCode string `default:"AmazonEC2" choices:"AmazonEC2|AmazonElastiCache"`
	}
	shellutils.R(&PriceListOptions{}, "price-list", "List price", func(cli *aws.SRegion, args *PriceListOptions) error {
		prices, err := cli.ListPriceLists(args.ServiceCode)
		if err != nil {
			return err
		}
		printList(prices, 0, 0, 0, nil)
		return nil
	})

	type PriceFileOptions struct {
		ARN string
	}

	shellutils.R(&PriceFileOptions{}, "price-file-url", "List price", func(cli *aws.SRegion, args *PriceFileOptions) error {
		url, err := cli.GetPriceListFileUrl(args.ARN)
		if err != nil {
			return err
		}
		fmt.Println(url)
		return nil
	})
}
