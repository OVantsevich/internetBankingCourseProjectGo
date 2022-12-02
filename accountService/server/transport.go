package server

import (
	"context"
	"encoding/json"
	"github.com/OVantsevich/internetBankingCourseProjectGo/accountService/services"
	"github.com/go-kit/kit/auth/jwt"
	httptransport "github.com/go-kit/kit/transport/http"
	basicjwt "github.com/golang-jwt/jwt/v4"
	"github.com/gorilla/mux"
	"github.com/labstack/gommon/log"
	"net/http"
	"strconv"
)

func MakeHTTPHandler() http.Handler {
	r := mux.NewRouter()
	e := MakeServerEndpoints()
	options := []httptransport.ServerOption{
		httptransport.ServerBefore(jwt.HTTPToContext()),
		httptransport.ServerErrorEncoder(ErrorEncoder),
	}
	// POST    /createAccount              creat user account
	// GET     /getUserAccounts            get all user accounts
	// POST     /createTransaction          create account transaction
	// GET  /getAccountTransactions	   get account transactions

	r.Methods("POST").Path("/createAccount").Handler(httptransport.NewServer(
		jwt.NewParser(services.Key, basicjwt.SigningMethodHS256, jwt.StandardClaimsFactory)(e.CreateAccountEndpoint),
		DecodeCreateAccountRequest,
		EncodeCreateAccountResponse,
		options...,
	))
	r.Methods("GET").Path("/getUserAccounts").Handler(httptransport.NewServer(
		jwt.NewParser(services.Key, basicjwt.SigningMethodHS256, jwt.StandardClaimsFactory)(e.GetUserAccountsEndpoint),
		DecodeGetUserAccountsRequest,
		EncodeGetUserAccountsResponse,
		options...,
	))
	r.Methods("POST").Path("/createTransaction").Handler(httptransport.NewServer(
		jwt.NewParser(services.Key, basicjwt.SigningMethodHS256, jwt.StandardClaimsFactory)(e.CreateTransactionEndpoint),
		DecodeCreateTransactionRequest,
		EncodeCreateTransactionResponse,
		options...,
	))
	r.Methods("GET").Path("/getAccountTransactions").Handler(httptransport.NewServer(
		jwt.NewParser(services.Key, basicjwt.SigningMethodHS256, jwt.StandardClaimsFactory)(e.GetAccountTransactionsEndpoint),
		DecodeGetAccountTransactionsRequest,
		EncodeGetAccountTransactionsResponse,
		options...,
	))
	return r
}

func ErrorEncoder(_ context.Context, err error, w http.ResponseWriter) {
	code := http.StatusInternalServerError
	var msg string = ""

	log.Errorf("server error: %v", err)

	if sc, ok := err.(httptransport.StatusCoder); ok {
		msg = strconv.Itoa(sc.StatusCode())
	}

	if msg != "" {
		json.NewEncoder(w).Encode(msg)
	}
	w.WriteHeader(code)
}
