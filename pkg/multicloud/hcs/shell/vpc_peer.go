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

	"yunion.io/x/cloudmux/pkg/cloudprovider"
	"yunion.io/x/cloudmux/pkg/multicloud/hcs"
)

func init() {
	type VpcPeerListOptions struct {
	}
	shellutils.R(&VpcPeerListOptions{}, "vpc-peer-list", "List vpc peers", func(cli *hcs.SRegion, args *VpcPeerListOptions) error {
		peers, err := cli.GetVpcPeerings()
		if err != nil {
			return nil
		}
		printList(peers, 0, 0, 0, nil)
		return nil
	})

	type VpcPeerIdOption struct {
		ID string
	}

	shellutils.R(&VpcPeerIdOption{}, "vpc-peer-show", "Show vpc peer", func(cli *hcs.SRegion, args *VpcPeerIdOption) error {
		ret, err := cli.GetVpcPeering(args.ID)
		if err != nil {
			return err
		}
		printObject(ret)
		return nil
	})

	shellutils.R(&VpcPeerIdOption{}, "vpc-peer-delete", "Delete vpc peer", func(cli *hcs.SRegion, args *VpcPeerIdOption) error {
		return cli.DeleteVpcPeering(args.ID)
	})

	type PeerCreateOptions struct {
		NAME    string
		VPC_ID  string
		PEER_ID string
	}

	shellutils.R(&PeerCreateOptions{}, "vpc-peer-create", "Create vpc peer", func(cli *hcs.SRegion, args *PeerCreateOptions) error {
		ret, err := cli.CreateVpcPeering(args.VPC_ID, &cloudprovider.VpcPeeringConnectionCreateOptions{
			Name:      args.NAME,
			PeerVpcId: args.PEER_ID,
		})
		if err != nil {
			return err
		}
		printObject(ret)
		return nil
	})

}
