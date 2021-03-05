package adapters

import (
	"context"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopher-cache/internal/common/emulators"
	"gopher-cache/internal/games/app/query"
	"gopher-cache/internal/games/domain/game"
	"testing"
)

func TestFirestoreGameRepository_AddGame(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	ctx := context.Background()

	client, cleanup := emulators.NewFirestoreClient(ctx)
	defer func() {
		_ = client.Close()
		cleanup()
	}()

	repo, err := NewFirestoreGameRepository(client)
	require.NoError(t, err)

	userID, err := uuid.NewRandom()
	require.NoError(t, err)

	userNumber := "15734497033"

	u, err := game.NewUser(userID.String(), userNumber)
	require.NoError(t, err)

	expectedGame, err := game.NewUrbanGame(
		u,
		"An Awesome Game",
		"This is an awesome game",
		"The end!",
		"Austin",
		"Texas",
		"USA",
		game.NewLevelAdder(
			"Level One",
			"This is Level One",
			[]string{"Who is the best?", "Level One is the best", "Say I am the best"},
			[]string{"Level One is the best"}),
		game.NewLevelAdder(
			"Level Two",
			"This is Level Two",
			[]string{"Who is the best?", "Level Two is the best", "Say I am the best"},
			[]string{"Level Two is the best"}),
		game.NewLevelAdder(
			"Level Three",
			"This is Level Three",
			[]string{"Who is the best?", "Level Three is the best", "Say I am the best"},
			[]string{"Level Three is the best"}),
	)
	require.NoError(t, err)

	err = repo.AddGame(ctx, expectedGame)
	assert.NoError(t, err)

	gotGame, err := repo.GetGame(ctx, expectedGame.UUID())
	assert.NoError(t, err)

	assert.Equal(t, expectedGame, gotGame)
}

func TestFirestoreGameRepository_AddPlayer(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	ctx := context.Background()

	client, cleanup := emulators.NewFirestoreClient(ctx)
	defer func() {
		_ = client.Close()
		cleanup()
	}()

	repo, err := NewFirestoreGameRepository(client)
	require.NoError(t, err)

	userID, err := uuid.NewRandom()
	require.NoError(t, err)

	userNumber := "15734497033"

	u, err := game.NewUser(userID.String(), userNumber)
	require.NoError(t, err)

	expectedPlayer, err := game.NewPlayerFromUser(u)
	require.NoError(t, err)

	// Add player
	err = repo.AddPlayer(ctx, expectedPlayer)
	assert.NoError(t, err)

	// Get player
	gotPlayer, err := repo.GetPlayer(ctx, expectedPlayer.UUID())
	assert.NoError(t, err)

	assert.Equal(t, expectedPlayer, gotPlayer)

	// Get player by number
	gotPlayer, err = repo.GetPlayerByNumber(ctx, expectedPlayer.Number())
	assert.NoError(t, err)

	assert.Equal(t, expectedPlayer, gotPlayer)
}

func TestFirestoreGameRepository_AddState(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	ctx := context.Background()

	client, cleanup := emulators.NewFirestoreClient(ctx)
	defer func() {
		_ = client.Close()
		cleanup()
	}()

	repo, err := NewFirestoreGameRepository(client)
	require.NoError(t, err)

	userID, err := uuid.NewRandom()
	require.NoError(t, err)

	userNumber := "15734497033"

	u, err := game.NewUser(userID.String(), userNumber)
	require.NoError(t, err)

	p, err := game.NewPlayerFromUser(u)
	require.NoError(t, err)

	err = repo.AddPlayer(ctx, p)
	require.NoError(t, err)

	g, err := game.NewUrbanGame(
		u,
		"An Awesome Game",
		"This is an awesome game",
		"The end!",
		"Austin",
		"Texas",
		"USA",
		game.NewLevelAdder(
			"Level One",
			"This is Level One",
			[]string{"Who is the best?", "Level One is the best", "Say I am the best"},
			[]string{"Level One is the best"}),
		game.NewLevelAdder(
			"Level Two",
			"This is Level Two",
			[]string{"Who is the best?", "Level Two is the best", "Say I am the best"},
			[]string{"Level Two is the best"}),
		game.NewLevelAdder(
			"Level Three",
			"This is Level Three",
			[]string{"Who is the best?", "Level Three is the best", "Say I am the best"},
			[]string{"Level Three is the best"}),
	)
	require.NoError(t, err)

	expectedState, _, err := game.Start(g, p)
	require.NoError(t, err)

	// Add state
	err = repo.AddState(ctx, expectedState)
	require.NoError(t, err)

	// Get state
	gotState, err := repo.GetState(ctx, expectedState.UUID())
	require.NoError(t, err)
	assert.Equal(t, expectedState, gotState)
}

