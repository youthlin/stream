package stream

import (
	"errors"
	"reflect"

	"github.com/youthlin/stream/types"
)

var (
	// ErrNotSlice a error to panic when call Slice but argument is not slice
	ErrNotSlice = errors.New("not slice")
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

// Of create a Stream from some element
// It's recommend to pass pointer type cause the element may be copy at each operate
func Of(elements ...types.T) Stream {
	return newHead(it(elements...))
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
