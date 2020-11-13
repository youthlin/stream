# Go Stream
[![PkgGoDev](https://pkg.go.dev/badge/github.com/youthlin/stream)](https://pkg.go.dev/github.com/youthlin/stream)
[![Go Report Card](https://goreportcard.com/badge/github.com/youthlin/stream)](https://goreportcard.com/report/github.com/youthlin/stream)
[![Build Status](https://travis-ci.org/youthlin/stream.svg?branch=main)](https://travis-ci.org/youthlin/stream)
[![codecov](https://codecov.io/gh/youthlin/stream/branch/main/graph/badge.svg?token=1CqmLWbsYL)](https://codecov.io/gh/youthlin/stream)

Go Stream, like Java 8 Stream.

Blog Post: https://youthlin.com/?p=1755

## How to get
```shell script
go get github.com/youthlin/stream
```

国内镜像: https://gitee.com/youthlin/stream  
在 `go.mod` 中引入模块路径 `github.com/youthlin/stream` 及版本后，  
再添加 replace 即可：
```go
// go.mod

require github.com/youthlin/stream latest

replace github.com/youthlin/stream latest => gitee.com/youthlin/stream latest

```

## Play online
https://play.golang.org/p/vO8NEkdNXzY
```go
package main

import (
	"fmt"

	"github.com/youthlin/stream"
	"github.com/youthlin/stream/types"
)

func main() {
	m := stream.IntRange(0, 10).
		Filter(func(e types.T) bool {
			return e.(int)%2 == 0
		}).
		Map(func(e types.T) types.R {
			return e.(int) * 2
		}).
		ReduceWith(map[int]string{}, func(m types.R, e types.T) types.R {
			m.(map[int]string)[e.(int)] = fmt.Sprintf("<%d>", e)
			return m
		})
	fmt.Println(m)
}

```

## Examples
```go

type Stream interface {
	// stateless operate 无状态操作

	Filter(types.Predicate) Stream         // 过滤
	Map(types.Function) Stream             // 转换
	FlatMap(func(t types.T) Stream) Stream // 打平
	Peek(types.Consumer) Stream            // peek 每个元素

	// stateful operate 有状态操作

	Distinct(types.IntFunction) Stream // 去重
	Sorted(types.Comparator) Stream    // 排序
	Limit(int64) Stream                // 限制个数
	Skip(int64) Stream                 // 跳过个数

	// terminal operate 终止操作

	// 遍历
	ForEach(types.Consumer)
	// return []T 转为切片
	ToSlice() []types.T
	// return []X which X is the type of some
	ToElementSlice(some types.T) types.R
	// return []X which X is same as the `typ` representation
	ToSliceOf(typ reflect.Type) types.R
	// 测试是否所有元素满足条件
	AllMatch(types.Predicate) bool
	// 测试是否没有元素满足条件
	NoneMatch(types.Predicate) bool
	// 测试是否有任意元素满足条件
	AnyMatch(types.Predicate) bool
	// Reduce return optional.Empty if no element. calculate result by (T, T) -> T from first element
	Reduce(accumulator types.BinaryOperator) optional.Optional
	// type of initValue is same as element.  (T, T) -> T
	ReduceFrom(initValue types.T, accumulator types.BinaryOperator) types.T
	// type of initValue is different from element. (R, T) -> R
	ReduceWith(initValue types.R, accumulator func(types.R, types.T) types.R) types.R
	FindFirst() optional.Optional
	// 返回元素个数
	Count() int64
}


func ExampleOf() {
	fmt.Println(stream.Of().Count())
	fmt.Println(stream.Of(1).Count())
	fmt.Println(stream.Of("a", "b").Count())
	var s = []int{1, 2, 3, 4}
	stream.Of(stream.Slice(s)...).ForEach(func(t types.T) {
		fmt.Printf("%d,", t)
	})
	// Output:
	// 0
	// 1
	// 2
	// 1,2,3,4,
}
func ExampleStream_Filter() {
	stream.Of(0, 1, 2, 3, 4, 5, 6, 7, 8, 9).
		Filter(func(e types.T) bool {
			return e.(int)%3 == 0
		}).
		ForEach(func(e types.T) {
			fmt.Println(e)
		})
	// Output:
	// 0
	// 3
	// 6
	// 9
}
func ExampleStream_Map() {
	stream.IntRange(0, 5).
		Map(func(t types.T) types.R {
			return fmt.Sprintf("<%d>", t.(int))
		}).
		ForEach(func(t types.T) {
			fmt.Printf("%v", t)
		})
	// Output:
	// <0><1><2><3><4>
}
func ExampleStream_FlatMap() {
	stream.Of([]int{0, 2, 4, 6, 8}, []int{1, 3, 5, 7, 9}).
		FlatMap(func(t types.T) stream.Stream {
			return stream.Of(stream.Slice(t)...)
		}).
		ForEach(func(t types.T) {
			fmt.Printf("%d", t)
		})
	// Output:
	// 0246813579
}
func ExampleStream_Sorted() {
	stream.IntRange(1, 10).
		Sorted(types.ReverseOrder(types.IntComparator)).
		ForEach(func(t types.T) {
			fmt.Printf("%d,", t)
		})
	// Output:
	// 9,8,7,6,5,4,3,2,1,
}
func TestToMap(t *testing.T) {
	m := stream.IntRange(0, 10).ReduceWith(make(map[int]int), func(acc types.R, t types.T) types.R {
		acc.(map[int]int)[t.(int)] = t.(int) * 10
		return acc
	})
	t.Log(m)
	// Output:
	// map[0:0 1:10 2:20 3:30 4:40 5:50 6:60 7:70 8:80 9:90]
}

```
