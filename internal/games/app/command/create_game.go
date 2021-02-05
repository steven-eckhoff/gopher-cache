// Package command contains all the commands or use cases for the games, players and game states.
package command

import (
	"context"
	"errors"
	"gopher-cache/internal/common/logs"
	"gopher-cache/internal/games/domain/game"
)

// CreateGame represents the command input for creating a game.
// All fields are required unless specified otherwise.
type CreateGame struct {
	Creator     game.User   `json:"-"`
	Title       string      `json:"title"`
	Description string      `json:"description"`
	Levels      []GameLevel `json:"levels"`
	Ending      string      `json:"ending"`
	// Kind can only be set to urban for now.
	Kind string `json:"kind"`
	// Required if the kind is urban.
	City string `json:"city"`
	// Required if the kind is urban
	State string `json:"state"`
	// Required if the kind is urban
	Country string `json:"country"`
}

type GameLevel struct {
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Clues       []string `json:"clues"`
	Answers     []string `json:"answers"`
}

// CreateGameHandler handles creating games.
type CreateGameHandler struct {
	repo game.Repository
}

// NewCreateGameHandler creates a new game handler.
func NewCreateGameHandler(repo game.Repository) CreateGameHandler {
	if repo == nil {
		panic("nil repo")
	}

	return CreateGameHandler{repo: repo}
}

// Handle handles the use case of creating games.
func (h CreateGameHandler) Handle(ctx context.Context, cmd CreateGame) (err error) {
	defer func() {
		logs.LogCommandExecution("CreateGame", cmd, err)
	}()

	var levelAdders []game.LevelAdder
	for _, l := range cmd.Levels {
		levelAdders = append(levelAdders, game.NewLevelAdder(l.Title, l.Description, l.Clues, l.Answers))
	}

	switch cmd.Kind {
	case "urban":
		g, err := game.NewUrbanGame(
			cmd.Creator,
			cmd.Title,
			cmd.Description,
			cmd.Ending,
			cmd.City,
			cmd.State,
			cmd.Country,
			levelAdders...)
		if err != nil {
			return err
		}

		return h.repo.AddGame(ctx, g)
	default:
		return errors.New("unknown game kind")
	}
}
