package server

import (
	"github.com/OVantsevich/internetBankingCourseProjectGo/accountService/services"
	"github.com/go-kit/kit/auth/jwt"
	httptransport "github.com/go-kit/kit/transport/http"
	basicjwt "github.com/golang-jwt/jwt/v4"
	"github.com/gorilla/mux"
	"net/http"
)

func MakeHTTPHandler() http.Handler {
	r := mux.NewRouter()
	e := MakeServerEndpoints()

	// POST    /createAccount              creat user account
	// GET     /getUserAccounts            get all user accounts
	// POST     /createTransaction          create account transaction
	// GET  /getAccountTransactions	   get account transactions

	r.Methods("POST").Path("/createAccount").Handler(httptransport.NewServer(
		jwt.NewParser(services.Key, basicjwt.SigningMethodHS256, jwt.StandardClaimsFactory)(e.CreateAccountEndpoint),
		DecodeCreateAccountRequest,
		EncodeCreateAccountResponse,
		[]httptransport.ServerOption{
			httptransport.ServerBefore(jwt.HTTPToContext()),
		}...,
	))
	r.Methods("GET").Path("/getUserAccounts").Handler(httptransport.NewServer(
		jwt.NewParser(services.Key, basicjwt.SigningMethodHS256, jwt.StandardClaimsFactory)(e.GetUserAccountsEndpoint),
		DecodeGetUserAccountsRequest,
		EncodeGetUserAccountsResponse,
		[]httptransport.ServerOption{
			httptransport.ServerBefore(jwt.HTTPToContext()),
		}...,
	))
	r.Methods("POST").Path("/createTransaction").Handler(httptransport.NewServer(
		jwt.NewParser(services.Key, basicjwt.SigningMethodHS256, jwt.StandardClaimsFactory)(e.CreateTransactionEndpoint),
		DecodeCreateTransactionRequest,
		EncodeCreateTransactionResponse,
		[]httptransport.ServerOption{
			httptransport.ServerBefore(jwt.HTTPToContext()),
		}...,
	))
	r.Methods("GET").Path("/getAccountTransactions").Handler(httptransport.NewServer(
		jwt.NewParser(services.Key, basicjwt.SigningMethodHS256, jwt.StandardClaimsFactory)(e.GetAccountTransactionsEndpoint),
		DecodeGetAccountTransactionsRequest,
		EncodeGetAccountTransactionsResponse,
		[]httptransport.ServerOption{
			httptransport.ServerBefore(jwt.HTTPToContext()),
		}...,
	))
	return r
}
