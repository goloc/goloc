// Copyright 2015 Mathieu MAST. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.
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

func MSplit(source string) []string {
	strs := Split(source)
	mapsplit := make(map[string]bool)
	for _, str := range strs {
		mapsplit[str] = true
		k := 0

		for i, _ := range str {
			l := 0
			for j, _ := range str {
				if j > i {
					if i == 0 {
						mapsplit[str[i:j]] = true
					} else if l > k+1 {
						mapsplit[str[i:j]] = true
					}
				}
				l++
			}
			k++
		}
	}
	mplit := make([]string, len(mapsplit))
	i := 0
	for str, _ := range mapsplit {
		mplit[i] = str
		i++
	}
	return mplit
}

func UpperUnaccentUnpunctString(str string) string {
	bs := bytes.NewBufferString("")
	for _, r := range str {
		bs.WriteRune(UpperUnaccentUnpunctRune(r))
	}
	return bs.String()
}

func UpperUnaccentUnpunctRune(r rune) rune {
	if unicode.IsLetter(r) {
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

		case 'Ý', 'Ŷ', 'Ÿ':
			return 'Y'

		case 'Ñ':
			return 'N'

		case 'Ŵ':
			return 'W'

		default:
			return up
		}
	} else if unicode.IsNumber(r) {
		return r
	} else {
		return ' '
	}
}

func Distance(search, reference string, ignoreCase bool) int {
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
					if UpperUnaccentUnpunctRune(runeSearch) != UpperUnaccentUnpunctRune(runeRef) {
						cost = 2
					}
				} else {
					cost = 2
				}
			}
			column[y] = Min(Min(
				column[y]+1,    // insert on search
				column[y-1]+2), // delete on search
				lastdiag+cost) // substitution on search
			lastdiag = olddiag
			y++
		}
		x++
	}
	return column[lenSearch]
}

