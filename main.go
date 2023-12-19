package main

import (
	"encoding/json"
	"fmt"
	"gameapp/repository/mysql"
	"gameapp/service/userservice"
	"io"
	"log"
	"net/http"
)

func main() {

	mux := http.NewServeMux()
	mux.HandleFunc("/heath-check/", healthCheckHandler)
	mux.HandleFunc("/users/register", userRegisterHandler)
	mux.HandleFunc("/users/login", userLoginHandler)

	log.Println("server is lessening on localhost:1986")
	err := http.ListenAndServe(":1986", mux)
	if err != nil {
		panic(err)
	}
}

func userRegisterHandler(writer http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		fmt.Fprintf(writer, `{"error1": "invalid method"}`)

		return
	}

	data, err := io.ReadAll(req.Body)
	if err != nil {
		writer.Write([]byte(fmt.Sprintf(`{"error2": "%s"}`, err.Error())))

		return
	}

	var uReq userservice.RegisterRequest
	err = json.Unmarshal(data, &uReq)
	if err != nil {
		writer.Write([]byte(fmt.Sprintf(`{"error3": "%s"}`, err.Error())))

		return
	}

	mysqlRepo := mysql.New()
	userSvc := userservice.New(mysqlRepo)

	_, err = userSvc.Register(uReq)
	if err != nil {
		writer.Write([]byte(fmt.Sprintf(`{"error4": "%s"}`, err.Error())))

		return
	}

	writer.Write([]byte(fmt.Sprintf(`{"message": "user created"}`)))
}

func healthCheckHandler(writer http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(writer, `{"message": "every thing is fine!\n"}`)
}

func userLoginHandler(writer http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		fmt.Fprintf(writer, `{"error1": "invalid method"}`)

		return
	}

	data, err := io.ReadAll(req.Body)
	if err != nil {
		writer.Write([]byte(fmt.Sprintf(`{"error2": "%s"}`, err.Error())))

		return
	}

	var lReq userservice.LoginRequest
	err = json.Unmarshal(data, &lReq)
	if err != nil {
		writer.Write([]byte(fmt.Sprintf(`{"error3": "%s"}`, err.Error())))

		return
	}

	mysqlRepo := mysql.New()
	userSvc := userservice.New(mysqlRepo)

	_, err = userSvc.Login(lReq)
	if err != nil {
		writer.Write([]byte(fmt.Sprintf(`{"error4": "%s"}`, err.Error())))

		return
	}

	writer.Write([]byte(fmt.Sprintf(`{"message": "user credentional is okay"}`)))
}