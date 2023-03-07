package database

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func CreateRandomJoke(t *testing.T, username string) Joke {
	arg := CreateJokeParams{
		Author:      username,
		Title:       "my joke",
		Text:        "funny joke :o",
		Explanation: "pretty obvious",
	}

	joke, err := testQueries.CreateJoke(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, joke)

	require.Equal(t, arg.Author, joke.Author)
	require.Equal(t, arg.Title, joke.Title)
	require.Equal(t, arg.Text, joke.Text)
	require.Equal(t, arg.Explanation, joke.Explanation)
	require.NotZero(t, joke.ID)
	require.NotZero(t, joke.CreatedAt)

	return joke
}

func TestCreateJoke(t *testing.T) {
	user := CreateRandomUser(t)
	CreateRandomJoke(t, user.Username)
}

func TestGetJoke(t *testing.T) {
	user := CreateRandomUser(t)
	joke1 := CreateRandomJoke(t, user.Username)
	joke2, err := testQueries.GetJoke(context.Background(), joke1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, joke2)

	require.Equal(t, joke1.ID, joke2.ID)
	require.Equal(t, joke1.Title, joke2.Title)
	require.Equal(t, joke1.Text, joke2.Text)
	require.Equal(t, joke1.Explanation, joke2.Explanation)
	require.WithinDuration(t, joke1.CreatedAt, joke2.CreatedAt, time.Second)
	require.Equal(t, joke1.CreatedAt, joke2.CreatedAt, time.Second)
}

func TestListJokesByAuthor(t *testing.T) {
	user := CreateRandomUser(t)

	for i := 0; i < 15; i++ {
		CreateRandomJoke(t, user.Username)
	}

	arg := ListJokesByAuthorParams{
		Author: user.Username,
		Limit:  5,
		Offset: 5,
	}

	jokes, err := testQueries.ListJokesByAuthor(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, jokes)

	for _, joke := range jokes {
		require.NotEmpty(t, joke)
		require.Equal(t, joke.Author, user.Username)
	}
}

func TestDeleteJoke(t *testing.T) {
	user := CreateRandomUser(t)
	joke1 := CreateRandomJoke(t, user.Username)
	err := testQueries.DeleteJoke(context.Background(), joke1.ID)
	require.NoError(t, err)

	joke2, err := testQueries.GetJoke(context.Background(), joke1.ID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, joke2)
}

func TestDeleteJokesByAuthor(t *testing.T) {
	user := CreateRandomUser(t)

	for i := 0; i < 10; i++ {
		CreateRandomJoke(t, user.Username)
	}

	err := testQueries.DeleteJokesByAuthor(context.Background(), user.Username)
	require.NoError(t, err)

	err = testQueries.DeleteJokesByAuthor(context.Background(), user.Username)
	require.NoError(t, err)

	arg := ListJokesByAuthorParams{
		Author: user.Username,
		Limit:  5,
		Offset: 5,
	}
	jokes, err := testQueries.ListJokesByAuthor(context.Background(), arg)
	require.NoError(t, err)
	require.Empty(t, jokes)
}
