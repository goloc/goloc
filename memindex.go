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
	Keys      map[string]*container.LinkedList
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
	dataDecoder.Decode(&mi.Locations)
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
	if err = encoder.Encode(mi.Keys); err != nil {
		panic(err)
	}
	if err = encoder.Encode(mi.Locations); err != nil {
		panic(err)
	}
	if err = file.Close(); err != nil {
		panic(err)
	}
}
func (mi *Memindex) Add(loc Location) {
	internalAdd(mi, loc, mi.addLocationAndKeys)
}

func (mi *Memindex) SizeLocation() int {
	return len(mi.Locations)
}

func (mi *Memindex) SizeIndex() int {
	return len(mi.Keys)
}

func (mi *Memindex) Clear() {
	mi.Locations = make(map[string]Location)
	mi.Keys = make(map[string]*container.LinkedList)
}

func (mi *Memindex) Get(id string) Location {
	loc := mi.Locations[id]
	return loc
}

func (mi *Memindex) Search(search string, number int, scorer Scorer) *container.LimitedBinaryTree {
	return internalSearch(mi, search, number, scorer, mi.getIds)
}

func (mi *Memindex) getIds(key string) *container.LinkedList {
	ids := mi.Keys[key]
	return ids
}

func (mi *Memindex) addLocationAndKeys(loc Location, keys []string) {
	var ids *container.LinkedList
	var k string
	id := loc.GetId()
	mi.Locations[id] = loc
	for _, k = range keys {
		ids = mi.Keys[k]
		if ids == nil {
			ids = container.NewLinkedList()
			mi.Keys[k] = ids
		}
		ids.Push(id)
	}
}
