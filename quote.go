package sfapi

import (
	"fmt"
	"time"
)

type Quote struct {
	APIResponse
	Symbol    string
	Venue     string
	Bid       int64     // best price currently bid for the stock
	Ask       int64     // best price currently offered for the stock
	BidSize   int64     // aggregate size of all orders at the best bid
	AskSize   int64     // aggregate size of all orders at the best ask
	BidDepth  int64     // aggregate size of *all bids*
	AskDepth  int64     // aggregate size of *all asks*
	Last      int64     // price of last trade
	LastSize  int64     // quantity of last trade
	LastTrade Timestamp // timestamp of last trade
	QuoteTime Timestamp // ts we last updated quote at (server-side)
	Anomalous bool      // not part of API spec, but we use this in level4
}

func (q *Quote) Time() time.Time {
	result, err := time.Parse(time.RFC3339Nano, string(q.QuoteTime))
	Must(err)

	return result
}

func (a *API) GetQuote() (result Quote) {
	url := baseURL + "/venues/" + a.Venue + "/stocks/" + a.Stock + "/quote"
	a.doGet(url, &result)
	return
}

func (q Quote) Spread() int64 {
	return q.Ask - q.Bid
}

func (q Quote) Mid() int64 {
	return q.Bid + q.Spread()/2
}

// NB: this will be invoked whenever Quote is printed via the %v verb in the fmt package.
// See https://golang.org/pkg/fmt/
//
func (q Quote) String() string {
	return fmt.Sprintf("Bid: %s[%d], Ask: %s[%d], Last: %s, Spread: %s, Mid: %s", Dollars(q.Bid), q.BidSize, Dollars(q.Ask), q.AskSize, Dollars(q.Last), Dollars(q.Spread()), Dollars(q.Mid()))
}

type ByQuoteTime []Quote

func (q ByQuoteTime) Len() int           { return len(q) }
func (q ByQuoteTime) Swap(i, j int)      { q[i], q[j] = q[j], q[i] }
func (q ByQuoteTime) Less(i, j int) bool { return q[i].QuoteTime.Time().Before(q[j].QuoteTime.Time()) }
