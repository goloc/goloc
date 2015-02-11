package goloc

import ()

type Result struct {
	Score    int
	Search   string
	Id       string
	Priority uint8
	Name     string
	Lat      float32
	Lon      float32
}
