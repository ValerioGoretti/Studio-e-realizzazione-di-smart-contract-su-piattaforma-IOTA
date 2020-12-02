// hard coded implementation of the FairAuction smart contract
// The auction dApp is automatically run by committee, a distributed market for colored tokens
package fairauction

import (
	"bytes"
	"sort"
	"time"

	"wasp/packages/hashing"

	"wasp/packages/kv"
	"wasp/packages/sctransaction"
	"wasp/packages/util"

	"wasp/packages/vm/vmtypes"

	"github.com/iotaledger/goshimmer/dapps/valuetransfers/packages/address"
	"github.com/iotaledger/goshimmer/dapps/valuetransfers/packages/balance"
)

// program has is an id of the program
const ProgramHash = "4NbQFgvnsfgE3n9ZhtJ3p9hWZzfYUEDHfKU93wp8UowB"
const Description = "FairAuction, a PoC smart contract"

// implement Processor and EntryPoint interfaces

type fairAuctionProcessor map[sctransaction.RequestCode]fairAuctionEntryPoint

type fairAuctionEntryPoint func(ctx vmtypes.Sandbox)

const (
	RequestInitSC          = sctransaction.RequestCode(0) // NOP
	RequestStartAuction    = sctransaction.RequestCode(1)
	RequestFinalizeAuction = sctransaction.RequestCode(2)
	RequestPlaceBid        = sctransaction.RequestCode(3)
	RequestSetOwnerMargin  = sctransaction.RequestCode(4 | sctransaction.RequestCodeProtected)
)

// the processor is a map of entry points
var entryPoints = fairAuctionProcessor{
	RequestInitSC:          initSC,
	RequestStartAuction:    startAuction,
	RequestFinalizeAuction: finalizeAuction,
	RequestPlaceBid:        placeBid,
	RequestSetOwnerMargin:  setOwnerMargin,
}

// string constants for request arguments and state variable names
const (
	// request vars
	VarReqAuctionColor                = "color"
	VarReqStartAuctionDescription     = "dscr"
	VarReqStartAuctionDurationMinutes = "duration"
	VarReqStartAuctionMinimumBid      = "minimum" // in iotas
	VarReqOwnerMargin                 = "ownerMargin"

	// state vars
	VarStateAuctions            = "auctions"
	VarStateLog                 = "log"
	VarStateOwnerMarginPromille = "ownerMargin" // owner margin in percents
)

const (
	// minimum duration of auction
	MinAuctionDurationMinutes = 1
	MaxAuctionDurationMinutes = 120 // max 2 hours

	// default duration of the auction
	AuctionDurationDefaultMinutes = 60
	// Owner of the smart contract takes %% from the winning bid. The default, min, max
	OwnerMarginDefault = 50  // 5%
	OwnerMarginMin     = 5   // minimum 0.5%
	OwnerMarginMax     = 100 // max 10%
	MaxDescription     = 150
)

// validating constants at node boot
func init() {
	if OwnerMarginMax > 1000 ||
		OwnerMarginMin < 0 ||
		OwnerMarginDefault < OwnerMarginMin ||
		OwnerMarginDefault > OwnerMarginMax ||
		OwnerMarginMin > OwnerMarginMax {
		panic("wrong constants")
	}
}

// statical link point to the Wasp node
func GetProcessor() vmtypes.Processor {
	return entryPoints
}

func (v fairAuctionProcessor) GetDescription() string {
	return "FairAuction hard coded smart contract program"
}

func (v fairAuctionProcessor) GetEntryPoint(code sctransaction.RequestCode) (vmtypes.EntryPoint, bool) {
	f, ok := v[code]
	return f, ok
}

func (ep fairAuctionEntryPoint) Run(ctx vmtypes.Sandbox) {
	ep(ctx)
}

func (ep fairAuctionEntryPoint) WithGasLimit(_ int) vmtypes.EntryPoint {
	return ep
}

// AuctionInfo describes active auction
type AuctionInfo struct {
	// color of the tokens for sale. Max one auction per color at same time is allowed
	// all tokens are being sold as one lot
	Color balance.Color
	// number of tokens for sale
	NumTokens int64
	// minimum bid. Set by the auction initiator
	MinimumBid int64
	// any text, like "AuctionOwner of the token have a right to call me for a date". Set by auction initiator
	Description string
	// timestamp when auction started
	WhenStarted int64
	// duration of the auctions in minutes. Should be >= MinAuctionDurationMinutes
	DurationMinutes int64
	// address which issued StartAuction transaction
	AuctionOwner address.Address
	// total deposit by the auction owner. Iotas sent by the auction owner together with the tokens for sale in the same
	// transaction.
	TotalDeposit int64
	// AuctionOwner's margin in promilles, taken at the moment of creation of smart contract
	OwnerMargin int64
	// list of bids to the auction
	Bids []*BidInfo
}

