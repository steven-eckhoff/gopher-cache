package command

import (
	"context"
	"gopher-cache/internal/common/logs"
	"gopher-cache/internal/games/domain/game"
)

// UpdateGameState represents the command input for updating a game state.
// All fields are required unless specified otherwise.
type UpdateGameState struct {
	PlayerNumber string `json:"-"`
	Input        string `json:"input"`
}

// UpdateGameStateHandler handles updating the game state.
type UpdateGameStateHandler struct {
	repo game.Repository
}

// NewUpdateGameStateHandler creates a new handler.
func NewUpdateGameStateHandler(repo game.Repository) UpdateGameStateHandler {
	if repo == nil {
		panic("nil repo")
	}

	return UpdateGameStateHandler{repo: repo}
}

// Handle handles the use case of updating an existing game state.
func (h UpdateGameStateHandler) Handle(ctx context.Context, cmd UpdateGameState) (resp *game.Response, err error) {
	defer func() {
		logs.LogCommandExecution("UpdateGameState", cmd, err)
	}()

	p, err := h.repo.GetPlayerByNumber(ctx, cmd.PlayerNumber)
	if err != nil {
		return nil, err
	}

	s, err := h.repo.GetState(ctx, p.CurrentGameStateUUID())
	if err != nil {
		return nil, err
	}

	g, err := h.repo.GetGame(ctx, s.GameUUID())
	if err != nil {
		return nil, err
	}

	resp, err = s.Update(g, cmd.Input, p)
	if err != nil {
		return nil, err
	}

	err = h.repo.UpdateStateAndPlayer(ctx, s, p)
	if err != nil {
		return nil, err
	}

	return resp, err
}
