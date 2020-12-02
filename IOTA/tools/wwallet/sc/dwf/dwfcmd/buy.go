package dwfcmd

import (
	"fmt"
	"os"
	"strconv"

	"wasp/tools/wwallet/sc/dwf"
)

func buyCmd(args []string) {
	if len(args) != 1 {
		fmt.Printf("Usage: %s buy iota <amount>\n", os.Args[0])
		os.Exit(1)
	}

	amount, err := strconv.Atoi(args[0])
	check(err)

	_, err = dwf.Client().Buy(int64(amount))
	check(err)
}
