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

	type SSlCertificateIdOptions struct {
		ID string
	}

	shellutils.R(
		&SSlCertificateIdOptions{},
		"ssl-certificate-show",
		"List ssl certificates",
		func(cli *huawei.SRegion, args *SSlCertificateIdOptions) error {
			cert, err := cli.GetClient().GetSSLCertificate(args.ID)
			if err != nil {
				return err
			}
			printObject(cert)
			return nil
		})

}
