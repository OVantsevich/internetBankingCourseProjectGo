package repository

import (
	"context"
	"github.com/OVantsevich/internetBankingCourseProjectGo/accountService/domain"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/labstack/gommon/log"
)

var pool *pgxpool.Pool = nil

func Pool(ctx context.Context) (string, error) {
	if str, err := domain.InitConfig(); err != nil {
		return str, err
	}

	if pool == nil {
		var err error
		pool, err = pgxpool.Connect(ctx, domain.Config.DatabaseUrl)
		if err != nil {
			log.Errorf("database connection error: %v", err)
			return "database connection error", err
		}
	}
	return "", nil
}

func Close() {
	if pool != nil {
		pool.Close()
	}
}

func CreateUser(ctx context.Context, user *domain.User) error {

	if _, err := Pool(ctx); err != nil {
		return err
	}

	var userName, surname string
	if err := pool.QueryRow(ctx,
		"INSERT INTO users (user_login, user_name, surname) SELECT $1, $2, $3 WHERE NOT EXISTS(SELECT 1 FROM users WHERE user_login=$4) RETURNING user_name, surname",
		user.UserLogin, user.UserName, user.Surname, user.UserLogin).Scan(&userName, &surname); err != nil {
		log.Errorf("user with this login already exist: %v", err)
		return err
	}

	return nil
}

func GetUserByLogin(ctx context.Context, login string) (*domain.User, string, error) {

	if str, err := Pool(ctx); err != nil {
		return nil, str, err
	}

	var user domain.User
	row := pool.QueryRow(ctx, "select * from users where user_login=$1 and is_deleted=false", login)
	if err := row.Scan(
		&user.ID, &user.UserLogin, &user.UserName,
		&user.Surname, &user.IsDeleted); err != nil {
		log.Errorf("database error with login user: %v", err)
		return nil, "user with this login doesn't exist", err
	}

	return &user, "", nil
}
