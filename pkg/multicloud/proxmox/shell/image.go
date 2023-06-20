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
	"os"

	"yunion.io/x/pkg/util/shellutils"

	"yunion.io/x/cloudmux/pkg/multicloud/proxmox"
)

func init() {
	type ImageListOptions struct {
		NODE    string
		STORAGE string
	}

	shellutils.R(&ImageListOptions{}, "image-list", "list images", func(cli *proxmox.SRegion, args *ImageListOptions) error {
		images, err := cli.GetImages(args.NODE, args.STORAGE)
		if err != nil {
			return err
		}
		printList(images, 0, 0, 0, []string{})
		return nil
	})

	type ImageUploadOptions struct {
		NODE     string
		STORAGE  string
		FILENAME string
	}

	shellutils.R(&ImageUploadOptions{}, "image-upload", "upload image", func(cli *proxmox.SRegion, args *ImageUploadOptions) error {
		file, err := os.Open(args.FILENAME)
		if err != nil {
			return err
		}
		defer file.Close()
		image, err := cli.UploadImage(args.NODE, args.STORAGE, args.FILENAME, file)
		if err != nil {
			return err
		}
		printObject(image)
		return nil
	})

}
