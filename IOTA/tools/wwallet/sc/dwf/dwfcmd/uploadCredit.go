package dwfcmd

import (
	"fmt"
	"os"
	"strconv"

	"wasp/packages/txutil/vtxbuilder"
	"wasp/tools/wwallet/config"
	clientutil "wasp/tools/wwallet/util"
	"wasp/tools/wwallet/wallet"

	"github.com/iotaledger/goshimmer/dapps/valuetransfers/packages/address"
	"github.com/iotaledger/goshimmer/dapps/valuetransfers/packages/balance"
)

func uploadCreditCmd(args []string) {
	if len(args) != 1 {
		fmt.Printf("Usage: %s uploadCredit <amount>\n", os.Args[0])
		os.Exit(1)
	}

	wallet := wallet.Load()
	persona := wallet.Address() //indirizzo destinatario

	scAdd, err := address.FromBase58(args[0]) //indirizzo mittente, devo vedere come mettere quello dello sc
	check(err)

	amount, err := strconv.Atoi(args[2])
	check(err)

	bals, err := config.GoshimmerClient().GetConfirmedAccountOutputs(&scAdd) //vede se lo Sc ha denaro sufficente
	check(err)

	vtxb, err := vtxbuilder.NewFromOutputBalances(bals)
	check(err)

	check(vtxb.MoveToAddress(persona, balance.ColorIOTA, int64(amount*10)))

	tx := vtxb.Build(false)
	tx.Sign(wallet.SignatureScheme()) //mettere la signatura dello sc

	clientutil.PostTransaction(tx) //vedere cosa fa... Sono rimasto qua
}

/**
func decodeColor2(s string) *balance.Color {
	color, err := util.ColorFromString(s)
	check(err)
	return &color
}
*/
