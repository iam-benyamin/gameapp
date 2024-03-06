package main

import (
	"fmt"
	presenceClient "gameapp/adapter/presence"
	"gameapp/adapter/redis"
	"gameapp/config"
	"gameapp/repository/redis/redismatching"
	"gameapp/scheduler"
	"gameapp/service/matchingservice"
	"os"
	"os/signal"
	"sync"
)

func main() {
	cfg := config.Load("config.yml")

	redisAdapter := redis.New(cfg.Redis)
	matchingRepo := redismatching.New(redisAdapter)
	presenceAdapter := presenceClient.New(":8086")

	matchingSvc := matchingservice.New(cfg.MatchingService, matchingRepo, presenceAdapter, redisAdapter)

	done := make(chan bool)
	var wg sync.WaitGroup

	go func() {
		sch := scheduler.New(matchingSvc, cfg.Scheduler)

		wg.Add(1)
		sch.Start(done, &wg)
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	fmt.Println("\nreceived Interrupt signal, shutting down gracefully...")
	done <- true
}
