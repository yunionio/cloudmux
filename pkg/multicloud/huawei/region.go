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

package huawei

import (
	"fmt"
	"net/url"
	"strings"
	"time"

	"yunion.io/x/jsonutils"
	"yunion.io/x/log"
	"yunion.io/x/pkg/errors"
	"yunion.io/x/pkg/util/secrules"

	api "yunion.io/x/cloudmux/pkg/apis/compute"
	"yunion.io/x/cloudmux/pkg/cloudprovider"
	"yunion.io/x/cloudmux/pkg/multicloud"
	"yunion.io/x/cloudmux/pkg/multicloud/huawei/client"
	"yunion.io/x/cloudmux/pkg/multicloud/huawei/obs"
)

type Locales struct {
	EnUs string `json:"en-us"`
	ZhCN string `json:"zh-cn"`
}

// https://support.huaweicloud.com/api-iam/zh-cn_topic_0067148043.html
type SRegion struct {
	multicloud.SRegion

	client    *SHuaweiClient
	ecsClient *client.Client
	obsClient *obs.ObsClient // 对象存储client.请勿直接引用。

	Description    string  `json:"description"`
	ID             string  `json:"id"`
	Locales        Locales `json:"locales"`
	ParentRegionID string  `json:"parent_region_id"`
	Type           string  `json:"type"`

	zones []SZone
	ivpcs []cloudprovider.ICloudVpc

	storageCache *SStoragecache
}

func (self *SRegion) GetClient() *SHuaweiClient {
	return self.client
}

func (self *SRegion) getECSClient() (*client.Client, error) {
	var err error

	if len(self.client.projectId) > 0 {
		project, err := self.client.GetProjectById(self.client.projectId)
		if err != nil {
			return nil, err
		}

		if !strings.Contains(project.Name, self.ID) {
			return nil, errors.Errorf("region %s and project %s mismatch", self.ID, project.Name)
		}
	}

	if self.ecsClient == nil {
		self.ecsClient, err = self.client.newRegionAPIClient(self.ID)
		if err != nil {
			return nil, err
		}
	}

	return self.ecsClient, err
}

func (self *SRegion) getOBSEndpoint() string {
	return getOBSEndpoint(self.GetId())
}

func (self *SRegion) getOBSClient() (*obs.ObsClient, error) {
	if self.obsClient == nil {
		obsClient, err := self.client.getOBSClient(self.GetId())
		if err != nil {
			return nil, err
		}

		self.obsClient = obsClient
	}

	return self.obsClient, nil
}

func (self *SRegion) list(service, resource string, query url.Values) (jsonutils.JSONObject, error) {
	return self.client.list(service, self.ID, resource, query)
}

func (self *SRegion) delete(service, resource string) (jsonutils.JSONObject, error) {
	return self.client.delete(service, self.ID, resource)
}

func (self *SRegion) put(service, resource string, params map[string]interface{}) (jsonutils.JSONObject, error) {
	return self.client.put(service, self.ID, resource, params)
}

func (self *SRegion) post(service, resource string, params map[string]interface{}) (jsonutils.JSONObject, error) {
	return self.client.post(service, self.ID, resource, params)
}

func (self *SRegion) GetZones() ([]SZone, error) {
	if len(self.zones) > 0 {
		return self.zones, nil
	}
	resp, err := self.list(SERVICE_ECS, "os-availability-zone", nil)
	if err != nil {
		return nil, errors.Wrapf(err, "list")
	}
	self.zones = []SZone{}
	err = resp.Unmarshal(&self.zones, "availabilityZoneInfo")
	if err != nil {
		return nil, errors.Wrapf(err, "Unmarshal")
	}

	for i := range self.zones {
		self.zones[i].region = self
	}
	return self.zones, nil
}

func (self *SRegion) GetIVMById(id string) (cloudprovider.ICloudVM, error) {
	if len(id) == 0 {
		return nil, errors.Wrap(cloudprovider.ErrNotFound, "SRegion.GetIVMById")
	}

	instance, err := self.GetInstanceByID(id)
	if err != nil {
		return nil, err
	}
	return &instance, err
}

func (self *SRegion) GetIDiskById(id string) (cloudprovider.ICloudDisk, error) {
	return self.GetDisk(id)
}

func (self *SRegion) GetGeographicInfo() cloudprovider.SGeographicInfo {
	if info, ok := LatitudeAndLongitude[self.ID]; ok {
		return info
	}
	return cloudprovider.SGeographicInfo{}
}

