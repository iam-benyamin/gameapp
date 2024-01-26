package userservice

import (
	"gameapp/dto"
	"gameapp/pkg/richerror"
)

// Profile all request should be sanitized
func (s Service) Profile(req dto.ProfileRequest) (dto.ProfileResponse, error) {
	const op = "userservice.Profile"

	user, err := s.repo.GetUserByID(req.UserID)
	if err != nil {
		return dto.ProfileResponse{}, richerror.New(op).WithErr(err).WithMeta(map[string]interface{}{"req": req})
	}

	return dto.ProfileResponse{Name: user.Name}, nil
}
