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
	"yunion.io/x/log"
	"yunion.io/x/pkg/util/netutils"

	"yunion.io/x/cloudmux/pkg/cloudprovider"
)

type SInstanceNic struct {
	Instance *SInstance
	Id       string
	IpAddr   string
	MacAddr  string

	Classic bool

	cloudprovider.DummyICloudNic
}

func (nic *SInstanceNic) GetId() string {
	return nic.Id
}

func (nic *SInstanceNic) mustGetId() string {
	if nic.Id == "" {
		panic("empty network interface id")
	}
	return nic.Id
}

func (nic *SInstanceNic) GetIP() string {
	return nic.IpAddr
}

func (nic *SInstanceNic) GetMAC() string {
	if len(nic.MacAddr) > 0 {
		return nic.MacAddr
	}
	ip, _ := netutils.NewIPV4Addr(nic.GetIP())
	return ip.ToMac("00:16:")
}

func (nic *SInstanceNic) InClassicNetwork() bool {
	return nic.Classic == true
}

func (nic *SInstanceNic) GetDriver() string {
	return "virtio"
}

func (nic *SInstanceNic) GetINetworkId() string {
	log.Infoln("this is networkId:", nic.Instance.SubnetID)
	return nic.Instance.SubnetID
}

// func (nic *SInstanceNic) GetSubAddress() ([]string, error) {
// 	nicId := nic.mustGetId()
// 	region := nic.instance.host.zone.region
// 	params := map[string]string{
// 		"RegionId":             region.GetId(),
// 		"NetworkInterfaceId.1": nicId,
// 	}
// 	body, err := region.ecsRequest("DescribeNetworkInterfaces", params)
// 	if err != nil {
// 		return nil, err
// 	}

// 	type DescribeNetworkInterfacesResponse struct {
// 		TotalCount           int    `json:"TotalCount"`
// 		RequestID            string `json:"RequestId"`
// 		PageSize             int    `json:"PageSize"`
// 		NextToken            string `json:"NextToken"`
// 		PageNumber           int    `json:"PageNumber"`
// 		NetworkInterfaceSets struct {
// 			NetworkInterfaceSet []struct {
// 				Status             string `json:"Status"`
// 				PrivateIPAddress   string `json:"PrivateIpAddress"`
// 				ZoneID             string `json:"ZoneId"`
// 				ResourceGroupID    string `json:"ResourceGroupId"`
// 				InstanceID         string `json:"InstanceId"`
// 				VSwitchID          string `json:"VSwitchId"`
// 				NetworkInterfaceID string `json:"NetworkInterfaceId"`
// 				MacAddress         string `json:"MacAddress"`
// 				SecurityGroupIds   struct {
// 					SecurityGroupID []string `json:"SecurityGroupId"`
// 				} `json:"SecurityGroupIds"`
// 				Type     string `json:"Type"`
// 				Ipv6Sets struct {
// 					Ipv6Set []struct {
// 						Ipv6Address string `json:"Ipv6Address"`
// 					} `json:"Ipv6Set"`
// 				} `json:"Ipv6Sets"`
// 				VpcID              string `json:"VpcId"`
// 				OwnerID            string `json:"OwnerId"`
// 				AssociatedPublicIP struct {
// 				} `json:"AssociatedPublicIp"`
// 				CreationTime time.Time `json:"CreationTime"`
// 				Tags         struct {
// 					Tag []struct {
// 						TagKey   string `json:"TagKey"`
// 						TagValue string `json:"TagValue"`
// 					} `json:"Tag"`
// 				} `json:"Tags"`
// 				PrivateIPSets struct {
// 					PrivateIPSet []struct {
// 						PrivateIPAddress   string `json:"PrivateIpAddress"`
// 						AssociatedPublicIP struct {
// 						} `json:"AssociatedPublicIp"`
// 						Primary bool `json:"Primary"`
// 					} `json:"PrivateIpSet"`
// 				} `json:"PrivateIpSets"`
// 			} `json:"NetworkInterfaceSet"`
// 		} `json:"NetworkInterfaceSets"`
// 	}
// 	var resp DescribeNetworkInterfacesResponse
// 	if err := body.Unmarshal(&resp); err != nil {
// 		return nil, errors.Wrapf(err, "unmarshal DescribeNetworkInterfacesResponse: %s", body)
// 	}
// 	if got := len(resp.NetworkInterfaceSets.NetworkInterfaceSet); got != 1 {
// 		return nil, errors.Errorf("got %d element(s) in interface set, expect 1", got)
// 	}
// 	var (
// 		ipAddrs          []string
// 		networkInterface = resp.NetworkInterfaceSets.NetworkInterfaceSet[0]
// 	)
// 	if got := networkInterface.NetworkInterfaceID; got != nicId {
// 		return nil, errors.Errorf("got interface data for %s, expect %s", got, nicId)
// 	}
// 	for _, privateIP := range networkInterface.PrivateIPSets.PrivateIPSet {
// 		if !privateIP.Primary {
// 			ipAddrs = append(ipAddrs, privateIP.PrivateIPAddress)
// 		}
// 	}
// 	return ipAddrs, nil
// }

// func (nic *SInstanceNic) ipAddrsParams(ipAddrs []string) map[string]string {
// 	region := nic.instance.host.zone.region
// 	params := map[string]string{
// 		"RegionId":           region.GetId(),
// 		"NetworkInterfaceId": nic.mustGetId(),
// 	}
// 	for i, ipAddr := range ipAddrs {
// 		k := fmt.Sprintf("PrivateIpAddress.%d", i+1)
// 		params[k] = ipAddr
// 	}
// 	return params
// }

func (nic *SInstanceNic) AssignAddress(ipAddrs []string) error {
	return cloudprovider.ErrNotImplemented
}

func (nic *SInstanceNic) UnassignAddress(ipAddrs []string) error {
	return cloudprovider.ErrNotImplemented
}
