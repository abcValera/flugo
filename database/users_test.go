package database

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/abc_valera/flugo/utils"
	"github.com/stretchr/testify/require"
)

func CreateRandomUser(t *testing.T) User {
	hashedPassword, err := utils.HashPassword(utils.RandomPassword())
	require.NoError(t, err)

	createArgs := CreateUserParams{
		Username:       utils.RandomUsername(),
		Email:          utils.RandomEmail(),
		HashedPassword: hashedPassword,
		Fullname:       utils.RandomFullname(),
		Status:         utils.RandomStatus(),
		Bio:            utils.RandomBio(),
	}

	user, err := testQueries.CreateUser(context.Background(), createArgs)
	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.Equal(t, user.Username, createArgs.Username)
	require.Equal(t, user.Email, createArgs.Email)
	require.Equal(t, user.HashedPassword, createArgs.HashedPassword)
	require.Equal(t, user.Fullname, createArgs.Fullname)
	require.Equal(t, user.Status, createArgs.Status)
	require.Equal(t, user.Bio, createArgs.Bio)

	require.NotZero(t, user.ID)
	require.NotZero(t, user.CreatedAt)

	return user
}

func TestCreateUser(t *testing.T) {
	CreateRandomUser(t)
}

func TestGetUserByID(t *testing.T) {
	user1 := CreateRandomUser(t)
	user2, err := testQueries.GetUserByID(context.Background(), user1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, user2)

	require.Equal(t, user1.ID, user2.ID)
	require.Equal(t, user1.Username, user2.Username)
	require.Equal(t, user1.Email, user2.Email)
	require.Equal(t, user1.HashedPassword, user2.HashedPassword)
	require.Equal(t, user1.Fullname, user2.Fullname)
	require.Equal(t, user1.Status, user2.Status)
	require.Equal(t, user1.Bio, user2.Bio)
	require.WithinDuration(t, user1.CreatedAt, user2.CreatedAt, time.Second)
	require.WithinDuration(t, user1.UpdatedAt, user2.UpdatedAt, time.Second)
}

func TestListUsers(t *testing.T) {
	for i := 0; i < 15; i++ {
		CreateRandomUser(t)
	}

	arg := ListUsersParams{
		Limit:  5,
		Offset: 5,
	}

	users, err := testQueries.ListUsers(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, users)

	for _, user := range users {
		require.NotEmpty(t, user)
	}
}

func TestDeleteUser(t *testing.T) {
	user1 := CreateRandomUser(t)
	err := testQueries.DeleteUser(context.Background(), user1.ID)
	require.NoError(t, err)

	user2, err := testQueries.GetUserByID(context.Background(), user1.ID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, user2)
}
