package client

import (
	"net/http"

	"wasp/packages/kv"

	"github.com/iotaledger/goshimmer/dapps/valuetransfers/packages/address"
)

type SCStateDump struct {
	Index     uint32            `json:"index"`
	Variables map[kv.Key][]byte `json:"variables"`
}

func DumpSCStateRoute(address string) string {
	return "sc/" + address + "/dumpstate"
}

func (c *WaspClient) DumpSCState(addr *address.Address) (*SCStateDump, error) {
	res := &SCStateDump{}
	if err := c.do(http.MethodGet, AdminRoutePrefix+"/"+DumpSCStateRoute(addr.String()), nil, res); err != nil {
		return nil, err
	}
	return res, nil
}
