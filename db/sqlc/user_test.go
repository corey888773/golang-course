package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/corey888773/golang-course/util"
	"github.com/stretchr/testify/require"
)

func createRandomUser(t *testing.T) User {
	hashedPassword, err := util.HashPassword(util.RandomString(6))
	require.NoError(t, err)

	arg := CreateUserParams{
		Username:       util.RandomOwner(),
		HashedPassword: hashedPassword,
		FullName:       util.RandomOwner(),
		Email:          util.RandomEmail(),
	}

	user, err := testQueries.CreateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.Equal(t, arg.Username, user.Username)
	require.Equal(t, arg.HashedPassword, user.HashedPassword)
	require.Equal(t, arg.FullName, user.FullName)
	require.Equal(t, arg.Email, user.Email)

	require.True(t, user.PasswordChangedAt.IsZero())
	require.NotZero(t, user.CreatedAt)

	return user
}

func TestCreateUser(t *testing.T) {
	createRandomUser(t)
}

func TestGetUser(t *testing.T) {
	user1 := createRandomUser(t)
	user2, err := testQueries.GetUser(context.Background(), user1.Username)

	require.NoError(t, err)
	require.NotEmpty(t, user2)

	require.Equal(t, user1.Username, user2.Username)
	require.Equal(t, user1.Email, user2.Email)
	require.Equal(t, user1.HashedPassword, user2.HashedPassword)
	require.Equal(t, user1.FullName, user2.FullName)
	require.WithinDuration(t, user1.CreatedAt, user2.CreatedAt, time.Second)
	require.WithinDuration(t, user1.PasswordChangedAt, user2.PasswordChangedAt, time.Second)
}

func TestUpdateUserOnlyFullName(t *testing.T) {
	oldUser := createRandomUser(t)
	newFullName := util.RandomOwner()

	arg := UpdateUserParams{
		Username: oldUser.Username,
		FullName: sql.NullString{
			String: newFullName,
			Valid:  true,
		},
	}

	user, err := testQueries.UpdateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, user)
	require.Equal(t, user.Username, oldUser.Username)
	require.Equal(t, user.Email, oldUser.Email)
	require.Equal(t, user.HashedPassword, oldUser.HashedPassword)
	require.Equal(t, user.CreatedAt, oldUser.CreatedAt)
	require.Equal(t, user.PasswordChangedAt, oldUser.PasswordChangedAt)
	require.Equal(t, user.FullName, newFullName)
}

func TestUpdateUserOnlyEmail(t *testing.T) {
	oldUser := createRandomUser(t)
	newEmail := util.RandomEmail()

	arg := UpdateUserParams{
		Username: oldUser.Username,
		Email: sql.NullString{
			String: newEmail,
			Valid:  true,
		},
	}

	user, err := testQueries.UpdateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, user)
	require.Equal(t, user.Username, oldUser.Username)
	require.Equal(t, user.FullName, oldUser.FullName)
	require.Equal(t, user.HashedPassword, oldUser.HashedPassword)
	require.Equal(t, user.CreatedAt, oldUser.CreatedAt)
	require.Equal(t, user.PasswordChangedAt, oldUser.PasswordChangedAt)
	require.Equal(t, user.Email, newEmail)
}

func TestUpdateUserOnlyHashedPassword(t *testing.T) {
	oldUser := createRandomUser(t)
	newPassword := util.RandomString(6)
	newHashedPassword, err := util.HashPassword(newPassword)
	require.NoError(t, err)

	arg := UpdateUserParams{
		Username: oldUser.Username,
		HashedPassword: sql.NullString{
			String: newHashedPassword,
			Valid:  true,
		},
	}

	user, err := testQueries.UpdateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, user)
	require.Equal(t, user.Username, oldUser.Username)
	require.Equal(t, user.FullName, oldUser.FullName)
	require.Equal(t, user.Email, oldUser.Email)
	require.Equal(t, user.CreatedAt, oldUser.CreatedAt)
	require.Equal(t, user.PasswordChangedAt, oldUser.PasswordChangedAt)
	require.Equal(t, user.HashedPassword, newHashedPassword)
}

func TestUpdateAllFields(t *testing.T) {
	oldUser := createRandomUser(t)
	newPassword := util.RandomString(6)
	newHashedPassword, err := util.HashPassword(newPassword)
	require.NoError(t, err)

	newFullName := util.RandomOwner()
	newEmail := util.RandomEmail()

	arg := UpdateUserParams{
		Username: oldUser.Username,
		HashedPassword: sql.NullString{
			String: newHashedPassword,
			Valid:  true,
		},
		Email: sql.NullString{
			String: newEmail,
			Valid:  true,
		},
		FullName: sql.NullString{
			String: newFullName,
			Valid:  true,
		},
	}

	user, err := testQueries.UpdateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, user)
	require.Equal(t, user.Username, oldUser.Username)
	require.Equal(t, user.FullName, newFullName)
	require.Equal(t, user.Email, newEmail)
	require.Equal(t, user.CreatedAt, oldUser.CreatedAt)
	require.Equal(t, user.PasswordChangedAt, oldUser.PasswordChangedAt)
	require.Equal(t, user.HashedPassword, newHashedPassword)
}
