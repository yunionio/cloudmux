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
	"context"
	"yunion.io/x/cloudmux/pkg/cloudprovider"
)

func (i *SLoadbalancer) GetIEIPs() ([]cloudprovider.ICloudEIP, error) {
	return nil, cloudprovider.ErrNotSupported
}
func (i *SLoadbalancer) Delete(ctx context.Context) error { return cloudprovider.ErrNotSupported }
func (i *SLoadbalancer) Start() error                     { return cloudprovider.ErrNotSupported }
func (i *SLoadbalancer) Stop() error                      { return cloudprovider.ErrNotSupported }
func (i *SLoadbalancer) CreateILoadBalancerBackendGroup(group *cloudprovider.SLoadbalancerBackendGroup) (cloudprovider.ICloudLoadbalancerBackendGroup, error) {
	return nil, cloudprovider.ErrNotSupported
}
func (i *SLoadbalancer) CreateILoadBalancerListener(ctx context.Context, listener *cloudprovider.SLoadbalancerListenerCreateOptions) (cloudprovider.ICloudLoadbalancerListener, error) {
	return nil, cloudprovider.ErrNotSupported
}
func (i *SLoadbalancer) GetSecurityGroupIds() ([]string, error) {
	return nil, cloudprovider.ErrNotSupported
}
func (i *SLoadbalancerListener) CreateILoadBalancerListenerRule(rule *cloudprovider.SLoadbalancerListenerRule) (cloudprovider.ICloudLoadbalancerListenerRule, error) {
	return nil, cloudprovider.ErrNotSupported
}
func (i *SLoadbalancerListener) GetILoadBalancerListenerRuleById(ruleId string) (cloudprovider.ICloudLoadbalancerListenerRule, error) {
	return nil, cloudprovider.ErrNotSupported
}
func (i *SLoadbalancerListener) GetILoadbalancerListenerRules() ([]cloudprovider.ICloudLoadbalancerListenerRule, error) {
	if i.Protocol == "TCP" || i.Protocol == "UDP" {
		return []cloudprovider.ICloudLoadbalancerListenerRule{}, nil
	}
	return nil, cloudprovider.ErrNotSupported
}
func (i *SLoadbalancerListener) Start() error { return cloudprovider.ErrNotSupported }
func (i *SLoadbalancerListener) Stop() error  { return cloudprovider.ErrNotSupported }
func (i *SLoadbalancerListener) ChangeScheduler(ctx context.Context, opts *cloudprovider.ChangeListenerSchedulerOptions) error {
	return cloudprovider.ErrNotSupported
}
func (i *SLoadbalancerListener) SetHealthCheck(ctx context.Context, opts *cloudprovider.ListenerHealthCheckOptions) error {
	return cloudprovider.ErrNotSupported
}
func (i *SLoadbalancerListener) ChangeCertificate(ctx context.Context, opts *cloudprovider.ListenerCertificateOptions) error {
	return cloudprovider.ErrNotSupported
}
func (i *SLoadbalancerListener) SetAcl(ctx context.Context, opts *cloudprovider.ListenerAclOptions) error {
	return cloudprovider.ErrNotSupported
}
func (i *SLoadbalancerListener) Delete(ctx context.Context) error {
	return cloudprovider.ErrNotSupported
}
func (i *SLoadbalancerBackendGroup) AddBackendServer(opts *cloudprovider.SLoadbalancerBackend) (cloudprovider.ICloudLoadbalancerBackend, error) {
	return nil, cloudprovider.ErrNotSupported
}
func (i *SLoadbalancerBackendGroup) RemoveBackendServer(opts *cloudprovider.SLoadbalancerBackend) error {
	return cloudprovider.ErrNotSupported
}
func (i *SLoadbalancerBackendGroup) Delete(ctx context.Context) error {
	return cloudprovider.ErrNotSupported
}
func (i *SLoadbalancerBackendGroup) Update(ctx context.Context, opts *cloudprovider.SLoadbalancerBackendGroup) error {
	return cloudprovider.ErrNotSupported
}
func (i *SLoadbalancerBackend) Update(ctx context.Context, opts *cloudprovider.SLoadbalancerBackend) error {
	return cloudprovider.ErrNotSupported
}
