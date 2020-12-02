package wallet

import (
	"fmt"
	"os"
	"strconv"

	"wasp/packages/txutil/vtxbuilder"
	"wasp/tools/wwallet/config"
	"wasp/tools/wwallet/util"
)

func mintCmd(args []string) {
	if len(args) < 1 {
		fmt.Printf("Usage: %s mint <amount>\n", os.Args[0])
		os.Exit(1)
	}

	wallet := Load()

	amount, err := strconv.Atoi(args[0])
	check(err)

	tx, err := vtxbuilder.NewColoredTokensTransaction(config.GoshimmerClient(), wallet.SignatureScheme(), int64(amount))
	check(err)

	util.PostTransaction(tx)

	fmt.Printf("Minted %d tokens of color %s\n", amount, tx.ID())
}
