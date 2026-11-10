package utils

import "unicode"

func MixCase(s string) string {
	runes := []rune(s)
	upper := true
	for i, r := range runes {
		if !unicode.IsLetter(r) {
			continue
		}
		if upper {
			runes[i] = unicode.ToUpper(r)
		} else {
			runes[i] = unicode.ToLower(r)
		}
		upper = !upper
	}
	return string(runes)
}
