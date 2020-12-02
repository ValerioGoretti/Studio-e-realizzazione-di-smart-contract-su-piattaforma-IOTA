package wasptest2

import (
	"fmt"
	"testing"
	"time"

	"github.com/iotaledger/goshimmer/dapps/valuetransfers/packages/balance"
	"github.com/iotaledger/wasp/client/scclient"
	"github.com/iotaledger/wasp/packages/kv"
	"github.com/iotaledger/wasp/packages/testutil"
	"github.com/iotaledger/wasp/packages/vm/examples/donatewithfeedback"
	"github.com/iotaledger/wasp/packages/vm/examples/donatewithfeedback/dwfclient"
	"github.com/iotaledger/wasp/packages/vm/examples/donatewithfeedback/dwfimpl"
	"github.com/iotaledger/wasp/packages/vm/vmconst"
)

func TestDeployDWF(t *testing.T) {
	wasps := setup(t, "TestDeployDWF")

	err := requestFunds(wasps, scOwnerAddr, "sc owner")
	check(err, t)

	scAddr, scColor, err := startSmartContract(wasps, dwfimpl.ProgramHash, dwfimpl.Description)
	checkSuccess(err, t, "smart contract has been created and activated")

	if !wasps.VerifyAddressBalances(scOwnerAddr, testutil.RequestFundsAmount-1, map[balance.Color]int64{
		balance.ColorIOTA: testutil.RequestFundsAmount - 1,
	}, "sc owner in the end") {
		t.Fail()
		return
	}

	if !wasps.VerifyAddressBalances(scAddr, 1, map[balance.Color]int64{
		*scColor: 1,
	}, "sc in the end") {
		t.Fail()
		return
	}

	if !wasps.VerifySCStateVariables2(scAddr, map[kv.Key]interface{}{
		vmconst.VarNameOwnerAddress: scOwnerAddr[:],
		vmconst.VarNameProgramHash:  programHash[:],
		vmconst.VarNameDescription:  dwfimpl.Description,
	}) {
		t.Fail()
	}
}

const numDonations = 5

func TestDWFDonateNTimes(t *testing.T) {
	wasps := setup(t, "TestDeployDWF")

	err := requestFunds(wasps, scOwnerAddr, "sc owner")
	check(err, t)

	donor := wallet.WithIndex(1)
	donorAddr := donor.Address()
	err = requestFunds(wasps, donorAddr, "donor")
	check(err, t)

	scAddr, scColor, err := startSmartContract(wasps, dwfimpl.ProgramHash, dwfimpl.Description)
	checkSuccess(err, t, "smart contract has been created and activated")

	dwfClient := dwfclient.NewClient(scclient.New(
		wasps.NodeClient,
		wasps.WaspClient(0),
		scAddr,
		donor.SigScheme(),
	))

	for i := 0; i < numDonations; i++ {
		_, err = dwfClient.Donate(42, fmt.Sprintf("Donation #%d: well done, I give you 42 iotas", i))
		check(err, t)
		time.Sleep(1 * time.Second)
		checkSuccess(err, t, "donated 42")
	}

	status, err := dwfClient.FetchStatus()
	text := ""
	if err == nil {
		text = fmt.Sprintf("[test] Status fetched succesfully: num rec: %d, "+
			"total donations: %d, max donation: %d, last donation: %v, num rec returned: %d\n",
			status.NumRecords,
			status.TotalDonations,
			status.MaxDonation,
			status.LastDonated,
			len(status.LastRecordsDesc),
		)
		for i, di := range status.LastRecordsDesc {
			text += fmt.Sprintf("           %d: ts: %s, amount: %d, fb: '%s', donor: %s, err:%v\n",
				i,
				di.When.Format("2006-01-02 15:04:05"),
				di.Amount,
				di.Feedback,
				di.Sender.String(),
				di.Error,
			)
		}
	}
	checkSuccess(err, t, text)

	if !wasps.VerifyAddressBalances(scOwnerAddr, testutil.RequestFundsAmount-1, map[balance.Color]int64{
		balance.ColorIOTA: testutil.RequestFundsAmount - 1,
	}, "sc owner in the end") {
		t.Fail()
	}

	if !wasps.VerifyAddressBalances(scAddr, 1+42*numDonations, map[balance.Color]int64{
		*scColor:          1,
		balance.ColorIOTA: 42 * numDonations,
	}, "sc in the end") {
		t.Fail()
	}

	if !wasps.VerifyAddressBalances(donorAddr, testutil.RequestFundsAmount-42*numDonations, map[balance.Color]int64{
		balance.ColorIOTA: testutil.RequestFundsAmount - 42*numDonations,
	}, "donor in the end") {
		t.Fail()
	}

	if !wasps.VerifySCStateVariables2(scAddr, map[kv.Key]interface{}{
		vmconst.VarNameOwnerAddress:               scOwnerAddr[:],
		vmconst.VarNameProgramHash:                programHash[:],
		vmconst.VarNameDescription:                dwfimpl.Description,
		donatewithfeedback.VarStateMaxDonation:    42,
		donatewithfeedback.VarStateTotalDonations: 42 * numDonations,
	}) {
		t.Fail()
	}
}

