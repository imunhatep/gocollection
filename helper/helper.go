package helper

import (
	"golang.org/x/exp/constraints"
	"reflect"
)

func Range(s, e int) (rez []int) {
	for ; s < e; s++ {
		rez = append(rez, s)
	}

	return rez
}

func StrSort(s1, s2 string) bool {
	return s1 < s2
}

func IntSort(s1, s2 int) bool {
	return s1 < s2
}

func CompareAny[T any](t1 T, t2 T) bool {
	if reflect.TypeOf((*T)(nil)).Elem().Comparable() {
		if any(t1) == any(t2) {
			return true
		}

		return false
	}

	return reflect.DeepEqual(t1, t2)
}

func Min[T constraints.Ordered](a, b T) T {
	if a < b {
		return a
	}
	return b
}
