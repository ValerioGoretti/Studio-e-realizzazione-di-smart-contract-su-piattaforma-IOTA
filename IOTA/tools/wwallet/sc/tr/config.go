package tr

import (
	"wasp/packages/vm/examples/tokenregistry"
	"wasp/packages/vm/examples/tokenregistry/trclient"
	"wasp/tools/wwallet/sc"
	"wasp/tools/wwallet/wallet"
)

var Config = &sc.Config{
	ShortName:   "tr",
	Name:        "TokenRegistry",
	ProgramHash: tokenregistry.ProgramHash,
}

func Client() *trclient.TokenRegistryClient {
	return trclient.NewClient(Config.MakeClient(wallet.Load().SignatureScheme()))
}
