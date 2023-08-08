// Copyright 2018 JDCLOUD.COM
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
//
// NOTE: This class is auto generated by the jdcloud code generator program.

package apis

import (
    "github.com/jdcloud-api/jdcloud-sdk-go/core"
    vm "github.com/jdcloud-api/jdcloud-sdk-go/services/vm/models"
    common "github.com/jdcloud-api/jdcloud-sdk-go/services/common/models"
)

type DescribeInstancePrivateIpAddressRequest struct {

    core.JDCloudRequest

    /* 地域ID  */
    RegionId string `json:"regionId"`

    /* 页码；默认为1 (Optional) */
    PageNumber *int `json:"pageNumber"`

    /* 分页大小；默认为20；取值范围[10, 100] (Optional) */
    PageSize *int `json:"pageSize"`

    /* instanceId - 云主机ID，精确匹配，支持多个
privateIpAddress - 主网卡内网主IP地址，模糊匹配，支持多个
vpcId - 私有网络ID，精确匹配，支持多个
status - 云主机状态，精确匹配，支持多个，<a href="http://docs.jdcloud.com/virtual-machines/api/vm_status">参考云主机状态</a>
name - 云主机名称，模糊匹配，支持单个
imageId - 镜像ID，精确匹配，支持多个
networkInterfaceId - 弹性网卡ID，精确匹配，支持多个
subnetId - 子网ID，精确匹配，支持多个
chargeOnStopped - 停机不计费标志，keepCharging、stopCharging 或者 notApplicable
 (Optional) */
    Filters []common.Filter `json:"filters"`
}

/*
 * param regionId: 地域ID (Required)
 *
 * @Deprecated, not compatible when mandatory parameters changed
 */
func NewDescribeInstancePrivateIpAddressRequest(
    regionId string,
) *DescribeInstancePrivateIpAddressRequest {

	return &DescribeInstancePrivateIpAddressRequest{
        JDCloudRequest: core.JDCloudRequest{
			URL:     "/regions/{regionId}/instancePrivateIpAddress",
			Method:  "GET",
			Header:  nil,
			Version: "v1",
		},
        RegionId: regionId,
	}
}

/*
 * param regionId: 地域ID (Required)
 * param pageNumber: 页码；默认为1 (Optional)
 * param pageSize: 分页大小；默认为20；取值范围[10, 100] (Optional)
 * param filters: instanceId - 云主机ID，精确匹配，支持多个
privateIpAddress - 主网卡内网主IP地址，模糊匹配，支持多个
vpcId - 私有网络ID，精确匹配，支持多个
status - 云主机状态，精确匹配，支持多个，<a href="http://docs.jdcloud.com/virtual-machines/api/vm_status">参考云主机状态</a>
name - 云主机名称，模糊匹配，支持单个
imageId - 镜像ID，精确匹配，支持多个
networkInterfaceId - 弹性网卡ID，精确匹配，支持多个
subnetId - 子网ID，精确匹配，支持多个
chargeOnStopped - 停机不计费标志，keepCharging、stopCharging 或者 notApplicable
 (Optional)
 */
func NewDescribeInstancePrivateIpAddressRequestWithAllParams(
    regionId string,
    pageNumber *int,
    pageSize *int,
    filters []common.Filter,
) *DescribeInstancePrivateIpAddressRequest {

    return &DescribeInstancePrivateIpAddressRequest{
        JDCloudRequest: core.JDCloudRequest{
            URL:     "/regions/{regionId}/instancePrivateIpAddress",
            Method:  "GET",
            Header:  nil,
            Version: "v1",
        },
        RegionId: regionId,
        PageNumber: pageNumber,
        PageSize: pageSize,
        Filters: filters,
    }
}

/* This constructor has better compatible ability when API parameters changed */
func NewDescribeInstancePrivateIpAddressRequestWithoutParam() *DescribeInstancePrivateIpAddressRequest {

    return &DescribeInstancePrivateIpAddressRequest{
            JDCloudRequest: core.JDCloudRequest{
            URL:     "/regions/{regionId}/instancePrivateIpAddress",
            Method:  "GET",
            Header:  nil,
            Version: "v1",
        },
    }
}

/* param regionId: 地域ID(Required) */
func (r *DescribeInstancePrivateIpAddressRequest) SetRegionId(regionId string) {
    r.RegionId = regionId
}

/* param pageNumber: 页码；默认为1(Optional) */
func (r *DescribeInstancePrivateIpAddressRequest) SetPageNumber(pageNumber int) {
    r.PageNumber = &pageNumber
}

/* param pageSize: 分页大小；默认为20；取值范围[10, 100](Optional) */
func (r *DescribeInstancePrivateIpAddressRequest) SetPageSize(pageSize int) {
    r.PageSize = &pageSize
}

/* param filters: instanceId - 云主机ID，精确匹配，支持多个
privateIpAddress - 主网卡内网主IP地址，模糊匹配，支持多个
vpcId - 私有网络ID，精确匹配，支持多个
status - 云主机状态，精确匹配，支持多个，<a href="http://docs.jdcloud.com/virtual-machines/api/vm_status">参考云主机状态</a>
name - 云主机名称，模糊匹配，支持单个
imageId - 镜像ID，精确匹配，支持多个
networkInterfaceId - 弹性网卡ID，精确匹配，支持多个
subnetId - 子网ID，精确匹配，支持多个
chargeOnStopped - 停机不计费标志，keepCharging、stopCharging 或者 notApplicable
(Optional) */
func (r *DescribeInstancePrivateIpAddressRequest) SetFilters(filters []common.Filter) {
    r.Filters = filters
}

// GetRegionId returns path parameter 'regionId' if exist,
// otherwise return empty string
func (r DescribeInstancePrivateIpAddressRequest) GetRegionId() string {
    return r.RegionId
}

type DescribeInstancePrivateIpAddressResponse struct {
    RequestID string `json:"requestId"`
    Error core.ErrorResponse `json:"error"`
    Result DescribeInstancePrivateIpAddressResult `json:"result"`
}

type DescribeInstancePrivateIpAddressResult struct {
    InstancePrivateIpAddress []vm.InstancePrivateIpAddress `json:"instancePrivateIpAddress"`
    TotalCount int `json:"totalCount"`
}