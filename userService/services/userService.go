package services

import (
	"context"
	"fmt"
	"github.com/OVantsevich/internetBankingCourseProjectGo/userService/domain"
	"github.com/OVantsevich/internetBankingCourseProjectGo/userService/eventStreaming"
	"github.com/OVantsevich/internetBankingCourseProjectGo/userService/repository"
	"github.com/golang-jwt/jwt/v4"
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
		if err := eventStreaming.CreatingUser(user); err != nil {
			return "something went wrong", err
		}

	}

	return str, err
}

func SignIn(ctx context.Context, signIn *SignInRequest) (string, error) {

	if signIn == nil {
		return "", fmt.Errorf("user not found")
	}

	if _, err := ValidLogin(signIn.UserLogin); err != nil {
		return "Invalid login", err
	}

	user, str, err := repository.GetUserByLogin(ctx, signIn.UserLogin)
	if err != nil {
		return str, err
	}
	if user.IsDeleted {
		return "user with this login doesn't exist", fmt.Errorf("user with this login doesn't exist")
	}
	if !CheckPasswordHash(user.UserPassword, signIn.UserPassword) {
		return "Incorrect password", fmt.Errorf("incorrect password")
	}

	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["exp"] = time.Now().Add(time.Hour * 2).Unix()
	claims["login"] = signIn.UserLogin
	tokenString, err := token.SignedString([]byte(repository.Config.JwtKey))
	if err != nil {
		return "something went wrong", err
	}

	return tokenString, nil
}

func UpdateUser(ctx context.Context, userUpdateRequest *UpdateUserRequest) (string, error) {

	claims := jwt.MapClaims{}
	str, err := ParseToken(userUpdateRequest.Token, &claims)
	if err != nil {
		return str, nil
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

	user, str, err := repository.GetUserByLogin(ctx, claims["login"].(string))
	if err != nil {
		return str, err
	}
	if user.IsDeleted {
		return "user with this login doesn't exist", fmt.Errorf("user with this login doesn't exist")
	}

	user.CompareAndSet(&domain.User{UserName: userUpdateRequest.UserName, UserEmail: userUpdateRequest.UserEmail,
		Surname: userUpdateRequest.Surname, UserPassword: userUpdateRequest.UserPassword})

	user.ModificationDate = time.Now()

	return repository.UpdateUser(ctx, user)
}

func DeleteUser(ctx context.Context, userDeleteRequest *DeleteUserRequest) (string, error) {

	claims := jwt.MapClaims{}
	str, err := ParseToken(userDeleteRequest.Token, &claims)
	if err != nil {
		return str, nil
	}
	return repository.DeleteUser(ctx, claims["login"].(string), time.Now())
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
	return []byte(repository.Config.JwtKey), nil
}
