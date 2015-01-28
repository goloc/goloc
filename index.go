package goloc

import ()

type Index interface {
	Add(loc Localisation)
	Clear()
	SizeLocalisation() int
	SizeIndex() int
	Get(id string) Localisation
	Remove(id string)
	Search(str string) *map[string]*Result
}
