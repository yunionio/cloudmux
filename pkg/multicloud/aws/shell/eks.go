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

	"yunion.io/x/cloudmux/pkg/cloudprovider"
	"yunion.io/x/cloudmux/pkg/multicloud/aws"
)

func init() {
	type KubeClusterListOptions struct {
		NextToken string
	}
	shellutils.R(&KubeClusterListOptions{}, "kube-cluster-list", "List kube cluster", func(cli *aws.SRegion, args *KubeClusterListOptions) error {
		ret, _, err := cli.GetKubeClusters(args.NextToken)
		if err != nil {
			return err
		}
		printList(ret, 0, 0, 0, []string{})
		return nil
	})

	type KubeClusterNameOptions struct {
		NAME string
	}

	shellutils.R(&KubeClusterNameOptions{}, "kube-cluster-show", "Show kube cluster", func(cli *aws.SRegion, args *KubeClusterNameOptions) error {
		ret, err := cli.GetKubeCluster(args.NAME)
		if err != nil {
			return err
		}
		printObject(ret)
		return nil
	})

	shellutils.R(&KubeClusterNameOptions{}, "kube-cluster-delete", "Delete kube cluster", func(cli *aws.SRegion, args *KubeClusterNameOptions) error {
		return cli.DeleteKubeCluster(args.NAME)
	})

	shellutils.R(&cloudprovider.KubeClusterCreateOptions{}, "kube-cluster-create", "Create kube cluster", func(cli *aws.SRegion, args *cloudprovider.KubeClusterCreateOptions) error {
		cluster, err := cli.CreateKubeCluster(args)
		if err != nil {
			return err
		}
		printObject(cluster)
		return nil
	})

}
