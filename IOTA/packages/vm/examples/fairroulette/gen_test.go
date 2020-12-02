package fairroulette

import (
	"testing"
	"wasp/packages/hashing"
)

// it is needed only to generate dummy hash code
func TestGenData(t *testing.T) {
	h := hashing.HashStrings("FairRoulette smart contract")
	t.Logf("FairRulette program hash = %s", h.String())
}
