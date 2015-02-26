// Copyright 2015 Mathieu MAST. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.
package goloc

import (
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
	if target != name {
		t.Fail()
	}
}
