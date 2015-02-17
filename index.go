package goloc

import (
	"strconv"
)

type Index interface {
	Add(loc Localisation)
	Get(string) Localisation
	Search(string, int, Scorer) *LinkedList

	getInternalIdsForKey(string) *LinkedList
	addInternalLocalisation(Localisation, map[string]bool)
}

func internalAdd(index Index, loc Localisation) {
	keys := Nkeys(Split(Partialphone(loc.GetName())))
	index.addInternalLocalisation(loc, keys)
}

func internalSearch(index Index, search string, number int, scorer Scorer) *LinkedList {
	if scorer == nil {
		scorer = DefaultScorer
	}
	results := NewLinkedList()
	keys := Nkeys(Split(Partialphone(search)))
	mapRes := make(map[string]*Result)
	jobChan := make(chan bool, maxRoutine)
	resultChan := make(chan *Result, maxRoutine)
	var maxScore, tmpScore, i, l int
	var result *Result
	var id, k string
	var ok bool
	var ids *LinkedList
	var elem *LinkedElement
	var err error
	var loc Localisation
	numKeyInternal := 2147483647
	for k = range keys {
		ids = index.getInternalIdsForKey(k)
		if ids != nil && ids.Size < numKeyInternal {
			numKeyInternal = ids.Size
		}
	}
	numKeyInternal = Max(numKeyInternal, maxKeyInternal)
	for k = range keys {
		ids = index.getInternalIdsForKey(k)
		if ids != nil && ids.Size <= numKeyInternal {
			if _, err = strconv.Atoi(k); err != nil {
				// is not num
				tmpScore = 3 + len(k)*len(k)
			} else {
				// is num
				tmpScore = 1
			}
			for elem = ids.First; elem != nil; elem = elem.Next {
				id = elem.Element.(string)
				result, ok = mapRes[id]
				if ok {
					result.Score += tmpScore
				} else {
					result = new(Result)
					result.Score = tmpScore
					result.Id = id
					mapRes[id] = result
				}
				if result.Score > maxScore {
					maxScore = result.Score
				}
			}
		}
	}

	// remove num score
	maxScore -= 3

	for id, result = range mapRes {
		if result.Score < maxScore {
			delete(mapRes, id)
		}
	}

	l = len(mapRes)
	if l > maxInternal {
		return results
	}

	go func() {
		for id, result = range mapRes {
			loc = index.Get(result.Id)
			jobChan <- true
			result.Search = search
			if loc != nil {
				result.Name = loc.GetName()
				result.Lat = loc.GetLat()
				result.Lon = loc.GetLon()
				result.Type = loc.GetType()
			}
			go scoreWorker(result, scorer, jobChan, resultChan)
		}
	}()

	maxScore = 1
	l = len(mapRes)
	for i = 0; i < l; i++ {
		select {
		case result = <-resultChan:
			if result.Score > maxScore {
				maxScore = result.Score
			}
		}
	}
	close(resultChan)
	close(jobChan)

	l = Min(len(mapRes), number)
	for results.Size < l && maxScore > 0 {
		tmpScore = 0
		for _, result = range mapRes {
			if result.Score == maxScore {
				results.AddLast(result)
			} else if result.Score < maxScore {
				if result.Score > tmpScore {
					tmpScore = result.Score
				}
			}
		}
		maxScore = tmpScore
	}

	return results
}
