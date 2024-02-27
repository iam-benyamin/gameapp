package redis

import (
	"context"
	"fmt"
	"gameapp/entity"
	"github.com/labstack/gommon/log"
	"time"
)

func (a Adapter) Publish(event entity.Event, payload string) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := a.client.Publish(ctx, string(event), payload).Err(); err != nil {
		log.Errorf(fmt.Sprintf("publish err : %v", err))
		// TODO: log
		// TODO: update metrics
	}

	// TODO: update metrics
}
