package shell

import (
	"yunion.io/x/pkg/util/shellutils"

	"yunion.io/x/cloudmux/pkg/multicloud/aliyun"
)

func init() {
	type SSlCertificateListOptions struct {
		Page int
		Size int
	}
	shellutils.R(
		&SSlCertificateListOptions{},
		"sslcertificate-list",
		"List ssl certificates",
		func(cli *aliyun.SRegion, args *SSlCertificateListOptions) error {
			certs, _, err := cli.GetClient().GetSSLCertificates(args.Size, args.Page)
			if err != nil {
				return err
			}
			printList(certs, 0, 0, 0, nil)
			return nil
		})
}