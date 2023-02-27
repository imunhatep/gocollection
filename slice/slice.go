package slice

import (
	"github.com/barweiss/go-tuple"
	"github.com/imunhatep/gocollection/helper"
	"github.com/samber/mo"
	"sort"
)

func Get[V any](data []V, i int) mo.Option[V] {
	if i >= 0 && len(data) > i {
		return mo.Some(data[i])
	}

	return mo.None[V]()
}

func GetOrElse[V any](data []V, i int, def V) V {
	if i >= 0 && len(data) > i {
		return data[i]
	}

	return def
}

func IsEmpty[V any](data []V) bool {
	return len(data) == 0
}

func Sort[V any](data []V, f func(V, V) bool) []V {
	rez := Copy(data)
	sort.Slice(rez, func(p, q int) bool { return f(rez[p], rez[q]) })

	return rez
}

func Append[V any](d1 []V, d2 ...V) []V {
	return append(d1, d2...)
}

func Copy[V any](data []V) []V {
	rez := make([]V, len(data))
	copy(rez, data)

	return rez
}

// Head returns first element as an Option.
func Head[V any](data []V) mo.Option[V] {
	if len(data) == 0 {
		return mo.None[V]()
	}

	return mo.Some(data[0])
}

// Tail returns last elements as an []V.
func Tail[V any](data []V) []V {
	if len(data) == 0 {
		return []V{}
	}

	cpy := Copy(data)
	return cpy[1:]
}

func Find[V any](data []V, p func(V) bool) mo.Option[V] {
	for _, v := range data {
		if p(v) {
			return mo.Some(v)
		}
	}

	return mo.None[V]()
}

func FindWithIndex[V any](data []V, p func(int, V) bool) mo.Option[tuple.T2[int, V]] {
	for i, v := range data {
		if p(i, v) {
			return mo.Some(tuple.New2(i, v))
		}
	}

	return mo.None[tuple.T2[int, V]]()
}

func Map[V any, Z any](data []V, p func(V) Z) []Z {
	rez := make([]Z, 0, len(data))
	for _, v := range data {
		rez = append(rez, p(v))
	}

	return rez
}

func MapWithIndex[V any, Z any](data []V, p func(int, V) Z) []Z {
	rez := make([]Z, 0, len(data))
	for i, v := range data {
		rez = append(rez, p(i, v))
	}

	return rez
}

func Reversed[V any](data []V) []V {
	rez := make([]V, 0, len(data))
	for i := len(data) - 1; i > -1; i-- {
		rez = append(rez, data[i])
	}

	return rez
}

func Filter[V any](data []V, p func(V) bool) []V {
	var rez []V
	for _, v := range data {
		if p(v) {
			rez = append(rez, v)
		}
	}

	return rez
}

func FilterNot[V any](data []V, f func(V) bool) []V {
	fn := func(v V) bool { return !f(v) }
	return Filter(data, fn)
}

func FilterWithIndex[V any](data []V, p func(int, V) bool) []V {
	var rez []V
	for i, v := range data {
		if p(i, v) {
			rez = append(rez, v)
		}
	}

	return rez
}

func FoldLeft[V any, Z any](data []V, z Z, p func(Z, int, V) Z) Z {
	for i, v := range data {
		z = p(z, i, v)
	}

	return z
}

func FoldRight[V any, Z any](data []V, z Z, m func(Z, int, V) Z) Z {
	return FoldLeft(Reversed(data), z, m)
}

// Size returns len() of items.
func Size[V any](data []V) int {
	return len(data)
}

// Limit returns first N elements.
func Limit[V any](data []V, c int) []V {
	return data[0:helper.Min(c, Size(data))]
}

// IndexOf Search the list for a given value and return position, if not found returns -1
func IndexOf[V comparable](data []V, elem V) (int, bool) {
	rez := FindWithIndex(data, func(i int, v V) bool { return v == elem })
	if rez.IsPresent() {
		return rez.MustGet().V1, true
	}

	return -1, false
}

func IndexOfAny[V any](data []V, elem V) (int, bool) {
	rez := FindWithIndex(data, func(i int, v V) bool { return helper.CompareAny(v, elem) })
	if rez.IsPresent() {
		return rez.MustGet().V1, true
	}

	return -1, false
}

// Contains Tests whether this list contains a given items as an element.
func Contains[V comparable](data []V, elem V) bool {
	rez := Find(data, func(v V) bool { return v == elem })
	return rez.IsPresent()
}

func ContainsAny[V any](data []V, elem V) bool {
	rez := Find(data, func(v V) bool { return helper.CompareAny(v, elem) })
	return rez.IsPresent()
}

func Unique[V comparable](data []V) []V {
	u := make([]V, 0, len(data))
	m := make(map[V]bool)

	for _, val := range data {
		if _, ok := m[val]; !ok {
			m[val] = true
			u = append(u, val)
		}
	}

	return u
}

func UniqueAny[V any](data []V) []V {
	contains := func(s []V, n V) bool {
		for _, v := range s {
			if helper.CompareAny(v, n) {
				return true
			}
		}

		return false
	}

	var rez []V
	for _, v := range Reversed(data) {
		if !contains(rez, v) {
			rez = append(rez, v)
		}
	}

	return Reversed(rez)
}