func TestFirestoreGameRepository_AddStateAndUpdatePlayer(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	ctx := context.Background()

	client, cleanup := emulators.NewFirestoreClient(ctx)
	defer func() {
		_ = client.Close()
		cleanup()
	}()

	repo, err := NewFirestoreGameRepository(client)
	require.NoError(t, err)

	userID, err := uuid.NewRandom()
	require.NoError(t, err)

	userNumber := "15734497033"

	u, err := game.NewUser(userID.String(), userNumber)
	require.NoError(t, err)

	expectedPlayer, err := game.NewPlayerFromUser(u)
	require.NoError(t, err)

	err = repo.AddPlayer(ctx, expectedPlayer)
	require.NoError(t, err)

	g, err := game.NewUrbanGame(
		u,
		"An Awesome Game",
		"This is an awesome game",
		"The end!",
		"Austin",
		"Texas",
		"USA",
		game.NewLevelAdder(
			"Level One",
			"This is Level One",
			[]string{"Who is the best?", "Level One is the best", "Say I am the best"},
			[]string{"Level One is the best"}),
		game.NewLevelAdder(
			"Level Two",
			"This is Level Two",
			[]string{"Who is the best?", "Level Two is the best", "Say I am the best"},
			[]string{"Level Two is the best"}),
		game.NewLevelAdder(
			"Level Three",
			"This is Level Three",
			[]string{"Who is the best?", "Level Three is the best", "Say I am the best"},
			[]string{"Level Three is the best"}),
	)
	require.NoError(t, err)

	expectedState, _, err := game.Start(g, expectedPlayer)
	require.NoError(t, err)

	// Add state and update player
	err = repo.AddStateAndUpdatePlayer(ctx, expectedState, expectedPlayer)
	require.NoError(t, err)

	// Get state
	gotState, err := repo.GetState(ctx, expectedState.UUID())
	require.NoError(t, err)
	assert.Equal(t, expectedState, gotState)

	// Get player
	gotPlayer, err := repo.GetPlayer(ctx, expectedPlayer.UUID())
	require.NoError(t, err)
	assert.Equal(t, expectedPlayer, gotPlayer)
}

func TestFirestoreGameRepository_UpdateState(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	ctx := context.Background()

	client, cleanup := emulators.NewFirestoreClient(ctx)
	defer func() {
		_ = client.Close()
		cleanup()
	}()

	repo, err := NewFirestoreGameRepository(client)
	require.NoError(t, err)

	userID, err := uuid.NewRandom()
	require.NoError(t, err)

	userNumber := "15734497033"

	u, err := game.NewUser(userID.String(), userNumber)
	require.NoError(t, err)

	p, err := game.NewPlayerFromUser(u)
	require.NoError(t, err)

	err = repo.AddPlayer(ctx, p)
	require.NoError(t, err)

	g, err := game.NewUrbanGame(
		u,
		"An Awesome Game",
		"This is an awesome game",
		"The end!",
		"Austin",
		"Texas",
		"USA",
		game.NewLevelAdder(
			"Level One",
			"This is Level One",
			[]string{"Who is the best?", "Level One is the best", "Say I am the best"},
			[]string{"Level One is the best"}),
		game.NewLevelAdder(
			"Level Two",
			"This is Level Two",
			[]string{"Who is the best?", "Level Two is the best", "Say I am the best"},
			[]string{"Level Two is the best"}),
		game.NewLevelAdder(
			"Level Three",
			"This is Level Three",
			[]string{"Who is the best?", "Level Three is the best", "Say I am the best"},
			[]string{"Level Three is the best"}),
	)
	require.NoError(t, err)

	expectedState, _, err := game.Start(g, p)
	require.NoError(t, err)

	// Add state and update player
	err = repo.AddStateAndUpdatePlayer(ctx, expectedState, p)
	require.NoError(t, err)

	resp, err := expectedState.Update(g, "Level One is the best", p)
	require.NoError(t, err)
	// Make sure the state has changed.
	require.Equal(t, game.LevelResponse, resp.Kind)

	err = repo.UpdateState(ctx, expectedState)

	gotState, err := repo.GetState(ctx, expectedState.UUID())
	require.NoError(t, err)
	assert.Equal(t, expectedState, gotState)
}

