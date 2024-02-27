package redispresence

import (
	"context"
	"gameapp/pkg/richerror"
	"gameapp/pkg/timestamp"
	"time"
)

func (d DB) Upsert(ctx context.Context, key string, timestamp int64, expirationTime time.Duration) error {
	const op = "redispresence.Upsert"
	_, err := d.adapter.Client().Set(ctx, key, timestamp, expirationTime).Result()
	if err != nil {
		return richerror.New(op).WithErr(err).WithKind(richerror.KindUnexpected)
	}

	return nil
}

func (d DB) GetPresence(ctx context.Context, prefixKey string, userIDs []uint) (map[uint]int64, error) {
	// TODO: implement me
	// TODO: how to gte multiple redis key at same time
	m := make(map[uint]int64)

	for _, u := range userIDs {
		m[u] = timestamp.Add(-100 * time.Millisecond)
	}

	return m, nil
}
