package main

import (
	"github.com/OVantsevich/internetBankingCourseProjectGo/accountService/eventStreaming"
	"github.com/OVantsevich/internetBankingCourseProjectGo/accountService/repository"
	"github.com/OVantsevich/internetBankingCourseProjectGo/accountService/server"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/labstack/gommon/log"
	"net/http"
)

func main() {
	e := server.MakeServerEndpoints()

	err := eventStreaming.JetStreamInit()
	if err != nil {
		log.Errorf("error with init stream: %v", err)
	}
	_, err = eventStreaming.JetStream.Subscribe(eventStreaming.SubjectNameUserCreated, eventStreaming.CreatingUserHandler)
	if err != nil {
		log.Printf("Subscribe for " + eventStreaming.SubjectNameUserCreated + " failed")
	}

	createAccount := httptransport.NewServer(
		e.CreateAccountEndpoint,
		server.DecodeCreateAccountRequest,
		server.EncodeCreateAccountResponse,
	)

	getUserAccounts := httptransport.NewServer(
		e.GetUserAccountsEndpoint,
		server.DecodeGetUserAccountsRequest,
		server.EncodeGetUserAccountsResponse,
	)

	CreateTransaction := httptransport.NewServer(
		e.CreateTransactionEndpoint,
		server.DecodeCreateTransactionRequest,
		server.EncodeCreateTransactionResponse,
	)

	getAccountTransactions := httptransport.NewServer(
		e.GetAccountTransactionsEndpoint,
		server.DecodeGetAccountTransactionsRequest,
		server.EncodeGetAccountTransactionsResponse,
	)

	http.Handle("/createAccount", createAccount)
	http.Handle("/getUserAccounts", getUserAccounts)
	http.Handle("/createTransaction", CreateTransaction)
	http.Handle("/getAccountTransactions", getAccountTransactions)
	log.Fatal(http.ListenAndServe(":12344", nil))
	repository.Close()
}
