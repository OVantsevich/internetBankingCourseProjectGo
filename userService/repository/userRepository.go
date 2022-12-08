package repository

import (
	"context"
	"github.com/OVantsevich/internetBankingCourseProjectGo/userService/domain"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/labstack/gommon/log"
	"time"
)

var pool *pgxpool.Pool = nil

func Pool(ctx context.Context) error {
	if pool == nil {
		var err error
		pool, err = pgxpool.Connect(ctx, domain.Config.DatabaseUrl)
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

func CreateUser(ctx context.Context, user *domain.User, key string) (string, error) {

	if err := Pool(ctx); err != nil {
		return "something went wrong", err
	}

	var userEmail string
	if err := pool.QueryRow(ctx,
		"INSERT INTO users (user_login, user_email, user_password, user_name, surname, is_verified, verification_date) SELECT $1, $2, $3, $4, $5, $6, $7 WHERE NOT EXISTS(SELECT 1 FROM users WHERE user_login=$8) RETURNING user_email",
		user.UserLogin, user.UserEmail, user.UserPassword, user.UserName, user.Surname, key, time.Now().Add(time.Minute*5), user.UserLogin).Scan(&userEmail); err != nil {
		log.Errorf("database error with create user: %v", err)
		return "user with this login already exist", err
	}

	return "Please verify your account using email: " + userEmail, nil
}

func Verify(ctx context.Context, key string) (*domain.User, string, error) {

	if err := Pool(ctx); err != nil {
		return nil, "something went wrong", err
	}

	var user = &domain.User{}

	if err := pool.QueryRow(ctx, "update users "+
		"set is_verified='' "+
		"where is_verified = $1 "+
		"and verification_date > $2 "+
		"returning user_login, user_email, user_name, surname",
		key, time.Now()).Scan(&user.UserLogin, &user.UserEmail, &user.UserName, &user.Surname); err != nil {
		log.Errorf("database error with create user: %v", err)
		return nil, "user with this login already verified, or verification time is expired", err
	}

	return user, "Greetings, " + user.UserName + " " + user.Surname, nil
}

func GetUserByLogin(ctx context.Context, login string) (*domain.User, string, error) {

	if err := Pool(ctx); err != nil {
		return nil, "something went wrong", err
	}

	var user domain.User
	if err := pool.QueryRow(ctx, "select * from users where user_login=$1 and not is_deleted and is_verified=''", login).Scan(
		&user.ID, &user.UserLogin, &user.UserEmail, &user.UserPassword, &user.UserName,
		&user.Surname, &user.IsVerified, &user.VerificationDate, &user.IsDeleted, &user.CreationDate, &user.ModificationDate); err != nil {
		log.Errorf("database error with login user: %v", err)
		return nil, "user with this login doesn't exist", err
	}

	return &user, "", nil
}

func UpdateUser(ctx context.Context, user *domain.User) (string, error) {

	if err := Pool(ctx); err != nil {
		return "something went wrong", err
	}

	var id int
	if err := pool.QueryRow(ctx, "update users set user_name=$1, surname=$2, modification_date=$3, user_password=$4 where user_login=$5 and is_deleted=false returning id",
		user.UserName, user.Surname, user.ModificationDate, user.UserPassword, user.UserLogin).Scan(&id); err != nil {
		log.Errorf("database error with update user: %v", err)
		return "user with this login doesn't exist", err
	}

	return user.UserLogin + " your account has been changed", nil
}

func DeleteUser(ctx context.Context, userLogin string) (string, error) {

	if err := Pool(ctx); err != nil {
		return "something went wrong", err
	}

	var id int
	if err := pool.QueryRow(ctx, "update users set is_deleted=true, modification_date=$1 where user_login=$2 and is_deleted=false returning id",
		time.Now(), userLogin).Scan(&id); err != nil {
		log.Errorf("database error with delete user: %v", err)
		return "user with this login doesn't exist", err
	}

	return userLogin + " your account is deleted", nil
}
