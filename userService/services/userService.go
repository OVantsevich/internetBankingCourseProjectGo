package services

import (
	"context"
	"fmt"
	"github.com/OVantsevich/internetBankingCourseProjectGo/userService/domain"
	"github.com/OVantsevich/internetBankingCourseProjectGo/userService/repository"
	"github.com/dgrijalva/jwt-go"
	passwordvalidator "github.com/wagslane/go-password-validator"
	"golang.org/x/crypto/bcrypt"
	"net/mail"
	"time"
	"unicode"
)

type SignInRequest struct {
	UserLogin    string `json:"user_login" sql:"type:varchar(50);not null"`
	UserPassword string `json:"user_password" sql:"type:varchar(50);not null"`
}

type UpdateUserRequest struct {
	UserEmail    string `json:"user_email" sql:"type:varchar(50);not null"`
	UserPassword string `json:"user_password" sql:"type:varchar(50);not null"`
	UserName     string `json:"user_name" sql:"type:varchar(50);not null"`
	Surname      string `json:"surname" sql:"type:varchar(50);not null"`
	Token        string `json:"token"`
}

type DeleteUserRequest struct {
	Token string `json:"token"`
}

func CreateUser(ctx context.Context, user *domain.User) (string, error) {

	if str, err := ValidLogin(user.UserLogin); err != nil {
		return str, err
	}

	if user, str, err := repository.GetUserByLogin(ctx, user.UserLogin); err != nil {
		return str, err
	} else {
		if user != nil {
			return "user with this login already exists", fmt.Errorf("database error with create user: " +
				"user with this login already exists")
		}
	}

	if user, str, err := repository.GetUserByEmail(ctx, user.UserEmail); err != nil {
		return str, err
	} else {
		if user != nil {
			return "this email is already in use", fmt.Errorf("database error with create user: " +
				"this email is already in use")
		}
	}

	if str, err := ValidPassword(user.UserPassword); err != nil {
		return str, err
	}
	if str, err := ValidEmail(user.UserEmail); err != nil {
		return str, err
	}
	if str, err := ValidName(user.UserName, "name"); err != nil {
		return str, err
	}
	if str, err := ValidName(user.Surname, "surname"); err != nil {
		return str, err
	}

	var err error
	if user.UserPassword, err = HashingPassword(user.UserPassword); err != nil {
		return user.UserPassword, err
	}

	return repository.CreateUser(ctx, user)
}

func SignIn(ctx context.Context, signIn *SignInRequest) (string, error) {

	if signIn == nil {
		return "", fmt.Errorf("user not found")
	}

	if _, err := ValidLogin(signIn.UserLogin); err != nil {
		return "Invalid login", err
	}

	var err error
	var password string
	if password, err = HashingPassword(signIn.UserPassword); err != nil {
		return password, err
	}

	var str string
	var user *domain.User
	if user, str, err = repository.GetUserByLogin(ctx, signIn.UserLogin); err != nil {
		return str, err
	}
	if user.IsDeleted {
		return "user with this login doesn't exist", fmt.Errorf("user with this login doesn't exist")
	}
	if CheckPasswordHash(password, user.UserPassword) {
		return "Incorrect password", fmt.Errorf("incorrect password")
	}

	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["exp"] = time.Now().Add(20 * time.Minute)
	claims["authorized"] = true
	claims["user"] = signIn.UserLogin
	tokenString, err := token.SignedString([]byte(repository.Config.JwtKey))
	if err != nil {
		return "something went wrong", err
	}

	return tokenString, nil
}

func UpdateUser(ctx context.Context, userUpdateRequest *UpdateUserRequest) (string, error) {

	claims := jwt.MapClaims{}
	_, err := jwt.ParseWithClaims(userUpdateRequest.Token, claims,
		func(token *jwt.Token) (interface{}, error) {
			return []byte(repository.Config.JwtKey), nil
		})
	if err != nil {
		return "Invalid token", err
	}

	if str, err := ValidName(userUpdateRequest.UserName, "name"); userUpdateRequest.UserName != "" && err != nil {
		return str, err
	}
	if str, err := ValidName(userUpdateRequest.Surname, "surname"); userUpdateRequest.Surname != "" && err != nil {
		return str, err
	}
	if str, err := ValidPassword(userUpdateRequest.UserPassword); userUpdateRequest.UserPassword != "" && err != nil {
		return str, err
	}
	if str, err := ValidEmail(userUpdateRequest.UserEmail); userUpdateRequest.UserEmail != "" && err != nil {
		return str, err
	}

	if userUpdateRequest.UserPassword, err = HashingPassword(userUpdateRequest.UserPassword); err != nil {
		return userUpdateRequest.UserPassword, err
	}

	var str string
	var user *domain.User
	if user, str, err = repository.GetUserByLogin(ctx, claims["user"].(string)); err != nil {
		return str, err
	} else {
		if user == nil {
			return "user with this login doesn't exists", fmt.Errorf("database error with create user: " +
				"user with this login doesn't exists")
		}
	}
	if user.IsDeleted {
		return "user with this login doesn't exist", fmt.Errorf("user with this login doesn't exist")
	}

	if user, str, err := repository.GetUserByEmail(ctx, userUpdateRequest.UserEmail); err != nil {
		return str, err
	} else {
		if user != nil {
			return "this email is already in use", fmt.Errorf("database error with create user: " +
				"this email is already in use")
		}
	}

	user.CompareAndSet(&domain.User{UserName: userUpdateRequest.UserName, UserEmail: userUpdateRequest.UserEmail,
		Surname: userUpdateRequest.Surname, UserPassword: userUpdateRequest.UserPassword})

	user.ModificationDate = time.Now()

	return repository.UpdateUser(ctx, user)
}

func DeleteUser(ctx context.Context, userDeleteRequest *DeleteUserRequest) (string, error) {

	claims := jwt.MapClaims{}
	_, err := jwt.ParseWithClaims(userDeleteRequest.Token, claims,
		func(token *jwt.Token) (interface{}, error) {
			return []byte(repository.Config.JwtKey), nil
		})
	if err != nil {
		return "Invalid token", err
	}

	return repository.DeleteUser(ctx, claims["user"].(string), time.Now())
}

func ValidName(name string, fieldName string, isNil ...bool) (string, error) {
	if name == "" {
		return fieldName + "can't be empty", fmt.Errorf("validation error:" + fieldName + "can't be empty")
	}
	for _, r := range name {
		if !unicode.IsLetter(r) {
			return fieldName + "isn't valid", fmt.Errorf("validation error:" + fieldName + " isn't valid")
		}
	}
	return "", nil
}

func ValidLogin(login string) (string, error) {
	if len(login) < 4 || len(login) > 15 {
		return "login should be greater than 3 and less than 16 symbols", fmt.Errorf("validation error: login should be greater than 3 and less than 16 symbols")
	}
	return "", nil
}

func ValidPassword(password string) (string, error) {
	err := passwordvalidator.Validate(password, 50)
	if err != nil {
		return err.Error(), err
	}
	return "", nil
}

func ValidEmail(email string) (string, error) {
	_, err := mail.ParseAddress(email)
	if err != nil {
		return err.Error(), err
	}
	return "", nil
}

func HashingPassword(password string) (string, error) {
	hashedBytesPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "something went wrong", err
	}
	return string(hashedBytesPassword), nil
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
