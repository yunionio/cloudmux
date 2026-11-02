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

package huawei

import (
	"net/http"
	"testing"

	"yunion.io/x/pkg/errors"

	api "yunion.io/x/cloudmux/pkg/apis/compute"
	"yunion.io/x/cloudmux/pkg/cloudprovider"
)

func TestAgencyUrn(t *testing.T) {
	urn := agencyUrn("01a0eb06ab5475be94aa2b4112d35a81", "OrganizationAccountAccessAgency")
	if urn != "iam::01a0eb06ab5475be94aa2b4112d35a81:agency:OrganizationAccountAccessAgency" {
		t.Fatalf("urn %s", urn)
	}
}

func TestOuPathTags(t *testing.T) {
	if ouPathTags(nil) != nil {
		t.Fatalf("empty path should not produce tags")
	}
	tags := ouPathTags([]string{"prod", "app"})
	if tags["L1"] != "prod" || tags["L2"] != "app" || len(tags) != 2 {
		t.Fatalf("unexpected tags %#v", tags)
	}
}

func TestOrganizationAccountHealth(t *testing.T) {
	if organizationAccountHealth("active") != api.CLOUD_PROVIDER_HEALTH_NORMAL {
		t.Fatalf("active")
	}
	if organizationAccountHealth("pending_closure") != api.CLOUD_PROVIDER_HEALTH_PENDING {
		t.Fatalf("pending")
	}
	if organizationAccountHealth("suspended") != api.CLOUD_PROVIDER_HEALTH_SUSPENDED {
		t.Fatalf("suspended")
	}
}

func TestOrganizationSubAccount(t *testing.T) {
	cli := &SHuaweiClient{
		HuaweiClientConfig: &HuaweiClientConfig{
			accessKey: "ak",
			cpcfg:     cloudprovider.ProviderConfig{Name: "mgmt"},
		},
	}
	self := cli.organizationSubAccount(SOrganizationAccount{
		Id: "self", Name: "master", Status: "active",
	}, "self", nil)
	if self.Account != "ak" || self.Id != "self" || self.Name != "master/mgmt" {
		t.Fatalf("self account %#v", self)
	}
	if self.DefaultProjectId != "0" || self.HealthStatus != api.CLOUD_PROVIDER_HEALTH_NORMAL {
		t.Fatalf("self status %#v", self)
	}

	member := cli.organizationSubAccount(SOrganizationAccount{
		Id: "member", Name: "dev", Status: "suspended", Description: "desc",
	}, "self", map[string]string{"L1": "ou"})
	if member.Account != "ak/member" || member.Id != "member" || member.Name != "dev/member" {
		t.Fatalf("member account %#v", member)
	}
	if member.HealthStatus != api.CLOUD_PROVIDER_HEALTH_SUSPENDED || member.Desc != "desc" || member.Tags["L1"] != "ou" {
		t.Fatalf("member meta %#v", member)
	}
}

func TestIsOrgUnavailable(t *testing.T) {
	if !isOrgUnavailable(&sHuaweiError{ErrorCode: "Organizations.1100"}) {
		t.Fatalf("org not found")
	}
	if !isOrgUnavailable(&sHuaweiError{StatusCode: http.StatusForbidden}) {
		t.Fatalf("forbidden")
	}
	if !isOrgUnavailable(errors.Wrapf(cloudprovider.ErrNotFound, "missing")) {
		t.Fatalf("not found")
	}
	if isOrgUnavailable(&sHuaweiError{ErrorCode: "CBC.0150", StatusCode: http.StatusBadRequest}) {
		t.Fatalf("unrelated error should not fallback")
	}
}
