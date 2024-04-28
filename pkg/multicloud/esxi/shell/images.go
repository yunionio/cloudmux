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
	"context"
	"fmt"
	"path/filepath"

	"yunion.io/x/log"
	"yunion.io/x/pkg/errors"
	"yunion.io/x/pkg/util/shellutils"

	"yunion.io/x/cloudmux/pkg/multicloud/esxi"
)

func init() {
	type ImageListOptions struct {
	}
	shellutils.R(&ImageListOptions{}, "image-list", "List all images", func(cli *esxi.SESXiClient, args *ImageListOptions) error {
		caches, err := cli.GetIStoragecaches()
		if err != nil {
			return err
		}
		for i := range caches {
			images, err := caches[i].GetICloudImages()
			if err != nil {
				return err
			}
			fmt.Printf("storagecache: %s\n", caches[i].GetName())
			printList(images, nil)
		}
		return nil
	})

	type ImageUploadOptions struct {
		DATA_STORE string
		FILE       string
		Dir        string `default:"image_cache"`
	}
	shellutils.R(&ImageUploadOptions{}, "image-upload", "Upload image", func(cli *esxi.SESXiClient, args *ImageUploadOptions) error {
		dss, err := cli.GetIStorages()
		if err != nil {
			return errors.Wrapf(err, "GetDatacenter")
		}
		for i := range dss {
			if dss[i].GetName() == args.DATA_STORE {
				ds := dss[i].(*esxi.SDatastore)
				path := fmt.Sprintf("%s/%s", args.Dir, filepath.Base(args.FILE))
				log.Infof("upload file %s -> %s", args.FILE, path)
				return ds.ImportISO(context.Background(), args.FILE, path)
			}
		}
		return nil
	})

}
