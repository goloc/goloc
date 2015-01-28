package goloc

import ()

type Localisation interface {
	GetId() string
	GetName() string
	GetType() string
	GetPoint() *Point
}
