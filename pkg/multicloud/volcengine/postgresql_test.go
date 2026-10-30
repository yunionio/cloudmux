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
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"yunion.io/x/pkg/errors"

	billing "yunion.io/x/cloudmux/pkg/apis/billing"
	api "yunion.io/x/cloudmux/pkg/apis/compute"
	"yunion.io/x/cloudmux/pkg/cloudprovider"
)

type postgresqlRoundTripper func(*http.Request) (*http.Response, error)

func (f postgresqlRoundTripper) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }

func postgresqlResponse(status int, body interface{}) *http.Response {
	encoded, _ := json.Marshal(body)
	return &http.Response{StatusCode: status, Header: http.Header{}, Body: io.NopCloser(strings.NewReader(string(encoded)))}
}

func postgresqlFixture(id string) map[string]interface{} {
	return map[string]interface{}{
		"BasicInfo": map[string]interface{}{
			"InstanceId": id, "InstanceName": "database", "InstanceStatus": "Running",
			"DBEngineVersion": "PostgreSQL_17", "InstanceType": "HA", "NodeSpec": "rds.postgres.2c4g",
			"VCPU": 2, "Memory": 4, "StorageSpace": 50, "StorageUse": 3 * 1024 * 1024,
			"StorageType": "LocalSSD", "RegionId": "cn-beijing", "ZoneId": "cn-beijing-a",
			"VpcID": "vpc-test", "SubnetId": "subnet-test", "ProjectName": "default",
			"CreateTime":  "2026-01-02T03:04:05.123Z",
			"InstanceTag": []map[string]string{{"Key": "business", "Value": "demo"}},
		},
		"ChargeDetail": map[string]interface{}{"ChargeType": "PrePaid", "AutoRenew": false, "ChargeEndTime": "2027-01-02T03:04:05Z"},
		"Endpoints": []map[string]interface{}{
			{"ReadWriteMode": "ReadOnly", "Address": []map[string]string{{"NetworkType": "Private", "Domain": "replica.example.invalid", "Port": "5432"}}},
			{"ReadWriteMode": "ReadWrite", "Address": []map[string]string{
				{"NetworkType": "Private", "Domain": "primary.example.invalid", "IPAddress": "10.0.0.5", "Port": "5432"},
				{"NetworkType": "Public", "Domain": "public.example.invalid", "Port": "5432"},
			}},
		},
		"Nodes": []map[string]string{{"ZoneId": "cn-beijing-a"}, {"ZoneId": "cn-beijing-b"}, {"ZoneId": "cn-beijing-b"}},
	}
}

func newPostgresqlTestRegion(t *testing.T, handler func(string, map[string]interface{}) (int, interface{})) *SRegion {
	t.Helper()
	client := &SVolcEngineClient{VolcEngineClientConfig: NewVolcEngineClientConfig("test-access-key", "test-secret-key")}
	client.postgresqlClient = &sPostgresqlClient{httpClient: &http.Client{Transport: postgresqlRoundTripper(func(req *http.Request) (*http.Response, error) {
		if req.Method != http.MethodPost || req.URL.Scheme != "https" || req.URL.Host != "rds-postgresql.cn-beijing.volcengineapi.com" {
			t.Fatalf("unexpected PostgreSQL request: %s %s", req.Method, req.URL)
		}
		if req.Header.Get("Content-Type") != "application/json" || req.URL.Query().Get("Version") != postgresqlAPIVersion {
			t.Fatal("missing JSON request or API version")
		}
		if !strings.Contains(req.Header.Get("Authorization"), "/cn-beijing/rds_postgresql/request") {
			t.Fatal("request was not signed for the PostgreSQL service and region")
		}
		var params map[string]interface{}
		if err := json.NewDecoder(req.Body).Decode(&params); err != nil {
			t.Fatal(err)
		}
		status, body := handler(req.URL.Query().Get("Action"), params)
		return postgresqlResponse(status, body), nil
	})}}
	return &SRegion{client: client, RegionId: "cn-beijing"}
}

