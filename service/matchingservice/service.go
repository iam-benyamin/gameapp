package matchingservice

import (
	"gameapp/entity"
	"gameapp/param"
	"gameapp/pkg/richerror"
	"time"
)

type Repo interface {
	AddToWaitingList(UserID uint, category entity.Category) error
}

type Config struct {
	WaitingTimeout time.Duration `koanf:"waiting_timeout"`
}

type Service struct {
	config Config
	repo   Repo
}

func New(config Config, repo Repo) Service {
	return Service{config: config, repo: repo}
}

func (s Service) AddToWaitingList(req param.AddToWaitingListRequest) (param.AddToWaitingListResponse, error) {
	const op = richerror.Op("matchingservice.AddToWaitingList")

	// add user to the waiting list for the given category if not exists
	// also we can update the waiting timestamp

	err := s.repo.AddToWaitingList(req.UserID, req.Category)
	if err != nil {
		return param.AddToWaitingListResponse{}, richerror.New(op).WithErr(err).WithKind(richerror.KindUnexpected)
	}

	return param.AddToWaitingListResponse{Timeout: s.config.WaitingTimeout}, nil
}

func (s Service) MatchWaitedUsers(req param.MatchWaitedUsersRequest) (param.MatchWaitedUsersResponse, error) {
	return param.MatchWaitedUsersResponse{}, nil
}
