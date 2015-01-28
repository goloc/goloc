package goloc

import ()

type Point struct {
	Lat float64
	Lon float64
}

/* compute distance in km */
func (p Point) Distance(target Point) float64 {
	return 0
}

func NewPoint() *Point {
	p := new(Point)
	return p
}