func regionCatalog() map[string]interface{} {
	return map[string]interface{}{"Result": map[string]interface{}{"Regions": []map[string]string{{"RegionId": "cn-beijing"}}}}
}

func TestPostgresqlInventoryPaginationAndMapping(t *testing.T) {
	pages, catalogs := []int{}, 0
	region := newPostgresqlTestRegion(t, func(action string, params map[string]interface{}) (int, interface{}) {
		switch action {
		case "DescribeRegions":
			catalogs++
			return 200, regionCatalog()
		case "DescribeDBInstances":
			page := int(params["PageNumber"].(float64))
			pages = append(pages, page)
			if params["PageSize"] != float64(100) {
				t.Fatal("unexpected page size")
			}
			ids := []map[string]string{{"InstanceId": "postgres-a"}, {"InstanceId": "postgres-b"}}
			if page == 2 {
				ids = []map[string]string{{"InstanceId": "postgres-c"}}
			}
			return 200, map[string]interface{}{"Result": map[string]interface{}{"Total": 3, "Instances": ids}}
		case "DescribeDBInstanceDetail":
			return 200, map[string]interface{}{"Result": postgresqlFixture(params["InstanceId"].(string))}
		default:
			t.Fatalf("unexpected action %s", action)
			return 500, nil
		}
	})
	for iteration := 0; iteration < 2; iteration++ {
		instances, err := region.GetIDBInstances()
		if err != nil {
			t.Fatal(err)
		}
		if len(instances) != 3 {
			t.Fatalf("got %d instances", len(instances))
		}
		ids := []string{}
		for _, instance := range instances {
			ids = append(ids, instance.GetGlobalId())
		}
		if !reflect.DeepEqual(ids, []string{"postgres-a", "postgres-b", "postgres-c"}) {
			t.Fatal(ids)
		}
		instance := instances[0]
		if instance.GetEngine() != api.DBINSTANCE_TYPE_POSTGRESQL || instance.GetEngineVersion() != "17" || instance.GetStatus() != api.DBINSTANCE_RUNNING {
			t.Fatal("incorrect engine or status")
		}
		if instance.GetVcpuCount() != 2 || instance.GetVmemSizeMB() != 4096 || instance.GetDiskSizeGB() != 50 {
			t.Fatal("incorrect instance sizing")
		}
		if instance.GetProjectId() != "default" || instance.GetIVpcId() != "vpc-test" || instance.GetZone1Id() != "cn-beijing-a" || instance.GetZone2Id() != "cn-beijing-b" || instance.GetZone3Id() != "" {
			t.Fatal("incorrect project or network ownership")
		}
		if instance.GetPort() != 5432 || instance.GetInternalConnectionStr() != "primary.example.invalid" || instance.GetConnectionStr() != "public.example.invalid" {
			t.Fatal("wrong primary endpoint")
		}
		networks, err := instance.GetDBNetworks()
		if err != nil || len(networks) != 1 || networks[0].NetworkId != "subnet-test" || networks[0].IP != "10.0.0.5" {
			t.Fatalf("incorrect private networks: %+v %v", networks, err)
		}
		tags, err := instance.GetTags()
		if err != nil || tags["business"] != "demo" {
			t.Fatal("missing tags")
		}
		if instance.GetBillingType() != billing.BILLING_TYPE_PREPAID || instance.IsAutoRenew() || instance.GetCreatedAt().Year() != 2026 || instance.GetExpiredAt().Year() != 2027 {
			t.Fatal("incorrect lifecycle fields")
		}
	}
	if catalogs != 1 || !reflect.DeepEqual(pages, []int{1, 2, 1, 2}) {
		t.Fatalf("catalogs=%d pages=%v", catalogs, pages)
	}
}

