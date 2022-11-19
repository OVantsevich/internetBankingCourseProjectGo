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
	pool *pgxpool.Pool
}

func (repos *UserRepository) Pool(ctx context.Context, url string) *pgxpool.Pool {
	databaseUrl := "postgres://postgres:postgres@host.docker.internal:5432/userService"
	if repos.pool == nil {
		var err error
		repos.pool, err = pgxpool.Connect(context.Background(), databaseUrl)
		if err != nil {
			log.Errorf("database connection error: %v", err)
		}
	}
	return repos.pool
}

func (repos *UserRepository) Close() {
	if repos.pool != nil {
		repos.pool.Close()
	}
}

func (repos *UserRepository) CreateUser(ctx context.Context, user *domain.User) (string, error) {

	rows, err := repos.Pool(ctx, "").Query(ctx, "select user_login from users where user_login=$1", user.UserLogin)
	if err != nil {
		log.Errorf("database error with create user: %v", err)
		return "something went wrong", err
	}
	if rows.Next() {
		return "user with this login already exists", fmt.Errorf("database error with create user: " +
			"user with this login already exists")
	}

	_, err = repos.Pool(ctx, "").Exec(ctx, "insert into users(user_name, surname, user_login, user_password) "+
		"values($1, $2, $3, $4)", user.UserName, user.Surname, user.UserLogin, user.UserPassword)
	if err != nil {
		log.Errorf("database error with create user: %v", err)
		return "something went wrong", err
	}

	return "Greetings, " + user.UserLogin, nil
}

func (repos *UserRepository) SignIn(ctx context.Context, user *domain.User) (string, error) {

	var password string
	rows, err := repos.Pool(ctx, "").Query(ctx, "select user_password from users  where user_login=$1", user.UserLogin)
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
	if password != user.UserPassword {
		return "wrong password", fmt.Errorf("wrong password")
	}
	return "Welcome, " + user.UserLogin, nil
}

func (repos *UserRepository) UpdateUser(ctx context.Context, user *domain.User) (string, error) {

	rows, err := repos.Pool(ctx, "").Query(context.Background(), "select user_name, surname, is_deleted, user_password "+
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

	_, err = repos.Pool(ctx, "").Exec(context.Background(), "update users set user_name=$1, surname=$2, is_deleted=$3, modification_date=$4, user_password=$5 where user_login=$6",
		currentUser.UserName, currentUser.Surname, currentUser.IsDeleted, currentUser.ModificationDate, currentUser.UserPassword, currentUser.UserLogin)
	if err != nil {
		log.Errorf("database error with create user: %v", err)
		return "something went wrong", err
	}

	return user.UserLogin + " your account has been changed", nil
}
