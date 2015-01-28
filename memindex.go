package goloc

import (
	"encoding/gob"
	"os"
)

type LinkedId struct {
	Next *LinkedId
	Id   string
}

type Memindex struct {
	Localisations map[string]Localisation
	Phoneindex    map[string]*LinkedId
}

func (mi *Memindex) Add(loc Localisation) {
	id := loc.GetId()
	name := loc.GetName()
	mi.Localisations[id] = loc
	keys := splitSpacePunct(partialphone(name))
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
func (mi *Memindex) Search(search string) *map[string]*Result {
	keys := splitSpacePunct(partialphone(search))
	mapResult := make(map[string]*Result)
	maxScore := 0
	keysScore := 0
	for _, k := range keys {
		keysScore += len(k)
	}
	for _, k := range keys {
		for linkedId := mi.Phoneindex[k]; linkedId != nil; linkedId = linkedId.Next {
			id := linkedId.Id
			result, ok := mapResult[id]
			if ok {
				result.score += len(id)
			} else {
				result = new(Result)
				result.id = id
				result.score = len(id)
				result.loc = mi.Localisations[id]
				mapResult[id] = result
			}
			if result.score > maxScore {
				maxScore = result.score
			}
		}
	}
	targetScore := min(maxScore, keysScore)
	numResult := 0
	for id, result := range mapResult {
		if result.score >= targetScore {
			numResult++
			mapResult[id] = result
		} else {
			delete(mapResult, id)
		}
	}
	for _, result := range mapResult {
		loc := result.loc
		name := loc.GetName()
		result.score = score(search, name)
	}
	return &mapResult
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
