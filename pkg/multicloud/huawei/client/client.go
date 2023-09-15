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

package client

import (
	"net/http"

	"yunion.io/x/cloudmux/pkg/multicloud/huawei/client/auth"
	"yunion.io/x/cloudmux/pkg/multicloud/huawei/client/auth/credentials"
	"yunion.io/x/cloudmux/pkg/multicloud/huawei/client/modules"
)

type Client struct {
	cfg *SClientConfig
	// 标记初始化状态
	init bool

	Interface            *modules.SInterfaceManager
	Jobs                 *modules.SJobManager
	Orders               *modules.SOrderManager
	Servers              *modules.SServerManager
	ServersV2            *modules.SServerManager
	NovaServers          *modules.SServerManager
	VpcRoutes            *modules.SVpcRouteManager
	VpcPeerings          *modules.SVpcPeeringManager
	DBInstance           *modules.SDBInstanceManager
	DBInstanceBackup     *modules.SDBInstanceBackupManager
	DBInstanceFlavor     *modules.SDBInstanceFlavorManager
	DBInstanceJob        *modules.SDBInstanceJobManager
	CloudEye             *modules.SCloudEyeManager
	Roles                *modules.SRoleManager
	Groups               *modules.SGroupManager
	SAMLProviders        *modules.SAMLProviderManager
	SAMLProviderMappings *modules.SAMLProviderMappingManager
}

type SClientConfig struct {
	signer    auth.Signer
	endpoint  string // myhuaweicloud.com
	regionId  string
	domainId  string
	projectId string

	debug bool
}

func (self *SClientConfig) GetSigner() auth.Signer {
	return self.signer
}

func (self *SClientConfig) GetEndpoint() string {
	return self.endpoint
}

func (self *SClientConfig) GetRegionId() string {
	return self.regionId
}

func (self *SClientConfig) GetDomainId() string {
	return self.domainId
}

func (self *SClientConfig) GetProjectId() string {
	return self.projectId
}

func (self *SClientConfig) GetDebug() bool {
	return self.debug
}

func (self *Client) SetHttpClient(httpClient *http.Client) {
	self.Servers.SetHttpClient(httpClient)
	self.ServersV2.SetHttpClient(httpClient)
	self.NovaServers.SetHttpClient(httpClient)
	self.Orders.SetHttpClient(httpClient)
	self.Interface.SetHttpClient(httpClient)
	self.Jobs.SetHttpClient(httpClient)
	self.VpcRoutes.SetHttpClient(httpClient)
	self.DBInstance.SetHttpClient(httpClient)
	self.DBInstanceBackup.SetHttpClient(httpClient)
	self.DBInstanceFlavor.SetHttpClient(httpClient)
	self.DBInstanceJob.SetHttpClient(httpClient)
	self.CloudEye.SetHttpClient(httpClient)
	self.Roles.SetHttpClient(httpClient)
	self.Groups.SetHttpClient(httpClient)
	self.SAMLProviders.SetHttpClient(httpClient)
	self.SAMLProviderMappings.SetHttpClient(httpClient)
}

func (self *Client) InitWithAccessKey(endpoint, regionId, domainId, projectId, accessKey, secretKey string, debug bool) error {
	// accessKey signer
	credential := &credentials.AccessKeyCredential{
		AccessKeyId:     accessKey,
		AccessKeySecret: secretKey,
	}

	// 从signer中初始化
	signer, err := auth.NewSignerWithCredential(credential)
	if err != nil {
		return err
	}
	self.cfg = &SClientConfig{
		signer:    signer,
		endpoint:  endpoint,
		regionId:  regionId,
		domainId:  domainId,
		projectId: projectId,
		debug:     debug,
	}

	// 初始化 resource manager
	self.initManagers()
	return err
}

func (self *Client) initManagers() {
	if !self.init {
		self.Servers = modules.NewServerManager(self.cfg)
		self.ServersV2 = modules.NewServerV2Manager(self.cfg)
		self.NovaServers = modules.NewNovaServerManager(self.cfg)
		self.Orders = modules.NewOrderManager(self.cfg)
		self.Interface = modules.NewInterfaceManager(self.cfg)
		self.Jobs = modules.NewJobManager(self.cfg)
		self.VpcRoutes = modules.NewVpcRouteManager(self.cfg)
		self.VpcPeerings = modules.NewVpcPeeringManager(self.cfg)
		self.DBInstance = modules.NewDBInstanceManager(self.cfg)
		self.DBInstanceBackup = modules.NewDBInstanceBackupManager(self.cfg)
		self.DBInstanceFlavor = modules.NewDBInstanceFlavorManager(self.cfg)
		self.DBInstanceJob = modules.NewDBInstanceJobManager(self.cfg)
		self.CloudEye = modules.NewCloudEyeManager(self.cfg)
		self.Roles = modules.NewRoleManager(self.cfg)
		self.Groups = modules.NewGroupManager(self.cfg)
		self.SAMLProviders = modules.NewSAMLProviderManager(self.cfg)
		self.SAMLProviderMappings = modules.NewSAMLProviderMappingManager(self.cfg)
	}

	self.init = true
}

func NewClientWithAccessKey(endpoint, regionId, domainId, projectId, accessKey, secretKey string, debug bool) (*Client, error) {
	c := &Client{}
	err := c.InitWithAccessKey(endpoint, regionId, domainId, projectId, accessKey, secretKey, debug)
	return c, err
}

func NewPublicCloudClientWithAccessKey(regionId, domainId, projectId, accessKey, secretKey string, debug bool) (*Client, error) {
	return NewClientWithAccessKey("myhuaweicloud.com", regionId, domainId, projectId, accessKey, secretKey, debug)
}
