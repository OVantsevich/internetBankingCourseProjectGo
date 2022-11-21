package services

import (
	"context"
	"fmt"
	"github.com/OVantsevich/internetBankingCourseProjectGo/userService/domain"
	"github.com/OVantsevich/internetBankingCourseProjectGo/userService/repository"
	passwordvalidator "github.com/wagslane/go-password-validator"
)

type UpdateUserRequest struct {
	UserName     string `json:"user_name" sql:"type:varchar(20);not null"`
	Surname      string `json:"surname" sql:"type:varchar(50);not null"`
	UserLogin    string `json:"user_login" sql:"type:varchar(100);not null"`
	UserPassword string `json:"user_password" sql:"type:varchar(200);not null"`
	Password     string `json:"password"`
}

type Service struct {
	rps repository.UserRepository
}

func (s *Service) Close() {
	s.rps.Close()
}

func NewService(pool repository.UserRepository) *Service {
	return &Service{rps: pool}
}

func (s *Service) SignIn(ctx context.Context, user *domain.User) (string, error) {
	return s.rps.SignIn(ctx, user)
}

func (s *Service) CreateUser(ctx context.Context, user *domain.User) (string, error) {

	err := passwordvalidator.Validate(user.UserPassword, 50)
	if err != nil {
		return err.Error(), err
	}

	if len(user.UserLogin) < 4 || len(user.UserLogin) > 15 {
		return "login should be Greater than 3 and less than 16 symbols", fmt.Errorf("database error with create user: login should from 4 to 15 symbols")
	}

	return s.rps.CreateUser(ctx, user)
}

func (s *Service) UpdateUser(ctx context.Context, userUpdateRequest *UpdateUserRequest) (string, error) {

	if userUpdateRequest.UserPassword != "" {
		err := passwordvalidator.Validate(userUpdateRequest.UserPassword, 50)
		if err != nil {
			return err.Error(), fmt.Errorf("wrong password")
		}
	}

	var user = domain.User{}
	user.UserName = userUpdateRequest.UserName
	user.UserLogin = userUpdateRequest.UserLogin
	user.UserPassword = userUpdateRequest.UserPassword
	user.Surname = userUpdateRequest.Surname

	return s.rps.UpdateUser(ctx, &user, userUpdateRequest.Password)
}
