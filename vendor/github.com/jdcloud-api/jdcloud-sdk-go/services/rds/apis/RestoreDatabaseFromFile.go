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

type RestoreDatabaseFromFileRequest struct {

    core.JDCloudRequest

    /* 地域代码，取值范围参见[《各地域及可用区对照表》](../Enum-Definitions/Regions-AZ.md)  */
    RegionId string `json:"regionId"`

    /* RDS 实例ID，唯一标识一个RDS实例  */
    InstanceId string `json:"instanceId"`

    /* 库名称  */
    DbName string `json:"dbName"`

    /* 共享文件的全局ID，可从上传文件查询接口[describeImportFiles](../Cloud-on-Single-Database/describeImportFiles.md)获取；如果该文件不是共享文件，则不用输入该参数 (Optional) */
    SharedFileGid *string `json:"sharedFileGid"`

    /* 用户上传的备份文件名称（包括文件后缀名），例如mydb1.bak  */
    FileName string `json:"fileName"`
}

/*
 * param regionId: 地域代码，取值范围参见[《各地域及可用区对照表》](../Enum-Definitions/Regions-AZ.md) (Required)
 * param instanceId: RDS 实例ID，唯一标识一个RDS实例 (Required)
 * param dbName: 库名称 (Required)
 * param fileName: 用户上传的备份文件名称（包括文件后缀名），例如mydb1.bak (Required)
 *
 * @Deprecated, not compatible when mandatory parameters changed
 */
func NewRestoreDatabaseFromFileRequest(
    regionId string,
    instanceId string,
    dbName string,
    fileName string,
) *RestoreDatabaseFromFileRequest {

	return &RestoreDatabaseFromFileRequest{
        JDCloudRequest: core.JDCloudRequest{
			URL:     "/regions/{regionId}/instances/{instanceId}/databases/{dbName}:restoreDatabaseFromFile",
			Method:  "POST",
			Header:  nil,
			Version: "v1",
		},
        RegionId: regionId,
        InstanceId: instanceId,
        DbName: dbName,
        FileName: fileName,
	}
}

/*
 * param regionId: 地域代码，取值范围参见[《各地域及可用区对照表》](../Enum-Definitions/Regions-AZ.md) (Required)
 * param instanceId: RDS 实例ID，唯一标识一个RDS实例 (Required)
 * param dbName: 库名称 (Required)
 * param sharedFileGid: 共享文件的全局ID，可从上传文件查询接口[describeImportFiles](../Cloud-on-Single-Database/describeImportFiles.md)获取；如果该文件不是共享文件，则不用输入该参数 (Optional)
 * param fileName: 用户上传的备份文件名称（包括文件后缀名），例如mydb1.bak (Required)
 */
func NewRestoreDatabaseFromFileRequestWithAllParams(
    regionId string,
    instanceId string,
    dbName string,
    sharedFileGid *string,
    fileName string,
) *RestoreDatabaseFromFileRequest {

    return &RestoreDatabaseFromFileRequest{
        JDCloudRequest: core.JDCloudRequest{
            URL:     "/regions/{regionId}/instances/{instanceId}/databases/{dbName}:restoreDatabaseFromFile",
            Method:  "POST",
            Header:  nil,
            Version: "v1",
        },
        RegionId: regionId,
        InstanceId: instanceId,
        DbName: dbName,
        SharedFileGid: sharedFileGid,
        FileName: fileName,
    }
}

/* This constructor has better compatible ability when API parameters changed */
func NewRestoreDatabaseFromFileRequestWithoutParam() *RestoreDatabaseFromFileRequest {

    return &RestoreDatabaseFromFileRequest{
            JDCloudRequest: core.JDCloudRequest{
            URL:     "/regions/{regionId}/instances/{instanceId}/databases/{dbName}:restoreDatabaseFromFile",
            Method:  "POST",
            Header:  nil,
            Version: "v1",
        },
    }
}

/* param regionId: 地域代码，取值范围参见[《各地域及可用区对照表》](../Enum-Definitions/Regions-AZ.md)(Required) */
func (r *RestoreDatabaseFromFileRequest) SetRegionId(regionId string) {
    r.RegionId = regionId
}

/* param instanceId: RDS 实例ID，唯一标识一个RDS实例(Required) */
func (r *RestoreDatabaseFromFileRequest) SetInstanceId(instanceId string) {
    r.InstanceId = instanceId
}

/* param dbName: 库名称(Required) */
func (r *RestoreDatabaseFromFileRequest) SetDbName(dbName string) {
    r.DbName = dbName
}

/* param sharedFileGid: 共享文件的全局ID，可从上传文件查询接口[describeImportFiles](../Cloud-on-Single-Database/describeImportFiles.md)获取；如果该文件不是共享文件，则不用输入该参数(Optional) */
func (r *RestoreDatabaseFromFileRequest) SetSharedFileGid(sharedFileGid string) {
    r.SharedFileGid = &sharedFileGid
}

/* param fileName: 用户上传的备份文件名称（包括文件后缀名），例如mydb1.bak(Required) */
func (r *RestoreDatabaseFromFileRequest) SetFileName(fileName string) {
    r.FileName = fileName
}

// GetRegionId returns path parameter 'regionId' if exist,
// otherwise return empty string
func (r RestoreDatabaseFromFileRequest) GetRegionId() string {
    return r.RegionId
}

type RestoreDatabaseFromFileResponse struct {
    RequestID string `json:"requestId"`
    Error core.ErrorResponse `json:"error"`
    Result RestoreDatabaseFromFileResult `json:"result"`
}

type RestoreDatabaseFromFileResult struct {
}