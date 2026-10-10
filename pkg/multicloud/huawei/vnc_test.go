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

package huawei

import "testing"

func TestConvertNovncHtmlUrlToWebsocket(t *testing.T) {
	cases := []struct {
		in   string
		want string
	}{
		{
			in:   "https://console-intl.huaweicloud.com:443/vnc/region.ap-southeast-1/vnc_auto.html?token=0fda3eca-8232-4249-abcd-1234567890ab&lang=EN",
			want: "wss://console-intl.huaweicloud.com:443/vnc/region.ap-southeast-1/websockify?token=0fda3eca-8232-4249-abcd-1234567890ab",
		},
		{
			in:   "https://nova-novncproxy.az1.dc1.domainname.com:8002/vnc_auto.html?token=0fda3eca-8232-4249-abcd-1234567890ab&lang=EN&tLength=70",
			want: "wss://nova-novncproxy.az1.dc1.domainname.com:8002/websockify?token=0fda3eca-8232-4249-abcd-1234567890ab",
		},
		{
			in:   "http://192.168.1.10:6080/vnc_auto.html?token=8cb819e5-5882-41c9-9ab3-17609160e6d0",
			want: "ws://192.168.1.10:6080/websockify?token=8cb819e5-5882-41c9-9ab3-17609160e6d0",
		},
		{
			in:   "https://host:8002/vnc_auto.html?path=websockify%3Ftoken%3Dabc-123",
			want: "wss://host:8002/websockify?token=abc-123",
		},
	}
	for _, c := range cases {
		got, err := ConvertNovncHtmlUrlToWebsocket(c.in)
		if err != nil {
			t.Fatalf("ConvertNovncHtmlUrlToWebsocket(%q) err: %v", c.in, err)
		}
		if got != c.want {
			t.Fatalf("ConvertNovncHtmlUrlToWebsocket(%q)\n got  %q\n want %q", c.in, got, c.want)
		}
	}
}
