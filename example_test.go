package stream_test

import (
	"fmt"
	"reflect"
	"sort"
	"testing"

	"github.com/youthlin/stream"
	"github.com/youthlin/stream/types"
)

func ExampleSlice() {
	var ints = []int{1, 2, 3}
	fmt.Printf("%#v\n", stream.Slice(ints))
	var str = []string{"abc", "###"}
	fmt.Printf("%#v\n", stream.Slice(str))
	// Output:
	// []types.T{1, 2, 3}
	// []types.T{"abc", "###"}
}

func ExampleEntries() {
	var m1 = map[int]string{
		1: "a",
		2: "b",
		3: "c",
	}
	entries := stream.Entries(m1)
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].First.(int) < entries[j].First.(int)
	})
	fmt.Printf("%v\n", entries)
	stream.Of(stream.Slice(entries)...).ReduceWith(map[string]int{}, func(acc types.R, e types.T) types.R {
		pair := e.(types.Pair)
		(acc.(map[string]int))[pair.Second.(string)] = pair.First.(int)
		return acc
	})
	// Output:
	// [{1 a} {2 b} {3 c}]
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

func ExampleOfSlice() {
	var intArr = []int{1, 2, 3, 4}
	stream.OfSlice(intArr).ForEach(func(e types.T) {
		fmt.Printf("%d,", e)
	})
	var nilArr []int
	stream.OfSlice(nilArr).ForEach(func(e types.T) {
		fmt.Printf("should not print")
	})
	var strArr = []string{"a", "b"}
	stream.OfSlice(strArr).
		Map(func(e types.T) types.R {
			return fmt.Sprintf("<%s>", e)
		}).
		ForEach(func(e types.T) {
			fmt.Printf("%s,", e)
		})
	// Output:
	// 1,2,3,4,<a>,<b>,
}

func ExampleOfMap() {
	var m1 = map[int]string{
		3: "c",
		2: "b",
		1: "a",
	}
	s := stream.OfMap(m1).
		Map(func(e types.T) types.R {
			p := e.(types.Pair)
			p.First, p.Second = p.Second, p.First
			return p
		}).
		Sorted(func(left types.T, right types.T) int {
			p1 := left.(types.Pair)
			p2 := right.(types.Pair)
			return p1.Second.(int) - p2.Second.(int)
		}).
		ToSlice()
	fmt.Println(s)
	stream.OfMap(nil).ForEach(func(e types.T) {
		fmt.Println("not print")
	})
	// Output:
	// [{a 1} {b 2} {c 3}]
}

