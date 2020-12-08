package stream

import (
	"errors"
	"reflect"

	"github.com/youthlin/stream/optional"
	"github.com/youthlin/stream/types"
)

var (
	// ErrNotSlice a error to panic when call Slice but argument is not slice
	ErrNotSlice = errors.New("not slice")
	ErrNotMap   = errors.New("not map")
)

// Slice 把任意的切片类型转为[]T类型. 可用作 Of() 入参.
// Slice convert any slice type to []types.T e.g. []int -> []types.T. may used with Of().
// Note: cannot use ints (type []int) as type []types.T in argument to streams.Of
func Slice(slice types.T) []types.T {
	if reflect.TypeOf(slice).Kind() != reflect.Slice {
		panic(ErrNotSlice)
	}
	var result []types.T
	value := reflect.ValueOf(slice)
	count := value.Len()
	for i := 0; i < count; i++ {
		result = append(result, value.Index(i).Interface())
	}
	return result
}

// Entries 把任意的 map 类型转为 []Pair
// Entries return entries of a map as []types.Pair which `First` field is key, `Second` field is value
func Entries(mapValue types.T) []types.Pair {
	if reflect.TypeOf(mapValue).Kind() != reflect.Map {
		panic(ErrNotMap)
	}
	value := reflect.ValueOf(mapValue)
	var result []types.Pair
	var it = value.MapRange()
	for it.Next() {
		result = append(result, types.Pair{
			First:  it.Key().Interface(),
			Second: it.Value().Interface(),
		})
	}
	return result
}

// Of create a Stream from some element
// It's recommend to pass pointer type cause the element may be copy at each operate
func Of(elements ...types.T) Stream {
	return newHead(it(elements...))
}

func OfInts(element ...int) Stream {
	return newHead(&intsIt{
		base: &base{
			current: 0,
			size:    len(element),
		},
		elements: element,
	})
}
func OfInt64s(element ...int64) Stream {
	return newHead(&int64sIt{
		base: &base{
			current: 0,
			size:    len(element),
		},
		elements: element,
	})
}
func OfFloat32s(element ...float32) Stream {
	return newHead(&float32sIt{
		base: &base{
			current: 0,
			size:    len(element),
		},
		elements: element,
	})
}
func OfFloat64s(element ...float64) Stream {
	return newHead(&float64sIt{
		base: &base{
			current: 0,
			size:    len(element),
		},
		elements: element,
	})
}
func OfStrings(element ...string) Stream {
	return newHead(&stringIt{
		base: &base{
			current: 0,
			size:    len(element),
		},
		elements: element,
	})
}

// OfSlice return a Stream. the input parameter `slice` must be a slice.
// if input is nil, return a empty Stream( same as Of() )
func OfSlice(slice types.T) Stream {
	if optional.IsNil(slice) {
		return Of()
	}
	if reflect.TypeOf(slice).Kind() != reflect.Slice {
		panic(ErrNotSlice)
	}
	value := reflect.ValueOf(slice)
	it := &sliceIt{
		base: &base{
			current: 0,
			size:    value.Len(),
		},
		sliceValue: value,
	}
	return newHead(it)
}

// OfMap return a Stream which element type is types.Pair.
// the input parameter `mapValue` must be a map or it will panic
// if mapValue is nil, return a empty Stream ( same as Of() )
func OfMap(mapValue types.T) Stream {
	if optional.IsNil(mapValue) {
		return Of()
	}
	if reflect.TypeOf(mapValue).Kind() != reflect.Map {
		panic(ErrNotMap)
	}
	value := reflect.ValueOf(mapValue)
	it := &mapIt{
		base: &base{
			current: 0,
			size:    value.Len(),
		},
		mapValue: value.MapRange(),
	}
	return newHead(it)
}

// Iterate create a Stream by a seed and an UnaryOperator
func Iterate(seed types.T, operator types.UnaryOperator) Stream {
	return newHead(withSeed(seed, operator))
}

// Generate generates a infinite Stream which each element is generate by Supplier
func Generate(get types.Supplier) Stream {
	return newHead(withSupplier(get))
}

// Repeat returns a infinite Stream which all element is same
func Repeat(e types.T) Stream {
	return newHead(withSupplier(func() types.T {
		return e
	}))
}

// RepeatN returns a Stream which has `count` element and all the element is the given `e`
func RepeatN(e types.T, count int64) Stream {
	return Repeat(e).Limit(count)
}

// IntRange creates a Stream which element is the given range
func IntRange(fromInclude, toExclude int) Stream {
	return IntRangeStep(fromInclude, toExclude, 1)
}

// IntRangeStep creates a Stream which element is the given range by step
func IntRangeStep(fromInclude, toExclude, step int) Stream {
	return newHead(withRange(epInt(fromInclude), epInt(toExclude), step)).Map(func(t types.T) types.R {
		// streams.epInt is not int
		// 所以转回 int 让调用方不至于迷惑
		return int(t.(epInt))
	})
}

// Int64Range like IntRange
func Int64Range(fromInclude, toExclude int64) Stream {
	return Int64RangeStep(fromInclude, toExclude, 1)
}

// Int64RangeStep like IntRangeStep
func Int64RangeStep(fromInclude, toExclude int64, step int) Stream {
	return newHead(withRange(epInt64(fromInclude), epInt64(toExclude), step)).Map(func(t types.T) types.R {
		return int64(t.(epInt64))
	})
}
