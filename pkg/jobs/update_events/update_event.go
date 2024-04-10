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
		ev:            ev,
		tournamentsID: tournamentsID,
	}

}

type UpdateHockeyEvents struct {
	tournamentsID []tournaments.ID
	dailyGetTime  time.Time
	ev            *events.EventsService
}

func (job *UpdateHockeyEvents) UpdateDuration(ctx context.Context) time.Duration {
	tournInfo, err := job.ev.GetTournamentsByNextDay(ctx, 1)
	if err != nil {
		switch err {
		case events.NotFoundTour:
			if time.Now().After(job.dailyGetTime) {
				job.dailyGetTime = job.dailyGetTime.Add(24 * time.Hour)
			}
		default:
			log.Println("Job GetTournamentsByNextDay:", err)
			if time.Now().After(job.dailyGetTime) {
				job.dailyGetTime = job.dailyGetTime.Add(24 * time.Hour)
			}
		}
	} else {
		job.dailyGetTime = time.UnixMilli(tournInfo[0].TimeStart)
		//log.Println(job.dailyGetTime)
		job.tournamentsID[0] = tournInfo[0].TournamentId
		//log.Println(job.tournamentsID[0])
		job.tournamentsID[1] = tournInfo[1].TournamentId
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
			err := job.ev.UpdateStatusTournaments(ctx, job.tournamentsID)
			if err != nil {
				log.Println("Job UpdateStatusTournaments:", err)
			}

			durationNext := job.UpdateDuration(ctx)
			timer.Reset(durationNext)
			//timer.Reset(time.Hour * 24)
		}
	}

}
