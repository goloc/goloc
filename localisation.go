package goloc

import ()

type Localisation interface {
	GetId() string
	GetName() string
	GetType() string
	GetPoint() Point
}

type Point struct {
	lat float32
	lon float32
}

type Address struct {
	num   string
	point Point
}

type Street struct {
	id         string
	streetName string
	postcode   string
	settlement string
	city       string
	region     string
	country    string
	point      Point
}

func (l Street) GetId() string {
	return l.id
}
func (l Street) GetName() string {
	var name = l.streetName
	if len(l.postcode) > 0 {
		name += " " + l.postcode
	}
	if len(l.settlement) > 0 {
		name += " " + l.settlement
	}
	if len(l.region) > 0 {
		name += " " + l.region
	}
	if len(l.country) > 0 {
		name += " " + l.country
	}
	return name
}
func (l Street) GetType() string {
	return "street"
}
func (l Street) GetPoint() Point {
	return l.point
}

func clear() {

}

func index(l *Localisation) {

}
