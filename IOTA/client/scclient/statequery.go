package scclient

import "wasp/client/statequery"

func (sc *SCClient) StateQuery(query *statequery.Request) (*statequery.Results, error) {
	return sc.WaspClient.StateQuery(sc.Address, query)
}
