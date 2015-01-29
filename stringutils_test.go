package goloc

import (
	"testing"
)

func TestSplitSpacePunct(t *testing.T) {
	res := splitSpacePunct("start er  ee r '-('([[-||]]) zée-ee -_ end")
	if len(res) != 7 {
		t.Fail()
	}
	res = splitSpacePunct(" ")
	if len(res) != 0 {
		t.Fail()
	}
	res = splitSpacePunct("")
	if len(res) != 0 {
		t.Fail()
	}
}

func TestStripAccents(t *testing.T) {
	source := "zùeèàüî~Ýa erÆ a"
	res := stripAccents(source)
	if res != "zueeaui~Ya erAE a" {
		t.Fail()
	}
	res = stripAccents("")
	if res != "" {
		t.Fail()
	}
}

var levTests = []struct {
	source     string
	target     string
	ignoreCase bool
	distance   int
}{
	{"", "", false, 0},
	{"PARIS", "", false, 5},
	{"", "PARIS", false, 5},
	{"PARIS", "PARIS", false, 0},
	{"PaRIS", "PARiS", false, 2},
	{"PaRIS", "PARiS", true, 0},
	{"PARIS", "PARI", false, 1},
	{"PARS", "PARIS", false, 1},
	{"PAR", "PARIS", false, 2},
	{"PR", "PARIS", false, 3},
	{"PARIS", "FRANCE", false, 5},
	{"PĂRIS", "PARIŞ", false, 2},
	{"CAR", "Carpeaux", true, 5},
}

func TestLevenshteinDistance(t *testing.T) {
	for _, tt := range levTests {
		d := levenshteinDistance(tt.source, tt.target, tt.ignoreCase)
		if d != tt.distance {
			t.Fail()
		}
	}
}

var partialphoneTests = []struct {
	source string
	target string
}{
	{"", ""},
	{"Avenue des Champs-Élysées 75008 Paris France", "AVN D S ALS 7 PR RS"},
	{"Rue du Square Carpeaux 75018 Paris France", "R D KR SP 7 PR RS"},
	{"a23 34u 78004R 345TE", "A2 3A 7 3S"},
}

func TestPartialphone(t *testing.T) {
	for _, tt := range partialphoneTests {
		target := partialphone(tt.source)
		if target != tt.target {
			t.Fail()
		}
	}
}
