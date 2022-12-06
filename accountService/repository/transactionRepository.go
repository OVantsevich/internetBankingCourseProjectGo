package repository

import (
	"context"
	"fmt"
	"github.com/OVantsevich/internetBankingCourseProjectGo/accountService/domain"
	"github.com/labstack/gommon/log"
	"strconv"
)

func CreateTransaction(ctx context.Context, transaction *domain.Transaction) (string, error) {

	if err := Pool(ctx); err != nil {
		return "something went wrong", err
	}

	var id int
	row := pool.QueryRow(ctx, "insert into transactions (account_sender_id, account_receiver_id, amount) values ($1, $2, $3) returning id",
		transaction.AccountSenderId, transaction.AccountReceiverId, transaction.Amount)
	if err := row.Scan(&id); err != nil {
		log.Errorf("database error with create transaction: %v", err)
		return "something went wrong", err
	}

	var isCompleted bool
	if err := pool.QueryRow(ctx, "select is_completed from transactions where id=$1 and is_deleted=false", id).Scan(
		&isCompleted); err != nil {
		log.Errorf("database error with create transaction: %v", err)
		return "something went wrong", err
	}

	if isCompleted {
		return "Transaction amount is: " + strconv.Itoa(transaction.Amount), nil
	}
	var amount int
	if err := pool.QueryRow(ctx, "select amount from accounts where id=$1 and is_deleted=false", transaction.AccountSenderId).Scan(
		&amount); err != nil {
		log.Errorf("database error with create transaction: %v", err)
		return "something went wrong", err
	}
	return "Insufficient funds. Current amount is: " + strconv.Itoa(amount), nil
}

func GetAccountTransactions(ctx context.Context, accountId int) ([]string, []string, []domain.Transaction, string, error) {

	if err := Pool(ctx); err != nil {
		return nil, nil, nil, "something went wrong", err
	}

	rows, err := pool.Query(ctx, "select ("+
		"select account_name from accounts where accounts.id = t.account_sender_id) sender_name,"+
		"(select account_name from accounts where accounts.id = t.account_receiver_id) receiver_name,"+
		"t.amount, t.is_completed, t.creation_date "+
		"from transactions t where t.account_receiver_id = $1 or t.account_sender_id = $2", accountId, accountId)
	defer rows.Close()
	if err != nil {
		log.Errorf("database error with create user: %v", err)
		return nil, nil, nil, "something went wrong", err
	}
	if !rows.Next() {
		return nil, nil, nil, "you don't have transactions on this account yet", fmt.Errorf("you don't have transactions on this account yet")
	}

	var transactions []domain.Transaction
	var senderNames []string
	var receiverNames []string
	var i = 0
	for rows.Next() {
		transactions = append(transactions, domain.Transaction{})
		senderNames = append(senderNames, "")
		receiverNames = append(receiverNames, "")
		err = rows.Scan(&senderNames[i], &receiverNames[i], &transactions[i].Amount, &transactions[i].IsCompleted, &transactions[i].CreationDate)
		if err != nil {
			log.Errorf("database error with execution from rows: %v", err)
			return nil, nil, nil, "something went wrong", err
		}
		i++
	}
	return senderNames, receiverNames, transactions, "", nil
}
