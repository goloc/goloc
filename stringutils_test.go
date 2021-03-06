// Copyright 2015 Mathieu MAST. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.
package goloc

import (
	"testing"

	"github.com/goloc/container"
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
	{"PaRIS", "PARiS", false, 4},
	{"PaRIS", "PARiS", true, 0},
	{"PARI", "PARIS", false, 1},
	{"PARIS", "PARI", false, 2},
	{"PARS", "PARIS", false, 1},
	{"PAROS", "PARIS", false, 2},
	{"PAR", "PARIS", false, 2},
	{"PR", "PARIS", false, 3},
	{"PARIS", "FRANCE", false, 9},
	{"PĂRIS", "PARIŞ", false, 4},
	{"CAR", "Carpeaux", true, 5},
	{"Élysées", "elysees", false, 4},
	{"Élysées", "elysees", true, 0},
	{"eco", "Écaillon", true, 5},
	{"eco", "Ecole", true, 2},
	{"eco", "C", true, 3},
}

func TestDistance(t *testing.T) {
	for _, tt := range levTests {
		d := 0
		if tt.ignoreCase {
			d = Distance(UpperUnaccentUnpunctString(tt.source), UpperUnaccentUnpunctString(tt.target))
		} else {
			d = Distance(tt.source, tt.target)
		}
		if d != tt.distance {
			t.Logf("Distance of %v and %v should be %v but was %v.",
				tt.source, tt.target, tt.distance, d)
			t.Fail()
		}
	}
}

var cleanTests = []struct {
	source string
	target string
}{
	{"Avenue des Champs-Élysées 75008 Paris France", "CHAMPS ELYSEES 75008 PARIS FRANCE"},
	{"Rue du Square Carpeaux 75018 Paris France", "SQUARE CARPEAUX 75018 PARIS FRANCE"},
	{"Place Carnot 69002 Lyon", "CARNOT 69002 LYON"},
	{"Rue Tissot 69009 Lyon", "TISSOT 69009 LYON"},
	{"Tissot Rue Lyon", "TISSOT LYON"},
	{"Champs-Élysées Avenue 75008 Paris France des", "CHAMPS ELYSEES 75008 PARIS FRANCE"},
}

func TestClean(t *testing.T) {
	stopWords := container.NewLinkedList()
	stopWords.Add("Des")
	stopWords.Add("Du")
	stopWords.Add("AVENUE")
	stopWords.Add("RUE")
	stopWords.Add("PLACE")
	for _, tt := range cleanTests {
		target := Clean(tt.source, stopWords)
		if target != tt.target {
			t.Logf("Clean of %v should be %v but was %v.",
				tt.source, tt.target, target)
			t.Fail()
		}
	}
}

var partialphoneTests = []struct {
	source string
	target string
}{
	{"", ""},
	{"Brest", "BRST"},
	{"Avenue des Champs-Élysées 75008 Paris France", "AVN DS SMPS ALSS 7 PRS FRNS"},
	{"Rue du Square Carpeaux 75018 Paris France", "R D SKR KRPS 7 PRS FRNS"},
	{"Place Carnot 69002 Lyon", "PLS KRNT 6 LN"},
	{"Rue Tissot 69009 Lyon", "R TST 6 LN"},
	{"Grandclément", "JRNDKLMNT"},
	{"Grand Clément", "JRND KLMNT"},
	{"Grandchemin", "JRNDSMN"},
	{"Grand chemin", "JRND SMN"},
	{"Grandcarmeaux", "JRNDKRMS"},
	{"Grand carmeaux", "JRND KRMS"},
	{"Piscine Immeuble", "PSN AMBL"},
	{"Straßer STRAẞER", "STRSR STRSR"},
	{"Wagon", "VJN"},
	{"6", "6"},
	{"62", "6"},
	{"6 a23 T34 TER34 3423u 78045R 34567TEERT", "6 A2 T3 TR3 3A 7R 3TRT"},
	{" ", ""},
	{"!", ""},
	{"? !", ""},
	{"t", "T"},
	{"i", "A"},
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

func TestScore(t *testing.T) {
	score1 := ContainerScore(Split(UpperUnaccentUnpunctString("champs elysees paris")), UpperUnaccentUnpunctString("Avenue des Champs-Élysées 75008 Paris France"))
	score2 := ContainerScore(Split(UpperUnaccentUnpunctString("paris champs elysees")), UpperUnaccentUnpunctString("Avenue des Champs-Élysées 75008 Paris France"))
	score3 := ContainerScore(Split(UpperUnaccentUnpunctString("champs elyse paris")), UpperUnaccentUnpunctString("Avenue des Champs-Élysées 75008 Paris France"))
	if score1 <= score2 || score1 <= score3 {
		t.Logf("score1=%v score2=%v score3=%v", score1, score2, score3)
		t.Fail()
	}

	score4 := ContainerScore(Split(UpperUnaccentUnpunctString("eco")), UpperUnaccentUnpunctString("Ecole – Aubers"))
	score5 := ContainerScore(Split(UpperUnaccentUnpunctString("eco")), UpperUnaccentUnpunctString("Cite Arbrisseau All C 59176 Écaillon"))
	if score4 <= score5 {
		t.Logf("score4=%v score5=%v", score4, score5)
		t.Fail()
	}
}

func TestMSplit(t *testing.T) {
	msplit := MSplit("A LN VNS JRND KLMNT")
	if msplit.Size() != 15 {
		t.Logf("msplit=%v", msplit)
		t.Fail()
	}
}
