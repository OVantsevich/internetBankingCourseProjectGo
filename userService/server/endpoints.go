package server

import (
	"context"
	"github.com/OVantsevich/internetBankingCourseProjectGo/userService/domain"
	"github.com/OVantsevich/internetBankingCourseProjectGo/userService/services"
	"github.com/go-kit/kit/endpoint"
)

type Endpoints struct {
	RegisterUserEndpoint endpoint.Endpoint
	VerificationEndpoint endpoint.Endpoint
	SignInEndpoint       endpoint.Endpoint
	UpdateUserEndpoint   endpoint.Endpoint
	DeleteUserEndpoint   endpoint.Endpoint
}

func MakeServerEndpoints() Endpoints {
	return Endpoints{
		RegisterUserEndpoint: MakeRegisterUserEndpoint(),
		VerificationEndpoint: MakeVerificationEndpoint(),
		SignInEndpoint:       MakeSignInEndpoint(),
		UpdateUserEndpoint:   MakeUpdateUserEndpoint(),
		DeleteUserEndpoint:   MakeDeleteUserEndpoint(),
	}
}

func MakeRegisterUserEndpoint() endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(domain.User)
		response, _ := services.CreateUser(ctx, &req)
		return response, nil
	}
}

func MakeVerificationEndpoint() endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(services.VerificationRequest)
		response, _ := services.Verification(ctx, &req)
		return response, nil
	}
}

func MakeSignInEndpoint() endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(services.SignInRequest)
		response, _ := services.SignIn(ctx, &req)
		return response, nil
	}
}

func MakeUpdateUserEndpoint() endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(services.UpdateUserRequest)
		response, _ := services.UpdateUser(ctx, &req)
		return response, nil
	}
}

func MakeDeleteUserEndpoint() endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(services.DeleteUserRequest)
		response, _ := services.DeleteUser(ctx, &req)
		return response, nil
	}
}
