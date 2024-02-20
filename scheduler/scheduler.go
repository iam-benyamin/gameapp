package scheduler

import (
	"fmt"
	"time"
)

// Scheduler is long-running process
type Scheduler struct{}

func New() Scheduler {
	return Scheduler{}
}

func (s Scheduler) Start(done <-chan bool) {
	fmt.Println("start Scheduler")

	for {
		select {
		case <-done:
			fmt.Println("stopping schedulers..")
			return
		default:
			now := time.Now()
			fmt.Println("scheduler now: ", now)
			time.Sleep(3 * time.Second)
		}
	}
}
