// Copyright 2015 Mathieu MAST. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.
package goloc

import (
	"encoding/gob"

	"github.com/goloc/container"
)

func GobRegister() {
	gob.RegisterName("g.S", &Street{})
	gob.RegisterName("g.SNP", &StreetNumberedPoint{})
	gob.RegisterName("g.Z", &Zone{})
	gob.RegisterName("g.P", &Poi{})
	gob.RegisterName("c.m", &container.Map{})
	gob.RegisterName("c.s", &container.Set{})
	gob.RegisterName("c.ll", &container.LinkedList{})
}
