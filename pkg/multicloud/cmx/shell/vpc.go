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

// import (
// 	"yunion.io/x/cloudmux/pkg/cloudprovider"
// )
//
// func init() {
// 	type VpcListOptions struct {
// 		// Limit  int `help:"page size"`
// 		// Offset int `help:"page offset"`
// 	}
// 	RegionRList(&VpcListOptions{}, "vpc-list", "List vpcs", func(cli cloudprovider.ICloudRegion, args *VpcListOptions) error {
// 		vpcs, e := cli.GetIVpcs()
// 		if e != nil {
// 			return e
// 		}
// 		PrintList(vpcs, 0, 0, 0, []string{})
// 		return nil
// 	})
//
// 	type VpcCreateOptions struct {
// 		Name string
// 		Desc string
// 		CIDR string
// 	}
//
// 	RegionR(&VpcCreateOptions{}, "vpc-create", "Create vpc", func(cli cloudprovider.ICloudRegion, args *VpcCreateOptions) error {
// 		opts := cloudprovider.VpcCreateOptions{
// 			NAME: args.Name,
// 			CIDR: args.CIDR,
// 			Desc: args.Desc,
// 		}
// 		vpc, err := cli.CreateIVpc(&opts)
// 		if err != nil {
// 			return err
// 		}
// 		PrintObject(vpc)
// 		return nil
// 	})
//
// 	type VpcOptions struct {
// 		ID string `help:"VPC id"`
// 	}
//
// 	RegionR(&VpcOptions{}, "vpc-delete", "Delete vpc", func(cli cloudprovider.ICloudRegion, args *VpcOptions) error {
// 		vpc, err := cli.GetIVpcById(args.ID)
// 		if err != nil {
// 			return err
// 		}
// 		return vpc.Delete()
// 	})
//
// 	// type VpcMoveResourceGroup struct {
// 	// 	ResourceType    string `choices:"vpc|eip|bandwidthpackage" default:"vpc"`
// 	// 	ResourceGroupId string
// 	// 	ResourceId      string
// 	// }
// 	//
// 	// R(&VpcMoveResourceGroup{}, "vpc-mv-resource-group", "Delete vpc", func(cli *aliyun.SRegion, args *VpcMoveResourceGroup) error {
// 	// 	return cli.VpcMoveResourceGroup(args.ResourceType, args.ResourceGroupId, args.ResourceId)
// 	// })
//
// }
