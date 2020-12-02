package nodeconn

import (
	"wasp/plugins/peering"

	"github.com/iotaledger/goshimmer/dapps/valuetransfers/packages/address"
	valuetransaction "github.com/iotaledger/goshimmer/dapps/valuetransfers/packages/transaction"
	"github.com/iotaledger/goshimmer/dapps/waspconn/packages/waspconn"
)

func SendWaspIdToNode() error {
	data, err := waspconn.EncodeMsg(&waspconn.WaspToNodeSetIdMsg{
		Waspid: peering.MyNetworkId(),
	})
	if err != nil {
		return err
	}
	if err := SendDataToNode(data); err != nil {
		return err
	}
	return nil
}

func RequestOutputsFromNode(addr *address.Address) error {
	data, err := waspconn.EncodeMsg(&waspconn.WaspToNodeGetOutputsMsg{
		Address: *addr,
	})
	if err != nil {
		return err
	}
	if err := SendDataToNode(data); err != nil {
		return err
	}
	return nil
}

func RequestConfirmedTransactionFromNode(txid *valuetransaction.ID) error {
	data, err := waspconn.EncodeMsg(&waspconn.WaspToNodeGetConfirmedTransactionMsg{
		TxId: *txid,
	})
	if err != nil {
		return err
	}
	if err := SendDataToNode(data); err != nil {
		return err
	}
	return nil
}

func RequestInclusionLevelFromNode(txid *valuetransaction.ID, addr *address.Address) error {
	log.Debugf("RequestInclusionLevelFromNode. txid %s", txid.String())

	data, err := waspconn.EncodeMsg(&waspconn.WaspToNodeGetTxInclusionLevelMsg{
		TxId:      *txid,
		SCAddress: *addr,
	})
	if err != nil {
		return err
	}
	if err := SendDataToNode(data); err != nil {
		return err
	}
	return nil

}

func PostTransactionToNode(tx *valuetransaction.Transaction, fromSc *address.Address, fromLeader uint16) error {
	data, err := waspconn.EncodeMsg(&waspconn.WaspToNodeTransactionMsg{
		Tx:        tx,
		SCAddress: *fromSc, // just for tracing
		Leader:    fromLeader,
	})
	if err != nil {
		return err
	}
	if err = SendDataToNode(data); err != nil {
		return err
	}
	return nil
}
