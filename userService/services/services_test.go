package services

import (
	"context"
	"github.com/OVantsevich/internetBankingCourseProjectGo/userService/domain"
	"github.com/OVantsevich/internetBankingCourseProjectGo/userService/repository"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/labstack/gommon/log"
	"github.com/stretchr/testify/require"
	"testing"
)

var (
	databaseUrl = "postgres://postgres:postgres@localhost:5432/courseProject"
)

func Pool() *pgxpool.Pool {
	dbPool, err := pgxpool.Connect(context.Background(), databaseUrl)
	if err != nil {
		log.Errorf("database connection error: %v", err)
	}
	return dbPool
}

func TestCreate(t *testing.T) {
	testValidData := []domain.User{
		{
			UserName:     `Oleg`,
			Surname:      `Vantsevich`,
			UserLogin:    `Login1234`,
			UserPassword: `oleg23102002`,
		},
		{
			UserName:     `Oleg`,
			Surname:      `Vantsevich`,
			UserLogin:    `Login12345`,
			UserPassword: `password123`,
		},
	}
	testNoValidData := []domain.User{
		{
			UserName:     `Oleg`,
			Surname:      `Vantsevich`,
			UserLogin:    `Login12345`,
			UserPassword: `oleg23102002`,
		},
		{
			UserName:     `Oleg`,
			Surname:      `Vantsevich`,
			UserLogin:    `log`,
			UserPassword: `oleg23102002`,
		},
		{
			UserName:     `Oleg`,
			Surname:      `Vantsevich`,
			UserLogin:    `login1234567890login`,
			UserPassword: `oleg23102002`,
		},
		{
			UserName:     `Oleg`,
			Surname:      `Vantsevich`,
			UserLogin:    `Login123456`,
			UserPassword: `oleg`,
		},
		{
			UserName:     `Oleg`,
			Surname:      `Vantsevich`,
			UserLogin:    `Login1234567`,
			UserPassword: `1234567890`,
		},
		{
			UserName:     `Oleg`,
			Surname:      `Vantsevich`,
			UserLogin:    `Login12345678`,
			UserPassword: `oleg2310`,
		},
	}
	rps := NewService(repository.UserRepository{Pool: Pool()})
	for _, p := range testValidData {
		_, err := rps.CreateUser(&p)
		require.NoError(t, err, "create error")
	}
	for _, p := range testNoValidData {
		_, err := rps.CreateUser(&p)
		require.Error(t, err, "create error")
	}
}

func TestSignIn(t *testing.T) {
	testValidData := []domain.User{
		{
			UserLogin:    `oleg3`,
			UserPassword: `1`,
		},
		{
			UserLogin:    `oleg4`,
			UserPassword: `1`,
		},
	}
	testNoValidData := []domain.User{
		{
			UserLogin:    `olegggg`,
			UserPassword: `1`,
		},
		{
			UserLogin:    `oleg4`,
			UserPassword: `11241`,
		},
	}
	rps := NewService(repository.UserRepository{Pool: Pool()})
	for _, p := range testValidData {
		_, err := rps.SignIn(p.UserLogin, p.UserPassword)
		require.NoError(t, err, "create error")
	}
	for _, p := range testNoValidData {
		_, err := rps.SignIn(p.UserLogin, p.UserPassword)
		require.Error(t, err, "create error")
	}
}
