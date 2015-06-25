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
	indexTest(memindex, t)
}

func indexTest(index Index, t *testing.T) {
	index.AddStopWord("D", "DE", "DU", "DES", "L", "LE", "LA", "LES")
	index.AddStopWord("RUE", "ROUTE", "ALLEE", "PLACE", "CHEMIN", "IMPASSE", "AVENUE", "BOULEVARD")

	paris8 := NewZone("Z1", "75008", "", "Paris", "", "France", 0, 0, 0, 0)
	index.Add(paris8)

	street1 := NewStreet("S1", "Avenue des Champs-Élysées", paris8, 0, 0)
	index.Add(street1)

	paris18 := NewZone("Z2", "75018", "", "Paris", "", "France", 0, 0, 0, 0)
	index.Add(paris18)

	street2 := NewStreet("S2", "Rue du Square Carpeaux", paris18, 0, 0)
	street2.AddNumberedPoint(NewStreetNumberedPoint("8", 0, 0))
	street2.AddNumberedPoint(NewStreetNumberedPoint("9", 0, 0))
	street2.AddNumberedPoint(NewStreetNumberedPoint("10", 0, 0))
	index.Add(street2)

	paris12 := NewZone("Z3", "75012", "", "Paris", "", "France", 0, 0, 0, 0)
	index.Add(paris12)

	street3 := NewStreet("S3", "Rue de Lyon", paris12, 0, 0)
	index.Add(street3)

	poi1 := NewPoi("P1", "Gare de Lyon", "Gare", paris12, 0, 0)
	index.Add(poi1)

	results, err := index.Search(NewMapParameters(KeyValue{Key: "search", Value: "rue"}))
	if err != nil {
		t.Logf("%v\n", err)
		t.Fail()
	} else if results.GetSize() != 0 {
		t.Logf("Result number should be %v but was %v.", 0,
			results.GetSize())
		t.Fail()
	}

	results, err = index.Search(NewMapParameters(KeyValue{Key: "search", Value: "paris"}))
	if err != nil {
		t.Logf("%v\n", err)
		t.Fail()
	} else if results.GetSize() != 7 {
		t.Logf("Result number should be %v but was %v.", 7,
			results.GetSize())
		t.Fail()
	}

	results, err = index.Search(NewMapParameters(KeyValue{Key: "search", Value: "avenue champs"}))
	if err != nil {
		t.Logf("%v\n", err)
		t.Fail()
	} else if results.GetSize() != 1 {
		t.Logf("Result number should be %v but was %v.", 1,
			results.GetSize())
		t.Fail()
	}

	results, err = index.Search(NewMapParameters(KeyValue{Key: "search", Value: "carpe"}))
	if err != nil {
		t.Logf("%v\n", err)
		t.Fail()
	} else if results.GetSize() != 1 {
		t.Logf("Result number should be %v but was %v.", 1,
			results.GetSize())
		t.Fail()
	} else if results.ToArrayOfType(reflect.TypeOf(new(Result))).([]*Result)[0].Name != "Rue du Square Carpeaux 75018 Paris France" {
		t.Fail()
	}

	results, err = index.Search(NewMapParameters(KeyValue{Key: "search", Value: "10 carpe"}))
	if err != nil {
		t.Logf("%v\n", err)
		t.Fail()
	} else if results.GetSize() != 1 {
		t.Logf("Result number should be %v but was %v.", 1,
			results.GetSize())
		t.Fail()
	} else if results.ToArrayOfType(reflect.TypeOf(new(Result))).([]*Result)[0].Name != "Rue du Square Carpeaux 75018 Paris France" {
		t.Fail()
	}

	results, err = index.Search(NewMapParameters(KeyValue{Key: "search", Value: "rue lyon paris"}))
	if err != nil {
		t.Logf("%v\n", err)
		t.Fail()
	} else if results.GetSize() != 2 {
		t.Logf("Result number should be %v but was %v.", 2,
			results.GetSize())
		t.Fail()
	}

	results, err = index.Search(NewMapParameters(KeyValue{Key: "search", Value: "lyon"}))
	if err != nil {
		t.Logf("%v\n", err)
		t.Fail()
	} else if results.GetSize() != 2 {
		t.Logf("Result number should be %v but was %v.", 2,
			results.GetSize())
		t.Fail()
	}

	// Search only poi:Gare
	results, err = index.Search(NewMapParameters(KeyValue{Key: "search", Value: "lyon"}, KeyValue{Key: "filter", Value: func(result *Result) bool {
		if result.Type == "poi:Gare" {
			return true
		}
		return false
	}}))
	if err != nil {
		t.Logf("%v\n", err)
		t.Fail()
	} else if results.GetSize() != 1 {
		t.Logf("Result number should be %v but was %v.", 1,
			results.GetSize())
		t.Fail()
	}
}
