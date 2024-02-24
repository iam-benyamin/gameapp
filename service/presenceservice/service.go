package presenceservice

import (
	"context"
	"fmt"
	"gameapp/param"
	"gameapp/pkg/richerror"
	"time"
)

type Config struct {
	ExpirationTime time.Duration `koanf:"expiration_time"`
	Prefix         string        `koanf:"prefix"`
}

type Repo interface {
	Upsert(ctx context.Context, key string, timestamp int64, expirationTime time.Duration) error
}

type Service struct {
	config Config
	repo   Repo
}

func New(config Config, repo Repo) Service {
	return Service{config: config, repo: repo}
}

func (s Service) UpsertPresence(ctx context.Context, req param.UpsertPresenceRequest) (param.UpsertPresenceResponse, error) {
	const op = "presencesservice.UpsertPresence"

	err := s.repo.Upsert(ctx, fmt.Sprintf("%s:%d", s.config.Prefix, req.UserID), req.Timestamp, s.config.ExpirationTime)
	if err != nil {
		return param.UpsertPresenceResponse{}, richerror.New(op).WithErr(err)
	}

	return param.UpsertPresenceResponse{}, nil
}

func (s Service) GetPresence(ctx context.Context, request param.GetPresenceRequest) (param.GetPresenceResponse, error) {
	fmt.Println("req: ", request)
	// TODO: implement me
	return param.GetPresenceResponse{Items: []param.GetPresenceItem{
		{UserID: 1, Timestamp: 123452151},
		{UserID: 2, Timestamp: 234234512},
	}}, nil
}
