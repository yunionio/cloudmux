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

	"yunion.io/x/pkg/util/shellutils"

	"yunion.io/x/cloudmux/pkg/cloudprovider"
	"yunion.io/x/cloudmux/pkg/multicloud/aws"
)

func init() {
	type InstanceListOptions struct {
		Id      []string `help:"IDs of instances to show"`
		ImageId string
		Zone    string `help:"Zone ID"`
	}
	shellutils.R(&InstanceListOptions{}, "instance-list", "List intances", func(cli *aws.SRegion, args *InstanceListOptions) error {
		instances, err := cli.GetInstances(args.Zone, args.ImageId, args.Id)
		if err != nil {
			return err
		}
		printList(instances, 0, 0, 0, []string{})
		return nil
	})

	type InstanceAttributeShowOptions struct {
		ID string
		aws.InstanceAttributeInput
	}

	shellutils.R(&InstanceAttributeShowOptions{}, "instance-attribute-show", "Show intance attribute", func(cli *aws.SRegion, args *InstanceAttributeShowOptions) error {
		attr, err := cli.DescribeInstanceAttribute(args.ID, &args.InstanceAttributeInput)
		if err != nil {
			return err
		}
		printObject(attr)
		return nil
	})

	type InstanceAttributeChangeOptions struct {
		ID string
		aws.SInstanceAttr
	}

	shellutils.R(&InstanceAttributeChangeOptions{}, "instance-attribute-change", "Change intance attribute", func(cli *aws.SRegion, args *InstanceAttributeChangeOptions) error {
		return cli.ModifyInstanceAttribute(args.ID, &args.SInstanceAttr)
	})

	type InstanceDiskOperationOptions struct {
		ID   string `help:"instance ID"`
		DISK string `help:"disk ID"`
	}

	type InstanceDiskAttachOptions struct {
		ID     string `help:"instance ID"`
		DISK   string `help:"disk ID"`
		DEVICE string `help:"disk device name. eg. /dev/sdb"`
	}

	shellutils.R(&InstanceDiskAttachOptions{}, "instance-attach-disk", "Attach a disk to instance", func(cli *aws.SRegion, args *InstanceDiskAttachOptions) error {
		return cli.AttachDisk(args.ID, args.DISK, args.DEVICE)
	})

	shellutils.R(&InstanceDiskOperationOptions{}, "instance-detach-disk", "Detach a disk to instance", func(cli *aws.SRegion, args *InstanceDiskOperationOptions) error {
		return cli.DetachDisk(args.ID, args.DISK)
	})

	type InstanceOperationOptions struct {
		ID string `help:"instance ID"`
	}
	shellutils.R(&InstanceOperationOptions{}, "instance-start", "Start a instance", func(cli *aws.SRegion, args *InstanceOperationOptions) error {
		return cli.StartVM(args.ID)
	})

	shellutils.R(&InstanceOperationOptions{}, "instance-password", "Show instance passowrd", func(cli *aws.SRegion, args *InstanceOperationOptions) error {
		password, err := cli.GetPasswordData(args.ID)
		if err != nil {
			return err
		}
		fmt.Println(password)
		return nil
	})

	type InstanceUpdateNameOptions struct {
		ID string `help:"Instance ID"`
		cloudprovider.SInstanceUpdateOptions
	}

	shellutils.R(&InstanceUpdateNameOptions{}, "instance-update", "Set new name", func(cli *aws.SRegion, args *InstanceUpdateNameOptions) error {
		err := cli.UpdateVM(args.ID, args.SInstanceUpdateOptions)
		if err != nil {
			return err
		}
		return nil
	})

	type InstanceStopOptions struct {
		ID    string `help:"instance ID"`
		Force bool   `help:"Force stop instance"`
	}
	shellutils.R(&InstanceStopOptions{}, "instance-stop", "Stop a instance", func(cli *aws.SRegion, args *InstanceStopOptions) error {
		return cli.StopVM(args.ID, args.Force)
	})
	shellutils.R(&InstanceOperationOptions{}, "instance-delete", "Delete a instance", func(cli *aws.SRegion, args *InstanceOperationOptions) error {
		return cli.DeleteVM(args.ID)
	})

	type InstanceRebuildRootOptions struct {
		ID    string `help:"instance ID"`
		Image string `help:"Image ID"`
		Size  int    `help:"system disk size in GB"`
	}

	shellutils.R(&InstanceRebuildRootOptions{}, "instance-rebuild-root", "Reinstall virtual server system image", func(cli *aws.SRegion, args *InstanceRebuildRootOptions) error {
		ctx := context.Background()
		img, err := cli.GetImage(args.Image)
		if err != nil {
			return err
		}

		diskID, err := cli.ReplaceSystemDisk(ctx, args.ID, img, args.Size, "", "")
		if err != nil {
			return err
		}
		fmt.Printf("New diskID is %s", diskID)
		return nil
	})

	type InstanceChangeConfigOptions struct {
		ID string `help:"instance ID"`
		aws.SInstanceAttr
	}

	shellutils.R(&InstanceChangeConfigOptions{}, "instance-change-config", "Deploy keypair/password to a stopped virtual server", func(cli *aws.SRegion, args *InstanceChangeConfigOptions) error {
		return cli.ModifyInstanceAttribute(args.ID, &args.SInstanceAttr)
	})

	type InstanceSaveImageOptions struct {
		ID         string `help:"Instance ID"`
		IMAGE_NAME string `help:"Image name"`
		Notes      string `hlep:"Image desc"`
	}
	shellutils.R(&InstanceSaveImageOptions{}, "instance-save-image", "Save instance to image", func(cli *aws.SRegion, args *InstanceSaveImageOptions) error {
		opts := cloudprovider.SaveImageOptions{
			Name:  args.IMAGE_NAME,
			Notes: args.Notes,
		}
		image, err := cli.SaveImage(args.ID, &opts)
		if err != nil {
			return err
		}
		printObject(image)
		return nil
	})

	type SInstanceNicAddrsOptions struct {
		NIC_ID string
		IpAddr []string
	}

	shellutils.R(&SInstanceNicAddrsOptions{}, "instance-nic-assign-ips", "Assing nic ipaddrs", func(cli *aws.SRegion, args *SInstanceNicAddrsOptions) error {
		return cli.AssignAddres(args.NIC_ID, args.IpAddr)
	})

	shellutils.R(&SInstanceNicAddrsOptions{}, "instance-nic-unassign-ips", "Unassing nic ipaddrs", func(cli *aws.SRegion, args *SInstanceNicAddrsOptions) error {
		return cli.UnassignAddress(args.NIC_ID, args.IpAddr)
	})

	type SInstanceNicSubAddrsOptions struct {
		NIC_ID string
	}

	shellutils.R(&SInstanceNicSubAddrsOptions{}, "instance-nic-sub-addrs", "Show nic subaddr", func(cli *aws.SRegion, args *SInstanceNicSubAddrsOptions) error {
		addrs, err := cli.GetSubAddress(args.NIC_ID)
		if err != nil {
			return err
		}
		printObject(addrs)
		return nil
	})

}
