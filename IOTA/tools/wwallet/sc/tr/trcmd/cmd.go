package trcmd

import (
	"fmt"
	"os"

	"wasp/tools/wwallet/sc/tr"
)

func InitCommands(commands map[string]func([]string)) {
	commands["tr"] = cmd
}

var subcmds = map[string]func([]string){
	"set":    tr.Config.HandleSetCmd,
	"admin":  adminCmd,
	"status": statusCmd,
	"query":  queryCmd,
	"mint":   mintCmd,
}

func cmd(args []string) {
	tr.Config.HandleCmd(args, subcmds)
}

func check(err error) {
	if err != nil {
		fmt.Printf("error: %s\n", err)
		os.Exit(1)
	}
}
