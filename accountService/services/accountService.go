package services

import (
	"context"
	"fmt"
	"github.com/OVantsevich/internetBankingCourseProjectGo/accountService/domain"
	"github.com/OVantsevich/internetBankingCourseProjectGo/accountService/repository"
	"github.com/dgrijalva/jwt-go"
	"time"
	"unicode"
)

type CreateAccountRequest struct {
	AccountName string `json:"account_name" sql:"type:varchar(40);default 'account'::character varying;not null"`
	Token       string `json:"token"`
}

type GetUserAccountsRequest struct {
	Token string `json:"token"`
}

type GetUserAccountsResponse struct {
	AccountName  string    `json:"account_name" sql:"type:varchar(40);default 'account'::character varying;not null"`
	Amount       int       `json:"amount" sql:"type:integer;default 0;not null"`
	CreationDate time.Time `json:"creation_date" sql:"type:timestamp(6);default CURRENT_TIMESTAMP(6);not null"`
}

func CreateAccount(ctx context.Context, request *CreateAccountRequest) (string, error) {
	if str, err := domain.InitConfig(); err != nil {
		return str, err
	}

	if str, err := ValidAccountName(request.AccountName, ""); err != nil {
		return str, err
	}

	claims := jwt.MapClaims{}
	_, err := jwt.ParseWithClaims(request.Token, claims,
		func(token *jwt.Token) (interface{}, error) {
			return []byte(domain.Config.JwtKey), nil
		})
	if err != nil {
		return "Invalid token", err
	}

	user, str, err := repository.GetUserByLogin(ctx, claims["login"].(string))
	if err != nil {
		return str, nil
	}

	return repository.CreateAccount(ctx, &domain.Account{UserId: user.ID, AccountName: request.AccountName})
}

func GetUserAccounts(ctx context.Context, request *GetUserAccountsRequest) ([]domain.Account, string, error) {
	if str, err := domain.InitConfig(); err != nil {
		return nil, str, err
	}

	claims := jwt.MapClaims{}
	_, err := jwt.ParseWithClaims(request.Token, claims,
		func(token *jwt.Token) (interface{}, error) {
			return []byte(domain.Config.JwtKey), nil
		})
	if err != nil {
		return nil, "Invalid token", err
	}

	user, str, err := repository.GetUserByLogin(ctx, claims["login"].(string))
	if err != nil {
		return nil, str, nil
	}

	return repository.GetUserAccounts(ctx, user.ID)
}

func ValidAccountName(name string, filedName string) (string, error) {
	if name == "" {
		return "account " + filedName + " name can't be empty", fmt.Errorf("validation error: account " + filedName + " name can't be empty")
	}
	for _, r := range name {
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) {
			return "account " + filedName + " name isn't valid", fmt.Errorf("validation error: account " + filedName + " name isn't valid")
		}
	}
	return "", nil
}
