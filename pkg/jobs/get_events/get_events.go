package get_events

import (
	"context"
	"github.com/Frozen-Fantasy/fantasy-backend.git/pkg/service"
	"log"
	"time"
)

func NewGetHokeyEvents(
	//cfg config.ServiceConfiguration, //Потом брать интервал из конфига
	events *service.Services,
) *GetHokeyEvents {
	curTime := time.Now()
	return &GetHokeyEvents{
		dailyGetTime: time.Date(curTime.Year(), curTime.Month(), curTime.Day(), 0, 0, 1, 0, time.Local),
		events:       events,
	}

}

type GetHokeyEvents struct {
	dailyGetTime time.Time
	events       *service.Services
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
			err := job.events.Teams.AddEventsKHL(ctx)
			if err != nil {
				log.Println("Job AddEventsKHL:", err)
			}
			timer.Reset(24 * time.Hour)
		}
	}
	//err := job.events.Teams.AddEventsKHL(ctx)
	//if err != nil {
	//	log.Println("Job AddEventsKHL:", err)
	//}

}
