package sfapi

type AllOrdersForStock struct {
	APIResponse
	Venue  string
	Orders []Order
}

func (a *API) GetStatusForAllOrdersInAStock() (result AllOrdersForStock) {
	url := baseURL + "/venues/" + a.Venue + "/accounts/" + a.Account + "/stocks/" + a.Stock + "/orders"
	a.doGet(url, &result)
	return
}
