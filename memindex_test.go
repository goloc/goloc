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
	paris18.Postcode = "75018"
	paris18.City = "Paris"
	paris18.Country = "France"

	street2 := new(Street)
	street2.Id = "2"
	street2.StreetName = "Rue du Square Carpeaux"
	street2.Zone = paris18
	memindex.Add(street2)

	paris12 := new(Zone)
	paris12.Postcode = "75012"
	paris12.City = "Paris"
	paris12.Country = "France"

	street3 := new(Street)
	street3.Id = "3"
	street3.StreetName = "Rue de Lyon"
	street3.Zone = paris12
	memindex.Add(street3)

	sizeLoc := memindex.SizeLocalisation()
	if sizeLoc != 3 {
		t.Fail()
	}

	sizeIndex := memindex.SizeIndex()
	if sizeIndex != 16 {
		t.Fail()
	}

	results := memindex.Search("paris", 10, 600, 300)
	if len(results) != 3 {
		t.Fail()
	}

	results = memindex.Search("avenue", 10, 600, 300)
	if len(results) != 1 {
		t.Fail()
	}

	results = memindex.Search("carpe", 10, 600, 300)
	if len(results) != 1 {
		t.Fail()
	}

	results = memindex.Search("rue lyon paris", 10, 600, 300)
	if len(results) != 1 {
		t.Fail()
	}
}
