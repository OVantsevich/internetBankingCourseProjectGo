package server

import (
	"context"
	"encoding/json"
	"github.com/OVantsevich/internetBankingCourseProjectGo/accountService/services"
	"net/http"
)

func DecodeCreateAccountRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request services.CreateAccountRequest
	request.Token = r.Header.Get("Authorization")
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

func DecodeGetUserAccountsRequest(_ context.Context, r *http.Request) (interface{}, error) {
	token := r.Header.Get("Authorization")
	return services.GetUserAccountsRequest{Token: token}, nil
}

func DecodeCreateTransactionRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request services.CreateTransactionRequest
	request.Token = r.Header.Get("Authorization")
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

func DecodeGetAccountTransactionsRequest(_ context.Context, r *http.Request) (interface{}, error) {
	token := r.Header.Get("Authorization")
	number := r.Header.Get("Account-Number")
	return services.GetAccountTransactionsRequest{Token: token, AccountNumber: number}, nil
}
