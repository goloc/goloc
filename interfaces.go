// Copyright 2015 Mathieu MAST. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.
package goloc

import "github.com/goloc/container"

type Parameters interface {
	Get(string) interface{}
	Set(string, interface{})
}

type Sniffer interface {
	// Search
	Search(Parameters) (container.Container, error)
}

type Index interface {
	Sniffer
	// Add new location
	Add(Location)
	// Get location
	Get(string) Location
	// Get ids number for key
	GetNbIds(string) int
	// // Get ids for key
	GetIds(string) container.Container
	// Add new stop word
	AddStopWord(...string)
	// Get stop words
	GetStopWords() container.Container
}

// Location definition.
type Location interface {
	// Unique identifier.
	GetId() string
	// Name to indexing it.
	GetName() string
	// Type (street, poi, zone...).
	GetType() string
	// Latitude
	GetLat() float32
	// Longitude
	GetLon() float32
}

// Numbered definition.
type NumberedPointBag interface {
	// Get number for search input parameter
	GetNumberedPoint(string) NumberedPoint
}

// Numbered point definition.
type NumberedPoint interface {
	// String number
	GetNumber() string
	// Latitude
	GetLat() float32
	// Longitude
	GetLon() float32
}
