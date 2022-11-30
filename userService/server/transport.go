package server

import (
	"github.com/OVantsevich/internetBankingCourseProjectGo/userService/services"
	"github.com/go-kit/kit/auth/jwt"
	httptransport "github.com/go-kit/kit/transport/http"
	basicjwt "github.com/golang-jwt/jwt/v4"
	"github.com/gorilla/mux"
	"net/http"
)

func MakeHTTPHandler() http.Handler {
	r := mux.NewRouter()
	e := MakeServerEndpoints()

	// POST    /createUser              creat user
	// GET     /signIn                  signIn user, return token
	// PUT     /updateUser              update user
	// DELETE  /deleteUser  			remove user

	r.Methods("POST").Path("/createUser").Handler(httptransport.NewServer(
		e.RegisterUserEndpoint,
		DecodeCreateUserRequest,
		EncodeResponse,
	))
	r.Methods("GET").Path("/signIn").Handler(httptransport.NewServer(
		e.SignInEndpoint,
		DecodeSignInRequest,
		EncodeResponse,
	))
	r.Methods("PUT").Path("/updateUser").Handler(httptransport.NewServer(
		jwt.NewParser(services.Key, basicjwt.SigningMethodHS256, jwt.StandardClaimsFactory)(e.UpdateUserEndpoint),
		DecodeUpdateUserRequest,
		EncodeResponse,
		[]httptransport.ServerOption{
			httptransport.ServerBefore(jwt.HTTPToContext()),
		}...,
	))
	r.Methods("DELETE").Path("/deleteUser").Handler(httptransport.NewServer(
		jwt.NewParser(services.Key, basicjwt.SigningMethodHS256, jwt.StandardClaimsFactory)(e.DeleteUserEndpoint),
		DecodeDeleteUserRequest,
		EncodeResponse,
		[]httptransport.ServerOption{
			httptransport.ServerBefore(jwt.HTTPToContext()),
		}...,
	))
	return r
}
