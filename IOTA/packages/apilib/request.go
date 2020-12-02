package apilib

import (
	"errors"
	"fmt"
	"strconv"
	"time"
	"wasp/packages/hashing"
	"wasp/packages/kv"
	"wasp/packages/nodeclient"
	"wasp/packages/sctransaction"
	"wasp/packages/sctransaction/txbuilder"
	"wasp/packages/subscribe"

	"github.com/iotaledger/goshimmer/dapps/valuetransfers/packages/address"
	"github.com/iotaledger/goshimmer/dapps/valuetransfers/packages/address/signaturescheme"
	"github.com/iotaledger/goshimmer/dapps/valuetransfers/packages/balance"
)

type RequestBlockParams struct {
	TargetSCAddress *address.Address
	RequestCode     sctransaction.RequestCode
	Timelock        uint32
	Transfer        map[balance.Color]int64 // should not not include request token. It is added automatically
	Vars            map[string]interface{}  ` `
}

type CreateRequestTransactionParams struct {
	NodeClient          nodeclient.NodeClient
	SenderSigScheme     signaturescheme.SignatureScheme
	BlockParams         []RequestBlockParams
	Mint                map[address.Address]int64
	Post                bool
	WaitForConfirmation bool
	WaitForCompletion   bool
	PublisherHosts      []string
	PublisherQuorum     int
	Timeout             time.Duration
}

func CreateRequestTransaction(par CreateRequestTransactionParams) (*sctransaction.Transaction, error) {

	senderAddr := par.SenderSigScheme.Address()
	allOuts, err := par.NodeClient.GetConfirmedAccountOutputs(&senderAddr)
	fmt.Println("senderAddr -> ", senderAddr)
	fmt.Println("allOuts -> ", allOuts)
	if err != nil {
		//1
		return nil, fmt.Errorf("can't get outputs from the node: %v", err)
	}

	txb, err := txbuilder.NewFromOutputBalances(allOuts)
	if err != nil {
		//2
		return nil, err
	}

	for targetAddress, amount := range par.Mint {
		// TODO: check that targetAddress is not any target address in request blocks
		err = txb.MintColor(targetAddress, balance.ColorIOTA, amount)
		if err != nil {
			//3
			return nil, err
		}
	}

	for _, blockPar := range par.BlockParams {
		reqBlk := sctransaction.NewRequestBlock(*blockPar.TargetSCAddress, blockPar.RequestCode).
			WithTimelock(blockPar.Timelock)

		args := convertArgs(blockPar.Vars)
		if args == nil {
			//4
			return nil, errors.New("wrong arguments")
		}
		reqBlk.SetArgs(args)

		err = txb.AddRequestBlockWithTransfer(reqBlk, blockPar.TargetSCAddress, blockPar.Transfer)
		if err != nil {
			//5
			return nil, err
		}
	}

	tx, err := txb.Build(false)

	//dump := txb.Dump()

	if err != nil {
		//6
		return nil, err
	}
	tx.Sign(par.SenderSigScheme)

	// semantic check just in case
	if _, err := tx.Properties(); err != nil {
		//7
		return nil, err
	}
	//fmt.Printf("$$$$ dumping builder for %s\n%s\n", tx.ID().String(), dump)

	if !par.Post {
		//8
		return tx, nil
	}
	if !par.WaitForConfirmation {
		if err = par.NodeClient.PostTransaction(tx.Transaction); err != nil {
			//9
			return nil, err
		}
		//10
		return tx, nil
	}
	var subs *subscribe.Subscription
	if par.WaitForCompletion {
		// post and wait for completion
		subs, err = subscribe.SubscribeMulti(par.PublisherHosts, "request_out", par.PublisherQuorum)
		if err != nil {
			//11
			return nil, err
		}
		defer subs.Close()
	}

	err = par.NodeClient.PostAndWaitForConfirmation(tx.Transaction)
	if err != nil {
		//12
		return nil, err
	}
	if par.WaitForCompletion {
		patterns := make([][]string, len(par.BlockParams))
		for i := range patterns {
			patterns[i] = []string{"request_out", par.BlockParams[i].TargetSCAddress.String(), tx.ID().String(), strconv.Itoa(i)}
		}
		if !subs.WaitForPatterns(patterns, par.Timeout, par.PublisherQuorum) {
			//13
			return nil, fmt.Errorf("didn't receive completion message after %v", par.Timeout)
		}
	}
	//14
	return tx, nil
}

func convertArgs(vars map[string]interface{}) kv.Map {
	args := kv.NewMap()
	codec := args.Codec()
	for k, v := range vars {
		key := kv.Key(k)
		switch vt := v.(type) {
		case int:
			codec.SetInt64(key, int64(vt))
		case byte:
			codec.SetInt64(key, int64(vt))
		case int16:
			codec.SetInt64(key, int64(vt))
		case int32:
			codec.SetInt64(key, int64(vt))
		case int64:
			codec.SetInt64(key, vt)
		case uint16:
			codec.SetInt64(key, int64(vt))
		case uint32:
			codec.SetInt64(key, int64(vt))
		case uint64:
			codec.SetInt64(key, int64(vt))
		case string:
			codec.SetString(key, vt)
		case []byte:
			codec.Set(key, vt)
		case *hashing.HashValue:
			args.Codec().SetHashValue(key, vt)
		case *address.Address:
			args.Codec().Set(key, vt.Bytes())
		case *balance.Color:
			args.Codec().Set(key, vt.Bytes())
		default:
			return nil
		}
	}
	return args
}
