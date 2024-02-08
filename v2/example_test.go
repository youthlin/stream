package stream_test

import (
	"fmt"
	"iter"
	"strings"

	"github.com/youthlin/stream/v2"
	"github.com/youthlin/stream/v2/types"
)

func ExampleOf() {
	fmt.Printf("Of(): %v\n", stream.Of[int]().Count())
	fmt.Printf("Of(1): %v\n", stream.Of(1).Count())
	var s = []int{1, 2, 3, 4}
	fmt.Printf("Of(...): %v\n", stream.Of(s...).Count())
	stream.Of(s...).Skip(1).Limit(2).ForEach(func(i int) {
		fmt.Printf("%d,", i)
	})
	// Output:
	// Of(): 0
	// Of(1): 1
	// Of(...): 4
	// 2,3,
}

func ExampleCountFrom() {
	s := stream.CountFrom(0).Limit(4).Peek(func(i int) {
		fmt.Printf("%d;", i)
	}).Skip(2).Collect()
	fmt.Println(s)
	// Output:
	// 0;1;2;3;[2 3]
}

func ExampleRange() {
	fmt.Println(stream.Range(0, 10).Limit(9).Count())
	stream.RangeStep(5, 0, -1).Limit(4).ForEach(func(i int) {
		fmt.Printf("%d,", i)
	})
	// Output:
	// 9
	// 5,4,3,2,
}

func ExampleRepeat() {
	fmt.Println(strings.Join(stream.Repeat("a").Limit(5).Collect(), ""))
	// Output:
	// aaaaa
}

func fib() types.Supplier[int] {
	a := -1
	b := 1
	return func() int {
		n := a + b
		a = b
		b = n
		return n
	}
}
func ExampleGenerate() {
	stream.Generate(fib()).Limit(10).ForEach(func(i int) {
		fmt.Printf("%d,", i)
	})
	// Output:
	// 0,1,1,2,3,5,8,13,21,34,
}

func ExampleStream() {
	stream.Generate(fib()).Filter(func(i int) bool {
		if i%2 == 0 {
			fmt.Printf("[%dok]", i)
			return true
		} else {
			fmt.Printf("[%dno]", i)
			return false
		}
	}).Peek(func(i int) {
		fmt.Printf("<%d>", i)
	}).Map(func(i int) int {
		return i * 10
	}).Limit(3).ForEach(func(i int) {
		fmt.Printf("got=%d,", i)
	})
	// Output:
	// [0ok]<0>got=0,[1no][1no][2ok]<2>got=20,[3no][5no][8ok]<8>got=80,
}

func ExampleStreamFlatMap() {
	var i int64 = 0
	stream.Of("a", "c", "b").FlatMap(func(s string) iter.Seq[string] {
		i++
		return stream.Repeat(s).Limit(i).Seq()
	}).Sorted(func(s1, s2 string) int {
		return strings.Compare(s1, s2)
	}).ForEach(func(s string) {
		fmt.Printf("%v,", s)
	})
	// Output:
	// a,b,b,b,c,c,
}

func ExampleDistict() {
	s := stream.Repeat(1).Limit(10).Distinct(func(i int) int {
		return i
	}).Collect()
	fmt.Println(s)
	// Output:
	// [1]
}

func ExampleMatch() {
	fmt.Println(stream.Repeat(1).Limit(10).AllMatch(func(i int) bool {
		return i == 1
	}))
	fmt.Println(stream.Range(1, 10).Limit(10).AnyMatch(func(i int) bool {
		return i == 5
	}))
	fmt.Println(stream.Repeat(1).Limit(10).NoneMatch(func(i int) bool {
		return i > 0
	}))
	// Output:
	// true
	// true
	// false
}

func ExampleReduce() {
	fmt.Println(stream.Range(1, 101).Reduce(func(i1, i2 int) int {
		return i1 + i2
	}).Value())
	fmt.Println(stream.Range(1, 2).ReduceFrom(100, func(i1, i2 int) int {
		return i1 + i2
	}))
	fmt.Println(stream.Of[int]().FindFirst().IsAbsent())
	// Output:
	// 5050
	// 101
	// true
}
