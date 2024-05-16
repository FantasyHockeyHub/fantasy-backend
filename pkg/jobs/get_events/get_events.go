package get_events

import (
	"context"
	"github.com/Frozen-Fantasy/fantasy-backend.git/pkg/service/events"
	"log"
	"time"
)

func NewGetHokeyEvents(
	//cfg config.ServiceConfiguration, //Потом брать интервал из конфига
	//events *service.Services,
	ev *events.EventsService,
) *GetHokeyEvents {
	curTime := time.Now()
	return &GetHokeyEvents{
		dailyGetTime: time.Date(curTime.Year(), curTime.Month(), curTime.Day(), 1, 0, 0, 0, time.Local),
		//events:       events,
		ev: ev,
	}

}

type GetHokeyEvents struct {
	dailyGetTime time.Time
	//events       *service.Services
	ev *events.EventsService
}

func (job *GetHokeyEvents) Start(ctx context.Context) {

	if time.Now().After(job.dailyGetTime) {
		job.dailyGetTime = job.dailyGetTime.Add(24 * time.Hour)
	}

	durationTillNextExecution := job.dailyGetTime.Sub(time.Now())

	timer := time.NewTimer(durationTillNextExecution)

	for {
		select {
		case <-ctx.Done():
			return
		case <-timer.C:
			err := job.ev.AddEventsKHL(ctx)
			//err := job.events.Teams.AddEventsKHL(ctx)
			if err != nil {
				log.Println("Job AddEventsKHL:", err)
			}
			err = job.ev.AddEventsNHL(ctx)
			if err != nil {
				log.Println("Job AddEventsNHL:", err)
			}
			err = job.ev.CreateTournaments(ctx)
			if err != nil {
				log.Println("Job CreateTournaments:", err)
			}

			timer.Reset(24 * time.Hour)
		}
	}
	//err := job.events.Teams.AddEventsKHL(ctx)
	//if err != nil {
	//	log.Println("Job AddEventsKHL:", err)
	//}

}
