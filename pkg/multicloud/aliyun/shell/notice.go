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

package shell

import (
	"yunion.io/x/pkg/util/shellutils"

	"yunion.io/x/cloudmux/pkg/multicloud/aliyun"
)

type sNoticeRow struct {
	Title   string
	Content string
}

func init() {
	type NoticeListOptions struct {
	}
	shellutils.R(&NoticeListOptions{}, "notice-list", "list announcements", func(cli *aliyun.SRegion, args *NoticeListOptions) error {
		notices, err := cli.GetClient().GetNotices()
		if err != nil {
			return err
		}
		rows := make([]sNoticeRow, 0, len(notices))
		for _, notice := range notices {
			rows = append(rows, sNoticeRow{
				Title:   notice.GetTitle(),
				Content: notice.GetContent(),
			})
		}
		printList(rows, len(rows), 0, 0, []string{})
		return nil
	})
}
