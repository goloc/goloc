package goloc

import (
	"fmt"
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
	fmt.Printf("size localisation %d\n", sizeLoc)
	if sizeLoc != 2 {
		t.Fail()
	}

	sizeIndex := memindex.SizeIndex()
	fmt.Printf("size index %d\n", sizeIndex)
	if sizeIndex != 15 {
		t.Fail()
	}

	mapResult := memindex.Search("paris")
	if len(*mapResult) != 2 {
		t.Fail()
	}
	for _, result := range *mapResult {
		name := result.loc.GetName()
		fmt.Printf("%t %d\n", name, result.score)
	}

	mapResult = memindex.Search("avenue")
	if len(*mapResult) != 1 {
		t.Fail()
	}
	for _, result := range *mapResult {
		name := result.loc.GetName()
		fmt.Printf("%t %d\n", name, result.score)
	}

	mapResult = memindex.Search("carpo")
	if len(*mapResult) != 1 {
		t.Fail()
	}
	for _, result := range *mapResult {
		name := result.loc.GetName()
		fmt.Printf("%t %d\n", name, result.score)
	}

	memindex.SaveInFile("testSaveMemindex.gob")

	memindex2 := NewMemindexFromFile("testSaveMemindex.gob")

	sizeLoc2 := memindex2.SizeLocalisation()
	fmt.Printf("size localisation %d\n", sizeLoc2)
	if sizeLoc2 != 2 {
		t.Fail()
	}

	sizeIndex2 := memindex2.SizeIndex()
	fmt.Printf("size index %d\n", sizeIndex2)
	if sizeIndex2 != 15 {
		t.Fail()
	}
}
