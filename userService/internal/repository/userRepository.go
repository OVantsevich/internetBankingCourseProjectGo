package repository

import (
	"context"
	"github.com/OVantsevich/internetBankingCourseProjectGo/userService/internal/domain"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/labstack/gommon/log"
	passwordvalidator "github.com/wagslane/go-password-validator"
)

type UserRepository struct {
	Pool *pgxpool.Pool
}

func (repos *UserRepository) CreateUser(user *domain.User) (string, error) {

	rows, err := repos.Pool.Query(context.Background(), "select user_login from users where user_login=$1", user.UserLogin)
	if err != nil {
		log.Errorf("database error with create user: %v", err)
		return "", err
	}
	if rows.Next() {
		return "A user with this login already exists", nil
	}

	err = passwordvalidator.Validate(user.UserPassword, 50)
	if err != nil {
		return err.Error(), nil
	}

	_, err = repos.Pool.Exec(context.Background(), "insert into users(user_name, surname, user_login, user_password) values($1, $2, $3, $4)",
		user.UserName, user.Surname, user.UserLogin, user.UserPassword)
	if err != nil {
		log.Errorf("database error with create user: %v", err)
		return "", err
	}

	return "Greetings: " + user.UserName + " " + user.Surname, nil
}

func (repos *UserRepository) SignIn(userLogin, userPassword string) (string, error) {

	var password string
	rows, err := repos.Pool.Query(context.Background(), "select user_password from users  where user_login=$1", userLogin)
	if err != nil {
		log.Errorf("database error with create user: %v", err)
		return "", err
	}
	if !rows.Next() {
		return "User with this login does not exist", nil
	}
	err = rows.Scan(&password)
	if err != nil {
		log.Errorf("database error with create user: %v", err)
		return "", err
	}
	if password != userPassword {
		return "Wrong password" + userPassword, nil
	}
	return "Welcome", nil
}
