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
	location, _ := time.LoadLocation("Europe/Moscow")
	return &UpdateHockeyEvents{
		dailyGetTime:  time.Date(curTime.Year(), curTime.Month(), curTime.Day(), 22, 0, 0, 0, location),
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
	//tournInfo1 := tournaments.Tournament{TournamentId: 523294174, TimeStart: 1715532060000, TimeEnd: 1715532180000}
	//tournInfo2 := tournaments.Tournament{TournamentId: 1631395586, TimeStart: 1715532060000, TimeEnd: 1715532180000}
	//tournInfo := []tournaments.Tournament{tournInfo1, tournInfo2}
	//var err error
	//err = nil
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
		job.tournamentsID[0] = tournInfo[0].TournamentId
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
			var durationTournament time.Duration
			if job.tournamentsID[0] != 0 {
				err := job.ev.UpdateStatusTournaments(ctx, job.tournamentsID, "started")
				if err != nil {
					log.Println("Job UpdateStatusTournaments:", err)
				}
				durationTournament = job.dailyEndTime.Sub(job.dailyGetTime)
			}

			durationNext := job.UpdateDuration(ctx)
			timer.Reset(durationNext)
			//timer.Reset(time.Hour * 24)

			if durationTournament != 0 {
				//запускаем получение данных о матчах каждые 15 минут
				ctx2, cancel := context.WithTimeout(ctx, durationTournament)
				job.GetMatchesResult(ctx2, cancel)

				err := job.ev.UpdateStatusTournaments(ctx, job.tournamentsID, "finished")
				if err != nil {
					log.Println("Job UpdateStatusTournaments:", err)
				}
			}

		}
	}
}

func (job *UpdateHockeyEvents) GetMatchesResult(ctx context.Context, cancel context.CancelFunc) {
	defer cancel()
	ticker := time.NewTicker(10 * time.Minute)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			err := job.ev.UpdateMatches(ctx, job.tournamentsID)
			if err != nil {
				log.Println("Job UpdateMatches:", err)
			}
		}
	}
}
