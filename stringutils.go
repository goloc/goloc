package goloc

import (
	"bytes"
	"strings"
	"unicode"
)

func Split(source string) []string {
	return strings.FieldsFunc(source, func(r rune) bool {
		return !unicode.IsLetter(r) && !unicode.IsNumber(r)
	})
}

func LevenshteinDistance(search string, reference string, ignoreCase bool) int {
	var cost, lastdiag, olddiag int
	lenSearch := 0
	for range search {
		lenSearch++
	}
	column := make([]int, lenSearch+1)
	for y := 0; y <= lenSearch; y++ {
		column[y] = y
	}
	x := 1
	for _, runeRef := range reference {
		column[0] = x
		lastdiag = x - 1
		y := 1
		for _, runeSearch := range search {
			olddiag = column[y]
			cost = 0
			if runeSearch != runeRef {
				if ignoreCase == true {
					if unicode.ToUpper(runeSearch) != unicode.ToUpper(runeRef) {
						cost = 1
					}
				} else {
					cost = 1
				}
			}
			column[y] = Min(Min(
				column[y]+1,
				column[y-1]+1),
				lastdiag+cost)
			lastdiag = olddiag
			y++
		}
		x++
	}
	return column[lenSearch]
}

func Score(search string, reference string) int {
	defaultMaxScore := 1000
	if search == "" || reference == "" {
		return 0
	}
	searchWords := Split(search)
	referenceWords := Split(reference)
	var match, topMatch, m, bestIndex int
	lastIndex := -1
	for _, currentSearchWord := range searchWords {
		topMatch = 0
		bestIndex = 0
		for i, currentRefenceWord := range referenceWords {
			m = len(currentRefenceWord) - LevenshteinDistance(currentSearchWord, currentRefenceWord, true)
			if m > topMatch {
				topMatch = m
				bestIndex = i
			}
		}
		if lastIndex == -1 {
			lastIndex = bestIndex
		}
		match += topMatch
		if bestIndex < lastIndex {
			match--
		}
		lastIndex = bestIndex
	}
	if match < 0 {
		return 0
	} else {
		return (defaultMaxScore * match) / len(search)
	}
}

func Partialphone(source string) string {
	b := bytes.NewBufferString("")
	lastRune := ' '

	for _, currentRune := range source {
		switch currentRune {
		case 'B', 'b':
			lastRune = 'B'

		case 'C', 'c':
			lastRune = 'C'

		case 'Ç', 'ç':
			lastRune = 'S'

		case 'D', 'd':
			lastRune = 'D'

		case 'F', 'f':
			lastRune = 'F'

		case 'G', 'g':
			lastRune = 'J'

		case 'H', 'h':
			if lastRune == 'P' {
				lastRune = 'F'
			} else if lastRune == 'C' {
				lastRune = 'S'
			} else {
				// Silent
			}

		case 'J', 'j':
			lastRune = 'J'

		case 'K', 'k':
			lastRune = 'K'

		case 'L', 'l':
			lastRune = 'L'

		case 'M', 'm':
			lastRune = 'M'

		case 'N', 'n', 'Ñ', 'ñ':
			lastRune = 'N'

		case 'P', 'p':
			lastRune = 'P'

		case 'Q', 'q':
			lastRune = 'K'

		case 'R', 'r':
			lastRune = 'R'

		case 'S', 's', 'ẞ', 'ß':
			lastRune = 'S'

		case 'T', 't':
			lastRune = 'S'

		case 'V', 'v':
			lastRune = 'V'

		case 'W', 'w':
			lastRune = 'V'

		case 'X', 'x':
			lastRune = 'S'

		case 'Z', 'z':
			lastRune = 'S'

		default:
			if unicode.IsPunct(currentRune) || unicode.IsSpace(currentRune) || unicode.IsControl(currentRune) || unicode.IsSymbol(currentRune) || unicode.IsMark(currentRune) {
				switch lastRune {
				case ' ':

				default:
					b.WriteRune(' ')
					lastRune = ' '
				}
			} else if unicode.IsNumber(currentRune) {
				if lastRune == ' ' || lastRune == 'A' {
					b.WriteRune(currentRune)
					lastRune = currentRune
				} else if unicode.IsNumber(lastRune) {
					lastRune = currentRune
				} else {
					if lastRune != 0 {
						b.WriteRune(lastRune)
					}
					lastRune = 0
				}

			} else {
				if lastRune == ' ' || unicode.IsNumber(lastRune) {
					b.WriteRune('A')
					lastRune = 'A'
				} else if lastRune == 'A' {
					lastRune = currentRune
				} else if lastRune == 'C' {
					switch unicode.ToUpper(currentRune) {
					case 'E', 'È', 'É', 'Ê', 'Ë', 'I', 'Ì', 'Í', 'Î', 'Ï', 'Y', 'Ŷ', 'Ÿ':
						b.WriteRune('S')
					default:
						b.WriteRune('K')
					}
					lastRune = 0
				} else {
					if lastRune != 0 {
						b.WriteRune(lastRune)
					}
					lastRune = 0
				}
			}

		}
	}

	return b.String()
}
