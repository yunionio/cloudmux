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

	"yunion.io/x/cloudmux/pkg/cloudprovider"
	"yunion.io/x/cloudmux/pkg/multicloud/ksyun"
)

func init() {
	type BucketListOptions struct {
	}
	shellutils.R(&BucketListOptions{}, "bucket-list", "List buckets", func(cli *ksyun.SRegion, args *BucketListOptions) error {
		buckets, err := cli.GetIBuckets()
		if err != nil {
			return err
		}
		printList(buckets)
		return nil
	})

	type BucketAclOptions struct {
		BUCKET string `help:"name of bucket"`
	}
	shellutils.R(&BucketAclOptions{}, "bucket-acl", "Show bucket ACL", func(cli *ksyun.SRegion, args *BucketAclOptions) error {
		fmt.Println(cli.GetBucketAcl(args.BUCKET))
		return nil
	})

	type BucketAclSetOptions struct {
		BUCKET string `help:"name of bucket"`
		ACL    string `help:"ACL string" choices:"private|public-read|public-read-write"`
	}
	shellutils.R(&BucketAclSetOptions{}, "bucket-set-acl", "Set bucket ACL", func(cli *ksyun.SRegion, args *BucketAclSetOptions) error {
		return cli.SetBucketAcl(args.BUCKET, cloudprovider.TBucketACLType(args.ACL))
	})
}