// BidInfo represents one bid to the auction
type BidInfo struct {
	// total sum of the bid = total amount of iotas available in the request - 1 - SC reward - ServiceFeeBid
	// the total is a cumulative sum of all bids from the same bidder
	Total int64
	// originator of the bid
	Bidder address.Address
	// timestamp Unix nano
	When int64
}

func (ai *AuctionInfo) SumOfBids() int64 {
	sum := int64(0)
	for _, bid := range ai.Bids {
		sum += bid.Total
	}
	return sum
}

func (ai *AuctionInfo) WinningBid() *BidInfo {
	var winner *BidInfo
	for _, bi := range ai.Bids {
		if bi.Total < ai.MinimumBid {
			continue
		}
		if winner == nil || bi.WinsAgainst(winner) {
			winner = bi
		}
	}
	return winner
}

func (ai *AuctionInfo) Due() int64 {
	return ai.WhenStarted + ai.DurationMinutes*time.Minute.Nanoseconds()
}

func (bi *BidInfo) WinsAgainst(other *BidInfo) bool {
	if bi.Total < other.Total {
		return false
	}
	if bi.Total > other.Total {
		return true
	}
	return bi.When < other.When
}

// NOP
func initSC(ctx vmtypes.Sandbox) {
	ctx.Publish("initSC")
}

