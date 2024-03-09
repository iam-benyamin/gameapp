package main

import (
	"context"
	"fmt"
	presenceClient "gameapp/adapter/presence"
	"gameapp/param"
	"google.golang.org/grpc"
)

func main() {
	// TODO: fix grpc server bug with lesson 22 fix bug and this commit https://github.com/gocasts-bootcamp/gameapp/commit/3482896ba71b044da19c8acb397df308d4c893d0
	conn, err := grpc.Dial(":8086", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	client := presenceClient.New("8086")
	//client := presence.NewPresenceServiceClient(conn)
	//resp, err := client.GetPresence(context.Background(), &presence.GetPresenceRequest{UserIds: []uint64{1, 2, 3, 4}})
	resp, err := client.GetPresence(context.Background(), param.GetPresenceRequest{UserIDs: []uint{1, 2, 3, 4, 5, 6}})
	if err != nil {
		fmt.Println(err)
	}

	for _, item := range resp.Items {
		fmt.Println("item", item.UserID, item.Timestamp)
	}
}
