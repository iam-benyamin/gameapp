package main

import (
	"gameapp/config"
	"gameapp/delivery/httpserver"
	"gameapp/repository/mysql"
	"gameapp/service/authservice"
	"gameapp/service/userservice"
	"time"
)

const (
	JwtSignKey                 = "jwt_secret"
	AccessTokenSubject         = "ac"
	RefreshTokenSubject        = "rt"
	AccessTokenExpireDuration  = time.Hour * 24
	RefreshTokenExpireDuration = time.Hour * 24 * 7
)

func main() {
	cfg := config.Config{
		HTTPServer: config.HTTPServer{Port: 1986},
		Auth: authservice.Config{
			SignKey:               JwtSignKey,
			AccessExpirationTime:  AccessTokenExpireDuration,
			RefreshExpirationTime: RefreshTokenExpireDuration,
			AccessSubject:         AccessTokenSubject,
			RefreshSubject:        RefreshTokenSubject,
		},
		Mysql: mysql.Config{
			Username: "gameapp",
			Password: "gameappt0lk2o20",
			Port:     3308,
			Host:     "localhost",
			DBName:   "gameapp_db",
		},
	}
	authSvc, userSvc := setupServices(cfg)

	server := httpserver.New(cfg, authSvc, userSvc)
	server.Serve()
}

//func userLoginHandler(writer http.ResponseWriter, req *http.Request) {
//	writer.Header().Add("Content-Type", "application/json")
//
//	if req.Method != http.MethodPost {
//		fmt.Fprintf(writer, `{"error": "invalid method"}`)
//
//		return
//	}
//
//	data, err := io.ReadAll(req.Body)
//	if err != nil {
//		writer.Write([]byte(fmt.Sprintf(`{"error": "%s"}`, err.Error())))
//
//		return
//	}
//
//	var lReq userservice.LoginRequest
//	err = json.Unmarshal(data, &lReq)
//	if err != nil {
//		writer.Write([]byte(fmt.Sprintf(`{"error": "%s"}`, err.Error())))
//
//		return
//	}
//
//	authSvc := authservice.New(JwtSignKey, AccessTokenSubject, RefreshTokenSubject,
//		AccessTokenExpireDuration, RefreshTokenExpireDuration)
//
//	mysqlRepo := mysql.New()
//	userSvc := userservice.New(authSvc, mysqlRepo)
//
//	resp, err := userSvc.Login(lReq)
//	if err != nil {
//		writer.WriteHeader(http.StatusBadRequest)
//		writer.Write([]byte(fmt.Sprintf(`{"error": "%s"}`, err.Error())))
//
//		return
//	}
//
//	data, err = json.Marshal(resp)
//	if err != nil {
//		writer.Write([]byte(fmt.Sprintf(`{"error": "%s"}`, err.Error())))
//
//		return
//	}
//
//	//writer.Header().Add("Content-Type", "application/json")
//	writer.Write(data)
//}
//
//func userProfileHandler(writer http.ResponseWriter, req *http.Request) {
//	if req.Method != http.MethodGet {
//		fmt.Fprintf(writer, `{"error": "invalid method"}`)
//
//		return
//	}
//
//	authSvc := authservice.New(JwtSignKey, AccessTokenSubject, RefreshTokenSubject,
//		AccessTokenExpireDuration, RefreshTokenExpireDuration)
//
//	authToken := req.Header.Get("Authorization")
//	claims, err := authSvc.ParseToken(authToken)
//	if err != nil {
//		fmt.Fprintf(writer, `{"error": "token is not valid"}`)
//	}
//
//	mysqlRepo := mysql.New()
//	userSvc := userservice.New(authSvc, mysqlRepo)
//
//	resp, err := userSvc.Profile(userservice.ProfileRequest{UserID: claims.UserID})
//	if err != nil {
//		writer.Write([]byte(fmt.Sprintf(`{"error": "%s"}`, err.Error())))
//
//		return
//	}
//
//	data, err := json.Marshal(resp)
//	if err != nil {
//		writer.Write([]byte(fmt.Sprintf(`{"error": "%s"}`, err.Error())))
//
//		return
//	}
//
//	writer.Write(data)
//}

func setupServices(cfg config.Config) (authservice.Service, userservice.Service) {
	authSvc := authservice.New(cfg.Auth)

	MysqlRepo := mysql.New(cfg.Mysql)
	UserSvc := userservice.New(authSvc, MysqlRepo)

	return authSvc, UserSvc
}
