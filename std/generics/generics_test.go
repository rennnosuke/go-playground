package generics

import "testing"

type Number interface {
	int64 | int
}

type CompatibleNumber interface {
	~int64 | ~int
}

type Int int
type Int64 int64

func Sum[T Number](values ...T) T {
	var sum T
	for _, i := range values {
		sum += i
	}
	return sum
}

func SumCompatible[T CompatibleNumber](values ...T) T {
	var sum T
	for _, i := range values {
		sum += i
	}
	return sum
}

func TestSum(t *testing.T) {
	t.Log(Sum([]int64{1, 2, 3}...)) // ok
	t.Log(Sum([]int{1, 2, 3}...))   // ok
	//t.Log(Sum([]Int{1, 2, 3}...)) // ng
	//t.Log(Sum([]Int64{1, 2, 3}...))   // ng

	t.Log(SumCompatible([]int64{1, 2, 3}...)) // ok
	t.Log(SumCompatible([]int{1, 2, 3}...))   // ok
	t.Log(SumCompatible([]Int{1, 2, 3}...))   // ok
	t.Log(SumCompatible([]Int64{1, 2, 3}...)) // ok
}
