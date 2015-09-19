// Copyright 2015 Mathieu MAST. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.
package goloc

import (
	"bytes"

	"github.com/goloc/container"
)

type Street struct {
	Id          string
	EncodedName string
	CleanedName string
	StreetName  string
	Zone        *Zone
	Point
	NumberedPoints *container.LinkedList
}

func (s *Street) GetId() string {
	return s.Id
}

func (s *Street) GetEncodedName() string {
	return s.EncodedName
}

func (s *Street) SetEncodedName(encodedName string) {
	s.EncodedName = encodedName
}

func (s *Street) GetCleanedName() string {
	return s.CleanedName
}

func (s *Street) SetCleanedName(cleanedName string) {
	s.CleanedName = cleanedName
}

func (s *Street) GetName() string {
	b := bytes.NewBufferString("")
	if len(s.StreetName) > 0 {
		b.WriteString(s.StreetName)
	}
	if s.Zone != nil {
		b.WriteRune(' ')
		b.WriteString(s.Zone.GetName())
	}
	return b.String()
}

func (s *Street) AddNumberedPoint(numberedPoint *StreetNumberedPoint) {
	s.NumberedPoints.Add(numberedPoint)
}

func (s *Street) GetNumberedPoints() container.Container {
	return s.NumberedPoints
}

func (s *Street) GetType() string {
	return "street"
}

func (s *Street) GetLat() float32 {
	return s.Lat
}

func (s *Street) GetLon() float32 {
	return s.Lon
}

func NewStreet(id string, streetName string, zone *Zone, lat float32, lon float32) *Street {
	s := new(Street)
	s.NumberedPoints = container.NewLinkedList()
	s.Id = id
	s.StreetName = streetName
	s.Zone = zone
	s.Lat = lat
	s.Lon = lon
	return s
}

type StreetNumberedPoint struct {
	Number string
	Point
}

func (np *StreetNumberedPoint) GetNumber() string {
	return np.Number
}

func (np *StreetNumberedPoint) GetLat() float32 {
	return np.Lat
}

func (np *StreetNumberedPoint) GetLon() float32 {
	return np.Lon
}

func NewStreetNumberedPoint(number string, lat float32, lon float32) *StreetNumberedPoint {
	np := new(StreetNumberedPoint)
	np.Number = number
	np.Lat = lat
	np.Lon = lon
	return np
}
