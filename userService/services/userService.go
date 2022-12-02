package services

import (
	"context"
	"fmt"
	"github.com/OVantsevich/internetBankingCourseProjectGo/userService/domain"
	"github.com/OVantsevich/internetBankingCourseProjectGo/userService/eventStreaming"
	"github.com/OVantsevich/internetBankingCourseProjectGo/userService/repository"
	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/gommon/log"
	passwordvalidator "github.com/wagslane/go-password-validator"
	"golang.org/x/crypto/bcrypt"
	"net/mail"
	"strings"
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

	err := eventStreaming.JetStreamInit()
	if err != nil {
		return "something went wrong", err
	}

	if user == nil {
		log.Errorf("creating user, services: %v", fmt.Errorf("user is nil"))
		return "something went wrong", fmt.Errorf("user is nil")
	}

	if str, err := ValidLogin(user.UserLogin); err != nil {
		return str, err
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

	if user.UserPassword, err = HashingPassword(user.UserPassword); err != nil {
		return user.UserPassword, err
	}

	str, err := repository.CreateUser(ctx, user)
	if err == nil {
		if errLocal := eventStreaming.CreatingUser(user); errLocal != nil {
			log.Errorf("creating user, services, event streaming down: %v", errLocal)
			return str, err
		}
	}

	return str, err
}

func SignIn(ctx context.Context, request *SignInRequest) (string, error) {

	if request == nil {
		log.Errorf("signing user, services: %v", fmt.Errorf("request is nil"))
		return "something went wrong", fmt.Errorf("request is nil")
	}

	if _, err := ValidLogin(request.UserLogin); err != nil {
		return "Invalid login", err
	}

	user, str, err := repository.GetUserByLogin(ctx, request.UserLogin)
	if err != nil {
		return str, err
	}
	if !CheckPasswordHash(user.UserPassword, request.UserPassword) {
		return "Incorrect password", fmt.Errorf("incorrect password")
	}

	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["exp"] = time.Now().Add(time.Hour * 2).Unix()
	claims["login"] = request.UserLogin
	tokenString, err := token.SignedString([]byte(domain.Config.JwtKey))
	if err != nil {
		log.Errorf("signing user, services: %v", err)
		return "something went wrong", err
	}

	return tokenString, nil
}

func UpdateUser(ctx context.Context, request *UpdateUserRequest) (string, error) {

	if request == nil {
		log.Errorf("updating user, services: %v", fmt.Errorf("request is nil"))
		return "something went wrong", fmt.Errorf("request is nil")
	}

	claims := jwt.MapClaims{}
	str, err := ParseToken(request.Token, &claims)
	if err != nil {
		return str, nil
	}

	if str, err := ValidName(request.UserName, "name"); request.UserName != "" && err != nil {
		return str, err
	}
	if str, err := ValidName(request.Surname, "surname"); request.Surname != "" && err != nil {
		return str, err
	}
	if str, err := ValidPassword(request.UserPassword); request.UserPassword != "" && err != nil {
		return str, err
	}
	if str, err := ValidEmail(request.UserEmail); request.UserEmail != "" && err != nil {
		return str, err
	}

	if request.UserPassword, err = HashingPassword(request.UserPassword); err != nil {
		return request.UserPassword, err
	}

	user, str, err := repository.GetUserByLogin(ctx, claims["login"].(string))
	if err != nil {
		return str, err
	}
	if user.IsDeleted {
		return "user with this login doesn't exist", fmt.Errorf("user with this login doesn't exist")
	}

	user.UpdateUser(&domain.User{UserName: request.UserName, UserEmail: request.UserEmail,
		Surname: request.Surname, UserPassword: request.UserPassword})

	user.ModificationDate = time.Now()

	str, err = repository.UpdateUser(ctx, user)
	if err == nil {
		if errLocal := eventStreaming.CreatingUser(user); errLocal != nil {
			log.Errorf("creating user, services, event streaming down: %v", errLocal)
			return str, err
		}
	}

	return str, err
}

func DeleteUser(ctx context.Context, request *DeleteUserRequest) (string, error) {

	claims := jwt.MapClaims{}
	str, err := ParseToken(request.Token, &claims)
	if err != nil {
		return str, nil
	}
	return repository.DeleteUser(ctx, claims["login"].(string))
}

func ValidName(name string, fieldName string) (string, error) {
	if name == "" {
		return fieldName + " can't be empty", fmt.Errorf("validation error: " + fieldName + " can't be empty")
	}
	for _, r := range name {
		if !unicode.IsLetter(r) {
			return fieldName + " isn't valid", fmt.Errorf("validation error: " + fieldName + " isn't valid")
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

func CheckPasswordHash(hash, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func ParseToken(val string, claims *jwt.MapClaims) (string, error) {
	authHeaderParts := strings.Split(val, " ")
	if len(authHeaderParts) != 2 || !strings.EqualFold(authHeaderParts[0], "bearer") {
		return "not bearer auth", fmt.Errorf("not bearer auth")
	}

	_, err := jwt.ParseWithClaims(authHeaderParts[1], *claims, Key)
	if err != nil {
		return "invalid token", err
	}

	return "", nil
}

func Key(_ *jwt.Token) (interface{}, error) {
	return []byte(domain.Config.JwtKey), nil
}
