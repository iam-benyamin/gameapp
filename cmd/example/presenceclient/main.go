package main

import (
	"context"
	"fmt"
	presenceClient "gameapp/adapter/presence"
	"gameapp/param"
	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial(":8086", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	client := presenceClient.New(conn)
	//client := presence.NewPresenceServiceClient(conn)
	//resp, err := client.GetPresence(context.Background(), &presence.GetPresenceRequest{UserIds: []uint64{1, 2, 3, 4}})
	resp, err := client.GetPresence(context.Background(), param.GetPresenceRequest{UserIDs: []uint{1, 2, 3, 4, 5, 6}})
	if err != nil {
		panic(err)
	}

	for _, item := range resp.Items {
		fmt.Println("item", item.UserID, item.Timestamp)
	}
}
