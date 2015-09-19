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
	// Add new locations
	Add(...Location)
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
	// Get encoded stop words
	GetEncodedStopWords() container.Container
}

// Location definition.
type Location interface {
	// Get unique identifier.
	GetId() string
	// Get cleaned name
	GetCleanedName() string
	// Set cleaned name
	SetCleanedName(string)
	// Get phonetic encoded name
	GetEncodedName() string
	// Set phonetic encoded name
	SetEncodedName(string)
	// Get name to indexing it.
	GetName() string
	// Get type (street, poi, zone...).
	GetType() string
	// Get latitude
	GetLat() float32
	// Get lLongitude
	GetLon() float32
}

// Numbered definition.
type NumberedPointBag interface {
	// Get number points
	GetNumberedPoints() container.Container
}

// Numbered point definition.
type NumberedPoint interface {
	// Get string number
	GetNumber() string
	// Get latitude
	GetLat() float32
	// Get longitude
	GetLon() float32
}
