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
	"yunion.io/x/cloudmux/pkg/multicloud/baidu"
)

func init() {
	type BucketListOptions struct {
	}
	shellutils.R(&BucketListOptions{}, "bucket-list", "list buckets", func(cli *baidu.SRegion, args *BucketListOptions) error {
		buckets, err := cli.ListBuckets()
		if err != nil {
			return err
		}
		printList(buckets)
		return nil
	})

	type BucketNameOptions struct {
		BUCKET string
	}
	shellutils.R(&BucketNameOptions{}, "bucket-storage-class", "get bucket storage class", func(cli *baidu.SRegion, args *BucketNameOptions) error {
		storageClass, err := cli.GetBucketStorageClass(args.BUCKET)
		if err != nil {
			return err
		}
		fmt.Println(storageClass)
		return nil
	})

	shellutils.R(&BucketNameOptions{}, "bucket-delete", "delete bucket", func(cli *baidu.SRegion, args *BucketNameOptions) error {
		return cli.DeleteBucket(args.BUCKET)
	})

	shellutils.R(&BucketNameOptions{}, "bucket-acl-list", "list bucket acl", func(cli *baidu.SRegion, args *BucketNameOptions) error {
		acl, err := cli.GetBucketAcl(args.BUCKET)
		if err != nil {
			return err
		}
		printObject(acl)
		return nil
	})

	type BucketSetStorageClassOptions struct {
		BUCKET        string
		STORAGE_CLASS string
	}
	shellutils.R(&BucketSetStorageClassOptions{}, "bucket-set-storage-class", "set bucket storage class", func(cli *baidu.SRegion, args *BucketSetStorageClassOptions) error {
		return cli.SetBucketStorageClass(args.BUCKET, args.STORAGE_CLASS)
	})

	type BucketCreateOptions struct {
		BUCKET       string
		StorageClass string `choices:"STANDARD|COLDLINE"`
		Acl          string `choices:"private|public-read|public-read-write"`
	}
	shellutils.R(&BucketCreateOptions{}, "bucket-create", "create bucket", func(cli *baidu.SRegion, args *BucketCreateOptions) error {
		return cli.CreateBucket(args.BUCKET, args.StorageClass, args.Acl)
	})

	type BucketSetAclOptions struct {
		BUCKET string
		Acl    string `choices:"private|public-read|public-read-write"`
	}
	shellutils.R(&BucketSetAclOptions{}, "bucket-set-acl", "set bucket acl", func(cli *baidu.SRegion, args *BucketSetAclOptions) error {
		return cli.SetBucketAcl(args.BUCKET, cloudprovider.TBucketACLType(args.Acl))
	})

}
