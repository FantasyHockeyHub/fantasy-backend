package main

import (
	"encoding/json"
	"fmt"
	"github.com/Frozen-Fantasy/fantasy-backend.git/pkg/models/tournaments"
	"net/http"
)

func main() {
	url := "https://api-web.nhle.com/v1/standings/now"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Print(err.Error())
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Print(err.Error())
	}
	defer res.Body.Close()
	decoder := json.NewDecoder(res.Body)

	var standings tournaments.StandingsResponse

	err = decoder.Decode(&standings)
	if err != nil {
		fmt.Println("Error decoding JSON:", err)
		return
	}
	//bodyByte, readErr := ioutil.ReadAll(res.Body)
	//if readErr != nil {
	//	fmt.Print(err.Error())
	//}
	fmt.Println(len(standings.Standings))
	for _, curStand := range standings.Standings {
		fmt.Printf("Team Name: %s\n", curStand.TeamName.Default)
		fmt.Printf("Team Abbrev: %s\n", curStand.TeamAbbrev.Default)
		fmt.Printf("Conference Name: %s\n", curStand.ConferenceName)
		fmt.Printf("Logo: %s\n", curStand.TeamLogo)
		fmt.Println("-------------------")
	}

}
