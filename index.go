// Copyright 2015 Mathieu MAST. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.
package goloc

import (
	"github.com/goloc/container"
	"strconv"
	"sync"
)

type Index interface {
	Add(loc Location)
	Get(string) Location
	Search(string, int, Scorer) *container.LimitedBinaryTree

	getInternalIdsForKey(string) *container.LinkedList
	addInternalLocation(Location, []string)
}

func internalAdd(index Index, loc Location) {
	keys := Nkeys(Split(Partialphone(loc.GetName())))
	index.addInternalLocation(loc, keys)
}

func internalSearch(index Index, search string, number int, scorer Scorer) *container.LimitedBinaryTree {
	if scorer == nil {
		scorer = DefaultScorer
	}
	results := container.NewLimitedBinaryTree(CompareScoreResult, number, true)
	keys := Nkeys(Split(Partialphone(search)))
	mapRes := make(map[string]*Result)
	var waitgroup sync.WaitGroup
	var mutex sync.Mutex
	var maxKeyScore, tmpScore int
	var ids *container.LinkedList
	var nums *container.LinkedList
	for _, key := range keys {
		ids = index.getInternalIdsForKey(key)
		if ids != nil && ids.Size <= maxInternal {
			if n, err := strconv.Atoi(key); err != nil {
				// is not num
				tmpScore = 3 + len(key)*len(key)
			} else {
				// is num
				tmpScore = 1
				nums.Push(n)
			}
			ids.Visit(func(element interface{}, i int) {
				id := element.(string)
				result, ok := mapRes[id]
				if ok {
					result.Score += tmpScore
				} else {
					result = new(Result)
					result.Score = tmpScore
					result.Id = id
					mapRes[id] = result
				}
				if result.Score > maxKeyScore {
					maxKeyScore = result.Score
				}
			})
		}
	}

	// remove num score
	maxKeyScore -= 3

	for _, result := range mapRes {
		if result.Score >= maxKeyScore {
			loc := index.Get(result.Id)
			result.Search = search
			result.Score = 0
			if loc != nil {
				result.Name = loc.GetResultName(search)
				result.Lat = loc.GetLat()
				result.Lon = loc.GetLon()
				result.Type = loc.GetType()
				waitgroup.Add(1)
				go func(result *Result) {
					defer waitgroup.Done()
					result.Score = scorer(result)
					mutex.Lock()
					defer mutex.Unlock()
					if result.Score > 0 {
						results.Add(result)
					}
				}(result)
			}
		}
	}

	waitgroup.Wait()

	return results
}
