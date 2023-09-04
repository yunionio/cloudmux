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
	type InternetGatewayCreateOptions struct {
	}
	shellutils.R(&InternetGatewayCreateOptions{}, "igw-create", "Create igw", func(cli *aws.SRegion, args *InternetGatewayCreateOptions) error {
		ret, err := cli.CreateIgw()
		if err != nil {
			return err
		}
		printObject(ret)
		return nil
	})

	type InternetGatewayListOptions struct {
		VpcId string
	}

	shellutils.R(&InternetGatewayListOptions{}, "igw-list", "List igw", func(cli *aws.SRegion, args *InternetGatewayListOptions) error {
		igws, err := cli.GetInternetGateways(args.VpcId)
		if err != nil {
			return err
		}
		printList(igws, 0, 0, 0, nil)
		return nil
	})

	type InternetGatewayIdOptions struct {
		ID string
	}

	shellutils.R(&InternetGatewayIdOptions{}, "igw-delete", "Delete igw", func(cli *aws.SRegion, args *InternetGatewayIdOptions) error {
		return cli.DeleteInternetGateway(args.ID)
	})

	type InternetGatewayDetachOptions struct {
		VPC_ID string
		IGW_ID string
	}

	shellutils.R(&InternetGatewayDetachOptions{}, "igw-detach", "Detach igw", func(cli *aws.SRegion, args *InternetGatewayDetachOptions) error {
		return cli.DetachInternetGateway(args.VPC_ID, args.IGW_ID)
	})

}
