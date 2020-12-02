package frclient

import (
	"fmt"
	"time"

	"wasp/client/scclient"
	"wasp/client/statequery"
	"wasp/packages/sctransaction"
	"wasp/packages/util"
	"wasp/packages/vm/examples/fairroulette"

	"github.com/iotaledger/goshimmer/dapps/valuetransfers/packages/address"
	"github.com/iotaledger/goshimmer/dapps/valuetransfers/packages/balance"
)

type FairRouletteClient struct {
	*scclient.SCClient
}

func NewClient(scClient *scclient.SCClient) *FairRouletteClient {
	return &FairRouletteClient{scClient}
}

type Status struct {
	*scclient.SCStatus

	CurrentBetsAmount uint16
	CurrentBets       []*fairroulette.BetInfo

	LockedBetsAmount uint16
	LockedBets       []*fairroulette.BetInfo

	LastWinningColor int64

	PlayPeriodSeconds int64

	NextPlayTimestamp time.Time

	PlayerStats map[address.Address]*fairroulette.PlayerStats

	WinsPerColor []uint32
}

func (s *Status) NextPlayIn() string {
	diff := s.NextPlayTimestamp.Sub(s.FetchedAt)
	// round to the second
	diff -= diff % time.Second
	if diff < 0 {
		return "unknown"
	}
	return diff.String()
}

func (frc *FairRouletteClient) FetchStatus() (*Status, error) {
	scStatus, results, err := frc.FetchSCStatus(func(query *statequery.Request) {
		query.AddArray(fairroulette.StateVarBets, 0, 100)
		query.AddArray(fairroulette.StateVarLockedBets, 0, 100)
		query.AddScalar(fairroulette.StateVarLastWinningColor)
		query.AddScalar(fairroulette.ReqVarPlayPeriodSec)
		query.AddScalar(fairroulette.StateVarNextPlayTimestamp)
		query.AddDictionary(fairroulette.StateVarPlayerStats, 100)
		query.AddArray(fairroulette.StateArrayWinsPerColor, 0, fairroulette.NumColors)
	})
	if err != nil {
		return nil, err
	}

	status := &Status{SCStatus: scStatus}

	lastWinningColor, _ := results.Get(fairroulette.StateVarLastWinningColor).MustInt64()
	status.LastWinningColor = lastWinningColor

	playPeriodSeconds, _ := results.Get(fairroulette.ReqVarPlayPeriodSec).MustInt64()
	status.PlayPeriodSeconds = playPeriodSeconds
	if status.PlayPeriodSeconds == 0 {
		status.PlayPeriodSeconds = fairroulette.DefaultPlaySecondsAfterFirstBet
	}

	nextPlayTimestamp, _ := results.Get(fairroulette.StateVarNextPlayTimestamp).MustInt64()
	status.NextPlayTimestamp = time.Unix(0, nextPlayTimestamp).UTC()

	status.PlayerStats, err = decodePlayerStats(results.Get(fairroulette.StateVarPlayerStats).MustDictionaryResult())
	if err != nil {
		return nil, err
	}

	status.WinsPerColor, err = decodeWinsPerColor(results.Get(fairroulette.StateArrayWinsPerColor).MustArrayResult())
	if err != nil {
		return nil, err
	}

	status.CurrentBetsAmount, status.CurrentBets, err = decodeBets(results.Get(fairroulette.StateVarBets).MustArrayResult())
	if err != nil {
		return nil, err
	}

	status.LockedBetsAmount, status.LockedBets, err = decodeBets(results.Get(fairroulette.StateVarLockedBets).MustArrayResult())
	if err != nil {
		return nil, err
	}

	return status, nil
}

func decodeBets(result *statequery.ArrayResult) (uint16, []*fairroulette.BetInfo, error) {
	size := result.Len
	bets := make([]*fairroulette.BetInfo, 0)
	for _, b := range result.Values {
		bet, err := fairroulette.DecodeBetInfo(b)
		if err != nil {
			return 0, nil, err
		}
		bets = append(bets, bet)
	}
	return size, bets, nil
}

func decodeWinsPerColor(result *statequery.ArrayResult) ([]uint32, error) {
	ret := make([]uint32, 0)
	for _, b := range result.Values {
		var n uint32
		if b != nil {
			n = util.Uint32From4Bytes(b)
		}
		ret = append(ret, n)
	}
	return ret, nil
}

func decodePlayerStats(result *statequery.DictResult) (map[address.Address]*fairroulette.PlayerStats, error) {
	playerStats := make(map[address.Address]*fairroulette.PlayerStats)
	for _, e := range result.Entries {
		if len(e.Key) != address.Length {
			return nil, fmt.Errorf("not an address: %v", e.Key)
		}
		addr, _, err := address.FromBytes(e.Key)
		if err != nil {
			return nil, err
		}
		ps, err := fairroulette.DecodePlayerStats(e.Value)
		if err != nil {
			return nil, err
		}
		playerStats[addr] = ps
	}
	return playerStats, nil
}

func (frc *FairRouletteClient) Bet(color int, amount int) (*sctransaction.Transaction, error) {
	return frc.PostRequest(
		fairroulette.RequestPlaceBet,
		nil,
		map[balance.Color]int64{balance.ColorIOTA: int64(amount)},
		map[string]interface{}{fairroulette.ReqVarColor: int64(color)},
	)
}

func (frc *FairRouletteClient) SetPeriod(seconds int) (*sctransaction.Transaction, error) {
	return frc.PostRequest(
		fairroulette.RequestSetPlayPeriod,
		nil,
		nil,
		map[string]interface{}{fairroulette.ReqVarPlayPeriodSec: int64(seconds)},
	)
}
