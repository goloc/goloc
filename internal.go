// Copyright 2015 Mathieu MAST. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.
package goloc

import (
	"fmt"
	"github.com/goloc/container"
	"strings"
	"sync"
)

type internal struct {
	tolerance          int
	keyLimit           int
	locLimit           int
	get                func(string) Location
	getIds             func(string) container.Container
	addLocationAndKeys func(loc Location, keys container.Container)
	getStopWords       func() container.Container
}

func (inter *internal) add(loc Location) {
	name := " " + UpperUnaccentUnpunctString(loc.GetName()) + " "
	if inter.getStopWords() != nil {
		inter.getStopWords().Visit(func(element interface{}, i int) {
			word := " " + element.(string) + " "
			name = strings.Join(strings.Split(name, word), " ")
		})
	}
	mkeys := MSplit(Partialphone(name))
	inter.addLocationAndKeys(loc, mkeys)
}

func (inter *internal) search(search string, number int, filter Filter) container.Container {
	if filter == nil {
		filter = DefaultFilter
	}

	cleansearch := UpperUnaccentUnpunctString(" " + search + " ")
	if inter.getStopWords() != nil {
		inter.getStopWords().Visit(func(element interface{}, i int) {
			word := " " + element.(string) + " "
			cleansearch = strings.Join(strings.Split(cleansearch, word), " ")
		})
	}

	words := Split(cleansearch)
	mwords := MSplit(cleansearch)
	mkeys := MSplit(Partialphone(cleansearch))

	var waitgroup sync.WaitGroup
	var mutex sync.Mutex

	counter := container.NewCounter()
	mkeys.Visit(func(element interface{}, i int) {
		waitgroup.Add(1)
		go func(key string) {
			defer waitgroup.Done()
			ids := inter.getIds(key)
			if ids != nil && ids.GetSize() > 0 {
				ids.Visit(func(element interface{}, i int) {
					id := element.(string)
					counter.Incr(id)
				})
			}
		}(element.(string))
	})
	waitgroup.Wait()

	tmpResults := container.NewLimitedBinaryTree(CompareScoreResult, inter.locLimit, true)
	minKeyScore := int(counter.GetMax()) - inter.tolerance
	counter.Visit(func(id string, count uint32, i int) {
		if int(count) >= minKeyScore {
			waitgroup.Add(1)
			go func(id string) {
				defer waitgroup.Done()
				loc := inter.get(id)
				if loc != nil {
					result := new(Result)
					result.Id = id
					result.Search = search
					result.Name = loc.GetName()
					result.Lat = loc.GetLat()
					result.Lon = loc.GetLon()
					result.Type = loc.GetType()
					bag, ok := loc.(NumberedPointBag)
					if ok {
						numbered := bag.GetNumberedPoint(search)
						if numbered != nil {
							result.Number = numbered.GetNumber()
							result.Lat = numbered.GetLat()
							result.Lon = numbered.GetLon()
						}
					}
					if filter(result) {
						result.Score += QuickScore(mwords, UpperUnaccentUnpunctString(result.Name))
						if result.Score > 0 {
							mutex.Lock()
							tmpResults.Add(result)
							mutex.Unlock()
						}
					}
				}
			}(id)
		}
	})
	waitgroup.Wait()

	results := container.NewLimitedBinaryTree(CompareScoreResult, number, true)
	tmpResults.Visit(func(element interface{}, i int) {
		waitgroup.Add(1)
		go func(result *Result) {
			defer waitgroup.Done()
			result.Score += Score(words, UpperUnaccentUnpunctString(result.Name))
			if result.Score > 0 {
				mutex.Lock()
				results.Add(result)
				mutex.Unlock()
			}
		}(element.(*Result))
	})
	waitgroup.Wait()

	return results
}
