package events

import "time"

func GetTimeForNextDay() (int64, int64, error) {
	location, err := time.LoadLocation("Europe/Moscow")
	if err != nil {
		curTime := time.Now().UTC().Truncate(24 * time.Hour)
		tomorrow := curTime.Add(24 * time.Hour)
		startDay := time.Date(tomorrow.Year(), tomorrow.Month(), tomorrow.Day(), 0, 0, 0, 0, location)
		endDay := time.Date(tomorrow.Year(), tomorrow.Month(), tomorrow.Day(), 23, 59, 59, 0, location)
		return startDay.UnixMilli(), endDay.UnixMilli(), err
	}
	curTime := time.Now().In(location)
	tomorrow := curTime.Add(24 * time.Hour)
	startDay := time.Date(tomorrow.Year(), tomorrow.Month(), tomorrow.Day(), 0, 0, 0, 0, location)
	endDay := time.Date(tomorrow.Year(), tomorrow.Month(), tomorrow.Day(), 23, 59, 59, 0, location)

	return startDay.UnixMilli(), endDay.UnixMilli(), nil
}

func GetTimeFor2Days() (int64, int64, error) {
	location, err := time.LoadLocation("Europe/Moscow")
	if err != nil {
		curTime := time.Now()
		tomorrowTime := curTime.Add(24 * time.Hour)
		startDay := time.Date(curTime.Year(), curTime.Month(), curTime.Day(), 0, 0, 0, 0, time.UTC)
		endDay := time.Date(tomorrowTime.Year(), tomorrowTime.Month(), tomorrowTime.Day(), 23, 59, 59, 0, time.UTC)
		return startDay.UnixMilli(), endDay.UnixMilli(), err
	}

	curTime := time.Now().In(location)
	tomorrowTime := curTime.Add(24 * time.Hour)
	startDay := time.Date(curTime.Year(), curTime.Month(), curTime.Day(), 0, 0, 0, 0, location)
	endDay := time.Date(tomorrowTime.Year(), tomorrowTime.Month(), tomorrowTime.Day(), 23, 59, 59, 0, location)

	return startDay.UnixMilli(), endDay.UnixMilli(), nil
}
