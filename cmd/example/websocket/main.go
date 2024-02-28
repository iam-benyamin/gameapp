package main

import (
	"fmt"
	"gameapp/entity"
	"gameapp/pkg/protobufencoder"
	"github.com/labstack/gommon/log"
	"net"
	"net/http"

	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
)

func printDecodedNotification() {
	d := protobufencoder.EncodeNotification(entity.Notification{
		EventType: "ping",
		Payload:   "hello",
	})

	fmt.Println("d", d)
}

func main() {
	printDecodedNotification()

	http.ListenAndServe(":8080", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, _, _, err := ws.UpgradeHTTP(r, w)
		if err != nil {
			// handle error
		}

		defer conn.Close()
		done := make(chan bool)

		go readMessage(conn, done)

		<-done

		//go func() {
		//	defer conn.Close()
		//
		//	for {
		//		msg, op, err := wsutil.ReadClientData(conn)
		//		if err != nil {
		//			fmt.Println("read client err :", err)
		//		}
		//		err = wsutil.WriteServerMessage(conn, op, msg)
		//		if err != nil {
		//			fmt.Println("server message err: ", err)
		//		}
		//		fmt.Println("msg is :", string(msg))
		//		fmt.Println("op is :", string(op))
		//	}
		//}()
	}))
}

func readMessage(conn net.Conn, done chan<- bool) {
	for {
		msg, opCode, err := wsutil.ReadClientData(conn)
		if err != nil {
			log.Print(err)
			done <- true
		}
		fmt.Println("msg", string(msg))
		notif := protobufencoder.DecodeNotification(string(msg))

		fmt.Println("opCode", opCode)
		fmt.Println("notif", notif)
	}
}
