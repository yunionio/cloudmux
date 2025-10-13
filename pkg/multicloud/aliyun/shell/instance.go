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

	"yunion.io/x/pkg/util/billing"
	"yunion.io/x/pkg/util/shellutils"

	"yunion.io/x/cloudmux/pkg/cloudprovider"
	"yunion.io/x/cloudmux/pkg/multicloud/aliyun"
)

func init() {
	type InstanceListOptions struct {
		Id   []string `help:"IDs of instances to show"`
		Zone string   `help:"Zone ID"`
	}
	shellutils.R(&InstanceListOptions{}, "instance-list", "List intances", func(cli *aliyun.SRegion, args *InstanceListOptions) error {
		instances, e := cli.GetInstances(args.Zone, args.Id)
		if e != nil {
			return e
		}
		printList(instances, 0, 0, 0, []string{})
		return nil
	})

	type InstanceCreateOptions struct {
		NAME      string `help:"name of instance"`
		IMAGE     string `help:"image ID"`
		CPU       int    `help:"CPU count"`
		MEMORYGB  int    `help:"MemoryGB"`
		Disk      []int  `help:"Data disk sizes int GB"`
		STORAGE   string `help:"Storage type"`
		VSWITCH   string `help:"Vswitch ID"`
		PASSWD    string `help:"password"`
		PublicKey string `help:"PublicKey"`
	}
	shellutils.R(&InstanceCreateOptions{}, "instance-create", "Create a instance", func(cli *aliyun.SRegion, args *InstanceCreateOptions) error {
		instance, e := cli.CreateInstanceSimple(args.NAME, args.IMAGE, args.CPU, args.MEMORYGB, args.STORAGE, args.Disk, args.VSWITCH, args.PASSWD, args.PublicKey)
		if e != nil {
			return e
		}
		printObject(instance)
		return nil
	})

	type InstanceDiskOperationOptions struct {
		ID   string `help:"instance ID"`
		DISK string `help:"disk ID"`
	}

	shellutils.R(&InstanceDiskOperationOptions{}, "instance-attach-disk", "Attach a disk to instance", func(cli *aliyun.SRegion, args *InstanceDiskOperationOptions) error {
		err := cli.AttachDisk(args.ID, args.DISK)
		if err != nil {
			return err
		}
		return nil
	})

	shellutils.R(&InstanceDiskOperationOptions{}, "instance-detach-disk", "Detach a disk to instance", func(cli *aliyun.SRegion, args *InstanceDiskOperationOptions) error {
		err := cli.DetachDisk(args.ID, args.DISK)
		if err != nil {
			return err
		}
		return nil
	})

	type InstanceOperationOptions struct {
		ID string `help:"instance ID"`
	}
	shellutils.R(&InstanceOperationOptions{}, "instance-start", "Start a instance", func(cli *aliyun.SRegion, args *InstanceOperationOptions) error {
		err := cli.StartVM(args.ID)
		if err != nil {
			return err
		}
		return nil
	})

	shellutils.R(&InstanceOperationOptions{}, "instance-status", "Show a instance status", func(cli *aliyun.SRegion, args *InstanceOperationOptions) error {
		status, err := cli.DescribeInstancesFullStatus(args.ID)
		if err != nil {
			return err
		}
		printObject(status)
		return nil
	})

	shellutils.R(&InstanceOperationOptions{}, "instance-auto-renew-info", "Show instance auto renew info", func(cli *aliyun.SRegion, args *InstanceOperationOptions) error {
		info, err := cli.GetInstanceAutoRenewAttribute(args.ID)
		if err != nil {
			return err
		}
		printObject(info)
		return nil
	})

	shellutils.R(&InstanceOperationOptions{}, "instance-eip-convert", "Convert instance public ip to eip", func(cli *aliyun.SRegion, args *InstanceOperationOptions) error {
		err := cli.ConvertPublicIpToEip(args.ID)
		if err != nil {
			return err
		}
		return nil
	})

	shellutils.R(&InstanceOperationOptions{}, "instance-vnc", "Get a instance VNC url", func(cli *aliyun.SRegion, args *InstanceOperationOptions) error {
		url, err := cli.GetInstanceVNCUrl(args.ID)
		if err != nil {
			return err
		}
		fmt.Println(url)
		return nil
	})

	type InstanceStopOptions struct {
		ID           string `help:"instance ID"`
		Force        bool   `help:"Force stop instance"`
		StopCharging bool   `help:"Stop Charging"`
	}
	shellutils.R(&InstanceStopOptions{}, "instance-stop", "Stop a instance", func(cli *aliyun.SRegion, args *InstanceStopOptions) error {
		err := cli.StopVM(args.ID, args.Force, args.StopCharging)
		if err != nil {
			return err
		}
		return nil
	})
	shellutils.R(&InstanceOperationOptions{}, "instance-delete", "Delete a instance", func(cli *aliyun.SRegion, args *InstanceOperationOptions) error {
		err := cli.DeleteVM(args.ID)
		if err != nil {
			return err
		}
		return nil
	})

	/*
		server-change-config 更改系统配置
		server-reset
	*/
	type InstanceDeployOptions struct {
		ID            string `help:"instance ID"`
		PublicKey     string `help:"Keypair Name"`
		DeleteKeypair bool   `help:"Remove SSH keypair"`
		Password      string `help:"new password"`
	}

	shellutils.R(&InstanceDeployOptions{}, "instance-deploy", "Deploy keypair/password to a stopped virtual server", func(cli *aliyun.SRegion, args *InstanceDeployOptions) error {
		err := cli.DeployVM(args.ID, &cloudprovider.SInstanceDeployOptions{DeleteKeypair: args.DeleteKeypair, Password: args.Password, PublicKey: args.PublicKey})
		if err != nil {
			return err
		}
		return nil
	})

	type InstanceRebuildRootOptions struct {
		ID       string `help:"instance ID"`
		Image    string `help:"Image ID"`
		Password string `help:"pasword"`
		Keypair  string `help:"keypair name"`
		Size     int    `help:"system disk size in GB"`
	}

	shellutils.R(&InstanceRebuildRootOptions{}, "instance-rebuild-root", "Reinstall virtual server system image", func(cli *aliyun.SRegion, args *InstanceRebuildRootOptions) error {
		diskID, err := cli.ReplaceSystemDisk(args.ID, args.Image, args.Password, args.Keypair, args.Size)
		if err != nil {
			return err
		}
		fmt.Printf("New diskID is %s", diskID)
		return nil
	})

	type InstanceChangeConfigOptions struct {
		ID           string `help:"instance ID"`
		InstanceType string `help:"instance type"`
	}

	shellutils.R(&InstanceChangeConfigOptions{}, "instance-change-config", "Deploy keypair/password to a stopped virtual server", func(cli *aliyun.SRegion, args *InstanceChangeConfigOptions) error {
		err := cli.ChangeVMConfig(args.ID, args.InstanceType)
		if err != nil {
			return err
		}
		return nil
	})

	type InstanceUpdatePasswordOptions struct {
		ID     string `help:"Instance ID"`
		PASSWD string `help:"new password"`
	}
	shellutils.R(&InstanceUpdatePasswordOptions{}, "instance-update-password", "Update instance password", func(cli *aliyun.SRegion, args *InstanceUpdatePasswordOptions) error {
		err := cli.UpdateInstancePassword(args.ID, args.PASSWD)
		return err
	})

	type InstanceSetAutoRenewOptions struct {
		ID        string `help:"Instance ID"`
		AutoRenew bool   `help:"Is auto renew instance"`
		Duration  string
	}

	shellutils.R(&InstanceSetAutoRenewOptions{}, "instance-set-auto-renew", "Set instance auto renew", func(cli *aliyun.SRegion, args *InstanceSetAutoRenewOptions) error {
		bc, _ := billing.ParseBillingCycle(args.Duration)
		bc.AutoRenew = args.AutoRenew
		return cli.SetInstanceAutoRenew(args.ID, bc)
	})

	type InstanceSaveImageOptions struct {
		ID         string `help:"Instance ID"`
		IMAGE_NAME string `help:"Image name"`
		Notes      string `hlep:"Image desc"`
	}
	shellutils.R(&InstanceSaveImageOptions{}, "instance-save-image", "Save instance to image", func(cli *aliyun.SRegion, args *InstanceSaveImageOptions) error {
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

	type InstanceChangeBillingType struct {
		ID          string
		BillingType string `choices:"postpaid|prepaid" default:"prepaid"`
	}

	shellutils.R(&InstanceChangeBillingType{}, "instance-change-billing-type", "change instance billing type", func(cli *aliyun.SRegion, args *InstanceChangeBillingType) error {
		return cli.ModifyInstanceChargeType(args.ID, args.BillingType)
	})

}
