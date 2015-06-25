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
		cleansearch := UpperUnaccentUnpunctString(" " + parameters.Get("search").(string) + " ")
		if s.index.GetStopWords() != nil {
			s.index.GetStopWords().Visit(func(element interface{}, i int) {
				word := " " + element.(string) + " "
				cleansearch = strings.Join(strings.Split(cleansearch, word), " ")
			})
		}
		parameters.Set("cleansearch", cleansearch)

		res := s.searchInternal(parameters)
		future.Resolve(res)
	}()

	return future.GetPromise()
}

func (s *ConcurrentSniffer) searchInternal(parameters Parameters) container.Container {
	search := parameters.Get("search").(string)
	cleansearch := parameters.Get("cleansearch").(string)
	words := Split(cleansearch)
	mwords := MSplit(cleansearch)
	mkeys := MSplit(Partialphone(cleansearch))

	var waitgroup sync.WaitGroup

	keysCounter := container.NewCounter()
	mkeys.Visit(func(element interface{}, i int) {
		waitgroup.Add(1)
		go func(key string) {
			defer waitgroup.Done()
			val := s.index.GetNbIds(key)
			if val > 0 {
				keysCounter.Incr(key, val)
			}
		}(element.(string))
	})
	waitgroup.Wait()

	filter := parameters.Get("filter").(func(*Result) bool)

	tmpResults := container.NewLimitedBinaryTree(CompareScoreResult, parameters.Get("workLimit").(int), true)
	keysCounter.Visit(func(element interface{}, i int) {
		if i <= parameters.Get("tolerance").(int) {
			waitgroup.Add(1)
			go func(count *container.Count) {
				defer waitgroup.Done()
				ids := s.index.GetIds(count.Key)
				if ids != nil && ids.GetSize() > 0 {
					ids.Visit(func(element interface{}, i int) {
						id := element.(string)
						loc := s.index.Get(id)
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
								numbered := bag.GetNumberedPoint(cleansearch)
								if numbered != nil {
									result.Number = numbered.GetNumber()
									result.Lat = numbered.GetLat()
									result.Lon = numbered.GetLon()
								}
							}
							if filter(result) {
								result.Score += QuickScore(mwords, UpperUnaccentUnpunctString(result.Name))
								if result.Score > 0 {
									tmpResults.Add(result)
								}
							}
						}
					})
				}
			}(element.(*container.Count))
		}
	})
	waitgroup.Wait()

	results := container.NewLimitedBinaryTree(CompareScoreResult, parameters.Get("limit").(int), true)
	tmpResults.Visit(func(element interface{}, i int) {
		waitgroup.Add(1)
		go func(result *Result) {
			defer waitgroup.Done()
			result.Score += Score(words, UpperUnaccentUnpunctString(result.Name))
			if result.Score > 0 {
				results.Add(result)
			}
		}(element.(*Result))
	})
	waitgroup.Wait()

	return results
}
