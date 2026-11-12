// Copyright 2026 Yunion
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package volcengine

import (
	"strings"
	"time"
	"yunion.io/x/cloudmux/pkg/cloudprovider"
	"yunion.io/x/cloudmux/pkg/multicloud"
	"yunion.io/x/pkg/errors"
)

type SLoadbalancer struct {
	multicloud.SLoadbalancerBase
	VolcEngineTags
	region                                                                                                  *SRegion
	LoadBalancerId, LoadBalancerName, Status, Type, VpcId, SubnetId, MasterZoneId, SlaveZoneId, ProjectName string
	EipAddress, EniAddress, LoadBalancerSpec, CreateTime                                                    string
	LoadBalancerBillingType                                                                                 int
	ServerGroups                                                                                            []struct{ ServerGroupId string }
}

var _ cloudprovider.ICloudLoadbalancer = (*SLoadbalancer)(nil)

func (r *SRegion) GetILoadBalancers() ([]cloudprovider.ICloudLoadbalancer, error) {
	items, e := r.inventoryPages("clb", "DescribeLoadBalancers", "LoadBalancers", "TotalCount", "LoadBalancerId", nil)
	if e != nil {
		return nil, e
	}
	ret := make([]cloudprovider.ICloudLoadbalancer, 0, len(items))
	for _, item := range items {
		id, _ := item.GetString("LoadBalancerId")
		lb, e := r.GetLoadbalancer(id)
		if e != nil {
			return nil, e
		}
		ret = append(ret, lb)
	}
	return ret, nil
}
func (r *SRegion) GetILoadBalancerById(id string) (cloudprovider.ICloudLoadbalancer, error) {
	return r.GetLoadbalancer(id)
}
func (r *SRegion) GetLoadbalancer(id string) (*SLoadbalancer, error) {
	if id == "" {
		return nil, cloudprovider.ErrNotFound
	}
	data, e := r.client.inventoryRequest(r.RegionId, "clb", "DescribeLoadBalancerAttributes", map[string]interface{}{"LoadBalancerId": id})
	if e != nil {
		return nil, e
	}
	if _, e := data.GetArray("ServerGroups"); e != nil {
		return nil, errors.Wrap(e, "missing CLB server group references")
	}
	lb := &SLoadbalancer{region: r}
	if e := data.Unmarshal(lb); e != nil {
		return nil, e
	}
	if lb.LoadBalancerId != id {
		return nil, errors.Errorf("CLB detail ID mismatch")
	}
	return lb, nil
}
func (l *SLoadbalancer) GetId() string        { return l.LoadBalancerId }
func (l *SLoadbalancer) GetGlobalId() string  { return l.LoadBalancerId }
func (l *SLoadbalancer) GetName() string      { return l.LoadBalancerName }
func (l *SLoadbalancer) GetProjectId() string { return l.ProjectName }
func (l *SLoadbalancer) GetStatus() string {
	switch l.Status {
	case "Active":
		return "enabled"
	case "Stopped", "Inactive":
		return "disabled"
	default:
		return "unknown"
	}
}
func (l *SLoadbalancer) GetAddress() string {
	if l.Type == "public" {
		return l.EipAddress
	}
	return l.EniAddress
}
func (l *SLoadbalancer) GetAddressType() string {
	if l.Type == "public" {
		return "internet"
	}
	return "intranet"
}
func (l *SLoadbalancer) GetNetworkType() string { return "vpc" }
func (l *SLoadbalancer) GetNetworkIds() []string {
	if l.SubnetId == "" {
		return nil
	}
	return []string{l.SubnetId}
}
func (l *SLoadbalancer) GetVpcId() string            { return l.VpcId }
func (l *SLoadbalancer) GetZoneId() string           { return l.MasterZoneId }
func (l *SLoadbalancer) GetZone1Id() string          { return l.SlaveZoneId }
func (l *SLoadbalancer) GetLoadbalancerSpec() string { return l.LoadBalancerSpec }
func (l *SLoadbalancer) GetChargeType() string       { return "" }
func (l *SLoadbalancer) GetEgressMbps() int          { return 0 }
func (l *SLoadbalancer) GetCreatedAt() time.Time {
	t, _ := time.Parse(time.RFC3339Nano, l.CreateTime)
	return t
}
func (l *SLoadbalancer) Refresh() error {
	lb, e := l.region.GetLoadbalancer(l.LoadBalancerId)
	if e != nil {
		return e
	}
	*l = *lb
	return nil
}
func (l *SLoadbalancer) GetILoadBalancerListeners() ([]cloudprovider.ICloudLoadbalancerListener, error) {
	items, e := l.region.inventoryPages("clb", "DescribeListeners", "Listeners", "TotalCount", "ListenerId", map[string]interface{}{"LoadBalancerId": l.LoadBalancerId})
	if e != nil {
		return nil, e
	}
	ret := make([]cloudprovider.ICloudLoadbalancerListener, 0, len(items))
	for _, item := range items {
		x := &SLoadbalancerListener{lb: l}
		if e := item.Unmarshal(x); e != nil {
			return nil, e
		}
		if x.LoadBalancerId != l.LoadBalancerId {
			return nil, errors.Errorf("CLB listener parent mismatch")
		}
		ret = append(ret, x)
	}
	return ret, nil
}
func (l *SLoadbalancer) GetILoadBalancerListenerById(id string) (cloudprovider.ICloudLoadbalancerListener, error) {
	items, e := l.GetILoadBalancerListeners()
	if e != nil {
		return nil, e
	}
	for _, i := range items {
		if i.GetGlobalId() == id {
			return i, nil
		}
	}
	return nil, cloudprovider.ErrNotFound
}
func (l *SLoadbalancer) GetILoadBalancerBackendGroups() ([]cloudprovider.ICloudLoadbalancerBackendGroup, error) {
	// Refresh the authoritative parent before reading its server group references.
	lb, e := l.region.GetLoadbalancer(l.LoadBalancerId)
	if e != nil {
		return nil, e
	}
	ret := make([]cloudprovider.ICloudLoadbalancerBackendGroup, 0, len(lb.ServerGroups))
	seen := map[string]bool{}
	for _, g := range lb.ServerGroups {
		if g.ServerGroupId == "" || seen[g.ServerGroupId] {
			return nil, errors.Errorf("invalid CLB server group reference")
		}
		seen[g.ServerGroupId] = true
		x, e := l.GetILoadBalancerBackendGroupById(g.ServerGroupId)
		if e != nil {
			return nil, e
		}
		ret = append(ret, x)
	}
	return ret, nil
}
func (l *SLoadbalancer) GetILoadBalancerBackendGroupById(id string) (cloudprovider.ICloudLoadbalancerBackendGroup, error) {
	if id == "" {
		return nil, cloudprovider.ErrNotFound
	}
	data, e := l.region.client.inventoryRequest(l.region.RegionId, "clb", "DescribeServerGroupAttributes", map[string]interface{}{"ServerGroupId": id})
	if e != nil {
		return nil, e
	}
	g := &SLoadbalancerBackendGroup{}
	if e := data.Unmarshal(g); e != nil {
		return nil, e
	}
	if g.ServerGroupId != id {
		return nil, errors.Errorf("CLB server group identity mismatch")
	}
	if _, e := data.GetArray("Servers"); e != nil {
		return nil, errors.Wrap(e, "missing CLB backend servers")
	}
	return g, nil
}

