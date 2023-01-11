package mutable

import (
	"bytes"
	"encoding/gob"
	"github.com/barweiss/go-tuple"
	"github.com/imunhatep/gocollection/dict"
	"github.com/samber/mo"
)

// HashMap is a container for an optional Values of type V. If Values exists, HashMap is
// of type Some. If the Values is absent, HashMap is of type None.
type HashMap[K comparable, V any] struct {
	items map[K]V
}

// NewMap Builds a HashMap
func NewMap[K comparable, V any]() HashMap[K, V] {
	mapVal := HashMap[K, V]{}
	mapVal.items = make(map[K]V)

	return mapVal
}

func ToMap[K comparable, V any](data map[K]V) HashMap[K, V] {
	return HashMap[K, V]{items: data}
}

func (o HashMap[K, V]) GobEncode() ([]byte, error) {
	store := bytes.NewBuffer([]byte{})
	encoder := gob.NewEncoder(store)
	err := encoder.Encode(o.items)

	return store.Bytes(), err
}

func (o *HashMap[K, V]) GobDecode(data []byte) error {
	o.items = map[K]V{}
	err := gob.NewDecoder(bytes.NewBuffer(data)).Decode(&o.items)

	return err
}

func (o HashMap[K, V]) Get(k K) mo.Option[tuple.T2[K, V]] {
	return dict.Get(o.items, k)
}

func (o HashMap[K, V]) GetOrElse(k K, def V) V {
	return dict.GetOrElse(o.items, k, def)
}

func (o HashMap[K, V]) Join(m HashMap[K, V]) HashMap[K, V] {
	o.items = dict.Merge(o.items, m.ToMap())
	return o
}

func (o HashMap[K, V]) Update(key K, value V) HashMap[K, V] {
	o.items = dict.Set(o.items, key, value)
	return o
}

func (o HashMap[K, V]) Remove(key K) HashMap[K, V] {
	o.items = dict.Remove(o.items, key)
	return o
}

func (o HashMap[K, V]) Contains(s V) bool {
	return dict.Contains(o.items, s)
}

func (o HashMap[K, V]) Find(f func(K, V) bool) mo.Option[tuple.T2[K, V]] {
	return dict.Find(o.items, f)
}

func (o HashMap[K, V]) HasKey(k K) bool {
	return dict.ContainsKey(o.items, k)
}

// Size returns 1 when Values is present or 0 instead.
func (o HashMap[K, V]) Size() int {
	return dict.Size(o.items)
}

func (o HashMap[K, V]) IsEmpty() bool {
	return dict.IsEmpty(o.items)
}

func (o HashMap[K, V]) Keys() Sequence[K] {
	return NewSequence(dict.Keys(o.items)...)
}

func (o HashMap[K, V]) Values() Sequence[V] {
	return NewSequence(dict.Values(o.items)...)
}

func (o HashMap[K, V]) Filter(f func(K, V) bool) HashMap[K, V] {
	o.items = dict.Filter(o.items, f)
	return o
}

func (o HashMap[K, V]) FoldLeft(z HashMap[K, V], f func(HashMap[K, V], K, V) HashMap[K, V]) HashMap[K, V] {
	return dict.Fold(o.items, z, f)
}

// HashMap executes the mapper function if Values is present or returns None if absent.
func (o HashMap[K, V]) Map(f func(K, V) V) HashMap[K, V] {
	o.items = dict.Map(o.items, f)
	return o
}

func (o HashMap[K, V]) ToMap() map[K]V {
	return dict.Copy(o.items)
}

func (o HashMap[K, V]) ToSlice() []tuple.T2[K, V] {
	return dict.ToSlice(o.items)
}

func (o HashMap[K, V]) ToSequence() Sequence[tuple.T2[K, V]] {
	return NewSequence(o.ToSlice()...)
}

func HashMapFoldLeft[K comparable, V any, T any](
	src HashMap[K, V],
	dst T,
	p func(T, K, V) T,
) T {
	for k, v := range src.ToMap() {
		dst = p(dst, k, v)
	}

	return dst
}
