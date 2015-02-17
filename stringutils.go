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

func ToUpper(r rune) rune {
	up := unicode.ToUpper(r)
	switch up {
	case 'Á', 'À', 'Ã', 'Ä', 'Å', 'Ā', 'Æ':
		return 'A'

	case 'È', 'É', 'Ê', 'Ë', 'Ē', 'Œ':
		return 'E'

	case 'Ì', 'Í', 'Î', 'Ï', 'Ī':
		return 'I'

	case 'Ó', 'Ò', 'Ô', 'Õ', 'Ö', 'Ø', 'Ō':
		return 'O'

	case 'Ú', 'Ù', 'Û', 'Ü', 'Ū':
		return 'U'

	case 'Y', 'Ý', 'Ŷ', 'Ÿ':
		return 'Y'
	default:
		return up
	}
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
					if ToUpper(runeSearch) != ToUpper(runeRef) {
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
	var match, topMatch, m, bestIndex, l, lTotal, i int
	var currentSearchWord, currentRefenceWord string
	lastIndex := -1
	for _, currentSearchWord = range searchWords {
		topMatch = 0
		bestIndex = 0
		l = len(currentSearchWord)
		lTotal += l
		for i, currentRefenceWord = range referenceWords {
			m = l - LevenshteinDistance(currentSearchWord, currentRefenceWord, true)
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
		return (defaultMaxScore * match) / lTotal
	}
}

func Partialphone(source string) string {
	b := bytes.NewBufferString("")
	lastRune := ' '
	penultimateRune := ' '

	for _, currentRune := range source {
		if unicode.IsPunct(currentRune) || unicode.IsSpace(currentRune) || unicode.IsControl(currentRune) || unicode.IsSymbol(currentRune) || unicode.IsMark(currentRune) {
			currentRune = ' '
		} else if unicode.IsNumber(currentRune) {
			if !unicode.IsNumber(lastRune) && lastRune != ' ' {
				b.WriteRune(' ')
			}
		} else {
			if unicode.IsNumber(lastRune) {
				b.WriteRune(' ')
			}
		}

		switch currentRune {
		case 'B', 'b':
			if lastRune != 0 && penultimateRune != lastRune {
				penultimateRune = lastRune
			}
			lastRune = 'B'

		case 'C', 'c':
			if lastRune != 0 && penultimateRune != lastRune {
				penultimateRune = lastRune
			}
			lastRune = 'C'

		case 'Ç', 'ç':
			if lastRune != 0 && penultimateRune != lastRune {
				penultimateRune = lastRune
			}
			lastRune = 'S'

		case 'D', 'd':
			if lastRune != 0 && penultimateRune != lastRune {
				penultimateRune = lastRune
			}
			lastRune = 'D'

		case 'F', 'f':
			if lastRune != 0 && penultimateRune != lastRune {
				penultimateRune = lastRune
			}
			lastRune = 'F'

		case 'G', 'g':
			if lastRune != 0 && penultimateRune != lastRune {
				penultimateRune = lastRune
			}
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
			if lastRune != 0 && penultimateRune != lastRune {
				penultimateRune = lastRune
			}
			lastRune = 'J'

		case 'K', 'k':
			if lastRune != 0 && penultimateRune != lastRune {
				penultimateRune = lastRune
			}
			lastRune = 'K'

		case 'L', 'l':
			if lastRune != 0 && penultimateRune != lastRune {
				penultimateRune = lastRune
			}
			lastRune = 'L'

		case 'M', 'm':
			if lastRune != 0 && penultimateRune != lastRune {
				penultimateRune = lastRune
			}
			lastRune = 'M'

		case 'N', 'n', 'Ñ', 'ñ':
			if lastRune != 0 && penultimateRune != lastRune {
				penultimateRune = lastRune
			}
			lastRune = 'N'

		case 'P', 'p':
			if lastRune != 0 && penultimateRune != lastRune {
				penultimateRune = lastRune
			}
			lastRune = 'P'

		case 'Q', 'q':
			if lastRune != 0 && penultimateRune != lastRune {
				penultimateRune = lastRune
			}
			lastRune = 'K'

		case 'R', 'r':
			if lastRune != 0 && penultimateRune != lastRune {
				penultimateRune = lastRune
			}
			lastRune = 'R'

		case 'S', 's', 'ẞ', 'ß':
			if lastRune != 0 && penultimateRune != lastRune {
				penultimateRune = lastRune
			}
			lastRune = 'S'

		case 'T', 't':
			if lastRune != 0 && penultimateRune != lastRune {
				penultimateRune = lastRune
			}
			lastRune = 'T'

		case 'V', 'v':
			if lastRune != 0 && penultimateRune != lastRune {
				penultimateRune = lastRune
			}
			lastRune = 'V'

		case 'W', 'w':
			if lastRune != 0 && penultimateRune != lastRune {
				penultimateRune = lastRune
			}
			lastRune = 'V'

		case 'X', 'x':
			if lastRune != 0 && penultimateRune != lastRune {
				penultimateRune = lastRune
			}
			lastRune = 'S'

		case 'Z', 'z':
			if lastRune != 0 && penultimateRune != lastRune {
				penultimateRune = lastRune
			}
			lastRune = 'S'

		default:
			if unicode.IsPunct(currentRune) || unicode.IsSpace(currentRune) || unicode.IsControl(currentRune) || unicode.IsSymbol(currentRune) || unicode.IsMark(currentRune) {
				switch lastRune {
				case ' ':

				default:
					if penultimateRune != 0 && penultimateRune != lastRune && penultimateRune != ' ' && penultimateRune != 'A' && !unicode.IsNumber(penultimateRune) {
						b.WriteRune(penultimateRune)
						penultimateRune = ' '
					}
					if lastRune != 0 && lastRune != ' ' && lastRune != 'A' && !unicode.IsNumber(lastRune) {
						b.WriteRune(lastRune)
					}
					b.WriteRune(' ')
					lastRune = ' '
				}
			} else if unicode.IsNumber(currentRune) {
				if lastRune == ' ' {
					b.WriteRune(currentRune)
					lastRune = currentRune
				} else if lastRune == 'A' {
					b.WriteRune(currentRune)
					lastRune = currentRune
				}

			} else {
				if lastRune == ' ' {
					b.WriteRune('A')
					lastRune = 'A'
				} else if unicode.IsNumber(lastRune) {
					b.WriteRune('A')
					lastRune = 'A'
				} else if lastRune == 'A' {
					lastRune = 'A'
				} else if lastRune == 'C' {
					switch unicode.ToUpper(currentRune) {
					case 'E', 'È', 'É', 'Ê', 'Ë', 'Ē', 'I', 'Ì', 'Í', 'Î', 'Ï', 'Ī', 'Y', 'Ý', 'Ŷ', 'Ÿ':
						if penultimateRune != 0 && penultimateRune != lastRune && penultimateRune != ' ' && penultimateRune != 'A' && !unicode.IsNumber(penultimateRune) {
							b.WriteRune(penultimateRune)
							penultimateRune = ' '
						}
						b.WriteRune('S')
					default:
						if penultimateRune != 0 && penultimateRune != lastRune && penultimateRune != ' ' && penultimateRune != 'A' && !unicode.IsNumber(penultimateRune) {
							b.WriteRune(penultimateRune)
							penultimateRune = ' '
						}
						b.WriteRune('K')
					}
					lastRune = 0
				} else {
					if lastRune != 0 {
						if penultimateRune != 0 && penultimateRune != lastRune && penultimateRune != ' ' && penultimateRune != 'A' && !unicode.IsNumber(penultimateRune) {
							b.WriteRune(penultimateRune)
							penultimateRune = ' '
						}
						b.WriteRune(lastRune)
					}
					lastRune = 0
				}
			}

		}
	}

	if lastRune != ' ' && lastRune != 0 && lastRune != 'A' {
		if penultimateRune != 0 && penultimateRune != lastRune && penultimateRune != ' ' && penultimateRune != 'A' && !unicode.IsNumber(penultimateRune) {
			b.WriteRune(penultimateRune)
			penultimateRune = ' '
		}
		b.WriteRune(lastRune)
	}

	return b.String()
}

func Nkeys(keys []string) map[string]bool {
	var i, j, l int
	mapKeys := make(map[string]bool)
	for _, k := range keys {
		l = len(k)
		i = l
		if l >= 2 {
			i = 1
		}
		for ; i <= l; i++ {
			for j = 0; i+j <= l; j++ {
				subk := k[j : j+i]
				mapKeys[subk] = true
			}
		}
	}
	return mapKeys
}
