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
		Content string `default:"iso"`
	}

	shellutils.R(&ImageListOptions{}, "image-list", "list images", func(cli *proxmox.SProxmoxClient, args *ImageListOptions) error {
		images, err := cli.GetImages(args.NODE, args.STORAGE, args.Content)
		if err != nil {
			return err
		}
		printList(images, 0, 0, 0, []string{})
		return nil
	})

	type TemplateImageListOptions struct {
		Node string
	}
	shellutils.R(&TemplateImageListOptions{}, "template-image-list", "list template images", func(cli *proxmox.SProxmoxClient, args *TemplateImageListOptions) error {
		images, err := cli.GetTemplateImages(args.Node)
		if err != nil {
			return err
		}
		printList(images, 0, 0, 0, []string{})
		return nil
	})

	type ImportImageOptions struct {
		HOST    string
		STORAGE string
		IMAGEID string
		FORMAT  string
		FILE    string
	}
	shellutils.R(&ImportImageOptions{}, "image-upload", "upload image", func(cli *proxmox.SProxmoxClient, args *ImportImageOptions) error {
		file, err := os.Open(args.FILE)
		if err != nil {
			return err
		}
		defer file.Close()
		err = cli.ImportImage(args.HOST, args.STORAGE, args.IMAGEID, args.FORMAT, file)
		if err != nil {
			return err
		}
		return nil
	})
}
