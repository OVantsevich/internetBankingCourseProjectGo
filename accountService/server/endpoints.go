package server

import (
	"context"
	"github.com/OVantsevich/internetBankingCourseProjectGo/accountService/services"
	"github.com/go-kit/kit/endpoint"
)

type Endpoints struct {
	CreateAccountEndpoint          endpoint.Endpoint
	GetUserAccountsEndpoint        endpoint.Endpoint
	CreateTransactionEndpoint      endpoint.Endpoint
	GetAccountTransactionsEndpoint endpoint.Endpoint
}

func MakeServerEndpoints() Endpoints {
	return Endpoints{
		CreateAccountEndpoint:          MakeCreateAccountEndpoint(),
		GetUserAccountsEndpoint:        MakeGetUserAccountsEndpoint(),
		CreateTransactionEndpoint:      MakeCreateTransactionEndpoint(),
		GetAccountTransactionsEndpoint: MakeGetAccountTransactionsEndpoint(),
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

func MakeCreateTransactionEndpoint() endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(services.CreateTransactionRequest)
		response, _ := services.CreateTransaction(ctx, &req)
		return response, nil
	}
}

func MakeGetAccountTransactionsEndpoint() endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(services.GetAccountTransactionsRequest)
		transactions, response, err := services.GetAccountTransactions(ctx, &req)
		if transactions != nil {
			return transactions, nil
		}
		return response, err
	}
}
