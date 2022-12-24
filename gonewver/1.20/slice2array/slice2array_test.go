package main

import "testing"

func TestSliceToArray(t *testing.T) {
	sl := make([]int, 3, 5)

	// before 1.20
	// sliceはarray参照にキャストできる
	// 但し、キャスト先の配列長 < len(slice)でないとruntime errorとなる（capは関係なし）
	arr := *(*[3]int)(sl)
	// arr := *(*[2]int)(sl) // ok
	// arr := *(*[5]int)(sl) // runtime error
	t.Log(arr)

	// after 1.20
	// array参照を経由しなくても、そのままキャストできる
	// newCastArr := [4]int(sl)
	// t.Log(newCastArr)
}
