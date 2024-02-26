package main

import (
	"context"
	"fmt"
	"gameapp/adapter/redis"
	"gameapp/config"
	"gameapp/entity"
	"gameapp/pkg/protobufencoder"
)

func main() {
	cfg := config.Load("config.yml")

	redisAdapter := redis.New(cfg.Redis)

	topic := entity.MatchingUsersMatchedEvent
	mu := entity.MatchedUsers{
		Category: entity.FootballCategory,
		UserIDs:  []uint{1, 4},
	}

	payload := protobufencoder.EncodeMatchingUsersMatchedEvent(mu)
	if err := redisAdapter.Client().Publish(context.Background(), string(topic), payload).Err(); err != nil {
		panic(fmt.Sprintf("publish err : %v", err))
	}
}
