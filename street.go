package goloc

import (
	"bytes"
)

type Street struct {
	Id            string
	StreetName    string
	Zone          *Zone
	Lat           float64
	Lon           float64
	LinkedAddress *LinkedElement
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

func (s *Street) GetType() string {
	return "street"
}

func (z *Street) GetPriority() uint8 {
	return 2
}

func (s *Street) GetLat() float64 {
	return s.Lat
}

func (s *Street) GetLon() float64 {
	return s.Lon
}

func NewStreet() *Street {
	s := new(Street)
	return s
}
