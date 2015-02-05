package goloc

import ()

type Index interface {
	Add(loc Localisation)
	Clear()
	SizeLocalisation() int
	SizeIndex() int
	Get(string) Localisation
	Remove(string)
	Search(string, int) *LinkedList
}
