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

	"yunion.io/x/cloudmux/pkg/multicloud/aws"
)

func init() {
	type ElasticacheClusterListOption struct {
		Id         string
		IsMemcache bool
	}
	shellutils.R(&ElasticacheClusterListOption{}, "elasti-cache-list", "List elasticacheCluster", func(cli *aws.SRegion, args *ElasticacheClusterListOption) error {
		clusters, e := cli.GetElasticaches(args.Id, args.IsMemcache)
		if e != nil {
			return e
		}
		printList(clusters, len(clusters), 0, len(clusters), []string{})
		return nil
	})

	type ElasticacheEngineVersionListOption struct {
		Engine string `default:"redis" choices:"redis|memcached"`
	}
	shellutils.R(&ElasticacheEngineVersionListOption{}, "elasti-cache-engine-version-list", "List elasticache engine version", func(cli *aws.SRegion, args *ElasticacheEngineVersionListOption) error {
		versions, err := cli.GetElastiCacheEngineVersion(args.Engine)
		if err != nil {
			return err
		}
		printList(versions, 0, 0, 0, []string{})
		return nil
	})

	type ElasticacheClusterIdOption struct {
		ID string
	}

	shellutils.R(&ElasticacheClusterIdOption{}, "elasti-cache-delete", "Delete elasticache", func(cli *aws.SRegion, args *ElasticacheClusterIdOption) error {
		return cli.DeleteElastiCache(args.ID)
	})

	shellutils.R(&ElasticacheClusterIdOption{}, "replication-group-delete", "Delete replication group", func(cli *aws.SRegion, args *ElasticacheClusterIdOption) error {
		return cli.DeleteReplicationGroup(args.ID)
	})

	type ElasticacheReplicaGroupListOption struct {
		Id string
	}
	shellutils.R(&ElasticacheReplicaGroupListOption{}, "replication-group-list", "List elasticaReplicaGroup", func(cli *aws.SRegion, args *ElasticacheReplicaGroupListOption) error {
		clusters, e := cli.GetReplicationGroups(args.Id)
		if e != nil {
			return e
		}
		printList(clusters, len(clusters), 0, len(clusters), []string{})
		return nil
	})

	type ElasticacheSubnetGroupOption struct {
		Id string `help:"subnetgroupId"`
	}
	shellutils.R(&ElasticacheSubnetGroupOption{}, "elasti-cache-subnet-show", "List elasticacheSubnetGroup", func(cli *aws.SRegion, args *ElasticacheSubnetGroupOption) error {
		subnetGroups, e := cli.DescribeCacheSubnetGroups(args.Id)
		if e != nil {
			return e
		}
		printList(subnetGroups, len(subnetGroups), 0, len(subnetGroups), []string{})
		return nil
	})

	type ElasticacheSnapshotOption struct {
		ReplicaGroupId string `help:"replicaGroupId"`
		SnapshotId     string `help:"SnapshotId"`
	}
	shellutils.R(&ElasticacheSnapshotOption{}, "elasti-cache-snapshot-list", "List elasticacheSnapshot", func(cli *aws.SRegion, args *ElasticacheSnapshotOption) error {
		snapshots, e := cli.GetCacheSnapshots(args.ReplicaGroupId, args.SnapshotId)
		if e != nil {
			return e
		}
		printList(snapshots, len(snapshots), 0, len(snapshots), []string{})
		return nil
	})

	type ElasticacheParameterOption struct {
		ParameterGroupId string
	}
	shellutils.R(&ElasticacheParameterOption{}, "elasti-cache-parameter-list", "List elasticacheParameter", func(cli *aws.SRegion, args *ElasticacheParameterOption) error {
		parameters, e := cli.GetCacheParameters(args.ParameterGroupId)
		if e != nil {
			return e
		}
		printList(parameters, len(parameters), 0, len(parameters), []string{})
		return nil
	})

	type ElasticacheUserOption struct {
		Engine string `help:"redis"`
		UserId string
	}
	shellutils.R(&ElasticacheUserOption{}, "elasti-cache-user-list", "List elasticacheUser", func(cli *aws.SRegion, args *ElasticacheUserOption) error {
		users, e := cli.GetCacheUsers(args.Engine, args.UserId)
		if e != nil {
			return e
		}
		printList(users, len(users), 0, len(users), []string{})
		return nil
	})
}
