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

	"yunion.io/x/cloudmux/pkg/multicloud/aliyun"
)

func init() {
	type ServiceListOptions struct {
	}
	shellutils.R(&ServiceListOptions{}, "fc-service-list", "List Service", func(cli *aliyun.SRegion, args *ServiceListOptions) error {
		ret, err := cli.GetFcServices()
		if err != nil {
			return err
		}
		printList(ret, 0, 0, 0, []string{})
		return nil
	})

	type FunctionListOptions struct {
		SERVICE string
	}

	shellutils.R(&FunctionListOptions{}, "fc-function-list", "List Function", func(cli *aliyun.SRegion, args *FunctionListOptions) error {
		ret, err := cli.GetFcFunctions(args.SERVICE)
		if err != nil {
			return err
		}
		printList(ret, 0, 0, 0, []string{})
		return nil
	})

	type InstanceListOptions struct {
		SERVICE  string
		FUNCTION string
	}

	shellutils.R(&InstanceListOptions{}, "fc-instance-list", "List Instance", func(cli *aliyun.SRegion, args *InstanceListOptions) error {
		ret, err := cli.GetFcInstances(args.SERVICE, args.FUNCTION)
		if err != nil {
			return err
		}
		printList(ret, 0, 0, 0, []string{})
		return nil
	})

}
