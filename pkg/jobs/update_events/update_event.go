package update_events

import (
	"context"
	"github.com/Frozen-Fantasy/fantasy-backend.git/pkg/models/tournaments"
	"github.com/Frozen-Fantasy/fantasy-backend.git/pkg/service/events"
	"log"
	"time"
)

func NewUpdateHockeyEvents(
	//cfg config.ServiceConfiguration, //Потом брать интервал из конфига
	ev *events.EventsService,
) *UpdateHockeyEvents {
	curTime := time.Now()
	tournamentsID := make([]tournaments.ID, 2)
	return &UpdateHockeyEvents{
		dailyGetTime:  time.Date(curTime.Year(), curTime.Month(), curTime.Day(), 22, 0, 0, 0, time.UTC),
		dailyEndTime:  curTime,
		ev:            ev,
		tournamentsID: tournamentsID,
	}

}

type UpdateHockeyEvents struct {
	tournamentsID []tournaments.ID
	dailyGetTime  time.Time
	dailyEndTime  time.Time
	ev            *events.EventsService
}

func (job *UpdateHockeyEvents) UpdateDuration(ctx context.Context) time.Duration {
	tournInfo, err := job.ev.GetTournamentsByNextDay(ctx, 1)
	if err != nil {
		switch err {
		case events.NotFoundTour:
			if time.Now().After(job.dailyGetTime) {
				job.dailyGetTime = job.dailyGetTime.Add(24 * time.Hour)
				job.tournamentsID[0] = 0
				job.tournamentsID[1] = 0
			}
		default:
			log.Println("Job GetTournamentsByNextDay:", err)
			if time.Now().After(job.dailyGetTime) {
				job.dailyGetTime = job.dailyGetTime.Add(24 * time.Hour)
				job.tournamentsID[0] = 0
				job.tournamentsID[1] = 0
			}
		}
	} else {
		job.dailyGetTime = time.UnixMilli(tournInfo[0].TimeStart)
		//log.Println(job.dailyGetTime)
		job.tournamentsID[0] = tournInfo[0].TournamentId
		//log.Println(job.tournamentsID[0])
		job.tournamentsID[1] = tournInfo[1].TournamentId
		job.dailyEndTime = time.UnixMilli(tournInfo[0].TimeEnd)
	}

	return job.dailyGetTime.Sub(time.Now())
}

func (job *UpdateHockeyEvents) Start(ctx context.Context) {

	durationTillNextExecution := job.UpdateDuration(ctx)
	//log.Println(durationTillNextExecution)

	timer := time.NewTimer(durationTillNextExecution)

	for {
		select {
		case <-ctx.Done():
			return
		case <-timer.C:
			if job.tournamentsID[0] != 0 {
				err := job.ev.UpdateStatusTournaments(ctx, job.tournamentsID)
				if err != nil {
					log.Println("Job UpdateStatusTournaments:", err)
				}
				//запускаем получение данных о матчах каждые 15 минут
				timeEndTour := job.dailyEndTime
				durationTournament := job.dailyGetTime.Sub(timeEndTour)
				ctx2, cancel := context.WithTimeout(ctx, durationTournament)
				job.GetMatchesResult(ctx2, cancel)
			}
			durationNext := job.UpdateDuration(ctx)
			timer.Reset(durationNext)
			//timer.Reset(time.Hour * 24)

		}
	}
}

func (job *UpdateHockeyEvents) GetMatchesResult(ctx context.Context, cancel context.CancelFunc) {
	defer cancel()
	ticker := time.NewTicker(15 * time.Minute)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			err := job.ev.UpdateMatches(ctx)
			if err != nil {
				log.Println("Job UpdateMatches:", err)
			}
		}
	}
}
