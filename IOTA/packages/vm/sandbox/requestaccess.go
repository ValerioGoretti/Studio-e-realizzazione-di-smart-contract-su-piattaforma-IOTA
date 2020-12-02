package sandbox

import (
	"wasp/packages/kv"
	"wasp/packages/sctransaction"

	"github.com/iotaledger/goshimmer/dapps/valuetransfers/packages/address"
)

// access to the request block
type requestWrapper struct {
	ref *sctransaction.RequestRef
}

func (r *requestWrapper) ID() sctransaction.RequestId {
	return *r.ref.RequestId()
}

func (r *requestWrapper) Code() sctransaction.RequestCode {
	return r.ref.RequestBlock().RequestCode()
}

func (r *requestWrapper) Args() kv.RCodec {
	return r.ref.RequestBlock().Args()
}

// addresses of request transaction inputs
func (r *requestWrapper) Sender() address.Address {
	return *r.ref.Tx.MustProperties().Sender()
}

//MintedBalances return total minted tokens minus number of
func (r *requestWrapper) NumFreeMintedTokens() int64 {
	return r.ref.Tx.MustProperties().NumFreeMintedTokens()
}
