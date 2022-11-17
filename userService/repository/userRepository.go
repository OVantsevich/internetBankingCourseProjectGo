package repository

import (
	"context"
	"fmt"
	"github.com/OVantsevich/internetBankingCourseProjectGo/userService/domain"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/labstack/gommon/log"
	"time"
)

type UserRepository struct {
	Pool *pgxpool.Pool
}

func (repos *UserRepository) CreateUser(user *domain.User) (string, error) {

	rows, err := repos.Pool.Query(context.Background(), "select user_login from users where user_login=$1", user.UserLogin)
	if err != nil {
		log.Errorf("database error with create user: %v", err)
		return "something went wrong", err
	}
	if rows.Next() {
		return "user with this login already exists", fmt.Errorf("database error with create user: " +
			"user with this login already exists")
	}

	_, err = repos.Pool.Exec(context.Background(), "insert into users(user_name, surname, user_login, user_password) "+
		"values($1, $2, $3, $4)", user.UserName, user.Surname, user.UserLogin, user.UserPassword)
	if err != nil {
		log.Errorf("database error with create user: %v", err)
		return "something went wrong", err
	}

	return "Greetings, " + user.UserLogin, nil
}

func (repos *UserRepository) SignIn(userLogin, userPassword string) (string, error) {

	var password string
	rows, err := repos.Pool.Query(context.Background(), "select user_password from users  where user_login=$1", userLogin)
	if err != nil {
		log.Errorf("database error with create user: %v", err)
		return "something went wrong", err
	}
	if !rows.Next() {
		return "user with this login does not exist", fmt.Errorf("database error with create user: " +
			"user with this login does not exist")
	}
	err = rows.Scan(&password)
	if err != nil {
		log.Errorf("database error with create user: %v", err)
		return "something went wrong", err
	}
	if password != userPassword {
		return "wrong password" + userPassword, fmt.Errorf("wrong password")
	}
	return "Welcome, " + userLogin, nil
}

func (repos *UserRepository) UpdateUser(user *domain.User) (string, error) {

	rows, err := repos.Pool.Query(context.Background(), "select user_name, surname, is_deleted, user_password "+
		"from users where user_login=$1", user.UserLogin)
	if err != nil {
		log.Errorf("database error with create user: %v", err)
		return "something went wrong", err
	}
	if !rows.Next() {
		return "user with this login does not exist", fmt.Errorf("database error with create user: " +
			"user with this login does not exist")
	}

	var currentUser domain.User
	err = rows.Scan(currentUser.UserName, currentUser.Surname, currentUser.IsDeleted, currentUser.UserPassword)

	currentUser.CompareAndSet(user)
	currentUser.ModificationDate = time.Now()

	_, err = repos.Pool.Exec(context.Background(), "update users set user_name=$1, surname=$2, is_deleted=$3, modification_date=$4, user_password=$5 where user_login=$6",
		currentUser.UserName, currentUser.Surname, currentUser.IsDeleted, currentUser.ModificationDate, currentUser.UserPassword, currentUser.UserLogin)
	if err != nil {
		log.Errorf("database error with create user: %v", err)
		return "something went wrong", err
	}

	return user.UserLogin + " your account has been changed", nil
}
