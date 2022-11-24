package services

import (
	"context"
	"github.com/OVantsevich/internetBankingCourseProjectGo/userService/domain"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/stretchr/testify/require"
	"testing"
)

var pool *pgxpool.Pool = nil

func Pool(ctx context.Context) error {

	if pool == nil {
		var err error
		pool, err = pgxpool.Connect(ctx, "postgres://postgres:postgres@host.docker.internal:5432/userService?sslmode=disable")
		if err != nil {
			return err
		}
	}
	return nil
}

func TestCreate(t *testing.T) {
	testValidData := []domain.User{
		{
			UserName:     `NAME`,
			Surname:      `SURNAME`,
			UserLogin:    `CreateLOGIN1`,
			UserEmail:    `LOGIN1@gmail.com`,
			UserPassword: `LOGIN23102002`,
		},
		{
			UserName:     `NAME`,
			Surname:      `SURNAME`,
			UserLogin:    `CreateLOGIN2`,
			UserEmail:    `LOGIN2@gmail.com`,
			UserPassword: `PASSWORD123`,
		},
	}
	testNoValidData := []domain.User{
		//----------------------------------------------------------------- Not valid
		{
			Surname:      `SURNAME`,
			UserLogin:    `LOGIN1`,
			UserEmail:    `LOGIN1@gmail.com`,
			UserPassword: `LOGIN23102002`,
		},
		{
			UserName:     `NAME`,
			UserLogin:    `LOGIN2`,
			UserEmail:    `LOGIN2@gmail.com`,
			UserPassword: `PASSWORD123`,
		}, {
			UserName:     `NAME`,
			Surname:      `SURNAME`,
			UserEmail:    `LOGIN1@gmail.com`,
			UserPassword: `LOGIN23102002`,
		},
		{
			UserName:     `NAME`,
			Surname:      `SURNAME`,
			UserLogin:    `LOGIN2`,
			UserPassword: `PASSWORD123`,
		},
		{
			UserName:  `NAME`,
			Surname:   `SURNAME`,
			UserLogin: `LOGIN1`,
			UserEmail: `LOGIN1@gmail.com`,
		},
		{
			UserName:     `222NAME`,
			Surname:      `SURNAME`,
			UserLogin:    `LOGIN2`,
			UserEmail:    `LOGIN2@gmail.com`,
			UserPassword: `PASSWORD123`,
		},
		{
			UserName:     `NAME`,
			Surname:      `222SURNAME`,
			UserLogin:    `LOGIN1`,
			UserEmail:    `LOGIN1@gmail.com`,
			UserPassword: `LOGIN23102002`,
		},
		{
			UserName:     `NAME`,
			Surname:      `SURNAME`,
			UserLogin:    `LOG`,
			UserEmail:    `LOGIN2@gmail.com`,
			UserPassword: `PASSWORD123`,
		},
		{
			UserName:     `NAME`,
			Surname:      `SURNAME`,
			UserLogin:    `LOGIN123456789012345`,
			UserEmail:    `LOGIN1@gmail.com`,
			UserPassword: `LOGIN23102002`,
		},
		{
			UserName:     `NAME`,
			Surname:      `SURNAME`,
			UserLogin:    `LOGIN2`,
			UserEmail:    `LOGIN2`,
			UserPassword: `PASSWORD123`,
		},
		{
			UserName:     `NAME`,
			Surname:      `SURNAME`,
			UserLogin:    `LOGIN1`,
			UserEmail:    `@gmail`,
			UserPassword: `LOGIN23102002`,
		},
		{
			UserName:     `NAME`,
			Surname:      `SURNAME`,
			UserLogin:    `LOGIN2`,
			UserEmail:    `LOGIN2@gmail.com`,
			UserPassword: `weak`,
		},
		{
			UserName:     `NAME`,
			Surname:      `SURNAME`,
			UserLogin:    `LOGIN1`,
			UserEmail:    `LOGIN1@gmail.com`,
			UserPassword: `11111111111111111111111111111111111111`,
		},
		{
			UserName:     `NAME`,
			Surname:      `SURNAME`,
			UserLogin:    `LOGIN2`,
			UserEmail:    `LOGIN2@gmail.com`,
			UserPassword: `password`,
		},
		//----------------------------------------------------------------- Already exist
		{
			UserName:     `NAME`,
			Surname:      `SURNAME`,
			UserLogin:    `CreateLOGIN1`,
			UserEmail:    `LOGIN1@gmail.com`,
			UserPassword: `LOGIN23102002`,
		},
		{
			UserName:     `NAME`,
			Surname:      `SURNAME`,
			UserLogin:    `CreateLOGIN2`,
			UserEmail:    `LOGIN2@gmail.com`,
			UserPassword: `PASSWORD123`,
		},
	}

	ctx := context.Background()
	err := Pool(ctx)
	require.NoError(t, err, "create error")

	for _, u := range testValidData {

		_, err := pool.Exec(ctx, "delete from users where user_login=$1 ", u.UserLogin)
		require.NoError(t, err, "delete error")

		_, err = CreateUser(ctx, &u)
		require.NoError(t, err, "create error")
	}
	for _, u := range testNoValidData {

		_, err := CreateUser(ctx, &u)
		require.Error(t, err, "create error "+u.UserLogin)
	}

	for _, u := range testNoValidData {

		_, err = pool.Exec(ctx, "delete from users where user_login=$1 ", u.UserLogin)
		require.NoError(t, err, "delete error")
	}
}