type SLoadbalancerListener struct {
	multicloud.SResourceBase
	VolcEngineTags
	lb                                                                                            *SLoadbalancer
	ListenerId, ListenerName, LoadBalancerId, Protocol, ServerGroupId, Status, Enabled, Scheduler string
	AclStatus, AclType, CertificateId, SecurityPolicyId, Http2Enabled, Cookie, PersistenceType    string
	AclIds                                                                                        []string
	Port, Bandwidth, EstablishedTimeout, ProxyConnectTimeout, PersistenceTimeout                  int
	HealthCheck                                                                                   struct {
		Enabled, Domain, Uri, Method, HttpCode, UdpRequest, UdpExpect string
		Port, Interval, Timeout, HealthyThreshold, UnHealthyThreshold int
	}
}

var _ cloudprovider.ICloudLoadbalancerListener = (*SLoadbalancerListener)(nil)

func (l *SLoadbalancerListener) GetId() string       { return l.ListenerId }
func (l *SLoadbalancerListener) GetGlobalId() string { return l.ListenerId }
func (l *SLoadbalancerListener) GetName() string     { return l.ListenerName }
func (l *SLoadbalancerListener) GetStatus() string {
	if l.Enabled == "on" && l.Status == "Active" {
		return "enabled"
	}
	if l.Enabled == "off" {
		return "disabled"
	}
	return "unknown"
}
func (l *SLoadbalancerListener) GetListenerType() string { return strings.ToLower(l.Protocol) }
func (l *SLoadbalancerListener) GetListenerPort() int    { return l.Port }
func (l *SLoadbalancerListener) GetScheduler() string    { return l.Scheduler }
func (l *SLoadbalancerListener) GetAclStatus() string    { return l.AclStatus }
func (l *SLoadbalancerListener) GetAclType() string      { return l.AclType }
func (l *SLoadbalancerListener) GetAclId() string {
	if len(l.AclIds) == 1 {
		return l.AclIds[0]
	}
	return ""
}
func (l *SLoadbalancerListener) GetEgressMbps() int            { return l.Bandwidth }
func (l *SLoadbalancerListener) GetBackendGroupId() string     { return l.ServerGroupId }
func (l *SLoadbalancerListener) GetBackendServerPort() int     { return 0 }
func (l *SLoadbalancerListener) GetClientIdleTimeout() int     { return l.EstablishedTimeout }
func (l *SLoadbalancerListener) GetBackendConnectTimeout() int { return l.ProxyConnectTimeout }
func (l *SLoadbalancerListener) GetStickySession() string {
	if l.PersistenceType != "" && l.PersistenceType != "none" {
		return "on"
	}
	return "off"
}
func (l *SLoadbalancerListener) GetStickySessionType() string       { return l.PersistenceType }
func (l *SLoadbalancerListener) GetStickySessionCookie() string     { return l.Cookie }
func (l *SLoadbalancerListener) GetStickySessionCookieTimeout() int { return l.PersistenceTimeout }
func (l *SLoadbalancerListener) XForwardedForEnabled() bool         { return false }
func (l *SLoadbalancerListener) GzipEnabled() bool                  { return false }
func (l *SLoadbalancerListener) GetCertificateId() string           { return l.CertificateId }
func (l *SLoadbalancerListener) GetTLSCipherPolicy() string         { return l.SecurityPolicyId }
func (l *SLoadbalancerListener) HTTP2Enabled() bool                 { return l.Http2Enabled == "on" }
func (l *SLoadbalancerListener) GetRedirect() string                { return "" }
func (l *SLoadbalancerListener) GetRedirectCode() int64             { return 0 }
func (l *SLoadbalancerListener) GetRedirectScheme() string          { return "" }
func (l *SLoadbalancerListener) GetRedirectHost() string            { return "" }
func (l *SLoadbalancerListener) GetRedirectPath() string            { return "" }
func (l *SLoadbalancerListener) GetHealthCheck() string             { return l.HealthCheck.Enabled }
func (l *SLoadbalancerListener) GetHealthCheckType() string         { return strings.ToLower(l.Protocol) }
func (l *SLoadbalancerListener) GetHealthCheckTimeout() int         { return l.HealthCheck.Timeout }
func (l *SLoadbalancerListener) GetHealthCheckInterval() int        { return l.HealthCheck.Interval }
func (l *SLoadbalancerListener) GetHealthCheckRise() int            { return l.HealthCheck.HealthyThreshold }
func (l *SLoadbalancerListener) GetHealthCheckFail() int            { return l.HealthCheck.UnHealthyThreshold }
func (l *SLoadbalancerListener) GetHealthCheckReq() string          { return l.HealthCheck.UdpRequest }
func (l *SLoadbalancerListener) GetHealthCheckExp() string          { return l.HealthCheck.UdpExpect }
func (l *SLoadbalancerListener) GetHealthCheckDomain() string       { return l.HealthCheck.Domain }
func (l *SLoadbalancerListener) GetHealthCheckURI() string          { return l.HealthCheck.Uri }
func (l *SLoadbalancerListener) GetHealthCheckMethod() string       { return l.HealthCheck.Method }
func (l *SLoadbalancerListener) GetHealthCheckPort() int            { return l.HealthCheck.Port }
func (l *SLoadbalancerListener) GetHealthCheckCode() string         { return l.HealthCheck.HttpCode }

