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
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"testing"

	"yunion.io/x/cloudmux/pkg/cloudprovider"
)

func inventoryTestRegion(t *testing.T, handler func(string, string, map[string]interface{}) (int, interface{})) *SRegion {
	t.Helper()
	c := &SVolcEngineClient{VolcEngineClientConfig: NewVolcEngineClientConfig("test", "test")}
	c.postgresqlClient = &sPostgresqlClient{httpClient: &http.Client{Transport: postgresqlRoundTripper(func(req *http.Request) (*http.Response, error) {
		host := strings.Split(req.URL.Host, ".")
		service := host[0]
		if req.URL.Scheme != "https" || req.URL.Host != service+".cn-beijing.volcengineapi.com" {
			t.Fatalf("unexpected host: %s", req.URL.Host)
		}
		if req.Header.Get("Authorization") == "" {
			t.Fatal("unsigned request")
		}
		p := map[string]interface{}{}
		if service == "clb" {
			if req.Method != "GET" || req.URL.Query().Get("Version") != "2020-04-01" {
				t.Fatal("CLB protocol")
			}
			for k, v := range req.URL.Query() {
				p[k] = v[0]
			}
		} else {
			if req.Method != "POST" || req.Header.Get("Content-Type") != "application/json" {
				t.Fatal("JSON protocol")
			}
			if e := json.NewDecoder(req.Body).Decode(&p); e != nil {
				t.Fatal(e)
			}
		}
		status, result := handler(service, req.URL.Query().Get("Action"), p)
		return postgresqlResponse(status, result), nil
	})}}
	return &SRegion{client: c, RegionId: "cn-beijing"}
}
func result(v interface{}) map[string]interface{} { return map[string]interface{}{"Result": v} }
func TestInventoryPaginationFailsClosed(t *testing.T) {
	for _, body := range []interface{}{
		map[string]interface{}{}, map[string]interface{}{"TotalCount": -1, "Items": []interface{}{}},
		map[string]interface{}{"TotalCount": 0, "Items": nil},
		map[string]interface{}{"TotalCount": 1, "Items": []interface{}{}},
		map[string]interface{}{"TotalCount": 2, "Items": []map[string]string{{"Id": "same"}, {"Id": "same"}}},
	} {
		t.Run(fmt.Sprint(body), func(t *testing.T) {
			r := inventoryTestRegion(t, func(s, a string, p map[string]interface{}) (int, interface{}) { return 200, result(body) })
			if items, e := r.inventoryPages("vke", "ListClusters", "Items", "TotalCount", "Id", nil); e == nil {
				t.Fatalf("accepted invalid inventory: %v", items)
			}
		})
	}
}
func TestInventoryPaginationAndPermissions(t *testing.T) {
	for _, mode := range []string{"ok", "denied", "changed"} {
		t.Run(mode, func(t *testing.T) {
			calls := 0
			r := inventoryTestRegion(t, func(s, a string, p map[string]interface{}) (int, interface{}) {
				calls++
				if calls == 2 && mode == "denied" {
					return 403, map[string]interface{}{"ResponseMetadata": map[string]interface{}{"Error": map[string]interface{}{"Code": "AccessDenied", "Message": "denied"}}}
				}
				count := 2
				if calls == 2 && mode == "changed" {
					count = 3
				}
				return 200, result(map[string]interface{}{"TotalCount": count, "Items": []map[string]string{{"Id": fmt.Sprint(calls)}}})
			})
			items, e := r.inventoryPages("vke", "ListClusters", "Items", "TotalCount", "Id", nil)
			if mode == "ok" {
				if e != nil || len(items) != 2 {
					t.Fatalf("%v %v", items, e)
				}
			} else if e == nil || items != nil {
				t.Fatal("partial inventory accepted")
			}
		})
	}
}
func TestInventoryWritesAndEndpointsRejected(t *testing.T) {
	r := inventoryTestRegion(t, func(s, a string, p map[string]interface{}) (int, interface{}) {
		t.Fatal("unexpected request")
		return 200, nil
	})
	for _, v := range [][3]string{{"cn-beijing", "redis", "DeleteDBInstance"}, {"cn-beijing", "clb", "DeleteLoadBalancer"}, {"cn-beijing", "vke", "CreateKubeconfig"}, {"cn-beijing.example.invalid", "redis", "DescribeRegions"}, {"cn-beijing", "unknown", "ListClusters"}} {
		if _, e := r.client.inventoryRequest(v[0], v[1], v[2], nil); e == nil {
			t.Fatalf("accepted %v", v)
		}
	}
}
func redisFixture(id string) map[string]interface{} {
	return map[string]interface{}{"InstanceId": id, "RegionId": "cn-beijing", "InstanceName": "test", "Status": "Running", "Capacity": map[string]int{"Total": 2048}, "InstanceClass": "redis.shard.2g", "ProjectName": "test-project", "VpcId": "vpc-test", "SubnetId": "subnet-test", "ZoneIds": []string{"cn-beijing-a"}, "EngineVersion": "7.0", "NodeNumber": 2, "VpcAuthMode": "close", "Tags": []map[string]string{{"Key": "owner", "Value": "test"}}, "VisitAddrs": []map[string]string{{"AddrType": "Public", "Address": "public.example.invalid", "Port": "6380"}, {"AddrType": "Private", "Address": "private.example.invalid", "VIP": "10.0.0.1", "Port": "6379"}}}
}
func TestRedisRepeatedDiscoveryAndReadOnly(t *testing.T) {
	r := inventoryTestRegion(t, func(s, a string, p map[string]interface{}) (int, interface{}) {
		switch a {
		case "DescribeRegions":
			return 200, result(map[string]interface{}{"Regions": []map[string]string{{"RegionId": "cn-beijing"}}})
		case "DescribeDBInstances":
			return 200, result(map[string]interface{}{"TotalInstancesNum": 1, "Instances": []map[string]string{{"InstanceId": "redis-test"}}})
		default:
			return 200, result(redisFixture("redis-test"))
		}
	})
	for n := 0; n < 2; n++ {
		items, e := r.GetIElasticcaches()
		if e != nil || len(items) != 1 {
			t.Fatalf("%v %v", items, e)
		}
		i := items[0]
		tags, e := i.GetTags()
		if e != nil || tags["owner"] != "test" || i.GetGlobalId() != "redis-test" || i.GetCapacityMB() != 2048 || i.GetZoneId() != r.GetGlobalId()+"/cn-beijing-a" || i.GetPrivateConnectPort() != 6379 || i.GetPrivateDNS() != "private.example.invalid" || i.GetProjectId() != "test-project" || i.GetArchType() != "master" || i.GetAuthMode() != "on" {
			t.Fatalf("incorrect Redis mapping: %+v", i)
		}
		if i.Delete() == nil || i.Restart() == nil || i.SetTags(nil, false) == nil {
			t.Fatal("write accepted")
		}
	}
}
func TestRedisDetailMismatchFails(t *testing.T) {
	r := inventoryTestRegion(t, func(s, a string, p map[string]interface{}) (int, interface{}) {
		return 200, result(redisFixture("wrong"))
	})
	if _, e := r.GetRedisInstance("expected"); e == nil {
		t.Fatal("accepted wrong instance")
	}
}
func TestVKEClusterNodePoolAndNode(t *testing.T) {
	r := inventoryTestRegion(t, func(s, a string, p map[string]interface{}) (int, interface{}) {
		var v interface{}
		switch a {
		case "ListClusters":
			v = map[string]interface{}{"Id": "cluster-test", "Name": "cluster", "KubernetesVersion": "v1.32.0", "ProjectName": "project", "Status": map[string]string{"Phase": "Running"}, "ClusterConfig": map[string]interface{}{"VpcId": "vpc-test", "SubnetIds": []string{"subnet-test"}}}
		case "ListNodes":
			v = map[string]interface{}{"Id": "node-test", "ClusterId": "cluster-test", "NodePoolId": "pool-test", "InstanceId": "ecs-test", "MetadataName": "node.example.invalid"}
		case "ListNodePools":
			v = map[string]interface{}{"Id": "pool-test", "ClusterId": "cluster-test", "AutoScaling": map[string]int{"MinReplicas": 1, "MaxReplicas": 3, "DesiredReplicas": 2}, "NodeConfig": map[string]interface{}{"InstanceTypeIds": []string{"ecs.g3i.large"}, "SubnetIds": []string{"subnet-test"}, "SystemVolume": map[string]int{"Size": 60}}}
		default:
			t.Fatal(a)
		}
		return 200, result(map[string]interface{}{"TotalCount": 1, "Items": []interface{}{v}})
	})
	cs, e := r.GetICloudKubeClusters()
	if e != nil || len(cs) != 1 {
		t.Fatalf("%v %v", cs, e)
	}
	c := cs[0]
	if c.GetVersion() != "v1.32.0" || c.GetVpcId() != "vpc-test" {
		t.Fatal("cluster mapping")
	}
	ns, e := c.GetIKubeNodes()
	if e != nil || len(ns) != 1 || ns[0].GetINodePoolId() != "pool-test" {
		t.Fatalf("node mapping %v", e)
	}
	ps, e := c.GetIKubeNodePools()
	if e != nil || len(ps) != 1 || ps[0].GetDesiredInstanceCount() != 2 || ps[0].GetRootDiskSizeGb() != 60 {
		t.Fatalf("pool mapping %v", e)
	}
	if _, e := c.GetKubeConfig(true, 30); e != cloudprovider.ErrNotSupported {
		t.Fatal("credential retrieval allowed")
	}
	if c.Delete(false) == nil || ps[0].Delete() == nil {
		t.Fatal("write allowed")
	}
}
func TestCLBRelationsAndReadOnly(t *testing.T) {
	r := inventoryTestRegion(t, func(s, a string, p map[string]interface{}) (int, interface{}) {
		switch a {
		case "DescribeLoadBalancers":
			return 200, result(map[string]interface{}{"TotalCount": 1, "LoadBalancers": []map[string]string{{"LoadBalancerId": "clb-test"}}})
		case "DescribeLoadBalancerAttributes":
			return 200, result(map[string]interface{}{"LoadBalancerId": "clb-test", "Status": "Active", "Type": "private", "VpcId": "vpc-test", "SubnetId": "subnet-test", "EniAddress": "10.0.0.2", "ServerGroups": []map[string]string{{"ServerGroupId": "group-test"}}})
		case "DescribeListeners":
			if p["LoadBalancerId"] != "clb-test" {
				t.Fatal("missing parent")
			}
			return 200, result(map[string]interface{}{"TotalCount": 1, "Listeners": []map[string]interface{}{{"ListenerId": "lsn-test", "LoadBalancerId": "clb-test", "ServerGroupId": "group-test", "Port": 443, "Protocol": "TCP", "Enabled": "on", "Status": "Active"}}})
		case "DescribeServerGroupAttributes":
			return 200, result(map[string]interface{}{"ServerGroupId": "group-test", "Servers": []map[string]interface{}{{"ServerId": "backend-test", "InstanceId": "ecs-test", "Ip": "10.0.0.3", "Port": 8443, "Weight": 100, "Type": "ecs"}}})
		}
		t.Fatal(a)
		return 500, nil
	})
	lbs, e := r.GetILoadBalancers()
	if e != nil || len(lbs) != 1 {
		t.Fatalf("%v %v", lbs, e)
	}
	lb := lbs[0]
	if lb.GetAddress() != "10.0.0.2" || lb.GetVpcId() != "vpc-test" {
		t.Fatal("lb mapping")
	}
	ls, e := lb.GetILoadBalancerListeners()
	if e != nil || len(ls) != 1 || ls[0].GetBackendGroupId() != "group-test" || ls[0].GetListenerPort() != 443 {
		t.Fatalf("listener mapping %v", e)
	}
	gs, e := lb.GetILoadBalancerBackendGroups()
	if e != nil || len(gs) != 1 {
		t.Fatalf("groups %v", e)
	}
	bs, e := gs[0].GetILoadbalancerBackends()
	if e != nil || len(bs) != 1 || bs[0].GetBackendId() != "ecs-test" || bs[0].GetPort() != 8443 {
		t.Fatalf("backends %v", e)
	}
	if lb.Start() == nil || ls[0].Start() == nil {
		t.Fatal("write allowed")
	}
}
func TestInventoryValidEmpty(t *testing.T) {
	r := inventoryTestRegion(t, func(s, a string, p map[string]interface{}) (int, interface{}) {
		return 200, result(map[string]interface{}{"Items": []interface{}{}, "TotalCount": 0})
	})
	items, e := r.inventoryPages("vke", "ListClusters", "Items", "TotalCount", "Id", nil)
	if e != nil || len(items) != 0 {
		t.Fatalf("%v %v", items, e)
	}
}

func TestRedisOmittedEmptyInstances(t *testing.T) {
	for _, total := range []int{0, 1} {
		r := inventoryTestRegion(t, func(s, a string, p map[string]interface{}) (int, interface{}) {
			return 200, result(map[string]interface{}{"TotalInstancesNum": total})
		})
		items, e := r.inventoryPages("redis", "DescribeDBInstances", "Instances", "TotalInstancesNum", "InstanceId", nil)
		if total == 0 {
			if e != nil || len(items) != 0 {
				t.Fatalf("explicit empty region %v", e)
			}
		} else if e == nil {
			t.Fatal("missing positive-total list accepted")
		}
	}
}

func TestRedisPasswordFreeAccessSemantics(t *testing.T) {
	for input, want := range map[string]string{"open": "off", "OPEN": "off", "close": "on", "Close": "on", "": "unknown"} {
		i := SRedisInstance{VpcAuthMode: input}
		if got := i.GetAuthMode(); got != want {
			t.Fatalf("VpcAuthMode %q mapped to %q, want %q", input, got, want)
		}
	}
}
