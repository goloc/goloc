// Copyright 2015 Mathieu MAST. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.
package goloc

import (
	"bytes"
)

type Street struct {
	Id         string
	StreetName string
	Zone       *Zone
	Point
	NumberedPoints map[string]*StreetNumberedPoint
}

func (s *Street) GetId() string {
	return s.Id
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

func (s *Street) GetNumberedPoint(search string) NumberedPoint {
	if len(s.NumberedPoints) > 0 {
		strs := Split(search)
		for _, str := range strs {
			if number, ok := s.NumberedPoints[str]; ok == true {
				return number
			}
		}
	}
	return nil
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

func NewStreet(id string) *Street {
	s := new(Street)
	s.NumberedPoints = make(map[string]*StreetNumberedPoint)
	s.Id = id
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

func NewStreetNumberedPoint(number string) *StreetNumberedPoint {
	np := new(StreetNumberedPoint)
	np.Number = number
	return np
}
