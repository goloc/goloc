package goloc

import (
	"bytes"
)

type Poi struct {
	Id      string
	PoiName string
	PoiType string
	Zone    *Zone
	Point
}

func (p *Poi) GetId() string {
	return p.Id
}

func (p *Poi) GetName() string {
	b := bytes.NewBufferString("")
	if len(p.PoiName) > 0 {
		b.WriteString(p.PoiName)
	}
	if p.Zone != nil {
		b.WriteRune(' ')
		b.WriteString(p.Zone.GetName())
	}
	return b.String()
}

func (p *Poi) GetType() string {
	return "zone"
}

func (p *Poi) GetPriority() uint8 {
	return 0
}

func (p *Poi) GetLat() float32 {
	return p.Lat
}

func (p *Poi) GetLon() float32 {
	return p.Lon
}

func NewPoi() *Poi {
	p := new(Poi)
	return p
}
