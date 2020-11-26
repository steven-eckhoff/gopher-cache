package game

import (
	"errors"
)

// Player holds all information about a player.
type Player struct {
	uuid                 string
	number               string
	gamesStarted         int
	gamesFinished        int
	totalPoints          int
	currentGameStateUUID string
}

func (p *Player) UUID() string                 { return p.uuid }
func (p *Player) Number() string               { return p.number }
func (p *Player) GamesStarted() int            { return p.gamesStarted }
func (p *Player) GamesFinished() int           { return p.gamesFinished }
func (p *Player) TotalPoints() int             { return p.totalPoints }
func (p *Player) CurrentGameStateUUID() string { return p.currentGameStateUUID }

// NewPlayer creates a new Player from a User.
func NewPlayerFromUser(u User) (*Player, error) {
	if u.UUID() == "" {
		return &Player{}, errors.New("invalid user")
	}

	return &Player{
		uuid:   u.UUID(),
		number: u.Number(),
	}, nil
}

func (p *Player) finishGame(g *Game, s *State) error {
	if s.completed == false {
		return errors.New("game is not completed")
	}

	if s.gameUUID != g.uuid {
		return errors.New("invalid game")
	}

	if p.uuid != s.playerUUID {
		return errors.New("invalid state")
	}

	p.gamesFinished++
	p.totalPoints += g.value

	return nil
}

// UnmarshalPlayerFromDatabase should only be used in repo implementations to unmarshal data from a database
// into a domain game player.
func UnmarshalPlayerFromDatabase(
	uuid,
	number string,
	gamesStarted,
	gamesFinished,
	totalPoints int,
	currentGameStateUUID string) *Player {
	return &Player{
		uuid:                 uuid,
		number:               number,
		gamesStarted:         gamesStarted,
		gamesFinished:        gamesFinished,
		totalPoints:          totalPoints,
		currentGameStateUUID: currentGameStateUUID,
	}
}