func TestSignIn(t *testing.T) {
	testValidData := []domain.User{
		{
			UserName:     `NAME`,
			Surname:      `SURNAME`,
			UserLogin:    `SignInLOGIN1`,
			UserEmail:    `LOGIN1@gmail.com`,
			UserPassword: `LOGIN23102002`,
		},
		{
			UserName:     `NAME`,
			Surname:      `SURNAME`,
			UserLogin:    `SignInLOGIN2`,
			UserEmail:    `LOGIN2@gmail.com`,
			UserPassword: `PASSWORD123`,
		},
	}
	testNoValidPassword := domain.User{
		UserName:     `NAME`,
		Surname:      `SURNAME`,
		UserLogin:    `SignInLOGIN1`,
		UserEmail:    `LOGIN1@gmail.com`,
		UserPassword: `LOGIN23102002`,
	}
	testNoValidLogin := domain.User{
		UserName:     `NAME`,
		Surname:      `SURNAME`,
		UserLogin:    `SignInLOGIN2`,
		UserEmail:    `LOGIN2@gmail.com`,
		UserPassword: `PASSWORD123`,
	}
	testNoExistData := []domain.User{
		{
			UserName:     `NAME`,
			Surname:      `SURNAME`,
			UserLogin:    `SignInLOGIN1NE`,
			UserEmail:    `LOGIN1@gmail.com`,
			UserPassword: `LOGIN23102002`,
		},
		{
			UserName:     `NAME`,
			Surname:      `SURNAME`,
			UserLogin:    `SignInLOGIN2NE`,
			UserEmail:    `LOGIN2@gmail.com`,
			UserPassword: `PASSWORD123`,
		},
	}
	testDeletedData := []domain.User{
		{
			UserName:     `NAME`,
			Surname:      `SURNAME`,
			UserLogin:    `SignInLOGIN1NE`,
			UserEmail:    `LOGIN1@gmail.com`,
			UserPassword: `LOGIN23102002`,
		},
		{
			UserName:     `NAME`,
			Surname:      `SURNAME`,
			UserLogin:    `SignInLOGIN2NE`,
			UserEmail:    `LOGIN2@gmail.com`,
			UserPassword: `PASSWORD123`,
		},
	}

	ctx := context.Background()
	err := Pool(ctx)
	require.NoError(t, err, "create error")

	for _, u := range testValidData {

		_, err := pool.Exec(ctx, "delete from users where user_login=$1 ", u.UserLogin)
		require.NoError(t, err, "delete error")

		password := u.UserPassword
		_, err = CreateUser(ctx, &u)
		require.NoError(t, err, "create error")

		_, err = SignIn(ctx, &SignInRequest{UserLogin: u.UserLogin, UserPassword: password})
		require.NoError(t, err, "signIn error")

		_, err = pool.Exec(ctx, "delete from users where user_login=$1 ", u.UserLogin)
		require.NoError(t, err, "delete error")
	}

	_, err = pool.Exec(ctx, "delete from users where user_login=$1 ", testNoValidLogin.UserLogin)
	require.NoError(t, err, "delete error")

	password := testNoValidLogin.UserPassword
	_, err = CreateUser(ctx, &testNoValidLogin)
	require.NoError(t, err, "create error")

	_, err = SignIn(ctx, &SignInRequest{UserLogin: "", UserPassword: password})
	require.Error(t, err, "signIn error")

	_, err = pool.Exec(ctx, "delete from users where user_login=$1 ", testNoValidLogin.UserLogin)
	require.NoError(t, err, "delete error")

	_, err = pool.Exec(ctx, "delete from users where user_login=$1 ", testNoValidPassword.UserLogin)
	require.NoError(t, err, "delete error")

	_, err = CreateUser(ctx, &testNoValidPassword)
	require.NoError(t, err, "create error")

	_, err = SignIn(ctx, &SignInRequest{UserLogin: testNoValidPassword.UserLogin, UserPassword: ""})
	require.Error(t, err, "signIn error")

	_, err = pool.Exec(ctx, "delete from users where user_login=$1 ", testNoValidPassword.UserLogin)
	require.NoError(t, err, "delete error")

	for _, u := range testNoExistData {

		_, err := pool.Exec(ctx, "delete from users where user_login=$1 ", u.UserLogin)
		require.NoError(t, err, "delete error")

		_, err = SignIn(ctx, &SignInRequest{UserLogin: u.UserLogin, UserPassword: u.UserPassword})
		require.Error(t, err, "signIn error")
	}

	for _, u := range testDeletedData {

		_, err := pool.Exec(ctx, "delete from users where user_login=$1 ", u.UserLogin)
		require.NoError(t, err, "delete error")

		password := u.UserPassword
		_, err = CreateUser(ctx, &u)
		require.NoError(t, err, "create error")

		_, err = pool.Exec(ctx, "update users set is_deleted=$1 where user_login=$2", true, u.UserLogin)
		require.NoError(t, err, "delete error")

		_, err = SignIn(ctx, &SignInRequest{UserLogin: u.UserLogin, UserPassword: password})
		require.Error(t, err, "signIn error")

		_, err = pool.Exec(ctx, "delete from users where user_login=$1 ", u.UserLogin)
		require.NoError(t, err, "delete error")
	}
}

