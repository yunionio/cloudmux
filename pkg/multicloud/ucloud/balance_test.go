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

import (
	"testing"

	"yunion.io/x/jsonutils"
)

func TestParseBalanceFromResp(t *testing.T) {
	cases := []struct {
		name string
		raw  string
		want float64
	}{
		{
			name: "flat response",
			raw: `{
				"Action": "GetBalanceResponse",
				"Amount": 3.99529,
				"AmountAvailable": 2.5,
				"AmountFreeze": 6.51374,
				"RetCode": 0
			}`,
			want: 2.5,
		},
		{
			name: "nested AccountInfo with string amounts",
			raw: `{
				"Action": "GetBalanceResponse",
				"AccountInfo": {
					"Amount": "100.12",
					"AmountAvailable": "88.88"
				},
				"RetCode": 0
			}`,
			want: 88.88,
		},
		{
			name: "fallback to Amount",
			raw: `{
				"Action": "GetBalanceResponse",
				"Amount": 10.5,
				"RetCode": 0
			}`,
			want: 10.5,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			obj, err := jsonutils.ParseString(tc.raw)
			if err != nil {
				t.Fatalf("ParseString() error = %v", err)
			}
			got := parseBalanceFromResp(obj).GetAvailableAmount()
			if got != tc.want {
				t.Fatalf("GetAvailableAmount() = %v, want %v", got, tc.want)
			}
		})
	}
}
