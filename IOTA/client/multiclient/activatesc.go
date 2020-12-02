package multiclient

import (
	"wasp/client"

	"github.com/iotaledger/goshimmer/dapps/valuetransfers/packages/address"
)

func (m *MultiClient) ActivateSC(addr *address.Address) error {
	return m.Do(func(i int, w *client.WaspClient) error {
		return w.ActivateSC(addr)
	})
}

func (m *MultiClient) DeactivateSC(addr *address.Address) error {
	return m.Do(func(i int, w *client.WaspClient) error {
		return w.DeactivateSC(addr)
	})
}
