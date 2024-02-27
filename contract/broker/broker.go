package broker

import "gameapp/entity"

type Publisher interface {
	Publish(event entity.Event, payload string)
}
