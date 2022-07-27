package slice

import (
	"fmt"
	"testing"
)

func TestSliceTrick(t *testing.T) {
	// concat
	a := []int{1, 2, 3, 4, 5}
	b := []int{6, 7, 8, 9, 10}
	fmt.Println(append(a, b...)) // [1 2 3 4 5 6 7 8 9 10]

	// copy
	a = []int{1, 2, 3, 4, 5}
	b = make([]int, len(a))
	c := make([]int, 0, len(a))
	copy(b, a)
	copy(c, a)
	fmt.Println(b) // [1 2 3 4 5]
	fmt.Println(c) // これはコピーされない

	// deletion in range
	a = []int{1, 2, 3, 4, 5}
	i, j := 2, 3
	fmt.Println(append(a[:i], a[j:]...)) // [1 2 4 5]

	// deletion an element
	a = []int{1, 2, 3, 4, 5}
	fmt.Println(append(a[:i], a[i+1:]...)) // [1 2 4 5]

	a = []int{1, 2, 3, 4, 5}
	a = a[:i+copy(a[i:], a[i+1:])]
	fmt.Println(a) // [1 2 4 5]

}
