package main

import (
	"flag"
	"log"
	"os"

	"github.com/armhold/sfapi"
)

var (
	stop              = flag.Bool("stop", false, "-stop")
	start3            = flag.Bool("start3", false, "-start3")
	start4            = flag.Bool("start4", false, "-start4")
	start5            = flag.Bool("start5", false, "-start5")
	start6            = flag.Bool("start6", false, "-start6")
	submitEvidence    = flag.Bool("evidence", false, "-evidence --instance 12345 --account ABC1232 --link  http://example.com --summary \"executive summary yada yada...\"")
	evidenceInstance  = flag.Int("instance", -1, "")
	evidenceAccount   = flag.String("account", "", "")
	evidenceLink      = flag.String("link", "", "")
	evidenceSummary   = flag.String("summary", "", "")
	restartInstanceId = flag.Int("restart", -1, "-restart INSTANCE_ID")
	authKey           = os.Getenv("STARFIGHTER_API_KEY")
	a                 *sfapi.API
)

/*
 A somewhat clunky way to stop/start the various levels, and save the game info to /tmp/gamelevel.
 This obviates the need to copy/paste the account/venue/symbol by hand each time you restart your client.
*/

func main() {
	if authKey == "" {
		log.Fatal(os.Stderr, "STARFIGHTER_API_KEY environment variable not set")
	}

	flag.Parse()

	a = sfapi.NewAPI()

	if *stop {
		a.StopLastGame()
	} else if *start3 {
		gameLevel := a.StartLevel3()
		a.WriteGameLevel(gameLevel)
	} else if *start4 {
		gameLevel := a.StartLevel4()
		a.WriteGameLevel(gameLevel)
	} else if *start5 {
		gameLevel := a.StartLevel5()
		a.WriteGameLevel(gameLevel)
	} else if *start6 {
		gameLevel := a.StartLevel6()
		a.WriteGameLevel(gameLevel)
	} else if *restartInstanceId != -1 {
		gameLevel := a.RestartLevel(*restartInstanceId)
		log.Printf("GameLevel: %+v\n", gameLevel)
		log.Printf("%s\n", gameLevel.CopyPaste())
	} else if *submitEvidence {
		log.Printf("submit evidence, instance: %d, account: \"%s\", link: %s, summary: %s\n", *evidenceInstance, *evidenceAccount, *evidenceLink, *evidenceSummary)
		ev := sfapi.Evidence{Account: *evidenceAccount, Link: *evidenceLink, Summary: *evidenceSummary}

		result := a.SubmitEvidence(ev, *evidenceInstance)
		log.Printf("result of submission: %+v\n", result)
	} else {
		flag.Usage()
		os.Exit(1)
	}
}
