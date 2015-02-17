package goloc

import ()

type Localisation interface {
	GetId() string
	GetName() string
	GetType() string
	GetLat() float32
	GetLon() float32
}