type SLoadbalancerBackendGroup struct {
	multicloud.SResourceBase
	VolcEngineTags
	ServerGroupId, ServerGroupName, Type string
	Servers                              []SLoadbalancerBackend
}

var _ cloudprovider.ICloudLoadbalancerBackendGroup = (*SLoadbalancerBackendGroup)(nil)

func (g *SLoadbalancerBackendGroup) GetId() string            { return g.ServerGroupId }
func (g *SLoadbalancerBackendGroup) GetGlobalId() string      { return g.ServerGroupId }
func (g *SLoadbalancerBackendGroup) GetName() string          { return g.ServerGroupName }
func (g *SLoadbalancerBackendGroup) GetStatus() string        { return "unknown" }
func (g *SLoadbalancerBackendGroup) IsDefault() bool          { return false }
func (g *SLoadbalancerBackendGroup) GetType() string          { return "normal" }
func (g *SLoadbalancerBackendGroup) GetScheduler() string     { return "" }
func (g *SLoadbalancerBackendGroup) GetHealthCheckId() string { return "" }
func (g *SLoadbalancerBackendGroup) GetILoadbalancerBackends() ([]cloudprovider.ICloudLoadbalancerBackend, error) {
	ret := make([]cloudprovider.ICloudLoadbalancerBackend, 0, len(g.Servers))
	seen := map[string]bool{}
	for idx := range g.Servers {
		x := &g.Servers[idx]
		if x.ServerId == "" || seen[x.ServerId] {
			return nil, errors.Errorf("invalid CLB backend identity")
		}
		seen[x.ServerId] = true
		ret = append(ret, x)
	}
	return ret, nil
}
func (g *SLoadbalancerBackendGroup) GetILoadbalancerBackendById(id string) (cloudprovider.ICloudLoadbalancerBackend, error) {
	items, e := g.GetILoadbalancerBackends()
	if e != nil {
		return nil, e
	}
	for _, x := range items {
		if x.GetGlobalId() == id {
			return x, nil
		}
	}
	return nil, cloudprovider.ErrNotFound
}

