package sctransaction

import (
	"fmt"
	"io"
	"wasp/packages/hashing"
	"wasp/packages/util"

	"github.com/iotaledger/goshimmer/dapps/valuetransfers/packages/balance"
)

// state block of the SC transaction. Represents SC state update
// previous state block can be determined by the chain transfer of the SC token in the UTXO part of the
// transaction
type StateBlock struct {
	// color of the SC which is updated
	// color contains balance.NEW_COLOR for the origin transaction
	color balance.Color
	// stata index is 0 for the origin transaction
	// consensus maintains incremental sequence of state indexes
	stateIndex uint32
	// timestamp of the transaction. 0 means transaction is not timestamped
	timestamp int64
	// requestId = tx hash + requestId index which originated this state update
	// the list is needed for batches of requests
	// this reference makes requestIds (inputs to state update) immutable part of the state update
	stateHash hashing.HashValue
}

type NewStateBlockParams struct {
	Color      balance.Color
	StateIndex uint32
	StateHash  hashing.HashValue
	Timestamp  int64
}

func NewStateBlock(par NewStateBlockParams) *StateBlock {
	return &StateBlock{
		color:      par.Color,
		stateIndex: par.StateIndex,
		stateHash:  par.StateHash,
		timestamp:  par.Timestamp,
	}
}

func (sb *StateBlock) Clone() *StateBlock {
	if sb == nil {
		return nil
	}
	return NewStateBlock(NewStateBlockParams{
		Color:      sb.color,
		StateIndex: sb.stateIndex,
		StateHash:  sb.stateHash,
		Timestamp:  sb.timestamp,
	})
}

func (sb *StateBlock) Color() balance.Color {
	return sb.color
}

func (sb *StateBlock) StateIndex() uint32 {
	return sb.stateIndex
}

func (sb *StateBlock) Timestamp() int64 {
	return sb.timestamp
}

func (sb *StateBlock) StateHash() hashing.HashValue {
	return sb.stateHash
}

func (sb *StateBlock) WithStateParams(stateIndex uint32, h *hashing.HashValue, ts int64) *StateBlock {
	sb.stateIndex = stateIndex
	sb.stateHash = *h
	sb.timestamp = ts
	return sb
}

// encoding
// important: each block starts with 65 bytes of scid

func (sb *StateBlock) Write(w io.Writer) error {
	if _, err := w.Write(sb.color[:]); err != nil {
		return err
	}
	if err := util.WriteUint32(w, sb.stateIndex); err != nil {
		return err
	}
	if err := util.WriteUint64(w, uint64(sb.timestamp)); err != nil {
		return err
	}
	if err := sb.stateHash.Write(w); err != nil {
		return err
	}
	return nil
}

func (sb *StateBlock) Read(r io.Reader) error {
	if n, err := r.Read(sb.color[:]); err != nil || n != balance.ColorLength {
		return fmt.Errorf("error while reading color: %v", err)
	}
	if err := util.ReadUint32(r, &sb.stateIndex); err != nil {
		return err
	}
	var timestamp uint64
	if err := util.ReadUint64(r, &timestamp); err != nil {
		return err
	}
	sb.timestamp = int64(timestamp)
	if err := sb.stateHash.Read(r); err != nil {
		return err
	}
	return nil
}
