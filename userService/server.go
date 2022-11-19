package main

import (
	"github.com/OVantsevich/internetBankingCourseProjectGo/userService/repository"
	"github.com/OVantsevich/internetBankingCourseProjectGo/userService/server"
	"github.com/OVantsevich/internetBankingCourseProjectGo/userService/services"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/labstack/gommon/log"
	"net/http"
)

//func RegisterUserEndpoint(svc services.Service) endpoint.Endpoint {
//	return func(_ context.Context, request interface{}) (interface{}, error) {
//		req := request.(domain.User)
//
//		v, err := svc.CreateUser(&req)
//		if err != nil {
//			return v, nil
//		}
//		return v, nil
//	}
//}
//
//func SignInEndpoint(svc services.Service) endpoint.Endpoint {
//	return func(_ context.Context, request interface{}) (interface{}, error) {
//		req := request.(domain.User)
//
//		v, err := svc.SignIn(req.UserLogin, req.UserPassword)
//		if err != nil {
//			return v, nil
//		}
//		return v, nil
//	}
//}

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

	http.Handle("/registerUser", registerUser)
	http.Handle("/signIn", signIn)
	log.Fatal(http.ListenAndServe(":12345", nil))
	srv.Close()
}
