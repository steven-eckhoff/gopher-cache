// Package rand implements a library that extends the stdlib rand package.
package rand

import (
	"math/rand"
	"time"
)

var sRand *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))
