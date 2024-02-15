package userservice

import (
	"fmt"
	"gameapp/entity"
	"gameapp/param"
)

func (s Service) Register(req param.RegisterRequest) (param.RegisterResponse, error) {

	// TODO: use bcrypt
	//passwd := []byte(req.Password)
	//bcrypt.GenerateFromPassword(passwd, 0)

	// create new user in storage
	user := entity.User{
		ID:          0,
		Name:        req.Name,
		PhoneNumber: req.PhoneNumber,
		Password:    GetMD5Hash(req.Password),
		Role:        entity.UserRole,
	}
	createdUser, err := s.repo.Register(user)
	if err != nil {
		return param.RegisterResponse{}, fmt.Errorf("unexpected error: %w", err)
	}

	return param.RegisterResponse{param.UserInfo{
		ID:          createdUser.ID,
		PhoneNumber: createdUser.PhoneNumber,
		Name:        createdUser.Name,
	}}, nil
}