type SLoadbalancerBackend struct {
	multicloud.SResourceBase
	ServerId, InstanceId, Ip, Type, Description string
	Port, Weight                                int
}

var _ cloudprovider.ICloudLoadbalancerBackend = (*SLoadbalancerBackend)(nil)

func (s *SLoadbalancerBackend) GetId() string          { return s.ServerId }
func (s *SLoadbalancerBackend) GetGlobalId() string    { return s.ServerId }
func (s *SLoadbalancerBackend) GetName() string        { return s.ServerId }
func (s *SLoadbalancerBackend) GetDescription() string { return s.Description }
func (s *SLoadbalancerBackend) GetStatus() string      { return "unknown" }
func (s *SLoadbalancerBackend) GetWeight() int         { return s.Weight }
func (s *SLoadbalancerBackend) GetPort() int           { return s.Port }
func (s *SLoadbalancerBackend) GetBackendType() string {
	if s.Type == "ecs" {
		return "guest"
	}
	if s.Type == "ip" {
		return "ip"
	}
	return "unknown"
}
func (s *SLoadbalancerBackend) GetBackendRole() string { return "default" }
func (s *SLoadbalancerBackend) GetBackendId() string {
	if s.Type == "ecs" {
		return s.InstanceId
	}
	return s.Ip
}
func (s *SLoadbalancerBackend) GetIpAddress() string { return s.Ip }
