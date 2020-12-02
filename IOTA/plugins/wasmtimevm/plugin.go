// Wasp can have several VM types. Each of them can be represented by separate plugin
// Plugin name serves as a VM type during dynamic loading of the binary.
// VM plugins can be enabled/disabled in the configuration of the node instance
// wasmtimevm plugin statically links VM implemented with Wasmtime to Wasp
// be registering wasmhost.GetProcessor as function
package wasmtimevm

import (
	"wasp/packages/vm/vmtypes"
	"wasp/packages/vm/wasmhost"

	"github.com/iotaledger/hive.go/logger"
	"github.com/iotaledger/hive.go/node"
)

// PluginName is the name of the plugin.
const PluginName = "wasmtimevm"

var log *logger.Logger

func Init() *node.Plugin {
	return node.NewPlugin(PluginName, node.Enabled, configure, run)
}

func configure(_ *node.Plugin) {
	log = logger.NewLogger(PluginName)

	// register VM type(s)
	err := vmtypes.RegisterVMType(PluginName, wasmhost.GetProcessor)
	if err != nil {
		log.Panicf("%v: %v", PluginName, err)
	}
	log.Infof("registered VM type: '%s'", PluginName)
}

func run(_ *node.Plugin) {
}
