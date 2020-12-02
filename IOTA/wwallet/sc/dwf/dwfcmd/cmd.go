package dwfcmd

import (
	"fmt"
	"os"

	"wasp/tools/wwallet/sc/dwf"
)

func InitCommandsBuy(commands map[string]func([]string)) {
	commands["buy"] = cmd
}

func InitCommands(commands map[string]func([]string)) {
	commands["dwf"] = cmd
}

var subcmds = map[string]func([]string){
	"set":          dwf.Config.HandleSetCmd,
	"iota":         buyCmd,          //Funzione per comprare iota passando un intero (numero euro caricati) (NON FUNZIONA, fino alla fine ma poi non accredita)
	"payIota":      payIotaCmd,      //Funzione per pagare il biglietto con gli IOTA (FUNZIONA)
	"payPl":        payPlCmd,        //Funzione per pagare biglietti con la plastica (NON FUNZIONA)
	"uploadCredit": uploadCreditCmd, //Funzione per ricevere  biglietti con la plastica (con send request)
	"admin":        adminCmd,
	"donate":       donateCmd,
	"withdraw":     withdrawCmd,
	"status":       statusCmd,
}

func cmd(args []string) {
	dwf.Config.HandleCmd(args, subcmds)
}

func check(err error) {
	if err != nil {
		fmt.Printf("error: %s\n", err)
		os.Exit(1)
	}
}
