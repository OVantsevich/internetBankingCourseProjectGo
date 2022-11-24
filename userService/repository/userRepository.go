package repository

import (
	"context"
	"github.com/OVantsevich/internetBankingCourseProjectGo/userService/domain"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/labstack/gommon/log"
	"time"
)

var pool *pgxpool.Pool = nil
var Config domain.Config

func Pool(ctx context.Context) error {

	//if err := env.Parse(&Config); err != nil {
	//	log.Fatalf("something went wrong with environment, %e", err)
	//	return err
	//}
	Config.JwtKey = "874967EC3EA3490F8F2EF6478B72A756"
	Config.DatabaseUrl = "postgres://postgres:postgres@host.docker.internal:5432/userService?sslmode=disable"

	if pool == nil {
		var err error
		pool, err = pgxpool.Connect(ctx, Config.DatabaseUrl)
		if err != nil {
			log.Errorf("database connection error: %v", err)
			return err
		}
	}
	return nil
}

func Close() {
	if pool != nil {
		pool.Close()
	}
}

func CreateUser(ctx context.Context, user *domain.User) (string, error) {

	if err := Pool(ctx); err != nil {
		return "something went wrong", err
	}

	var userName, surname string
	if err := pool.QueryRow(ctx,
		"INSERT INTO users (user_login, user_email, user_password, user_name, surname) SELECT $1, $2, $3, $4, $5 WHERE NOT EXISTS(SELECT 1 FROM users WHERE user_login=$6) RETURNING user_name, surname",
		user.UserLogin, user.UserEmail, user.UserPassword, user.UserName, user.Surname, user.UserLogin).Scan(&userName, &surname); err != nil {
		log.Errorf("database error with create user: %v", err)
		return "user with this login already exist", err
	}

	return "Greetings, " + userName + " " + surname, nil
}

func GetUserByLogin(ctx context.Context, login string) (*domain.User, string, error) {

	if err := Pool(ctx); err != nil {
		return nil, "something went wrong", err
	}

	var user domain.User
	if err := pool.QueryRow(ctx, "select * from users where user_login=$1", login).Scan(
		&user.ID, &user.UserLogin, &user.UserEmail, &user.UserPassword, &user.UserName,
		&user.Surname, &user.IsDeleted, &user.CreationDate, &user.ModificationDate); err != nil {
		log.Errorf("database error with login user: %v", err)
		return nil, "user with this login doesn't exist", err
	}

	return &user, "", nil
}

func UpdateUser(ctx context.Context, user *domain.User) (string, error) {

	if err := Pool(ctx); err != nil {
		return "something went wrong", err
	}

	_, err := pool.Exec(ctx, "update users set user_name=$1, surname=$2, modification_date=$3, user_password=$4, user_email=$5 where user_login=$6",
		user.UserName, user.Surname, user.ModificationDate, user.UserPassword, user.UserEmail, user.UserLogin)
	if err != nil {
		log.Errorf("database error with update user: %v", err)
		return "user with this login doesn't exist", err
	}

	return user.UserLogin + " your account has been changed", nil
}

func DeleteUser(ctx context.Context, userLogin string, modificationDate time.Time) (string, error) {

	if err := Pool(ctx); err != nil {
		return "something went wrong", err
	}

	_, err := pool.Exec(ctx, "update users set is_deleted=$1, modification_date=$2 where user_login=$3",
		true, modificationDate, userLogin)
	if err != nil {
		log.Errorf("database error with delete user: %v", err)
		return "user with this login doesn't exist", err
	}

	return userLogin + " your account is deleted", nil
}
