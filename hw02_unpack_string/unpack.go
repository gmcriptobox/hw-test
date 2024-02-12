package hw02unpackstring

import (
	"errors"
	"strings"
)

var ErrInvalidString = errors.New("invalid string")

var isNumber = func(c rune) bool { return c >= '0' && c <= '9' }

func appendToStrBuilder(strBuild *strings.Builder, char rune, count int) {
	for i := 0; i < count; i++ {
		strBuild.WriteRune(char)
	}
}

func unpackOneCicle(strBuild *strings.Builder, runeSlice []rune, i int, c rune) bool {
	if isNumber(c) { // nolint:nestif
		if i > 0 {
			if isNumber(runeSlice[i-1]) {
				return false
			}
			appendToStrBuilder(strBuild, runeSlice[i-1], int(c-'0'))
		} else {
			// крайний случай, если первым символом идет цифра
			return false
		}
	} else if i > 0 && !isNumber(runeSlice[i-1]) {
		strBuild.WriteRune(runeSlice[i-1])
	}
	return true
}

func Unpack(str string) (string, error) {
	// проверка крайнего случая, когда строка пуста
	if len(str) == 0 {
		return "", nil
	}
	builder := strings.Builder{}
	runeSlice := []rune(str)
	// идем по слайсу символов, при встрече числа добавляем n раз предыдущий символ
	for i, c := range runeSlice {
		if !unpackOneCicle(&builder, runeSlice, i, c) {
			return "", ErrInvalidString
		}
	}
	lenSlice := len(runeSlice)
	if !isNumber(runeSlice[lenSlice-1]) {
		// добавляем последний символ, если там не цифра
		builder.WriteRune(runeSlice[lenSlice-1])
	} else if lenSlice > 2 {
		// проверка крайнего случая, когда последний и предпоследний символ цифры
		if isNumber(runeSlice[lenSlice-2]) {
			return "", ErrInvalidString
		}
	}
	return builder.String(), nil
}
