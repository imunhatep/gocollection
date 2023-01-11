package slice

import (
	"github.com/imunhatep/gocollection/helper"
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
)

type testStruct struct {
	Some string
}

func NewStrTestSeq() []string {
	return []string{"test1", "test2", "test3"}
}

func NewIntTestSeq() []int {
	return []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
}

func TestSeqHeadTail(t *testing.T) {
	l1 := NewStrTestSeq()

	// Head / Tail
	for _, v := range NewStrTestSeq() {
		// assert equality
		assert.Equal(t, Head(l1).MustGet(), v, "they should be equal")
		l1 = Tail(l1)
	}
}

func TestSeqMap(t *testing.T) {
	l1 := NewIntTestSeq()

	stringify := func(p int) string { return string(rune(p * 2)) }

	// map
	r1 := Map(l1, stringify)
	for i, v := range l1 {
		assert.Equal(t, stringify(v), Get(r1, i).OrEmpty(), "they should be equal")
	}
}

func TestSeqUnique(t *testing.T) {
	l1 := NewIntTestSeq()
	l2 := Unique(append(l1, NewIntTestSeq()...))

	assert.Equal(t, l1, l2, "they should be equal")
}

func TestSeqCompare(t *testing.T) {
	l1 := NewStrTestSeq()
	v1 := Head(l1).OrEmpty()

	assert.True(t, Contains(l1, v1), "seq must contain a value")

	i1 := Find(l1, func(v string) bool { return v == v1 }).OrEmpty()
	assert.Equal(t, v1, i1, "they should be equal")

	e1 := Find([]string{}, func(v string) bool { return v == "empty" })
	assert.Empty(t, e1.OrEmpty())

	i2 := FindWithIndex(l1, func(i int, v string) bool { return v == v1 }).OrEmpty()
	assert.Equal(t, v1, i2.V2, "they should be equal")

	e2 := FindWithIndex([]string{}, func(i int, v string) bool { return v == "empty" })
	assert.Empty(t, e2.OrEmpty())

	_, ok := IndexOfAny(l1, v1)
	assert.True(t, ok)

	s1 := FoldLeft(
		l1,
		[]testStruct{},
		func(z []testStruct, i int, v string) []testStruct { return append(z, testStruct{v}) },
	)
	sv1 := Head(s1).OrEmpty()
	assert.True(t, Contains(s1, sv1))

	_, ok = IndexOfAny(s1, sv1)
	assert.True(t, ok)
}

func TestSeqFilter(t *testing.T) {
	l1 := NewIntTestSeq()
	l2Size := Size(l1) - 4
	l2 := FilterWithIndex(l1, func(i int, v int) bool { return i < l2Size })

	assert.Equal(t, l2Size, len(l2), "they should be equal")
}

func TestSeqFolding(t *testing.T) {
	l1 := NewIntTestSeq()
	l2 := FoldRight(l1, []int{}, func(s []int, i int, v int) []int { return append(s, v) })

	assert.Equal(t, Reversed(l1), l2, "they should be equal")
}

func TestSeqEmpty(t *testing.T) {
	l1 := []int{}

	assert.True(t, IsEmpty(l1))
	assert.Empty(t, Tail(l1))
	assert.Empty(t, Head(l1).OrEmpty())
	assert.False(t, Contains(l1, 0))

	_, ok := IndexOf(l1, 0)
	assert.False(t, ok)
}

func TestSeqTail(t *testing.T) {
	l1 := NewStrTestSeq()

	// Tail
	r1 := Tail(l1)
	assert.NotEmpty(t, r1)
	l1[1] = "new"
	assert.NotEqual(t, l1[1:], r1, "should not be equal")
}

func TestSeqReversed(t *testing.T) {
	l1 := NewStrTestSeq()

	// Reversed
	r1 := Reversed(l1)
	for i, v := range Reversed(r1) {
		idx, _ := IndexOf(l1, v)
		assert.Equal(t, i, idx, "they should be equal")
	}
}

func TestSeqCopy(t *testing.T) {
	l1 := NewStrTestSeq()
	assert.Equal(t, l1, Copy(l1), "they should be equal")
}

func TestSeqSort(t *testing.T) {
	l1 := NewStrTestSeq()
	assert.NotEmpty(t, Reversed(l1), "should not be empty should be equal")
	assert.Equal(t, l1, Sort(Reversed(l1), helper.StrSort), "they should be equal")
}

func TestSeqRace(t *testing.T) {
	update := func(lst []int, s int) []int {
		a := 30
		for i := s * a; i < (s+1)*a; i++ {
			lst = Append(lst, i)
		}

		return lst
	}

	var mx sync.Mutex
	l1 := []int{}
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

	l2 := []int{}
	l2 = update(l2, 0)
	l2 = update(l2, 1)
	l2 = update(l2, 2)

	assert.Equal(t, l2, Sort(l1, helper.IntSort), "they should be equal")
}
