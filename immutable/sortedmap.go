package immutable

import (
	"bytes"
	"encoding/gob"
	"github.com/barweiss/go-tuple"
	"github.com/imunhatep/gocollection/dict"
	"github.com/imunhatep/gocollection/helper"
	"github.com/imunhatep/gocollection/slice"
	"github.com/samber/mo"
	"sort"
)

// SortedMap is a container for an optional Values of type V. If Values exists, HashMap is
// of type Some. If the Values is absent, HashMap is of type None.
type SortedMap[K comparable, V any] struct {
	index map[K]int
	pairs []tuple.T2[K, V]
}

// NewSortedMap Builds a HashMap
func NewSortedMap[K comparable, V any]() SortedMap[K, V] {
	mapVal := SortedMap[K, V]{
		pairs: []tuple.T2[K, V]{},
		index: map[K]int{},
	}

	return mapVal
}

func ToSortedMap[K comparable, V any](data map[K]V) SortedMap[K, V] {
	pairs := []tuple.T2[K, V]{}
	for k, v := range data {
		pairs = append(pairs, tuple.New2(k, v))
	}

	return SortedMap[K, V]{
		index: buildIndexFromPairs(pairs),
		pairs: pairs,
	}
}

func (o SortedMap[K, V]) SetKeys(keys ...K) SortedMap[K, V] {
	var empty V
	pairs := []tuple.T2[K, V]{}
	for _, k := range keys {
		pairs = append(pairs, tuple.New2(k, o.GetOrElse(k, empty)))
	}

	mapVal := SortedMap[K, V]{
		pairs: pairs,
		index: buildIndexFromPairs(pairs),
	}

	return mapVal
}

func (o SortedMap[K, V]) SetValues(values ...V) SortedMap[K, V] {
	var empty V
	pairs := []tuple.T2[K, V]{}
	for i, v := range o.pairs {
		pairs = append(pairs, tuple.New2(v.V1, slice.GetOrElse(values, i, empty)))
	}

	mapVal := SortedMap[K, V]{
		pairs: pairs,
		index: buildIndexFromPairs(pairs),
	}

	return mapVal
}

func newSortedMapWithList[K comparable, V any](values []tuple.T2[K, V]) SortedMap[K, V] {
	rez := SortedMap[K, V]{
		pairs: values,
		index: buildIndexFromPairs(values),
	}

	return rez
}

func (o SortedMap[K, V]) GobEncode() ([]byte, error) {
	store := bytes.NewBuffer([]byte{})
	encoder := gob.NewEncoder(store)
	err := encoder.Encode(o.pairs)

	return store.Bytes(), err
}

func (o *SortedMap[K, V]) GobDecode(data []byte) error {
	o.pairs = []tuple.T2[K, V]{}
	err := gob.NewDecoder(bytes.NewBuffer(data)).Decode(&o.pairs)
	o.index = buildIndexFromPairs(o.pairs)

	return err
}

func (o SortedMap[K, V]) Get(k K) mo.Option[tuple.T2[K, V]] {
	return slice.Get(o.pairs, dict.GetOrElse(o.index, k, -1))
}

func (o SortedMap[K, V]) GetOrElse(k K, def V) V {
	pair := slice.Get(o.pairs, dict.GetOrElse(o.index, k, -1))
	if pair.IsPresent() {
		return pair.MustGet().V2
	}

	return def
}

func (o SortedMap[K, V]) Head() mo.Option[tuple.T2[K, V]] {
	return slice.Head(o.pairs)
}

func (o SortedMap[K, V]) Tail() SortedMap[K, V] {
	return newSortedMapWithList(slice.Tail(o.pairs))
}

func (o SortedMap[K, V]) Join(m SortedMap[K, V]) SortedMap[K, V] {
	rez := append(o.pairs, m.ToSequence().ToSlice()...)
	return newSortedMapWithList(slice.UniqueAny(rez))
}

func (o SortedMap[K, V]) indexOf(key K) (int, bool) {
	if idx, ok := o.index[key]; ok {
		return idx, true
	}

	return -1, false
}

func (o SortedMap[K, V]) Update(key K, value V) SortedMap[K, V] {
	lst := slice.Copy(o.pairs)
	if idx, ok := o.indexOf(key); ok {
		lst[idx].V2 = value
	} else {
		lst = append(lst, tuple.New2(key, value))
	}

	return newSortedMapWithList(lst)
}

