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
	"github.com/pkg/errors"
	"yunion.io/x/cloudmux/pkg/cloudprovider"
	"yunion.io/x/cloudmux/pkg/multicloud/ksyun"
	"yunion.io/x/jsonutils"
	"yunion.io/x/log"
	"yunion.io/x/pkg/util/shellutils"
)

func init() {
	type InstanceListOptions struct {
	}
	shellutils.R(&InstanceListOptions{}, "instance-list", "list regions", func(cli *ksyun.SRegion, args *InstanceListOptions) error {
		res, err := cli.GetInstances()
		if err != nil {
			return errors.Wrap(err, "GetInstances")
		}
		for _, ins := range res {
			disks, err := cli.GetDisks(nil)
			if err != nil {
				log.Errorln("get disk err:", err)
				return err
			}
			res := []cloudprovider.ICloudDisk{}
			for _, disk := range disks {
				if disk.InstanceID == ins.InstanceID {
					res = append(res, &disk)
				}
			}
			log.Infoln("this is disks:", jsonutils.Marshal(res))
		}
		printList(res, 0, 0, 0, []string{})
		return nil
	})

}
