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
)

type ModifyInstanceNameRequest struct {

    core.JDCloudRequest

    /* 地域代码，取值范围参见[《各地域及可用区对照表》](../Enum-Definitions/Regions-AZ.md)  */
    RegionId string `json:"regionId"`

    /* RDS 实例ID，唯一标识一个RDS实例  */
    InstanceId string `json:"instanceId"`

    /* 实例名称，名称支持中文，实例名的具体规则可参见帮助中心文档:[名称及密码限制](../../../documentation/Database-and-Cache-Service/RDS/Introduction/Restrictions/SQLServer-Restrictions.md)  */
    InstanceName string `json:"instanceName"`
}

/*
 * param regionId: 地域代码，取值范围参见[《各地域及可用区对照表》](../Enum-Definitions/Regions-AZ.md) (Required)
 * param instanceId: RDS 实例ID，唯一标识一个RDS实例 (Required)
 * param instanceName: 实例名称，名称支持中文，实例名的具体规则可参见帮助中心文档:[名称及密码限制](../../../documentation/Database-and-Cache-Service/RDS/Introduction/Restrictions/SQLServer-Restrictions.md) (Required)
 *
 * @Deprecated, not compatible when mandatory parameters changed
 */
func NewModifyInstanceNameRequest(
    regionId string,
    instanceId string,
    instanceName string,
) *ModifyInstanceNameRequest {

	return &ModifyInstanceNameRequest{
        JDCloudRequest: core.JDCloudRequest{
			URL:     "/regions/{regionId}/instances/{instanceId}:modifyInstanceName",
			Method:  "POST",
			Header:  nil,
			Version: "v1",
		},
        RegionId: regionId,
        InstanceId: instanceId,
        InstanceName: instanceName,
	}
}

/*
 * param regionId: 地域代码，取值范围参见[《各地域及可用区对照表》](../Enum-Definitions/Regions-AZ.md) (Required)
 * param instanceId: RDS 实例ID，唯一标识一个RDS实例 (Required)
 * param instanceName: 实例名称，名称支持中文，实例名的具体规则可参见帮助中心文档:[名称及密码限制](../../../documentation/Database-and-Cache-Service/RDS/Introduction/Restrictions/SQLServer-Restrictions.md) (Required)
 */
func NewModifyInstanceNameRequestWithAllParams(
    regionId string,
    instanceId string,
    instanceName string,
) *ModifyInstanceNameRequest {

    return &ModifyInstanceNameRequest{
        JDCloudRequest: core.JDCloudRequest{
            URL:     "/regions/{regionId}/instances/{instanceId}:modifyInstanceName",
            Method:  "POST",
            Header:  nil,
            Version: "v1",
        },
        RegionId: regionId,
        InstanceId: instanceId,
        InstanceName: instanceName,
    }
}

/* This constructor has better compatible ability when API parameters changed */
func NewModifyInstanceNameRequestWithoutParam() *ModifyInstanceNameRequest {

    return &ModifyInstanceNameRequest{
            JDCloudRequest: core.JDCloudRequest{
            URL:     "/regions/{regionId}/instances/{instanceId}:modifyInstanceName",
            Method:  "POST",
            Header:  nil,
            Version: "v1",
        },
    }
}

/* param regionId: 地域代码，取值范围参见[《各地域及可用区对照表》](../Enum-Definitions/Regions-AZ.md)(Required) */
func (r *ModifyInstanceNameRequest) SetRegionId(regionId string) {
    r.RegionId = regionId
}

/* param instanceId: RDS 实例ID，唯一标识一个RDS实例(Required) */
func (r *ModifyInstanceNameRequest) SetInstanceId(instanceId string) {
    r.InstanceId = instanceId
}

/* param instanceName: 实例名称，名称支持中文，实例名的具体规则可参见帮助中心文档:[名称及密码限制](../../../documentation/Database-and-Cache-Service/RDS/Introduction/Restrictions/SQLServer-Restrictions.md)(Required) */
func (r *ModifyInstanceNameRequest) SetInstanceName(instanceName string) {
    r.InstanceName = instanceName
}

// GetRegionId returns path parameter 'regionId' if exist,
// otherwise return empty string
func (r ModifyInstanceNameRequest) GetRegionId() string {
    return r.RegionId
}

type ModifyInstanceNameResponse struct {
    RequestID string `json:"requestId"`
    Error core.ErrorResponse `json:"error"`
    Result ModifyInstanceNameResult `json:"result"`
}

type ModifyInstanceNameResult struct {
}