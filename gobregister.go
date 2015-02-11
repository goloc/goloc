package goloc

import (
	"encoding/gob"
)

func GobRegister() {
	gob.RegisterName("goloc.Street", &Street{})
	gob.RegisterName("goloc.Zone", &Zone{})
	gob.RegisterName("goloc.Poi", &Poi{})
	gob.RegisterName("goloc.Point", &Point{})
}
