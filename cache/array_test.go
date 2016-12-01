package cache

import (
	"fmt"
	"testing"
)

func TestArray(t *testing.T) {
	a := NewArray(2)
	a.Push("王海涛")
	a.Push("吴静涛")
	a.Push("周龙伟")
	a.Unshift("李靖秋")
	fmt.Println(a.Length())

	fmt.Println("抛出", a.Pop())
	fmt.Println("抛出", a.Pop())
	fmt.Println("抛出", a.Pop())
	fmt.Println("抛出", a.Pop())
}
