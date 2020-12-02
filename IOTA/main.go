package main

import (
	"wasp/packages/parameters"
	"wasp/packages/registry"
	"wasp/plugins/banner"
	"wasp/plugins/cli"
	"wasp/plugins/committees"
	"wasp/plugins/config"
	"wasp/plugins/dashboard"
	"wasp/plugins/database"
	"wasp/plugins/dispatcher"
	"wasp/plugins/examples"
	"wasp/plugins/gracefulshutdown"
	"wasp/plugins/logger"
	"wasp/plugins/nodeconn"
	"wasp/plugins/peering"
	"wasp/plugins/publisher"
	"wasp/plugins/testplugins/nodeping"
	"wasp/plugins/testplugins/roundtrip"
	"wasp/plugins/wasmtimevm"
	"wasp/plugins/webapi"

	"github.com/iotaledger/hive.go/node"
)

func main() {
	registry.InitFlags()
	parameters.InitFlags()

	plugins := node.Plugins(
		banner.Init(),
		config.Init(),
		logger.Init(),
		gracefulshutdown.Init(),
		webapi.Init(),
		cli.Init(),
		database.Init(),
		peering.Init(),
		nodeconn.Init(),
		dispatcher.Init(),
		committees.Init(),
		wasmtimevm.Init(),
		publisher.Init(),
		dashboard.Init(),
		examples.Init(),
	)

	testPlugins := node.Plugins(
		roundtrip.Init(),
		nodeping.Init(),
	)

	node.Run(
		plugins,
		testPlugins,
	)
}
