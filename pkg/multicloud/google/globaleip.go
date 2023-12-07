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

package google

import (
	"fmt"
	"time"

	"yunion.io/x/jsonutils"
	"yunion.io/x/log"
	"yunion.io/x/pkg/errors"

	api "yunion.io/x/cloudmux/pkg/apis/compute"
	"yunion.io/x/cloudmux/pkg/cloudprovider"
	"yunion.io/x/cloudmux/pkg/multicloud"
)

type SGlobalAddress struct {
	region *SGlobalRegion
	SResourceBase
	multicloud.SEipBase
	GoogleTags
	instanceId string

	CreationTimestamp time.Time
	Description       string
	Address           string
	Status            string
	Region            string
	Users             []string
	NetworkTier       string
	AddressType       string
	Kind              string
}

func (addr *SGlobalAddress) GetStatus() string {
	switch addr.Status {
	case "RESERVING":
		return api.EIP_STATUS_ASSOCIATE
	case "RESERVED":
		return api.EIP_STATUS_READY
	case "IN_USE":
		return api.EIP_STATUS_READY
	default:
		log.Errorf("Unknown eip status: %s", addr.Status)
		return api.EIP_STATUS_UNKNOWN
	}
}

func (addr *SGlobalAddress) GetIpAddr() string {
	return addr.Address
}

func (addr *SGlobalAddress) GetMode() string {
	if addr.IsEmulated() {
		return api.EIP_MODE_INSTANCE_PUBLICIP
	}
	return api.EIP_MODE_STANDALONE_EIP
}

func (addr *SGlobalAddress) GetBandwidth() int {
	return 0
}

func (addr *SGlobalAddress) GetInternetChargeType() string {
	return api.EIP_CHARGE_TYPE_BY_TRAFFIC
}

func (addr *SGlobalAddress) Delete() error {
	return addr.region.Delete(addr.SelfLink)
}

func (addr *SGlobalAddress) Associate(conf *cloudprovider.AssociateConfig) error {
	return addr.region.AssociateInstanceEip(conf.InstanceId, addr.Address)
}

func (addr *SGlobalAddress) Dissociate() error {
	if len(addr.Users) > 0 {
		return addr.region.DissociateInstanceEip(addr.Users[0], addr.Address)
	}
	return nil
}

func (addr *SGlobalAddress) ChangeBandwidth(bw int) error {
	return cloudprovider.ErrNotSupported
}

func (addr *SGlobalAddress) GetProjectId() string {
	return addr.region.GetProjectId()
}

func (region *SGlobalRegion) GetEips(address string, maxResults int, pageToken string) ([]SGlobalAddress, error) {
	eips := []SGlobalAddress{}
	params := map[string]string{}
	if len(address) > 0 {
		params["filter"] = fmt.Sprintf(`address="%s"`, address)
	}
	resource := "global/addresses"

	err := region.List(resource, params, maxResults, pageToken, &eips)
	if err != nil {
		return nil, err
	}

	for i := range eips {
		eips[i].region = region
	}
	return eips, nil
}

func (addr *SGlobalAddress) GetAssociationExternalId() string {
	if len(addr.instanceId) > 0 {
		return addr.instanceId
	}
	if len(addr.Users) > 0 {
		res := &SResourceBase{}
		err := addr.region.GetBySelfId(addr.Users[0], res)
		if err != nil {
			return ""
		}
		return res.GetGlobalId()
	}
	return ""
}

func (addr *SGlobalAddress) GetAssociationType() string {
	if len(addr.GetAssociationExternalId()) > 0 {
		return api.EIP_ASSOCIATE_TYPE_SERVER
	}
	return ""
}

func (self *SGlobalRegion) AssociateInstanceEip(instanceId string, eip string) error {
	instance, err := self.GetInstance(instanceId)
	if err != nil {
		return errors.Wrap(err, "region.GetInstance")
	}
	for _, networkInterface := range instance.NetworkInterfaces {
		body := map[string]interface{}{
			"type":  "ONE_TO_ONE_NAT",
			"name":  "External NAT",
			"natIP": eip,
		}
		params := map[string]string{"networkInterface": networkInterface.Name}
		return self.Do(instance.SelfLink, "addAccessConfig", params, jsonutils.Marshal(body))
	}
	return fmt.Errorf("no valid networkinterface to associate")
}

func (self *SGlobalRegion) DissociateInstanceEip(instanceId string, eip string) error {
	instance := SInstance{}
	err := self.GetBySelfId(instanceId, &instance)
	if err != nil {
		return errors.Wrap(err, "region.GetInstance")
	}
	for _, networkInterface := range instance.NetworkInterfaces {
		for _, accessConfig := range networkInterface.AccessConfigs {
			if accessConfig.NatIP == eip {
				body := map[string]string{}
				params := map[string]string{
					"networkInterface": networkInterface.Name,
					"accessConfig":     accessConfig.Name,
				}
				return self.Do(instance.SelfLink, "deleteAccessConfig", params, jsonutils.Marshal(body))
			}
		}
	}
	return nil
}

func (region *SGlobalRegion) listAll(method string, resource string, params map[string]string, retval interface{}) error {
	return region.client._ecsListAll(method, resource, params, retval)
}

func (region *SGlobalRegion) ListAll(resource string, params map[string]string, retval interface{}) error {
	return region.listAll("GET", resource, params, retval)
}

func (region *SGlobalRegion) List(resource string, params map[string]string, maxResults int, pageToken string, retval interface{}) error {
	if maxResults == 0 && len(pageToken) == 0 {
		return region.ListAll(resource, params, retval)
	}
	if params == nil {
		params = map[string]string{}
	}
	params["maxResults"] = fmt.Sprintf("%d", maxResults)
	params["pageToken"] = pageToken
	resp, err := region.client.ecsList(resource, params)
	if err != nil {
		return errors.Wrap(err, "ecsList")
	}
	if resp.Contains("items") && retval != nil {
		err = resp.Unmarshal(retval, "items")
		if err != nil {
			return errors.Wrap(err, "resp.Unmarshal")
		}
	}
	return nil
}
