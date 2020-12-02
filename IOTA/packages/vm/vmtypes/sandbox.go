package vmtypes

import (
	"wasp/packages/hashing"
	"wasp/packages/kv"
	"wasp/packages/sctransaction"

	"github.com/iotaledger/goshimmer/dapps/valuetransfers/packages/address"
	"github.com/iotaledger/goshimmer/dapps/valuetransfers/packages/balance"
	"github.com/iotaledger/hive.go/logger"
)

// Sandbox is an interface given to the processor to access the VMContext
// and virtual state, transaction builder and request parameters through it.
type Sandbox interface {
	// general function
	IsOriginState() bool
	GetSCAddress() *address.Address
	GetOwnerAddress() *address.Address
	GetTimestamp() int64
	GetEntropy() hashing.HashValue // 32 bytes of deterministic and unpredictably random data

	// Same as panic(), but added as a Sandbox method to emphasize that it's ok to panic from a SC.
	// A panic will be recovered, and Rollback() will be automatically called after.
	Panic(v interface{})

	// clear all updates, restore same context as in the beginning of the VM call
	Rollback()

	// sub interfaces
	// access to the request block
	AccessRequest() RequestAccess
	// base level of virtual state access
	AccessState() kv.MustCodec
	// AccessSCAccount
	AccessSCAccount() AccountAccess
	// Send request
	SendRequest(par NewRequestParams) bool
	// Send request to itself
	SendRequestToSelf(reqCode sctransaction.RequestCode, args kv.Map) bool
	// Send request to itself with timelock for some seconds after the current timestamp
	SendRequestToSelfWithDelay(reqCode sctransaction.RequestCode, args kv.Map, deferForSec uint32) bool
	// for testing
	// Publish "vmmsg" message through Publisher
	Publish(msg string)
	Publishf(format string, args ...interface{})

	GetWaspLog() *logger.Logger
	DumpAccount() string
}

// access to request parameters (arguments)
type RequestAccess interface {
	//request id
	ID() sctransaction.RequestId
	// request code
	Code() sctransaction.RequestCode
	// sender address (exactly 1)
	Sender() address.Address
	// arguments
	Args() kv.RCodec // TODO must return MustCodec
	// number of free minted tokens in the request transaction
	// it is equal to total minted tokens minus number of requests
	NumFreeMintedTokens() int64
}

// access to token operations (txbuilder)
// mint (create new color) is not here on purpose: ColorNew is used for request tokens
type AccountAccess interface {
	// access to total available outputs/balances
	AvailableBalance(col *balance.Color) int64
	MoveTokens(targetAddr *address.Address, col *balance.Color, amount int64) bool
	EraseColor(targetAddr *address.Address, col *balance.Color, amount int64) bool
	// part of the outputs/balances which are coming from the current request transaction
	AvailableBalanceFromRequest(col *balance.Color) int64
	MoveTokensFromRequest(targetAddr *address.Address, col *balance.Color, amount int64) bool
	EraseColorFromRequest(targetAddr *address.Address, col *balance.Color, amount int64) bool
	// send iotas to the smart contract owner
	HarvestFees(amount int64) int64
	HarvestFeesFromRequest(amount int64) bool
}

type NewRequestParams struct {
	TargetAddress *address.Address
	RequestCode   sctransaction.RequestCode
	Timelock      uint32
	Args          kv.Map
	IncludeReward int64
}
