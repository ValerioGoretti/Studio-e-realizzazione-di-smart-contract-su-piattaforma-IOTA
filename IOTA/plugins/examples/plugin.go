package examples

import (
	"wasp/packages/hashing"
	"wasp/packages/registry"
	"wasp/packages/vm/examples/donatewithfeedback/dwfimpl"
	"wasp/packages/vm/examples/fairauction"
	"wasp/packages/vm/examples/fairroulette"
	"wasp/packages/vm/examples/inccounter"
	"wasp/packages/vm/examples/logsc"
	"wasp/packages/vm/examples/sc7"
	"wasp/packages/vm/examples/sc8"
	"wasp/packages/vm/examples/sc9"
	"wasp/packages/vm/examples/tokenregistry"
	"wasp/packages/vm/examples/vmnil"
	"wasp/packages/vm/processor"
	"wasp/packages/vm/vmtypes"

	"github.com/iotaledger/hive.go/node"
)

const PluginName = "Examples"

type example struct {
	programHash  string
	getProcessor func() vmtypes.Processor
	name         string
}

func Init() *node.Plugin {
	return node.NewPlugin(PluginName, node.Enabled, configure, run)
}

func configure(ctx *node.Plugin) {
	allExamples := []example{
		{vmnil.ProgramHash, vmnil.GetProcessor, "vmnil"},
		{logsc.ProgramHash, logsc.GetProcessor, "logsc"},
		{inccounter.ProgramHash, inccounter.GetProcessor, "inccounter"},
		{fairroulette.ProgramHash, fairroulette.GetProcessor, "FairRoulette"},
		//{wasmhost.ProgramHash, wasmhost.GetProcessor, "wasmpoc"},
		{fairauction.ProgramHash, fairauction.GetProcessor, "FairAuction"},
		{tokenregistry.ProgramHash, tokenregistry.GetProcessor, "TokenRegistry"},
		{sc7.ProgramHash, sc7.GetProcessor, "sc7"},
		{sc8.ProgramHash, sc8.GetProcessor, "sc8"},
		{sc9.ProgramHash, sc9.GetProcessor, "sc9"},
		{dwfimpl.ProgramHash, dwfimpl.GetProcessor, "DonateWithFeedback"},
	}

	for _, ex := range allExamples {
		hash, _ := hashing.HashValueFromBase58(ex.programHash)
		registry.RegisterBuiltinProgramMetadata(&hash, ex.name+" (Built-in Smart Contract example)")
		processor.RegisterBuiltinProcessor(&hash, ex.getProcessor)
	}
}

func run(ctx *node.Plugin) {
}
