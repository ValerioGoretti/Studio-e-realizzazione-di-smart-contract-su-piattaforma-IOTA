package fr

import (
	"wasp/packages/vm/examples/fairroulette"
	"wasp/packages/vm/examples/fairroulette/frclient"
	"wasp/tools/wwallet/sc"
	"wasp/tools/wwallet/wallet"
)

var Config = &sc.Config{
	ShortName:   "fr",
	Name:        "FairRoulette",
	ProgramHash: fairroulette.ProgramHash,
}

func Client() *frclient.FairRouletteClient {
	return frclient.NewClient(Config.MakeClient(wallet.Load().SignatureScheme()))
}
