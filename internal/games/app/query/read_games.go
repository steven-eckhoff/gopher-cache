// Package query contains all the queries for games, players, and game states.
package query

import (
	"context"
)

// ReadGamesHandler handles the reading of games.
type ReadGamesHandler struct {
	readModel GamesReadModel
}

// NewReadGamesHandler creates a new handler.
func NewReadGamesHandler(readModel GamesReadModel) ReadGamesHandler {
	if readModel == nil {
		panic("nil readModel")
	}

	return ReadGamesHandler{readModel: readModel}
}

// GamesReadModel is the interface used for reading Games for a client query.
type GamesReadModel interface {
	// ReadGames will return an empty non-nil slice if no games are found.
	ReadGames(ctx context.Context, limit, offset int, options ...GameOption) ([]*Game, error)
}

// Handle is the use case for reading games.
func (h ReadGamesHandler) Handle(ctx context.Context, limit, offset int, options ...GameOption) ([]*Game, error) {
	return h.readModel.ReadGames(ctx, limit, offset, options...)
}

// GameOption is used for matching games in a query.
type GameOption struct {
	Key   string
	Op    string
	Value interface{}
}
