package main

import (
	"encoding/json"
	"fmt"
	"gameapp/entity"
	"github.com/labstack/gommon/log"
	"net"
	"net/http"
	"time"

	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
)

func producer(remoteAddr string, channel chan string) {
	for {
		channel <- remoteAddr
		time.Sleep(5 * time.Second)
	}
}

func main() {
	http.ListenAndServe(":8080", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, _, _, err := ws.UpgradeHTTP(r, w)
		if err != nil {
			// handle error
		}

		defer conn.Close()
		done := make(chan bool)

		go readMessage(conn, done)

		channel := make(chan string)
		go producer(r.RemoteAddr, channel)
		go writeMessage(conn, channel)

		<-done
	}))
}

func readMessage(conn net.Conn, done chan<- bool) {
	for {
		msg, opCode, err := wsutil.ReadClientData(conn)
		if err != nil {
			log.Print(err)
			done <- true

			return
		}

		var notif entity.Notification
		//notif := protobufencoder.DecodeNotification(string(msg))
		if err := json.Unmarshal(msg, &notif); err != nil {
			panic(err)
		}

		fmt.Println("opCode", opCode)
		fmt.Println("notif", notif)
	}
}

func writeMessage(conn net.Conn, channel <-chan string) {
	for data := range channel {
		if err := wsutil.WriteServerMessage(conn, ws.OpText, []byte(data)); err != nil {
			panic(err)
		}
	}
}
