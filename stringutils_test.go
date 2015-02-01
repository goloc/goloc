package goloc

import (
	"testing"
)

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
	{"Élysées", "elysees", true, 2},
}

func TestLevenshteinDistance(t *testing.T) {
	for _, tt := range levTests {
		d := LevenshteinDistance(tt.source, tt.target, tt.ignoreCase)
		if d != tt.distance {
			t.Logf("Distance of %v and %v should be %v but was %v.",
				tt.source, tt.target, tt.distance, d)
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
	{"Rue du Square Carpeaux 75018 Paris France", "R D KR KP 7 PR RS"},
	{"a23 34u 78004R 345TE", "A2 3A 7 3S"},
	{"cacécîcocücy", "KSSKKS"},
}

func TestPartialphone(t *testing.T) {
	for _, tt := range partialphoneTests {
		target := Partialphone(tt.source)
		if target != tt.target {
			t.Logf("Partialphone of %v should be %v but was %v.",
				tt.source, tt.target, target)
			t.Fail()
		}
	}
}
