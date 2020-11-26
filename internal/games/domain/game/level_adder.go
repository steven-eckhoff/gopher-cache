package game

import "errors"

// LevelAdder checks a level for errors and if it is error free adds it to a game.
type LevelAdder func(g *Game) error

// NewLevelAdder creates a new LevelAdder.
func NewLevelAdder(title, description string, clues, answers []string) LevelAdder {
	return func(g *Game) error {
		if title == "" {
			return errors.New("level has not title")
		}

		if len(title) > MaxTitleLength {
			return errors.New("level title length greater than 64")
		}

		if description == "" {
			return errors.New("level has no description")
		}

		if len(description) > MaxDescriptionLength {
			return errors.New("level description length greater than 200")
		}

		l := Level{
			title:       title,
			description: description,
		}

		for _, clue := range clues {
			if clue == "" {
				return errors.New("clue is empty")
			}

			if len(clue) > MaxClueLength {
				return errors.New("clue length greater than 64")
			}

			l.clues = append(l.clues, clue)
		}

		if len(answers) == 0 {
			return errors.New("level has no answers")
		}

		for _, answer := range answers {
			if answer == "" {
				return errors.New("answer is empty")
			}

			if len(answer) > MaxAnswerLength {
				return errors.New("answer length greater than 64")
			}

			l.answers = append(l.answers, answer)
		}

		g.levels = append(g.levels, &l)

		return nil
	}
}
