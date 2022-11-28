package server

import (
	"context"
	"encoding/json"
	"github.com/OVantsevich/internetBankingCourseProjectGo/accountService/domain"
	"github.com/OVantsevich/internetBankingCourseProjectGo/accountService/services"
	"net/http"
)

func EncodeCreateAccountResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}

func EncodeGetUserAccountsResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	encoder := json.NewEncoder(w)

	accounts := response.([]domain.Account)

	var err error
	for _, r := range accounts {
		err = encoder.Encode(services.GetUserAccountsResponse{AccountName: r.AccountName, Amount: r.Amount, CreationDate: r.CreationDate})
		if err != nil {
			return err
		}
	}

	return err
}

func EncodeCreateTransactionResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}

func EncodeGetAccountTransactionsResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	encoder := json.NewEncoder(w)

	accounts := response.([]services.GetAccountTransactionsResponse)

	var err error
	for _, r := range accounts {
		err = encoder.Encode(r)
		if err != nil {
			return err
		}
	}

	return err
}
