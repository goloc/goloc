package goloc

import ()

type Localisation interface {
	GetId() string
	GetName() string
	GetType() string
	GetPriority() uint8
	GetLat() float32
	GetLon() float32
}
