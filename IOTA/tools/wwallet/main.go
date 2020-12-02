package main

import (
	"fmt"
	"os"
	"strings"

	config "wasp/tools/wwallet/config"

	"wasp/tools/wwallet/dashboard/dashboardcmd"
	"wasp/tools/wwallet/program"
	"wasp/tools/wwallet/sc/dwf/dwfcmd"
	"wasp/tools/wwallet/sc/fa/facmd"
	"wasp/tools/wwallet/sc/fr/frcmd"
	"wasp/tools/wwallet/sc/sccmd"
	"wasp/tools/wwallet/sc/tr/trcmd"
	"wasp/tools/wwallet/wallet"

	"github.com/spf13/pflag"
)

func check(err error) {
	if err != nil {
		fmt.Printf("error: %s\n", err)
		os.Exit(1)
	}
}

func usage(commands map[string]func([]string), flags *pflag.FlagSet) {
	cmdNames := make([]string, 0)
	for k := range commands {
		cmdNames = append(cmdNames, k)
	}

	fmt.Printf("Usage: %s [options] [%s]\n", os.Args[0], strings.Join(cmdNames, "|"))
	flags.PrintDefaults()
	os.Exit(1)
}

func main() {
	commands := map[string]func([]string){}
	flags := pflag.NewFlagSet("global flags", pflag.ExitOnError)

	config.InitCommands(commands, flags)
	wallet.InitCommands(commands, flags)
	frcmd.InitCommands(commands)
	facmd.InitCommands(commands)
	trcmd.InitCommands(commands)
	dwfcmd.InitCommands(commands)
	dwfcmd.InitCommandsBuy(commands) //mia funzione di inserimento della funzionalità buy, è una sfaccettatuta del dwf
	dashboardcmd.InitCommands(commands, flags)
	sccmd.InitCommands(commands, flags)
	program.InitCommands(commands, flags)
	check(flags.Parse(os.Args[1:]))

	config.Read()

	if flags.NArg() < 1 {
		usage(commands, flags)
	}

	cmd, ok := commands[flags.Arg(0)]
	if !ok {
		usage(commands, flags)
	}
	cmd(flags.Args()[1:])
}
