// Copyright 2015 Mathieu MAST. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.
package goloc

import (
	"encoding/gob"

	"github.com/goloc/container"
)

func GobRegister() {
	gob.RegisterName("goloc.Street", &Street{})
	gob.RegisterName("goloc.Zone", &Zone{})
	gob.RegisterName("goloc.Poi", &Poi{})
	gob.RegisterName("goloc.Point", &Point{})
	gob.RegisterName("goloc.Bound", &Bound{})
	gob.RegisterName("goloc.ConcurrentSniffer", &ConcurrentSniffer{})
	gob.RegisterName("container.LinkedList", &container.LinkedList{})
	gob.RegisterName("container.BinaryTree", &container.BinaryTree{})
	gob.RegisterName("container.LimitedBinaryTree", &container.LimitedBinaryTree{})
}
