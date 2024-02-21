package redispresence

import "gameapp/adapter/redis"

type DB struct {
	adapter redis.Adapter
}

func New(adapter redis.Adapter) DB {
	return DB{adapter: adapter}
}