func (self *SRegion) GetILoadBalancers() ([]cloudprovider.ICloudLoadbalancer, error) {
	elbs, err := self.GetLoadBalancers()
	if err != nil {
		return nil, err
	}

	ielbs := make([]cloudprovider.ICloudLoadbalancer, len(elbs))
	for i := range elbs {
		elbs[i].region = self
		ielbs[i] = &elbs[i]
	}

	return ielbs, nil
}

func (self *SRegion) GetLoadBalancers() ([]SLoadbalancer, error) {
	lbs := []SLoadbalancer{}
	params := url.Values{}
	return lbs, self.lbListAll("elb/loadbalancers", params, "loadbalancers", &lbs)
}

func (self *SRegion) GetILoadBalancerById(id string) (cloudprovider.ICloudLoadbalancer, error) {
	elb, err := self.GetLoadbalancer(id)
	if err != nil {
		return nil, err
	}
	return elb, nil
}

func (self *SRegion) GetILoadBalancerAclById(aclId string) (cloudprovider.ICloudLoadbalancerAcl, error) {
	acl, err := self.GetLoadBalancerAcl(aclId)
	if err != nil {
		return nil, err
	}
	return acl, nil
}

func (self *SRegion) GetILoadBalancerCertificateById(certId string) (cloudprovider.ICloudLoadbalancerCertificate, error) {
	cert, err := self.GetLoadBalancerCertificate(certId)
	if err != nil {
		return nil, err
	}
	return cert, nil
}

