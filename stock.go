package sfapi

type Stock struct {
	Name   string
	Symbol string
}

type StocksOnVenue struct {
	APIResponse
	Symbols []Stock
}

func (a *API) GetStocksOnAVenue() (result StocksOnVenue) {
	url := baseURL + "/venues/" + a.Venue + "/stocks/"
	a.doGet(url, &result)
	return
}
