package shell

import (
	"yunion.io/x/pkg/util/shellutils"

	"yunion.io/x/cloudmux/pkg/multicloud/ecloud"
)

func init() {
	type BalanceOptions struct {
	}
	shellutils.R(&BalanceOptions{}, "balance-show", "Show account balance (MOPC)", func(cli *ecloud.SRegion, args *BalanceOptions) error {
		balance, err := cli.GetClient().GetBalance()
		if err != nil {
			return err
		}
		printObject(balance)
		return nil
	})
}
