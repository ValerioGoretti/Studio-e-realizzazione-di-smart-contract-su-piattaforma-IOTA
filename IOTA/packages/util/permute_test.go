package util

import (
	"testing"
	"wasp/packages/hashing"
)

func TestPermute(t *testing.T) {

	for n := uint16(1); n < 1000; n = n + 3 {
		for k := 0; k < 10; k++ {
			perm := NewPermutation16(n, hashing.RandomHash(nil).Bytes())
			if !ValidPermutation(perm.GetArray()) {
				t.Fatalf("invalid permutation %+v", perm)
			}
		}
	}
}
