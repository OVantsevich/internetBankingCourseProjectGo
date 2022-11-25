package server

import (
	"context"
	"github.com/OVantsevich/internetBankingCourseProjectGo/accountService/services"
	"github.com/go-kit/kit/endpoint"
)

type Endpoints struct {
	CreateAccountEndpoint   endpoint.Endpoint
	GetUserAccountsEndpoint endpoint.Endpoint
}

func MakeServerEndpoints() Endpoints {
	return Endpoints{
		CreateAccountEndpoint:   MakeCreateAccountEndpoint(),
		GetUserAccountsEndpoint: MakeGetUserAccountsEndpoint(),
	}
}

func MakeCreateAccountEndpoint() endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(services.CreateAccountRequest)
		response, _ := services.CreateAccount(ctx, &req)
		return response, nil
	}
}

func MakeGetUserAccountsEndpoint() endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(services.GetUserAccountsRequest)
		accounts, response, err := services.GetUserAccounts(ctx, &req)
		if accounts != nil {
			return accounts, nil
		}
		return response, err
	}
}
