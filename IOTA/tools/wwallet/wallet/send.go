package wallet

import (
	"fmt"
	"os"
	"strconv"

	"wasp/packages/txutil/vtxbuilder"
	"wasp/packages/util"
	"wasp/tools/wwallet/config"
	clientutil "wasp/tools/wwallet/util"

	"github.com/iotaledger/goshimmer/dapps/valuetransfers/packages/address"
	"github.com/iotaledger/goshimmer/dapps/valuetransfers/packages/balance"
)

func sendFundsCmd(args []string) {
	if len(args) < 3 {
		fmt.Printf("Usage: %s send-funds <target-address> <color> <amount>\n", os.Args[0])
		os.Exit(1)
	}

	wallet := Load()
	sourceAddress := wallet.Address()

	targetAddress, err := address.FromBase58(args[0])
	check(err)

	color := decodeColor(args[1])

	amount, err := strconv.Atoi(args[2])
	check(err)

	bals, err := config.GoshimmerClient().GetConfirmedAccountOutputs(&sourceAddress)
	check(err)

	vtxb, err := vtxbuilder.NewFromOutputBalances(bals)
	check(err)

	check(vtxb.MoveToAddress(targetAddress, *color, int64(amount)))

	tx := vtxb.Build(false)
	tx.Sign(wallet.SignatureScheme())

	clientutil.PostTransaction(tx)
}

func decodeColor(s string) *balance.Color {
	color, err := util.ColorFromString(s)
	check(err)
	return &color
}
