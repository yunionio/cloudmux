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
    charge "github.com/jdcloud-api/jdcloud-sdk-go/services/charge/models"
)

type DescribeBackupChargeRequest struct {

    core.JDCloudRequest

    /* 地域代码，取值范围参见[《各地域及可用区对照表》](../Enum-Definitions/Regions-AZ.md)  */
    RegionId string `json:"regionId"`

    /* 实例引擎类型  */
    Engine string `json:"engine"`
}

/*
 * param regionId: 地域代码，取值范围参见[《各地域及可用区对照表》](../Enum-Definitions/Regions-AZ.md) (Required)
 * param engine: 实例引擎类型 (Required)
 *
 * @Deprecated, not compatible when mandatory parameters changed
 */
func NewDescribeBackupChargeRequest(
    regionId string,
    engine string,
) *DescribeBackupChargeRequest {

	return &DescribeBackupChargeRequest{
        JDCloudRequest: core.JDCloudRequest{
			URL:     "/regions/{regionId}/instances:describeBackupCharge",
			Method:  "GET",
			Header:  nil,
			Version: "v1",
		},
        RegionId: regionId,
        Engine: engine,
	}
}

/*
 * param regionId: 地域代码，取值范围参见[《各地域及可用区对照表》](../Enum-Definitions/Regions-AZ.md) (Required)
 * param engine: 实例引擎类型 (Required)
 */
func NewDescribeBackupChargeRequestWithAllParams(
    regionId string,
    engine string,
) *DescribeBackupChargeRequest {

    return &DescribeBackupChargeRequest{
        JDCloudRequest: core.JDCloudRequest{
            URL:     "/regions/{regionId}/instances:describeBackupCharge",
            Method:  "GET",
            Header:  nil,
            Version: "v1",
        },
        RegionId: regionId,
        Engine: engine,
    }
}

/* This constructor has better compatible ability when API parameters changed */
func NewDescribeBackupChargeRequestWithoutParam() *DescribeBackupChargeRequest {

    return &DescribeBackupChargeRequest{
            JDCloudRequest: core.JDCloudRequest{
            URL:     "/regions/{regionId}/instances:describeBackupCharge",
            Method:  "GET",
            Header:  nil,
            Version: "v1",
        },
    }
}

/* param regionId: 地域代码，取值范围参见[《各地域及可用区对照表》](../Enum-Definitions/Regions-AZ.md)(Required) */
func (r *DescribeBackupChargeRequest) SetRegionId(regionId string) {
    r.RegionId = regionId
}

/* param engine: 实例引擎类型(Required) */
func (r *DescribeBackupChargeRequest) SetEngine(engine string) {
    r.Engine = engine
}

// GetRegionId returns path parameter 'regionId' if exist,
// otherwise return empty string
func (r DescribeBackupChargeRequest) GetRegionId() string {
    return r.RegionId
}

type DescribeBackupChargeResponse struct {
    RequestID string `json:"requestId"`
    Error core.ErrorResponse `json:"error"`
    Result DescribeBackupChargeResult `json:"result"`
}

type DescribeBackupChargeResult struct {
    Charge charge.Charge `json:"charge"`
}