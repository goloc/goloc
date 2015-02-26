// Copyright 2015 Mathieu MAST. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.
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
