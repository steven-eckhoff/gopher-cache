package game

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestStartGame(t *testing.T) {
	g := newValidTestUrbanGame()
	p := newValidTestPlayer()

	t.Run("valid", func(t *testing.T) {
		resp := newLevelResponse(g.levels[0].title, g.levels[0].description)

		s, resp, err := Start(g, p)
		require.NoError(t, err)

		assert.Equal(t, p.UUID(), s.PlayerUUID())
		assert.Equal(t, g.UUID(), s.GameUUID())
		assert.Equal(t, -1, s.clue)
		assert.Equal(t, len(g.levels), s.GameLevels())
		assert.Equal(t, s.UUID(), p.currentGameStateUUID)
		assert.Equal(t, *resp, s.CurrentResponse())
	})

	t.Run("invalid game", func(t *testing.T) {
		_, _, err := Start(&Game{}, p)
		assert.NotNil(t, err)
		assert.Equal(t, 1, p.gamesStarted)
	})

	t.Run("invalid player", func(t *testing.T) {
		_, _, err := Start(g, &Player{})
		assert.NotNil(t, err)
	})
}

func TestState_Update(t *testing.T) {
	g := newValidTestUrbanGame()
	p := newValidTestPlayer()
	s, _, err := Start(g, p)
	require.NoError(t, err)

	l1 := g.levels[0]
	l2 := g.levels[1]
	l3 := g.levels[2]

	t.Run("invalid game", func(t *testing.T) {
		_, err := s.Update(&Game{}, l1.answers[0], p)
		assert.NotNil(t, err)
	})

	t.Run("level 1 wrong answer", func(t *testing.T) {
		resp, err := s.Update(g, "wrong answer", p)
		assert.Nil(t, err)
		assert.Equal(t, ClueResponse, resp.Kind)
		assert.Equal(t, "", resp.LevelTitle)
		assert.Equal(t, "", resp.LevelDescription)
		assert.NotEqual(t, "", resp.Clue)
		assert.Equal(t, "", resp.EndMessage)
		assert.Equal(t, *resp, s.CurrentResponse())
		assert.Equal(t, 0, p.TotalPoints())

		clue := resp.Clue

		resp, err = s.Update(g, "wrong answer", p)
		assert.Nil(t, err)
		assert.NotEqual(t, "", resp.Clue)
		assert.NotEqual(t, clue, resp.Clue)
		assert.Equal(t, *resp, s.CurrentResponse())
		assert.Equal(t, 0, p.TotalPoints())

		clue = resp.Clue

		resp, err = s.Update(g, "wrong answer", p)
		assert.Nil(t, err)
		assert.NotEqual(t, "", resp.Clue)
		assert.NotEqual(t, clue, resp.Clue)
		assert.Equal(t, *resp, s.CurrentResponse())
		assert.Equal(t, 0, p.TotalPoints())

		clue = resp.Clue

		resp, err = s.Update(g, "wrong answer", p)
		assert.Nil(t, err)
		assert.NotEqual(t, "", resp.Clue)
		assert.Equal(t, clue, resp.Clue) // The last clue should repeat.
		assert.Equal(t, *resp, s.CurrentResponse())
		assert.Equal(t, 0, p.TotalPoints())
	})

	t.Run("level 1 correct answer", func(t *testing.T) {
		resp, err := s.Update(g, l1.answers[0], p)
		assert.Nil(t, err)
		assert.Equal(t, LevelResponse, resp.Kind)
		assert.NotEqual(t, "", resp.LevelTitle)
		assert.NotEqual(t, "", resp.LevelDescription)
		assert.Equal(t, "", resp.Clue)
		assert.Equal(t, "", resp.EndMessage)
		assert.Equal(t, *resp, s.CurrentResponse())
		assert.Equal(t, 0, p.TotalPoints())
	})

	t.Run("end game", func(t *testing.T) {
		resp, err := s.Update(g, l2.answers[0], p)
		assert.Nil(t, err)
		assert.Equal(t, LevelResponse, resp.Kind)
		assert.NotEqual(t, "", resp.LevelTitle)
		assert.NotEqual(t, "", resp.LevelDescription)
		assert.Equal(t, "", resp.Clue)
		assert.Equal(t, "", resp.EndMessage)
		assert.Equal(t, *resp, s.CurrentResponse())

		resp, err = s.Update(g, l3.answers[0], p)
		assert.Nil(t, err)
		assert.Equal(t, EndResponse, resp.Kind)
		assert.Equal(t, "", resp.LevelTitle)
		assert.Equal(t, "", resp.LevelDescription)
		assert.Equal(t, "", resp.Clue)
		assert.NotEqual(t, "", resp.EndMessage)
		assert.Equal(t, *resp, s.CurrentResponse())
		assert.Equal(t, g.Value(), p.TotalPoints())
	})
}
