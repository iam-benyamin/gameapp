package main

import (
	"fmt"
	"gameapp/config"
	"gameapp/delivery/httpserver"
	"gameapp/repository/migrator"
	"gameapp/repository/mysql"
	"gameapp/repository/mysql/mysqlaccesscontrol"
	"gameapp/repository/mysql/mysqluser"
	"gameapp/service/authorizationservice"
	"gameapp/service/authservice"
	"gameapp/service/backofficeuserservice"
	"gameapp/service/userservice"
	"gameapp/validator/uservalidator"
)

// TODO: show migration table name
// TODO: add limit to Up and Down
func main() {
	// TODO: read config path from command line
	// TODO: merge  cfg with cfg2exi
	cfg2 := config.Load("config.yml")
	fmt.Printf("cfg2 : %+v\n", cfg2)

	cfg := config.Config{
		HTTPServer: config.HTTPServer{Port: 1986},
		Auth: authservice.Config{
			SignKey:               config.JwtSignKey,
			AccessExpirationTime:  config.AccessTokenExpireDuration,
			RefreshExpirationTime: config.RefreshTokenExpireDuration,
			AccessSubject:         config.AccessTokenSubject,
			RefreshSubject:        config.RefreshTokenSubject,
		},
		Mysql: mysql.Config{
			Username: "gameapp",
			Password: "gameappt0lk2o20",
			Port:     3308,
			Host:     "localhost",
			DBName:   "gameapp_db",
		},
	}

	mgr := migrator.New(cfg.Mysql)
	//mgr.Down()
	mgr.Up()

	authSvc, userSvc, userValidator, backofficeSvc, authorizationSvc := setupServices(cfg)

	server := httpserver.New(cfg, authSvc, userSvc, userValidator, backofficeSvc, authorizationSvc)

	server.Serve()
}

func setupServices(cfg config.Config) (
	authservice.Service, userservice.Service, uservalidator.Validator,
	backofficeuserservice.Service, authorizationservice.Service) {
	authSvc := authservice.New(cfg.Auth)

	MysqlRepo := mysql.New(cfg.Mysql)
	userMysql := mysqluser.New(MysqlRepo)
	userSvc := userservice.New(authSvc, userMysql)
	backofficeUserSvc := backofficeuserservice.New()

	aclMysql := mysqlaccesscontrol.New(MysqlRepo)
	authorizationSvc := authorizationservice.New(aclMysql)

	uV := uservalidator.New(userMysql)

	return authSvc, userSvc, uV, backofficeUserSvc, authorizationSvc
}
