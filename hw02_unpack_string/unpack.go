package hw02unpackstring

import (
	"errors"  //nolint:depguard
	"strings" //nolint:depguard
	"unicode" //nolint:depguard
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(str string) (string, error) {
	var result strings.Builder
	runes := []rune(str) // Преобразуем строку в слайс рун для корректной обработки Unicode
	i := 0

	for i < len(runes) {
		letter := runes[i]

		if unicode.IsDigit(rune(letter)) {
			return "", ErrInvalidString
		}

		if letter == 92 && i+1 < len(str) {
			i++
			letter = runes[i]
		} else if letter == 92 {
			return "", ErrInvalidString
		}

		if i+1 < len(str) {
			digit := rune(str[i+1])
			if unicode.IsDigit(digit) {
				for amount := int(digit - '0'); amount > 0; amount-- {
					result.WriteRune(letter)
				}
				i++
			} else {
				result.WriteRune(letter)
			}
		} else {
			result.WriteRune(letter)
		}
		i++
	}
	return result.String(), nil
}