func TestFirestoreGameRepository_UpdateStateAndPlayer(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	ctx := context.Background()

	client, cleanup := emulators.NewFirestoreClient(ctx)
	defer func() {
		_ = client.Close()
		cleanup()
	}()

	repo, err := NewFirestoreGameRepository(client)
	require.NoError(t, err)

	userID, err := uuid.NewRandom()
	require.NoError(t, err)

	userNumber := "15734497033"

	u, err := game.NewUser(userID.String(), userNumber)
	require.NoError(t, err)

	expectedPlayer, err := game.NewPlayerFromUser(u)
	require.NoError(t, err)

	err = repo.AddPlayer(ctx, expectedPlayer)
	require.NoError(t, err)

	g, err := game.NewUrbanGame(
		u,
		"An Awesome Game",
		"This is an awesome game",
		"The end!",
		"Austin",
		"Texas",
		"USA",
		game.NewLevelAdder(
			"Level One",
			"This is Level One",
			[]string{"Who is the best?", "Level One is the best", "Say I am the best"},
			[]string{"Level One is the best"}),
		game.NewLevelAdder(
			"Level Two",
			"This is Level Two",
			[]string{"Who is the best?", "Level Two is the best", "Say I am the best"},
			[]string{"Level Two is the best"}),
		game.NewLevelAdder(
			"Level Three",
			"This is Level Three",
			[]string{"Who is the best?", "Level Three is the best", "Say I am the best"},
			[]string{"Level Three is the best"}),
	)
	require.NoError(t, err)

	expectedState, _, err := game.Start(g, expectedPlayer)
	require.NoError(t, err)

	// Add state and update player
	err = repo.AddStateAndUpdatePlayer(ctx, expectedState, expectedPlayer)
	require.NoError(t, err)

	resp, err := expectedState.Update(g, "Level One is the best", expectedPlayer)
	require.NoError(t, err)

	resp, err = expectedState.Update(g, "Level Two is the best", expectedPlayer)
	require.NoError(t, err)

	resp, err = expectedState.Update(g, "Level Three is the best", expectedPlayer)
	require.NoError(t, err)

	// Make sure the state has changed.
	require.Equal(t, game.EndResponse, resp.Kind)
	require.Equal(t, 1, expectedPlayer.GamesFinished())

	err = repo.UpdateStateAndPlayer(ctx, expectedState, expectedPlayer)

	gotState, err := repo.GetState(ctx, expectedState.UUID())
	require.NoError(t, err)
	assert.Equal(t, expectedState, gotState)

	gotPlayer, err := repo.GetPlayer(ctx, expectedPlayer.UUID())
	require.NoError(t, err)
	assert.Equal(t, expectedPlayer, gotPlayer)
}

func TestFirestoreGameRepository_ReadGames(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	ctx := context.Background()

	client, cleanup := emulators.NewFirestoreClient(ctx)
	defer func() {
		_ = client.Close()
		cleanup()
	}()

	repo, err := NewFirestoreGameRepository(client)
	require.NoError(t, err)

	userID, err := uuid.NewRandom()
	require.NoError(t, err)

	userNumber := "15734497033"

	u, err := game.NewUser(userID.String(), userNumber)
	require.NoError(t, err)

	austinGame, err := game.NewUrbanGame(
		u,
		"An Awesome Game",
		"This is an awesome game",
		"The end!",
		"Austin",
		"Texas",
		"USA",
		game.NewLevelAdder(
			"Level One",
			"This is Level One",
			[]string{"Who is the best?", "Level One is the best", "Say I am the best"},
			[]string{"Level One is the best"}),
		game.NewLevelAdder(
			"Level Two",
			"This is Level Two",
			[]string{"Who is the best?", "Level Two is the best", "Say I am the best"},
			[]string{"Level Two is the best"}),
		game.NewLevelAdder(
			"Level Three",
			"This is Level Three",
			[]string{"Who is the best?", "Level Three is the best", "Say I am the best"},
			[]string{"Level Three is the best"}),
	)
	require.NoError(t, err)

	err = repo.AddGame(ctx, austinGame)
	require.NoError(t, err)

	dallasGame, err := game.NewUrbanGame(
		u,
		"An Awesome Game",
		"This is an awesome game",
		"The end!",
		"Dallas",
		"Texas",
		"USA",
		game.NewLevelAdder(
			"Level One",
			"This is Level One",
			[]string{"Who is the best?", "Level One is the best", "Say I am the best"},
			[]string{"Level One is the best"}),
		game.NewLevelAdder(
			"Level Two",
			"This is Level Two",
			[]string{"Who is the best?", "Level Two is the best", "Say I am the best"},
			[]string{"Level Two is the best"}),
		game.NewLevelAdder(
			"Level Three",
			"This is Level Three",
			[]string{"Who is the best?", "Level Three is the best", "Say I am the best"},
			[]string{"Level Three is the best"}),
	)
	require.NoError(t, err)

	err = repo.AddGame(ctx, dallasGame)
	require.NoError(t, err)

	chicagoGame, err := game.NewUrbanGame(
		u,
		"An Awesome Game",
		"This is an awesome game",
		"The end!",
		"Chicago",
		"Illinois",
		"USA",
		game.NewLevelAdder(
			"Level One",
			"This is Level One",
			[]string{"Who is the best?", "Level One is the best", "Say I am the best"},
			[]string{"Level One is the best"}),
		game.NewLevelAdder(
			"Level Two",
			"This is Level Two",
			[]string{"Who is the best?", "Level Two is the best", "Say I am the best"},
			[]string{"Level Two is the best"}),
		game.NewLevelAdder(
			"Level Three",
			"This is Level Three",
			[]string{"Who is the best?", "Level Three is the best", "Say I am the best"},
			[]string{"Level Three is the best"}),
	)
	require.NoError(t, err)

	err = repo.AddGame(ctx, chicagoGame)
	require.NoError(t, err)

	games, err := repo.ReadGames(ctx, 10, 0)
	assert.NoError(t, err)
	assert.Equal(t, 3, len(games))

	games, err = repo.ReadGames(ctx, 2, 0)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(games))

	games, err = repo.ReadGames(ctx, 10, 0, query.GameOption{
		Key:   "country",
		Op:    "==",
		Value: "USA",
	})
	assert.NoError(t, err)
	assert.Equal(t, 3, len(games))

	games, err = repo.ReadGames(ctx, 10, 0, query.GameOption{
		Key:   "state",
		Op:    "==",
		Value: "Texas",
	})
	assert.NoError(t, err)
	assert.Equal(t, 2, len(games))

	games, err = repo.ReadGames(ctx, 10, 0, query.GameOption{
		Key:   "state",
		Op:    "==",
		Value: "Illinois",
	})
	assert.NoError(t, err)
	assert.Equal(t, 1, len(games))

	games, err = repo.ReadGames(ctx, 10, 0, query.GameOption{
		Key:   "city",
		Op:    "==",
		Value: "Austin",
	})
	assert.NoError(t, err)
	assert.Equal(t, 1, len(games))
}

