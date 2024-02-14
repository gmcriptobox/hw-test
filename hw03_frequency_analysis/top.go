package hw03frequencyanalysis

import (
	"sort"
	"strings"
	"unicode"
)

type pair struct {
	key   string
	value int
}

func Top10(str string) []string {
	if len(str) == 0 {
		return nil
	}
	words := strings.Fields(strings.ToLower(str))
	countMap := make(map[string]int, len(words))
	for _, word := range words {
		runeWord := []rune(word)
		left := 0
		right := 0
		if unicode.IsPunct(runeWord[left]) {
			left++
		}
		if unicode.IsPunct(runeWord[len(runeWord)-1-right]) {
			right++
		}
		if len(runeWord)-right >= left {
			countMap[string(runeWord[left:len(runeWord)-right])]++
		}
	}
	pairs := make([]pair, 0, len(words))
	i := 0
	for key, value := range countMap {
		pairs[i] = pair{key, value}
		i++
	}
	sort.Slice(pairs, func(i, j int) bool {
		if pairs[i].value == pairs[j].value {
			return pairs[i].key < pairs[j].key
		}
		return pairs[i].value > pairs[j].value
	})
	var result = make([]string, 10)
	for i := range result {
		if i >= len(pairs) {
			break
		}
		result[i] = pairs[i].key
	}
	return result
}
