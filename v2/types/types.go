package types

import "iter"

type Optional[T any] interface {
	Get() (T, bool)
	Must() T // panic when nil
	IsPresent() bool
	IfPresent(Consumer[T])
	IsAbsent() bool
	IfAbsent(func())
	Value() T // nillable
	Or(v T) T
	OrZero() T
	OrGet(Supplier[T]) T
	Ptr() *T
	Filter(Predicate[T]) Optional[T]
	Map(Function[T, T]) Optional[T]
	FlatMap(Function[T, Optional[T]]) Optional[T]
}

type Stream[T any] interface {
	Seq() iter.Seq[T]

	Filter(Predicate[T]) Stream[T]
	Map(Function[T, T]) Stream[T]
	FlatMap(Function[T, iter.Seq[T]]) Stream[T]
	Peek(Consumer[T]) Stream[T]

	Distinct(Function[T, int]) Stream[T]
	Sorted(Comparator[T]) Stream[T]
	Limit(int64) Stream[T]
	Skip(int64) Stream[T]

	ForEach(Consumer[T])
	Collect() []T
	AllMatch(Predicate[T]) bool
	NoneMatch(Predicate[T]) bool
	AnyMatch(Predicate[T]) bool
	Reduce(acc BinaryOperator[T]) Optional[T]
	ReduceFrom(initVal T, acc BinaryOperator[T]) T
	FindFirst() Optional[T]
	Count() int64
}

// Supplier 产生一个元素
type Supplier[T any] func() T

// Function 将一个类型转为另一个类型
type Function[T, R any] func(T) R

// Predicate 断言是否满足指定条件
type Predicate[T any] Function[T, bool]

// UnaryOperator 对输入进行一元运算返回相同类型的结果
type UnaryOperator[T any] Function[T, T]

// BiFunction 将两个类型转为第三个类型
type BiFunction[T, R, U any] func(T, R) U

// BinaryOperator 输入两个相同类型的参数，对其做二元运算，返回相同类型的结果
type BinaryOperator[T any] BiFunction[T, T, T]

// Comparator 比较两个元素.
// 第一个元素大于第二个元素时，返回正数;
// 第一个元素小于第二个元素时，返回负数;
// 否则返回 0.
type Comparator[T any] BiFunction[T, T, int]

// Consumer 消费一个元素
type Consumer[T any] func(T)

// 整数类型
type Int interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64
}

// Pair 表示一对关联元素
type Pair[T, R any] struct {
	First  T
	Second R
}

// ReverseOrder .
// 反转顺序
func ReverseOrder[T any](cmp Comparator[T]) Comparator[T] {
	return func(t1, t2 T) int {
		return cmp(t2, t1)
	}
}
