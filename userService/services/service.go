package services

import (
	"fmt"
	"github.com/OVantsevich/internetBankingCourseProjectGo/userService/domain"
	"github.com/OVantsevich/internetBankingCourseProjectGo/userService/repository"
	passwordvalidator "github.com/wagslane/go-password-validator"
)

type Service struct {
	rps repository.UserRepository
}

func NewService(pool repository.UserRepository) *Service {
	return &Service{rps: pool}
}

func (se *Service) SignIn(login, password string) (string, error) {
	return se.rps.SignIn(login, password)
}

func (se *Service) CreateUser(user *domain.User) (string, error) {

	err := passwordvalidator.Validate(user.UserPassword, 50)
	if err != nil {
		return err.Error(), err
	}

	if len(user.UserLogin) < 4 || len(user.UserLogin) > 15 {
		return "login should from 4 to 15 symbols", fmt.Errorf("database error with create user: login should from 4 to 15 symbols")
	}

	return se.rps.CreateUser(user)
}
