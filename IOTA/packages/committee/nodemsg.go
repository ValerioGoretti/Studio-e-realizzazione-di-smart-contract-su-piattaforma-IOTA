package committee

import (
	"wasp/packages/sctransaction"

	"github.com/iotaledger/goshimmer/dapps/valuetransfers/packages/balance"
	valuetransaction "github.com/iotaledger/goshimmer/dapps/valuetransfers/packages/transaction"
)

type StateTransactionMsg struct {
	*sctransaction.Transaction
}

type TransactionInclusionLevelMsg struct {
	TxId  *valuetransaction.ID
	Level byte
}

type BalancesMsg struct {
	Balances map[valuetransaction.ID][]*balance.Balance
}

type RequestMsg struct {
	*sctransaction.Transaction
	Index uint16
}

func (reqMsg *RequestMsg) RequestId() *sctransaction.RequestId {
	ret := sctransaction.NewRequestId(reqMsg.Transaction.ID(), reqMsg.Index)
	return &ret
}

func (reqMsg *RequestMsg) RequestBlock() *sctransaction.RequestBlock {
	return reqMsg.Requests()[reqMsg.Index]
}

func (reqMsg *RequestMsg) Timelock() uint32 {
	return reqMsg.RequestBlock().Timelock()
}
