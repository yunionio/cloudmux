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

package harvester

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"yunion.io/x/jsonutils"
	"yunion.io/x/pkg/errors"
	"yunion.io/x/pkg/gotypes"
	"yunion.io/x/pkg/util/httputils"

	api "yunion.io/x/cloudmux/pkg/apis/compute"
	"yunion.io/x/cloudmux/pkg/cloudprovider"
)

const (
	CLOUD_PROVIDER_HARVESTER = api.CLOUD_PROVIDER_HARVESTER
	HARVESTER_DEFAULT_REGION = "Harvester"
)

type HarvesterClientConfig struct {
	cpcfg cloudprovider.ProviderConfig

	host            string
	accessKey       string
	accessKeySecret string

	debug bool
}

func NewHarvesterClientConfig(host, accessKey, accessKeySecret string) *HarvesterClientConfig {
	cfg := &HarvesterClientConfig{
		host:            host,
		accessKey:       accessKey,
		accessKeySecret: accessKeySecret,
	}
	return cfg
}

func (cfg *HarvesterClientConfig) CloudproviderConfig(cpcfg cloudprovider.ProviderConfig) *HarvesterClientConfig {
	cfg.cpcfg = cpcfg
	return cfg
}

func (cfg *HarvesterClientConfig) Debug(debug bool) *HarvesterClientConfig {
	cfg.debug = debug
	return cfg
}

type SHarvesterClient struct {
	*HarvesterClientConfig

	httpClient *http.Client
}

func NewHarvesterClient(cfg *HarvesterClientConfig) (*SHarvesterClient, error) {
	httpClient := cfg.cpcfg.AdaptiveTimeoutHttpClient()
	ts, _ := httpClient.Transport.(*http.Transport)
	httpClient.Transport = cloudprovider.GetCheckTransport(ts, func(req *http.Request) (func(resp *http.Response) error, error) {
		if cfg.cpcfg.ReadOnly {
			if req.Method == "GET" || req.Method == "HEAD" {
				return nil, nil
			}
			return nil, errors.Wrapf(cloudprovider.ErrAccountReadOnly, "%s %s", req.Method, req.URL.Path)
		}
		return nil, nil
	})

	cli := &SHarvesterClient{
		HarvesterClientConfig: cfg,
		httpClient:            httpClient,
	}
	return cli, cli.getUser()
}

func (cli *SHarvesterClient) GetCloudRegionExternalIdPrefix() string {
	return fmt.Sprintf("%s/%s", CLOUD_PROVIDER_HARVESTER, cli.cpcfg.Id)
}

func (cli *SHarvesterClient) GetSubAccounts() ([]cloudprovider.SSubAccount, error) {
	subAccount := cloudprovider.SSubAccount{
		Id:           cli.cpcfg.Id,
		Account:      cli.cpcfg.Id,
		Name:         cli.cpcfg.Name,
		HealthStatus: api.CLOUD_PROVIDER_HEALTH_NORMAL,
	}
	return []cloudprovider.SSubAccount{subAccount}, nil
}

func (cli *SHarvesterClient) GetIRegions() ([]cloudprovider.ICloudRegion, error) {
	return nil, cloudprovider.ErrNotImplemented
	//return cli.iregions, nil
}

func (cli *SHarvesterClient) GetIRegionById(id string) (cloudprovider.ICloudRegion, error) {
	//for i := 0; i < len(cli.iregions); i++ {
	//	if cli.iregions[i].GetGlobalId() == id {
	//		return cli.iregions[i], nil
	//	}
	//}
	return nil, cloudprovider.ErrNotFound
}

func (cli *SHarvesterClient) Do(req *http.Request) (*http.Response, error) {
	req.Header.Set("Authorization", fmt.Sprintf("bearer %s:%s", cli.accessKey, cli.accessKeySecret))
	return cli.httpClient.Do(req)
}

func (cli *SHarvesterClient) getUser() error {
	params := url.Values{}
	params.Set("me", "true")
	_, err := cli.list("v3/users", params)
	return err
}

func (cli *SHarvesterClient) list(res string, params url.Values) (jsonutils.JSONObject, error) {
	if gotypes.IsNil(params) {
		params = url.Values{}
	}
	ret := jsonutils.NewArray()
	for {
		resp, err := cli._list(res, params)
		if err != nil {
			return nil, err
		}
		part := struct {
			Data     []jsonutils.JSONObject
			Continue string
		}{}
		err = resp.Unmarshal(&part)
		if err != nil {
			return nil, err
		}
		ret.Add(part.Data...)
		if len(part.Continue) == 0 {
			break
		}
		params.Set("continue", part.Continue)
	}
	return ret, nil
}

func (cli *SHarvesterClient) _list(res string, params url.Values) (jsonutils.JSONObject, error) {
	return cli.request(httputils.GET, res, params, nil)
}

func (cli *SHarvesterClient) request(method httputils.THttpMethod, res string, params url.Values, body jsonutils.JSONObject) (jsonutils.JSONObject, error) {
	uri := fmt.Sprintf("https://%s/%s", cli.host, strings.TrimPrefix(res, "/"))
	if !gotypes.IsNil(params) && len(params) > 0 {
		uri = fmt.Sprintf("%s?%s", uri, params.Encode())
	}
	_, resp, err := httputils.JSONRequest(cli, context.Background(), method, uri, nil, body, cli.debug)
	return resp, err
}

func (cli *SHarvesterClient) GetRegion() *SRegion {
	return &SRegion{client: cli}
}

func (cli *SHarvesterClient) GetIProjects() ([]cloudprovider.ICloudProject, error) {
	return nil, cloudprovider.ErrNotImplemented
}

func (self *SHarvesterClient) GetCapabilities() []string {
	caps := []string{
		cloudprovider.CLOUD_CAPABILITY_PROJECT,
		cloudprovider.CLOUD_CAPABILITY_COMPUTE,
		cloudprovider.CLOUD_CAPABILITY_NETWORK,
	}
	return caps
}
