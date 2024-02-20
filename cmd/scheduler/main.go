package main

import (
	"fmt"
	"gameapp/scheduler"
	"os"
	"os/signal"
)

func main() {
	done := make(chan bool)

	go func() {
		sch := scheduler.New()
		sch.Start(done)
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	fmt.Println("\nreceived Interrupt signal, shutting down gracefully...")
	done <- true
}