// startAuction processes the StartAuction request
// Arguments:
// - VarReqAuctionColor: color of the tokens for sale
// - VarReqStartAuctionDescription: description of the lot
// - VarReqStartAuctionMinimumBid: minimum price for the whole lot
// - VarReqStartAuctionDurationMinutes: duration of auction
// Request transaction must contain at least number of iotas >= of current owner margin from the minimum bid
// (not including node reward with request token)
// Tokens for sale must be included into the request transaction
func startAuction(ctx vmtypes.Sandbox) {
	ctx.Publish("startAuction begin")

	sender := ctx.AccessRequest().Sender()
	reqArgs := ctx.AccessRequest().Args()
	account := ctx.AccessSCAccount()

	// check how many iotas the request contains
	totalDeposit := account.AvailableBalanceFromRequest(&balance.ColorIOTA)
	if totalDeposit < 1 {
		// it is expected at least 1 iota in deposit
		// this 1 iota is needed as a "operating capital for the time locked request to itself"
		// refund iotas
		refundFromRequest(ctx, &balance.ColorIOTA, 1)

		ctx.Publish("startAuction: exit 0: must be at least 1i in deposit")
		return
	}

	// take current setting of the smart contract owner margin
	ownerMargin := GetOwnerMarginPromille(ctx.AccessState().GetInt64(VarStateOwnerMarginPromille))

	// determine color of the token for sale
	colh, ok, err := reqArgs.GetHashValue(VarReqAuctionColor)
	if err != nil || !ok {
		// incorrect request arguments, colore for sale is not determined
		// refund half of the deposit in iotas
		refundFromRequest(ctx, &balance.ColorIOTA, totalDeposit/2)

		ctx.Publish("startAuction: exit 1")
		return
	}
	colorForSale := balance.Color(*colh)
	if colorForSale == balance.ColorIOTA || colorForSale == balance.ColorNew {
		// reserved color code are not allowed
		// refund half
		refundFromRequest(ctx, &balance.ColorIOTA, totalDeposit/2)

		ctx.Publish("startAuction: exit 2")
		return
	}

	// determine amount of colored tokens for sale. They must be in the outputs of the request transaction
	tokensForSale := account.AvailableBalanceFromRequest(&colorForSale)
	if tokensForSale == 0 {
		// no tokens transferred. Refund half of deposit
		refundFromRequest(ctx, &balance.ColorIOTA, totalDeposit/2)

		ctx.Publish("startAuction exit 3: no tokens for sale")
		return
	}

	// determine minimum bid
	minimumBid, _, err := reqArgs.GetInt64(VarReqStartAuctionMinimumBid)
	if err != nil {
		// wrong argument. Hard reject, no refund

		ctx.Publish("startAuction: exit 4")
		return
	}
	// ensure tokens are not sold for the minimum price less than 1 iota per token!
	if minimumBid < tokensForSale {
		minimumBid = tokensForSale
	}

	// check if enough iotas for service fees to create the auction
	expectedDeposit := GetExpectedDeposit(minimumBid, ownerMargin)

	if totalDeposit < expectedDeposit {
		// not enough fees
		// return half of expected deposit and all tokens for sale (if any)
		harvest := expectedDeposit / 2
		if harvest < 1 {
			harvest = 1
		}
		refundFromRequest(ctx, &balance.ColorIOTA, harvest)
		refundFromRequest(ctx, &colorForSale, 0)

		ctx.Publishf("startAuction: not enough iotas for the fee. Expected %d, got %d", expectedDeposit, totalDeposit)
		return
	}

	// determine duration of the auction. Take default if no set in request and ensure minimum
	duration, ok, err := reqArgs.GetInt64(VarReqStartAuctionDurationMinutes)
	if err != nil {
		// fatal error
		return
	}
	if !ok {
		duration = AuctionDurationDefaultMinutes
	}
	if duration < MinAuctionDurationMinutes {
		duration = MinAuctionDurationMinutes
	}
	if duration > MaxAuctionDurationMinutes {
		duration = MaxAuctionDurationMinutes
	}

	// read description text from the request
	description, ok, err := reqArgs.GetString(VarReqStartAuctionDescription)
	if err != nil {
		return
	}
	if !ok {
		description = "N/A"
	}
	description = util.GentleTruncate(description, MaxDescription)

	// find out if auction for this color already exist in the dictionary
	auctions := ctx.AccessState().GetDictionary(VarStateAuctions)
	if b := auctions.GetAt(colorForSale.Bytes()); b != nil {
		// auction already exists. Ignore sale auction.
		// refund iotas less fee
		refundFromRequest(ctx, &balance.ColorIOTA, expectedDeposit/2)
		// return all tokens for sale
		refundFromRequest(ctx, &colorForSale, 0)

		ctx.Publish("startAuction: exit 6")
		return
	}

	// create record for the new auction in the dictionary
	aiData := util.MustBytes(&AuctionInfo{
		Color:           colorForSale,
		NumTokens:       tokensForSale,
		MinimumBid:      minimumBid,
		Description:     description,
		WhenStarted:     ctx.GetTimestamp(),
		DurationMinutes: duration,
		AuctionOwner:    sender,
		TotalDeposit:    totalDeposit,
		OwnerMargin:     ownerMargin,
	})
	auctions.SetAt(colorForSale.Bytes(), aiData)

	ctx.Publishf("New auction record. color: %s, numTokens: %d, minBid: %d, ownerMargin: %d duration %d minutes",
		colorForSale.String(), tokensForSale, minimumBid, ownerMargin, duration)

	// prepare and send request FinalizeAuction to self time-locked for the duration
	// the FinalizeAuction request will be time locked for the duration and then auction will be run
	args := kv.NewMap()
	args.Codec().SetHashValue(VarReqAuctionColor, (*hashing.HashValue)(&colorForSale))
	ctx.SendRequestToSelfWithDelay(RequestFinalizeAuction, args, uint32(duration*60))

	//logToSC(ctx, fmt.Sprintf("start auction. For sale %d tokens of color %s. Minimum bid: %di. Duration %d minutes",
	//	tokensForSale, colorForSale.String(), minimumBid, duration))

	ctx.Publishf("startAuction: success. Auction: '%s', color: %s, duration: %d",
		description, colorForSale.String(), duration)
}

