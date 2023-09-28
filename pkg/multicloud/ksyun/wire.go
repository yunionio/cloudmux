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

package ksyun

import (
	"fmt"

	api "yunion.io/x/cloudmux/pkg/apis/compute"
	"yunion.io/x/cloudmux/pkg/cloudprovider"
	"yunion.io/x/cloudmux/pkg/multicloud"
	"yunion.io/x/pkg/errors"
)

type SWire struct {
	multicloud.SResourceBase
	SKsTag
	zone   *SZone
	region *SRegion
	vpc    *SVpc

	Id        string
	WireId    string
	ZoneId    string
	inetworks []cloudprovider.ICloudNetwork
}

func (wire *SWire) GetId() string {
	return fmt.Sprintf("%s-%s", wire.vpc.GetId(), wire.zone.GetId())
}

func (wire *SWire) GetName() string {
	return wire.GetId()
}

func (wire *SWire) IsEmulated() bool {
	return true
}

func (wire *SWire) GetStatus() string {
	return api.WIRE_STATUS_AVAILABLE
}

func (wire *SWire) Refresh() error {
	return nil
}

func (wire *SWire) GetGlobalId() string {
	return fmt.Sprintf("%s-%s", wire.vpc.GetGlobalId(), wire.zone.GetGlobalId())
}

func (wire *SWire) GetIVpc() cloudprovider.ICloudVpc {
	return wire.vpc
}

func (wire *SWire) GetIZone() cloudprovider.ICloudZone {
	return wire.zone
}

func (wire *SWire) GetINetworks() ([]cloudprovider.ICloudNetwork, error) {
	networks, err := wire.vpc.region.GetNetworks(wire.vpc.VpcID)
	if err != nil {
		return nil, errors.Wrap(err, "GetNetworks")
	}
	inetworks := []cloudprovider.ICloudNetwork{}
	for i := 0; i < len(networks); i++ {
		if networks[i].AvailabilityZoneName == wire.zone.GetName() {
			inetworks = append(inetworks, &networks[i])
		}
	}
	return inetworks, nil
}

func (wire *SWire) GetBandwidth() int {
	return 10000
}

func (wire *SWire) CreateINetwork(opts *cloudprovider.SNetworkCreateOptions) (cloudprovider.ICloudNetwork, error) {
	return nil, cloudprovider.ErrNotImplemented
}

func (wire *SWire) GetINetworkById(netid string) (cloudprovider.ICloudNetwork, error) {
	networks, err := wire.GetINetworks()
	if err != nil {
		return nil, err
	}
	for i := 0; i < len(networks); i += 1 {
		if networks[i].GetGlobalId() == netid {
			return networks[i], nil
		}
	}
	return nil, cloudprovider.ErrNotFound
}

func (wire *SWire) GetDescription() string {
	return ""
}
