package util

import (
	"fmt"
	"os"

	"wasp/tools/wwallet/config"

	"github.com/iotaledger/goshimmer/dapps/valuetransfers/packages/transaction"
)

func PostTransaction(tx *transaction.Transaction) {
	WithTransaction(func() (*transaction.Transaction, error) {
		return tx, config.GoshimmerClient().PostTransaction(tx)
	})
}

func WithTransaction(f func() (*transaction.Transaction, error)) {
	tx, err := f()
	check(err)

	if config.WaitForCompletion {
		check(config.GoshimmerClient().WaitForConfirmation(tx.ID()))
	}
}

func check(err error) {
	if err != nil {
		fmt.Printf("error: %s\n", err)
		os.Exit(1)
	}
}
