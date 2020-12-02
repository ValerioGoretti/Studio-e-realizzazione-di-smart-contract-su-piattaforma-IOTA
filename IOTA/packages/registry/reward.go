package registry

import (
	"wasp/plugins/config"

	"github.com/iotaledger/goshimmer/dapps/valuetransfers/packages/address"
	flag "github.com/spf13/pflag"
)

const (
	// CfgBindAddress defines the config flag of the web API binding address.
	CfgRewardAddress = "reward.address"
)

func InitFlags() {
	flag.String(CfgRewardAddress, "", "reward address for this Wasp node. Empty (default) means no rewards are collected")
}

func GetRewardAddress(scaddr *address.Address) address.Address {
	//TODO
	ret, err := address.FromBase58(config.Node.GetString(CfgRewardAddress))
	if err != nil {
		ret = address.Address{}
	}
	return ret
}
