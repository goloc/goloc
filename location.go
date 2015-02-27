// Copyright 2015 Mathieu MAST. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.
package goloc

import ()

// Location definition.
type Location interface {
	// Unique identifier.
	GetId() string
	// Name to indexing it.
	GetName() string
	// Name for result search presentation.
	// In some cases, it's equivalent to GetName().
	GetResultName(string) string
	// Type (street, poi, zone...).
	GetType() string
	// Latitude
	GetLat() float32
	// Longitude
	GetLon() float32
}
