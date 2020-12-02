package fa

import (
	"wasp/packages/vm/examples/fairauction"
	"wasp/packages/vm/examples/fairauction/faclient"
	"wasp/tools/wwallet/sc"
	"wasp/tools/wwallet/wallet"
)

var Config = &sc.Config{
	ShortName:   "fa",
	Name:        "FairAuction",
	ProgramHash: fairauction.ProgramHash,
}

func Client() *faclient.FairAuctionClient {
	return faclient.NewClient(Config.MakeClient(wallet.Load().SignatureScheme()))
}
