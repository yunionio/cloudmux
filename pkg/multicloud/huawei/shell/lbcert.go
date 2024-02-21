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

	"yunion.io/x/cloudmux/pkg/multicloud/huawei"
)

func init() {
	type ElbCertListOptions struct {
	}
	shellutils.R(&ElbCertListOptions{}, "lbcert-list", "List lb certs", func(cli *huawei.SRegion, args *ElbCertListOptions) error {
		certs, err := cli.GetLoadBalancerCertificates()
		if err != nil {
			return err
		}
		printList(certs, len(certs), 0, 0, []string{})
		return nil
	})

	type ElbIdOptions struct {
		ID string
	}

	shellutils.R(&ElbIdOptions{}, "lbcert-show", "Show lb cert", func(cli *huawei.SRegion, args *ElbIdOptions) error {
		cert, err := cli.GetLoadBalancerCertificate(args.ID)
		if err != nil {
			return err
		}
		printObject(cert)
		return nil
	})

}
