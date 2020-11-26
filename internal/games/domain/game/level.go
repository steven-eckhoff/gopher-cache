package game

// Level holds all information for a level in a game.
type Level struct {
	title       string
	description string
	clues       []string
	answers     []string
}

func (l *Level) Title() string       { return l.title }
func (l *Level) Description() string { return l.description }
func (l *Level) Clues() []string     { return l.clues }
func (l *Level) Answers() []string   { return l.answers }

func (l *Level) isAnswer(input string) bool {
	for _, ans := range l.answers {
		if ans == input {
			return true
		}
	}

	return false
}

// UnmarshalLevelFromDatabase should only be used in repo implementations to unmarshal data from a database
// into a domain game level.
func UnmarshalLevelFromDatabase(title, description string, clues, answers []string) *Level {
	return &Level{
		title:       title,
		description: description,
		clues:       clues,
		answers:     answers,
	}
}
