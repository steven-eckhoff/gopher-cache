package command

import (
	"context"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopher-cache/internal/common/emulators"
	"gopher-cache/internal/games/adapters"
	"gopher-cache/internal/games/domain/game"
	"testing"
)

func TestCreateGameHandler_Handle(t *testing.T) {
	ctx := context.Background()

	client, cleanup := emulators.NewFirestoreClient(ctx)
	defer func() {
		_ = client.Close()
		cleanup()
	}()

	repo, err := adapters.NewFirestoreGameRepository(client)
	if err != nil {
		panic(err)
	}

	createGameHandler := NewCreateGameHandler(repo)

	userID, err := uuid.NewRandom()
	require.NoError(t, err)

	userNumber := "15734497033"

	user, err := game.NewUser(userID.String(), userNumber)
	require.NoError(t, err)

	createGame := CreateGame{
		Creator:     user,
		Title:       "An Awesome Game",
		Description: "This is an awesome game",
		Levels: []GameLevel{
			{
				Title:       "Level One",
				Description: "This is Level One",
				Clues:       []string{"Level One Clue One", "Level One Clue Two", "Level One Clue Three"},
				Answers:     []string{"Level One is the best"},
			},
			{
				Title:       "Level Two",
				Description: "This is Level Two",
				Clues:       []string{"Level Two Clue One", "Level Two Clue Two", "Level Two Clue Three"},
				Answers:     []string{"Level Two is the best"},
			},
			{
				Title:       "Level Three",
				Description: "This is Level Three",
				Clues:       []string{"Level Three Clue One", "Level Three Clue Two", "Level Three Clue Three"},
				Answers:     []string{"Level Three is the best"},
			},
		},
		Ending:  "The end",
		Kind:    "urban",
		City:    "Austin",
		State:   "Texas",
		Country: "USA",
	}

	err = createGameHandler.Handle(ctx, createGame)
	assert.NoError(t, err)

	games, err := repo.ReadGames(ctx, 10, 0)
	require.NoError(t, err)

	assert.Equal(t, 1, len(games))
}
