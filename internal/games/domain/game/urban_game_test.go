package game

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopher-cache/internal/common/rand"
	"testing"
)

func TestNewUrbanGame(t *testing.T) {
	creator := newTestUser()

	title := "game title"
	description := "game description"
	ending := "game ending"

	kind := "urban"
	city := "austin"
	state := "texas"
	country := "united states"

	levelOneTitle := "level one title"
	levelOneDescription := "level one description"
	levelOneClues := []string{"level one clue one", "level one clue two", "level one clue three"}
	levelOneAnswers := []string{"level one answer one", "level one answer two", "level one answer three"}

	levelTwoTitle := "level two title"
	levelTwoDescription := "level two description"
	levelTwoClues := []string{"level two clue one", "level two clue two", "level two clue three"}
	levelTwoAnswers := []string{"level two answer one", "level two answer two", "level two answer three"}

	levelThreeTitle := "level three title"
	levelThreeDescription := "level three description"
	levelThreeClues := []string{"level three clue one", "level three clue two", "level three clue three"}
	levelThreeAnswers := []string{"level three answer one", "level three answer two", "level three answer three"}

	t.Run("valid", func(t *testing.T) {
		g, err := NewUrbanGame(creator, title, description, ending, city, state, country,
			NewLevelAdder(levelOneTitle, levelOneDescription, levelOneClues, levelOneAnswers),
			NewLevelAdder(levelTwoTitle, levelTwoDescription, levelTwoClues, levelTwoAnswers),
			NewLevelAdder(levelThreeTitle, levelThreeDescription, levelThreeClues, levelThreeAnswers),
		)
		require.NoError(t, err)

		assert.Equal(t, creator.UUID(), g.creatorUUID)
		assert.Equal(t, title, g.Title())
		assert.Equal(t, description, g.Description())
		assert.Equal(t, ending, g.Ending())
		assert.Equal(t, kind, g.Kind())
		assert.Equal(t, city, g.City())
		assert.Equal(t, state, g.State())
		assert.Equal(t, country, g.Country())

		levels := g.Levels()
		require.Equal(t, 3, len(levels))

		l := levels[0]
		assert.Equal(t, levelOneTitle, l.Title())
		assert.Equal(t, levelOneDescription, l.Description())
		assert.Equal(t, levelOneClues, l.Clues())
		assert.Equal(t, levelOneAnswers, l.Answers())

		l = levels[1]
		assert.Equal(t, levelTwoTitle, l.Title())
		assert.Equal(t, levelTwoDescription, l.Description())
		assert.Equal(t, levelTwoClues, l.Clues())
		assert.Equal(t, levelTwoAnswers, l.Answers())

		l = levels[2]
		assert.Equal(t, levelThreeTitle, l.Title())
		assert.Equal(t, levelThreeDescription, l.Description())
		assert.Equal(t, levelThreeClues, l.Clues())
		assert.Equal(t, levelThreeAnswers, l.Answers())
	})

	t.Run("invalid creator", func(t *testing.T) {
		_, err := NewUrbanGame(User{}, title, description, ending, city, state, country,
			NewLevelAdder(levelOneTitle, levelOneDescription, levelOneClues, levelOneAnswers),
		)
		assert.NotNil(t, err)
	})

	t.Run("invalid title", func(t *testing.T) {
		_, err := NewUrbanGame(creator, "", description, ending, city, state, country,
			NewLevelAdder(levelOneTitle, levelOneDescription, levelOneClues, levelOneAnswers),
		)
		assert.NotNil(t, err)

		_, err = NewUrbanGame(creator, rand.String(MaxTitleLength+1, 1), description, ending, city, state, country,
			NewLevelAdder(levelOneTitle, levelOneDescription, levelOneClues, levelOneAnswers),
		)
		assert.NotNil(t, err)
	})

	t.Run("invalid description", func(t *testing.T) {
		_, err := NewUrbanGame(creator, title, "", ending, city, state, country,
			NewLevelAdder(levelOneTitle, levelOneDescription, levelOneClues, levelOneAnswers),
		)
		assert.NotNil(t, err)

		_, err = NewUrbanGame(creator, title, rand.String(MaxDescriptionLength+1, 1), ending, city, state, country,
			NewLevelAdder(levelOneTitle, levelOneDescription, levelOneClues, levelOneAnswers),
		)
		assert.NotNil(t, err)
	})

	t.Run("invalid ending", func(t *testing.T) {
		_, err := NewUrbanGame(creator, title, description, "", city, state, country,
			NewLevelAdder(levelOneTitle, levelOneDescription, levelOneClues, levelOneAnswers),
		)
		assert.NotNil(t, err)

		_, err = NewUrbanGame(creator, title, description, rand.String(MaxEndingLength+1, 1), city, state, country,
			NewLevelAdder(levelOneTitle, levelOneDescription, levelOneClues, levelOneAnswers),
		)
		assert.NotNil(t, err)
	})

	t.Run("missing city", func(t *testing.T) {
		_, err := NewUrbanGame(creator, title, description, ending, "", state, country,
			NewLevelAdder(levelOneTitle, levelOneDescription, levelOneClues, levelOneAnswers),
		)
		assert.NotNil(t, err)
	})

	t.Run("missing state", func(t *testing.T) {
		_, err := NewUrbanGame(creator, title, description, ending, city, "", country,
			NewLevelAdder(levelOneTitle, levelOneDescription, levelOneClues, levelOneAnswers),
		)
		assert.NotNil(t, err)
	})

	t.Run("missing country", func(t *testing.T) {
		_, err := NewUrbanGame(creator, title, description, ending, city, state, "",
			NewLevelAdder(levelOneTitle, levelOneDescription, levelOneClues, levelOneAnswers),
		)
		assert.NotNil(t, err)
	})

	t.Run("invalid level title", func(t *testing.T) {
		_, err := NewUrbanGame(creator, title, description, ending, city, state, country,
			NewLevelAdder("", levelOneDescription, levelOneClues, levelOneAnswers),
		)
		assert.NotNil(t, err)

		_, err = NewUrbanGame(creator, title, description, ending, city, state, country,
			NewLevelAdder(rand.String(MaxTitleLength+1, 1), levelOneDescription, levelOneClues, levelOneAnswers),
		)
		assert.NotNil(t, err)
	})

	t.Run("invalid level description", func(t *testing.T) {
		_, err := NewUrbanGame(creator, title, description, ending, city, state, country,
			NewLevelAdder(title, "", levelOneClues, levelOneAnswers),
		)
		assert.NotNil(t, err)

		_, err = NewUrbanGame(creator, title, description, ending, city, state, country,
			NewLevelAdder(levelOneTitle, rand.String(MaxDescriptionLength+1, 1), levelOneClues, levelOneAnswers),
		)
		assert.NotNil(t, err)
	})

	t.Run("invalid level clue", func(t *testing.T) {
		_, err := NewUrbanGame(creator, title, description, ending, city, state, country,
			NewLevelAdder(levelOneTitle, levelOneDescription, []string{""}, levelOneAnswers),
		)
		assert.NotNil(t, err)

		_, err = NewUrbanGame(creator, title, description, ending, city, state, "",
			NewLevelAdder(levelOneTitle, levelOneDescription, []string{rand.String(MaxClueLength+1, 1)}, levelOneAnswers),
		)
		assert.NotNil(t, err)
	})

	t.Run("invalid level answer", func(t *testing.T) {
		_, err := NewUrbanGame(creator, title, description, ending, city, state, country,
			NewLevelAdder(levelOneTitle, levelOneDescription, levelOneClues, []string{""}),
		)
		assert.NotNil(t, err)

		_, err = NewUrbanGame(creator, title, description, ending, city, state, "",
			NewLevelAdder(levelOneTitle, levelOneDescription, levelOneClues, []string{rand.String(MaxClueLength+1, 1)}),
		)
		assert.NotNil(t, err)
	})

	t.Run("missing answers", func(t *testing.T) {
		_, err := NewUrbanGame(creator, title, description, ending, city, state, country,
			NewLevelAdder(levelOneTitle, levelOneDescription, levelOneClues, []string{}),
		)
		assert.NotNil(t, err)
	})
}

