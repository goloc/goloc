package goloc

import (
	"encoding/gob"
	"os"
	"sort"
	"strconv"
)

type LinkedId struct {
	Next *LinkedId
	Id   string
}

type Memindex struct {
	Localisations map[string]Localisation
	Phoneindex    map[string]*LinkedId
}

type ByScore []*Result

func (a ByScore) Len() int           { return len(a) }
func (a ByScore) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByScore) Less(i, j int) bool { return a[i].Score > a[j].Score }

func (mi *Memindex) Add(loc Localisation) {
	id := loc.GetId()
	mi.Localisations[id] = loc
	keys := Split(Partialphone(loc.GetName()))
	for _, k := range keys {
		l := len(k)
		i := l
		if l >= 2 {
			i = 1
		}
		for ; i <= l; i++ {
			subk := k[0:i]
			linkedId := new(LinkedId)
			linkedId.Id = id
			nextLinkedId, ok := mi.Phoneindex[subk]
			if ok {
				linkedId.Next = nextLinkedId
			}
			mi.Phoneindex[subk] = linkedId
		}
	}
}
func (mi *Memindex) SizeLocalisation() int {
	return len(mi.Localisations)
}
func (mi *Memindex) SizeIndex() int {
	return len(mi.Phoneindex)
}
func (mi *Memindex) Clear() {
	mi.Localisations = make(map[string]Localisation)
	mi.Phoneindex = make(map[string]*LinkedId)
}
func (mi *Memindex) Get(id string) Localisation {
	loc := mi.Localisations[id]
	return loc
}
func (mi *Memindex) Remove(id string) {
	delete(mi.Localisations, id)
}
func (mi *Memindex) Search(search string, max int, minScore int, maxDeviation int) []*Result {
	keys := Split(Partialphone(search))
	mapResult := make(map[string]*Result)
	var maxScore, keysScore, numResult, tmpScore, i int
	var result *Result
	var id string
	var ok bool
	for _, k := range keys {
		if _, err := strconv.Atoi(k); err != nil {
			// is not num
			keysScore += 2 + len(k)*len(k)
		}
	}
	for _, k := range keys {
		for linkedId := mi.Phoneindex[k]; linkedId != nil; linkedId = linkedId.Next {
			id = linkedId.Id
			result, ok = mapResult[id]
			if _, err := strconv.Atoi(k); err != nil {
				// is not num
				tmpScore = 2 + len(k)*len(k)
			} else {
				// is num
				tmpScore = 1
			}
			if ok {
				result.Score += tmpScore
			} else {
				result = new(Result)
				result.Score = tmpScore
				result.Localisation = mi.Localisations[id]
				if result.Localisation != nil {
					mapResult[id] = result
				}
			}
		}
	}

	numResult = 0
	for id, result = range mapResult {
		if result.Score >= keysScore {
			numResult++
		} else {
			delete(mapResult, id)
		}
	}

	maxScore = 0
	numResult = 0
	for id, result = range mapResult {
		tmpScore = Score(search, result.Localisation.GetName())
		if tmpScore < minScore {
			delete(mapResult, id)
		} else {
			result.Score = tmpScore
			if tmpScore > maxScore {
				maxScore = tmpScore
			}
			numResult++
		}
	}

	numResult = 0
	for id, result = range mapResult {
		if result.Score < maxScore-maxDeviation {
			delete(mapResult, id)
		} else {
			numResult++
		}
	}

	nb := Min(numResult, max)
	results := make([]*Result, nb)

	i = 0
	for _, result := range mapResult {
		results[i] = result
		i++
		if i >= nb {
			break
		}
	}

	sort.Sort(ByScore(results))

	return results
}
func (mi *Memindex) SaveInFile(filename string) {
	file, err := os.Create(filename)
	if err != nil {
		panic(err)
		return
	}
	defer file.Close()

	encoder := gob.NewEncoder(file)
	if err := encoder.Encode(&(mi.Phoneindex)); err != nil {
		panic(err)
	}
	if err := encoder.Encode(&(mi.Localisations)); err != nil {
		panic(err)
	}

}
func NewMemindex() *Memindex {
	mi := new(Memindex)
	mi.Clear()
	gob.RegisterName("core.Street", &Street{})
	gob.RegisterName("core.Address", &Address{})
	gob.RegisterName("core.Point", &Point{})
	gob.RegisterName("core.Zone", &Zone{})
	return mi
}
func NewMemindexFromFile(filename string) *Memindex {
	mi := NewMemindex()
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
		return mi
	}
	defer file.Close()

	dataDecoder := gob.NewDecoder(file)
	dataDecoder.Decode(&(mi.Phoneindex))
	dataDecoder.Decode(&(mi.Localisations))
	return mi
}
