// Package app provides the games application.
package app

import (
	"gopher-cache/internal/games/app/command"
	"gopher-cache/internal/games/app/query"
)

// Application is the games application.
type Application struct {
	Commands Commands
	Queries  Queries
}

// Commands for the games application.
type Commands struct {
	CreateGame      command.CreateGameHandler
	CreateGameState command.CreateGameStateHandler
	UpdateGameState command.UpdateGameStateHandler
}

// Queries for the games application.
type Queries struct {
	GetGames  query.ReadGamesHandler
	GetPlayer query.ReadPlayerHandler
	GetState  query.ReadStateHandler
}
