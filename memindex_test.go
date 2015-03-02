// Copyright 2015 Mathieu MAST. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.
package goloc

import (
	"fmt"
	"runtime"
	"testing"
	"time"
)

func TestMemindex(t *testing.T) {
	memindex := NewMemindex()

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

	results := memindex.Search("paris", 10, nil)
	if results.Size != 6 {
		t.Fail()
	}

	results = memindex.Search("avenue", 10, nil)
	if results.Size != 1 {
		t.Fail()
	}

	results = memindex.Search("carpe", 10, nil)
	if results.Size != 1 {
		t.Fail()
	}
	if results.Head.Element.(*Result).Name != "Rue du Square Carpeaux 75018 Paris France" {
		t.Fail()
	}

	results = memindex.Search("10 carpe", 10, nil)
	if results.Size != 1 {
		t.Fail()
	}
	if results.Head.Element.(*Result).Name != "Rue du Square Carpeaux 75018 Paris France" {
		t.Fail()
	}

	results = memindex.Search("rue lyon paris", 10, nil)
	if results.Size != 1 {
		t.Fail()
	}

	results = memindex.Search("lyon", 10, nil)
	if results.Size != 2 {
		t.Fail()
	}

	// Search only poi:Gare
	results = memindex.Search("lyon", 10, func(result *Result) int {
		if result.Type == "poi:Gare" {
			return DefaultScorer(result)
		} else {
			return 0
		}
	})
	if results.Size != 1 {
		t.Fail()
	}
}

func TestPerfMemindex(t *testing.T) {
	runtime.GOMAXPROCS(16)
	memindex := NewMemindex()

	for i := 1; i <= 9; i++ {
		zone := new(Zone)
		zone.Id = fmt.Sprintf("Z%v", i)
		zone.Postcode = fmt.Sprintf("6900%v", i)
		zone.City = "Lyon"
		zone.Country = "France"
		memindex.Add(zone)
	}

	for i := 1; i < 10000; i++ {
		for j := 1; j <= 9; j++ {
			street := NewStreet(fmt.Sprintf("S%v%v", i, j))
			street.NumberedPoints["8"] = NewStreetNumberedPoint("8")
			street.NumberedPoints["9"] = NewStreetNumberedPoint("9")
			street.NumberedPoints["10"] = NewStreetNumberedPoint("10")
			street.StreetName = fmt.Sprintf("Rue du numéro %v%v", i, j)
			street.Zone = memindex.Get(fmt.Sprintf("Z%v", j)).(*Zone)
			memindex.Add(street)
		}
	}

	t0 := time.Now()
	results := memindex.Search("9 rue numero lyon", 10, nil)
	if results.Size != 10 {
		t.Fail()
	}
	t1 := time.Now()
	fmt.Printf("The search took %v to run.\n", t1.Sub(t0))
}
