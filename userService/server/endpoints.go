package server

import (
	"context"
	"github.com/OVantsevich/internetBankingCourseProjectGo/userService/domain"
	"github.com/OVantsevich/internetBankingCourseProjectGo/userService/services"
	"github.com/go-kit/kit/endpoint"
)

type Endpoints struct {
	RegisterUserEndpoint endpoint.Endpoint
	SignInEndpoint       endpoint.Endpoint
}

func MakeServerEndpoints(s *services.Service) Endpoints {
	return Endpoints{
		RegisterUserEndpoint: MakeRegisterUserEndpoint(s),
		SignInEndpoint:       MakeSignInEndpoint(s),
	}
}

func MakeRegisterUserEndpoint(s *services.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(domain.User)
		response, _ := s.CreateUser(ctx, &req)
		return response, nil
	}
}

func MakeSignInEndpoint(s *services.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(domain.User)
		response, _ := s.SignIn(ctx, &req)
		return response, nil
	}
}
