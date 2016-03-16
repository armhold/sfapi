package sfapi

type Heartbeat struct {
	APIResponse
}

func (a *API) CheckAPIUp() (result Heartbeat) {
	url := baseURL + "/heartbeat"
	a.doGet(url, &result)
	return
}

func (a *API) CheckVenueUp() (result Heartbeat) {
	url := baseURL + "/venues/" + a.Venue + "/heartbeat"
	a.doGet(url, &result)
	return
}
