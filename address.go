package goloc

import ()

type Address struct {
	Num   string
	Point *Point
}

func NewAddress() *Address {
	a := new(Address)
	a.Point = NewPoint()
	return a
}
