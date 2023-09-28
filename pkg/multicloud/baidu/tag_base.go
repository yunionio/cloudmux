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

package baidu

type SBaiduTag struct{}

func (tag SBaiduTag) GetName() string {
	return ""
}

func (tag SBaiduTag) GetDescription() string {
	return ""
}

func (tag *SBaiduTag) GetTags() (map[string]string, error) {
	return nil, nil
}

func (tag *SBaiduTag) GetSysTags() map[string]string {
	return nil
}

func (tag *SBaiduTag) SetTags(tags map[string]string, replace bool) error {
	return nil
}
