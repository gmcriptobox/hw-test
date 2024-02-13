package hw02unpackstring

import (
	"errors"
	"strings"
)

var ErrInvalidString = errors.New("invalid string")

var isNumber = func(c rune) bool { return c >= '0' && c <= '9' }
var isSlash = func(c rune) bool { return c == '\\' }
var isChar = func(c rune) bool { return !isNumber(c) && !isSlash(c) }

func appendToStrBuilder(strBuild *strings.Builder, c rune, count int) {
	for i := 0; i < count; i++ {
		strBuild.WriteRune(c)
	}
}

func slashHandle(i int, runeSlice []rune, strBuild *strings.Builder) int {
	c1 := runeSlice[i+1]
	firdIndex := i + 2
	if firdIndex < len(runeSlice) {
		if isNumber(runeSlice[firdIndex]) {
			appendToStrBuilder(strBuild, c1, int((runeSlice[firdIndex] - '0')))
			i += 2
		} else {
			strBuild.WriteRune(c1)
			i++
		}
	} else {
		strBuild.WriteRune(c1)
		i++
	}
	return i
}

func Unpack(str string) (string, error) {
	// проверка крайнего случая, когда строка пуста
	if len(str) == 0 {
		return "", nil
	}
	runeSlice := []rune(str)
	lenSlice := len(runeSlice)
	// проверка крайних случаев, одиночный слэш и число 1 символом
	if lenSlice == 1 && isSlash(runeSlice[0]) || isNumber(runeSlice[0]) {
		return "", ErrInvalidString
	}
	builder := strings.Builder{}
	for i := 0; i < lenSlice-1; i++ {
		c := runeSlice[i]
		c1 := runeSlice[i+1]
		if isNumber(c) || isSlash(c) && isChar(c1) {
			return "", ErrInvalidString
		} else if isSlash(c) && !isChar(c1) {
			i = slashHandle(i, runeSlice, &builder)
		} else if isChar(c) {
			if isNumber(c1) {
				appendToStrBuilder(&builder, c, int(c1-'0'))
				i++
			} else {
				builder.WriteRune(c)
			}
		}
	}
	if !isNumber(runeSlice[lenSlice-1]) {
		builder.WriteRune(runeSlice[lenSlice-1])
	}
	return builder.String(), nil
}
