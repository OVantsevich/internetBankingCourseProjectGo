package server

import (
	"context"
	"encoding/json"
	"github.com/OVantsevich/internetBankingCourseProjectGo/userService/services"
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
		options...,
	))
	r.Methods("DELETE").Path("/deleteUser").Handler(httptransport.NewServer(
		jwt.NewParser(services.Key, basicjwt.SigningMethodHS256, jwt.StandardClaimsFactory)(e.DeleteUserEndpoint),
		DecodeDeleteUserRequest,
		EncodeResponse,
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
