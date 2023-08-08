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


type ServiceInfoB struct {

    /*  (Optional) */
    DataMapping DataMappingB `json:"dataMapping"`

    /*  (Optional) */
    Filters interface{} `json:"filters"`

    /*  (Optional) */
    Ipv4Regions []LvInfoB `json:"ipv4Regions"`

    /*  (Optional) */
    Ipv6Regions []LvInfoB `json:"ipv6Regions"`

    /*  (Optional) */
    ListApi ListAPIB `json:"listApi"`

    /*  (Optional) */
    Products []ProductB `json:"products"`

    /*  (Optional) */
    Search []SearchB `json:"search"`

    /*  (Optional) */
    ServiceCode string `json:"serviceCode"`

    /* 中英文适配 (Optional) */
    ServiceName string `json:"serviceName"`

    /*  (Optional) */
    ServiceNameCn string `json:"serviceNameCn"`

    /*  (Optional) */
    ServiceNameEn string `json:"serviceNameEn"`

    /*  (Optional) */
    Status []StatusB `json:"status"`

    /*  (Optional) */
    TableColumns []TableColumnB `json:"tableColumns"`
}
