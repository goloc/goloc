// Copyright 2015 Mathieu MAST. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.
package goloc

import (
	"fmt"
	"github.com/goloc/container"
	"sync"
)

type internal struct {
	tolerance          int
	keyLimit           int
	locLimit           int
	get                func(string) Location
	getIds             func(string) container.Container
	addLocationAndKeys func(loc Location, keys []string)
}

func (inter *internal) add(loc Location) {
	name := loc.GetName()
	mkeys := MSplit(Partialphone(name))
	inter.addLocationAndKeys(loc, mkeys)
}

func (inter *internal) search(search string, number int, filter Filter) container.Container {
	if filter == nil {
		filter = DefaultFilter
	}

	words := Split(UpperUnaccentUnpunctString(search))
	mwords := MSplit(UpperUnaccentUnpunctString(search))
	mkeys := MSplit(Partialphone(search))

	var waitgroup sync.WaitGroup
	var mutex sync.Mutex

	max := 0
	mapRes := make(map[string]*Result)
	for _, key := range mkeys {
		waitgroup.Add(1)
		go func(key string) {
			defer waitgroup.Done()
			ids := inter.getIds(key)
			if ids != nil && ids.GetSize() > 0 {
				if ids.GetSize() <= inter.keyLimit {
					ids.Visit(func(element interface{}, i int) {
						id := element.(string)
						mutex.Lock()
						result, ok := mapRes[id]
						if ok {
							result.Score++
						} else {
							result = new(Result)
							result.Score = 1
							result.Id = id
							mapRes[id] = result
						}
						if result.Score > max {
							max = result.Score
						}
						mutex.Unlock()
					})
				} else {
					fmt.Printf("Key %v has to many elements (%v).\n", key, ids.GetSize())
				}
			}
		}(key)
	}
	waitgroup.Wait()

	tmpResults := container.NewLimitedBinaryTree(CompareScoreResult, inter.locLimit, true)
	minKeyScore := max - inter.tolerance
	for _, result := range mapRes {
		if result.Score >= minKeyScore {
			waitgroup.Add(1)
			go func(result *Result) {
				defer waitgroup.Done()
				loc := inter.get(result.Id)
				if loc != nil {
					result.Score = 0
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
			}(result)
		}
	}
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
