package frcmd

import (
	"wasp/tools/wwallet/sc/fr"
)

func InitCommands(commands map[string]func([]string)) {
	commands["fr"] = cmd
}

var subcmds = map[string]func([]string){
	"set":    fr.Config.HandleSetCmd,
	"admin":  adminCmd,
	"status": statusCmd,
	"bet":    betCmd,
}

func cmd(args []string) {
	fr.Config.HandleCmd(args, subcmds)
}
