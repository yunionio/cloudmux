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

package harvester

import (
	"net/url"
	"time"
)

type SDisk struct {
	Id       string
	Metadata struct {
		CreationTimestamp time.Time
		Fields            []string
		Labels            map[string]string
		Name              string
		Uid               string
	}
	Spec struct {
		AccessModes []string
		Resources   struct {
			Requests struct {
				Storage string
			}
		}
		StorageClassName string
		VolumeMode       string
		VolumeName       string
	}
}

func (region *SRegion) GetDisks() ([]SDisk, error) {
	params := url.Values{}
	params.Set("exclude", "metadata.managedFields")
	resp, err := region.list("v1/harvester/persistentvolumeclaims", params)
	if err != nil {
		return nil, err
	}
	ret := []SDisk{}
	err = resp.Unmarshal(&ret)
	if err != nil {
		return nil, err
	}
	return ret, nil
}
