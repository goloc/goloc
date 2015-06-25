// Copyright 2015 Mathieu MAST. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.
package goloc

type Bound struct {
	PointMin Point
	PointMax Point
}

func NewBound() *Bound {
	b := new(Bound)
	return b
}
