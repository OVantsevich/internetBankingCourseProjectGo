package main

import (
	"github.com/OVantsevich/internetBankingCourseProjectGo/userService/repository"
	"github.com/OVantsevich/internetBankingCourseProjectGo/userService/server"
	"github.com/OVantsevich/internetBankingCourseProjectGo/userService/services"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/labstack/gommon/log"
	"net/http"
)

func main() {
	srv := services.NewService(repository.UserRepository{})

	e := server.MakeServerEndpoints(srv)

	registerUser := httptransport.NewServer(
		e.RegisterUserEndpoint,
		server.DecodeUserRequest,
		server.EncodeResponse,
	)

	signIn := httptransport.NewServer(
		e.SignInEndpoint,
		server.DecodeUserRequest,
		server.EncodeResponse,
	)

	updateUser := httptransport.NewServer(
		e.UpdateUserEndpoint,
		server.DecodeUpdateUserRequest,
		server.EncodeResponse,
	)

	http.Handle("/registerUser", registerUser)
	http.Handle("/signIn", signIn)
	http.Handle("/updateUser", updateUser)
	log.Fatal(http.ListenAndServe(":12345", nil))
	srv.Close()
}
