package goloc

import (
	"encoding/gob"
	"fmt"
	"os"
	//	"strconv"
)

type Memindex struct {
	Localisations map[string]Localisation
	Keys          map[string]*LinkedList
}

func NewMemindex() *Memindex {
	mi := new(Memindex)
	mi.Clear()
	GobRegister()
	return mi
}

func NewMemindexFromFile(filename string) *Memindex {
	fmt.Printf("load %v\n", filename)
	mi := NewMemindex()
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
		return mi
	}
	defer file.Close()

	dataDecoder := gob.NewDecoder(file)
	dataDecoder.Decode(&mi.Keys)
	dataDecoder.Decode(&mi.Localisations)
	fmt.Printf("%v localisations\n", mi.SizeLocalisation())
	fmt.Printf("%v keys\n", mi.SizeIndex())

	return mi
}

func (mi *Memindex) SaveInFile(filename string) {
	fmt.Printf("save %v\n", filename)
	file, err := os.Create(filename)
	if err != nil {
		panic(err)
		return
	}
	encoder := gob.NewEncoder(file)
	if err = encoder.Encode(mi.Keys); err != nil {
		panic(err)
	}
	if err = encoder.Encode(mi.Localisations); err != nil {
		panic(err)
	}
	if err = file.Close(); err != nil {
		panic(err)
	}
}
func (mi *Memindex) Add(loc Localisation) {
	internalAdd(mi, loc)
}

func (mi *Memindex) SizeLocalisation() int {
	return len(mi.Localisations)
}

func (mi *Memindex) SizeIndex() int {
	return len(mi.Keys)
}

func (mi *Memindex) Clear() {
	mi.Localisations = make(map[string]Localisation)
	mi.Keys = make(map[string]*LinkedList)
}

func (mi *Memindex) Get(id string) Localisation {
	loc := mi.Localisations[id]
	return loc
}

func (mi *Memindex) getInternalIdsForKey(key string) *LinkedList {
	ids := mi.Keys[key]
	return ids
}

func (mi *Memindex) addInternalLocalisation(loc Localisation, keys map[string]bool) {
	var ids *LinkedList
	var k string
	id := loc.GetId()
	mi.Localisations[id] = loc
	for k, _ = range keys {
		ids = mi.Keys[k]
		if ids == nil {
			ids = NewLinkedList()
			mi.Keys[k] = ids
		}
		ids.AddLast(id)
	}
}

func (mi *Memindex) Search(search string, number int, scorer Scorer) *LinkedList {
	return internalSearch(mi, search, number, scorer)
}
