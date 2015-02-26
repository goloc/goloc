// Copyright 2015 Mathieu MAST. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.
package goloc

import ()

type Localisation interface {
	GetId() string
	GetName() string
	GetType() string
	GetLat() float32
	GetLon() float32
}
