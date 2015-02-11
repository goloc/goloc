package goloc

import (
	"encoding/gob"
	"fmt"
	"os"
	"strconv"
)

type Memindex struct {
	Localisations map[string]Localisation
	Keys          map[string]*LinkedList
}

func NewMemindex() *Memindex {
	mi := new(Memindex)
	mi.Clear()
	GobRegister()
	return mi
}

func NewMemindexFromFile(filename string) *Memindex {
	fmt.Printf("load %v\n", filename)
	mi := NewMemindex()
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
		return mi
	}
	defer file.Close()

	dataDecoder := gob.NewDecoder(file)
	dataDecoder.Decode(&mi.Keys)
	dataDecoder.Decode(&mi.Localisations)
	fmt.Printf("%v localisations\n", mi.SizeLocalisation())
	fmt.Printf("%v keys\n", mi.SizeIndex())

	return mi
}

func (mi *Memindex) SaveInFile(filename string) {
	fmt.Printf("save %v\n", filename)
	file, err := os.Create(filename)
	if err != nil {
		panic(err)
		return
	}
	encoder := gob.NewEncoder(file)
	if err = encoder.Encode(mi.Keys); err != nil {
		panic(err)
	}
	if err = encoder.Encode(mi.Localisations); err != nil {
		panic(err)
	}
	if err = file.Close(); err != nil {
		panic(err)
	}
}
func (mi *Memindex) Add(loc Localisation) {
	keys := Nkeys(Split(Partialphone(loc.GetName())))
	var id, k string
	id = loc.GetId()
	mi.Localisations[id] = loc
	var ids *LinkedList
	for k, _ = range keys {
		ids = mi.Keys[k]
		if ids == nil {
			ids = NewLinkedList()
			mi.Keys[k] = ids
		}
		ids.AddLast(id)
	}
}

func (mi *Memindex) SizeLocalisation() int {
	return len(mi.Localisations)
}

func (mi *Memindex) SizeIndex() int {
	return len(mi.Keys)
}

func (mi *Memindex) Clear() {
	mi.Localisations = make(map[string]Localisation)
	mi.Keys = make(map[string]*LinkedList)
}

func (mi *Memindex) Get(id string) Localisation {
	loc := mi.Localisations[id]
	return loc
}

func (mi *Memindex) Remove(id string) {
	delete(mi.Localisations, id)
}

func (mi *Memindex) Search(search string, number int) *LinkedList {
	fmt.Printf("\nSearch '%v'\n", search)
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
		ids = mi.Keys[k]
		if ids != nil && ids.Size < numKeyInternal {
			numKeyInternal = ids.Size
		}
	}
	numKeyInternal = Max(numKeyInternal, maxKeyInternal)
	fmt.Printf("Keys: ", search)
	for k = range keys {
		fmt.Printf("%v ", k)
		ids = mi.Keys[k]
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
					loc = mi.Localisations[id]
					if loc != nil {
						result = new(Result)
						result.Score = tmpScore
						result.Id = id
						mapRes[id] = result
					}
				}
				if result.Score > maxScore {
					maxScore = result.Score
				}
			}
		}
	}
	fmt.Println()

	// remove num score
	maxScore -= 3

	fmt.Printf("1 - found=%v maxScore=%v\n", len(mapRes), maxScore)

	for id, result = range mapRes {
		if result.Score < maxScore {
			delete(mapRes, id)
		}
	}

	l = len(mapRes)
	if l > maxInternal {
		fmt.Printf("2 - Too much found=%v\n", l)
		return results
	} else {
		fmt.Printf("2 - found=%v\n", l)
	}

	go func() {
		for _, result = range mapRes {
			jobChan <- true
			loc = mi.Localisations[result.Id]
			result.Search = search
			result.Name = loc.GetName()
			result.Lat = loc.GetLat()
			result.Lon = loc.GetLon()
			result.Priority = loc.GetPriority()
			go scoreWorker(result, jobChan, resultChan)
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

	fmt.Printf("3 - maxScore=%v\n", maxScore)

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

	fmt.Printf("4 - found=%v\n", results.Size)

	return results
}
