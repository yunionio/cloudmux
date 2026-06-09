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
)

func TestParseInstanceType(t *testing.T) {
	cases := []struct {
		name         string
		instanceType string
		want         SInstanceType
		wantErr      bool
	}{
		{
			name:         "O.c1.m1",
			instanceType: "O.c1.m1",
			want: SInstanceType{
				UHostType: "O",
				CPU:       1,
				MemoryMB:  1024,
				GPU:       0,
			},
		},
		{
			name:         "T4S.c4.m8.g1",
			instanceType: "T4S.c4.m8.g1",
			want: SInstanceType{
				UHostType: "G",
				GpuType:   "T4S",
				CPU:       4,
				MemoryMB:  8192,
				GPU:       1,
			},
		},
		{
			name:         "2080Ti-4C.c8.m32.g2",
			instanceType: "2080Ti-4C.c8.m32.g2",
			want: SInstanceType{
				UHostType: "G",
				GpuType:   "2080Ti-4C",
				CPU:       8,
				MemoryMB:  32768,
				GPU:       2,
			},
		},
		{
			name:         "invalid format",
			instanceType: "O.c1",
			wantErr:      true,
		},
		{
			name:         "extra segment",
			instanceType: "O.c1.m1.extra",
			wantErr:      true,
		},
		{
			name:         "missing prefix",
			instanceType: "c1.m1",
			wantErr:      true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := ParseInstanceType(tc.instanceType)
			if tc.wantErr {
				if err == nil {
					t.Fatalf("ParseInstanceType(%q) error = nil, want error", tc.instanceType)
				}
				return
			}
			if err != nil {
				t.Fatalf("ParseInstanceType(%q) error = %v", tc.instanceType, err)
			}
			if got != tc.want {
				t.Fatalf("ParseInstanceType(%q) = %+v, want %+v", tc.instanceType, got, tc.want)
			}
		})
	}
}
