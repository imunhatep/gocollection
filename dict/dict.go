package dict

import (
	"github.com/barweiss/go-tuple"
	"github.com/imunhatep/gocollection/helper"
	"github.com/samber/mo"
)

func Keys[K comparable, V any](data map[K]V) []K {
	keys := []K{}
	for k, _ := range data {
		keys = append(keys, k)
	}

	return keys
}

func Values[K comparable, V any](data map[K]V) []V {
	values := []V{}
	for _, v := range data {
		values = append(values, v)
	}

	return values
}

func Copy[K comparable, V any](data map[K]V) map[K]V {
	cpy := map[K]V{}
	for k, v := range data {
		cpy[k] = v
	}
	return cpy
}

func Get[K comparable, V any](data map[K]V, k K) mo.Option[tuple.T2[K, V]] {
	if v, ok := data[k]; ok {
		return mo.Some(tuple.New2(k, v))
	}

	return mo.None[tuple.T2[K, V]]()
}

func GetOrElse[K comparable, V any](data map[K]V, k K, def V) V {
	if v, ok := data[k]; ok {
		return v
	}

	return def
}

func Merge[K comparable, V any](data map[K]V, m map[K]V) map[K]V {
	rez := Copy(data)
	for k, v := range m {
		rez[k] = v
	}

	return rez
}

func Set[K comparable, V any](data map[K]V, key K, value V) map[K]V {
	rez := Copy(data)
	rez[key] = value

	return rez
}

func Remove[K comparable, V any](data map[K]V, key K) map[K]V {
	if !ContainsKey(data, key) {
		return data
	}

	rez := Copy(data)
	delete(rez, key)

	return rez
}

func Find[K comparable, V any](data map[K]V, f func(K, V) bool) mo.Option[tuple.T2[K, V]] {
	for k, v := range data {
		if f(k, v) {
			return mo.Some(tuple.New2(k, v))
		}
	}

	return mo.None[tuple.T2[K, V]]()
}

func Contains[K comparable, V any](data map[K]V, s V) bool {
	exists := func(k K, v V) bool { return helper.CompareAny(v, s) }
	return Find(data, exists).IsPresent()
}

func ContainsKey[K comparable, V any](data map[K]V, k K) bool {
	_, ok := data[k]
	return ok
}

func Size[K comparable, V any](data map[K]V) int {
	return len(data)
}

func IsEmpty[K comparable, V any](data map[K]V) bool {
	return Size(data) == 0
}

func Filter[K comparable, V any](data map[K]V, f func(K, V) bool) map[K]V {
	rez := map[K]V{}
	for k, v := range data {
		if f(k, v) {
			rez[k] = v
		}
	}

	return rez
}

func FilterNot[K comparable, V any](data map[K]V, f func(K, V) bool) map[K]V {
	fn := func(k K, v V) bool { return !f(k, v) }
	return Filter(data, fn)
}

func Fold[K comparable, V any, Z any](data map[K]V, z Z, f func(Z, K, V) Z) Z {
	for k, v := range data {
		z = f(z, k, v)
	}

	return z
}

func Map[K comparable, V any, Z any](data map[K]V, f func(K, V) Z) map[K]Z {
	rez := map[K]Z{}
	for k, v := range data {
		rez[k] = f(k, v)
	}

	return rez
}

func Limit[K comparable, V any](data map[K]V, l int) map[K]V {
	rez := map[K]V{}
	for k, v := range data {
		rez[k] = v

		l -= 1
		if l <= 0 {
			break
		}
	}

	return rez
}

func ToSlice[K comparable, V any](data map[K]V) []tuple.T2[K, V] {
	rez := []tuple.T2[K, V]{}
	for k, v := range data {
		rez = append(rez, tuple.New2(k, v))
	}

	return rez
}
