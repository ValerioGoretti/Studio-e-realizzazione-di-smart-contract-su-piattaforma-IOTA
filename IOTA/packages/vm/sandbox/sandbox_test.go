package sandbox

import (
	"testing"

	"wasp/packages/state"

	"github.com/iotaledger/goshimmer/dapps/valuetransfers/packages/address"
	"github.com/iotaledger/hive.go/kvstore/mapdb"
	"github.com/stretchr/testify/assert"
)

func TestSetThenGet(t *testing.T) {
	db := mapdb.NewMapDB()
	addr := address.Random()
	s := stateWrapper{
		virtualState: state.NewVirtualState(db, &addr),
		stateUpdate:  state.NewStateUpdate(nil),
	}

	s.Set("x", []byte{1})
	v, err := s.Get("x")

	assert.NoError(t, err)
	assert.Equal(t, []byte{1}, v)

	s.Del("x")
	v, err = s.Get("x")

	assert.NoError(t, err)
	assert.Nil(t, v)

	s.Set("x", []byte{2})
	s.Set("x", []byte{3})
	v, err = s.Get("x")

	assert.NoError(t, err)
	assert.Equal(t, []byte{3}, v)
}
