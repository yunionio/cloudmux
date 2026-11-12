// Copyright 2026 Yunion
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package volcengine

import (
	"time"
	api "yunion.io/x/cloudmux/pkg/apis/compute"
	"yunion.io/x/cloudmux/pkg/cloudprovider"
	"yunion.io/x/cloudmux/pkg/multicloud"
	"yunion.io/x/pkg/errors"
)

type SKubeCluster struct {
	multicloud.SResourceBase
	VolcEngineTags
	region                                               *SRegion
	Id, Name, ProjectName, KubernetesVersion, CreateTime string
	Status                                               struct{ Phase string }
	ClusterConfig                                        struct {
		VpcId     string
		SubnetIds []string
	}
}

var _ cloudprovider.ICloudKubeCluster = (*SKubeCluster)(nil)

func (r *SRegion) GetICloudKubeClusters() ([]cloudprovider.ICloudKubeCluster, error) {
	items, err := r.inventoryPages("vke", "ListClusters", "Items", "TotalCount", "Id", nil)
	if err != nil {
		return nil, err
	}
	ret := make([]cloudprovider.ICloudKubeCluster, 0, len(items))
	for _, item := range items {
		c := &SKubeCluster{region: r}
		if err := item.Unmarshal(c); err != nil {
			return nil, err
		}
		ret = append(ret, c)
	}
	return ret, nil
}
func (r *SRegion) GetICloudKubeClusterById(id string) (cloudprovider.ICloudKubeCluster, error) {
	clusters, e := r.GetICloudKubeClusters()
	if e != nil {
		return nil, e
	}
	for _, c := range clusters {
		if c.GetGlobalId() == id {
			return c, nil
		}
	}
	return nil, cloudprovider.ErrNotFound
}
func (c *SKubeCluster) GetId() string        { return c.Id }
func (c *SKubeCluster) GetGlobalId() string  { return c.Id }
func (c *SKubeCluster) GetName() string      { return c.Name }
func (c *SKubeCluster) GetProjectId() string { return c.ProjectName }
func (c *SKubeCluster) GetCreatedAt() time.Time {
	t, _ := time.Parse(time.RFC3339Nano, c.CreateTime)
	return t
}
func (c *SKubeCluster) GetStatus() string {
	if c.Status.Phase == "Running" {
		return api.KUBE_CLUSTER_STATUS_RUNNING
	}
	return c.Status.Phase
}
func (c *SKubeCluster) GetEnabled() bool        { return true }
func (c *SKubeCluster) GetVersion() string      { return c.KubernetesVersion }
func (c *SKubeCluster) GetVpcId() string        { return c.ClusterConfig.VpcId }
func (c *SKubeCluster) GetNetworkIds() []string { return c.ClusterConfig.SubnetIds }

// Fetching credentials is deliberately excluded from read-only resource discovery.
func (c *SKubeCluster) GetKubeConfig(private bool, expireMinutes int) (*cloudprovider.SKubeconfig, error) {
	return nil, cloudprovider.ErrNotSupported
}
func (c *SKubeCluster) Delete(retain bool) error { return cloudprovider.ErrNotSupported }
func (c *SKubeCluster) CreateIKubeNodePool(opts *cloudprovider.KubeNodePoolCreateOptions) (cloudprovider.ICloudKubeNodePool, error) {
	return nil, cloudprovider.ErrNotSupported
}
func (c *SKubeCluster) Refresh() error {
	r, e := c.region.GetICloudKubeClusterById(c.Id)
	if e != nil {
		return e
	}
	*c = *(r.(*SKubeCluster))
	return nil
}
func (c *SKubeCluster) GetIKubeNodes() ([]cloudprovider.ICloudKubeNode, error) {
	items, e := c.region.inventoryPages("vke", "ListNodes", "Items", "TotalCount", "Id", map[string]interface{}{"ClusterId": c.Id})
	if e != nil {
		return nil, e
	}
	ret := make([]cloudprovider.ICloudKubeNode, 0, len(items))
	for _, item := range items {
		n := &SKubeNode{}
		if e := item.Unmarshal(n); e != nil {
			return nil, e
		}
		if n.ClusterId != c.Id {
			return nil, errors.Errorf("VKE node cluster mismatch")
		}
		ret = append(ret, n)
	}
	return ret, nil
}
func (c *SKubeCluster) GetIKubeNodePools() ([]cloudprovider.ICloudKubeNodePool, error) {
	items, e := c.region.inventoryPages("vke", "ListNodePools", "Items", "TotalCount", "Id", map[string]interface{}{"ClusterId": c.Id})
	if e != nil {
		return nil, e
	}
	ret := make([]cloudprovider.ICloudKubeNodePool, 0, len(items))
	for _, item := range items {
		n := &SKubeNodePool{}
		if e := item.Unmarshal(n); e != nil {
			return nil, e
		}
		if n.ClusterId != c.Id {
			return nil, errors.Errorf("VKE node pool cluster mismatch")
		}
		ret = append(ret, n)
	}
	return ret, nil
}

type SKubeNode struct {
	multicloud.SResourceBase
	VolcEngineTags
	Id, Name, ClusterId, NodePoolId, InstanceId, MetadataName string
	Status                                                    struct{ Phase string }
}

var _ cloudprovider.ICloudKubeNode = (*SKubeNode)(nil)

func (n *SKubeNode) GetId() string       { return n.Id }
func (n *SKubeNode) GetGlobalId() string { return n.Id }
func (n *SKubeNode) GetName() string {
	if n.MetadataName != "" {
		return n.MetadataName
	}
	return n.Name
}
func (n *SKubeNode) GetStatus() string      { return n.Status.Phase }
func (n *SKubeNode) GetINodePoolId() string { return n.NodePoolId }

type SKubeNodePool struct {
	multicloud.SResourceBase
	VolcEngineTags
	Id, Name, ClusterId string
	Status              struct{ Phase string }
	AutoScaling         struct{ MinReplicas, MaxReplicas, DesiredReplicas int }
	NodeConfig          struct {
		InstanceTypeIds, SubnetIds []string
		SystemVolume               struct{ Size int }
	}
}

var _ cloudprovider.ICloudKubeNodePool = (*SKubeNodePool)(nil)

func (n *SKubeNodePool) GetId() string                { return n.Id }
func (n *SKubeNodePool) GetGlobalId() string          { return n.Id }
func (n *SKubeNodePool) GetName() string              { return n.Name }
func (n *SKubeNodePool) GetStatus() string            { return n.Status.Phase }
func (n *SKubeNodePool) GetMinInstanceCount() int     { return n.AutoScaling.MinReplicas }
func (n *SKubeNodePool) GetMaxInstanceCount() int     { return n.AutoScaling.MaxReplicas }
func (n *SKubeNodePool) GetDesiredInstanceCount() int { return n.AutoScaling.DesiredReplicas }
func (n *SKubeNodePool) GetRootDiskSizeGb() int       { return n.NodeConfig.SystemVolume.Size }
func (n *SKubeNodePool) GetInstanceTypes() []string   { return n.NodeConfig.InstanceTypeIds }
func (n *SKubeNodePool) GetNetworkIds() []string      { return n.NodeConfig.SubnetIds }
func (n *SKubeNodePool) Delete() error                { return cloudprovider.ErrNotSupported }
