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

	"yunion.io/x/cloudmux/pkg/cloudprovider"
	"yunion.io/x/cloudmux/pkg/multicloud/aws"
)

func init() {
	type LbListOptions struct {
		Id     string
		Marker string
	}
	shellutils.R(&LbListOptions{}, "lb-list", "List loadbalancer", func(cli *aws.SRegion, args *LbListOptions) error {
		ret, _, e := cli.GetLoadbalancers(args.Id, args.Marker)
		if e != nil {
			return e
		}
		printList(ret, 0, 0, 0, []string{})
		return nil
	})

	shellutils.R(&LbListOptions{}, "lb-lis-list", "List loadbalancer", func(cli *aws.SRegion, args *LbListOptions) error {
		ret, _, e := cli.GetElbListeners(args.Id, "", args.Marker)
		if e != nil {
			return e
		}
		printList(ret, 0, 0, 0, []string{})
		return nil
	})

	type LbIdOptions struct {
		ID string
	}
	shellutils.R(&LbIdOptions{}, "lb-show", "Show loadbalancer attribute", func(cli *aws.SRegion, args *LbIdOptions) error {
		ret, err := cli.GetElbAttributes(args.ID)
		if err != nil {
			return err
		}
		printObject(ret)
		return nil
	})

	shellutils.R(&cloudprovider.SLoadbalancerCreateOptions{}, "lb-create", "Create loadbalancer", func(cli *aws.SRegion, args *cloudprovider.SLoadbalancerCreateOptions) error {
		ret, err := cli.CreateLoadbalancer(args)
		if err != nil {
			return err
		}
		printObject(ret)
		return nil
	})

	shellutils.R(&LbIdOptions{}, "lb-delete", "Delete loadbalancer", func(cli *aws.SRegion, args *LbIdOptions) error {
		return cli.DeleteElb(args.ID)
	})

	shellutils.R(&LbIdOptions{}, "lb-tag-list", "Show loadbalancer tags", func(cli *aws.SRegion, args *LbIdOptions) error {
		ret, err := cli.DescribeElbTags(args.ID)
		if err != nil {
			return err
		}
		printObject(ret)
		return nil
	})

	type LbBackendGroupListOptions struct {
		ElbId string
		Id    string
	}

	shellutils.R(&LbBackendGroupListOptions{}, "lb-backend-group-list", "List loadbalancer backend groups", func(cli *aws.SRegion, args *LbBackendGroupListOptions) error {
		ret, err := cli.GetElbBackendgroups(args.ElbId, args.Id)
		if err != nil {
			return err
		}
		printList(ret, 0, 0, 0, []string{})
		return nil
	})

	shellutils.R(&cloudprovider.SLoadbalancerCertificate{}, "lb-cert-create", "Create loadbalancer cert", func(cli *aws.SRegion, args *cloudprovider.SLoadbalancerCertificate) error {
		arn, err := cli.CreateLoadbalancerCertifacate(args)
		if err != nil {
			return err
		}
		fmt.Println(arn)
		return nil
	})

	type LbCertListOption struct {
	}

	shellutils.R(&LbCertListOption{}, "lb-cert-list", "Create loadbalancer cert", func(cli *aws.SRegion, args *LbCertListOption) error {
		certs, err := cli.ListServerCertificates()
		if err != nil {
			return err
		}
		printList(certs, 0, 0, 0, nil)
		return nil
	})

	type LbListenerRuleListOptions struct {
		ListenerId string
		RuleId     string
	}

	shellutils.R(&LbListenerRuleListOptions{}, "lb-listener-rule-list", "List loadbalancer listener rules", func(cli *aws.SRegion, args *LbListenerRuleListOptions) error {
		rules, err := cli.GetElbListenerRules(args.ListenerId, args.RuleId)
		if err != nil {
			return err
		}
		printList(rules, 0, 0, 0, nil)
		return nil
	})

}
