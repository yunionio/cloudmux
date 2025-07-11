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

	"yunion.io/x/pkg/util/shellutils"

	"yunion.io/x/cloudmux/pkg/multicloud/qcloud"
)

func init() {
	type CallerShowOptions struct {
	}
	shellutils.R(&CallerShowOptions{}, "caller-show", "Show caller", func(cli *qcloud.SRegion, args *CallerShowOptions) error {
		caller, err := cli.GetClient().GetCallerIdentity()
		if err != nil {
			return err
		}
		printObject(caller)
		return nil
	})

	shellutils.R(&CallerShowOptions{}, "app-id-show", "Show app id", func(cli *qcloud.SRegion, args *CallerShowOptions) error {
		appId, err := cli.GetClient().GetAppId()
		if err != nil {
			return err
		}
		fmt.Println(appId)
		return nil
	})

}
