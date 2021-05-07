package aws

import (
	"github.com/pkg/errors"

	"yunion.io/x/cloudmux/pkg/cloudprovider"
	"yunion.io/x/cloudmux/pkg/multicloud"
)

type SInternetGateway struct {
	multicloud.SResourceBase
	region *SRegion

	Attachments       []InternetGatewayAttachment `json:"Attachments"`
	InternetGatewayID string                      `json:"InternetGatewayId"`
	OwnerID           string                      `json:"OwnerId"`
}

type InternetGatewayAttachment struct {
	State string `json:"State"`
	VpcID string `json:"VpcId"`
}

func (i *SInternetGateway) GetId() string {
	return i.InternetGatewayID
}

func (i *SInternetGateway) GetName() string {
	return i.InternetGatewayID
}

func (i *SInternetGateway) GetGlobalId() string {
	return i.GetId()
}

func (i *SInternetGateway) GetStatus() string {
	return ""
}

func (i *SInternetGateway) Refresh() error {
	return errors.Wrap(cloudprovider.ErrNotImplemented, "Refresh")
}

func (i *SInternetGateway) IsEmulated() bool {
	return false
}
