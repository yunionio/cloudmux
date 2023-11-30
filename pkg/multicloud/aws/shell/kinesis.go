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

	"yunion.io/x/cloudmux/pkg/multicloud/aws"
)

func init() {
	type KinesisListOptions struct {
	}
	shellutils.R(&KinesisListOptions{}, "kinesis-stream-list", "List kinesis stream", func(cli *aws.SRegion, args *KinesisListOptions) error {
		streams, err := cli.ListStreams()
		if err != nil {
			return err
		}
		printList(streams, 0, 0, 0, []string{})
		return nil
	})

	type KinesisNameOptions struct {
		NAME string
	}
	shellutils.R(&KinesisNameOptions{}, "kinesis-stream-show", "Show kinesis stream", func(cli *aws.SRegion, args *KinesisNameOptions) error {
		stream, err := cli.DescribeStream(args.NAME)
		if err != nil {
			return err
		}
		printObject(stream)
		return nil
	})

}
