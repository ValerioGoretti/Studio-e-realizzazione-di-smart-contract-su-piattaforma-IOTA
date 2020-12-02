package scclient

import (
	"wasp/packages/registry"
)

func (sc *SCClient) GetBootupData() (*registry.BootupData, error) {
	return sc.WaspClient.GetBootupData(sc.Address)
}
