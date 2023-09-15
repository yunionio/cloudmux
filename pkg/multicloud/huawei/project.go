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

import (
	"strings"

	api "yunion.io/x/cloudmux/pkg/apis/compute"
	"yunion.io/x/cloudmux/pkg/cloudprovider"
	"yunion.io/x/pkg/errors"
)

// https://support.huaweicloud.com/api-iam/zh-cn_topic_0057845625.html
type SProject struct {
	client *SHuaweiClient

	IsDomain    bool   `json:"is_domain"`
	Description string `json:"description"`
	Enabled     bool   `json:"enabled"`
	ID          string `json:"id"`
	ParentID    string `json:"parent_id"`
	DomainID    string `json:"domain_id"`
	Name        string `json:"name"`
}

func (self *SProject) GetRegionID() string {
	return strings.Split(self.Name, "_")[0]
}

func (self *SProject) GetDescription() string {
	return self.Description
}

func (self *SProject) GetHealthStatus() string {
	if self.Enabled {
		return api.CLOUD_PROVIDER_HEALTH_NORMAL
	}
	return api.CLOUD_PROVIDER_HEALTH_SUSPENDED
}

// obs 权限必须赋予到mos project之上
func (self *SHuaweiClient) GetMosProjectId() string {
	projects, err := self.GetProjects()
	if err != nil {
		return ""
	}
	for i := range projects {
		if strings.ToLower(projects[i].Name) == "mos" {
			return projects[i].ID
		}
	}
	return ""
}

func (self *SHuaweiClient) GetProjectById(projectId string) (*SProject, error) {
	projects, err := self.GetProjects()
	if err != nil {
		return nil, err
	}

	for i := range projects {
		if projects[i].ID == projectId {
			return &projects[i], nil
		}
	}
	return nil, errors.Wrapf(cloudprovider.ErrNotFound, projectId)
}

func (self *SHuaweiClient) GetProjects() ([]SProject, error) {
	if len(self.projects) > 0 {
		return self.projects, nil
	}
	var err error
	self.projects, err = self.getProjects()
	return self.projects, err
}
