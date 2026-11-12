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
	"testing"
)

func rdsFixture(service, id string) map[string]interface{} {
	engine := "MySQL_8_0"
	port := "3306"
	if service == "rds-mssql" {
		engine = "SQLServer_2022_Std"
		port = "1433"
	}
	base := map[string]interface{}{"InstanceId": id, "RegionId": "cn-beijing", "InstanceName": "database-test", "InstanceStatus": "Running", "DBEngineVersion": engine, "VCPU": 4, "Memory": 8, "StorageSpace": 100, "ZoneId": "cn-beijing-a", "VpcId": "vpc-test", "SubnetId": "subnet-test", "Tags": []map[string]string{{"Key": "owner", "Value": "test"}}}
	endpoint := map[string]interface{}{"EndpointType": "Primary", "ReadWriteMode": "ReadWrite", "Address": []map[string]string{{"NetworkType": "Private", "Domain": "database.example.invalid", "Port": port, "IPAddress": "10.0.0.2", "SubnetId": "subnet-test"}}}
	f := map[string]interface{}{"BasicInfo": base, "ChargeDetail": map[string]string{"ChargeType": "PostPaid"}}
	if service == "rds-mssql" {
		f["ConnectionInfo"] = []interface{}{endpoint}
	} else {
		f["Endpoints"] = []interface{}{endpoint}
	}
	return f
}
func TestRDSAllEnginesAndAtomicFailure(t *testing.T) {
	for _, denied := range []bool{false, true} {
		r := inventoryTestRegion(t, func(service, action string, params map[string]interface{}) (int, interface{}) {
			if denied && service == "rds-mysql" {
				return 403, map[string]interface{}{"ResponseMetadata": map[string]interface{}{"Error": map[string]string{"Code": "AccessDenied", "Message": "denied"}}}
			}
			if action == "DescribeRegions" {
				return 200, result(map[string]interface{}{"Regions": []map[string]string{{"RegionId": "cn-beijing"}}})
			}
			id, array := "postgres-test", "Instances"
			if service == "rds-mysql" {
				id = "mysql-test"
			}
			if service == "rds-mssql" {
				id, array = "mssql-test", "InstancesInfo"
			}
			if action == "DescribeDBInstances" {
				return 200, result(map[string]interface{}{"Total": 1, array: []map[string]string{{"InstanceId": id}}})
			}
			if service == "rds-postgresql" {
				return 200, result(postgresqlFixture(id))
			}
			return 200, result(rdsFixture(service, id))
		})
		rows, e := r.GetIDBInstances()
		if denied {
			if e == nil || rows != nil {
				t.Fatal("returned partial engines after denied request")
			}
			continue
		}
		if e != nil || len(rows) != 3 {
			t.Fatalf("engines %v %v", rows, e)
		}
		mysql, mssql := rows[1], rows[2]
		if mysql.GetEngine() != "MySQL" || mysql.GetEngineVersion() != "8.0" || mysql.GetPort() != 3306 || mysql.GetVmemSizeMB() != 8192 {
			t.Fatalf("MySQL mapping %s %s %d", mysql.GetEngine(), mysql.GetEngineVersion(), mysql.GetPort())
		}
		if mssql.GetEngine() != "SQLServer" || mssql.GetPort() != 1433 || mssql.GetInternalConnectionStr() != "database.example.invalid" {
			t.Fatalf("SQL Server mapping %s %d", mssql.GetEngine(), mssql.GetPort())
		}
		tags, _ := mssql.GetTags()
		if tags["owner"] != "test" || mssql.Reboot() == nil || mysql.Delete() == nil {
			t.Fatal("RDS tags or unsupported writes")
		}
	}
}
func TestRDSMissingEngineRegionsPreservePostgreSQL(t *testing.T) {
	r := inventoryTestRegion(t, func(service, action string, p map[string]interface{}) (int, interface{}) {
		if action == "DescribeRegions" {
			region := "cn-beijing"
			if service != "rds-postgresql" {
				region = "cn-shanghai"
			}
			return 200, result(map[string]interface{}{"Regions": []map[string]string{{"RegionId": region}}})
		}
		if action == "DescribeDBInstances" {
			return 200, result(map[string]interface{}{"Total": 1, "Instances": []map[string]string{{"InstanceId": "postgres-test"}}})
		}
		return 200, result(postgresqlFixture("postgres-test"))
	})
	rows, e := r.GetIDBInstances()
	if e != nil || len(rows) != 1 || rows[0].GetEngine() != "PostgreSQL" {
		t.Fatalf("%v %v", rows, e)
	}
}
func TestRDSDetailIdentityAndIDRouting(t *testing.T) {
	for _, service := range []string{"rds-mysql", "rds-mssql"} {
		r := inventoryTestRegion(t, func(s, a string, p map[string]interface{}) (int, interface{}) {
			return 200, result(rdsFixture(service, "wrong"))
		})
		if _, e := r.getReadonlyRDSInstance(service, "expected"); e == nil {
			t.Fatal("wrong RDS detail ID accepted")
		}
	}
	r := inventoryTestRegion(t, func(s, a string, p map[string]interface{}) (int, interface{}) {
		return 200, result(rdsFixture(s, p["InstanceId"].(string)))
	})
	for _, id := range []string{"mysql-test", "mssql-test"} {
		i, e := r.GetIDBInstanceById(id)
		if e != nil || i.GetGlobalId() != id {
			t.Fatalf("RDS lookup %v", e)
		}
	}
}
