package repository

import (
	"context"
	"fmt"
	"github.com/OVantsevich/internetBankingCourseProjectGo/accountService/domain"
	"github.com/labstack/gommon/log"
	"strconv"
)

func CreateTransaction(ctx context.Context, amount int, accountSenderNumber, accountReceiverNumber string) (string, error) {

	if err := Pool(ctx); err != nil {
		return "something went wrong", err
	}

	var id int
	if err := pool.QueryRow(ctx, "insert into transactions (account_sender_id, account_receiver_id, amount) "+
		"select (select id from accounts where account_number = $1 and is_deleted = false), "+
		"(select id from accounts where account_number = $2 and is_deleted = false), "+
		"$3 "+
		"returning id ",
		accountSenderNumber, accountReceiverNumber, amount).Scan(&id); err != nil {
		log.Errorf("database error with create transaction: %v", err)
		return "Insufficient funds.", err
	}

	return "Transaction amount is: " + strconv.Itoa(amount), nil
}

func GetAccountTransactions(ctx context.Context, login, accountNumber string) ([]string, []string, []domain.Transaction, string, error) {

	if err := Pool(ctx); err != nil {
		return nil, nil, nil, "something went wrong", err
	}

	rows, err := pool.Query(ctx, "select (select account_number from accounts where accounts.id = t.account_sender_id)   sender_number, "+
		"(select account_number from accounts where accounts.id = t.account_receiver_id) receiver_number, "+
		"t.amount, "+
		"t.creation_date "+
		"from transactions t "+
		"where exists(select 1 "+
		"from accounts "+
		"where user_id = (select id from users where user_login = $1 and is_deleted = false) "+
		"and account_number = $2 "+
		"and is_deleted = false) "+
		"and (t.account_receiver_id = (select id from accounts where account_number = $3 and accounts.is_deleted = false) "+
		"or t.account_sender_id = (select id from accounts where account_number = $4 and accounts.is_deleted = false)) "+
		"and is_deleted = false", login, accountNumber, accountNumber, accountNumber)
	defer rows.Close()
	if err != nil {
		log.Errorf("database error with create user: %v", err)
		return nil, nil, nil, "something went wrong", err
	}
	if !rows.Next() {
		return nil, nil, nil, "you don't have transactions on this account yet", fmt.Errorf("you don't have transactions on this account yet")
	}

	var transactions []domain.Transaction
	var senderNumbers []string
	var receiverNumbers []string
	var i = 0
	for rows.Next() {
		transactions = append(transactions, domain.Transaction{})
		senderNumbers = append(senderNumbers, "")
		receiverNumbers = append(receiverNumbers, "")
		err = rows.Scan(&senderNumbers[i], &receiverNumbers[i], &transactions[i].Amount, &transactions[i].CreationDate)
		if err != nil {
			log.Errorf("database error with execution from rows: %v", err)
			return nil, nil, nil, "something went wrong", err
		}
		i++
	}
	return senderNumbers, receiverNumbers, transactions, "", nil
}
