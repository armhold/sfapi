package sfapi

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"

	"github.com/gorilla/websocket"
)

type QuoteAlert struct {
	APIResponse
	Quote Quote
}

type ExecutionAlert struct {
	APIResponse
	Account          string
	Venue            string
	Symbol           string
	Order            Order
	StandingId       int64
	IncomingId       int64
	Price            int64
	Filled           int64
	FilledAt         Timestamp
	StandingComplete bool
	IncomingComplete bool
}

func (ea *ExecutionAlert) Succinct() string {
	return fmt.Sprintf("%s: %s %d @ %s", PadString(ea.Account, 14), ea.Order.Direction, ea.Filled, Dollars(ea.Price))
}

func (a *API) RunQuoteTicker() {
	makeAlert := func() apiResponse { return &QuoteAlert{} }

	callback := func(alert apiResponse) {
		alertPtr := alert.(*QuoteAlert)
		a.QuoteAlertChan <- *alertPtr
	}

	a.runTicker(makeAlert, callback, "wss://api.stockfighter.io/ob/api/ws/"+a.Account+"/venues/"+a.Venue+"/tickertape/stocks/"+a.Stock)
}

func (a *API) RunExecutionTicker(reportingChan chan ExecutionAlert) {
	makeAlert := func() apiResponse { return &ExecutionAlert{} }

	callback := func(alert apiResponse) {
		alertPtr := alert.(*ExecutionAlert)
		reportingChan <- *alertPtr
	}

	a.runTicker(makeAlert, callback, "wss://api.stockfighter.io/ob/api/ws/"+a.Account+"/venues/"+a.Venue+"/executions/stocks/"+a.Stock)
}

// There is some ugliness here to avoid code duplication that would otherwise be solved by generics:
//
// 1. makeAlert() should return a pointer to an instance of the proper alert struct (e.g. ExecutionAlert, QuoteAlert).
//
// 2. callback() will be then be called with the same instance from (1), with its fields populated from
// the JSON. It can then run a type assertion on its arg, feed it to a properly-typed channel, etc.
//
func (a *API) runTicker(makeAlert func() apiResponse, callback func(alert apiResponse), urlEndpoint string) {
	u, err := url.Parse(urlEndpoint)
	Must(err)

	var c *websocket.Conn = nil

	openConnection := func() *websocket.Conn {
		c, _, err = websocket.DefaultDialer.Dial(u.String(), nil)
		if err != nil {
			log.Fatalf("error dialing websocket: %s", err)
		}

		return c
	}

	defer func() {
		if c != nil {
			c.Close()
		}
	}()

	for {
		if c == nil {
			c = openConnection()
		}

		_, message, err := c.ReadMessage()
		if err != nil {
			log.Printf("handling websocket connection error: \"%s\", continuing...\n", err)

			c = nil
			continue
		}

		alert := makeAlert()
		err = json.Unmarshal(message, alert)
		Must(err)

		if !alert.IsOk() || alert.GetError() != "" {
			log.Printf("received error from websocket stream: %s, TICKER EXITING\n", alert.GetError())
			return
		}

		callback(alert)
	}
}
