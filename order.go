package sfapi

import (
	"math"
	"strconv"
)

// return from creating a new order, or checking status of previous order
type Order struct {
	APIResponse
	Symbol      string
	Venue       string
	Direction   string
	OriginalQty int64 // original request
	Qty         int64 // quantity left outstanding
	Price       int64
	OrderType   string
	Id          OrderId
	Account     string    // TODO: make this an AccountId
	Ts          Timestamp // when order was placed
	Fills       []Fill
	TotalFilled int64
	Open        bool
}

func (a *API) NewOrderForStock(request *OrderRequest) (result Order) {
	//	log.Println(request.String())

	url := baseURL + "/venues/" + request.Venue + "/stocks/" + request.Stock + "/orders"
	a.doPost(url, request, &result)
	return
}

func (a *API) GetUpdatedOrderStatus(orderId OrderId) (result Order) {
	orderIdAsString := strconv.Itoa(int(orderId))
	url := baseURL + "/venues/" + a.Venue + "/stocks/" + a.Stock + "/orders/" + orderIdAsString
	a.doGet(url, &result)
	return
}

func (a *API) CancelOrder(orderId OrderId) (result Order) {
	orderIdAsString := strconv.Itoa(int(orderId))
	url := baseURL + "/venues/" + a.Venue + "/stocks/" + a.Stock + "/orders/" + orderIdAsString
	a.doDelete(url, &result)
	return
}

func (o *Order) GuessPrice() (result int64) {
	// for limit orders use the given order price; for market, use the "worst" fill price
	result = o.Price

	if o.OrderType == "market" {
		if o.Direction == "buy" {
			result = o.HighestFillPrice()
		} else {
			result = o.LowestFillPrice()
		}
	}

	return
}

func (o Order) HighestFillPrice() (result int64) {
	result = int64(0)

	for _, fill := range o.Fills {
		if fill.Price > result {
			result = fill.Price
		}
	}

	return
}

func (o Order) LowestFillPrice() (result int64) {
	result = math.MaxInt64

	for _, fill := range o.Fills {
		if fill.Price < result {
			result = fill.Price
		}
	}

	return
}

type ByTime []Order

func (o ByTime) Len() int           { return len(o) }
func (o ByTime) Swap(i, j int)      { o[i], o[j] = o[j], o[i] }
func (o ByTime) Less(i, j int) bool { return o[i].Ts.Time().Before(o[j].Ts.Time()) }
