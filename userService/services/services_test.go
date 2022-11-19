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
	databaseUrl = "postgres://postgres:postgres@192.168.100.5:5432/userService"
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
			UserLogin:    `Login1`,
			UserPassword: `oleg23102002`,
		},
		{
			UserName:     `Oleg`,
			Surname:      `Vantsevich`,
			UserLogin:    `Login2`,
			UserPassword: `password123`,
		},
	}
	testNoValidData := []domain.User{
		{
			UserName:     `Oleg`,
			Surname:      `Vantsevich`,
			UserLogin:    `Login1`,
			UserPassword: `oleg23102002`,
		},
		{
			UserName:     `Oleg`,
			Surname:      `Vantsevich`,
			UserLogin:    `Log`,
			UserPassword: `oleg23102002`,
		},
		{
			UserName:     `Oleg`,
			Surname:      `Vantsevich`,
			UserLogin:    `Login+1234567890`,
			UserPassword: `oleg23102002`,
		},
		{
			UserName:     `Oleg`,
			Surname:      `Vantsevich`,
			UserLogin:    `Login3`,
			UserPassword: `oleg`,
		},
		{
			UserName:     `Oleg`,
			Surname:      `Vantsevich`,
			UserLogin:    `Login3`,
			UserPassword: `1234567890`,
		},
		{
			UserName:     `Oleg`,
			Surname:      `Vantsevich`,
			UserLogin:    `Login3`,
			UserPassword: `oleg2310`,
		},
	}

	var ctx = context.Background()
	rps := NewService(repository.UserRepository{})
	for _, u := range testValidData {

		_, err := rps.rps.Pool(ctx, "").Exec(context.Background(), "delete from users where user_login=$1 ", u.UserLogin)
		require.NoError(t, err, "create error")

		_, err = rps.CreateUser(ctx, &u)
		require.NoError(t, err, "create error")

		_, err = rps.SignIn(ctx, &u)
		require.NoError(t, err, "create error")

	}
	for _, u := range testNoValidData {

		_, err := rps.CreateUser(ctx, &u)
		require.Error(t, err, "create error")

		_, err = rps.rps.Pool(ctx, "").Exec(ctx, "delete from users where user_login=$1 ", u.UserLogin)
		require.NoError(t, err, "create error")
	}
}

func TestSignIn(t *testing.T) {
	testValidData := []domain.User{
		{
			UserLogin:    `Login1`,
			UserPassword: `oleg23102002`,
		},
		{
			UserLogin:    `Login2`,
			UserPassword: `password123`,
		},
	}
	testNoExistentData := []domain.User{
		{
			UserLogin:    `Login3`,
			UserPassword: `oleg23102002`,
		},
		{
			UserLogin:    `Login4`,
			UserPassword: `password123`,
		},
	}
	testMismatchedPasswords := []domain.User{
		{
			UserLogin:    `Login1`,
			UserPassword: `password123`,
		},
		{
			UserLogin:    `Login2`,
			UserPassword: `oleg23102002`,
		},
	}

	var ctx = context.Background()
	rps := NewService(repository.UserRepository{})
	for _, u := range testValidData {
		_, err := rps.rps.Pool(ctx, "").Exec(context.Background(), "delete from users where user_login=$1 ", u.UserLogin)
		require.NoError(t, err, "create error")

		_, err = rps.CreateUser(ctx, &u)
		require.NoError(t, err, "create error")

		_, err = rps.SignIn(ctx, &u)
		require.NoError(t, err, "create error")
	}
	for _, u := range testNoExistentData {
		_, err := rps.SignIn(ctx, &u)
		require.Error(t, err, "create error")
	}
	for _, u := range testMismatchedPasswords {

		_, err := rps.SignIn(ctx, &u)
		require.Error(t, err, "create error")

		_, err = rps.rps.Pool(ctx, "").Exec(ctx, "delete from users where user_login=$1 ", u.UserLogin)
		require.NoError(t, err, "create error")
	}
}
