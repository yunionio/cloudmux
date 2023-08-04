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
	"yunion.io/x/cloudmux/pkg/cloudprovider"
)

type KsyunClientConfig struct {
	cpcfg           cloudprovider.ProviderConfig
	accessKeyId     string
	accessKeySecret string

	debug bool
}

type SKsyunClient struct {
	*KsyunClientConfig
}

func NewKsyunClientConfig(accessKeyId, accessKeySecret string) *KsyunClientConfig {
	cfg := &KsyunClientConfig{
		accessKeyId:     accessKeyId,
		accessKeySecret: accessKeySecret,
	}
	return cfg
}

func (self *KsyunClientConfig) Debug(debug bool) *KsyunClientConfig {
	self.debug = debug
	return self
}

func (self *KsyunClientConfig) CloudproviderConfig(cpcfg cloudprovider.ProviderConfig) *KsyunClientConfig {
	self.cpcfg = cpcfg
	return self
}

func NewKsyunClient(cfg *KsyunClientConfig) (*SKsyunClient, error) {
	client := &SKsyunClient{
		KsyunClientConfig: cfg,
	}
	return client, cloudprovider.ErrNotImplemented
}

func (self *SKsyunClient) GetRegions() ([]SRegion, error) {
	return nil, cloudprovider.ErrNotImplemented
}

func (self *SKsyunClient) GetRegion(id string) (*SRegion, error) {
	return nil, cloudprovider.ErrNotImplemented
}
