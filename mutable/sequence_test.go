package mutable

import (
	"github.com/imunhatep/gocollection/helper"
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
)

type testStruct struct {
	Some string
}

func NewStrTestSeq() Sequence[string] {
	return NewSequence([]string{"test1", "test2", "test3"}...)
}

func NewIntTestSeq() Sequence[int] {
	return NewSequence([]int{1, 2, 3, 4, 5, 6, 7, 8, 9}...)
}

func TestSeqHeadTail(t *testing.T) {
	l1 := NewStrTestSeq()

	// Head / Tail
	for _, v := range l1.ToSlice() {
		// assert equality
		assert.Equal(t, l1.Head().MustGet(), v, "they should be equal")
		l1 = l1.Tail()
	}
}

func TestSeqMap(t *testing.T) {
	l1 := NewIntTestSeq()

	double := func(p int) int { return p * 2 }
	// map
	r1 := l1.Map(double)
	for i, v := range l1.ToSlice() {
		assert.Equal(t, double(v), r1.Index(i).MustGet(), "they should be equal")
	}
}

func TestSeqUnique(t *testing.T) {
	l1 := NewIntTestSeq()
	l2 := l1.Join(NewIntTestSeq()).Unique()

	assert.Equal(t, l1.ToSlice(), l2.ToSlice(), "they should be equal")
}

func TestSeqCompare(t *testing.T) {
	l1 := NewStrTestSeq()
	v1 := l1.Head().MustGet()

	assert.True(t, l1.Contains(v1), "seq must contain a value")

	i1 := l1.Find(func(v string) bool { return v == v1 }).OrEmpty()
	assert.Equal(t, v1, i1, "they should be equal")

	e1 := NewSequence[string]().Find(func(v string) bool { return v == "empty" })
	assert.Empty(t, e1.OrEmpty())

	i2 := l1.FindWithIndex(func(i int, v string) bool { return v == v1 }).OrEmpty()
	assert.Equal(t, v1, i2.V2, "they should be equal")

	e2 := NewSequence[string]().FindWithIndex(func(i int, v string) bool { return v == "empty" })
	assert.Empty(t, e2.OrEmpty())

	s1 := SeqFoldLeft(
		l1,
		NewSequence[testStruct](),
		func(z Sequence[testStruct], i int, v string) Sequence[testStruct] { return z.Append(testStruct{v}) },
	)
	sv1 := s1.Head().MustGet()

	assert.True(t, s1.Contains(sv1))
}

func TestSeqFilter(t *testing.T) {
	l1 := NewIntTestSeq()
	l2Size := l1.Size() - 4
	l2 := l1.FilterWithIndex(func(i int, v int) bool { return i < l2Size })

	assert.Equal(t, l2Size, l2.Size(), "they should be equal")
}

func TestSeqFolding(t *testing.T) {
	l1 := NewIntTestSeq()
	l2 := l1.FoldRight(NewSequence[int](), func(s Sequence[int], i int, v int) Sequence[int] { return s.Append(v) })

	assert.Equal(t, l1.Reversed().ToSlice(), l2.ToSlice(), "they should be equal")
}

func TestSeqEmpty(t *testing.T) {
	l1 := NewSequence[int]()

	assert.True(t, l1.IsEmpty())
	assert.Empty(t, l1.Tail().ToSlice())
	assert.Empty(t, l1.Head().OrEmpty())
	assert.Empty(t, l1.Index(0).OrEmpty())
	assert.False(t, l1.Contains(0))
}

func TestSeqReversed(t *testing.T) {
	l1 := NewStrTestSeq()

	// Reversed
	r1 := l1.Reversed()
	for i, v := range r1.Reversed().ToSlice() {
		idx, _ := l1.IndexOf(v)
		assert.Equal(t, i, idx, "they should be equal")
	}
}

func TestSeqSort(t *testing.T) {
	l1 := NewStrTestSeq()
	assert.Equal(t, l1.ToSlice(), l1.Reversed().Sort(helper.StrSort).ToSlice(), "they should be equal")
}

func TestSeqRace(t *testing.T) {
	update := func(lst Sequence[int], s int) Sequence[int] {
		a := 30
		for i := s * a; i < (s+1)*a; i++ {
			lst = lst.Append(i)
		}

		return lst
	}

	var mx sync.Mutex
	l1 := NewSequence[int]()
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

	l2 := NewSequence[int]()
	l2 = update(l2, 0)
	l2 = update(l2, 1)
	l2 = update(l2, 2)

	assert.Equal(t, l2.ToSlice(), l1.Sort(helper.IntSort).ToSlice(), "they should be equal")
}
