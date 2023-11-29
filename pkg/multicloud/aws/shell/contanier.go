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

	"yunion.io/x/cloudmux/pkg/multicloud/aws"
)

func init() {
	type ContainerClusterListOptions struct {
	}
	shellutils.R(&ContainerClusterListOptions{}, "container-cluster-list", "List container clusters", func(cli *aws.SRegion, args *ContainerClusterListOptions) error {
		clusters, err := cli.ListClusters()
		if err != nil {
			return err
		}
		fmt.Println(clusters)
		return nil
	})

	type ContainerServiceListOptions struct {
		CLUSTER string
	}
	shellutils.R(&ContainerServiceListOptions{}, "container-service-list", "List cluster services", func(cli *aws.SRegion, args *ContainerServiceListOptions) error {
		services, err := cli.ListServices(args.CLUSTER)
		if err != nil {
			return err
		}
		fmt.Println(services)
		return nil
	})

	type ContainerTaskListOptions struct {
		CLUSTER string
	}
	shellutils.R(&ContainerTaskListOptions{}, "container-task-list", "List cluster tasks", func(cli *aws.SRegion, args *ContainerTaskListOptions) error {
		tasks, err := cli.ListTasks(args.CLUSTER)
		if err != nil {
			return err
		}
		fmt.Println(tasks)
		return nil
	})

}
