package vm

import (
	"wasp/packages/hashing"
	"wasp/packages/sctransaction"
	"wasp/packages/sctransaction/txbuilder"
	"wasp/packages/state"

	"github.com/iotaledger/goshimmer/dapps/valuetransfers/packages/address"
	"github.com/iotaledger/hive.go/logger"
)

// context of one VM call (for one request)
type VMContext struct {
	// invariant through the batch
	// address of the smart contract
	Address address.Address
	// programHash
	ProgramHash hashing.HashValue
	// owner address
	OwnerAddress address.Address
	// reward address
	RewardAddress address.Address
	// minimum reward
	MinimumReward int64
	// are node rewards enabled
	NodeRewardsEnabled bool
	// deterministic source of entropy. Equal the the hash of the previous
	Entropy hashing.HashValue
	// tx builder to build the final transaction
	TxBuilder *txbuilder.Builder
	// timestamp of the batch
	Timestamp int64
	// initial state of the call
	VirtualState state.VirtualState
	// set for each call
	RequestRef sctransaction.RequestRef
	// IsEmpty state update upon call, result of the call.
	StateUpdate state.StateUpdate
	// log
	Log *logger.Logger
}
