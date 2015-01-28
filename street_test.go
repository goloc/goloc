package goloc

import (
	"fmt"
	"testing"
)

func TestStreet(t *testing.T) {
	paris8 := new(Zone)
	paris8.Postcode = "75008"
	paris8.City = "Paris"
	paris8.Country = "France"

	street := new(Street)
	street.Id = "1"
	street.StreetName = "Avenue des Champs-Élysées"
	street.Zone = paris8
	name := street.GetName()
	target := "Avenue des Champs-Élysées 75008 Paris France"
	fmt.Printf("%t -> %t\n", name, target)
	if target != name {
		t.Fail()
	}
}