func TestUpdateUser(t *testing.T) {
	ctx := context.Background()
	err := Pool(ctx)
	require.NoError(t, err, "create error")

	testValidDataUser := domain.User{
		UserName:     `NAME`,
		Surname:      `SURNAME`,
		UserLogin:    `UpdateLOGIN1`,
		UserEmail:    `LOGIN1@gmail.com`,
		UserPassword: `LOGIN23102002`,
	}
	testValidDataUpdate := []UpdateUserRequest{
		{
			UserName: `UpdateNAME`,
			Token:    "",
		},
		{
			Surname: `UpdateSURNAME`,
			Token:   "",
		},
		{
			UserPassword: `UPDATE23102002`,
			Token:        "",
		},
		{
			UserEmail: `UpdateLOGIN1@gmail.com`,
			Token:     "",
		},
		{
			UserName:     `NAME`,
			Surname:      `SURNAME`,
			UserEmail:    `LOGIN1@gmail.com`,
			UserPassword: `LOGIN23102002`,
			Token:        "",
		},
	}
	_, err = pool.Exec(ctx, "delete from users where user_login=$1 ", testValidDataUser.UserLogin)
	require.NoError(t, err, "delete error")

	password := testValidDataUser.UserPassword
	_, err = CreateUser(ctx, &testValidDataUser)
	require.NoError(t, err, "create error")

	token, err := SignIn(ctx, &SignInRequest{UserLogin: testValidDataUser.UserLogin, UserPassword: password})
	require.NoError(t, err, "signIn error")

	for _, u := range testValidDataUpdate {
		u.Token = token
		_, err = UpdateUser(ctx, &u)
		require.NoError(t, err, "update error")
	}

	_, err = pool.Exec(ctx, "delete from users where user_login=$1 ", testValidDataUser.UserLogin)
	require.NoError(t, err, "delete error")

	testNoValidDataUser := domain.User{
		UserName:     `NAME`,
		Surname:      `SURNAME`,
		UserLogin:    `UpdateLOGIN1`,
		UserEmail:    `LOGIN1@gmail.com`,
		UserPassword: `LOGIN23102002`,
	}
	testNoValidDataUpdate := []UpdateUserRequest{
		{
			UserName:     `222NAME`,
			Surname:      `SURNAME`,
			UserEmail:    `LOGIN2@gmail.com`,
			UserPassword: `PASSWORD123`,
			Token:        "",
		},
		{
			UserName:     `NAME`,
			Surname:      `222SURNAME`,
			UserEmail:    `LOGIN1@gmail.com`,
			UserPassword: `LOGIN23102002`,
			Token:        "",
		},
		{
			UserName:     `NAME`,
			Surname:      `SURNAME`,
			UserEmail:    `LOGIN2`,
			UserPassword: `PASSWORD123`,
			Token:        "",
		},
		{
			UserName:     `NAME`,
			Surname:      `SURNAME`,
			UserEmail:    `@gmail`,
			UserPassword: `LOGIN23102002`,
			Token:        "",
		},
		{
			UserName:     `NAME`,
			Surname:      `SURNAME`,
			UserEmail:    `LOGIN2@gmail.com`,
			UserPassword: `weak`,
			Token:        "",
		},
		{
			UserName:     `NAME`,
			Surname:      `SURNAME`,
			UserEmail:    `LOGIN1@gmail.com`,
			UserPassword: `11111111111111111111111111111111111111`,
			Token:        "",
		},
		{
			UserName:     `NAME`,
			Surname:      `SURNAME`,
			UserEmail:    `LOGIN2@gmail.com`,
			UserPassword: `password`,
			Token:        "",
		},
	}
	_, err = pool.Exec(ctx, "delete from users where user_login=$1 ", testNoValidDataUser.UserLogin)
	require.NoError(t, err, "delete error")

	password = testNoValidDataUser.UserPassword
	_, err = CreateUser(ctx, &testNoValidDataUser)
	require.NoError(t, err, "create error")

	token, err = SignIn(ctx, &SignInRequest{UserLogin: testNoValidDataUser.UserLogin, UserPassword: password})
	require.NoError(t, err, "signIn error")

	for _, u := range testNoValidDataUpdate {
		u.Token = token
		_, err = UpdateUser(ctx, &u)
		require.Error(t, err, "update error")
	}

	_, err = pool.Exec(ctx, "delete from users where user_login=$1 ", testNoValidDataUser.UserLogin)
	require.NoError(t, err, "delete error")

	testNoValidTokenUser := domain.User{
		UserName:     `NAME`,
		Surname:      `SURNAME`,
		UserLogin:    `UpdateLOGIN1`,
		UserEmail:    `LOGIN1@gmail.com`,
		UserPassword: `LOGIN23102002`,
	}
	testNoValidTokenUpdate := UpdateUserRequest{
		UserName:     `NAME`,
		Surname:      `SURNAME`,
		UserEmail:    `LOGIN1@gmail.com`,
		UserPassword: `LOGIN23102002`,
		Token:        "",
	}
	_, err = pool.Exec(ctx, "delete from users where user_login=$1 ", testNoValidTokenUser.UserLogin)
	require.NoError(t, err, "delete error")

	password = testNoValidTokenUser.UserPassword
	_, err = CreateUser(ctx, &testNoValidTokenUser)
	require.NoError(t, err, "create error")

	_, err = UpdateUser(ctx, &testNoValidTokenUpdate)
	require.Error(t, err, "update error")

	_, err = pool.Exec(ctx, "delete from users where user_login=$1 ", testNoValidTokenUser.UserLogin)
	require.NoError(t, err, "delete error")
}

