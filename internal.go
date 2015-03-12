// Copyright 2015 Mathieu MAST. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.
package goloc

import (
	"github.com/goloc/container"
	"strings"
	"sync"
)

func internalAdd(index Index, loc Location, addLocationAndKeys func(loc Location, keys []string)) {
	name := loc.GetName()
	keys := MSplit(Partialphone(name))
	addLocationAndKeys(loc, keys)
}

func internalSearch(index Index, search string, number int, keyLimit int, locLimit int, filter Filter, getIds func(string) container.Container) container.Container {
	if filter == nil {
		filter = DefaultFilter
	}
	words := MSplit(UpperUnaccentUnpunctString(search))
	keys := MSplit(Partialphone(search))
	var waitgroup sync.WaitGroup
	var mutex sync.Mutex

	mapIds := make(map[string]container.Container)
	for _, key := range keys {
		waitgroup.Add(1)
		go func(key string) {
			defer waitgroup.Done()
			ids := getIds(key)
			mutex.Lock()
			defer mutex.Unlock()
			mapIds[key] = ids
		}(key)
	}
	waitgroup.Wait()

	maxKeyScore := 0
	bestIdsSize := 0
	var bestIds container.Container
	mapRes := make(map[string]*Result)
	for key, ids := range mapIds {
		if ids != nil && ids.GetSize() > 0 {
			if ids.GetSize() <= keyLimit {
				ids.Visit(func(element interface{}, i int) {
					id := element.(string)
					result, ok := mapRes[id]
					l := len(key)
					if ok {
						result.Score += l * l
					} else {
						result = new(Result)
						result.Score = l * l
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
	}

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

	tmpResults := container.NewLimitedBinaryTree(CompareScoreResult, locLimit, true)
	minKeyScore := maxKeyScore - 2
	maxKeyScore = 0
	for _, result := range mapRes {
		if result.Score >= minKeyScore {
			waitgroup.Add(1)
			go func(result *Result) {
				defer waitgroup.Done()
				loc := index.Get(result.Id)
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
