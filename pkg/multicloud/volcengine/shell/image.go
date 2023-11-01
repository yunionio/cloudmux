// Copyright 2023 Yunion
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package shell

import (
	"fmt"

	"yunion.io/x/cloudmux/pkg/multicloud/volcengine"
	"yunion.io/x/pkg/util/shellutils"
)

func init() {
	type ImageListOptions struct {
		Visibility string   `chocies:"public|private|shared"`
		Id         []string `help:"Image ID"`
		Name       string   `help:"image name"`
	}
	shellutils.R(&ImageListOptions{}, "image-list", "List images", func(cli *volcengine.SRegion, args *ImageListOptions) error {
		images, err := cli.GetImages(args.Visibility, args.Id, args.Name)
		if err != nil {
			return err
		}
		printList(images, 0, 0, 0, nil)
		return nil
	})

	type ImageShowOptions struct {
		ID string `help:"image ID"`
	}
	shellutils.R(&ImageShowOptions{}, "image-show", "Show image", func(cli *volcengine.SRegion, args *ImageShowOptions) error {
		img, err := cli.GetImage(args.ID)
		if err != nil {
			return err
		}
		printObject(img)
		return nil
	})

	type ImageDeleteOptions struct {
		ID string `help:"ID or Name to delete"`
	}
	shellutils.R(&ImageDeleteOptions{}, "image-delete", "Delete image", func(cli *volcengine.SRegion, args *ImageDeleteOptions) error {
		return cli.DeleteImage(args.ID)
	})

	type ImageExportOptions struct {
		ID     string `help:"ID or Name to export"`
		BUCKET string `help:"Bucket name"`
	}

	shellutils.R(&ImageExportOptions{}, "image-export", "Export image", func(cli *volcengine.SRegion, args *ImageExportOptions) error {
		exist, err := cli.IBucketExist(args.BUCKET)
		if err != nil {
			return err
		}
		if !exist {
			return fmt.Errorf("not exist bucket %s", args.BUCKET)
		}
		taskId, err := cli.ExportImage(args.ID, args.BUCKET)
		if err != nil {
			return err
		}
		println(taskId)
		return nil
	})
}
