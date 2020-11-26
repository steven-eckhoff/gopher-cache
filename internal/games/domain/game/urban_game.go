package game

import "errors"

// NewUrbanGame creates a new Game with urban game info.
func NewUrbanGame(creator User, title, description, ending, city, state, country string, levelAdders ...LevelAdder) (*Game, error) {
	g, err := newGame(creator, title, description, ending, "urban", levelAdders...)
	if err != nil {
		return nil, err
	}

	if city == "" {
		return nil, errors.New("missing city")
	}

	if state == "" {
		return nil, errors.New("missing state")
	}

	if country == "" {
		return nil, errors.New("missing country")
	}

	g.city = city
	g.state = state
	g.country = country

	return g, nil
}
