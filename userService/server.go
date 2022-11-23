package main

import (
	"github.com/OVantsevich/internetBankingCourseProjectGo/userService/repository"
	"github.com/OVantsevich/internetBankingCourseProjectGo/userService/server"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/labstack/gommon/log"
	"net/http"
)

func main() {
	e := server.MakeServerEndpoints()

	registerUser := httptransport.NewServer(
		e.RegisterUserEndpoint,
		server.DecodeCreateUserRequest,
		server.EncodeResponse,
	)

	signIn := httptransport.NewServer(
		e.SignInEndpoint,
		server.DecodeSignInRequest,
		server.EncodeResponse,
	)

	updateUser := httptransport.NewServer(
		e.UpdateUserEndpoint,
		server.DecodeUpdateUserRequest,
		server.EncodeResponse,
	)

	deleteUser := httptransport.NewServer(
		e.DeleteUserEndpoint,
		server.DecodeDeleteUserRequest,
		server.EncodeResponse,
	)

	http.Handle("/registerUser", registerUser)
	http.Handle("/signIn", signIn)
	http.Handle("/updateUser", updateUser)
	http.Handle("/deleteUser", deleteUser)
	log.Fatal(http.ListenAndServe(":12345", nil))
	repository.Close()
}
