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

package ucloud

import "testing"

func TestGetInstanceType(t *testing.T) {
	cases := []struct {
		name string
		inst SInstance
		want string
	}{
		{
			name: "O.c1.m1",
			inst: SInstance{
				UHostType: "O",
				CPU:       1,
				MemoryMB:  1024,
			},
			want: "O.c1.m1",
		},
		{
			name: "N2.c2.m4",
			inst: SInstance{
				UHostType: "N2",
				CPU:       2,
				MemoryMB:  4096,
			},
			want: "N2.c2.m4",
		},
		{
			name: "T4S.c4.m8.g1",
			inst: SInstance{
				MachineType: "G",
				GpuType:     "T4S",
				CPU:         4,
				MemoryMB:    8192,
				GPU:         1,
			},
			want: "T4S.c4.m8.g1",
		},
		{
			name: "machine type fallback",
			inst: SInstance{
				MachineType: "O",
				CPU:         1,
				MemoryMB:    1024,
			},
			want: "O.c1.m1",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := tc.inst.GetInstanceType()
			if got != tc.want {
				t.Fatalf("GetInstanceType() = %q, want %q", got, tc.want)
			}
		})
	}
}
