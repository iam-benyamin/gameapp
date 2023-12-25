package userservice

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"gameapp/entity"
	"gameapp/pkg/phonenumber"
)

type Repository interface {
	IsPhoneNumberUnique(phoneNumber string) (bool, error)
	Register(u entity.User) (entity.User, error)
	GetUserByPhoneNumber(phoneNumber string) (entity.User, bool, error)
	GetUserByID(UserID uint) (entity.User, error)
}

type AuthGenerator interface {
	CreateAccessToken(user entity.User) (string, error)
	CreateRefreshToken(user entity.User) (string, error)
}

type Service struct {
	auth AuthGenerator
	repo Repository
}

func New(authGenerator AuthGenerator, repo Repository) Service {
	return Service{auth: authGenerator, repo: repo}
}

type RegisterRequest struct {
	Name        string `json:"name"`         // struct tag are like meta information and compiler will ignore them
	PhoneNumber string `json:"phone_number"` // but some package like json marshal will look at them
	Password    string `json:"password"`
}

type UserInfo struct {
	ID          uint   `json:"id"`
	PhoneNumber string `json:"phone_number"`
	Name        string `json:"name"`
}

type RegisterResponse struct {
	User UserInfo `json:"user"`
}

func (s Service) Register(req RegisterRequest) (RegisterResponse, error) {
	// TODO: we should verify phone number by verification code

	// validate phone number
	if !phonenumber.IsValid(req.PhoneNumber) {
		return RegisterResponse{}, fmt.Errorf("phone number is not valid")
	}

	// check uniqueness of phone number
	if isUnique, err := s.repo.IsPhoneNumberUnique(req.PhoneNumber); err != nil || !isUnique {
		if err != nil {
			return RegisterResponse{}, fmt.Errorf("unexpected error: %w", err)
		}
		if !isUnique {
			return RegisterResponse{}, fmt.Errorf("phone number is not unique")
		}
	}

	// validate name
	if len(req.Name) < 3 {
		return RegisterResponse{}, fmt.Errorf("name lenght should be grater than 3")
	}

	// TODO: check the password with regex pattern
	// validate password
	if len(req.Password) < 8 {
		return RegisterResponse{}, fmt.Errorf("password lenght sholud be grater than 8")
	}

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
		return RegisterResponse{}, fmt.Errorf("unexpected error: %w", err)
	}

	return RegisterResponse{UserInfo{
		ID:          createdUser.ID,
		PhoneNumber: createdUser.PhoneNumber,
		Name:        createdUser.Name,
	}}, nil
}

type LoginRequest struct {
	PhoneNumber string `json:"phone_number"`
	Password    string `json:"password"`
}

type Tokens struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
type LoginResponse struct {
	User  UserInfo `json:"user"`
	Token Tokens   `json:"token"`
}

func (s Service) Login(req LoginRequest) (LoginResponse, error) {
	// TODO: it would be better to user two separated method for existence check and GetUserByPhoneNumber

	user, exist, err := s.repo.GetUserByPhoneNumber(req.PhoneNumber)
	if err != nil {
		return LoginResponse{}, fmt.Errorf("unexpected error: %w", err)
	}

	if !exist {
		return LoginResponse{}, fmt.Errorf("username or password is not correct")
	}

	if user.Password != GetMD5Hash(req.Password) {
		return LoginResponse{}, fmt.Errorf("username or password is not correct")
	}

	accessToken, err := s.auth.CreateAccessToken(user)
	if err != nil {
		return LoginResponse{}, fmt.Errorf("unexpected error: %w", err)
	}

	refreshToken, err := s.auth.CreateRefreshToken(user)
	if err != nil {
		return LoginResponse{}, fmt.Errorf("unexpected error : %w", err)
	}

	return LoginResponse{
		Token: Tokens{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		},
		User: UserInfo{
			ID:          user.ID,
			PhoneNumber: user.PhoneNumber,
			Name:        user.Name,
		},
	}, nil
}

func GetMD5Hash(text string) string {
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}

type ProfileRequest struct {
	UserID uint
}

type ProfileResponse struct {
	Name string `json:"name"`
}

// Profile all request should be sanitized
func (s Service) Profile(req ProfileRequest) (ProfileResponse, error) {
	user, err := s.repo.GetUserByID(req.UserID)
	if err != nil {
		// I don't expect the repository call return "not found" error.
		// because I assume the interactor input is sanitized
		// TODO: we can use Rich Error
		return ProfileResponse{}, fmt.Errorf("unexpected error : %w", err)
	}

	return ProfileResponse{Name: user.Name}, nil
}