func TestPostgresqlInventoryDoesNotReturnPartialOrMalformedSets(t *testing.T) {
	cases := []struct {
		name   string
		result interface{}
		status int
	}{
		{"missing total", map[string]interface{}{"Result": map[string]interface{}{"Instances": []interface{}{}}}, 200},
		{"missing instances", map[string]interface{}{"Result": map[string]interface{}{"Total": 0}}, 200},
		{"null instances", map[string]interface{}{"Result": map[string]interface{}{"Total": 0, "Instances": nil}}, 200},
		{"premature empty page", map[string]interface{}{"Result": map[string]interface{}{"Total": 1, "Instances": []interface{}{}}}, 200},
		{"negative total", map[string]interface{}{"Result": map[string]interface{}{"Total": -1, "Instances": []interface{}{}}}, 200},
		{"duplicate ID", map[string]interface{}{"Result": map[string]interface{}{"Total": 2, "Instances": []map[string]string{{"InstanceId": "a"}, {"InstanceId": "a"}}}}, 200},
		{"empty ID", map[string]interface{}{"Result": map[string]interface{}{"Total": 1, "Instances": []map[string]string{{"InstanceId": ""}}}}, 200},
		{"permission denied", map[string]interface{}{"ResponseMetadata": map[string]interface{}{"Error": map[string]string{"Code": "AccessDenied", "Message": "permission required"}}}, 403},
		{"invalid envelope", map[string]interface{}{}, 200},
		{"http failure", map[string]interface{}{}, 503},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			region := newPostgresqlTestRegion(t, func(action string, params map[string]interface{}) (int, interface{}) {
				if action == "DescribeRegions" {
					return 200, regionCatalog()
				}
				return tc.status, tc.result
			})
			instances, err := region.GetIDBInstances()
			if err == nil || instances != nil {
				t.Fatalf("unsafe successful inventory: %v %v", instances, err)
			}
		})
	}
	region := newPostgresqlTestRegion(t, func(action string, params map[string]interface{}) (int, interface{}) {
		if action == "DescribeRegions" {
			return 200, regionCatalog()
		}
		if params["PageNumber"] == float64(1) {
			return 200, map[string]interface{}{"Result": map[string]interface{}{"Total": 2, "Instances": []map[string]string{{"InstanceId": "a"}}}}
		}
		return 403, map[string]interface{}{"ResponseMetadata": map[string]interface{}{"Error": map[string]string{"Code": "AccessDenied"}}}
	})
	if instances, err := region.GetIDBInstances(); err == nil || instances != nil {
		t.Fatal("second page failure returned a partial inventory")
	}
}

func TestPostgresqlDetailIdentityAndPermissionFailures(t *testing.T) {
	for _, variant := range []string{"wrong ID", "wrong region", "permission"} {
		t.Run(variant, func(t *testing.T) {
			region := newPostgresqlTestRegion(t, func(action string, params map[string]interface{}) (int, interface{}) {
				if action == "DescribeRegions" {
					return 200, regionCatalog()
				}
				if action == "DescribeDBInstances" {
					return 200, map[string]interface{}{"Result": map[string]interface{}{"Total": 1, "Instances": []map[string]string{{"InstanceId": "postgres-a"}}}}
				}
				if variant == "permission" {
					return 403, map[string]interface{}{"ResponseMetadata": map[string]interface{}{"Error": map[string]string{"Code": "AccessDenied"}}}
				}
				data := postgresqlFixture("postgres-a")
				basic := data["BasicInfo"].(map[string]interface{})
				if variant == "wrong ID" {
					basic["InstanceId"] = "other"
				} else {
					basic["RegionId"] = "cn-shanghai"
				}
				return 200, map[string]interface{}{"Result": data}
			})
			if instances, err := region.GetIDBInstances(); err == nil || instances != nil {
				t.Fatal("invalid detail returned a successful inventory")
			}
		})
	}
}

func TestPostgresqlEmptyInventoryAndUnsupportedRegion(t *testing.T) {
	listCalls := 0
	region := newPostgresqlTestRegion(t, func(action string, params map[string]interface{}) (int, interface{}) {
		if action == "DescribeRegions" {
			return 200, regionCatalog()
		}
		listCalls++
		return 200, map[string]interface{}{"Result": map[string]interface{}{"Total": 0, "Instances": []interface{}{}}}
	})
	instances, err := region.GetIDBInstances()
	if err != nil || len(instances) != 0 || listCalls != 1 {
		t.Fatalf("empty list: %v %v", instances, err)
	}
	region.RegionId = "ap-southeast-3"
	instances, err = region.GetIDBInstances()
	if instances != nil || errors.Cause(err) != cloudprovider.ErrNotSupported || listCalls != 1 {
		t.Fatal("unsupported region was treated as an empty inventory")
	}
}

