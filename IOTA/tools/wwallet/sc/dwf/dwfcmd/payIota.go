package dwfcmd

import (
	"fmt"
	"os"
	"strconv"

	"wasp/tools/wwallet/sc/dwf"
)

func payIotaCmd(args []string) {
	if len(args) != 1 {
		fmt.Printf("Usage: %s bus payIota 15 \n", os.Args[0])
		os.Exit(1)
	} else if b, e := strconv.Atoi(args[0]); b != 15 || e != nil {
		fmt.Printf("Usage: %s bus payIota 15 \n", os.Args[0])
		fmt.Printf("Il costo del biglietto è di 15 IOTA\n")
		os.Exit(1)
	}

	amount, err := strconv.Atoi(args[0])
	check(err)

	feedback := ""

	dwf.Client().Donate(int64(amount), feedback)

	check(err)
	fmt.Printf("Biglietto acquistato! Puoi salire sull'autobus")
}
