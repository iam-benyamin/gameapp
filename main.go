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

	writer.Write([]byte(fmt.Sprintf(`{"message": "okay"}`)))
}

func healthCheckHandler(writer http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(writer, `{"message": "every thing is fine!\n"}`)
}

func test() {
	//mysqlRepo := mysql.New()

	//createdUser, err := mysqlRepo.Register(entity.User{
	//	Name:        "Ali Reza",
	//	PhoneNumber: "0923",
	//})
	//if err != nil {
	//	fmt.Println("register user : ", err)
	//}
	//fmt.Println("created user : ", createdUser)

	//	isUnique, err := mysqlRepo.IsPhoneNumberUnique("0922")
	//	if err != nil {
	//		fmt.Println("unique err: ", err)
	//	}
	//	fmt.Println("is unique : ", isUnique)
}
