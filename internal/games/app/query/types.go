package query

import "gopher-cache/internal/games/domain/game"

// Game represents how Game queries will be presented to clients.
type Game struct {
	UUID        string `json:"uuid"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

// Player represents how Player queries will be presented to clients.
type Player struct {
	GamesFinished int `json:"gamesFinished"`
	TotalPoints   int `json:"totalPoints"`
}

// State represents how State queries will be presented to clients.
type State struct {
	CurrentResponse game.Response `json:"currentResponse"`
}
