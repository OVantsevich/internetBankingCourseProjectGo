package services

import (
	"context"
	"fmt"
	"github.com/OVantsevich/internetBankingCourseProjectGo/userService/domain"
	"github.com/OVantsevich/internetBankingCourseProjectGo/userService/repository"
	passwordvalidator "github.com/wagslane/go-password-validator"
)

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
