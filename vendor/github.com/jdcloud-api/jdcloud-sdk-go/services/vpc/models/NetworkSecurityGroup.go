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


type NetworkSecurityGroup struct {

    /* 安全组ID (Optional) */
    NetworkSecurityGroupId string `json:"networkSecurityGroupId"`

    /* 安全组名称 (Optional) */
    NetworkSecurityGroupName string `json:"networkSecurityGroupName"`

    /* 安全组描述信息 (Optional) */
    Description string `json:"description"`

    /* 安全组所在vpc的Id (Optional) */
    VpcId string `json:"vpcId"`

    /* 安全组规则信息 (Optional) */
    SecurityGroupRules []SecurityGroupRule `json:"securityGroupRules"`

    /* 安全组创建时间 (Optional) */
    CreatedTime string `json:"createdTime"`

    /* 安全组类型, default：默认安全组，custom：自定义安全组 (Optional) */
    NetworkSecurityGroupType string `json:"networkSecurityGroupType"`

    /* 安全组绑定的弹性网卡列表 (Optional) */
    NetworkInterfaceIds []string `json:"networkInterfaceIds"`
}
