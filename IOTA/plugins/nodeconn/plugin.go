package nodeconn

import (
	"sync"
	"time"
	"wasp/packages/parameters"

	"github.com/iotaledger/goshimmer/dapps/valuetransfers/packages/address"
	"github.com/iotaledger/goshimmer/dapps/valuetransfers/packages/balance"
	"github.com/iotaledger/hive.go/daemon"
	"github.com/iotaledger/hive.go/logger"
	"github.com/iotaledger/hive.go/netutil/buffconn"
	"github.com/iotaledger/hive.go/node"
)

// PluginName is the name of the NodeConn plugin.
const PluginName = "NodeConn"

var (
	log *logger.Logger

	bconn             *buffconn.BufferedConnection
	bconnMutex        = &sync.Mutex{}
	subscriptions     = make(map[address.Address]balance.Color)
	subscriptionsSent bool
)

func Init() *node.Plugin {
	return node.NewPlugin(PluginName, node.Enabled, configure, run)
}

func configure(_ *node.Plugin) {
	log = logger.NewLogger(PluginName)
}

func run(_ *node.Plugin) {
	err := daemon.BackgroundWorker(PluginName, func(shutdownSignal <-chan struct{}) {
		go nodeConnect()
		go keepSendingSubscriptionIfNeeded(shutdownSignal)
		go keepSendingSubscriptionForced(shutdownSignal)

		<-shutdownSignal

		log.Info("Stopping node connection..")
		go func() {
			bconnMutex.Lock()
			defer bconnMutex.Unlock()

			if bconn != nil {
				log.Infof("Closing connection with node..")
				_ = bconn.Close()
				log.Infof("Closing connection with node.. Done")
			}
		}()

	}, parameters.PriorityNodeConnection)
	if err != nil {
		log.Errorf("failed to start NodeConn worker")
	}
}

// checking if need to be sent every second
func keepSendingSubscriptionIfNeeded(shutdownSignal <-chan struct{}) {
	for {
		select {
		case <-shutdownSignal:
			return
		case <-time.After(1 * time.Second):
			sendSubscriptions(false)
		}
	}
}

// will be sending subscriptions every minute to pull backlog
// needed in case node is not synced
func keepSendingSubscriptionForced(shutdownSignal <-chan struct{}) {
	for {
		select {
		case <-shutdownSignal:
			return
		case <-time.After(1 * time.Minute):
			sendSubscriptions(true)
		}
	}
}
