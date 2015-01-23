package goloc

import (
	"fmt"
	"testing"
)

func TestSplitSpacePunct(t *testing.T) {
	res := splitSpacePunct("start er  ee r '-('([[-||]]) zée-ee -_ end")
	fmt.Printf("nb = %d\n", len(res))
	if len(res) != 7 {
		t.Fail()
	}
	for i, e := range res {
		fmt.Printf("%d %t\n", i, e)
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
	res := stripAccents("zùeèàüî~Ýa erÆ a")
	fmt.Printf("%t -> %t\n", source, res)
	if res != "zueeaui~Ya erAE a" {
		t.Fail()
	}
	res = stripAccents("")
	if res != "" {
		t.Fail()
	}
}

var levTests = []struct {
	source   string
	target   string
	distance int
}{
	// two empty
	{"", "", 0},
	// deletion
	{"library", "librar", 1},
	// one empty, left
	{"", "library", 7},
	// one empty, right
	{"library", "", 7},
	{"car", "cars", 1},
	{"", "a", 1},
	{"a", "aa", 1},
	{"a", "aaa", 2},
	{"", "", 0},
	{"a", "b", 1},
	{"aaa", "aba", 1},
	{"aaa", "ab", 2},
	{"a", "a", 0},
	{"ab", "ab", 0},
	{"a", "", 1},
	{"aa", "a", 1},
	{"aaa", "a", 2},
	// unicode
	{"Schüßler", "Schübler", 1},
	{"Schüßler", "Schußler", 1},
	{"Schüßler", "Schüßler", 0},
	{"Schüßler", "Schüler", 1},
	{"Schüßler", "Schüßlers", 1},
}

func TestLevenshteinDistance(t *testing.T) {
	for _, tt := range levTests {
		d := levenshteinDistance(tt.source, tt.target)
		fmt.Printf("%t - %t -> %d\n", tt.source, tt.target, d)
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
	{"Avenue des Champs-Élysées 75008 Paris France", "AFN D S ALS 7 PR RS"},
	{"Rue du Square Carpeaux 75018 Paris France", "R D KR SP 7 PR RS"},
	{"i ii az tt rre ke_me-pe", "A A A  R K M P"},
	{"a23 34u 78004R 345TE", "A2 3A 7 3S"},
}

func TestPartialphone(t *testing.T) {
	for _, tt := range partialphoneTests {
		target := partialphone(tt.source)
		fmt.Printf("%t -> %t\n", tt.source, target)
		if target != tt.target {
			t.Fail()
		}
	}
}
