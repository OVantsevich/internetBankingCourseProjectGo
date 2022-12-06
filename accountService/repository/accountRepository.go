package repository

import (
	"context"
	"fmt"
	"github.com/OVantsevich/internetBankingCourseProjectGo/accountService/domain"
	"github.com/jackc/pgx/v4"
	"github.com/labstack/gommon/log"
	"strconv"
	"time"
)

func CreateAccount(ctx context.Context, account *domain.Account, userLogin string) (string, error) {

	if err := Pool(ctx); err != nil {
		return "something went wrong", err
	}

	var name string
	accountNumber := strconv.FormatInt(time.Now().Unix(), 10)
	if err := pool.QueryRow(ctx, "insert into accounts (user_id, account_name, amount, account_number)"+
		"select (select id from users where user_login = $1 and is_deleted = false), $2, $3, $4 "+
		"where exists(select 1 from users where user_login = $5 and is_deleted = false)"+
		"returning account_name",
		userLogin, account.AccountName, account.Amount, accountNumber, userLogin).Scan(&name); err != nil {
		if err == pgx.ErrNoRows {
			log.Errorf("database error with create account: %v", err)
			return "user with this login doesn't exist", err
		}
		log.Errorf("database error with create account: %v", err)
		return "account with this name already exist", err
	}

	return "Account: " + name + " created." + " Balance is: " + strconv.Itoa(account.Amount), nil
}

func GetUserAccounts(ctx context.Context, userLogin string) ([]domain.Account, string, error) {

	if err := Pool(ctx); err != nil {
		return nil, "something went wrong", err
	}

	rows, err := pool.Query(ctx, "select account_name, amount, creation_date, account_number "+
		"from accounts "+
		"where user_id=(select id from users where user_login = $1 and is_deleted = false)", userLogin)
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
		err = rows.Scan(&accounts[i].AccountName, &accounts[i].Amount, &accounts[i].CreationDate, &accounts[i].AccountNumber)
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
