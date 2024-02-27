package scheduler

import (
	"context"
	"fmt"
	"gameapp/param"
	"gameapp/service/matchingservice"
	"github.com/go-co-op/gocron"
	"sync"
	"time"
)

type Config struct {
	MatchWaitedUserIntervalInSeconds int64 `koanf:"match_waited_user_interval_in_seconds"`
}

// Scheduler long-running process
type Scheduler struct {
	sch      *gocron.Scheduler
	matchSvc matchingservice.Service
	config   Config
}

func New(matchSvc matchingservice.Service, config Config) Scheduler {
	return Scheduler{
		sch:      gocron.NewScheduler(time.UTC),
		matchSvc: matchSvc,
		config:   config,
	}
}

func (s Scheduler) Start(done <-chan bool, wg *sync.WaitGroup) {
	defer wg.Done()

	_, err := s.sch.Every(int(s.config.MatchWaitedUserIntervalInSeconds)).Second().Do(s.MatchWaitedUser)
	if err != nil {
		fmt.Println("scheduler error", err)
	}

	s.sch.StartAsync()

	<-done
	fmt.Println("stopping scheduler...")
	s.sch.Stop()

}

func (s Scheduler) MatchWaitedUser() {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	_, err := s.matchSvc.MatchWaitedUsers(ctx, param.MatchWaitedUsersRequest{})
	if err != nil {
		// TODO: log error
		// TODO: update metrics
		fmt.Println("matchSvc.MatchWaitedUser error : ", err)
	}
}
