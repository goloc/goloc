// Copyright 2015 Mathieu MAST. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.
package goloc

import (
	"bytes"
)

type Zone struct {
	Id         string
	Postcode   string
	Settlement string
	City       string
	Region     string
	Country    string
	PointMin   Point
	PointMax   Point
}

func (z *Zone) GetId() string {
	return z.Id
}

func (z *Zone) GetName() string {
	b := bytes.NewBufferString("")
	if len(z.Postcode) > 0 {
		if b.Len() > 0 {
			b.WriteRune(' ')
		}
		b.WriteString(z.Postcode)
	}
	if len(z.Settlement) > 0 {
		if b.Len() > 0 {
			b.WriteRune(' ')
		}
		b.WriteString(z.Settlement)
	}
	if len(z.City) > 0 {
		if b.Len() > 0 {
			b.WriteRune(' ')
		}
		b.WriteString(z.City)
	}
	if len(z.Region) > 0 {
		if b.Len() > 0 {
			b.WriteRune(' ')
		}
		b.WriteString(z.Region)
	}
	if len(z.Country) > 0 {
		if b.Len() > 0 {
			b.WriteRune(' ')
		}
		b.WriteString(z.Country)
	}
	return b.String()
}

func (z *Zone) GetType() string {
	return "zone"
}

func (z *Zone) GetLat() float32 {
	return (z.PointMin.Lat + z.PointMax.Lat) / 2
}

func (z *Zone) GetLon() float32 {
	return (z.PointMin.Lon + z.PointMax.Lon) / 2
}

func NewZone(id string) *Zone {
	z := new(Zone)
	z.Id = id
	return z
}
