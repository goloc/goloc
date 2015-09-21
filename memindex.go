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
	sniffer          Sniffer
	Locations        *container.Map
	Keys             *container.Map
	StopWords        *container.Set
	EncodedStopWords *container.Set
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
	fmt.Printf("%v Locations\n", mi.Locations.Size())
	fmt.Printf("%v Keys\n", mi.Keys.Size())
	fmt.Printf("%v Stop words\n", mi.StopWords.Size())
	fmt.Printf("%v Encoded stop words\n", mi.EncodedStopWords.Size())

	return mi
}

func (mi *Memindex) Search(parameters Parameters) (container.Container, error) {
	return mi.sniffer.Search(parameters)
}

func (mi *Memindex) SaveInFile(filename string) {
	fmt.Printf("save %v\n", filename)
	fmt.Printf("%v Locations\n", mi.Locations.Size())
	fmt.Printf("%v Keys\n", mi.Keys.Size())
	fmt.Printf("%v Stop words\n", mi.StopWords.Size())
	fmt.Printf("%v Encoded stop words\n", mi.EncodedStopWords.Size())

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
	cleanedName := Clean(name, mi.GetStopWords())
	loc.SetCleanedName(cleanedName)
	encodedName := Partialphone(cleanedName)
	loc.SetEncodedName(encodedName)
	mkeys := MSplit(encodedName)
	id := loc.GetId()
	mi.Locations.Add(&container.KeyValue{id, loc})
	mkeys.Visit(func(element interface{}, i int) {
		var ids *container.LinkedList
		key := element.(string)
		v, err := mi.Keys.Get(key)
		if err == nil && v != nil {
			ids = v.(*container.KeyValue).Value.(*container.LinkedList)
		} else {
			ids = container.NewLinkedList()
			mi.Keys.Add(&container.KeyValue{key, ids})
		}
		ids.Add(id)
	})
}

func (mi *Memindex) Clear() {
	mi.Locations = container.NewMap()
	mi.Keys = container.NewMap()
	mi.StopWords = container.NewSet()
	mi.EncodedStopWords = container.NewSet()
}

func (mi *Memindex) Get(id string) Location {
	v, err := mi.Locations.Get(id)
	if err == nil && v != nil {
		return v.(*container.KeyValue).Value.(Location)
	} else {
		return nil
	}
}

func (mi *Memindex) GetNbIds(key string) int {
	v, err := mi.Keys.Get(key)
	if err == nil && v != nil {
		return v.(*container.KeyValue).Value.(container.Container).Size()
	} else {
		return 0
	}
}

func (mi *Memindex) GetIds(key string) container.Container {
	v, err := mi.Keys.Get(key)
	if err == nil && v != nil {
		return v.(*container.KeyValue).Value.(container.Container)
	} else {
		return container.NewLinkedList()
	}
}

func (mi *Memindex) AddStopWord(words ...string) {
	for _, word := range words {
		w := UpperUnaccentUnpunctString(word)
		mi.StopWords.Add(w)
		mi.EncodedStopWords.Add(Partialphone(w))
	}
}

func (mi *Memindex) GetStopWords() container.Container {
	return mi.StopWords
}

func (mi *Memindex) GetEncodedStopWords() container.Container {
	return mi.EncodedStopWords
}
