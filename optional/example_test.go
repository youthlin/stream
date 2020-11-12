package optional_test

import (
	"errors"
	"fmt"

	"github.com/youthlin/stream/optional"
	"github.com/youthlin/stream/types"
)

func ExampleEmpty() {
	fmt.Println(optional.Empty().IsPresent()) // false
	optional.Empty().IfPresent(func(t types.T) {
		fmt.Println("This will not print")
	})
	fmt.Printf("%#v\n", optional.Empty().Filter(func(t types.T) bool {
		fmt.Println("Filter: Not Reach")
		return true
	}).Map(func(t types.T) types.R {
		fmt.Println("Map: Not Reach")
		return t
	}).FlatMap(func(t types.T) optional.Optional {
		fmt.Println("FlatMap: Not Reach")
		return optional.Empty()
	}).OrElse(1)) // 1
	fmt.Println(optional.Empty().OrElseGet(func() types.T {
		return "else"
	}))
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
		defer func() {
			if err := recover(); err != nil {
				fmt.Println(err)
				defer func() {
					if err := recover(); err != nil {
						fmt.Printf("OrPanicGet:%+v\n", err)
					}
				}()
				optional.Empty().OrPanicGet(func() types.T {
					return errors.New("err_msg")
				})
			}
		}()
		optional.Empty().OrPanic("error_msg")
	}()
	optional.Empty().Get()
	// Output:
	// false
	// 1
	// else
	// absent value
	// error_msg
	// OrPanicGet:err_msg
}
func ExampleOf() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
		fmt.Println(optional.Of(1).IsPresent()) // true
		fmt.Println(optional.Of(1).Filter(func(t types.T) bool {
			return t.(int) > 1
		}).IsPresent()) // false
		fmt.Println(optional.Of(1).Filter(func(t types.T) bool {
			return t.(int) == 1
		}).IsPresent()) // true
		fmt.Println(optional.Of(1).Map(func(t types.T) types.R {
			return fmt.Sprintf("<%d>", t)
		}).Get())
		optional.Of(1).FlatMap(optional.Of).IfPresent(func(t types.T) {
			fmt.Printf("FlatMap: %d\n", t)
		})
		fmt.Printf("OrElse: %d\n", optional.Of(1).OrElse(0))
		fmt.Printf("OrElseGet: %d\n", optional.Of(1).OrElseGet(func() types.T {
			return 0
		}))
		fmt.Printf("OrPanic: %d\n", optional.Of(1).OrPanic("boom"))

		defer func() {
			if err := recover(); err != nil {
				fmt.Println(err)
			}
		}()
		var s *string = nil
		optional.Of(s).IfPresent(func(t types.T) {
			fmt.Printf("%v\n", t)
		})
	}()
	of := optional.Of(nil)
	fmt.Println(of)
	// Output:
	// nil value
	// true
	// false
	// true
	// <1>
	// FlatMap: 1
	// OrElse: 1
	// OrElseGet: 1
	// OrPanic: 1
	// nil value
}
func ExampleOfNullable() {
	fmt.Printf("nil IsPresent: %t\n", optional.OfNullable(nil).IsPresent())
	fmt.Printf("nil(*string) IsPresent: %t\n", optional.OfNullable((*string)(nil)).IsPresent())
	fmt.Printf("IsPresent: %t\n", optional.OfNullable(1).IsPresent())
	// Output:
	// nil IsPresent: false
	// nil(*string) IsPresent: false
	// IsPresent: true
}
