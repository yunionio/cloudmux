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

package esxi

import (
	"testing"
)

func TestCompactIp(t *testing.T) {
	cases := []struct {
		input []string
		want  string
	}{
		{
			input: []string{},
			want:  "",
		},
		{
			input: []string{"192.168.200.2"},
			want:  "192.168.200.2",
		},
		{
			input: []string{"192.168.200.2", "192.168.200.3"},
			want:  "192.168.200.2,3",
		},
		{
			input: []string{"192.168.200.2", "192.168.200.3", "192.168.200.3"},
			want:  "192.168.200.2,3",
		},
		{
			input: []string{"192.168.200.2", "192.168.200.3", "192.168.200.3", "192.168.200.4"},
			want:  "192.168.200.2-4",
		},
		{
			input: []string{"192.168.200.2", "192.168.200.3", "192.168.200.3", "192.168.200.4", "192.168.200.5"},
			want:  "192.168.200.2-5",
		},
		{
			input: []string{"192.168.200.2", "192.168.200.3", "192.168.200.3", "192.168.200.4", "192.168.200.5", "192.168.200.7"},
			want:  "192.168.200.2-5,7",
		},
		{
			input: []string{"192.168.200.1", "192.168.200.3", "192.168.200.3", "192.168.200.4", "192.168.200.5", "192.168.200.7"},
			want:  "192.168.200.1,3-5,7",
		},
		{
			input: []string{"192.168.200.1", "192.168.200.3", "192.168.200.3", "192.168.200.4", "192.168.200.5", "192.168.200.7", "10.127.40.3"},
			want:  "10.127.40.3;192.168.200.1,3-5,7",
		},
		{
			input: []string{"192.168.200.1", "192.168.200.3", "192.168.200.3", "192.168.200.4", "192.168.200.5", "192.168.200.7", "10.127.40.3", "10.127.40.4"},
			want:  "10.127.40.3,4;192.168.200.1,3-5,7",
		},
		{
			input: []string{"192.168.200.1", "192.168.200.3", "192.168.200.3", "192.168.200.4", "192.168.200.5", "192.168.200.7", "10.127.40.3", "10.127.40.4", "10.127.40.5"},
			want:  "10.127.40.3-5;192.168.200.1,3-5,7",
		},
		{
			input: []string{"192.168.200.1", "192.168.200.3", "192.168.200.3", "192.168.200.4", "192.168.200.5", "192.168.200.7", "10.127.40.3", "10.127.40.4", "10.127.40.5", "192.168.201.1"},
			want:  "10.127.40.3-5;192.168.200.1,3-5,7;192.168.201.1",
		},
		{
			input: []string{"10.127.40.40", "10.127.40.43", "10.127.40.5", "10.127.40.35", "10.127.40.41"},
			want:  "10.127.40.5,35,40,41,43",
		},
	}
	for _, c := range cases {
		got := compactIPs(c.input)
		if got != c.want {
			t.Errorf("want: %s got: %s", c.want, got)
		}
	}
}
