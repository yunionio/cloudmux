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

package cloudprovider

import "yunion.io/x/pkg/util/billing"

type LoadbalancerType string

const (
	// 传统型负载均衡
	LoadbalancerTypeSLB LoadbalancerType = "slb"
	// 应用型负载均衡
	LoadbalancerTypeALB LoadbalancerType = "alb"
	// 网络型负载均衡
	LoadbalancerTypeNLB LoadbalancerType = "nlb"
	// 网关型负载均衡
	LoadbalancerTypeGWLB LoadbalancerType = "gwlb"
)

type SLoadbalancerCreateOptions struct {
	Name             string
	Desc             string
	ZoneId           string
	SlaveZoneId      string
	VpcId            string
	NetworkIds       []string
	EipId            string // eip id
	Address          string
	AddressType      string
	LoadbalancerSpec string
	ChargeType       string
	EgressMbps       int
	BillingCycle     *billing.SBillingCycle
	ProjectId        string
	Tags             map[string]string
}
