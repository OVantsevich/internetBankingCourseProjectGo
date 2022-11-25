package main

import (
	"github.com/OVantsevich/internetBankingCourseProjectGo/accountService/repository"
	"github.com/OVantsevich/internetBankingCourseProjectGo/accountService/server"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/labstack/gommon/log"
	"net/http"
)

func main() {
	e := server.MakeServerEndpoints()

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

	http.Handle("/createAccount", createAccount)
	http.Handle("/getUserAccounts", getUserAccounts)
	log.Fatal(http.ListenAndServe(":12344", nil))
	repository.Close()
}
