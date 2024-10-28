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

	"yunion.io/x/cloudmux/pkg/multicloud/cephfs"
)

func init() {
	type NasListOption struct {
	}
	shellutils.R(&NasListOption{}, "cephfs-list", "List cephfs", func(cli *cephfs.SCephFSClient, args *NasListOption) error {
		fs, err := cli.GetCephFSs()
		if err != nil {
			return err
		}
		printList(fs)
		return nil
	})

	type NasDirListOption struct {
		FS_ID string
	}
	shellutils.R(&NasDirListOption{}, "cephfs-dir-list", "List cephfs", func(cli *cephfs.SCephFSClient, args *NasDirListOption) error {
		dirs, err := cli.GetCephDirs(args.FS_ID)
		if err != nil {
			return err
		}
		printList(dirs)
		return nil
	})

	type NasDirCreateOption struct {
		FS_ID string
		PATH  string
	}
	shellutils.R(&NasDirCreateOption{}, "cephfs-dir-create", "Create cephfs dir", func(cli *cephfs.SCephFSClient, args *NasDirCreateOption) error {
		return cli.CreateDir(args.FS_ID, args.PATH)
	})

	shellutils.R(&NasDirCreateOption{}, "cephfs-dir-delete", "Delete cephfs dir", func(cli *cephfs.SCephFSClient, args *NasDirCreateOption) error {
		return cli.DeleteDir(args.FS_ID, args.PATH)
	})

	type NasDirSetQuotaOption struct {
		FS_ID    string
		PATH     string
		MaxMb    int64
		MaxFiles int64
	}
	shellutils.R(&NasDirSetQuotaOption{}, "cephfs-dir-set-quota", "Set cephfs dir quota", func(cli *cephfs.SCephFSClient, args *NasDirSetQuotaOption) error {
		return cli.SetQuota(args.FS_ID, args.PATH, args.MaxMb, args.MaxFiles)
	})

}
