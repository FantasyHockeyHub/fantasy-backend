package main

import (
	"encoding/json"
	"fmt"
	"github.com/Frozen-Fantasy/fantasy-backend.git/pkg/models/tournaments"
	"net/http"
	"time"
)

func main() {
	curTime := time.Now()
	curTime = curTime.Add(-24 * time.Hour)

	url := fmt.Sprint("https://api-web.nhle.com/v1/schedule/", curTime.Format("2006-01-02"))
	fmt.Println(url)

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

	var eventNHL tournaments.ScheduleNHL

	err = decoder.Decode(&eventNHL)
	if err != nil {
		fmt.Println("Error decoding JSON:", err)
		return
	}

	for _, curEnv := range eventNHL.GameWeeks[1].Games {
		startTime, err := time.Parse("2006-01-02T15:04:05Z", curEnv.StartTimeUTC)
		if err != nil {
			fmt.Println("Ошибка при парсинге времени:", err)
			return
		}

		// Выводим результат
		fmt.Println("Время в формате time.Time:", startTime.UnixMilli())
		fmt.Println("End time:", startTime.Add(3*time.Hour).UnixMilli())
		fmt.Printf("Away %d Home %d", curEnv.AwayTeam.Score, curEnv.HomeTeam.Score)
	}

	//fmt.Println(len(standings.Standings))
	//for _, curStand := range standings.Standings {
	//	fmt.Printf("Team Name: %s\n", curStand.TeamName.Default)
	//	fmt.Printf("Team Abbrev: %s\n", curStand.TeamAbbrev.Default)
	//	fmt.Printf("Conference Name: %s\n", curStand.ConferenceName)
	//	fmt.Printf("Logo: %s\n", curStand.TeamLogo)
	//	fmt.Println("-------------------")
	//}

	// Текущей UNIX timestamp в секундах

	//apochNow := time.Now().Unix()
	//fmt.Printf("Время эпохи в секундах: %d\n", apochNow)
	//
	//apochNano := time.Now().UnixNano()
	//fmt.Printf("Время эпохи в наносекундах: %d\n", apochNano)
	//
	//date := time.Date(2024, 2, 2, 0, 0, 0, 0, time.UTC)
	//fmt.Println(date.Unix())
	//fmt.Println(date.UnixMilli())
	//fmt.Println(date.UnixMicro())
	//fmt.Println(date.UnixNano())

}
