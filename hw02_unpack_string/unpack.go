package hw02unpackstring

import (
	"errors"
	"strings"
)

var ErrInvalidString = errors.New("invalid string")

var isNumber = func(c rune) bool { return c >= '0' && c <= '9' }

//var isSlash = func(c rune) bool { return c == '\\' }

func appendToStrBuilder(strBuild *strings.Builder, char rune, count int) {
	for i := 0; i < count; i++ {
		strBuild.WriteRune(char)
	}
}

func Unpack(str string) (string, error) {
	//проверка крайнего случая, когда строка пуста
	if len(str) == 0 {
		return "", nil
	}
	builder := strings.Builder{}
	runeSlice := []rune(str)
	slashCount := 0
	//идем по слайсу символов, при встрече числа добавляем n раз предыдущий символ
	for i, c := range runeSlice {
		if isNumber(c) {
			if i > 0 {
				if isNumber(runeSlice[i-1-slashCount]) {
					return "", ErrInvalidString
				}
				appendToStrBuilder(&builder, runeSlice[i-1-slashCount], int(c-'0'))
			} else {
				//крайний случай, если первым символом идет цифра
				return "", ErrInvalidString
			}
		} else if i > 0 && !isNumber(runeSlice[i-1]) {
			builder.WriteRune(runeSlice[i-1])
		}
	}
	lenSlice := len(runeSlice)
	if !isNumber(runeSlice[lenSlice-1]) {
		//добавляем последний символ, если там не цифра
		builder.WriteRune(runeSlice[lenSlice-1])
	} else if lenSlice > 2 {
		//проверка крайнего случая, когда последний и предпоследний символ цифры
		if isNumber(runeSlice[lenSlice-2]) {
			return "", ErrInvalidString
		}
	}
	return builder.String(), nil
}