func Score(search, reference string) int {
	if search == "" || reference == "" {
		return 0
	}
	searchWords := Split(search)
	referenceWords := Split(reference)
	var match, topMatch, m, bestIndex, lenCurrentSearchWord, lenTotal, i int
	var currentSearchWord, currentRefenceWord string
	lastIndex := -1
	for _, currentSearchWord = range searchWords {
		topMatch = 0
		bestIndex = 0
		lenCurrentSearchWord = len(currentSearchWord)
		lenTotal += lenCurrentSearchWord
		for i, currentRefenceWord = range referenceWords {
			m = lenCurrentSearchWord - Distance(currentSearchWord, currentRefenceWord, true)
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
		return 100 * match
	}
}

func PartialphoneWriteLast(b *bytes.Buffer, ptrCurrentRune, ptrLastRune, ptrPenultimateRune *rune) {
	if *ptrLastRune != ' ' && *ptrPenultimateRune != *ptrLastRune {
		if *ptrLastRune != 'A' {
			if *ptrLastRune == 'C' {
				switch *ptrCurrentRune {
				case 'H':
					*ptrLastRune = 'S'
				case 'E', 'I', 'Y':
					*ptrLastRune = 'S'
				default:
					*ptrLastRune = 'K'
				}
			}
			if *ptrLastRune == 'P' && *ptrCurrentRune == 'H' {
				*ptrLastRune = 'F'
			}
			b.WriteRune(*ptrLastRune)
		}
		*ptrPenultimateRune = *ptrLastRune
	}
}

func Partialphone(source string) string {
	b := bytes.NewBufferString("")
	lastRune := ' '
	penultimateRune := ' '

	for _, currentRune := range source {
		currentRune = UpperUnaccentUnpunctRune(currentRune)

		switch currentRune {

		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			if lastRune != ' ' && penultimateRune != lastRune {
				if lastRune != 'A' && lastRune != '1' {
					b.WriteRune(lastRune)
				}
				penultimateRune = lastRune
			}
			if lastRune != '1' {
				b.WriteRune(currentRune)
			}
			lastRune = '1'

		case 'A', 'E', 'I', 'O', 'U', 'Y':
			if lastRune == 'C' {
				switch currentRune {
				case 'E', 'I', 'Y':
					lastRune = 'S'
				default:
					lastRune = 'K'
				}
			}
			if lastRune != ' ' && penultimateRune != lastRune {
				if lastRune != 'A' && lastRune != '1' {
					b.WriteRune(lastRune)
				}
				penultimateRune = lastRune
			}
			if lastRune == ' ' || lastRune == '1' {
				b.WriteRune('A')
			}
			lastRune = 'A'

		case 'B':
			PartialphoneWriteLast(b, &currentRune, &lastRune, &penultimateRune)
			lastRune = currentRune

		case 'C':
			PartialphoneWriteLast(b, &currentRune, &lastRune, &penultimateRune)
			lastRune = currentRune

		case 'Ç':
			currentRune = 'S'
			PartialphoneWriteLast(b, &currentRune, &lastRune, &penultimateRune)
			lastRune = currentRune

		case 'D':
			PartialphoneWriteLast(b, &currentRune, &lastRune, &penultimateRune)
			lastRune = currentRune

		case 'F':
			PartialphoneWriteLast(b, &currentRune, &lastRune, &penultimateRune)
			lastRune = currentRune

		case 'G':
			currentRune = 'J'
			PartialphoneWriteLast(b, &currentRune, &lastRune, &penultimateRune)
			lastRune = currentRune

		case 'H':
			PartialphoneWriteLast(b, &currentRune, &lastRune, &penultimateRune)

		case 'J':
			PartialphoneWriteLast(b, &currentRune, &lastRune, &penultimateRune)
			lastRune = currentRune

		case 'K':
			PartialphoneWriteLast(b, &currentRune, &lastRune, &penultimateRune)
			lastRune = currentRune

		case 'L':
			PartialphoneWriteLast(b, &currentRune, &lastRune, &penultimateRune)
			lastRune = currentRune

		case 'M':
			PartialphoneWriteLast(b, &currentRune, &lastRune, &penultimateRune)
			lastRune = currentRune

		case 'N':
			PartialphoneWriteLast(b, &currentRune, &lastRune, &penultimateRune)
			lastRune = currentRune

		case 'P':
			PartialphoneWriteLast(b, &currentRune, &lastRune, &penultimateRune)
			lastRune = currentRune

		case 'Q':
			currentRune = 'K'
			PartialphoneWriteLast(b, &currentRune, &lastRune, &penultimateRune)
			lastRune = currentRune

		case 'R':
			PartialphoneWriteLast(b, &currentRune, &lastRune, &penultimateRune)
			lastRune = currentRune

		case 'S', 'ẞ', 'ß':
			currentRune = 'S'
			PartialphoneWriteLast(b, &currentRune, &lastRune, &penultimateRune)
			lastRune = currentRune

		case 'T':
			PartialphoneWriteLast(b, &currentRune, &lastRune, &penultimateRune)
			lastRune = currentRune

		case 'V':
			PartialphoneWriteLast(b, &currentRune, &lastRune, &penultimateRune)
			lastRune = currentRune

		case 'W':
			currentRune = 'W'
			PartialphoneWriteLast(b, &currentRune, &lastRune, &penultimateRune)
			lastRune = currentRune

		case 'X':
			currentRune = 'S'
			PartialphoneWriteLast(b, &currentRune, &lastRune, &penultimateRune)
			lastRune = currentRune

		case 'Z':
			currentRune = 'S'
			PartialphoneWriteLast(b, &currentRune, &lastRune, &penultimateRune)
			lastRune = currentRune

		default:
			if lastRune != ' ' && lastRune != 'A' && lastRune != '1' && penultimateRune != lastRune {
				b.WriteRune(lastRune)
				penultimateRune = lastRune
			}
			if lastRune != ' ' {
				b.WriteRune(currentRune)
			}
			lastRune = ' '
		}
	}

	if lastRune != ' ' && penultimateRune != lastRune {
		if lastRune != 'A' {
			b.WriteRune(lastRune)
		}
	}

	return b.String()
}
