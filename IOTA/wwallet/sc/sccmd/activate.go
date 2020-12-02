package sccmd

import (
	"fmt"
	"os"

	"wasp/client/multiclient"
	"wasp/tools/wwallet/config"

	"github.com/iotaledger/goshimmer/dapps/valuetransfers/packages/address"
)

func activateCmd(args []string) {
	if len(args) != 2 {
		activateUsage()
	}

	scAddress, err := address.FromBase58(args[0])
	check(err)
	committee := parseIntList(args[1])

	check(multiclient.New(config.CommitteeApi(committee)).ActivateSC(&scAddress))
}

func activateUsage() {
	fmt.Printf("Usage: %s sc activate <sc-address> <committee>\n", os.Args[0])
	fmt.Printf("Example: %s sc activate aBcD...wXyZ '0,1,2,3'\n", os.Args[0])
	os.Exit(1)
}
