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
	locations map[string]Location
	keys      map[string]container.Container
	internal
}

func NewMemindex() *Memindex {
	mi := new(Memindex)
	mi.tolerance = defaultTolerance
	mi.keyLimit = defaultKeyLimit
	mi.locLimit = defaultLocLimit
	mi.get = mi.Get
	mi.internal.getIds = mi.getIds
	mi.internal.addLocationAndKeys = mi.addLocationAndKeys
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
	dataDecoder.Decode(&mi.keys)
	dataDecoder.Decode(&mi.locations)
	fmt.Printf("%v locations\n", mi.SizeLocation())
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
	if err = encoder.Encode(mi.keys); err != nil {
		panic(err)
	}
	if err = encoder.Encode(mi.locations); err != nil {
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
	return len(mi.locations)
}

func (mi *Memindex) SizeIndex() int {
	return len(mi.keys)
}

func (mi *Memindex) Clear() {
	mi.locations = make(map[string]Location)
	mi.keys = make(map[string]container.Container)
}

func (mi *Memindex) Get(id string) Location {
	loc := mi.locations[id]
	return loc
}

func (mi *Memindex) Search(search string, number int, filter Filter) container.Container {
	return mi.search(search, number, filter)
}

func (mi *Memindex) getIds(key string) container.Container {
	ids := mi.keys[key]
	return ids
}

func (mi *Memindex) addLocationAndKeys(loc Location, keys []string) {
	var ids container.Container
	var k string
	id := loc.GetId()
	mi.locations[id] = loc
	for _, k = range keys {
		ids = mi.keys[k]
		if ids == nil {
			ids = container.NewLinkedList()
			mi.keys[k] = ids
		}
		ids.Add(id)
	}
}
