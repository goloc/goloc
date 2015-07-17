// Copyright 2015 Mathieu MAST. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.
package goloc

import (
	"encoding/gob"
	"fmt"
	"os"
	"runtime"
	"strings"

	"github.com/goloc/container"
)

type Memindex struct {
	sniffer   Sniffer
	Locations map[string]Location
	Keys      map[string]*container.LinkedList
	StopWords *container.LinkedList
}

func NewMemindex() *Memindex {
	mi := new(Memindex)
	mi.sniffer = NewConcurrentSniffer(mi)
	mi.Clear()
	GobRegister()
	runtime.GOMAXPROCS(runtime.NumCPU())
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

	mi := new(Memindex)
	mi.Clear()
	GobRegister()
	runtime.GOMAXPROCS(runtime.NumCPU())
	dataDecoder := gob.NewDecoder(file)
	dataDecoder.Decode(&mi)
	mi.sniffer = NewConcurrentSniffer(mi)
	fmt.Printf("%v Locations\n", len(mi.Locations))
	fmt.Printf("%v Keys\n", len(mi.Keys))
	fmt.Printf("%v Stop words\n", mi.StopWords.Size())

	return mi
}

func (mi *Memindex) Search(parameters Parameters) (container.Container, error) {
	return mi.sniffer.Search(parameters)
}

func (mi *Memindex) SaveInFile(filename string) {
	fmt.Printf("save %v\n", filename)
	fmt.Printf("%v Locations\n", len(mi.Locations))
	fmt.Printf("%v Keys\n", len(mi.Keys))
	fmt.Printf("%v Stop words\n", mi.StopWords.Size())

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

func (mi *Memindex) Add(locs ...Location) {
	for _, loc := range locs {
		mi.addOne(loc)
	}
}

func (mi *Memindex) addOne(loc Location) {
	name := " " + UpperUnaccentUnpunctString(loc.GetName()) + " "
	if mi.StopWords != nil {
		mi.StopWords.Visit(func(element interface{}, i int) {
			word := " " + element.(string) + " "
			name = strings.Join(strings.Split(name, word), " ")
		})
	}
	encodedName := Partialphone(name)
	loc.SetEncodedName(encodedName)
	mkeys := MSplit(encodedName)
	id := loc.GetId()
	mi.Locations[id] = loc
	mkeys.Visit(func(element interface{}, i int) {
		key := element.(string)
		ids := mi.Keys[key]
		if ids == nil {
			ids = container.NewLinkedList()
			mi.Keys[key] = ids
		}
		ids.Add(id)
	})
}

func (mi *Memindex) Clear() {
	mi.Locations = make(map[string]Location)
	mi.Keys = make(map[string]*container.LinkedList)
	mi.StopWords = container.NewLinkedList()
}

func (mi *Memindex) Get(id string) Location {
	loc := mi.Locations[id]
	return loc
}

func (mi *Memindex) GetNbIds(key string) int {
	ids := mi.Keys[key]
	if ids != nil {
		return ids.Size()
	} else {
		return 0
	}
}

func (mi *Memindex) GetIds(key string) container.Container {
	ids := mi.Keys[key]
	if ids == nil {
		return container.NewLinkedList()
	}
	return ids
}

func (mi *Memindex) AddStopWord(words ...string) {
	for _, word := range words {
		mi.StopWords.Add(UpperUnaccentUnpunctString(word))
	}
}

func (mi *Memindex) GetStopWords() container.Container {
	return mi.StopWords
}
