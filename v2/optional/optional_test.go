package optional_test

import (
	"fmt"

	"github.com/youthlin/stream/v2/optional"
	"github.com/youthlin/stream/v2/types"
)

func ExampleNil() {
	o := optional.Nil[int]()
	fmt.Println(o.Get())
	fmt.Printf("IsPresent=%v, IsAbsent=%v\n", o.IsPresent(), o.IsAbsent())
	o.IfPresent(func(i int) {})
	o.IfAbsent(func() {
		fmt.Println("OnIfNil")
	})
	fmt.Printf("Value()=%v\n", o.Value())
	fmt.Printf("Or(42)=%v\n", o.Or(42))
	fmt.Printf("OrZero()=%v\n", o.OrZero())
	fmt.Printf("OrGet()=%v\n", o.OrGet(func() int { return 1 }))
	fmt.Printf("Ptr()=%v\n", o.Ptr())
	fmt.Printf("Filter: %v\n", o.Filter(func(i int) bool { return i == 0 }).IsPresent())
	fmt.Printf("Map: %v\n", o.Map(func(i int) int { return i }).IsPresent())
	fmt.Printf("FlatMap: %v\n", o.FlatMap(func(i int) types.Optional[int] {
		return nil
	}).IsPresent())
	// Output:
	// 0 false
	// IsPresent=false, IsAbsent=true
	// OnIfNil
	// Value()=0
	// Or(42)=42
	// OrZero()=0
	// OrGet()=1
	// Ptr()=<nil>
	// Filter: false
	// Map: false
	// FlatMap: false
}

func ExampleOf() {
	o := optional.Of(1)
	fmt.Println(o.Get())
	fmt.Printf("IsPresent=%v, IsAbsent=%v\n", o.IsPresent(), o.IsAbsent())
	o.IfPresent(func(i int) {
		fmt.Println(i)
	})
	o.IfAbsent(func() {
		fmt.Println("OnIfNil")
	})
	fmt.Printf("Value()=%v\n", o.Value())
	fmt.Printf("Or(42)=%v\n", o.Or(42))
	fmt.Printf("OrZero()=%v\n", o.OrZero())
	fmt.Printf("OrGet()=%v\n", o.OrGet(func() int { return 1 }))
	fmt.Printf("Ptr: %T\n", o.Ptr())
	fmt.Printf("Filter(==0): %v\n", o.Filter(func(i int) bool { return i == 0 }).IsPresent())
	fmt.Printf("Map: %v\n", o.Map(func(i int) int { return i }).Value())
	fmt.Printf("FlatMap: %v\n", o.FlatMap(func(i int) types.Optional[int] {
		return optional.Of(i + 1)
	}).Value())
	// Output:
	// 1 true
	// IsPresent=true, IsAbsent=false
	// 1
	// Value()=1
	// Or(42)=1
	// OrZero()=1
	// OrGet()=1
	// Ptr: *int
	// Filter(==0): false
	// Map: 1
	// FlatMap: 2
}

func ExampleOfPtr() {
	var i int
	var a, b *int
	a = &i
	p := optional.OfPtr(a)
	n := optional.OfPtr(b)
	fmt.Println(p.IsPresent(), n.IsPresent())
	// Output:
	// true false
}
