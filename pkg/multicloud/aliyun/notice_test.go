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

package aliyun

import (
	"regexp"
	"strings"
	"testing"
)

var htmlTagLeftRe = regexp.MustCompile(`(?i)</?[a-z][a-z0-9]*\b[^>]*>`)

func TestStripHTMLAliyunRss(t *testing.T) {
	cases := []struct {
		name string
		in   string
		want []string
	}{
		{
			name: "br and headings",
			in:   `<h3>影响时间</h3><p>2026-08-18,&nbsp;调整<br/></p><p></p><h3>变更内容和影响</h3><div>第一行<br>第二行</div>`,
			want: []string{"影响时间", "2026-08-18, 调整", "变更内容和影响", "第一行", "第二行"},
		},
		{
			name: "broken nbsp attrs and link",
			in:   `<p>如您有任何问题，可随时通过<a&nbsp;href="https://smartservice.console.aliyun.com/service/create-ticket">工单</a>&nbsp;或者服务热线与我们联系。</p>`,
			want: []string{
				"工单 (https://smartservice.console.aliyun.com/service/create-ticket)",
				"或者服务热线与我们联系",
			},
		},
		{
			name: "list and nested list",
			in:   `<h4>影响范围</h4><ul><li>变更内容：软件维护</li><li>影响范围：<ul><li>张家口 可用区C</li><li>南通 可用区B</li></ul></li></ul>`,
			want: []string{"影响范围", "- 变更内容：软件维护", "- 影响范围：", "- 张家口 可用区C", "- 南通 可用区B"},
		},
		{
			name: "table cells",
			in:   `<table><tr><td&nbsp;colspan="1">商品及模块</td><td>能力</td><td>调整前</td><td>调整后</td></tr><tr><td>云安全中心</td><td>无代理检测</td><td>0.2元/GB</td><td>1元/GB</td></tr></table>`,
			want: []string{"商品及模块 | 能力 | 调整前 | 调整后", "云安全中心 | 无代理检测 | 0.2元/GB | 1元/GB"},
		},
		{
			name: "version compare entities",
			in:   `<p>6.8.0&nbsp;&lt;=&nbsp;WordPress&nbsp;&lt;=&nbsp;6.8.5</p><p>WordPress&nbsp;&gt;=&nbsp;6.8.6</p>`,
			want: []string{"6.8.0 <= WordPress <= 6.8.5", "WordPress >= 6.8.6"},
		},
		{
			name: "version compare must not eat text between ranges",
			in:   `<p>6.8.0&nbsp;&lt;=&nbsp;WordPress&nbsp;&lt;=&nbsp;6.8.5（仅受 CVE 影响）</p><p>WordPress&nbsp;&gt;=&nbsp;6.8.6</p>`,
			want: []string{"6.8.0 <= WordPress <= 6.8.5（仅受 CVE 影响）", "WordPress >= 6.8.6"},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := stripHTML(c.in)
			for _, part := range c.want {
				if !strings.Contains(got, part) {
					t.Fatalf("missing %q in:\n%s", part, got)
				}
			}
			if htmlTagLeftRe.MatchString(got) {
				t.Fatalf("html tags left: %s", got)
			}
		})
	}
}

func TestFormatRssNoticeContent(t *testing.T) {
	item := &sRssItem{
		Link:           "https://www.aliyun.com/notice/118539",
		ContentEncoded: `<h3>影响时间</h3><p>2026-08-18 调整</p>`,
	}
	got := formatRssNoticeContent(item)
	if !strings.Contains(got, "影响时间") || !strings.Contains(got, "2026-08-18 调整") {
		t.Fatalf("content lost: %s", got)
	}
	if !strings.HasSuffix(strings.TrimSpace(got), "Link: https://www.aliyun.com/notice/118539") {
		t.Fatalf("link missing: %s", got)
	}
}
