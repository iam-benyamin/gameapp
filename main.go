package main

import (
	"fmt"
	"gameapp/repository/mysql"
)

func main() {
	mysqlRepo := mysql.New()

	//createdUser, err := mysqlRepo.Register(entity.User{
	//	Name:        "Ali Reza",
	//	PhoneNumber: "0923",
	//})
	//if err != nil {
	//	fmt.Println("register user : ", err)
	//}
	//fmt.Println("created user : ", createdUser)

	isUnique, err := mysqlRepo.IsPhoneNumberUnique("0922")
	if err != nil {
		fmt.Println("unique err: ", err)
	}
	fmt.Println("is unique : ", isUnique)
}
