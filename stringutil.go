package chart

import "strings"

// SplitCSV splits a corpus by the `,`, dropping leading or trailing whitespace unless quoted.
func SplitCSV(text string) (output []string) {
	if len(text) == 0 {
		return
	}

	var state int
	var word []rune
	var opened rune
	for _, r := range text {
		switch state {
		case 0: // word
			switch {
			case isQuote(r):
				opened = r
				state = 1
			case isCSVDelim(r):
				output = append(output, strings.TrimSpace(string(word)))
				word = nil
			default:
				word = append(word, r)
			}
		case 1: // we're in a quoted section
			if matchesQuote(opened, r) {
				state = 0
				continue
			}
			word = append(word, r)
		}
	}

	if len(word) > 0 {
		output = append(output, strings.TrimSpace(string(word)))
	}
	return
}

func isCSVDelim(r rune) bool {
	return r == rune(',')
}

func isQuote(r rune) bool {
	return r == '"' || r == '\'' || r == '“' || r == '”' || r == '`'
}

func matchesQuote(a, b rune) bool {
	if a == '“' && b == '”' {
		return true
	}
	if a == '”' && b == '“' {
		return true
	}
	return a == b
}
