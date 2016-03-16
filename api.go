package sfapi

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

const (
	baseURL string = "https://api.stockfighter.io/ob/api"
	gmURL   string = "https://api.stockfighter.io/gm"
)

type API struct {
	Account        string
	Venue          string
	Stock          string
	AuthKey        string
	QuoteAlertChan chan QuoteAlert

	client http.Client
}

type Timestamp string

func (t Timestamp) Time() (result time.Time) {
	result, err := time.Parse(time.RFC3339Nano, string(t))
	if err != nil {
		log.Fatal(err)
	}

	return
}

type OrderId int64

type APIResponse struct {
	Ok    bool
	Error string
}

type Fill struct {
	Price int64
	Qty   int64
	Ts    Timestamp
}

type BidOrAsk struct {
	Price int64
	Qty   int64
	IsBuy bool
}

func (b *BidOrAsk) String() (result string) {
	if b.IsBuy {
		result = "buy "
	} else {
		result = "sell "
	}

	result += fmt.Sprintf(" %d @ %s", b.Qty, Dollars(b.Price))
	return
}

type apiResponse interface {
	IsOk() bool
	GetError() string
	AddErrorString(string)
	AddError(error)
}

func (b *APIResponse) IsOk() bool {
	return b.Ok
}

func (b *APIResponse) GetError() string {
	return b.Error
}

func (b *APIResponse) AddError(err error) {
	if err == nil {
		return
	}

	b.AddErrorString(err.Error())
}

func (b *APIResponse) AddErrorString(err string) {
	if err == "" {
		return
	}

	if b.Error == "" {
		b.Error = err
	} else {
		b.Error = fmt.Sprintf("%s, %s", b.Error, err)
	}
}

func NewAPI() (result *API) {
	result = &API{QuoteAlertChan: make(chan QuoteAlert, 1000)}

	result.AuthKey = os.Getenv("STARFIGHTER_API_KEY")

	if result.AuthKey == "" {
		log.Fatal("STARFIGHTER_API_KEY environment variable not set")
	}

	result.client = http.Client{}
	return
}

func APIForAccount(account, venue, stock, authKey string) (result *API) {
	result = NewAPI()
	result.Account = account
	result.Venue = venue
	result.Stock = stock
	result.AuthKey = authKey

	return
}

func APIFromCommandLineArgs() (result *API) {
	result = NewAPI()

	var venue, stock, account string
	var auto bool

	flag.BoolVar(&auto, "auto", false, "to read from saved GameLevel")
	flag.StringVar(&venue, "venue", "", "")
	flag.StringVar(&stock, "stock", "", "")
	flag.StringVar(&account, "account", "", "")

	flag.Parse()

	if auto {
		gl, err := result.ReadGameLevel()
		if err != nil {
			log.Fatalf("error reading previous gamelevel: %s", err)
		}
		venue = gl.Venues[0]
		stock = gl.Tickers[0]
		account = gl.Account
	}

	result.Venue = venue
	result.Stock = stock
	result.Account = account

	if result.Venue == "" || result.Stock == "" || result.Account == "" {
		flag.Usage()
		os.Exit(1)
	}

	return
}

func (a *API) doGet(url string, apiResponse apiResponse) {
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Starfighter-Authorization", a.AuthKey)

	resp, err := a.client.Do(req)
	if err != nil {
		apiResponse.AddError(err)
		return
	}
	defer resp.Body.Close()

	// special case where we return error immediately if not 200, before trying to decode result
	if resp.StatusCode != 200 {
		// preserve existing apiResponse.Error, if present
		s := fmt.Sprintf("%d - %s", resp.StatusCode, apiResponse.GetError())
		apiResponse.AddErrorString(s)
	}

	err = json.NewDecoder(resp.Body).Decode(&apiResponse)
	apiResponse.AddError(err)
}

func (a *API) doPost(url string, postData interface{}, apiResponse apiResponse) {
	b, err := json.Marshal(postData)
	if err != nil {
		apiResponse.AddError(err)
		return
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Starfighter-Authorization", a.AuthKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		apiResponse.AddError(err)
		return
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&apiResponse)
	apiResponse.AddError(err)
}

// TODO: almost exact copy/paste of doGet()
func (a *API) doDelete(url string, apiResponse apiResponse) {
	req, err := http.NewRequest("DELETE", url, nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Starfighter-Authorization", a.AuthKey)

	resp, err := a.client.Do(req)
	if err != nil {
		apiResponse.AddError(err)
		return
	}
	defer resp.Body.Close()

	// special case where we return error immediately if not 200, before trying to decode result
	if resp.StatusCode != 200 {
		// preserve existing apiResponse.Error, if present
		s := fmt.Sprintf("%d - %s", resp.StatusCode, apiResponse.GetError())
		apiResponse.AddErrorString(s)
	}

	err = json.NewDecoder(resp.Body).Decode(&apiResponse)
	apiResponse.AddError(err)
}
