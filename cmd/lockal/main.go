package main

import (
	"encoding/json"
	"fmt"
	"github.com/Frozen-Fantasy/fantasy-backend.git/pkg/models/tournaments"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func main() {
	curTime := time.Now()
	startDay := time.Date(curTime.Year(), curTime.Month(), curTime.Day(), 0, 0, 0, 0, time.UTC)
	endDay := time.Date(curTime.Year(), curTime.Month(), curTime.Day(), 23, 59, 59, 0, time.UTC)

	//fmt.Println(startDay.Unix())
	//fmt.Println(endDay.Unix())
	//url1 := fmt.Sprint("https://khl.api.webcaster.pro/api/khl_mobile/events_v2?q[start_at_lt_time_from_unixtime]=", endDay.Unix(), "&order_direction=desc&q[start_at_gt_time_from_unixtime]=", startDay.Unix())
	//fmt.Println(url1)
	url := fmt.Sprint("https://khl.api.webcaster.pro/api/khl_mobile/events_v2?q[start_at_lt_time_from_unixtime]=", endDay.Unix(), "&order_direction=desc&q[start_at_gt_time_from_unixtime]=", startDay.Unix())
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

	var eventKHL []tournaments.EventDataKHL

	err = decoder.Decode(&eventKHL)
	if err != nil {
		fmt.Println("Error decoding JSON:", err)
		return
	}

	for _, curEnv := range eventKHL {
		fmt.Println(time.UnixMilli(curEnv.Event.EventStartAt))
		fmt.Println(time.UnixMilli(curEnv.Event.EndAt))
		curEnv.Event.TeamA.Score, _ = strconv.Atoi(strings.Split(curEnv.Event.Score, ":")[0])
		curEnv.Event.TeamB.Score, _ = strconv.Atoi(strings.Split(curEnv.Event.Score, ":")[1])

		fmt.Println(curEnv.Event.TeamA.Score, curEnv.Event.TeamB.Score)
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
