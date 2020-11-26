package game

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNewPlayer(t *testing.T) {
	u := newTestUser()

	t.Run("valid", func(t *testing.T) {

		p, err := NewPlayerFromUser(u)
		require.NoError(t, err)

		assert.Equal(t, u.UUID(), p.UUID())
		assert.Equal(t, u.Number(), p.Number())
	})

	t.Run("invalid user", func(t *testing.T) {
		_, err := NewPlayerFromUser(User{})
		assert.NotNil(t, err)
	})
}

func newValidTestPlayer() *Player {
	p, err := NewPlayerFromUser(newTestUser())
	if err != nil {
		panic(err)
	}

	return p
}
