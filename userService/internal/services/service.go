package services

import (
	"github.com/OVantsevich/internetBankingCourseProjectGo/userService/internal/domain"
	"github.com/OVantsevich/internetBankingCourseProjectGo/userService/internal/repository"
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
	return se.rps.CreateUser(user)
}
