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

package models


type InstanceTemplateDisk struct {

    /* 云硬盘类型，取值为 ssd、premium-hdd、hdd.std1、ssd.gp1、ssd.io1 (Optional) */
    DiskType string `json:"diskType"`

    /* 云硬盘大小，单位为 GiB；ssd 类型取值范围[20,1000]GB，步长为10G，premium-hdd 类型取值范围[20,3000]GB，步长为10G，hdd.std1、ssd.gp1、ssd.io1 类型取值范围[20-16000]GB，步长为10GB (Optional) */
    DiskSizeGB int `json:"diskSizeGB"`

    /* 创建云硬盘的快照ID (Optional) */
    SnapshotId string `json:"snapshotId"`

    /* 策略ID (Optional) */
    PolicyId string `json:"policyId"`

    /* 是否加密，false:(默认)不加密；true:加密 (Optional) */
    Encrypt bool `json:"encrypt"`

    /* 云硬盘的iops值 (Optional) */
    Iops int `json:"iops"`
}
