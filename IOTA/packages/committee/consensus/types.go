package consensus

import (
	"sync"
	"time"
	"wasp/packages/state"

	"wasp/packages/committee"
	"wasp/packages/hashing"
	"wasp/packages/registry"
	"wasp/packages/sctransaction"
	"wasp/packages/tcrypto"
	"wasp/packages/tcrypto/tbdn"
	"wasp/packages/util"
	"wasp/packages/vm/vmconst"

	"github.com/iotaledger/goshimmer/dapps/valuetransfers/packages/address"
	"github.com/iotaledger/goshimmer/dapps/valuetransfers/packages/balance"
	valuetransaction "github.com/iotaledger/goshimmer/dapps/valuetransfers/packages/transaction"
	"github.com/iotaledger/hive.go/logger"
)

type operator struct {
	committee committee.Committee
	dkshare   *tcrypto.DKShare
	//currentSCState
	currentSCState state.VirtualState
	stateTx        *sctransaction.Transaction
	balances       map[valuetransaction.ID][]*balance.Balance

	// consensus stage
	consensusStage         int
	consensusStageDeadline time.Time
	//
	requestBalancesDeadline time.Time
	processorReady          bool

	// notifications with future currentSCState indices
	notificationsBacklog []*committee.NotifyReqMsg

	// backlog of requests with all information
	requests map[sctransaction.RequestId]*request

	peerPermutation *util.Permutation16

	leaderStatus            *leaderStatus
	sentResultToLeaderIndex uint16
	sentResultToLeader      *sctransaction.Transaction

	postedResultTxid       *valuetransaction.ID
	nextPullInclusionLevel time.Time // if postedResultTxid != nil

	log *logger.Logger

	// data for concurrent access, from APIs mostly
	concurrentAccessMutex sync.RWMutex
	requestIdsProtected   map[sctransaction.RequestId]bool
}

type leaderStatus struct {
	reqs          []*request
	batch         state.Batch
	batchHash     hashing.HashValue
	timestamp     int64
	balances      map[valuetransaction.ID][]*balance.Balance
	resultTx      *sctransaction.Transaction
	finalized     bool
	signedResults []*signedResult
}

type signedResult struct {
	essenceHash hashing.HashValue
	sigShare    tbdn.SigShare
}

// keeps stateTx of the request
type request struct {
	// id of the hash of request tx id and request block index
	reqId sctransaction.RequestId
	// from request message. nil if request message wasn't received yet
	reqTx *sctransaction.Transaction
	// time when request message was received by the operator
	whenMsgReceived time.Time
	// notification vector for the current currentSCState
	notifications []bool

	log *logger.Logger
}

func NewOperator(committee committee.Committee, dkshare *tcrypto.DKShare, log *logger.Logger) *operator {
	defer committee.SetReadyConsensus()

	ret := &operator{
		committee:           committee,
		dkshare:             dkshare,
		requests:            make(map[sctransaction.RequestId]*request),
		requestIdsProtected: make(map[sctransaction.RequestId]bool),
		peerPermutation:     util.NewPermutation16(committee.Size(), nil),
		log:                 log.Named("c"),
	}
	ret.setNextConsensusStage(consensusStageNoSync)
	return ret
}

func (op *operator) peerIndex() uint16 {
	return op.dkshare.Index
}

func (op *operator) quorum() uint16 {
	return op.dkshare.T
}

func (op *operator) size() uint16 {
	return op.dkshare.N
}

func (op *operator) stateIndex() (uint32, bool) {
	if op.currentSCState == nil {
		return 0, false
	}
	return op.currentSCState.StateIndex(), true
}

func (op *operator) mustStateIndex() uint32 {
	ret, ok := op.stateIndex()
	if !ok {
		panic("mustStateIndex")
	}
	return ret
}

func (op *operator) getProgramHash() (*hashing.HashValue, bool) {
	if op.currentSCState == nil {
		return nil, false
	}
	h, ok, err := op.currentSCState.Variables().Codec().GetHashValue(vmconst.VarNameProgramHash)
	if !ok || err != nil {
		return nil, false
	}
	return h, true
}

func (op *operator) getRewardAddress() address.Address {
	return registry.GetRewardAddress(op.committee.Address())
}

func (op *operator) getMinimumReward() int64 {
	if _, ok := op.stateIndex(); !ok {
		return 0
	}
	vt, ok, err := op.currentSCState.Variables().Codec().GetInt64(vmconst.VarNameMinimumReward)
	if err != nil {
		panic(err)
	}
	if !ok {
		return 0
	}
	return vt
}
