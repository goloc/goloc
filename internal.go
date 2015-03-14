// Copyright 2015 Mathieu MAST. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.
package goloc

import (
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
	addLocationAndKeys func(loc Location, keys []string)
}

func (inter *internal) add(loc Location) {
	name := loc.GetName()
	keys := MSplit(Partialphone(name), 1, 3)
	inter.addLocationAndKeys(loc, keys)
}

func (inter *internal) search(search string, number int, filter Filter) container.Container {
	if filter == nil {
		filter = DefaultFilter
	}
	words := MSplit(UpperUnaccentUnpunctString(search), 2, 2)
	keys := MSplit(Partialphone(search), 1, 3)
	var waitgroup sync.WaitGroup
	var mutex sync.Mutex

	maxKeyScore := 0
	bestIdsSize := 0
	var bestIds container.Container
	mapRes := make(map[string]*Result)
	for _, key := range keys {
		waitgroup.Add(1)
		go func(key string) {
			defer waitgroup.Done()
			ids := inter.getIds(key)
			mutex.Lock()
			defer mutex.Unlock()
			if ids != nil && ids.GetSize() > 0 {
				if ids.GetSize() <= inter.keyLimit {
					ids.Visit(func(element interface{}, i int) {
						id := element.(string)
						result, ok := mapRes[id]
						if ok {
							result.Score++
						} else {
							result = new(Result)
							result.Score++
							result.Id = id
							mapRes[id] = result
						}
						if result.Score > maxKeyScore {
							maxKeyScore = result.Score
						}
					})
				} else {
					if bestIdsSize == 0 || ids.GetSize() < bestIdsSize {
						bestIds = ids
						bestIdsSize = ids.GetSize()
					}
				}
			}
		}(key)
	}
	waitgroup.Wait()

	if len(mapRes) == 0 && bestIds != nil {
		bestIds.Visit(func(element interface{}, i int) {
			id := element.(string)
			result := new(Result)
			result.Score = 1
			result.Id = id
			mapRes[id] = result
		})
		maxKeyScore = 1
	}

	tmpResults := container.NewLimitedBinaryTree(CompareScoreResult, inter.locLimit, true)
	minKeyScore := maxKeyScore - inter.tolerance
	maxKeyScore = 0
	for _, result := range mapRes {
		if result.Score >= minKeyScore {
			waitgroup.Add(1)
			go func(result *Result) {
				defer waitgroup.Done()
				loc := inter.get(result.Id)
				if loc != nil {
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
						name := UpperUnaccentUnpunctString(result.Name)
						for _, word := range words {
							if strings.Contains(name, word) {
								result.Score += len(word)
							}
						}
						mutex.Lock()
						if result.Score > maxKeyScore {
							maxKeyScore = result.Score
						}
						tmpResults.Add(result)
						mutex.Unlock()
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
			result.Score += Score(result.Search, result.Name)
			mutex.Lock()
			defer mutex.Unlock()
			if result.Score > 0 {
				results.Add(result)
			}
		}(element.(*Result))
	})
	waitgroup.Wait()

	return results
}
