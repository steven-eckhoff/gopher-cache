package rand

// String returns a string with a specified rune count and max rune size.
func String(runeCount, maxRuneSize int) string {
	var codePointLimit int

	switch maxRuneSize {
	case 1:
		codePointLimit = 0x80
	case 2:
		codePointLimit = 0x800
	case 3:
		codePointLimit = 0x10000
	case 4:
		codePointLimit = 0x110000
	default:
		codePointLimit = 0x110000
	}

	codePoints := make([]rune, runeCount)

	for i := 0; i < runeCount; i++ {
		codePoints[i] = rune(sRand.Intn(codePointLimit))
	}

	return string(codePoints)
}
