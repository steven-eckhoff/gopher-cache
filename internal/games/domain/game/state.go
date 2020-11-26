package game

import (
	"errors"
	"github.com/google/uuid"
)

// State holds all the information for the state of a game.
type State struct {
	uuid            string
	playerUUID      string
	gameUUID        string
	gameLevels      int
	level           int
	clue            int
	completed       bool
	currentResponse Response
}

func (s State) UUID() string              { return s.uuid }
func (s State) PlayerUUID() string        { return s.playerUUID }
func (s State) GameUUID() string          { return s.gameUUID }
func (s State) GameLevels() int           { return s.gameLevels }
func (s State) Level() int                { return s.level }
func (s State) Clue() int                 { return s.clue }
func (s State) Completed() bool           { return s.completed }
func (s State) CurrentResponse() Response { return s.currentResponse }

// Update updates the state and player based on the current state of the game and the input from the player.
func (s *State) Update(g *Game, input string, p *Player) (*Response, error) {
	if s.gameUUID != g.UUID() {
		return nil, errors.New("invalid game")
	}

	// Check if game has been finished already.
	if s.completed {
		resp := newGameEndResponse(g.ending)
		return resp, nil
	}

	if s.level >= len(g.levels) {
		return nil, errors.New("invalid game state")
	}

	l := g.levels[s.level]

	if l.isAnswer(input) { // Is the input an answer to this level?
		s.level++
		s.clue = -1
		if s.level == len(g.levels) { // Have all levels been completed?
			s.completed = true
			resp := newGameEndResponse(g.ending)
			s.currentResponse = *resp
			if err := p.finishGame(g, s); err != nil {
				return nil, err
			}
			return resp, nil
		} else {
			l := g.levels[s.level]
			resp := newLevelResponse(l.title, l.description)
			s.currentResponse = *resp
			return resp, nil
		}
	} else {
		if len(l.clues) > 0 { // Does this level have any clues?
			if s.clue < len(l.clues)-1 {
				s.clue++
			}

			resp := newClueResponse(l.clues[s.clue])
			s.currentResponse = *resp
			return resp, nil
		} else {
			resp := newClueResponse("this level has no clues")
			s.currentResponse = *resp
			return resp, nil
		}
	}
}

// Start starts a game. It will update the player and return a new State.
func Start(g *Game, p *Player) (*State, *Response, error) {
	if g.uuid == "" {
		return nil, nil, errors.New("invalid game")
	}

	if p == nil {
		return nil, nil, errors.New("nil player")
	}

	if p.uuid == "" {
		return nil, nil, errors.New("invalid player")
	}

	id, err := uuid.NewRandom()
	if err != nil {
		return nil, nil, err
	}

	p.gamesStarted++
	p.currentGameStateUUID = id.String()

	resp := newLevelResponse(g.levels[0].title, g.levels[0].description)

	return &State{
		uuid:            id.String(),
		playerUUID:      p.uuid,
		gameUUID:        g.uuid,
		gameLevels:      len(g.levels),
		level:           0,
		clue:            -1,
		currentResponse: *resp,
	}, resp, nil
}

// UnmarshalGameStateFromDatabase should only be used in repo implementations to unmarshal data from a database
// into a domain game state.
func UnmarshalGameStateFromDatabase(
	uuid,
	playerUUID,
	gameUUID string,
	gameLevels,
	level,
	clue int,
	completed bool,
	currentResponse Response) *State {
	return &State{
		uuid:            uuid,
		playerUUID:      playerUUID,
		gameUUID:        gameUUID,
		gameLevels:      gameLevels,
		level:           level,
		clue:            clue,
		completed:       completed,
		currentResponse: currentResponse,
	}
}
