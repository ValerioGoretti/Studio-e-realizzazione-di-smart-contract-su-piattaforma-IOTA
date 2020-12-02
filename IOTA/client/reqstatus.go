package client

import (
	"net/http"

	"wasp/packages/sctransaction"

	"github.com/iotaledger/goshimmer/dapps/valuetransfers/packages/address"
)

func RequestStatusRoute(scAddr string, reqId string) string {
	return "sc/" + scAddr + "/request/" + reqId + "/status"
}

type RequestStatusResponse struct {
	IsProcessed bool
}

func (c *WaspClient) RequestStatus(scAddr *address.Address, requestId *sctransaction.RequestId) (*RequestStatusResponse, error) {
	res := &RequestStatusResponse{}
	if err := c.do(http.MethodGet, RequestStatusRoute(scAddr.String(), requestId.ToBase58()), nil, res); err != nil {
		return nil, err
	}
	return res, nil
}
