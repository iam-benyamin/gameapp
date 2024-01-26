package userservice

import (
	"fmt"
	"gameapp/dto"
	"gameapp/pkg/richerror"
)

func (s Service) Login(req dto.LoginRequest) (dto.LoginResponse, error) {
	// TODO: it would be better to user two separated method for existence check and GetUserByPhoneNumber
	const op = "userservice.Login"

	user, exist, err := s.repo.GetUserByPhoneNumber(req.PhoneNumber)
	if err != nil {
		return dto.LoginResponse{}, richerror.New(op).WithErr(err).
			WithMeta(map[string]interface{}{"phone_number": req.PhoneNumber})
	}

	if !exist {
		return dto.LoginResponse{}, fmt.Errorf("username or password is not correct")
	}

	if user.Password != GetMD5Hash(req.Password) {
		return dto.LoginResponse{}, fmt.Errorf("username or password is not correct")
	}

	accessToken, err := s.auth.CreateAccessToken(user)
	if err != nil {
		return dto.LoginResponse{}, fmt.Errorf("unexpected error: %w", err)
	}

	refreshToken, err := s.auth.CreateRefreshToken(user)
	if err != nil {
		return dto.LoginResponse{}, fmt.Errorf("unexpected error : %w", err)
	}

	return dto.LoginResponse{
		Token: dto.Tokens{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		},
		User: dto.UserInfo{
			ID:          user.ID,
			PhoneNumber: user.PhoneNumber,
			Name:        user.Name,
		},
	}, nil
}