func TestPostgresqlReadOnlyAndMissingOptionalValues(t *testing.T) {
	region := newPostgresqlTestRegion(t, func(action string, params map[string]interface{}) (int, interface{}) {
		t.Fatal("mutation must not issue an API request")
		return 500, nil
	})
	instance := &SPostgresqlInstance{region: region, sPostgresqlBasicInfo: sPostgresqlBasicInfo{InstanceId: "postgres-a", InstanceStatus: "FutureStatus"}}
	for _, mutate := range []func() error{instance.Delete, instance.Reboot, instance.OpenPublicConnection, instance.ClosePublicConnection, func() error { return instance.ChangeConfig(context.Background(), nil) }, func() error { return instance.SetTags(map[string]string{"x": "y"}, false) }} {
		if err := mutate(); err == nil {
			t.Fatal("write operation unexpectedly supported")
		}
	}
	if instance.GetName() != "postgres-a" || instance.GetPort() != 0 || !instance.GetCreatedAt().IsZero() || !instance.GetExpiredAt().IsZero() || instance.GetBillingType() != "" || instance.GetStatus() != api.DBINSTANCE_UNKNOWN {
		t.Fatal("invented optional values")
	}
	caps := region.client.GetCapabilities()
	ro, write := false, false
	for _, cap := range caps {
		ro = ro || cap == cloudprovider.CLOUD_CAPABILITY_RDS+cloudprovider.READ_ONLY_SUFFIX
		write = write || cap == cloudprovider.CLOUD_CAPABILITY_RDS
	}
	if !ro || write {
		t.Fatal("RDS capability is not read-only")
	}
	if _, err := region.client.postgresqlRequest("cn-beijing", "DeleteDBInstance", nil); errors.Cause(err) != cloudprovider.ErrNotSupported {
		t.Fatal("mutation action was not rejected")
	}
	if _, err := region.client.postgresqlRequest("cn-beijing@other.invalid/path", "DescribeRegions", nil); err == nil {
		t.Fatal("invalid region was not rejected")
	}
}

func TestPostgresqlTLSVerification(t *testing.T) {
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { fmt.Fprint(w, `{"Result":{}}`) }))
	defer server.Close()
	client := &SVolcEngineClient{VolcEngineClientConfig: NewVolcEngineClientConfig("test-access-key", "test-secret-key")}
	pg := client.getPostgresqlClient()
	transport := pg.httpClient.Transport.(*http.Transport)
	transport.Proxy = nil
	transport.DialContext = func(ctx context.Context, network, address string) (net.Conn, error) {
		return (&net.Dialer{}).DialContext(ctx, network, server.Listener.Addr().String())
	}
	pg.httpClient.Timeout = 3 * time.Second
	if _, err := client.postgresqlRequest("cn-beijing", "DescribeRegions", map[string]string{}); err == nil {
		t.Fatal("accepted an untrusted TLS certificate")
	}
}

func TestPostgresqlConcurrentRegionDiscovery(t *testing.T) {
	var calls int32
	region := newPostgresqlTestRegion(t, func(action string, params map[string]interface{}) (int, interface{}) {
		if action != "DescribeRegions" {
			t.Errorf("unexpected action %s", action)
		}
		atomic.AddInt32(&calls, 1)
		return 200, regionCatalog()
	})
	var workers sync.WaitGroup
	for i := 0; i < 8; i++ {
		workers.Add(1)
		go func() {
			defer workers.Done()
			supported, err := region.client.isPostgresqlRegion("cn-beijing")
			if err != nil || !supported {
				t.Errorf("region discovery: %v %v", supported, err)
			}
		}()
	}
	workers.Wait()
	if calls != 1 {
		t.Fatalf("region discovery made %d requests", calls)
	}
}
