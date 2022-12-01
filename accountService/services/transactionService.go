package services

import (
	"context"
	"github.com/OVantsevich/internetBankingCourseProjectGo/accountService/domain"
	"github.com/OVantsevich/internetBankingCourseProjectGo/accountService/repository"
	"github.com/golang-jwt/jwt/v4"
	"time"
)

type CreateTransactionRequest struct {
	AccountSenderName   string `json:"account_sender_name"`
	AccountReceiverName string `json:"account_receiver_name"`
	Amount              int    `json:"amount"`
	Token               string `json:"token"`
}

type GetAccountTransactionsRequest struct {
	AccountName string `json:"account_name"`
	Token       string `json:"token"`
}

type GetAccountTransactionsResponse struct {
	AccountSenderName   string    `json:"account_sender_name"`
	AccountReceiverName string    `json:"account_receiver_name"`
	Amount              int       `json:"amount"`
	CreationDate        time.Time `json:"creation_date"`
	IsCompleted         bool      `json:"is_completed"`
}

func CreateTransaction(ctx context.Context, request *CreateTransactionRequest) (string, error) {
	if str, err := domain.InitConfig(); err != nil {
		return str, err
	}

	if str, err := ValidAccountName(request.AccountSenderName, "sender"); err != nil {
		return str, err
	}
	if str, err := ValidAccountName(request.AccountReceiverName, "receiver"); err != nil {
		return str, err
	}

	claims := jwt.MapClaims{}
	str, err := ParseToken(request.Token, &claims)
	if err != nil {
		return str, err
	}

	user, str, err := repository.GetUserByLogin(ctx, claims["login"].(string))
	if err != nil {
		return str, nil
	}

	accountSender, str, err := repository.GetUserAccountByAccountName(ctx, user.ID, request.AccountSenderName)
	if accountSender == nil {
		return str, err
	}
	accountReceiver, str, err := repository.GetUserAccountByAccountName(ctx, user.ID, request.AccountReceiverName)
	if accountReceiver == nil {
		return str, err
	}

	return repository.CreateTransaction(ctx, &domain.Transaction{
		AccountSenderId: accountSender.ID, AccountReceiverId: accountReceiver.ID, Amount: request.Amount})
}

func GetAccountTransactions(ctx context.Context, request *GetAccountTransactionsRequest) ([]GetAccountTransactionsResponse, string, error) {
	if str, err := domain.InitConfig(); err != nil {
		return nil, str, err
	}

	claims := jwt.MapClaims{}
	str, err := ParseToken(request.Token, &claims)
	if err != nil {
		return nil, str, err
	}

	user, str, err := repository.GetUserByLogin(ctx, claims["login"].(string))
	if err != nil {
		return nil, str, err
	}

	account, str, err := repository.GetUserAccountByAccountName(ctx, user.ID, request.AccountName)
	if account == nil {
		return nil, str, err
	}

	senderNames, receiverNames, transaction, str, err := repository.GetAccountTransactions(ctx, account.ID)

	var response []GetAccountTransactionsResponse
	for i := range transaction {
		response = append(response, GetAccountTransactionsResponse{
			AccountSenderName:   senderNames[i],
			AccountReceiverName: receiverNames[i],
			Amount:              transaction[i].Amount,
			CreationDate:        transaction[i].CreationDate,
			IsCompleted:         transaction[i].IsCompleted,
		})
	}

	return response, "", nil
}
