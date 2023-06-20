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

package shell

type GlobalOptions struct {
	Debug      bool   `help:"Debug mode"`
	SUBCOMMAND string `help:"Cloudmux client subcommand" subcommand:"true"`

	Provider string `help:"Cloud provider" required:"true" choices:"Aliyun|Aws|Azure|Qcloud"`

	CloudEnv  string `help:"Cloud environment" default:"$CLOUDMUX_CLOUD_ENV" choices:"InternationalCloud|FinanceCloud|ChinaCloud|AzureGermanCloud|AzureChinaCloud|AzureUSGovernmentCloud|AzurePublicCloud" metavar:"CLOUDMUX_CLOUD_ENV"`
	AccessKey string `help:"Access key" default:"$CLOUDMUX_ACCESS_KEY" metavar:"CLOUDMUX_ACCESS_KEY"`
	Secret    string `help:"Secret" default:"$CLOUDMUX_SECRET" metavar:"CLOUDMUX_SECRET"`
	Region    string `help:"Default region" default:"$CLOUDMUX_REGION" metavar:"CLOUDMUX_REGION" short-token:"r"`
}

type EmptyOption struct{}

type IListOption interface {
	GetLimit() int
	GetOffset() int
	IsDetails() bool
	GetColumns() []string
}

type ListBaseOptions struct {
	Limit   int  `help:"Page size"`
	Offset  int  `help:"Page offset"`
	Details bool `help:"Show Details"`
}

func (o ListBaseOptions) GetLimit() int {
	return o.Limit
}

func (o ListBaseOptions) GetOffset() int {
	return o.Offset
}

func (o ListBaseOptions) IsDetails() bool {
	return o.Details
}

func (o ListBaseOptions) GetColumns() []string {
	return nil
}

type IZoneBaseOptions interface {
	GetZoneId() string
}

type ZoneBaseOptions struct {
	Zone string `help:"Zone ID"`
}

func (o ZoneBaseOptions) GetZoneId() string {
	return o.Zone
}

type IHostBaseOptions interface {
	IZoneBaseOptions
	GetHostId() string
}

type HostBaseOptions struct {
	ZoneBaseOptions
	Host string `help:"Host ID"`
}

func (o HostBaseOptions) GetHostId() string {
	return o.Host
}
