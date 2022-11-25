package server

import (
	"context"
	"encoding/json"
	"github.com/OVantsevich/internetBankingCourseProjectGo/accountService/services"
	"net/http"
)

func DecodeCreateAccountRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request services.CreateAccountRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

func DecodeGetUserAccountsRequest(_ context.Context, r *http.Request) (interface{}, error) {
	token := r.Header.Get("token")
	return services.GetUserAccountsRequest{Token: token}, nil
}
