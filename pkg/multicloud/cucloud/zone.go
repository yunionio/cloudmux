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

package cucloud

import (
	"yunion.io/x/cloudmux/pkg/cloudprovider"
	"yunion.io/x/pkg/errors"
	"yunion.io/x/pkg/util/httputils"
)

type SZone struct {
	region *SRegion

	host          *SHost
	ZoneName      string
	Status        string
	ZoneId        string
	CloudRegionId string
}

func (region *SRegion) GetZones() ([]SZone, error) {
	resp, err := region.client.request(httputils.GET, "/instance/v1/product/zones", nil)
	if err != nil {
		return nil, errors.Wrap(err, "request zone")
	}
	zones := []SZone{}
	return zones, resp.Unmarshal(&zones, "result", "list")
}

func (zone *SZone) GetIHosts() ([]cloudprovider.ICloudHost, error) {
	return []cloudprovider.ICloudHost{zone.getHost()}, nil
}

func (zone *SZone) getHost() *SHost {
	if zone.host == nil {
		zone.host = &SHost{zone: zone}
	}
	return zone.host
}
