package server

import (
	"context"
	"encoding/json"
	"github.com/OVantsevich/internetBankingCourseProjectGo/userService/domain"
	"net/http"
)

func DecodeUserRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request domain.User
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}
