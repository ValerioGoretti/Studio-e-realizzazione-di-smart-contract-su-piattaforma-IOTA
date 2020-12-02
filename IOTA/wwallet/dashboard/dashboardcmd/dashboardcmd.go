package dashboardcmd

import (
	"fmt"
	"os"

	"wasp/tools/wwallet/dashboard"
	"wasp/tools/wwallet/sc/dwf"
	"wasp/tools/wwallet/sc/dwf/dwfdashboard"
	"wasp/tools/wwallet/sc/fa"
	"wasp/tools/wwallet/sc/fa/fadashboard"
	"wasp/tools/wwallet/sc/fr"
	"wasp/tools/wwallet/sc/fr/frdashboard"
	"wasp/tools/wwallet/sc/tr"
	"wasp/tools/wwallet/sc/tr/trdashboard"

	"github.com/spf13/pflag"
)

func InitCommands(commands map[string]func([]string), flags *pflag.FlagSet) {
	commands["dashboard"] = cmd
}

func cmd(args []string) {
	listenAddr := ":10000"
	if len(args) > 0 {
		if len(args) != 1 {
			fmt.Printf("Usage: %s dashboard [listen-address]\n", os.Args[0])
			os.Exit(1)
		}
		listenAddr = args[0]
	}

	scs := make([]dashboard.SCDashboard, 0)
	if fr.Config.IsAvailable() {
		scs = append(scs, frdashboard.Dashboard())
		fmt.Printf("FairRoulette: %s\n", fr.Config.Href())
	} else {
		fmt.Println("FairRoulette not available")
	}
	if fa.Config.IsAvailable() {
		scs = append(scs, fadashboard.Dashboard())
		fmt.Printf("FairAuction: %s\n", fa.Config.Href())
	} else {
		fmt.Println("FairAuction not available")
	}
	if tr.Config.IsAvailable() {
		scs = append(scs, trdashboard.Dashboard())
		fmt.Printf("TokenRegistry: %s\n", tr.Config.Href())
	} else {
		fmt.Println("TokenRegistry not available")
	}
	if dwf.Config.IsAvailable() {
		fmt.Printf("DonateWithFeedback: %s\n", dwf.Config.Href())
		scs = append(scs, dwfdashboard.Dashboard())
	} else {
		fmt.Println("DonateWithFeedback not available")
	}

	dashboard.StartServer(listenAddr, scs)
}
