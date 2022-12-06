package services

import (
	"context"
	"fmt"
	"github.com/OVantsevich/internetBankingCourseProjectGo/accountService/domain"
	"github.com/OVantsevich/internetBankingCourseProjectGo/accountService/repository"
	"github.com/golang-jwt/jwt/v4"
	"strings"
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
	AccountName   string    `json:"account_name" sql:"type:varchar(40);default 'account'::character varying;not null"`
	Amount        int       `json:"amount" sql:"type:integer;default 0;not null"`
	AccountNumber string    `json:"account_number" sql:"type:bigint;not null"`
	CreationDate  time.Time `json:"creation_date" sql:"type:timestamp(6);default CURRENT_TIMESTAMP(6);not null"`
}

func CreateAccount(ctx context.Context, request *CreateAccountRequest) (string, error) {
	if str, err := ValidAccountName(request.AccountName, ""); err != nil {
		return str, err
	}

	claims := jwt.MapClaims{}
	str, err := ParseToken(request.Token, &claims)
	if err != nil {
		return str, err
	}

	return repository.CreateAccount(ctx, &domain.Account{AccountName: request.AccountName}, claims["login"].(string))
}

func GetUserAccounts(ctx context.Context, request *GetUserAccountsRequest) ([]GetUserAccountsResponse, string, error) {
	claims := jwt.MapClaims{}
	str, err := ParseToken(request.Token, &claims)
	if err != nil {
		return nil, str, err
	}

	accounts, str, err := repository.GetUserAccounts(ctx, claims["login"].(string))
	var response []GetUserAccountsResponse
	if accounts != nil {
		for _, account := range accounts {
			response = append(response, GetUserAccountsResponse{
				AccountName:   account.AccountName,
				Amount:        account.Amount,
				AccountNumber: account.AccountNumber,
				CreationDate:  account.CreationDate,
			})
		}
	}

	return response, str, err
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
