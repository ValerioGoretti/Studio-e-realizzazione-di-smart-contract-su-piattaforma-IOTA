package sctransaction

import (
	"errors"
	"io"

	"wasp/packages/txutil"

	"github.com/iotaledger/goshimmer/dapps/valuetransfers/packages/address"
	"github.com/iotaledger/goshimmer/dapps/valuetransfers/packages/balance"
)

func ReadRequestId(r io.Reader, reqid *RequestId) error {
	n, err := r.Read(reqid[:])
	if err != nil {
		return err
	}
	if n != RequestIdSize {
		return errors.New("error while reading request id")
	}
	return nil
}

func OutputValueOfColor(tx *Transaction, addr address.Address, color balance.Color) int64 {
	bals, ok := tx.Outputs().Get(addr)
	if !ok {
		return 0
	}

	return txutil.BalanceOfColor(bals.([]*balance.Balance), color)
}
