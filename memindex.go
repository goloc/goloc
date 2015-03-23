// Copyright 2015 Mathieu MAST. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.
package goloc

import (
	"encoding/gob"
	"fmt"
	"github.com/goloc/container"
	"os"
)

type Memindex struct {
	Locations map[string]Location
	Keys      map[string]container.Container
	StopWords container.Container
	internal
}

func NewMemindex() *Memindex {
	mi := new(Memindex)
	mi.tolerance = defaultTolerance
	mi.locLimit = defaultLocLimit
	mi.get = mi.Get
	mi.internal.getIds = mi.getIds
	mi.internal.addLocationAndKeys = mi.addLocationAndKeys
	mi.internal.getStopWords = mi.getStopWords
	mi.Clear()
	GobRegister()
	return mi
}

func NewMemindexFromFile(filename string) *Memindex {
	fmt.Printf("load %v\n", filename)
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
		return nil
	}
	defer file.Close()

	mi := NewMemindex()
	dataDecoder := gob.NewDecoder(file)
	dataDecoder.Decode(&mi)
	fmt.Printf("%v Locations\n", mi.SizeLocation())
	fmt.Printf("%v Keys\n", mi.SizeIndex())
	fmt.Printf("%v Stop words\n", mi.StopWords.GetSize())

	return mi
}

func (mi *Memindex) SaveInFile(filename string) {
	fmt.Printf("save %v\n", filename)
	fmt.Printf("%v Locations\n", mi.SizeLocation())
	fmt.Printf("%v Keys\n", mi.SizeIndex())
	fmt.Printf("%v Stop words\n", mi.StopWords.GetSize())

	file, err := os.Create(filename)
	if err != nil {
		panic(err)
		return
	}
	encoder := gob.NewEncoder(file)
	if err = encoder.Encode(mi); err != nil {
		panic(err)
	}
	if err = file.Close(); err != nil {
		panic(err)
	}
}

func (mi *Memindex) Add(loc Location) {
	mi.add(loc)
}

func (mi *Memindex) SizeLocation() int {
	return len(mi.Locations)
}

func (mi *Memindex) SizeIndex() int {
	return len(mi.Keys)
}

func (mi *Memindex) Clear() {
	mi.Locations = make(map[string]Location)
	mi.Keys = make(map[string]container.Container)
	mi.StopWords = container.NewLinkedList()
}

func (mi *Memindex) Get(id string) Location {
	loc := mi.Locations[id]
	return loc
}

func (mi *Memindex) Search(search string, number int, filter Filter) container.Container {
	return mi.search(search, number, filter)
}

func (mi *Memindex) AddStopWord(words ...string) {
	for _, word := range words {
		mi.StopWords.Add(UpperUnaccentUnpunctString(word))
	}
}

func (mi *Memindex) getIds(key string) container.Container {
	ids := mi.Keys[key]
	return ids
}

func (mi *Memindex) addLocationAndKeys(loc Location, Keys container.Container) {
	var ids container.Container
	id := loc.GetId()
	mi.Locations[id] = loc
	Keys.Visit(func(element interface{}, i int) {
		key := element.(string)
		ids = mi.Keys[key]
		if ids == nil {
			ids = container.NewLinkedList()
			mi.Keys[key] = ids
		}
		ids.Add(id)
	})
}

func (mi *Memindex) getStopWords() container.Container {
	return mi.StopWords
}
