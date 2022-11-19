package main

import (
	"context"
	"encoding/json"
	"github.com/OVantsevich/internetBankingCourseProjectGo/userService/domain"
	"github.com/OVantsevich/internetBankingCourseProjectGo/userService/repository"
	"github.com/OVantsevich/internetBankingCourseProjectGo/userService/server"
	"github.com/OVantsevich/internetBankingCourseProjectGo/userService/services"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/jackc/pgx/v4/pgxpool"
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
	databaseUrl := "postgres://postgres:postgres@host.docker.internal:5432/userService"
	dbPool, err := pgxpool.Connect(context.Background(), databaseUrl)
	if err != nil {
		log.Errorf("database connection error: %v", err)
	}
	serv := services.NewService(repository.UserRepository{Pool: dbPool})

	e := server.MakeServerEndpoints(serv)

	registerUser := httptransport.NewServer(
		e.RegisterUserEndpoint,
		DecodeUserRequest,
		EncodeResponse,
	)

	signIn := httptransport.NewServer(
		e.SignInEndpoint,
		DecodeUserRequest,
		EncodeResponse,
	)

	http.Handle("/registerUser", registerUser)
	http.Handle("/signIn", signIn)
	log.Fatal(http.ListenAndServe(":12345", nil))
	dbPool.Close()
}

func DecodeUserRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request domain.User
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

func EncodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}
