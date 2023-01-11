Simple GoLang collections implementation library
================
This library make suse of go generics, a.k.a type parameters with `Option` pattern. 

## Implemented structures
 - functions/dict
 - functions/slice
 - immutable/hashmap
 - immutable/sortedmap
 - immutable/sequence
 - mutable/hashmap
 - mutable/sequence

## Dict
List of functions for go `map`
```go
func Keys[K comparable, V any](data map[K]V) []K {}
func Values[K comparable, V any](data map[K]V) []V {}
func Copy[K comparable, V any](data map[K]V) map[K]V {}
func Get[K comparable, V any](data map[K]V, k K) mo.Option[tuple.T2[K, V]] {} 
func GetOrElse[K comparable, V any](data map[K]V, k K, def V) V {}
func Merge[K comparable, V any](data map[K]V, m map[K]V) map[K]V {}
func Set[K comparable, V any](data map[K]V, key K, value V) map[K]V {} 
func Remove[K comparable, V any](data map[K]V, key K) map[K]V {}
func Find[K comparable, V any](data map[K]V, f func(K, V) bool) mo.Option[tuple.T2[K, V]] {} 
func Contains[K comparable, V any](data map[K]V, s V) bool {}
func ContainsKey[K comparable, V any](data map[K]V, k K) bool {}
func Size[K comparable, V any](data map[K]V) int {}
func IsEmpty[K comparable, V any](data map[K]V) bool {} 
func Filter[K comparable, V any](data map[K]V, f func(K, V) bool) map[K]V {} 
func Fold[K comparable, V any, Z any](data map[K]V, z Z, f func(Z, K, V) Z) Z {} 
func Map[K comparable, V any, Z any](data map[K]V, f func(K, V) Z) map[K]Z {}
func ToSlice[K comparable, V any](data map[K]V) []tuple.T2[K, V] {} 
```

## Slice
List of functions for go `slice`
```go
func Get[V any](data []V, i int) mo.Option[V] {} 
func GetOrElse[V any](data []V, i int, def V) V {} 
func IsEmpty[V any](data []V) bool {} 
func Sort[V any](data []V, f func(V, V) bool) []V {} 
func Append[V any](d1 []V, d2 ...V) []V {} 
func Copy[V any](data []V) []V {} 
func Head[V any](data []V) mo.Option[V] {} 
func Tail[V any](data []V) []V {}
func Find[V any](data []V, p func(V) bool) mo.Option[V] {} 
func FindWithIndex[V any](data []V, p func(int, V) bool) mo.Option[tuple.T2[int, V]] {} 
func Map[V any, Z any](data []V, p func(V) Z) []Z {} 
func MapWithIndex[V any, Z any](data []V, p func(int, V) Z) []Z {} 
func Reversed[V any](data []V) []V {} 
func Filter[V any](data []V, p func(V) bool) []V {} 
func FilterWithIndex[V any](data []V, p func(int, V) bool) []V {} 
func FoldLeft[V any, Z any](data []V, z Z, p func(Z, int, V) Z) Z {}
func FoldRight[V any, Z any](data []V, z Z, m func(Z, int, V) Z) Z {} 
func Size[V any](data []V) int {}
func IndexOf[V comparable](data []V, elem V) (int, bool) {} 
func IndexOfAny[V any](data []V, elem V) (int, bool) {} 
func Contains[V comparable](data []V, elem V) bool {} 
func ContainsAny[V any](data []V, elem V) bool {} 
func Unique[V comparable](data []V) []V {} 
func UniqueAny[V any](data []V) []V {}
```


## Examples
More examples could be found in go tests

```go
func NewStrTestMap(size int) HashMap[string, string] {
    values := map[string]string{}
    for _, i := range helper.Range(1, size) {
        values[fmt.Sprintf("key_%d", i)] = fmt.Sprintf("value_%d", i)
    }

    return ToMap(values)
}

l1 := NewStrTestMap(5)
l2 := l1.Remove(l1.Keys().Head().MustGet())
```