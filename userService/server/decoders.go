package server

import (
	"context"
	"encoding/json"
	"github.com/OVantsevich/internetBankingCourseProjectGo/userService/domain"
	"github.com/OVantsevich/internetBankingCourseProjectGo/userService/services"
	"net/http"
)

func DecodeCreateUserRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request domain.User
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

func DecodeVerificationRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request services.VerificationRequest
	request.Verification = r.URL.Query().Get("verification")
	return request, nil
}

func DecodeSignInRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request services.SignInRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

func DecodeUpdateUserRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request services.UpdateUserRequest
	request.Token = r.Header.Get("Authorization")
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

func DecodeDeleteUserRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request services.DeleteUserRequest
	request.Token = r.Header.Get("Authorization")
	return request, nil
}