func TestDWFDonateWithdrawAuthorised(t *testing.T) {
	wasps := setup(t, "TestDeployDWF")

	err := requestFunds(wasps, scOwnerAddr, "sc owner")
	check(err, t)

	donor := wallet.WithIndex(1)
	donorAddr := donor.Address()
	err = requestFunds(wasps, donorAddr, "donor")
	check(err, t)

	scAddr, scColor, err := startSmartContract(wasps, dwfimpl.ProgramHash, dwfimpl.Description)
	checkSuccess(err, t, "smart contract has been created and activated")

	dwfDonorClient := dwfclient.NewClient(scclient.New(
		wasps.NodeClient,
		wasps.WaspClient(0),
		scAddr,
		donor.SigScheme(),
		15*time.Second,
	))
	_, err = dwfDonorClient.Donate(42, "well done, I give you 42i")
	check(err, t)
	checkSuccess(err, t, "donated 42")

	if !wasps.VerifyAddressBalances(scOwnerAddr, testutil.RequestFundsAmount-1, map[balance.Color]int64{
		balance.ColorIOTA: testutil.RequestFundsAmount - 1,
	}, "sc owner after donation") {
		t.Fail()
	}

	if !wasps.VerifyAddressBalances(scAddr, 1+42, map[balance.Color]int64{
		*scColor:          1,
		balance.ColorIOTA: 42,
	}, "sc after donation") {
		t.Fail()
	}

	if !wasps.VerifyAddressBalances(donorAddr, testutil.RequestFundsAmount-42, map[balance.Color]int64{
		balance.ColorIOTA: testutil.RequestFundsAmount - 42,
	}, "donor after donation") {
		t.Fail()
	}

	dwfOwnerClient := dwfclient.NewClient(scclient.New(
		wasps.NodeClient,
		wasps.WaspClient(0),
		scAddr,
		scOwner.SigScheme(),
		15*time.Second,
	))
	_, err = dwfOwnerClient.Withdraw(40)
	check(err, t)
	checkSuccess(err, t, "harvested 40")

	if !wasps.VerifyAddressBalances(scOwnerAddr, testutil.RequestFundsAmount-1+40, map[balance.Color]int64{
		balance.ColorIOTA: testutil.RequestFundsAmount - 1 + 40,
	}, "sc owner after withdraw") {
		t.Fail()
	}

	if !wasps.VerifyAddressBalances(scAddr, 1+2, map[balance.Color]int64{
		*scColor:          1,
		balance.ColorIOTA: 2,
	}, "sc after withdraw") {
		t.Fail()
	}

	if !wasps.VerifyAddressBalances(donorAddr, testutil.RequestFundsAmount-42, map[balance.Color]int64{
		balance.ColorIOTA: testutil.RequestFundsAmount - 42,
	}, "donor after withdraw") {
		t.Fail()
	}

	if !wasps.VerifySCStateVariables2(scAddr, map[kv.Key]interface{}{
		vmconst.VarNameOwnerAddress: scOwnerAddr[:],
		vmconst.VarNameProgramHash:  programHash[:],
		vmconst.VarNameDescription:  dwfimpl.Description,
	}) {
		t.Fail()
	}
}

func TestDWFDonateWithdrawNotAuthorised(t *testing.T) {
	wasps := setup(t, "TestDeployDWF")

	err := requestFunds(wasps, scOwnerAddr, "sc owner")
	check(err, t)

	donor := wallet.WithIndex(1)
	donorAddr := donor.Address()
	err = requestFunds(wasps, donorAddr, "donor")
	check(err, t)

	scAddr, scColor, err := startSmartContract(wasps, dwfimpl.ProgramHash, dwfimpl.Description)
	checkSuccess(err, t, "smart contract has been created and activated")

	dwfDonorClient := dwfclient.NewClient(scclient.New(
		wasps.NodeClient,
		wasps.WaspClient(0),
		scAddr,
		donor.SigScheme(),
		15*time.Second,
	))
	_, err = dwfDonorClient.Donate(42, "well done, I give you 42i")
	check(err, t)
	checkSuccess(err, t, "donated 42")

	if !wasps.VerifyAddressBalances(scOwnerAddr, testutil.RequestFundsAmount-1, map[balance.Color]int64{
		balance.ColorIOTA: testutil.RequestFundsAmount - 1,
	}, "sc owner in the end") {
		t.Fail()
	}

	if !wasps.VerifyAddressBalances(scAddr, 1+42, map[balance.Color]int64{
		*scColor:          1,
		balance.ColorIOTA: 42,
	}, "sc in the end") {
		t.Fail()
	}

	if !wasps.VerifyAddressBalances(donorAddr, testutil.RequestFundsAmount-42, map[balance.Color]int64{
		balance.ColorIOTA: testutil.RequestFundsAmount - 42,
	}, "donor in the end") {
		t.Fail()
	}

	// donor want to take back. Not authorised
	_, err = dwfDonorClient.Withdraw(40)
	check(err, t)
	checkSuccess(err, t, "sent harvest 40")

	if !wasps.VerifyAddressBalances(scOwnerAddr, testutil.RequestFundsAmount-1, map[balance.Color]int64{
		balance.ColorIOTA: testutil.RequestFundsAmount - 1,
	}, "sc owner in the end") {
		t.Fail()
	}

	if !wasps.VerifyAddressBalances(scAddr, 1+42, map[balance.Color]int64{
		*scColor:          1,
		balance.ColorIOTA: 42,
	}, "sc in the end") {
		t.Fail()
	}

	if !wasps.VerifyAddressBalances(donorAddr, testutil.RequestFundsAmount-42, map[balance.Color]int64{
		balance.ColorIOTA: testutil.RequestFundsAmount - 42,
	}, "donor in the end") {
		t.Fail()
	}

	if !wasps.VerifySCStateVariables2(scAddr, map[kv.Key]interface{}{
		vmconst.VarNameOwnerAddress: scOwnerAddr[:],
		vmconst.VarNameProgramHash:  programHash[:],
		vmconst.VarNameDescription:  dwfimpl.Description,
	}) {
		t.Fail()
	}
}
