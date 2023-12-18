package shell

import (
	"yunion.io/x/cloudmux/pkg/multicloud/huawei"
	"yunion.io/x/pkg/util/shellutils"
)

func init() {
	type SSlCertificateListOptions struct {
	}
	shellutils.R(
		&SSlCertificateListOptions{},
		"ssl-certificate-list",
		"List ssl certificates",
		func(cli *huawei.SRegion, args *SSlCertificateListOptions) error {
			certs, err := cli.GetClient().GetSSLCertificates()
			if err != nil {
				return err
			}
			printList(certs, 0, 0, 0, nil)
			return nil
		})
}
