// Copyright 2015 Mathieu MAST. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.
package goloc

import (
	"errors"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/goloc/concurrency"
	"github.com/goloc/container"
)

type ConcurrentSniffer struct {
	index     Index
	semaphore *concurrency.Semaphore
}

func NewConcurrentSniffer(index Index) *ConcurrentSniffer {
	s := new(ConcurrentSniffer)
	s.index = index
	return s
}

func (s *ConcurrentSniffer) Search(parameters Parameters) (container.Container, error) {
	if parameters.Get("search") == nil {
		return nil, errors.New("Search attribute is mandatory")
	}
	if parameters.Get("filter") == nil {
		parameters.Set("filter", DefaultFilter)
	}
	if parameters.Get("limit") == nil {
		parameters.Set("limit", defaultLimit)
	}
	if parameters.Get("tolerance") == nil {
		parameters.Set("tolerance", defaultTolerance)
	}
	if parameters.Get("maxWaitAcquire") == nil {
		parameters.Set("maxWaitAcquire", defaultMaxWaitAcquire)
	}
	if parameters.Get("maxWaitTraitment") == nil {
		parameters.Set("maxWaitTraitment", defaultMaxWaitTraitment)
	}
	if parameters.Get("workLimit") == nil {
		parameters.Set("workLimit", defaultWorkLimit)
	}

	if s.semaphore == nil {
		s.semaphore = concurrency.NewSemaphore(runtime.NumCPU())
	}
	if err := s.semaphore.Acquire(parameters.Get("maxWaitAcquire").(time.Duration)); err != nil {
		return nil, err
	}
	defer s.semaphore.Release()
	promise := s.searchPromise(parameters)
	element, err := promise.Wait(parameters.Get("maxWaitTraitment").(time.Duration))
	if element != nil {
		return element.(container.Container), err
	} else {
		return nil, err
	}
}

func (s *ConcurrentSniffer) searchPromise(parameters Parameters) *concurrency.Promise {
	future := concurrency.NewFuture()

	go func() {
		res := s.searchInternal(parameters)
		future.Resolve(res)
	}()

	return future.GetPromise()
}

func (s *ConcurrentSniffer) searchInternal(parameters Parameters) container.Container {
	search := parameters.Get("search").(string)
	filter := parameters.Get("filter").(func(*Result) bool)
	tolerance := parameters.Get("tolerance").(float32)
	cleanedSearch := Clean(search, s.index.GetStopWords())
	encodedSearch := Partialphone(cleanedSearch)
	keys := Split(encodedSearch)
	mkeys := MSplit(encodedSearch)
	mjkeys := MSplit(strings.Replace(encodedSearch, " ", "", -1))
	words := Split(cleanedSearch)

	var waitgroup sync.WaitGroup

	min1 := maxInt
	min2 := maxInt
	nbids := container.NewMap()
	mkeys.Visit(func(element interface{}, i int) {
		waitgroup.Add(1)
		go func(key string) {
			defer waitgroup.Done()
			nb := s.index.GetNbIds(key)
			nbids.Add(&container.KeyValue{Key: key, Value: nb})
			if !s.index.GetEncodedStopWords().Contains(key) {
				if nb > 0 && nb < min1 {
					min1 = nb
				}
			} else {
				if nb > 0 && nb < min2 {
					min2 = nb
				}
			}
		}(element.(string))
	})
	waitgroup.Wait()
	min := int((100.0 * float32(Min(min1, min2)) * (1 + tolerance)) / 100.0)

	mjkeys.Visit(func(element interface{}, i int) {
		waitgroup.Add(1)
		go func(key string) {
			defer waitgroup.Done()
			nb := s.index.GetNbIds(key)
			nbids.Add(&container.KeyValue{Key: key, Value: nb})
		}(element.(string))
	})
	waitgroup.Wait()

	ids := container.NewMap()
	nbids.Visit(func(element interface{}, i int) {
		waitgroup.Add(1)
		go func(keyValue *container.KeyValue) {
			defer waitgroup.Done()
			if keyValue.Value.(int) <= min {
				s.index.GetIds(keyValue.Key.(string)).Visit(func(element interface{}, i int) {
					ids.Add(&container.KeyValue{Key: element.(string), Value: true})
				})
			}
		}(element.(*container.KeyValue))
	})
	waitgroup.Wait()

	tmpResults := container.NewLimitedBinaryTree(CompareScoreResult, parameters.Get("workLimit").(int), true)
	ids.Visit(func(element interface{}, i int) {
		waitgroup.Add(1)
		go func(keyValue *container.KeyValue) {
			defer waitgroup.Done()
			loc := s.index.Get(keyValue.Key.(string))
			if loc != nil {
				result := new(Result)
				result.Id = keyValue.Key.(string)
				result.Search = search
				result.CleanedSearch = cleanedSearch
				result.Name = loc.GetName()
				result.CleanedName = loc.GetCleanedName()
				result.Lat = loc.GetLat()
				result.Lon = loc.GetLon()
				result.Type = loc.GetType()
				if filter(result) {
					result.Score += ContainerScore(keys, loc.GetEncodedName())
					if result.Score > 0 {
						tmpResults.Add(result)
					}
				}
			}
		}(element.(*container.KeyValue))
	})
	waitgroup.Wait()

	results := container.NewLimitedBinaryTree(CompareScoreResult, parameters.Get("limit").(int), true)
	tmpResults.Visit(func(element interface{}, i int) {
		waitgroup.Add(1)
		go func(result *Result) {
			defer waitgroup.Done()
			result.Score = ContainerScore(words, result.CleanedName)
			if result.Score > 0 {
				results.Add(result)
			}
		}(element.(*Result))
	})
	waitgroup.Wait()

	results.Visit(func(element interface{}, i int) {
		waitgroup.Add(1)
		go func(result *Result) {
			defer waitgroup.Done()
			loc := s.index.Get(result.Id)
			bag, ok := loc.(NumberedPointBag)
			if ok {
				minPos := maxInt
				bag.GetNumberedPoints().Visit(func(element interface{}, i int) {
					numbered := element.(NumberedPoint)
					num := UpperUnaccentUnpunctString(numbered.GetNumber())
					idx := strings.Index(" "+cleanedSearch+" ", " "+num+" ")
					if idx >= 0 && i < minPos {
						minPos = idx
						result.Number = numbered.GetNumber()
						result.Lat = numbered.GetLat()
						result.Lon = numbered.GetLon()
					}
				})
			}
		}(element.(*Result))
	})
	waitgroup.Wait()

	return results
}
