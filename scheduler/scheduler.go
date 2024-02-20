package scheduler

import (
	"fmt"
	"gameapp/param"
	"gameapp/service/matchingservice"
	"github.com/go-co-op/gocron"
	"sync"
	"time"
)

type Scheduler struct {
	sch      *gocron.Scheduler
	matchSvc matchingservice.Service
}

func New(matchSvc matchingservice.Service) Scheduler {
	return Scheduler{
		sch:      gocron.NewScheduler(time.UTC),
		matchSvc: matchSvc,
	}
}

// Start long-running process
func (s Scheduler) Start(done <-chan bool, wg *sync.WaitGroup) {
	defer wg.Done()

	s.sch.Every(5).Second().Do(s.MatchWaitedUser)

	s.sch.StartAsync()

	<-done
	fmt.Println("stopping scheduler...")
	s.sch.Stop()

}

func (s Scheduler) MatchWaitedUser() {
	s.matchSvc.MatchWaitedUsers(param.MatchWaitedUsersRequest{})
}
