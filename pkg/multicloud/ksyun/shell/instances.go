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
	"yunion.io/x/cloudmux/pkg/cloudprovider"
	"yunion.io/x/cloudmux/pkg/multicloud/ksyun"

	"yunion.io/x/pkg/errors"
	"yunion.io/x/pkg/util/shellutils"
)

func init() {
	type InstanceListOptions struct {
		Id       []string `help:"IDs of instances to show"`
		ZoneName string
	}
	shellutils.R(&InstanceListOptions{}, "instance-list", "list regions", func(cli *ksyun.SRegion, args *InstanceListOptions) error {
		res, err := cli.GetInstances(args.ZoneName, args.Id)
		if err != nil {
			return errors.Wrap(err, "GetInstances")
		}
		printList(res)
		return nil
	})

	type InstanceIdOptions struct {
		ID string
	}

	shellutils.R(&InstanceIdOptions{}, "instance-start", "start instance", func(cli *ksyun.SRegion, args *InstanceIdOptions) error {
		return cli.StartVM(args.ID)
	})

	shellutils.R(&InstanceIdOptions{}, "instance-show", "show instance", func(cli *ksyun.SRegion, args *InstanceIdOptions) error {
		instance, err := cli.GetInstance(args.ID)
		if err != nil {
			return err
		}
		printObject(instance)
		return nil
	})

	shellutils.R(&InstanceIdOptions{}, "instance-delete", "delete instance", func(cli *ksyun.SRegion, args *InstanceIdOptions) error {
		return cli.DeleteVM(args.ID)
	})

	type InstanceStopOptions struct {
		ID           string
		Force        bool
		StopCharging bool
	}
	shellutils.R(&InstanceStopOptions{}, "instance-stop", "stop instance", func(cli *ksyun.SRegion, args *InstanceStopOptions) error {
		return cli.StopVM(args.ID, args.Force, args.StopCharging)
	})

	type InstanceAttachDiskOptions struct {
		ID     string
		DiskID string
	}
	shellutils.R(&InstanceAttachDiskOptions{}, "instance-attach-disk", "attach disk to instance", func(cli *ksyun.SRegion, args *InstanceAttachDiskOptions) error {
		return cli.AttachDisk(args.ID, args.DiskID)
	})

	type InstanceDetachDiskOptions struct {
		ID     string
		DiskID string
	}
	shellutils.R(&InstanceDetachDiskOptions{}, "instance-detach-disk", "detach disk from instance", func(cli *ksyun.SRegion, args *InstanceDetachDiskOptions) error {
		return cli.DetachDisk(args.ID, args.DiskID)
	})

	type InstanceSetSecurityGroupsOptions struct {
		ID         string
		SecGroupID []string
		NicID      string
		SubnetID   string
	}
	shellutils.R(&InstanceSetSecurityGroupsOptions{}, "instance-set-security-groups", "set security groups for instance", func(cli *ksyun.SRegion, args *InstanceSetSecurityGroupsOptions) error {
		return cli.SetSecurityGroups(args.SecGroupID, args.ID, args.NicID, args.SubnetID)
	})

	type InstanceGetVNCInfoOptions struct {
		ID string
	}
	shellutils.R(&InstanceGetVNCInfoOptions{}, "instance-vnc", "get vnc info for instance", func(cli *ksyun.SRegion, args *InstanceGetVNCInfoOptions) error {
		vnc, err := cli.GetVNCInfo(args.ID)
		if err != nil {
			return errors.Wrap(err, "GetVNCInfo")
		}
		printObject(vnc)
		return nil
	})

	type InstanceUpdateOptions struct {
		ID          string
		Name        string
		HostName    string
		Description string
	}
	shellutils.R(&InstanceUpdateOptions{}, "instance-update", "update instance", func(cli *ksyun.SRegion, args *InstanceUpdateOptions) error {
		return cli.UpdateVM(args.ID, cloudprovider.SInstanceUpdateOptions{
			NAME:        args.Name,
			HostName:    args.HostName,
			Description: args.Description,
		})
	})

	type InstanceRebuildRootOptions struct {
		ID        string
		ImageId   string
		Password  string
		PublicKey string
		UserData  string
	}
	shellutils.R(&InstanceRebuildRootOptions{}, "instance-rebuild-root", "rebuild instance root", func(cli *ksyun.SRegion, args *InstanceRebuildRootOptions) error {
		opts := &cloudprovider.SManagedVMRebuildRootConfig{ImageId: args.ImageId, Password: args.Password, PublicKey: args.PublicKey, UserData: args.UserData}
		return cli.RebuildRoot(args.ID, opts)
	})

	type InstanceChangeConfigOptions struct {
		ID           string
		InstanceType string
	}
	shellutils.R(&InstanceChangeConfigOptions{}, "instance-change-config", "change instance config", func(cli *ksyun.SRegion, args *InstanceChangeConfigOptions) error {
		opts := &cloudprovider.SManagedVMChangeConfig{InstanceType: args.InstanceType}
		return cli.ChangeConfig(args.ID, opts)
	})

	type InstanceDeployOptions struct {
		ID            string
		PublicKey     string
		DeleteKeypair bool
		Password      string
	}
	shellutils.R(&InstanceDeployOptions{}, "instance-deploy", "deploy instance", func(cli *ksyun.SRegion, args *InstanceDeployOptions) error {
		return cli.DeployVM(args.ID, &cloudprovider.SInstanceDeployOptions{PublicKey: args.PublicKey, DeleteKeypair: args.DeleteKeypair, Password: args.Password})
	})
}
