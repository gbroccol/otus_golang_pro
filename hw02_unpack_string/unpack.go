package hw02unpackstring

import (
	"errors"
	"strings"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(str string) (string, error) {
	var result strings.Builder
	i := 0

	for i < len(str) {
		letter := str[i]

		if unicode.IsDigit(rune(letter)) {
			return "", ErrInvalidString
		}

		if letter == 92 && i+1 < len(str) {
			i++
			letter = str[i]
		} else if letter == 92 {
			return "", ErrInvalidString
		}

		if i+1 < len(str) {
			digit := rune(str[i+1])
			if unicode.IsDigit(digit) {
				for amount := int(digit - '0'); amount > 0; amount-- {
					result.WriteByte(letter)
				}
				i++
			} else {
				result.WriteByte(letter)
			}
		} else {
			result.WriteByte(letter)
		}
		i++
	}
	return result.String(), nil
}
