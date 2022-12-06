package repository

import (
	"context"
	"fmt"
	"github.com/OVantsevich/internetBankingCourseProjectGo/accountService/domain"
	"github.com/labstack/gommon/log"
	"strconv"
)

func CreateAccount(ctx context.Context, account *domain.Account) (string, error) {

	if err := Pool(ctx); err != nil {
		return "something went wrong", err
	}

	rows, err := pool.Query(ctx, "insert into accounts (account_name, amount, user_id) "+
		"select $1, $2, $3 where exists(select 1 from users where id = $4 and is_deleted = false) returning account_name",
		account.AccountName, account.Amount, account.UserId, account.UserId)
	if err != nil {
		log.Errorf("database error with create account: %v", rows.Err())
		return "something went wrong", rows.Err()
	}

	if !rows.Next() {
		rows.Close()
		if rows.Err() != nil {
			log.Errorf("database error with create account: %v", rows.Err())
			return "you already have account with this name", rows.Err()
		}
		return "user with this login doesn't exist", fmt.Errorf("user with this login doesn't exist")
	}

	return "Account: " + account.AccountName + " created." + " Balance is: " + strconv.Itoa(account.Amount), nil
}

func GetUserAccounts(ctx context.Context, userId int) ([]domain.Account, string, error) {

	if err := Pool(ctx); err != nil {
		return nil, "something went wrong", err
	}

	rows, err := pool.Query(ctx, "select account_name, amount, creation_date from accounts where user_id=$1 and is_deleted=false", userId)
	defer rows.Close()
	if err != nil {
		log.Errorf("database error with create user: %v", err)
		return nil, "something went wrong", err
	}
	if !rows.Next() {
		return nil, "you don't have accounts yet", fmt.Errorf("you don't have accounts yet")
	}

	var accounts []domain.Account
	var i = 0
	for {
		accounts = append(accounts, domain.Account{})
		err = rows.Scan(&accounts[i].AccountName, &accounts[i].Amount, &accounts[i].CreationDate)
		if err != nil {
			log.Errorf("database error with execution from rows: %v", err)
			return nil, "something went wrong", err
		}
		i++
		if !rows.Next() {
			break
		}
	}
	return accounts, "", nil
}

func GetUserAccountByAccountName(ctx context.Context, userId int, accountName string) (*domain.Account, string, error) {

	if err := Pool(ctx); err != nil {
		return nil, "something went wrong", err
	}

	rows, err := pool.Query(ctx, "select * from accounts where user_id=$1 and account_name=$2 and is_deleted=false", userId, accountName)
	defer rows.Close()
	if err != nil {
		log.Errorf("database error with create user: %v", err)
		return nil, "something went wrong", err
	}
	if !rows.Next() {
		return nil, "you don't have accounts yet", fmt.Errorf("you don't have accounts yet")
	}

	var account domain.Account
	err = rows.Scan(&account.ID, &account.UserId, &account.AccountName, &account.Amount, &account.IsDeleted, &account.CreationDate, &account.ModificationDate)
	if err != nil {
		log.Errorf("database error with execution from rows: %v", err)
		return nil, "something went wrong", err
	}

	return &account, "", nil
}
