package dict

import (
	"fmt"
	"github.com/imunhatep/gocollection/helper"
	"github.com/imunhatep/gocollection/slice"
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
)

type testStruct struct {
	Some string
}

func NewStrTestMap(size int) map[string]string {
	values := map[string]string{}
	for _, i := range helper.Range(1, size) {
		values[fmt.Sprintf("key_%d", i)] = fmt.Sprintf("value_%d", i)
	}

	return values
}

func NewIntTestMap(size int) map[string]int {
	values := map[string]int{}
	for _, i := range helper.Range(1, size) {
		values[fmt.Sprintf("key_%d", i)] = i
	}

	return values
}

func TestMapMap(t *testing.T) {
	double := func(i string, p int) int { return p * 2 }
	// map
	l1 := NewIntTestMap(5)
	r1 := Map(l1, double)

	for k, v := range l1 {
		assert.Equal(t, double("", v), r1[k], "they should be equal")
		assert.Equal(t, double("", v), GetOrElse(r1, k, -1), "they should be equal")
	}
}

func TestMapRemove(t *testing.T) {
	l1 := NewStrTestMap(5)
	l2 := Remove(l1, slice.Head(Keys(l1)).OrEmpty())

	assert.NotEqual(t, l1, l2, "they should not be equal")
	assert.Equal(t, Size(l1)-1, Size(l2), "map size should decrease")

	for k, _ := range l1 {
		l2 = Remove(l2, k)
	}

	assert.Empty(t, l2)
}

func TestMapUnique(t *testing.T) {
	l1 := NewIntTestMap(5)
	l2 := Merge(l1, NewIntTestMap(3))

	assert.Equal(t, l1, l2, "unique map should stay unchanged")
}

func TestMapCompare(t *testing.T) {
	l1 := NewStrTestMap(5)
	v1 := slice.Head(Values(l1)).MustGet()

	assert.True(t, Contains(l1, v1), "seq must contain a value")

	i1 := Find(l1, func(i, v string) bool { return v == v1 }).OrEmpty()
	assert.Equal(t, v1, i1.V2, "they should be equal")

	s1 := Fold(
		l1,
		map[string]testStruct{},
		func(z map[string]testStruct, k string, v string) map[string]testStruct {
			z[k] = testStruct{v}
			return z
		},
	)

	sv1 := slice.Head(Values(s1)).MustGet()
	assert.True(t, Contains(s1, sv1))
}

func TestMapFilter(t *testing.T) {
	l1 := NewIntTestMap(5)
	l2Size := Size(l1) - 4
	l2 := Filter(l1, func(i string, v int) bool { return v < l2Size })

	assert.Equal(t, l2Size, Size(l2), "they should be equal")
}

func TestMapFolding(t *testing.T) {
	l1 := NewIntTestMap(5)
	l2 := Fold(
		l1,
		map[string]int{},
		func(z map[string]int, k string, v int) map[string]int {
			z[k] = v
			return z
		},
	)

	assert.Equal(t, l1, l2, "they should be equal")
}

func TestMapEmpty(t *testing.T) {
	l1 := map[int]int{}

	assert.True(t, IsEmpty(l1))
	assert.Empty(t, Get(l1, 0).OrEmpty())
}

func TestMapRace(t *testing.T) {
	update := func(lst map[int]int, s int) map[int]int {
		a := 30
		for i := s * a; i < (s+1)*a; i++ {
			lst = Set(lst, i, i)
		}

		return lst
	}

	var mx sync.Mutex
	l1 := map[int]int{}
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

	l2 := map[int]int{}
	l2 = update(l2, 0)
	l2 = update(l2, 1)
	l2 = update(l2, 2)

	assert.Equal(t, l2, l1, "hashmap values should be equal")
}
