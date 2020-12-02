package multiclient

import (
	"net/http"
	"time"

	"wasp/client"
	"wasp/packages/util/multicall"
)

type MultiClient struct {
	nodes []*client.WaspClient

	Timeout time.Duration
}

func New(hosts []string, httpClient ...func() http.Client) *MultiClient {
	m := &MultiClient{
		nodes: make([]*client.WaspClient, len(hosts)),
	}
	for i, host := range hosts {
		if len(httpClient) > 0 {
			m.nodes[i] = client.NewWaspClient(host, httpClient[0]())
		} else {
			m.nodes[i] = client.NewWaspClient(host)
		}
	}
	m.Timeout = 30 * time.Second
	return m
}

func (m *MultiClient) Do(f func(int, *client.WaspClient) error) error {
	funs := make([]func() error, len(m.nodes))
	for i := range m.nodes {
		j := i // duplicate variable for closure
		funs[j] = func() error { return f(j, m.nodes[j]) }
	}
	ok, errs := multicall.MultiCall(funs, m.Timeout)
	if !ok {
		return multicall.WrapErrors(errs)
	}
	return nil
}