// placeBid is a request to place a bid in the auction for the particular color
// The request transaction must contain at least:
// - 1 request token + Bid/rise amount
// In case it is not the first bid by this bidder, respective iotas are treated as
// a rise of the bid and are added to the total
// Arguments:
// - VarReqAuctionColor: color of the tokens for sale
func placeBid(ctx vmtypes.Sandbox) {
	ctx.Publish("placeBid: begin")

	// all iotas in the request transaction are considered a bid/rise sum
	// it also means several bids can't be placed in the same transaction <-- TODO generic solution for it
	bidAmount := ctx.AccessSCAccount().AvailableBalanceFromRequest(&balance.ColorIOTA)
	if bidAmount == 0 {
		// no iotas sent
		ctx.Publish("placeBid: exit 0")
		return
	}

	reqArgs := ctx.AccessRequest().Args()
	// determine color of the bid
	colh, ok, err := reqArgs.GetHashValue(VarReqAuctionColor)
	if err != nil {
		// inconsistency. return all?
		ctx.Publish("placeBid: exit 1")
		return
	}
	if !ok {
		// missing argument
		ctx.Publish("placeBid: exit 2")
		refundFromRequest(ctx, &balance.ColorIOTA, 0)
		return
	}

	col := balance.Color(*colh)
	if col == balance.ColorIOTA || col == balance.ColorNew {
		// reserved color not allowed. Incorrect arguments
		refundFromRequest(ctx, &balance.ColorIOTA, 0)
		ctx.Publish("placeBid: exit 3")
		return
	}

	// find the auction
	auctions := ctx.AccessState().GetDictionary(VarStateAuctions)
	data := auctions.GetAt(col.Bytes())
	if data == nil {
		// no such auction. refund everything
		refundFromRequest(ctx, &balance.ColorIOTA, 0)
		ctx.Publish("placeBid: exit 4")
		return
	}
	// unmarshal auction data
	ai := &AuctionInfo{}
	if err := ai.Read(bytes.NewReader(data)); err != nil {
		// internal error
		ctx.Publish("placeBid: exit 6")
		return
	}
	// determine the sender of the bid
	sender := ctx.AccessRequest().Sender()

	// find bids of this bidder in the auction
	var bi *BidInfo
	for _, bitmp := range ai.Bids {
		if bitmp.Bidder == sender {
			bi = bitmp
			break
		}
	}
	if bi == nil {
		// first bid by the bidder. Create new bid record
		ai.Bids = append(ai.Bids, &BidInfo{
			Total:  bidAmount,
			Bidder: sender,
			When:   ctx.GetTimestamp(),
		})
		//logToSC(ctx, fmt.Sprintf("place bid. Auction color %s, total %di", col.String(), bidAmount))
	} else {
		// bidder has bid already. Treated it as a rise
		bi.Total += bidAmount
		bi.When = ctx.GetTimestamp()

		//logToSC(ctx, fmt.Sprintf("rise bid. Auction color %s, total %di", col.String(), bi.Total))
	}
	// marshal the whole auction info and save it into the state (the dictionary of auctions)
	data = util.MustBytes(ai)
	auctions.SetAt(col.Bytes(), data)

	ctx.Publishf("placeBid: success. Auction: '%s'", ai.Description)
}

