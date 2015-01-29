package goloc

import (
	"testing"
)

func TestMemindex(t *testing.T) {
	memindex := NewMemindex()

	paris8 := new(Zone)
	paris8.Postcode = "75008"
	paris8.City = "Paris"
	paris8.Country = "France"

	street1 := new(Street)
	street1.Id = "1"
	street1.StreetName = "Avenue des Champs-Élysées"
	street1.Zone = paris8
	memindex.Add(street1)

	paris18 := new(Zone)
	paris18.Postcode = "75008"
	paris18.City = "Paris"
	paris18.Country = "France"

	street2 := new(Street)
	street2.Id = "2"
	street2.StreetName = "Rue du Square Carpeaux"
	street2.Zone = paris18
	memindex.Add(street2)

	sizeLoc := memindex.SizeLocalisation()
	if sizeLoc != 2 {
		t.Fail()
	}

	sizeIndex := memindex.SizeIndex()
	if sizeIndex != 15 {
		t.Fail()
	}

	memindex.SaveInFile("golocTest.gob")
}

func TestReloadAndSearch(t *testing.T) {
	memindex := NewMemindexFromFile("golocTest.gob")

	sizeLoc := memindex.SizeLocalisation()
	if sizeLoc != 2 {
		t.Fail()
	}

	sizeIndex := memindex.SizeIndex()
	if sizeIndex != 15 {
		t.Fail()
	}

	results := memindex.Search("paris", 10, 500, 300)
	if len(results) != 2 {
		t.Fail()
	}

	results = memindex.Search("avenue", 10, 500, 300)
	if len(results) != 1 {
		t.Fail()
	}

	results = memindex.Search("carpe", 10, 500, 300)
	if len(results) != 1 {
		t.Fail()
	}
}
