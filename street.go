package goloc

import (
	"bytes"
)

type Street struct {
	Id         string
	StreetName string
	Zone       *Zone
	Point      *Point
	Addresses  []*Address
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
func (s *Street) GetPoint() *Point {
	return s.Point
}
func NewStreet() *Street {
	s := new(Street)
	return s
}
