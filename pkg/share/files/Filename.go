package files

import (
	"fmt"
	"strings"
)

// -- constants --

const (
	BASE             = 62
	DIGIT_OFFSET     = 48
	LOWERCASE_OFFSET = 97 - 10
	UPPERCASE_OFFSET = 65 - 36
)

// -- types --

// an index-based html filename
type Filename int

// -- impls --

// builds a base62-encoded filename for an html file
func (f *Filename) String() (string, error) {
	if s, err := f.Encode(); err != nil {
		return "", err
	} else {
		return fmt.Sprintf("%s.html", s), nil
	}
}

// returns the base62-encoded index
// see: https://gist.github.com/Pagliacii/dca0f6b732c19045d258eaee81917071
func (f *Filename) Encode() (string, error) {
	digits := int(*f)

	// return "0" when i == 0
	if digits == 0 {
		return "0", nil
	}

	// otherwise, build a base62 encoded string digit by digit
	var sb strings.Builder
	for digits > 0 {
		// grab the next digit
		digit := digits % BASE
		if c, err := f.GetChar(digit); err != nil {
			return "", err
		} else {
			sb.WriteRune(c)
		}

		// update with what's left
		digits = int(digits / BASE)
	}

	return sb.String(), nil
}

func (*Filename) GetChar(ord int) (rune, error) {
	switch {
	case ord < 10:
		return rune(ord + DIGIT_OFFSET), nil
	case ord >= 10 && ord <= 35:
		return rune(ord + LOWERCASE_OFFSET), nil
	case ord >= 36 && ord < 62:
		return rune(ord + UPPERCASE_OFFSET), nil
	default:
		return '*', fmt.Errorf("%d is not a valid integer in the range of base %d", ord, BASE)
	}
}
