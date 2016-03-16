package sfapi

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"strconv"
)

const (
	gameLevelFile = "/tmp/gamelevel"
)

type Instructions struct {
	Instructions string
	OrderTypes   string `json:"Order Types"`
}

type GameLevel struct {
	APIResponse
	Account              string
	InstanceId           int
	Instructions         *Instructions
	SecondsPerTradingDay int64
	Tickers              []string
	Venues               []string
	// TODO: Balances, etc...
}

func (a *API) StartLevel3() (result GameLevel) {
	url := gmURL + "/levels/sell_side"
	a.doPost(url, nil, &result)
	return
}

func (a *API) StartLevel4() (result GameLevel) {
	url := gmURL + "/levels/dueling_bulldozers"
	a.doPost(url, nil, &result)
	return
}

func (a *API) StartLevel5() (result GameLevel) {
	url := gmURL + "/levels/irrational_exuberance"
	a.doPost(url, nil, &result)
	return
}

func (a *API) StartLevel6() (result GameLevel) {
	url := gmURL + "/levels/making_amends"
	a.doPost(url, nil, &result)
	return
}

func (a *API) StopLevel(instanceId int) (result APIResponse) {
	instanceIdAsString := strconv.Itoa(instanceId)
	url := gmURL + "/instances/" + instanceIdAsString + "/stop"
	a.doPost(url, nil, &result)
	return
}

// NB: this is currently broken in stockfighter.
// See: https://discuss.starfighters.io/t/the-gm-api-how-to-start-stop-restart-resume-trading-levels-automagically/143/37?u=armhold
func (a *API) RestartLevel(instanceId int) (result GameLevel) {
	instanceIdAsString := strconv.Itoa(instanceId)
	url := gmURL + "/instances/" + instanceIdAsString + "/restart"
	a.doPost(url, nil, &result)
	return
}

func (g *GameLevel) CopyPaste() string {
	return fmt.Sprintf("copy paste this: -account %s -venue %s -stock %s\n", g.Account, g.Venues[0], g.Tickers[0])
}

func (a *API) StopLastGame() {
	gl, err := a.ReadGameLevel()
	if err != nil {
		log.Printf("error reading game level: %s", err)
		return
	}

	log.Printf("stopping: %d\n", gl.InstanceId)
	apiResponse := a.StopLevel(gl.InstanceId)

	// not fatal if game not running
	if !apiResponse.Ok || apiResponse.Error != "" {
		log.Println(apiResponse.Error)
	} else {
		log.Printf("stopped\n")
	}
}

func (a *API) SubmitEvidence(evidence Evidence, instanceId int) (result APIResponse) {
	instanceIdAsString := strconv.Itoa(instanceId)
	url := gmURL + "/instances/" + instanceIdAsString + "/judge"
	a.doPost(url, evidence, &result)

	return
}

func (a *API) WriteGameLevel(gameLevel GameLevel) {
	b, err := json.Marshal(gameLevel)
	Must(err)

	err = ioutil.WriteFile(gameLevelFile, b, 0644)
	Must(err)

	log.Printf("started: %+v\n", gameLevel)
	log.Printf("%s\n", gameLevel.CopyPaste())
}

func (a *API) ReadGameLevel() (result GameLevel, err error) {
	dat, err := ioutil.ReadFile(gameLevelFile)
	if err != nil {
		return
	}

	err = json.Unmarshal(dat, &result)
	return
}
