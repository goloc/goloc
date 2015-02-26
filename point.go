// Copyright 2015 Mathieu MAST. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.
package goloc

import ()

type Point struct {
	Lat float32
	Lon float32
}

func NewPoint() *Point {
	p := new(Point)
	return p
}
