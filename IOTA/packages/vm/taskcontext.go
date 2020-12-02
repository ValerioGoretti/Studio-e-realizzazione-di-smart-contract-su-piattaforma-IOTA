package vm

import (
	"bytes"
	"wasp/packages/hashing"
	"wasp/packages/sctransaction"
	"wasp/packages/state"
	"wasp/packages/util"

	"github.com/iotaledger/goshimmer/dapps/valuetransfers/packages/address"
	"github.com/iotaledger/goshimmer/dapps/valuetransfers/packages/balance"
	valuetransaction "github.com/iotaledger/goshimmer/dapps/valuetransfers/packages/transaction"
	"github.com/iotaledger/hive.go/logger"
)

// task context (for batch of requests)
type VMTask struct {
	// inputs (immutable)
	LeaderPeerIndex uint16
	ProgramHash     hashing.HashValue
	Address         address.Address
	Color           balance.Color
	// deterministic source of entropy (pseudorandom, unpredictable for parties)
	Entropy       hashing.HashValue
	Balances      map[valuetransaction.ID][]*balance.Balance
	OwnerAddress  address.Address
	RewardAddress address.Address
	MinimumReward int64
	Requests      []sctransaction.RequestRef
	Timestamp     int64
	VirtualState  state.VirtualState // input immutable
	Log           *logger.Logger
	// call when finished
	OnFinish func(error)
	// outputs
	ResultTransaction *sctransaction.Transaction
	ResultBatch       state.Batch
}

// BatchHash is used to uniquely identify the VM task
func BatchHash(reqids []sctransaction.RequestId, ts int64, leaderIndex uint16) hashing.HashValue {
	var buf bytes.Buffer
	for i := range reqids {
		buf.Write(reqids[i].Bytes())
	}
	_ = util.WriteInt64(&buf, ts)
	_ = util.WriteUint16(&buf, leaderIndex)

	return *hashing.HashData(buf.Bytes())
}