func TestDeleteUser(t *testing.T) {
	ctx := context.Background()
	err := Pool(ctx)
	require.NoError(t, err, "create error")

	testValidData := []domain.User{
		{
			UserName:     `NAME`,
			Surname:      `SURNAME`,
			UserLogin:    `SignInLOGIN1`,
			UserEmail:    `LOGIN1@gmail.com`,
			UserPassword: `LOGIN23102002`,
		},
	}
	for _, u := range testValidData {

		_, err := pool.Exec(ctx, "delete from users where user_login=$1 ", u.UserLogin)
		require.NoError(t, err, "delete error")

		password := u.UserPassword
		_, err = CreateUser(ctx, &u)
		require.NoError(t, err, "create error")

		token, err := SignIn(ctx, &SignInRequest{UserLogin: u.UserLogin, UserPassword: password})
		require.NoError(t, err, "signIn error")

		_, err = DeleteUser(ctx, &DeleteUserRequest{Token: token})
		require.NoError(t, err, "signIn error")

		_, err = pool.Exec(ctx, "delete from users where user_login=$1 ", u.UserLogin)
		require.NoError(t, err, "delete error")
	}

	_, err = DeleteUser(ctx, &DeleteUserRequest{Token: ""})
	require.Error(t, err, "delete error")

}
