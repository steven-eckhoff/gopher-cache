package rand

import (
	"github.com/stretchr/testify/assert"
	"math/rand"
	"testing"
	"unicode/utf8"
)

func TestString(t *testing.T) {
	// Test various rune lengths and string lengths.
	for inputRuneLen := 0; inputRuneLen <= 5; inputRuneLen++ {
		expectedRuneCount := rand.Intn(100)

		var expectedRuneLen int
		if inputRuneLen == 0 || inputRuneLen == 5 {
			expectedRuneLen = 4
		} else {
			expectedRuneLen = inputRuneLen
		}

		s := String(expectedRuneCount, inputRuneLen)

		assert.Equal(t, expectedRuneCount, utf8.RuneCountInString(s))

		for _, r := range s {
			assert.GreaterOrEqual(t, expectedRuneLen, utf8.RuneLen(r))
		}
	}

	// Test for sufficient randomness
	strings := make(map[string]bool)
	count := 1000
	for i := 0; i < count; i++ {
		strings[String(100, 1)] = true
	}

	assert.Equal(t, count, len(strings))
}
