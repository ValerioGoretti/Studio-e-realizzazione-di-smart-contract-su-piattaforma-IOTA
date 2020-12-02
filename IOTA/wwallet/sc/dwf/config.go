package dwf

import (
	"wasp/packages/vm/examples/donatewithfeedback/dwfclient"
	"wasp/packages/vm/examples/donatewithfeedback/dwfimpl"
	"wasp/tools/wwallet/sc"
	"wasp/tools/wwallet/wallet"
)

var Config = &sc.Config{
	ShortName:   "dwf",
	Name:        "DonateWithFeedback",
	ProgramHash: dwfimpl.ProgramHash,
}

func Client() *dwfclient.DWFClient {
	return dwfclient.NewClient(Config.MakeClient(wallet.Load().SignatureScheme()))
}
