package services

import (
	"context"
	"github.com/OVantsevich/internetBankingCourseProjectGo/accountService/repository"
	"github.com/golang-jwt/jwt/v4"
	"time"
)

type CreateTransactionRequest struct {
	AccountSenderNumber   string `json:"account_sender_number"`
	AccountReceiverNumber string `json:"account_receiver_number"`
	Amount                int    `json:"amount"`
	Token                 string `json:"token"`
}

type GetAccountTransactionsRequest struct {
	AccountNumber string `json:"account_number"`
	Token         string `json:"token"`
}

type GetAccountTransactionsResponse struct {
	AccountSenderName   string    `json:"account_sender_name"`
	AccountReceiverName string    `json:"account_receiver_name"`
	Amount              int       `json:"amount"`
	CreationDate        time.Time `json:"creation_date"`
	IsCompleted         bool      `json:"is_completed"`
}

func CreateTransaction(ctx context.Context, request *CreateTransactionRequest) (string, error) {
	claims := jwt.MapClaims{}
	str, err := ParseToken(request.Token, &claims)
	if err != nil {
		return str, err
	}

	return repository.CreateTransaction(ctx, request.Amount, request.AccountSenderNumber, request.AccountReceiverNumber)
}

func GetAccountTransactions(ctx context.Context, request *GetAccountTransactionsRequest) ([]GetAccountTransactionsResponse, string, error) {
	claims := jwt.MapClaims{}
	str, err := ParseToken(request.Token, &claims)
	if err != nil {
		return nil, str, err
	}

	senderNumbers, receiverNumber, transaction, str, err := repository.GetAccountTransactions(ctx, claims["login"].(string), request.AccountNumber)

	var response []GetAccountTransactionsResponse
	if err == nil {
		for i := range transaction {
			response = append(response, GetAccountTransactionsResponse{
				AccountSenderName:   senderNumbers[i],
				AccountReceiverName: receiverNumber[i],
				Amount:              transaction[i].Amount,
				CreationDate:        transaction[i].CreationDate,
			})
		}
	}

	return response, str, err
}
