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
	Point      *Point
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
func (z *Zone) GetPoint() *Point {
	return z.Point
}
func NewZone() *Zone {
	z := new(Zone)
	z.Point = NewPoint()
	return z
}