// finalizeAuction selects the winner and sends tokens to him.
// returns bid amounts to other bidders.
// The request is time locked for the period of the auction. It won't be executed if sent
// not by the smart contract instance itself
// Arguments:
// - VarReqAuctionColor: color of the auction
func finalizeAuction(ctx vmtypes.Sandbox) {
	ctx.Publish("finalizeAuction begin")

	accessReq := ctx.AccessRequest()
	if accessReq.Sender() != *ctx.GetSCAddress() {
		// finalizeAuction request can only be sent by the smart contract to itself. Otherwise it is NOP
		return
	}
	reqArgs := accessReq.Args()

	// determine color of the auction to finalize
	colh, ok, err := reqArgs.GetHashValue(VarReqAuctionColor)
	if err != nil || !ok {
		// wrong request arguments
		// internal error. Refund completely?
		ctx.Publish("finalizeAuction: exit 1")
		return
	}
	col := balance.Color(*colh)
	if col == balance.ColorIOTA || col == balance.ColorNew {
		// inconsistency
		ctx.Publish("finalizeAuction: exit 2")
		return
	}

	// find the record of the auction by color
	auctDict := ctx.AccessState().GetDictionary(VarStateAuctions)
	data := auctDict.GetAt(col.Bytes())
	if data == nil {
		// auction with this color does not exist. Inconsistency
		ctx.Publish("finalizeAuction: exit 3")
		return
	}

	// decode the Action record
	ai := &AuctionInfo{}
	if err := ai.Read(bytes.NewReader(data)); err != nil {
		// internal error. Refund completely?
		ctx.Publish("finalizeAuction: exit 4")
		return
	}

	account := ctx.AccessSCAccount()

	// find the winning amount and determine respective ownerFee
	winningAmount := int64(0)
	for _, bi := range ai.Bids {
		if bi.Total > winningAmount {
			winningAmount = bi.Total
		}
	}

	var winner *BidInfo
	var winnerIndex int

	// SC owner takes OwnerMargin (promille) fee from either minimum bid or from winning sum but not less than 1i
	ownerFee := (ai.MinimumBid * ai.OwnerMargin) / 1000
	if ownerFee < 1 {
		ownerFee = 1
	}

	// find the winner (if any). Take first if equal sums
	// minimum bid is always positive, at least 1 iota per colored token
	if winningAmount >= ai.MinimumBid {
		// there's winner. Select it.
		// Fee is re-calculated according to the winning sum
		ownerFee = (winningAmount * ai.OwnerMargin) / 1000
		if ownerFee < 1 {
			ownerFee = 1
		}

		winners := make([]*BidInfo, 0)
		for _, bi := range ai.Bids {
			if bi.Total == winningAmount {
				winners = append(winners, bi)
			}
		}
		sort.Slice(winners, func(i, j int) bool {
			return winners[i].When < winners[j].When
		})
		winner = winners[0]
		for i, bi := range ai.Bids {
			if bi == winner {
				winnerIndex = i
				break
			}
		}
	}

	// take fee for the smart contract owner
	feeTaken := ctx.AccessSCAccount().HarvestFees(ownerFee - 1)
	ctx.Publishf("finalizeAuction: harvesting SC owner fee: %d (+1 self request token left in SC)", feeTaken)

	if winner != nil {
		// send sold tokens to the winner
		account.MoveTokens(&ai.Bids[winnerIndex].Bidder, &ai.Color, ai.NumTokens)
		// send winning amount and return deposit sum less fees to the owner of the auction
		account.MoveTokens(&ai.AuctionOwner, &balance.ColorIOTA, winningAmount+ai.TotalDeposit-ownerFee)

		for i, bi := range ai.Bids {
			if i != winnerIndex {
				// return staked sum to the non-winner
				account.MoveTokens(&bi.Bidder, &balance.ColorIOTA, bi.Total)
			}
		}
		//logToSC(ctx, fmt.Sprintf("close auction. Color: %s. Winning bid: %di", col.String(), winner.Total))

		ctx.Publishf("finalizeAuction: winner is %s, winning amount = %d", winner.Bidder.String(), winner.Total)
	} else {
		// return unsold tokens to auction owner
		if account.MoveTokens(&ai.AuctionOwner, &ai.Color, ai.NumTokens) {
			ctx.Publishf("returned unsold tokens to auction owner. %s: %d", ai.Color.String(), ai.NumTokens)
		}

		// return deposit less fees less 1 iota
		if account.MoveTokens(&ai.AuctionOwner, &balance.ColorIOTA, ai.TotalDeposit-ownerFee) {
			ctx.Publishf("returned deposit less fees: %d", ai.TotalDeposit-ownerFee)
		}

		// return bids to bidders
		for _, bi := range ai.Bids {
			if account.MoveTokens(&bi.Bidder, &balance.ColorIOTA, bi.Total) {
				ctx.Publishf("returned bid to bidder: %d -> %s", bi.Total, bi.Bidder.String())
			} else {
				avail := ctx.AccessSCAccount().AvailableBalance(&balance.ColorIOTA)
				ctx.Publishf("failed to return bid to bidder: %d -> %s. Available: %d", bi.Total, bi.Bidder.String(), avail)
			}
		}
		//logToSC(ctx, fmt.Sprintf("close auction. Color: %s. No winner.", col.String()))

		ctx.Publishf("finalizeAuction: winner wasn't selected out of %d bids", len(ai.Bids))
	}

	// delete auction record
	auctDict.DelAt(col.Bytes())

	ctx.Publishf("finalizeAuction: success. Auction: '%s'", ai.Description)
}

// setOwnerMargin is a request to set the service fee to place a bid
// It is protected, i.e. must be sent by the owner of the smart contract
// Arguments:
// - VarReqOwnerMargin: the margin value in promilles
func setOwnerMargin(ctx vmtypes.Sandbox) {
	ctx.Publish("setOwnerMargin: begin")
	margin, ok, err := ctx.AccessRequest().Args().GetInt64(VarReqOwnerMargin)
	if err != nil || !ok {
		ctx.Publish("setOwnerMargin: exit 1")
		return
	}
	if margin < OwnerMarginMin {
		margin = OwnerMarginMin
	} else if margin > OwnerMarginMax {
		margin = OwnerMarginMax
	}
	ctx.AccessState().SetInt64(VarStateOwnerMarginPromille, margin)
	ctx.Publishf("setOwnerMargin: success. ownerMargin set to %d%%", margin/10)
}

// TODO implement universal 'refund' function to be used in rollback situations
// refundFromRequest returns all tokens of the given color to the sender minus sunkFee
func refundFromRequest(ctx vmtypes.Sandbox, color *balance.Color, harvest int64) {
	account := ctx.AccessSCAccount()
	ctx.AccessSCAccount().HarvestFeesFromRequest(harvest)
	available := account.AvailableBalanceFromRequest(color)
	sender := ctx.AccessRequest().Sender()
	ctx.AccessSCAccount().HarvestFeesFromRequest(harvest)
	account.MoveTokensFromRequest(&sender, color, available)
}

func logToSC(ctx vmtypes.Sandbox, msg string) {
	ctx.AccessState().GetTimestampedLog(VarStateLog).Append(ctx.GetTimestamp(), []byte(msg))
}
