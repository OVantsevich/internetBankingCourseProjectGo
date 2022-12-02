package server

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/OVantsevich/internetBankingCourseProjectGo/accountService/services"
	"github.com/labstack/gommon/log"
	"net/http"
)

func EncodeCreateAccountResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}

func EncodeGetUserAccountsResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	encoder := json.NewEncoder(w)

	accounts, ok := response.([]services.GetUserAccountsResponse)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		log.Errorf("error cast, getUserAccounts")
		return fmt.Errorf("something went wrong")
	}

	var err error
	for _, r := range accounts {
		err = encoder.Encode(r)
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

	transactions, ok := response.([]services.GetAccountTransactionsResponse)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		log.Errorf("error cast, getUserAccounts")
		return fmt.Errorf("something went wrong")
	}

	var err error
	for _, r := range transactions {
		err = encoder.Encode(r)
		if err != nil {
			return err
		}
	}

	return err
}
