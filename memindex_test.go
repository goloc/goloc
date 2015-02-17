package goloc

import (
	"testing"
)

func TestMemindex(t *testing.T) {
	memindex := NewMemindex()

	paris8 := new(Zone)
	paris8.Id = "Z1"
	paris8.Postcode = "75008"
	paris8.City = "Paris"
	paris8.Country = "France"
	memindex.Add(paris8)

	street1 := new(Street)
	street1.Id = "S1"
	street1.StreetName = "Avenue des Champs-Élysées"
	street1.Zone = paris8
	memindex.Add(street1)

	paris18 := new(Zone)
	paris18.Id = "Z2"
	paris18.Postcode = "75018"
	paris18.City = "Paris"
	paris18.Country = "France"
	memindex.Add(paris18)

	street2 := new(Street)
	street2.Id = "S2"
	street2.StreetName = "Rue du Square Carpeaux"
	street2.Zone = paris18
	memindex.Add(street2)

	paris12 := new(Zone)
	paris12.Id = "Z3"
	paris12.Postcode = "75012"
	paris12.City = "Paris"
	paris12.Country = "France"
	memindex.Add(paris12)

	street3 := new(Street)
	street3.Id = "S3"
	street3.StreetName = "Rue de Lyon"
	street3.Zone = paris12
	memindex.Add(street3)

	poi1 := NewPoi()
	poi1.Id = "P1"
	poi1.PoiName = "Gare de Lyon"
	poi1.PoiType = "Gare"
	poi1.Zone = paris12
	memindex.Add(poi1)

	list := memindex.Search("paris", 10, nil)
	if list.Size != 6 {
		t.Fail()
	}

	list = memindex.Search("avenue", 10, nil)
	if list.Size != 1 {
		t.Fail()
	}

	list = memindex.Search("carpe", 10, nil)
	if list.Size != 1 {
		t.Fail()
	}

	list = memindex.Search("rue lyon paris", 10, nil)
	if list.Size != 1 {
		t.Fail()
	}

	list = memindex.Search("lyon", 10, nil)
	if list.Size != 2 {
		t.Fail()
	}

	// Search only poi:Gare
	list = memindex.Search("lyon", 10, func(result *Result) int {
		if result.Type == "poi:Gare" {
			return DefaultScorer(result)
		} else {
			return 0
		}
	})
	if list.Size != 1 {
		t.Fail()
	}
}
