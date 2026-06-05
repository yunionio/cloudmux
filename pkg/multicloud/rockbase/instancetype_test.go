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

package rockbase

import (
	"testing"

	"yunion.io/x/jsonutils"
)

func TestParseAvailableInstanceTypes(t *testing.T) {
	raw := `{
		"AvailableInstanceTypes": [
			{
				"Zone": "hk-01",
				"Name": "N",
				"ParentType": "N2",
				"Status": "Normal",
				"Description": "General N",
				"MachineSizes": [
					{
						"Gpu": 0,
						"Collection": [
							{"Cpu": 1, "Memory": [1, 2, 4, 8]},
							{"Cpu": 2, "Memory": [2, 4, 8, 16]}
						]
					},
					{
						"Gpu": 1,
						"Collection": [
							{"Cpu": 8, "Memory": [32]}
						]
					}
				]
			}
		],
		"RetCode": 0
	}`
	obj, err := jsonutils.ParseString(raw)
	if err != nil {
		t.Fatalf("ParseString() error = %v", err)
	}
	types, err := parseAvailableInstanceTypes(obj)
	if err != nil {
		t.Fatalf("parseAvailableInstanceTypes() error = %v", err)
	}
	if len(types) != 1 {
		t.Fatalf("len(types) = %d, want 1", len(types))
	}
	if len(types[0].MachineSizes[0].Collection[0].MemoryGB) != 4 {
		t.Fatalf("MemoryGB len = %d, want 4", len(types[0].MachineSizes[0].Collection[0].MemoryGB))
	}

	packages := flattenInstancePackages(types)
	if len(packages) != 9 {
		t.Fatalf("flattenInstancePackages() len = %d, want 9", len(packages))
	}
	if packages[0].Spec != "N2.c1.m1" {
		t.Fatalf("packages[0].Spec = %s, want N2.c1.m1", packages[0].Spec)
	}
	if packages[8].Spec != "N2.c8.m32.g1" {
		t.Fatalf("packages[8].Spec = %s, want N2.c8.m32.g1", packages[8].Spec)
	}
}

func TestParseMemoryGBSingleInt(t *testing.T) {
	obj, err := jsonutils.ParseString(`{"Cpu": 1, "Memory": 1024}`)
	if err != nil {
		t.Fatalf("ParseString() error = %v", err)
	}
	memGB := parseMemoryGB(obj)
	if len(memGB) != 1 || memGB[0] != 1 {
		t.Fatalf("parseMemoryGB() = %v, want [1]", memGB)
	}
}

func TestParseMemoryGBWithMinimalCpuPlatform(t *testing.T) {
	obj, err := jsonutils.ParseString(`{
		"Cpu": 32,
		"Memory": [256],
		"MinimalCpuPlatform": ["Intel/CascadeLake", "Intel/IceLake"]
	}`)
	if err != nil {
		t.Fatalf("ParseString() error = %v", err)
	}
	col := parseInstanceCollection(obj)
	if len(col.MemoryGB) != 1 || col.MemoryGB[0] != 256 {
		t.Fatalf("MemoryGB = %v, want [256]", col.MemoryGB)
	}
	if len(col.MinimalCpuPlatform) != 2 {
		t.Fatalf("MinimalCpuPlatform len = %d, want 2", len(col.MinimalCpuPlatform))
	}
}
