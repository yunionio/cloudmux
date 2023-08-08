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


type CustomQuery struct {

    /* 快捷检索条件名称，长度为32个字符，只允许中文、数字、大小写字母、英文下划线“_”及中划线“-”，且不允许重名 (Optional) */
    Name string `json:"name"`

    /* 命名空间 (Optional) */
    Namespace string `json:"namespace"`

    /* 命名空间名字 (Optional) */
    NamespaceName string `json:"namespaceName"`

    /*  (Optional) */
    Query DescribeQueryOption `json:"query"`

    /* 快捷检索条件Uuid (Optional) */
    QueryUuid string `json:"queryUuid"`

    /* 地域 (Optional) */
    Region string `json:"region"`
}
