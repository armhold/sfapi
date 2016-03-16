package sfapi

type OrderBook struct {
	APIResponse
	Venue  string
	Symbol string
	Bids   []BidOrAsk
	Asks   []BidOrAsk
	Ts     Timestamp
}

func (a *API) GetOrderBook() (result OrderBook) {
	url := baseURL + "/venues/" + a.Venue + "/stocks/" + a.Stock
	a.doGet(url, &result)
	return
}

func (ob *OrderBook) LowestAskPrice() int64 {
	lowest := ob.Asks[0].Price

	for _, ask := range ob.Asks {
		if ask.Price < lowest {
			lowest = ask.Price
		}
	}

	return lowest
}

func (ob *OrderBook) HighestBidPrice() int64 {
	highest := ob.Bids[0].Price

	for _, bid := range ob.Bids {
		if bid.Price > highest {
			highest = bid.Price
		}
	}

	return highest
}

func (ob *OrderBook) String() (result string) {

	result = "*** ORDERBOOK START *** \n"
	for _, bid := range ob.Bids {
		result += "\t" + bid.String() + "\n"
	}

	for _, ask := range ob.Asks {
		result += "\t" + ask.String() + "\n"
	}

	result += "=== ORDERBOOK END ==="

	return
}
