package immutable

import (
	"fmt"
	"github.com/imunhatep/gocollection/helper"
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
)

func NewStrTestMap(size int) HashMap[string, string] {
	values := map[string]string{}
	for _, i := range helper.Range(1, size) {
		values[fmt.Sprintf("key_%d", i)] = fmt.Sprintf("value_%d", i)
	}

	return ToMap(values)
}

func NewIntTestMap(size int) HashMap[string, int] {
	values := map[string]int{}
	for _, i := range helper.Range(1, size) {
		values[fmt.Sprintf("key_%d", i)] = i
	}

	return ToMap(values)
}

func TestMapMap(t *testing.T) {
	double := func(i string, p int) int { return p * 2 }
	// map
	l1 := NewIntTestMap(5)
	r1 := l1.Map(double)

	for _, v := range l1.ToSlice() {
		assert.Equal(t, double("", v.V2), r1.Get(v.V1).MustGet().V2, "they should be equal")
		assert.Equal(t, double("", v.V2), r1.GetOrElse(v.V1, -1), "they should be equal")
	}
}

func TestMapUnique(t *testing.T) {
	l1 := NewIntTestMap(5)
	l2 := l1.Join(NewIntTestMap(3))

	assert.ElementsMatch(t, l1.ToSlice(), l2.ToSlice(), "unique map should stay unchanged")
}

func TestMapRemove(t *testing.T) {
	l1 := NewStrTestMap(5)
	l2 := l1.Remove(l1.Keys().Head().MustGet())

	assert.NotEqual(t, l1.ToSlice(), l2.ToSlice(), "they should not be equal")
	assert.Equal(t, l1.Size()-1, l2.Size(), "map size should decrease")

	for k, _ := range l1.ToMap() {
		l2 = l2.Remove(k)
	}

	assert.Empty(t, l2.ToSlice())
}

func TestMapCompare(t *testing.T) {
	l1 := NewStrTestMap(5)
	v1 := l1.Get(l1.Keys().Head().MustGet()).MustGet().V2

	assert.True(t, l1.Contains(v1), "seq must contain a value")

	i1 := l1.Find(func(i string, v string) bool { return v == v1 }).OrEmpty()
	assert.Equal(t, v1, i1.V2, "they should be equal")

	e1 := NewSequence[string]().Find(func(v string) bool { return v == "empty" })
	assert.Empty(t, e1.OrEmpty())

	s1 := HashMapFoldLeft(
		l1,
		NewMap[string, testStruct](),
		func(z HashMap[string, testStruct], k string, v string) HashMap[string, testStruct] {
			return z.Update(k, testStruct{v})
		},
	)

	sv1 := s1.Get(s1.Keys().Head().MustGet()).MustGet().V2
	assert.True(t, s1.Contains(sv1))
}

func TestMapFilter(t *testing.T) {
	l1 := NewIntTestMap(5)
	l2Size := l1.Size() - 4
	l2 := l1.Filter(func(i string, v int) bool { return v < l2Size })

	assert.Equal(t, l2Size, l2.Size(), "they should be equal")
}

func TestMapFolding(t *testing.T) {
	l1 := NewIntTestMap(5)
	l2 := l1.FoldLeft(
		NewMap[string, int](),
		func(s HashMap[string, int], k string, v int) HashMap[string, int] { return s.Update(k, v) },
	)

	assert.ElementsMatch(t, l1.ToSlice(), l2.ToSlice(), "they should be equal")
}

func TestMapEmpty(t *testing.T) {
	l1 := NewMap[int, int]()

	assert.True(t, l1.IsEmpty())
	assert.Empty(t, l1.Get(0).OrEmpty())
}

func TestMapRace(t *testing.T) {
	update := func(lst HashMap[int, int], s int) HashMap[int, int] {
		a := 30
		for i := s * a; i < (s+1)*a; i++ {
			lst = lst.Update(i, i)
		}

		return lst
	}

	var mx sync.Mutex
	l1 := NewMap[int, int]()
	updateRx := func(wg *sync.WaitGroup, s int) {
		mx.Lock()
		defer mx.Unlock()
		defer wg.Done()

		l1 = update(l1, s)
	}

	var wg sync.WaitGroup
	wg.Add(3)

	go updateRx(&wg, 0)
	go updateRx(&wg, 1)
	go updateRx(&wg, 2)

	wg.Wait()

	l2 := NewMap[int, int]()
	l2 = update(l2, 0)
	l2 = update(l2, 1)
	l2 = update(l2, 2)

	assert.ElementsMatch(t, l2.ToSlice(), l1.ToSlice(), "hashmap values should be equal")
}
