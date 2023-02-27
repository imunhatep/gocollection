package immutable

import (
	"github.com/barweiss/go-tuple"
	"github.com/imunhatep/gocollection/slice"
	"github.com/samber/mo"
)

// Sequence is a container for an optional items of type V. If items exists, Sequence is
// of type V. If the items is absent, Sequence is of type None.
type Sequence[V any] struct {
	items []V
}

func NewSequence[V any](values ...V) Sequence[V] {
	return Sequence[V]{values}
}

func (o Sequence[V]) Append(values ...V) Sequence[V] {
	res := append(o.items, values...)

	return NewSequence[V](res...)
}

func (o Sequence[V]) Join(seq Sequence[V]) Sequence[V] {
	return o.Append(seq.ToSlice()...)
}

// Head returns first element as an Option.
func (o Sequence[V]) Head() mo.Option[V] {
	return slice.Head(o.items)
}

// Tail returns last elements as an Sequence[V].
func (o Sequence[V]) Tail() Sequence[V] {
	return NewSequence[V](slice.Tail(o.items)...)
}

func (o Sequence[V]) Find(f func(V) bool) mo.Option[V] {
	return slice.Find(o.items, f)
}

func (o Sequence[V]) FindWithIndex(f func(int, V) bool) mo.Option[tuple.T2[int, V]] {
	return slice.FindWithIndex(o.items, f)
}

func (o Sequence[V]) Map(f func(V) V) Sequence[V] {
	return NewSequence(slice.Map(o.items, f)...)
}

func (o Sequence[V]) MapWithIndex(f func(int, V) V) Sequence[V] {
	return NewSequence(slice.MapWithIndex(o.items, f)...)
}

func (o Sequence[V]) Reversed() Sequence[V] {
	return NewSequence[V](slice.Reversed(o.items)...)
}

func (o Sequence[V]) Filter(p func(V) bool) Sequence[V] {
	return NewSequence(slice.Filter(o.items, p)...)
}

func (o Sequence[V]) FilterNot(p func(V) bool) Sequence[V] {
	return NewSequence(slice.FilterNot(o.items, p)...)
}

func (o Sequence[V]) FilterWithIndex(p func(int, V) bool) Sequence[V] {
	return NewSequence(slice.FilterWithIndex(o.items, p)...)
}

func (o Sequence[V]) FoldLeft(z Sequence[V], p func(Sequence[V], int, V) Sequence[V]) Sequence[V] {
	return slice.FoldLeft(o.items, z, p)
}

func (o Sequence[V]) FoldRight(z Sequence[V], p func(Sequence[V], int, V) Sequence[V]) Sequence[V] {
	return slice.FoldRight(o.items, z, p)
}

// Size returns len() of items.
func (o Sequence[V]) Size() int {
	return slice.Size(o.items)
}

func (o Sequence[V]) Limit(c int) Sequence[V] {
	return NewSequence(slice.Limit(o.items, c)...)
}

// IndexOf Search the sorted for a given value and return position, if not found returns -1
func (o Sequence[V]) IndexOf(e V) (int, bool) {
	return slice.IndexOfAny(o.items, e)
}

// Contains Tests whether this sorted contains a given items as an element.
func (o Sequence[V]) Contains(e V) bool {
	return slice.ContainsAny(o.items, e)
}

func (o Sequence[V]) Unique() Sequence[V] {
	return NewSequence(slice.UniqueAny(o.items)...)
}

func (o Sequence[V]) Index(i int) mo.Option[V] {
	return slice.Get(o.items, i)
}

func (o Sequence[V]) Sort(f func(v1, v2 V) bool) Sequence[V] {
	return NewSequence(slice.Sort(o.items, f)...)
}

func (o Sequence[V]) IsEmpty() bool {
	return slice.IsEmpty(o.items)
}

func (o Sequence[V]) ToSlice() []V {
	return o.items
}

func SeqFoldLeft[T1 any, T2 any](
	src Sequence[T1],
	dst T2,
	p func(T2, int, T1) T2,
) T2 {
	for i, v := range src.ToSlice() {
		dst = p(dst, i, v)
	}

	return dst
}
