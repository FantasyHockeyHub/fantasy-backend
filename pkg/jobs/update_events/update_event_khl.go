package update_events

import (
	"context"
	"github.com/Frozen-Fantasy/fantasy-backend.git/pkg/models/tournaments"
	"github.com/Frozen-Fantasy/fantasy-backend.git/pkg/service/events"
	"log"
	"time"
)

func NewUpdateHockeyEventsKHL(
	//cfg config.ServiceConfiguration, //Потом брать интервал из конфига
	ev *events.EventsService,
) *UpdateHockeyEventsKHL {
	curTime := time.Now()
	tournamentsID := make([]tournaments.ID, 2)
	return &UpdateHockeyEventsKHL{
		dailyGetTime:  time.Date(curTime.Year(), curTime.Month(), curTime.Day(), 9, 0, 0, 0, time.UTC),
		ev:            ev,
		tournamentsID: tournamentsID,
	}

}

type UpdateHockeyEventsKHL struct {
	tournamentsID []tournaments.ID
	dailyGetTime  time.Time
	ev            *events.EventsService
}

func (job *UpdateHockeyEventsKHL) UpdateDurationKHL(ctx context.Context) time.Duration {
	tournInfo, err := job.ev.GetTournamentsByNextDay(ctx, 2)
	if err != nil {
		switch err {
		case events.NotFoundTour:
			if time.Now().After(job.dailyGetTime) {
				job.dailyGetTime = job.dailyGetTime.Add(24 * time.Hour)
			}
		default:
			log.Println("Job GetTournamentsByNextDayKHL:", err)
			if time.Now().After(job.dailyGetTime) {
				job.dailyGetTime = job.dailyGetTime.Add(24 * time.Hour)
			}
		}
	} else {
		job.dailyGetTime = time.UnixMilli(tournInfo[0].TimeStart)
		job.tournamentsID[0] = tournInfo[0].TournamentId
		job.tournamentsID[1] = tournInfo[1].TournamentId
	}

	return job.dailyGetTime.Sub(time.Now())
}

func (job *UpdateHockeyEventsKHL) StartKHL(ctx context.Context) {

	durationTillNextExecution := job.UpdateDurationKHL(ctx)
	//log.Println("KHL", durationTillNextExecution)

	timer := time.NewTimer(durationTillNextExecution)

	for {
		select {
		case <-ctx.Done():
			return
		case <-timer.C:
			err := job.ev.UpdateStatusTournaments(ctx, job.tournamentsID, "started")
			if err != nil {
				log.Println("Job UpdateStatusTournamentsKHL:", err)
			}

			durationNext := job.UpdateDurationKHL(ctx)
			timer.Reset(durationNext)
			//timer.Reset(time.Hour * 24)
		}
	}

}
