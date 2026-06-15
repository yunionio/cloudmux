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
	"fmt"
	"strings"

	"yunion.io/x/pkg/util/shellutils"

	"yunion.io/x/cloudmux/pkg/multicloud/rockbase"
)

func init() {
	type TagGetOptions struct {
		ID []string `help:"resource id"`
	}
	shellutils.R(&TagGetOptions{}, "tag-show", "Show tags of resources", func(cli *rockbase.SRegion, args *TagGetOptions) error {
		tags, err := cli.FetchResourceTags(args.ID)
		if err != nil {
			return err
		}
		for id, tag := range tags {
			fmt.Println(id, tag)
		}
		return nil
	})

	type TagSetOptions struct {
		ID      string   `help:"resource id"`
		Tag     []string `help:"tag to set, key:value"`
		Replace bool     `help:"replace all tags"`
	}
	shellutils.R(&TagSetOptions{}, "tag-set", "Set tags of a resource", func(cli *rockbase.SRegion, args *TagSetOptions) error {
		tags := make(map[string]string)
		for _, t := range args.Tag {
			parts := strings.SplitN(t, ":", 2)
			if len(parts) != 2 {
				return fmt.Errorf("invalid tag format %q, want key:value", t)
			}
			tags[parts[0]] = parts[1]
		}
		return cli.SetResourceTags(args.ID, tags, args.Replace)
	})
}
