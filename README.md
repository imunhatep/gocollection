Simple GoLang collections implementation library
================
This library make suse of go generics, a.k.a type parameters with `Option` pattern. 

### Slice and Dictionary helpers
 - functions/dict
 - functions/slice

## Requirements
 - go 1.19.0

## Install
```bash
go get -u github.com/imunhatep/gocollection
```

## Dict functions
List of functions for go `map`
```go
package main

import "github.com/imunhatep/gocollection/dict"

func Contains[K comparable, V any](data map[K]V, s V) bool {}
func ContainsKey[K comparable, V any](data map[K]V, k K) bool {}
func Copy[K comparable, V any](data map[K]V) map[K]V {}
func Filter[K comparable, V any](data map[K]V, f func(K, V) bool) map[K]V {}
func FilterNot[K comparable, V any](data map[K]V, f func(K, V) bool) map[K]V {}
func Find[K comparable, V any](data map[K]V, f func(K, V) bool) mo.Option[tuple.T2[K, V]] {}
func Fold[K comparable, V any, Z any](data map[K]V, z Z, f func(Z, K, V) Z) Z {}
func Get[K comparable, V any](data map[K]V, k K) mo.Option[tuple.T2[K, V]] {}
func GetOrElse[K comparable, V any](data map[K]V, k K, def V) V {}
func IsEmpty[K comparable, V any](data map[K]V) bool {}
func Keys[K comparable, V any](data map[K]V) []K {}
func Limit[K comparable, V any](data map[K]V, c int) map[K]V {}
func Map[K comparable, V any, Z any](data map[K]V, f func(K, V) Z) map[K]Z {}
func Merge[K comparable, V any](data map[K]V, m map[K]V) map[K]V {}
func Remove[K comparable, V any](data map[K]V, key K) map[K]V {}
func Set[K comparable, V any](data map[K]V, key K, value V) map[K]V {}
func Size[K comparable, V any](data map[K]V) int {}
func ToSlice[K comparable, V any](data map[K]V) []tuple.T2[K, V] {}
func Values[K comparable, V any](data map[K]V) []V {}
```

## Slice functions
List of functions for go `slice`
```go
func Append[V any](d1 []V, d2 ...V) []V {}
func ContainsAny[V any](data []V, elem V) bool {}
func Contains[V comparable](data []V, elem V) bool {}
func Copy[V any](data []V) []V {}
func Filter[V any](data []V, p func(V) bool) []V {}
func FilterNot[V any](data []V, p func(V) bool) []V {}
func FilterWithIndex[V any](data []V, p func(int, V) bool) []V {}
func Find[V any](data []V, p func(V) bool) mo.Option[V] {}
func FindWithIndex[V any](data []V, p func(int, V) bool) mo.Option[tuple.T2[int, V]] {}
func FoldLeft[V any, Z any](data []V, z Z, p func(Z, int, V) Z) Z {}
func FoldRight[V any, Z any](data []V, z Z, m func(Z, int, V) Z) Z {}
func GetOrElse[V any](data []V, i int, def V) V {}
func Get[V any](data []V, i int) mo.Option[V] {}
func Head[V any](data []V) mo.Option[V] {}
func IndexOfAny[V any](data []V, elem V) (int, bool) {}
func IndexOf[V comparable](data []V, elem V) (int, bool) {}
func IsEmpty[V any](data []V) bool {}
func Limit[V any](data []V, n int) []V {}
func Map[V any, Z any](data []V, p func(V) Z) []Z {}
func MapWithIndex[V any, Z any](data []V, p func(int, V) Z) []Z {}
func Reversed[V any](data []V) []V {}
func Size[V any](data []V) int {}
func Sort[V any](data []V, f func(V, V) bool) []V {}
func Tail[V any](data []V) []V {}
func UniqueAny[V any](data []V) []V {}
func Unique[V comparable](data []V) []V {}
```


## Examples
More examples could be found in go tests

Dictionary
```go
package main

import "github.com/imunhatep/gocollection/dict"

func NewIntTestMap(size int) map[string]int {
	values := map[string]int{}
	for _, i := range helper.Range(1, size) {
		values[fmt.Sprintf("key_%d", i)] = i
	}

	return values
}

func main() {
	l1 := NewIntTestMap(5)
	
	l2 := dict.Fold(
		l1,
		map[string]int{},
		func(z map[string]int, k string, v int) map[string]int {
			z[k] = v
			return z
		},
	)

	double := func(p int) int { return p * 2 }
	l2.Map(double).Head().OrElse(0)
}
```

Slices
```go
package main

import "github.com/imunhatep/gocollection/slice"

func main() {
	double := func(p int) int { return p * 2 }

	// map
	l1 := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
	r1 := slice.Map(l1, double)
}
```