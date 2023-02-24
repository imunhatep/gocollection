package immutable

import (
	"fmt"
	"github.com/imunhatep/gocollection/helper"
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
)

func NewStrTestSortedMap(size int) SortedMap[string, string] {
	var keys, values []string
	for _, i := range helper.Range(1, size) {
		keys = append(keys, fmt.Sprintf("key_%d", i))
		values = append(values, fmt.Sprintf("value_%d", i))
	}

	return NewSortedMap[string, string]().SetKeys(keys...).SetValues(values...)
}

func NewIntTestSortedMap(size int) SortedMap[string, int] {
	values := map[string]int{}
	for _, i := range helper.Range(1, size) {
		values[fmt.Sprintf("key_%d", i)] = i
	}

	return ToSortedMap(values).SortByKey(helper.StrSort)
}

func TestSortedMapHeadTail(t *testing.T) {
	l1 := NewStrTestSortedMap(3)

	assert.False(t, l1.IsEmpty(), "map is not empty")
	assert.NotEmpty(t, l1.Head().OrEmpty(), "none empty map mst return head")

	// Head / Tail
	for _, p := range l1.ToSlice() {
		// assert equality
		assert.Equal(t, p.V2, l1.Head().MustGet().V2, "they should be equal")
		l1 = l1.Tail()
	}
}

func TestSortedMapSortedMap(t *testing.T) {
	double := func(i string, p int) int { return p * 2 }
	// map
	l1 := NewIntTestSortedMap(5)
	r1 := l1.Map(double)

	for _, v := range l1.ToSlice() {
		assert.Equal(t, double("", v.V2), r1.GetOrElse(v.V1, -1), "they should be equal")
		assert.Equal(t, double("", v.V2), r1.Get(v.V1).MustGet().V2, "they should be equal")
	}
}

func TestSortedMapUnique(t *testing.T) {
	l1 := NewIntTestSortedMap(5)
	l2 := l1.Join(NewIntTestSortedMap(3))

	assert.Equal(t, l1.ToMap(), l2.ToMap(), "they should be equal")
}

func TestSortedMapRemove(t *testing.T) {
	l1 := NewStrTestSortedMap(5)
	l2 := l1.Remove(l1.Head().MustGet().V1)

	assert.NotEqual(t, l1.ToSlice(), l2.ToSlice(), "they should be equal")
	assert.Equal(t, l1.Tail().ToSlice(), l2.ToSlice(), "they should be equal")

	for k, _ := range l1.ToMap() {
		l2 = l2.Remove(k)
	}

	assert.Empty(t, l2.ToSlice())
}

func TestSortedMapCompare(t *testing.T) {
	l1 := NewStrTestSortedMap(5)
	v1 := l1.Head().MustGet().V2

	assert.True(t, l1.Contains(v1), "seq must contain a value")

	i1 := l1.Find(func(i string, v string) bool { return v == v1 }).OrEmpty()
	assert.Equal(t, v1, i1.V2, "they should be equal")

	e1 := NewSequence[string]().Find(func(v string) bool { return v == "empty" })
	assert.Empty(t, e1.OrEmpty())

	s1 := SortedMapFoldLeft(
		l1,
		NewSortedMap[string, testStruct](),
		func(z SortedMap[string, testStruct], k string, v string) SortedMap[string, testStruct] {
			return z.Update(k, testStruct{v})
		},
	)
	sv1 := s1.Head().MustGet()

	assert.True(t, s1.Contains(sv1.V2))
}

func TestSortedMapFilter(t *testing.T) {
	l1 := NewIntTestSortedMap(5)
	l2Size := l1.Size() - 4
	l2 := l1.Filter(func(i string, v int) bool { return v < l2Size })

	assert.Equal(t, l2Size, l2.Size(), "they should be equal")
}

func TestSortedMapFolding(t *testing.T) {
	l1 := NewIntTestSortedMap(5)
	l2 := l1.FoldRight(
		NewSortedMap[string, int](),
		func(s SortedMap[string, int], k string, v int) SortedMap[string, int] { return s.Update(k, v) },
	)

	assert.Equal(t, l1.Reversed().ToSlice(), l2.ToSlice(), "they should be equal")
}

func TestSortedMapEmpty(t *testing.T) {
	l1 := NewSortedMap[int, int]()

	assert.True(t, l1.IsEmpty())
	assert.Empty(t, l1.Tail().ToSlice())
	assert.Empty(t, l1.Head().OrEmpty())
	assert.Empty(t, l1.Get(0).OrEmpty())
	assert.False(t, l1.Contains(0))
	assert.False(t, l1.ContainsKey(0))
}

func TestSortedMapReversed(t *testing.T) {
	l1 := NewStrTestSortedMap(5)
	for _, p := range l1.Reversed().Reversed().ToSlice() {
		assert.Equal(t, l1.Head().MustGet().V1, p.V1, "they should be equal")
		l1 = l1.Tail()
	}
}

func TestSortedMapSort(t *testing.T) {
	l1 := NewStrTestSortedMap(4)
	assert.Equal(t, l1.Values().ToSlice(), l1.Reversed().Sort(helper.StrSort).Values().ToSlice(), "they should be equal")
}

func TestSortedMapSortByKey(t *testing.T) {
	l1 := NewStrTestSortedMap(10)
	assert.Equal(t, l1.ToSlice(), l1.Reversed().SortByKey(helper.StrSort).ToSlice(), "they should be equal")
}

func TestSortedMapKeys(t *testing.T) {
	l1 := NewStrTestSortedMap(10)
	assert.Equal(t, l1.Keys().ToSlice(), l1.Reversed().SortByKey(helper.StrSort).Keys().ToSlice(), "they should be equal")
}

func TestSortedMapRace(t *testing.T) {
	update := func(lst SortedMap[int, int], s int) SortedMap[int, int] {
		a := 30
		for i := s * a; i < (s+1)*a; i++ {
			lst = lst.Update(i, i)
		}

		return lst
	}

	var mx sync.Mutex
	l1 := NewSortedMap[int, int]()
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

	l2 := NewSortedMap[int, int]()
	l2 = update(l2, 0)
	l2 = update(l2, 1)
	l2 = update(l2, 2)

	assert.Equal(t, l2.ToSlice(), l1.Sort(helper.IntSort).ToSlice(), "they should be equal")
}
