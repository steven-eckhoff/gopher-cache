// Package game provides types and methods for the game domain.
package game

import (
	"errors"
	"github.com/google/uuid"
)

// These are the limits imposed on games.
const (
	MaxTitleLength       = 64
	MaxDescriptionLength = 200
	MaxClueLength        = 64
	MaxAnswerLength      = 64
	MaxEndingLength      = 200
)

// Game holds all information about a game.
type Game struct {
	uuid        string
	creatorUUID string
	title       string
	description string
	levels      []*Level
	ending      string
	kind        string
	city        string
	state       string
	country     string
	value       int
}

func (g *Game) UUID() string        { return g.uuid }
func (g *Game) CreatorUUID() string { return g.creatorUUID }
func (g *Game) Title() string       { return g.title }
func (g *Game) Description() string { return g.description }
func (g *Game) Levels() []*Level    { return g.levels }
func (g *Game) Ending() string      { return g.ending }
func (g *Game) Kind() string        { return g.kind }
func (g *Game) City() string        { return g.city }
func (g *Game) State() string       { return g.state }
func (g *Game) Country() string     { return g.country }
func (g *Game) Value() int          { return g.value }

// newGame creates a new game for public constructors.
func newGame(creator User, title, description, ending string, kind string, levelAdders ...LevelAdder) (*Game, error) {
	if creator.UUID() == "" {
		return nil, errors.New("invalid creator")
	}

	if title == "" {
		return nil, errors.New("game has no title")
	}

	if len(title) > MaxTitleLength {
		return nil, errors.New("game title greater than 64")
	}

	if description == "" {
		return nil, errors.New("game has no description")
	}

	if len(description) > MaxDescriptionLength {
		return nil, errors.New("game description greater than 200")
	}

	if ending == "" {
		return nil, errors.New("game has no ending")
	}

	switch kind {
	case "urban":
	default:
		return nil, errors.New("unrecognized kind")
	}

	if len(ending) > MaxEndingLength {
		return nil, errors.New("ending length greater than 200")
	}

	if len(levelAdders) == 0 {
		return nil, errors.New("game has no levels")
	}

	id, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	g := &Game{
		uuid:        id.String(),
		creatorUUID: creator.UUID(),
		title:       title,
		description: description,
		ending:      ending,
		kind:        kind,
		value:       42,
	}

	for _, addLevel := range levelAdders {
		if err := addLevel(g); err != nil {
			return nil, err
		}
	}

	return g, nil
}

// UnmarshalFromDatabase should only be used in repo implementations to unmarshal data from a database
// into a domain game.
func UnmarshalFromDataBase(
	uuid,
	creatorUUID,
	title,
	description string,
	levels []*Level,
	ending,
	kind,
	city,
	state,
	country string,
	value int) (*Game, error) {
	return &Game{
		uuid:        uuid,
		creatorUUID: creatorUUID,
		title:       title,
		description: description,
		levels:      levels,
		ending:      ending,
		kind:        kind,
		city:        city,
		state:       state,
		country:     country,
		value:       value,
	}, nil
}
