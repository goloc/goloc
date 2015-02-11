package goloc

import ()

type Point struct {
	Lat float32
	Lon float32
}

func NewPoint() *Point {
	p := new(Point)
	return p
}