func newValidTestUrbanGame() *Game {
	creator := newTestUser()

	title := "game title"
	description := "game description"
	ending := "game ending"

	city := "austin"
	state := "texas"
	country := "united states"

	levelOneTitle := "level one title"
	levelOneDescription := "level one description"
	levelOneClues := []string{"level one clue one", "level one clue two", "level one clue three"}
	levelOneAnswers := []string{"level one answer one", "level one answer two", "level one answer three"}

	levelTwoTitle := "level two title"
	levelTwoDescription := "level two description"
	levelTwoClues := []string{"level two clue one", "level two clue two", "level two clue three"}
	levelTwoAnswers := []string{"level two answer one", "level two answer two", "level two answer three"}

	levelThreeTitle := "level three title"
	levelThreeDescription := "level three description"
	levelThreeClues := []string{"level three clue one", "level three clue two", "level three clue three"}
	levelThreeAnswers := []string{"level three answer one", "level three answer two", "level three answer three"}

	g, err := NewUrbanGame(creator, title, description, ending, city, state, country,
		NewLevelAdder(levelOneTitle, levelOneDescription, levelOneClues, levelOneAnswers),
		NewLevelAdder(levelTwoTitle, levelTwoDescription, levelTwoClues, levelTwoAnswers),
		NewLevelAdder(levelThreeTitle, levelThreeDescription, levelThreeClues, levelThreeAnswers),
	)

	if err != nil {
		panic(err)
	}

	return g
}
