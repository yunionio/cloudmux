// Copyright 2026 Yunion
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package volcengine

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"yunion.io/x/cloudmux/pkg/cloudprovider"
	"yunion.io/x/jsonutils"
	"yunion.io/x/pkg/errors"
)

// Explicit service/action allowlists keep inventory credentials out of arbitrary
// endpoints and prevent a future caller from accidentally issuing a write.
func (client *SVolcEngineClient) inventoryRequest(region, service, action string, params map[string]interface{}) (jsonutils.JSONObject, error) {
	if !postgresqlRegionPattern.MatchString(region) {
		return nil, errors.Errorf("invalid inventory region")
	}
	version, method := "", http.MethodPost
	switch service {
	case "redis":
		version = "2020-12-07"
		switch action {
		case "DescribeRegions", "DescribeDBInstances", "DescribeDBInstanceDetail":
		default:
			return nil, cloudprovider.ErrNotSupported
		}
	case "rds-mysql", "rds-mssql":
		version = "2022-01-01"
		switch action {
		case "DescribeRegions", "DescribeDBInstances", "DescribeDBInstanceDetail":
		default:
			return nil, cloudprovider.ErrNotSupported
		}
	case "clb":
		version, method = "2020-04-01", http.MethodGet
		switch action {
		case "DescribeLoadBalancers", "DescribeLoadBalancerAttributes", "DescribeListeners", "DescribeServerGroups", "DescribeServerGroupAttributes":
		default:
			return nil, cloudprovider.ErrNotSupported
		}
	case "vke":
		version = "2022-05-12"
		switch action {
		case "ListClusters", "ListNodes", "ListNodePools":
		default:
			return nil, cloudprovider.ErrNotSupported
		}
	default:
		return nil, cloudprovider.ErrNotSupported
	}
	query := url.Values{"Action": {action}, "Version": {version}}
	body := ""
	if method == http.MethodGet {
		for key, value := range params {
			query.Set(key, fmt.Sprint(value))
		}
	} else {
		body = jsonutils.Marshal(params).String()
	}
	req, err := http.NewRequest(method, fmt.Sprintf("https://%s.%s.volcengineapi.com/?%s", service, region, query.Encode()), strings.NewReader(body))
	if err != nil {
		return nil, err
	}
	if method == http.MethodPost {
		req.Header.Set("Content-Type", "application/json")
	}
	// Reuse the verified, proxy-aware transport. The PostgreSQL region catalog is
	// separate and is never used to decide support for another product.
	resp, err := client.getPostgresqlClient().httpClient.Do(client.getDefaultCredential(region, strings.ReplaceAll(service, "-", "_")).Sign(req))
	if err != nil {
		return nil, errors.Wrapf(err, "%s %s", service, action)
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(io.LimitReader(resp.Body, 8*1024*1024+1))
	if err != nil {
		return nil, err
	}
	if len(data) > 8*1024*1024 {
		return nil, errors.Errorf("inventory response too large")
	}
	obj, err := jsonutils.Parse(data)
	if err != nil {
		return nil, errors.Wrapf(err, "invalid %s %s response", service, action)
	}
	if obj.Contains("ResponseMetadata", "Error") {
		ve := &sVolcError{StatusCode: resp.StatusCode}
		if err := obj.Unmarshal(ve, "ResponseMetadata"); err != nil {
			return nil, err
		}
		return nil, ve
	}
	if resp.StatusCode != http.StatusOK {
		return nil, errors.Errorf("%s %s HTTP %d", service, action, resp.StatusCode)
	}
	return obj.Get("Result")
}

// Return a complete and consistent list, or an error. In particular, an error
// must never turn into a successful empty list during resource reconciliation.
func (region *SRegion) inventoryPages(service, action, arrayKey, totalKey, idKey string, filter map[string]interface{}) ([]jsonutils.JSONObject, error) {
	params := map[string]interface{}{}
	for key, value := range filter {
		params[key] = value
	}
	params["PageSize"] = 100
	items := []jsonutils.JSONObject{}
	seen := map[string]bool{}
	expected := int64(-1)
	for page := 1; page <= 1000; page++ {
		params["PageNumber"] = page
		result, err := region.client.inventoryRequest(region.RegionId, service, action, params)
		if err != nil {
			return nil, err
		}
		count, err := result.Int(totalKey)
		if err != nil || count < 0 {
			return nil, errors.Errorf("invalid %s %s", action, totalKey)
		}
		if expected < 0 {
			expected = count
		} else if expected != count {
			return nil, errors.Errorf("%s inventory changed during pagination; retry", action)
		}
		batch, err := result.GetArray(arrayKey)
		// Redis omits Instances entirely for explicitly empty regions. A missing
		// positive-total list, null array or API error is still rejected.
		if err != nil && (service == "redis" || service == "rds-mysql" || service == "rds-mssql") && expected == 0 && !result.Contains(arrayKey) {
			return items, nil
		}
		if err != nil {
			return nil, errors.Wrapf(err, "invalid %s %s", action, arrayKey)
		}
		for _, item := range batch {
			id, err := item.GetString(idKey)
			if err != nil || id == "" || seen[id] {
				return nil, errors.Errorf("%s empty or repeated inventory ID", action)
			}
			seen[id] = true
			items = append(items, item)
		}
		if int64(len(items)) == expected {
			return items, nil
		}
		if int64(len(items)) > expected || len(batch) == 0 {
			return nil, errors.Errorf("%s incomplete inventory: %s/%s", action, strconv.Itoa(len(items)), strconv.FormatInt(expected, 10))
		}
	}
	return nil, errors.Errorf("%s pagination exceeded safety limit", action)
}
