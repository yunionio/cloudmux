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
	"net/http"
	"strings"

	"github.com/pkg/errors"
	"yunion.io/x/pkg/util/shellutils"

	"yunion.io/x/cloudmux/pkg/multicloud/azure"
)

func init() {
	type StorageAccountListOptions struct {
	}
	shellutils.R(&StorageAccountListOptions{}, "storage-account-list", "List storage account", func(cli *azure.SRegion, args *StorageAccountListOptions) error {
		accounts, err := cli.ListStorageAccounts()
		if err != nil {
			return err
		}
		printList(accounts, len(accounts), 0, 0, []string{})
		return nil
	})

	type StorageAccountOptions struct {
		ID string `help:"StorageAccount ID"`
	}

	shellutils.R(&StorageAccountOptions{}, "storage-account-delete", "Delete storage account", func(cli *azure.SRegion, args *StorageAccountOptions) error {
		return cli.DeleteStorageAccount(args.ID)
	})

	shellutils.R(&StorageAccountOptions{}, "storage-account-show", "Show storage account detail", func(cli *azure.SRegion, args *StorageAccountOptions) error {
		account, err := cli.GetStorageAccountDetail(args.ID)
		if err != nil {
			return err
		}
		printObject(account)
		return nil
	})

	shellutils.R(&StorageAccountOptions{}, "storage-account-key", "Get storage account key", func(cli *azure.SRegion, args *StorageAccountOptions) error {
		if key, err := cli.GetStorageAccountKey(args.ID); err != nil {
			return err
		} else {
			fmt.Printf("Key: %s", key)
			return nil
		}
	})

	type StorageAccountObjectOptions struct {
		StorageAccountOptions
		Prefix     string `help:"prefix of object to list"`
		Delimiter  string `help:"delimiter of object to list"`
		MaxResults int    `help:"max results of object to list"`
		Marker     string `help:"marker of object to list"`
	}

	shellutils.R(&StorageAccountObjectOptions{}, "storage-account-object-list", "List objects of a storage account", func(cli *azure.SRegion, args *StorageAccountObjectOptions) error {
		account, err := cli.GetStorageAccountDetail(args.ID)
		if err != nil {
			return err
		}
		objects, err := account.ListObjects(args.Prefix, args.Marker, args.Delimiter, args.MaxResults)
		if err != nil {
			return err
		}
		printList(objects.Objects, len(objects.Objects), 0, 0, nil)
		return nil
	})

	shellutils.R(&StorageAccountOptions{}, "storage-container-list", "Get list of containers of a storage account", func(cli *azure.SRegion, args *StorageAccountOptions) error {
		account, err := cli.GetStorageAccountDetail(args.ID)
		if err != nil {
			return err
		}
		containers, err := account.GetContainers()
		if err != nil {
			return err
		}
		printList(containers, len(containers), 0, 0, nil)
		return nil
	})

	type StorageAccountCreateContainerOptions struct {
		ACCOUNT   string `help:"storage account ID"`
		CONTAINER string `help:"name of container to create"`
	}
	shellutils.R(&StorageAccountCreateContainerOptions{}, "storage-container-create", "Create a container in a storage account", func(cli *azure.SRegion, args *StorageAccountCreateContainerOptions) error {
		account, err := cli.GetStorageAccountDetail(args.ACCOUNT)
		if err != nil {
			return err
		}
		container, err := account.CreateContainer(args.CONTAINER)
		if err != nil {
			return err
		}
		printObject(container)
		return nil
	})

	shellutils.R(&StorageAccountCreateContainerOptions{}, "storage-container-list-objects", "Create a container in a storage account", func(cli *azure.SRegion, args *StorageAccountCreateContainerOptions) error {
		account, err := cli.GetStorageAccountDetail(args.ACCOUNT)
		if err != nil {
			return err
		}
		container, err := account.GetContainer(args.CONTAINER)
		if err != nil {
			return err
		}
		blobs, err := container.ListAllFiles()
		if err != nil {
			return err
		}
		printList(blobs, len(blobs), 0, 0, nil)
		return nil
	})

	type StorageAccountUploadOptions struct {
		ACCOUNT   string `help:"storage account ID"`
		CONTAINER string `help:"name of container to create"`
		FILE      string `help:"local file to upload"`
	}
	shellutils.R(&StorageAccountUploadOptions{}, "storage-container-upload", "Upload a container in a storage account", func(cli *azure.SRegion, args *StorageAccountUploadOptions) error {
		account, err := cli.GetStorageAccountDetail(args.ACCOUNT)
		if err != nil {
			return err
		}
		container, err := account.GetContainer(args.CONTAINER)
		if err != nil {
			return err
		}
		url, err := container.UploadFile(args.FILE, nil)
		if err != nil {
			return err
		}
		fmt.Println(url)
		return nil
	})

	type StorageAccountCreateOptions struct {
		NAME string `help:"StorageAccount NAME"`
	}

	shellutils.R(&StorageAccountCreateOptions{}, "storage-account-create", "Create a storage account", func(cli *azure.SRegion, args *StorageAccountCreateOptions) error {
		if account, err := cli.CreateStorageAccount(args.NAME); err != nil {
			return err
		} else {
			printObject(account)
			return nil
		}
	})

	type StorageAccountCheckeOptions struct {
	}

	shellutils.R(&StorageAccountCheckeOptions{}, "storage-uniq-name", "Get a uniqel storage account name", func(cli *azure.SRegion, args *StorageAccountCheckeOptions) error {
		uniqName, err := cli.GetUniqStorageAccountName()
		if err != nil {
			return err
		}
		fmt.Println(uniqName)
		return nil
	})

	type SStorageAccountSkuOptions struct {
	}
	shellutils.R(&SStorageAccountSkuOptions{}, "storage-account-skus", "List skus of storage account", func(cli *azure.SRegion, args *SStorageAccountSkuOptions) error {
		skus, err := cli.GetStorageAccountSkus()
		if err != nil {
			return err
		}
		printList(skus, 0, 0, 0, nil)
		return nil
	})

	type StorageAccountSetObjectMetaOptions struct {
		ACCOUNT string   `help:"storage account ID"`
		OBJECT  string   `help:"name of object to set meta"`
		META    []string `help:"meta to set"`
	}
	shellutils.R(&StorageAccountSetObjectMetaOptions{}, "storage-object-set-meta", "Set meta of a object in a storage account", func(cli *azure.SRegion, args *StorageAccountSetObjectMetaOptions) error {
		account, err := cli.GetStorageAccountDetail(args.ACCOUNT)
		if err != nil {
			return err
		}
		meta := http.Header{}
		for _, m := range args.META {
			parts := strings.SplitN(m, ":", 2)
			if len(parts) != 2 {
				return errors.Errorf("invalid meta: %s", m)
			}
			meta.Set(parts[0], parts[1])
		}
		err = account.SetObjectMeta(context.Background(), args.OBJECT, meta)
		if err != nil {
			return err
		}
		return nil
	})

	type StorageAccountGetObjectMetaOptions struct {
		ACCOUNT string `help:"storage account ID"`
		OBJECT  string `help:"name of object to get meta"`
	}
	shellutils.R(&StorageAccountGetObjectMetaOptions{}, "storage-object-get-meta", "Get meta of a object in a storage account", func(cli *azure.SRegion, args *StorageAccountGetObjectMetaOptions) error {
		account, err := cli.GetStorageAccountDetail(args.ACCOUNT)
		if err != nil {
			return err
		}
		meta, err := account.GetObjectMeta(args.OBJECT)
		if err != nil {
			return err
		}
		printObject(meta)
		return nil
	})
}
