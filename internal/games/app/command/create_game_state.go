package command

import (
	"context"
	"errors"
	"gopher-cache/internal/common/logs"
	"gopher-cache/internal/games/domain/game"
)

// CreateGameState represents the command input for creating a game state.
// All fields are required unless specified otherwise.
type CreateGameState struct {
	User     game.User `json:"-"`
	GameUUID string    `json:"gameUUID"`
}

// CreateGameStateHandler handles creating the game state.
type CreateGameStateHandler struct {
	repo game.Repository
}

// NewCreateGameStateHandler creates a new handler.
func NewCreateGameStateHandler(repo game.Repository) CreateGameStateHandler {
	if repo == nil {
		panic("nil repo")
	}

	return CreateGameStateHandler{repo: repo}
}

// Handle handles the use case of creating a new game state.
func (h CreateGameStateHandler) Handle(ctx context.Context, cmd CreateGameState) (resp *game.Response, err error) {
	defer func() {
		logs.LogCommandExecution("CreateGameState", cmd, err)
	}()

	g, err := h.repo.GetGame(ctx, cmd.GameUUID)
	if err != nil {
		return nil, err
	}

	p, err := h.repo.GetPlayer(ctx, cmd.User.UUID())
	if err != nil {
		if errors.Is(err, game.ErrorPlayerNotFound) {
			p, err = game.NewPlayerFromUser(cmd.User)
			if err != nil {
				return nil, err
			}

			err = h.repo.AddPlayer(ctx, p)
			if err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}

	state, resp, err := game.Start(g, p)
	if err != nil {
		return nil, err
	}

	err = h.repo.AddStateAndUpdatePlayer(ctx, state, p)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
