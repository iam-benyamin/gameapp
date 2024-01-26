package userservice

import (
	"fmt"
	"gameapp/dto"
	"gameapp/entity"
)

func (s Service) Register(req dto.RegisterRequest) (dto.RegisterResponse, error) {

	// TODO: use bcrypt
	//passwd := []byte(req.Password)
	//bcrypt.GenerateFromPassword(passwd, 0)

	// create new user in storage
	user := entity.User{
		ID:          0,
		Name:        req.Name,
		PhoneNumber: req.PhoneNumber,
		Password:    GetMD5Hash(req.Password),
	}
	createdUser, err := s.repo.Register(user)
	if err != nil {
		return dto.RegisterResponse{}, fmt.Errorf("unexpected error: %w", err)
	}

	return dto.RegisterResponse{dto.UserInfo{
		ID:          createdUser.ID,
		PhoneNumber: createdUser.PhoneNumber,
		Name:        createdUser.Name,
	}}, nil
}
