package goloc

import ()

type Address struct {
	Num string
	Lat float64
	Lon float64
}

func NewAddress() *Address {
	a := new(Address)
	return a
}