func ExampleIterate() {
	// 0   1   1   2 3 5 8
	// |   |  next
	// |   \curr(start)
	// \prev
	var prev = 0
	stream.Iterate(1, func(t types.T) types.T {
		curr := t.(int)
		next := prev + curr
		prev = curr
		return next
	}).
		Limit(6).
		ForEach(func(t types.T) {
			fmt.Printf("%d,", t)
		})
	// Output:
	// 1,1,2,3,5,8,
}
func ExampleGenerate() {
	// fibonacci 斐波那契数列 0, 1, 1, 2, 3, 5, 8, 13, 21, 34, 55, 89, 144, 233，377，610，987，1597，2584，4181，6765，10946，17711，28657，46368
	var fibonacci = func() types.Supplier {
		// -1 1 【0】 1 1 2 3
		a := -1
		b := 1
		return func() types.T {
			n := a + b
			a = b
			b = n
			return n
		}
	}
	stream.Generate(fibonacci()).Limit(20).ForEach(func(t types.T) {
		fmt.Printf("%d,", t)
	})
	// Output:
	// 0,1,1,2,3,5,8,13,21,34,55,89,144,233,377,610,987,1597,2584,4181,
}
func ExampleRepeat() {
	count := stream.Repeat("a").Peek(func(t types.T) {
		fmt.Printf("%s", t)
	}).Limit(10).Count()
	fmt.Printf("\n%d\n", count)
	// Output:
	// aaaaaaaaaa
	// 10
}
func ExampleRepeatN() {
	stream.RepeatN(1.0, 3).ForEach(func(t types.T) {
		fmt.Printf("<%T,%v>", t, t)
	})
	// Output:
	// <float64,1><float64,1><float64,1>
}
func ExampleIntRange() {
	stream.IntRange(0, 5).
		ForEach(func(t types.T) {
			fmt.Println(t)
		})
	// Output:
	// 0
	// 1
	// 2
	// 3
	// 4
}
func ExampleIntRangeStep() {
	stream.IntRangeStep(0, 10, 2).ForEach(func(t types.T) {
		fmt.Println(t)
	})
	// Output:
	// 0
	// 2
	// 4
	// 6
	// 8
}
func ExampleIntRangeStep_negStep() {
	stream.IntRangeStep(5, 0, -1).ForEach(func(t types.T) {
		fmt.Printf("%d,", t)
	})
	// Output:
	// 5,4,3,2,1,
}
func ExampleIntRangeStep_zeroStep() {
	stream.IntRangeStep(0, 5, 0).
		Limit(10).ForEach(func(t types.T) {
		fmt.Printf("%d,", t)
	})
	// Output:
	// 0,0,0,0,0,0,0,0,0,0,
}
func ExampleInt64Range() {
	stream.Int64Range(int64(0), int64(10)).ForEach(func(t types.T) {
		fmt.Printf("%d", t)
	})
	// Output:
	// 0123456789
}
func ExampleInt64RangeStep() {
	stream.Int64RangeStep(int64(0), int64(10), 3).ForEach(func(t types.T) {
		fmt.Printf("%d", t)
	})
	// Output:
	// 0369
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
func ExampleStream_Peek() {
	stream.Of(1, 2, 3, 4, 5).Peek(func(t types.T) {
		fmt.Printf("%d,", t)
	}).Count()
	// Output:
	// 1,2,3,4,5,
}
func ExampleStream_Peek_peekIsNotTerminalAction() {
	stream.Of(1, 2, 3, 4, 5).Peek(func(t types.T) {
		fmt.Printf("%d,", t)
	})
	stream.Of(1, 2, 3, 4, 5).Peek(func(t types.T) {
		fmt.Printf("%d,", t)
	}).Count()
	// Output:
	// 1,2,3,4,5,
}

func ExampleStream_Distinct() {
	fmt.Println(stream.RepeatN(1, 10).Distinct(func(t types.T) int {
		return t.(int)
	}).Count())
	// Output:
	// 1
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
func ExampleStream_Sorted_int64() {
	stream.Int64RangeStep(100, 0, -10).
		Peek(func(t types.T) {
			fmt.Printf("%d,", t)
		}).
		Sorted(types.Int64Comparator).
		ForEach(func(t types.T) {
			fmt.Printf("%d,", t)
		})
	// Output:
	// 100,90,80,70,60,50,40,30,20,10,10,20,30,40,50,60,70,80,90,100,
}
func ExampleStream_Limit() {
	fmt.Println(stream.Repeat(nil).Limit(100).Count())
	// Output:
	// 100
}
func ExampleStream_Skip() {
	stream.IntRange(0, 10).Skip(5).ForEach(func(t types.T) {
		fmt.Printf("%d,", t)
	})
	// Output:
	// 5,6,7,8,9,
}

func ExampleStream_ForEach() {
	stream.Of("hello", "world").ForEach(func(t types.T) {
		fmt.Println(t)
	})
	// Output:
	// hello
	// world
}
func ExampleStream_ToSlice() {
	slice := stream.Of(1, 2, 3).ToSlice()
	fmt.Printf("%#v\n", slice)
	// Output:
	// []types.T{1, 2, 3}
}
func ExampleStream_ToElementSlice() {
	slice := stream.Of(1, 2, 3).ToElementSlice(0)
	fmt.Printf("%#v\n", slice)
	// Output:
	// []int{1, 2, 3}
}
func ExampleStream_ToSliceOf() {
	slice := stream.Of(1, 2, 3).ToSliceOf(reflect.TypeOf(0))
	fmt.Printf("%#v\n", slice)
	// Output:
	// []int{1, 2, 3}
}
func ExampleStream_AllMatch() {
	allMatch := stream.IntRange(0, 10).AllMatch(func(t types.T) bool {
		i, ok := t.(int)
		return ok && i >= 0 && i < 10
	})
	fmt.Println(allMatch)
	allMatch = stream.IntRange(0, 10).AllMatch(func(t types.T) bool {
		i, ok := t.(int)
		return ok && i > 0 && i < 10
	})
	fmt.Println(allMatch)
	noElementSoResultShouldBeFalse := stream.Of().AnyMatch(func(t types.T) bool {
		return true
	})
	fmt.Println(noElementSoResultShouldBeFalse)
	// Output:
	// true
	// false
	// false
}
func ExampleStream_NoneMatch() {
	isStr := func(t types.T) bool {
		_, ok := t.(string)
		return ok
	}
	noneIsStr := stream.IntRange(0, 10).NoneMatch(isStr)
	fmt.Println(noneIsStr) // true
	noneIsStr = stream.Of(1, 2, "a", "b").NoneMatch(isStr)
	fmt.Println(noneIsStr) // false
	noElementSoResultShouldBeTrue := stream.Of().NoneMatch(func(t types.T) bool {
		return true
	})
	fmt.Println(noElementSoResultShouldBeTrue)
	// Output:
	// true
	// false
	// true
}
func ExampleStream_AnyMatch() {
	hasStr := stream.Of(0, 1, "a", "b").AnyMatch(func(t types.T) bool {
		_, ok := t.(string)
		return ok
	})
	hasInt := stream.Of(0, 1, "a", "b").AnyMatch(func(t types.T) bool {
		_, ok := t.(int)
		return ok
	})
	fmt.Println(hasStr, hasInt)
	noElementSoResultShouldBeFalse := stream.Of().AnyMatch(func(t types.T) bool {
		return true
	})
	fmt.Println(noElementSoResultShouldBeFalse)
	// Output:
	// true true
	// false
}
func ExampleStream_Reduce() {
	fmt.Println(stream.Of().Reduce(func(acc types.T, t types.T) types.T {
		return acc
	}).IsPresent())
	// Output:
	// false
}
func ExampleStream_ReduceFrom() {
	sum := stream.IntRange(1, 101).ReduceFrom(0, func(acc types.T, t types.T) types.T {
		return acc.(int) + t.(int)
	})
	fmt.Println(sum)
	mul := stream.IntRange(1, 5).ReduceFrom(1, func(acc types.T, t types.T) types.T {
		return acc.(int) * t.(int)
	})
	fmt.Println(mul)
	// Output:
	// 5050
	// 24
}
func ExampleStream_ReduceWith() {
	slice := stream.IntRange(0, 10).ReduceWith(make([]int, 0, 10), func(acc types.R, t types.T) types.R {
		return append(acc.([]int), t.(int))
	})
	fmt.Printf("%#v", slice)
	// Output:
	// []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
}
func ExampleStream_Count() {
	fmt.Println(stream.Of().Count())
	fmt.Println(stream.Of(1).Count())
	fmt.Println(stream.Of(1, 2).Count())
	// Output:
	// 0
	// 1
	// 2
}

type person struct {
	name string
	age  int
}

func TestStruct(t *testing.T) {
	stream.Of(&person{
		name: "Bob",
		age:  18,
	}, &person{
		name: "Alice",
		age:  20,
	}).Map(func(e types.T) types.R {
		p := e.(*person)
		p.age++
		t.Logf("%#v->%#v,", e, *p)
		return e
	}).Sorted(types.ReverseOrder(func(left, right types.T) int {
		l := left.(*person)
		r := right.(*person)
		return l.age - r.age
	})).ForEach(func(e types.T) {
		t.Log(e)
	})
	// Output:
	// &streams_test.person{name:"Bob", age:19}->streams_test.person{name:"Bob", age:19},
	// &streams_test.person{name:"Alice", age:21}->streams_test.person{name:"Alice", age:21},
	// &{Alice 21}
	// &{Bob 19}
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
