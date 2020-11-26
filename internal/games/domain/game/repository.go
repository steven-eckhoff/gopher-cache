package game

import (
	"context"
	"errors"
)

var (
	ErrorPlayerNotFound = errors.New("player not found")
)

// Repository is the interface used to persist domain types.
type Repository interface {
	AddGame(ctx context.Context, game *Game) error
	GetGame(ctx context.Context, uuid string) (*Game, error)

	AddPlayer(ctx context.Context, player *Player) error
	// GetPlayer returns ErrorPlayerNotFound if player does not exist.
	GetPlayer(ctx context.Context, uuid string) (*Player, error)
	// GetPlayerByNumber returns ErrorPlayerNotFound if player with number does not exist.
	GetPlayerByNumber(ctx context.Context, playerNumber string) (*Player, error)

	AddState(ctx context.Context, state *State) error
	GetState(ctx context.Context, uuid string) (*State, error)
	UpdateState(ctx context.Context, state *State) error
	AddStateAndUpdatePlayer(ctx context.Context, state *State, player *Player) error
	UpdateStateAndPlayer(ctx context.Context, state *State, player *Player) error
}
