// Copyright 2015 Mathieu MAST. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.
package goloc

import (
	"reflect"
	"testing"
)

func TestMemindex(t *testing.T) {
	memindex := NewMemindex()
	memindex.tolerance = 1
	memindex.AddStopWord("D", "DE", "DU", "DES", "L", "LE", "LA", "LES")
	memindex.AddStopWord("RUE", "ROUTE", "ALLEE", "PLACE", "CHEMIN", "IMPASSE", "AVENUE", "BOULEVARD")

	paris8 := NewZone("Z1")
	paris8.Postcode = "75008"
	paris8.City = "Paris"
	paris8.Country = "France"
	memindex.Add(paris8)

	street1 := NewStreet("S1")
	street1.StreetName = "Avenue des Champs-Élysées"
	street1.Zone = paris8
	memindex.Add(street1)

	paris18 := NewZone("Z2")
	paris18.Postcode = "75018"
	paris18.City = "Paris"
	paris18.Country = "France"
	memindex.Add(paris18)

	street2 := NewStreet("S2")
	street2.StreetName = "Rue du Square Carpeaux"
	street2.Zone = paris18
	street2.NumberedPoints["8"] = NewStreetNumberedPoint("8")
	street2.NumberedPoints["9"] = NewStreetNumberedPoint("9")
	street2.NumberedPoints["10"] = NewStreetNumberedPoint("10")
	memindex.Add(street2)

	paris12 := NewZone("Z3")
	paris12.Postcode = "75012"
	paris12.City = "Paris"
	paris12.Country = "France"
	memindex.Add(paris12)

	street3 := NewStreet("S3")
	street3.StreetName = "Rue de Lyon"
	street3.Zone = paris12
	memindex.Add(street3)

	poi1 := NewPoi("P1")
	poi1.PoiName = "Gare de Lyon"
	poi1.PoiType = "Gare"
	poi1.Zone = paris12
	memindex.Add(poi1)

	results := memindex.Search("rue", 10, nil)
	if results.GetSize() != 0 {
		t.Logf("Result number should be %v but was %v.", 0,
			results.GetSize())
		t.Fail()
	}

	results = memindex.Search("paris", 10, nil)
	if results.GetSize() != 6 {
		t.Logf("Result number should be %v but was %v.", 6,
			results.GetSize())
		t.Fail()
	}

	results = memindex.Search("avenue champs", 10, nil)
	if results.GetSize() != 1 {
		t.Logf("Result number should be %v but was %v.", 1,
			results.GetSize())
		t.Fail()
	}

	results = memindex.Search("carpe", 10, nil)
	if results.GetSize() != 1 {
		t.Logf("Result number should be %v but was %v.", 1,
			results.GetSize())
		t.Fail()
	} else if results.ToArrayOfType(reflect.TypeOf(new(Result))).([]*Result)[0].Name != "Rue du Square Carpeaux 75018 Paris France" {
		t.Fail()
	}

	results = memindex.Search("10 carpe", 10, nil)
	if results.GetSize() != 1 {
		t.Logf("Result number should be %v but was %v.", 1,
			results.GetSize())
		t.Fail()
	} else if results.ToArrayOfType(reflect.TypeOf(new(Result))).([]*Result)[0].Name != "Rue du Square Carpeaux 75018 Paris France" {
		t.Fail()
	}

	results = memindex.Search("rue lyon paris", 10, nil)
	if results.GetSize() != 2 {
		t.Logf("Result number should be %v but was %v.", 2,
			results.GetSize())
		t.Fail()
	}

	results = memindex.Search("lyon", 10, nil)
	if results.GetSize() != 2 {
		t.Logf("Result number should be %v but was %v.", 2,
			results.GetSize())
		t.Fail()
	}

	// Search only poi:Gare
	results = memindex.Search("lyon", 10, func(result *Result) bool {
		if result.Type == "poi:Gare" {
			return true
		} else {
			return false
		}
	})
	if results.GetSize() != 1 {
		t.Logf("Result number should be %v but was %v.", 1,
			results.GetSize())
		t.Fail()
	}
}