func (self *SRegion) CreateILoadBalancerCertificate(cert *cloudprovider.SLoadbalancerCertificate) (cloudprovider.ICloudLoadbalancerCertificate, error) {
	ret, err := self.CreateLoadBalancerCertificate(cert)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

func (self *SRegion) GetILoadBalancerAcls() ([]cloudprovider.ICloudLoadbalancerAcl, error) {
	ret, err := self.GetLoadBalancerAcls("")
	if err != nil {
		return nil, err
	}

	iret := make([]cloudprovider.ICloudLoadbalancerAcl, len(ret))
	for i := range ret {
		ret[i].region = self
		iret[i] = &ret[i]
	}
	return iret, nil
}

func (self *SRegion) GetILoadBalancerCertificates() ([]cloudprovider.ICloudLoadbalancerCertificate, error) {
	ret, err := self.GetLoadBalancerCertificates()
	if err != nil {
		return nil, err
	}

	iret := make([]cloudprovider.ICloudLoadbalancerCertificate, len(ret))
	for i := range ret {
		ret[i].region = self
		iret[i] = &ret[i]
	}
	return iret, nil
}

// https://support.huaweicloud.com/api-iam/zh-cn_topic_0057845622.html
func (self *SRegion) GetId() string {
	return self.ID
}

func (self *SRegion) GetName() string {
	return fmt.Sprintf("%s %s", CLOUD_PROVIDER_HUAWEI_CN, self.Locales.ZhCN)
}

func (self *SRegion) GetI18n() cloudprovider.SModelI18nTable {
	en := fmt.Sprintf("%s %s", CLOUD_PROVIDER_HUAWEI_EN, self.Locales.EnUs)
	table := cloudprovider.SModelI18nTable{}
	table["name"] = cloudprovider.NewSModelI18nEntry(self.GetName()).CN(self.GetName()).EN(en)
	return table
}

func (self *SRegion) GetGlobalId() string {
	return fmt.Sprintf("%s/%s", self.client.GetAccessEnv(), self.ID)
}

func (self *SRegion) GetStatus() string {
	return api.CLOUD_REGION_STATUS_INSERVER
}

func (self *SRegion) Refresh() error {
	return nil
}

func (self *SRegion) IsEmulated() bool {
	return false
}

func (self *SRegion) GetLatitude() float32 {
	if locationInfo, ok := LatitudeAndLongitude[self.ID]; ok {
		return locationInfo.Latitude
	}
	return 0.0
}

func (self *SRegion) GetLongitude() float32 {
	if locationInfo, ok := LatitudeAndLongitude[self.ID]; ok {
		return locationInfo.Longitude
	}
	return 0.0
}

func (self *SRegion) GetIZones() ([]cloudprovider.ICloudZone, error) {
	zones, err := self.GetZones()
	if err != nil {
		return nil, err
	}
	ret := []cloudprovider.ICloudZone{}
	for i := range zones {
		ret = append(ret, &zones[i])
	}
	return ret, nil
}

func (self *SRegion) GetIVpcs() ([]cloudprovider.ICloudVpc, error) {
	vpcs, err := self.GetVpcs()
	if err != nil {
		return nil, err
	}
	ret := []cloudprovider.ICloudVpc{}
	for i := range vpcs {
		vpcs[i].region = self
		ret = append(ret, &vpcs[i])
	}
	return ret, nil
}

func (self *SRegion) GetIEips() ([]cloudprovider.ICloudEIP, error) {
	eips, err := self.GetEips("", nil)
	if err != nil {
		return nil, err
	}

	ret := []cloudprovider.ICloudEIP{}
	for i := 0; i < len(eips); i += 1 {
		eips[i].region = self
		ret = append(ret, &eips[i])
	}
	return ret, nil
}

func (self *SRegion) GetIVpcById(id string) (cloudprovider.ICloudVpc, error) {
	vpc, err := self.GetVpc(id)
	if err != nil {
		return nil, err
	}
	return vpc, nil
}

func (self *SRegion) GetIZoneById(id string) (cloudprovider.ICloudZone, error) {
	izones, err := self.GetIZones()
	if err != nil {
		return nil, err
	}
	for i := 0; i < len(izones); i += 1 {
		if izones[i].GetGlobalId() == id {
			return izones[i], nil
		}
	}
	return nil, cloudprovider.ErrNotFound
}

func (self *SRegion) GetIEipById(eipId string) (cloudprovider.ICloudEIP, error) {
	eip, err := self.GetEip(eipId)
	if err != nil {
		return nil, err
	}
	return eip, nil
}

// https://support.huaweicloud.com/api-vpc/zh-cn_topic_0060595555.html
func (self *SRegion) DeleteSecurityGroup(id string) error {
	_, err := self.list(SERVICE_VPC, "security-groups/"+id, nil)
	return err
}

func (self *SRegion) GetISecurityGroupById(secgroupId string) (cloudprovider.ICloudSecurityGroup, error) {
	return self.GetSecurityGroupDetails(secgroupId)
}

func (self *SRegion) GetISecurityGroupByName(opts *cloudprovider.SecurityGroupFilterOptions) (cloudprovider.ICloudSecurityGroup, error) {
	secgroups, err := self.GetSecurityGroups(opts.VpcId, opts.Name)
	if err != nil {
		return nil, err
	}
	if len(secgroups) == 0 {
		return nil, cloudprovider.ErrNotFound
	}
	if len(secgroups) > 1 {
		return nil, cloudprovider.ErrDuplicateId
	}
	secgroups[0].region = self
	return &secgroups[0], nil
}

func (self *SRegion) CreateISecurityGroup(opts *cloudprovider.SecurityGroupCreateInput) (cloudprovider.ICloudSecurityGroup, error) {
	return self.CreateSecurityGroup(opts)
}

// https://support.huaweicloud.com/api-vpc/zh-cn_topic_0020090608.html
func (self *SRegion) CreateIVpc(opts *cloudprovider.VpcCreateOptions) (cloudprovider.ICloudVpc, error) {
	return self.CreateVpc(opts.NAME, opts.CIDR, opts.Desc)
}

func (self *SRegion) CreateVpc(name, cidr, desc string) (*SVpc, error) {
	params := map[string]interface{}{
		"vpc": map[string]string{
			"name":        name,
			"cidr":        cidr,
			"description": desc,
		},
	}
	vpc := &SVpc{region: self}
	resp, err := self.post(SERVICE_VPC, "vpcs", params)
	if err != nil {
		return nil, err
	}
	return vpc, resp.Unmarshal(vpc, "vpc")
}

// https://support.huaweicloud.com/api-vpc/zh-cn_topic_0020090596.html
// size: 1Mbit/s~2000Mbit/s
// bgpType: 5_telcom，5_union，5_bgp，5_sbgp.
// 东北-大连：5_telcom、5_union
// 华南-广州：5_sbgp
// 华东-上海二：5_sbgp
// 华北-北京一：5_bgp、5_sbgp
// 亚太-香港：5_bgp
func (self *SRegion) CreateEIP(opts *cloudprovider.SEip) (cloudprovider.ICloudEIP, error) {
	eip, err := self.AllocateEIP(opts)
	if err != nil {
		return nil, err
	}
	err = cloudprovider.WaitStatus(eip, api.EIP_STATUS_READY, 5*time.Second, time.Minute)
	return eip, err
}

func (self *SRegion) GetISnapshots() ([]cloudprovider.ICloudSnapshot, error) {
	snapshots, err := self.GetSnapshots("", "")
	if err != nil {
		log.Errorf("self.GetSnapshots fail %s", err)
		return nil, err
	}

	ret := make([]cloudprovider.ICloudSnapshot, len(snapshots))
	for i := 0; i < len(snapshots); i += 1 {
		snapshots[i].region = self
		ret[i] = &snapshots[i]
	}
	return ret, nil
}

func (self *SRegion) GetISnapshotById(snapshotId string) (cloudprovider.ICloudSnapshot, error) {
	snapshot, err := self.GetSnapshot(snapshotId)
	if err != nil {
		return nil, err
	}
	return snapshot, nil
}

func (self *SRegion) GetIHosts() ([]cloudprovider.ICloudHost, error) {
	iHosts := make([]cloudprovider.ICloudHost, 0)

	izones, err := self.GetIZones()
	if err != nil {
		return nil, err
	}
	for i := 0; i < len(izones); i += 1 {
		iZoneHost, err := izones[i].GetIHosts()
		if err != nil {
			return nil, err
		}
		iHosts = append(iHosts, iZoneHost...)
	}
	return iHosts, nil
}

func (self *SRegion) GetIHostById(id string) (cloudprovider.ICloudHost, error) {
	izones, err := self.GetIZones()
	if err != nil {
		return nil, err
	}
	for i := 0; i < len(izones); i += 1 {
		ihost, err := izones[i].GetIHostById(id)
		if err == nil {
			return ihost, nil
		} else if errors.Cause(err) != cloudprovider.ErrNotFound {
			return nil, err
		}
	}
	return nil, cloudprovider.ErrNotFound
}

func (self *SRegion) GetIStorages() ([]cloudprovider.ICloudStorage, error) {
	iStores := make([]cloudprovider.ICloudStorage, 0)

	izones, err := self.GetIZones()
	if err != nil {
		return nil, err
	}
	for i := 0; i < len(izones); i += 1 {
		iZoneStores, err := izones[i].GetIStorages()
		if err != nil {
			return nil, err
		}
		iStores = append(iStores, iZoneStores...)
	}
	return iStores, nil
}

func (self *SRegion) GetIStorageById(id string) (cloudprovider.ICloudStorage, error) {
	izones, err := self.GetIZones()
	if err != nil {
		return nil, err
	}
	for i := 0; i < len(izones); i += 1 {
		istore, err := izones[i].GetIStorageById(id)
		if err == nil {
			return istore, nil
		} else if errors.Cause(err) != cloudprovider.ErrNotFound {
			return nil, err
		}
	}
	return nil, cloudprovider.ErrNotFound
}

func (self *SRegion) GetProvider() string {
	return CLOUD_PROVIDER_HUAWEI
}

func (self *SRegion) GetCloudEnv() string {
	return self.client.cloudEnv
}

// https://support.huaweicloud.com/api-vpc/zh-cn_topic_0020090615.html
// 目前desc字段并没有用到
func (self *SRegion) CreateSecurityGroup(opts *cloudprovider.SecurityGroupCreateInput) (*SSecurityGroup, error) {
	info := map[string]interface{}{
		"name": opts.Name,
	}
	if len(opts.VpcId) > 0 && opts.VpcId != api.NORMAL_VPC_ID {
		info["vpc_id"] = opts.VpcId
	}
	params := map[string]interface{}{
		"security_group": info,
	}

	resp, err := self.post(SERVICE_VPC, "security-groups", params)
	if err != nil {
		return nil, err
	}

	secgroup := &SSecurityGroup{region: self}
	err = resp.Unmarshal(secgroup, "security_group")
	if err != nil {
		return nil, err
	}
	if opts.OnCreated != nil {
		opts.OnCreated(secgroup.ID)
	}
	for _, rule := range secgroup.SecurityGroupRules {
		if len(rule.RemoteGroupID) > 0 || rule.Ethertype != "IPv4" {
			continue
		}
		err := self.delSecurityGroupRule(rule.ID)
		if err != nil {
			return nil, errors.Wrapf(err, "delete rule %s", rule.ID)
		}
	}
	rules := opts.InRules.AllowList()
	rules = append(rules, opts.OutRules.AllowList()...)
	for i := range rules {
		err := self.addSecurityGroupRules(secgroup.ID, rules[i])
		if err != nil {
			return nil, errors.Wrapf(err, "addSecurityGroupRules")
		}
	}
	return secgroup, nil
}

// https://support.huaweicloud.com/api-vpc/zh-cn_topic_0087467071.html
func (self *SRegion) delSecurityGroupRule(id string) error {
	resource := fmt.Sprintf("vpc/security-group-rules/%s", id)
	_, err := self.delete(SERVICE_VPC, resource)
	return err
}

func (self *SRegion) DeleteSecurityGroupRule(ruleId string) error {
	return self.delSecurityGroupRule(ruleId)
}

func (self *SRegion) CreateSecurityGroupRule(secgroupId string, rule secrules.SecurityRule) error {
	return self.addSecurityGroupRules(secgroupId, rule)
}

// https://support.huaweicloud.com/api-vpc/zh-cn_topic_0087451723.html
// icmp port对应关系：https://support.huaweicloud.com/api-vpc/zh-cn_topic_0024109590.html
func (self *SRegion) addSecurityGroupRules(secGrpId string, rule secrules.SecurityRule) error {
	direction := ""
	if rule.Direction == secrules.SecurityRuleIngress {
		direction = "ingress"
	} else {
		direction = "egress"
	}

	protocal := rule.Protocol
	if rule.Protocol == secrules.PROTO_ANY {
		protocal = ""
	}

	// imcp协议默认为any
	if rule.Protocol == secrules.PROTO_ICMP {
		return self.addSecurityGroupRule(secGrpId, direction, "-1", "-1", protocal, rule.IPNet.String())
	}

	if len(rule.Ports) > 0 {
		for _, port := range rule.Ports {
			portStr := fmt.Sprintf("%d", port)
			err := self.addSecurityGroupRule(secGrpId, direction, portStr, portStr, protocal, rule.IPNet.String())
			if err != nil {
				return err
			}
		}
	} else {
		portStart := fmt.Sprintf("%d", rule.PortStart)
		portEnd := fmt.Sprintf("%d", rule.PortEnd)
		err := self.addSecurityGroupRule(secGrpId, direction, portStart, portEnd, protocal, rule.IPNet.String())
		if err != nil {
			return err
		}
	}

	return nil
}

func (self *SRegion) addSecurityGroupRule(secGrpId, direction, portStart, portEnd, protocol, ipNet string) error {
	obj := map[string]interface{}{
		"security_group_id": secGrpId,
		"direction":         direction,
		"remote_ip_prefix":  ipNet,
		"ethertype":         "IPV4",
	}
	// 端口为空或者1-65535
	if len(portStart) > 0 && portStart != "0" && portStart != "-1" {
		obj["port_range_min"] = portStart
	}
	if len(portEnd) > 0 && portEnd != "0" && portEnd != "-1" {
		obj["port_range_max"] = portEnd
	}
	if len(protocol) > 0 {
		obj["protocol"] = protocol
	}

	params := map[string]interface{}{
		"security_group_rule": obj,
	}
	_, err := self.post(SERVICE_VPC, "vpc/security-group-rules", params)
	return err
}

func (self *SRegion) CreateILoadBalancer(loadbalancer *cloudprovider.SLoadbalancerCreateOptions) (cloudprovider.ICloudLoadbalancer, error) {
	ret, err := self.CreateLoadBalancer(loadbalancer)
	if err != nil {
		return nil, err
	}

	return ret, nil
}

func (self *SRegion) CreateILoadBalancerAcl(acl *cloudprovider.SLoadbalancerAccessControlList) (cloudprovider.ICloudLoadbalancerAcl, error) {
	ret, err := self.CreateLoadBalancerAcl(acl)
	if err != nil {
		return nil, err
	}

	return ret, nil
}

func (region *SRegion) GetIBuckets() ([]cloudprovider.ICloudBucket, error) {
	iBuckets, err := region.client.getIBuckets()
	if err != nil {
		return nil, errors.Wrap(err, "getIBuckets")
	}
	ret := make([]cloudprovider.ICloudBucket, 0)
	for i := range iBuckets {
		// huawei OBS is shared across projects
		if iBuckets[i].GetLocation() == region.GetId() {
			ret = append(ret, iBuckets[i])
		}
	}
	return ret, nil
}

func str2StorageClass(storageClassStr string) (obs.StorageClassType, error) {
	if strings.EqualFold(storageClassStr, string(obs.StorageClassStandard)) {
		return obs.StorageClassStandard, nil
	} else if strings.EqualFold(storageClassStr, string(obs.StorageClassWarm)) {
		return obs.StorageClassWarm, nil
	} else if strings.EqualFold(storageClassStr, string(obs.StorageClassCold)) {
		return obs.StorageClassCold, nil
	} else {
		return obs.StorageClassStandard, errors.Error("unsupported storageClass")
	}
}

func (region *SRegion) CreateIBucket(name string, storageClassStr string, aclStr string) error {
	obsClient, err := region.getOBSClient()
	if err != nil {
		return errors.Wrap(err, "region.getOBSClient")
	}
	input := &obs.CreateBucketInput{}
	input.Bucket = name
	input.Location = region.GetId()
	if len(aclStr) > 0 {
		if strings.EqualFold(aclStr, string(obs.AclPrivate)) {
			input.ACL = obs.AclPrivate
		} else if strings.EqualFold(aclStr, string(obs.AclPublicRead)) {
			input.ACL = obs.AclPublicRead
		} else if strings.EqualFold(aclStr, string(obs.AclPublicReadWrite)) {
			input.ACL = obs.AclPublicReadWrite
		} else {
			return errors.Error("unsupported acl")
		}
	}
	if len(storageClassStr) > 0 {
		input.StorageClass, err = str2StorageClass(storageClassStr)
		if err != nil {
			return err
		}
	}
	_, err = obsClient.CreateBucket(input)
	if err != nil {
		return errors.Wrap(err, "obsClient.CreateBucket")
	}
	region.client.invalidateIBuckets()
	return nil
}

func obsHttpCode(err error) int {
	switch httpErr := err.(type) {
	case obs.ObsError:
		return httpErr.StatusCode
	case *obs.ObsError:
		return httpErr.StatusCode
	}
	return -1
}

func (region *SRegion) DeleteIBucket(name string) error {
	obsClient, err := region.getOBSClient()
	if err != nil {
		return errors.Wrap(err, "region.getOBSClient")
	}
	_, err = obsClient.DeleteBucket(name)
	if err != nil {
		if obsHttpCode(err) == 404 {
			return nil
		}
		log.Debugf("%#v %s", err, err)
		return errors.Wrap(err, "DeleteBucket")
	}
	region.client.invalidateIBuckets()
	return nil
}

func (region *SRegion) HeadBucket(name string) (*obs.BaseModel, error) {
	obsClient, err := region.getOBSClient()
	if err != nil {
		return nil, errors.Wrap(err, "region.getOBSClient")
	}
	return obsClient.HeadBucket(name)
}

func (region *SRegion) IBucketExist(name string) (bool, error) {
	_, err := region.HeadBucket(name)
	if err != nil {
		if obsHttpCode(err) == 404 {
			return false, nil
		} else {
			return false, errors.Wrap(err, "HeadBucket")
		}
	}
	return true, nil
}

func (region *SRegion) GetIBucketById(name string) (cloudprovider.ICloudBucket, error) {
	return cloudprovider.GetIBucketById(region, name)
}

func (region *SRegion) GetIBucketByName(name string) (cloudprovider.ICloudBucket, error) {
	return region.GetIBucketById(name)
}

func (self *SRegion) GetSkus(zoneId string) ([]cloudprovider.ICloudSku, error) {
	return nil, cloudprovider.ErrNotImplemented
}

func (self *SRegion) GetIElasticcaches() ([]cloudprovider.ICloudElasticcache, error) {
	caches, err := self.GetElasticCaches("")
	if err != nil {
		return nil, err
	}

	icaches := make([]cloudprovider.ICloudElasticcache, len(caches))
	for i := range caches {
		caches[i].region = self
		icaches[i] = &caches[i]
	}

	return icaches, nil
}

func (region *SRegion) GetCapabilities() []string {
	return region.client.GetCapabilities()
}

func (self *SRegion) GetDiskTypes() ([]SDiskType, error) {
	resp, err := self.list(SERVICE_EVS, "types", nil)
	if err != nil {
		return nil, err
	}
	dts := []SDiskType{}
	err = resp.Unmarshal(&dts, "volume_types")
	if err != nil {
		return nil, errors.Wrap(err, "Unmarshal")
	}
	return dts, nil
}

func (self *SRegion) GetZoneSupportedDiskTypes(zoneId string) ([]string, error) {
	dts, err := self.GetDiskTypes()
	if err != nil {
		return nil, errors.Wrap(err, "GetDiskTypes")
	}

	ret := []string{}
	for i := range dts {
		if dts[i].IsAvaliableInZone(zoneId) {
			ret = append(ret, dts[i].Name)
		}
	}

	return ret, nil
}
