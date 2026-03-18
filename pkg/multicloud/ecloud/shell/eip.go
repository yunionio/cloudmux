package shell

import (
	"yunion.io/x/pkg/util/shellutils"

	"yunion.io/x/cloudmux/pkg/multicloud/ecloud"
)

func init() {
	type EipListOptions struct {
	}
	shellutils.R(&EipListOptions{}, "eip-list", "List eips", func(cli *ecloud.SRegion, args *EipListOptions) error {
		eips, err := cli.GetIEips()
		if err != nil {
			return err
		}
		printList(eips)
		return nil
	})

	type EipShowOptions struct {
		ID string `help:"EIP ID"`
	}
	shellutils.R(&EipShowOptions{}, "eip-show", "Show eip detail", func(cli *ecloud.SRegion, args *EipShowOptions) error {
		eip, err := cli.GetEipById(args.ID)
		if err != nil {
			return err
		}
		printObject(eip)
		return nil
	})

	type EipShowByAddrOptions struct {
		Addr string `help:"EIP address"`
	}
	shellutils.R(&EipShowByAddrOptions{}, "eip-show-by-addr", "Show eip detail by address", func(cli *ecloud.SRegion, args *EipShowByAddrOptions) error {
		eip, err := cli.GetEipByAddr(args.Addr)
		if err != nil {
			return err
		}
		printObject(eip)
		return nil
	})
}

