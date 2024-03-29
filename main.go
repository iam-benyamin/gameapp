package main

import (
	"context"
	"fmt"
	presenceClient "gameapp/adapter/presence"
	"gameapp/adapter/redis"
	"gameapp/config"
	"gameapp/delivery/httpserver"
	"gameapp/logger"
	"gameapp/repository/migrator"
	"gameapp/repository/mysql"
	"gameapp/repository/mysql/mysqlaccesscontrol"
	"gameapp/repository/mysql/mysqluser"
	"gameapp/repository/redis/redismatching"
	"gameapp/repository/redis/redispresence"
	"gameapp/scheduler"
	"gameapp/service/authorizationservice"
	"gameapp/service/authservice"
	"gameapp/service/backofficeuserservice"
	"gameapp/service/matchingservice"
	"gameapp/service/presenceservice"
	"gameapp/service/userservice"
	"gameapp/validator/matchingvalidator"
	"gameapp/validator/uservalidator"
	"go.uber.org/zap"
	"net/http"
	"os"
	"os/signal"
	"sync"

	_ "net/http/pprof"
)

func main() {
	go func() {
		// TODO - add enabler config variable
		// curl http://localhost:8099/debug/pprof/goroutine --output goroutine.o
		//  go tool pprof -http=:8086 ./goroutine.o
		http.ListenAndServe(":8099", nil)
	}()

	cfg := config.Load("config.yml")
	//fmt.Printf("cfg : %+v\n", cfg)
	logger.Logger.Named("main").Info("config", zap.Any("config", cfg))

	mgr := migrator.New(cfg.Mysql)
	mgr.Up()
	//mgr.Down()

	// TODO: add struct and these returned items as struct field
	authSvc, userSvc, userValidator, backofficeSvc, authorizationSvc, matchingValidator, matchingSvc, presenceSvc := setupServices(cfg)

	done := make(chan bool)
	var wg sync.WaitGroup
	go func() {
		sch := scheduler.New(matchingSvc, cfg.Scheduler)

		wg.Add(1)
		sch.Start(done, &wg)
	}()

	server := httpserver.New(cfg, authSvc, userSvc, userValidator, backofficeSvc, authorizationSvc, matchingSvc, matchingValidator, presenceSvc)
	go func() {
		server.Serve()
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	fmt.Println("\nreceived Interrupt signal, shutting down gracefully...")

	ctx := context.Background()
	ctxWithTimeout, cancel := context.WithTimeout(ctx, cfg.Application.GracefulShutdownTimeout)
	defer cancel()

	if err := server.Router.Shutdown(ctxWithTimeout); err != nil {
		// TODO: replace all fmt.Print.. and std log calls with logger.Logger
		fmt.Println("http server shutdown error", err)
	}

	done <- true

	// TODO: dose order of ctx.Done & wg.Waite matter?
	<-ctxWithTimeout.Done()
	wg.Wait()
}

func setupServices(cfg config.Config) (
	authservice.Service, userservice.Service, uservalidator.Validator,
	backofficeuserservice.Service, authorizationservice.Service,
	matchingvalidator.Validator, matchingservice.Service,
	presenceservice.Service) {

	authSvc := authservice.New(cfg.Auth)

	MysqlRepo := mysql.New(cfg.Mysql)
	userMysql := mysqluser.New(MysqlRepo)
	userSvc := userservice.New(authSvc, userMysql)
	backofficeUserSvc := backofficeuserservice.New()

	aclMysql := mysqlaccesscontrol.New(MysqlRepo)
	authorizationSvc := authorizationservice.New(aclMysql)

	uV := uservalidator.New(userMysql)

	matchingV := matchingvalidator.New()

	redisAdapter := redis.New(cfg.Redis)

	presencesRepo := redispresence.New(redisAdapter)
	presenceSvc := presenceservice.New(cfg.PresenceService, presencesRepo)

	matchingRepo := redismatching.New(redisAdapter)

	// TODO: add address to config
	presenceAdapter := presenceClient.New(":8086")

	matchingSvc := matchingservice.New(cfg.MatchingService, matchingRepo, presenceAdapter, redisAdapter)

	return authSvc, userSvc, uV, backofficeUserSvc, authorizationSvc, matchingV, matchingSvc, presenceSvc
}
