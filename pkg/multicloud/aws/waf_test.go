package aws

import (
	"strings"
	"testing"

	"yunion.io/x/jsonutils"
)

func TestWafWebAclMarshal(t *testing.T) {
	da := &sWafDefaultAction{Allow: &sWafAllowAction{}}
	got := jsonutils.Marshal(da).String()
	if !strings.Contains(got, `"Allow":{}`) {
		t.Fatalf("DefaultAction marshal want Allow={}, got %s", got)
	}
	if strings.Contains(got, "Block") {
		t.Fatalf("nil Block should be omitted, got %s", got)
	}

	acl := &SWafWebACL{
		Name:          "test",
		Id:            "id",
		ARN:           "arn",
		DefaultAction: da,
		VisibilityConfig: &sWafVisibilityConfig{
			MetricName:               "test",
			CloudWatchMetricsEnabled: true,
			SampledRequestsEnabled:   true,
		},
	}
	sw := &SWebAcl{SWafWebACL: *acl, LockToken: "tok"}
	out := jsonutils.Marshal(sw).String()
	if !strings.Contains(out, `"Name":"test"`) || !strings.Contains(out, `"Id":"id"`) {
		t.Fatalf("SWebAcl marshal lost WebACL fields: %s", out)
	}
	if !strings.Contains(out, `"lock_token":"tok"`) && !strings.Contains(out, `"LockToken":"tok"`) {
		t.Fatalf("SWebAcl marshal lost LockToken: %s", out)
	}
}

func TestGetWebAclUnmarshal(t *testing.T) {
	body := `{
  "WebACL": {
    "ARN": "arn:aws:wafv2:us-east-1:123:regional/webacl/test/abc",
    "Capacity": 5,
    "DefaultAction": {"Allow": {}},
    "Description": "desc",
    "Id": "abc-id",
    "Name": "test",
    "Rules": [],
    "VisibilityConfig": {
      "CloudWatchMetricsEnabled": true,
      "MetricName": "test",
      "SampledRequestsEnabled": true
    }
  },
  "LockToken": "lock-token"
}`
	obj, err := jsonutils.ParseString(body)
	if err != nil {
		t.Fatal(err)
	}
	resp := struct {
		WebACL    *SWafWebACL
		LockToken string
	}{}
	if err := obj.Unmarshal(&resp); err != nil {
		t.Fatal(err)
	}
	if resp.WebACL == nil || resp.WebACL.Name != "test" {
		t.Fatalf("WebACL unmarshal failed: %+v", resp.WebACL)
	}
	if resp.WebACL.DefaultAction == nil || resp.WebACL.DefaultAction.Allow == nil {
		t.Fatalf("DefaultAction lost: %+v", resp.WebACL.DefaultAction)
	}
	ret := &SWebAcl{SWafWebACL: *resp.WebACL, LockToken: resp.LockToken}
	if ret.GetName() != "test" || ret.GetGlobalId() == "" {
		t.Fatalf("GetWebAcl accessors failed: name=%s id=%s", ret.GetName(), ret.GetGlobalId())
	}
	out := jsonutils.Marshal(ret).String()
	if !strings.Contains(out, `"Name":"test"`) {
		t.Fatalf("GetWebAcl result print empty: %s", out)
	}
}

func TestListWebACLsUnmarshal(t *testing.T) {
	body := `{
  "WebACLs": [
    {
      "ARN": "arn:aws:wafv2:us-east-1:123:regional/webacl/test/abc",
      "Description": "desc",
      "Id": "abc-id",
      "LockToken": "lock-token",
      "Name": "test"
    }
  ],
  "NextMarker": ""
}`
	obj, err := jsonutils.ParseString(body)
	if err != nil {
		t.Fatal(err)
	}
	resp := struct {
		WebACLs    []SWebAcl
		NextMarker string
	}{}
	if err := obj.Unmarshal(&resp); err != nil {
		t.Fatal(err)
	}
	if len(resp.WebACLs) != 1 {
		t.Fatalf("len=%d", len(resp.WebACLs))
	}
	acl := resp.WebACLs[0]
	if acl.Name != "test" {
		t.Fatalf("list item empty: %+v", acl)
	}
	if acl.LockToken != "lock-token" {
		t.Fatalf("LockToken=%q", acl.LockToken)
	}
}