func TestFirestoreGameRepository_ReadPlayer(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	ctx := context.Background()

	client, cleanup := emulators.NewFirestoreClient(ctx)
	defer func() {
		_ = client.Close()
		cleanup()
	}()

	repo, err := NewFirestoreGameRepository(client)
	require.NoError(t, err)

	userID, err := uuid.NewRandom()
	require.NoError(t, err)

	userNumber := "15734497033"

	u, err := game.NewUser(userID.String(), userNumber)
	require.NoError(t, err)

	commandPlayer, err := game.NewPlayerFromUser(u)
	require.NoError(t, err)

	err = repo.AddPlayer(ctx, commandPlayer)
	require.NoError(t, err)

	_, err = repo.ReadPlayer(ctx, commandPlayer.UUID())
	require.NoError(t, err)
}

func TestFirestoreGameRepository_ReadState(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	ctx := context.Background()

	client, cleanup := emulators.NewFirestoreClient(ctx)
	defer func() {
		_ = client.Close()
		cleanup()
	}()

	repo, err := NewFirestoreGameRepository(client)
	require.NoError(t, err)

	userID, err := uuid.NewRandom()
	require.NoError(t, err)

	userNumber := "15734497033"

	u, err := game.NewUser(userID.String(), userNumber)
	require.NoError(t, err)

	expectedPlayer, err := game.NewPlayerFromUser(u)
	require.NoError(t, err)

	err = repo.AddPlayer(ctx, expectedPlayer)
	require.NoError(t, err)

	g, err := game.NewUrbanGame(
		u,
		"An Awesome Game",
		"This is an awesome game",
		"The end!",
		"Austin",
		"Texas",
		"USA",
		game.NewLevelAdder(
			"Level One",
			"This is Level One",
			[]string{"Who is the best?", "Level One is the best", "Say I am the best"},
			[]string{"Level One is the best"}),
		game.NewLevelAdder(
			"Level Two",
			"This is Level Two",
			[]string{"Who is the best?", "Level Two is the best", "Say I am the best"},
			[]string{"Level Two is the best"}),
		game.NewLevelAdder(
			"Level Three",
			"This is Level Three",
			[]string{"Who is the best?", "Level Three is the best", "Say I am the best"},
			[]string{"Level Three is the best"}),
	)
	require.NoError(t, err)

	commandState, _, err := game.Start(g, expectedPlayer)
	require.NoError(t, err)

	err = repo.AddState(ctx, commandState)
	require.NoError(t, err)

	queryState, err := repo.ReadState(ctx, commandState.UUID())
	require.NoError(t, err)
	assert.Equal(t, commandState.CurrentResponse(), queryState.CurrentResponse)
}
