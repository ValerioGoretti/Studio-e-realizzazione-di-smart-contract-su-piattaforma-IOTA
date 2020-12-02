package dwfcmd

import (
	"os"
	"strconv"
	"wasp/tools/wwallet/sc/fa"

	"github.com/iotaledger/goshimmer/dapps/valuetransfers/packages/balance"
	"github.com/mr-tron/base58"
)

func payPlCmd(args []string) {
	if len(args) != 2 {
		fa.Config.PrintUsage("wwallet buy payPl <color> <amount> ")
		os.Exit(1)
	}

	description := ""

	color := decodeColor(args[1])

	amount, err := strconv.Atoi(args[2])
	check(err)

	minimumBid := 0

	durationMinutes := 0

	_, err = fa.Client().StartAuction(
		description,
		color,
		int64(amount),
		int64(minimumBid),
		int64(durationMinutes),
	)
	check(err)
}

func decodeColor(s string) *balance.Color {
	b, err := base58.Decode(s)
	check(err)
	color, _, err := balance.ColorFromBytes(b)
	check(err)
	return &color
}

/**
func payPlCmd(args []string) {
	if len(args) != 1 {
		fmt.Printf("Usage: %s bus payIota 150 \n", os.Args[0])
		os.Exit(1)
	} else if b, e := strconv.Atoi(args[0]); b != 150 || e != nil {
		fmt.Printf("Usage: %s bus payIota 150 \n", os.Args[0])
		fmt.Printf("Il costo del biglietto è di 150 IOTA\n")
		os.Exit(1)
	}

	amount, err := strconv.Atoi(args[0])
	check(err)

	feedback := ""

	tx, err := dwf.Client().DonateColor(int64(amount), feedback)
	check(err)
	fmt.Printf("Biglietto acquistato! Puoi salire sull'autobus, id transazione -> %v", tx)
}
*/