func (o SortedMap[K, V]) Remove(key K) SortedMap[K, V] {
	if !o.ContainsKey(key) {
		return o
	}

	// remove element
	lst := slice.Copy(o.pairs)
	idx, _ := o.indexOf(key)

	return newSortedMapWithList(append(lst[0:idx], lst[idx+1:]...))
}

func (o SortedMap[K, V]) Contains(s V) bool {
	exists := func(k K, v V) bool { return helper.CompareAny(v, s) }
	return o.Find(exists).IsPresent()
}

func (o SortedMap[K, V]) Find(f func(K, V) bool) mo.Option[tuple.T2[K, V]] {
	for _, p := range o.pairs {
		if f(p.V1, p.V2) {
			return mo.Some(p)
		}
	}

	return mo.None[tuple.T2[K, V]]()
}

func (o SortedMap[K, V]) ContainsKey(k K) bool {
	_, ok := o.index[k]
	return ok
}

// Size returns 1 when Values is present or 0 instead.
func (o SortedMap[K, V]) Size() int {
	return len(o.pairs)
}

func (o SortedMap[K, V]) IsEmpty() bool {
	return o.Size() == 0
}

func (o SortedMap[K, V]) Keys() Sequence[K] {
	keys := []K{}
	for _, p := range o.pairs {
		keys = append(keys, p.V1)
	}

	return NewSequence(keys...)
}

func (o SortedMap[K, V]) Values() Sequence[V] {
	values := []V{}
	for _, p := range o.pairs {
		values = append(values, p.V2)
	}

	return NewSequence(values...)
}

func (o SortedMap[K, V]) Reversed() SortedMap[K, V] {
	return newSortedMapWithList(slice.Reversed(o.pairs))
}

func (o SortedMap[K, V]) Sort(f func(V, V) bool) SortedMap[K, V] {
	lst := slice.Copy(o.pairs)
	sort.Slice(lst, func(p, q int) bool { return f(lst[p].V2, lst[q].V2) })
	return newSortedMapWithList(lst)
}

func (o SortedMap[K, V]) SortByKey(f func(K, K) bool) SortedMap[K, V] {
	lst := slice.Copy(o.pairs)
	sort.Slice(lst, func(p, q int) bool { return f(lst[p].V1, lst[q].V1) })
	return newSortedMapWithList(lst)
}

func (o SortedMap[K, V]) Filter(f func(K, V) bool) SortedMap[K, V] {
	rez := slice.Filter(o.pairs, func(p tuple.T2[K, V]) bool { return f(p.V1, p.V2) })
	return newSortedMapWithList(rez)
}

func (o SortedMap[K, V]) FoldLeft(z SortedMap[K, V], f func(SortedMap[K, V], K, V) SortedMap[K, V]) SortedMap[K, V] {
	for _, p := range o.pairs {
		z = f(z, p.V1, p.V2)
	}

	return z
}

func (o SortedMap[K, V]) FoldRight(z SortedMap[K, V], f func(SortedMap[K, V], K, V) SortedMap[K, V]) SortedMap[K, V] {
	return o.Reversed().FoldLeft(z, f)
}

// Map executes the mapper function if Values is present or returns None if absent.
func (o SortedMap[K, V]) Map(f func(K, V) V) SortedMap[K, V] {
	lst := []tuple.T2[K, V]{}
	for _, p := range o.pairs {
		lst = append(lst, tuple.New2(p.V1, f(p.V1, p.V2)))
	}

	return newSortedMapWithList(lst)
}

func (o SortedMap[K, V]) ToMap() map[K]V {
	rez := map[K]V{}
	for _, p := range o.pairs {
		rez[p.V1] = p.V2
	}

	return rez
}

func (o SortedMap[K, V]) ToSlice() []tuple.T2[K, V] {
	return slice.Copy(o.pairs)
}

func (o SortedMap[K, V]) ToSequence() Sequence[tuple.T2[K, V]] {
	return NewSequence(o.ToSlice()...)
}

func (o SortedMap[K, V]) getIndexMap() map[K]int {
	return o.index
}

func buildIndexFromPairs[K comparable, V any](values []tuple.T2[K, V]) map[K]int {
	index := map[K]int{}
	for i, v := range values {
		index[v.V1] = i
	}

	return index
}

func SortedMapFoldLeft[K comparable, V any, T any](
	src SortedMap[K, V],
	dst T,
	p func(T, K, V) T,
) T {
	for k, v := range src.ToMap() {
		dst = p(dst, k, v)
	}

	return dst
}
