package fairauction

func GetOwnerMarginPromille(ownerMargin int64, ok bool) int64 {
	if !ok {
		return OwnerMarginDefault
	}
	if ownerMargin > OwnerMarginMax {
		return OwnerMarginMax
	}
	if ownerMargin < OwnerMarginMin {
		return OwnerMarginMin
	}
	return ownerMargin
}

func GetExpectedDeposit(minimumBid int64, ownerMargin int64) int64 {
	// minimum deposit is owner margin from minimum bid
	expectedDeposit := (minimumBid * ownerMargin) / 1000
	// ensure that at least 1 iota is taken. It is needed for "operating capital"
	if expectedDeposit < 1 {
		return 1
	}
	return expectedDeposit
}
