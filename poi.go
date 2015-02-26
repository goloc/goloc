// Copyright 2015 Mathieu MAST. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.
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
	return b.String()
}

func (p *Poi) GetType() string {
	return "poi:" + p.PoiType
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
