package stream

import (
	"iter"

	"github.com/youthlin/stream/v2/types"
)

// Of build a Seq by input elements.
// 使用输入的任意个元素创建一个序列
func Of[T any](elements ...T) Seq[T] {
	return func(yield func(T) bool) {
		for _, e := range elements {
			if !yield(e) {
				return
			}
		}
	}
}

// OfSeq convert from iter.Seq.
// 类型转换, 从 iter.Seq 转为 Seq.
func OfSeq[T any](s iter.Seq[T]) Seq[T] {
	return Seq[T](s)
}

// CountFrom build a Seq count infinite from the input number.
// 构造一个从指定数字开始计数的序列
func CountFrom[T types.Int](from T) Seq[T] {
	return func(yield func(T) bool) {
		for i := from; ; i++ {
			if !yield(i) {
				return
			}
		}
	}
}

// Range build a Seq count from fromInclude to toExclude.
// 构造一个左闭右开的区间序列
func Range[T types.Int](fromInclude, toExclude T) Seq[T] {
	return RangeStep(fromInclude, toExclude, 1)
}

// RangeStep build a Seq count from fromInclude to toExclude.
// step may negetive.
// 按指定步进大小构造一个左闭右开的序列, 步进大小可以是负数
func RangeStep[T types.Int](fromInclude, toExclude, step T) Seq[T] {
	return func(yield func(T) bool) {
		if step >= 0 {
			for i := fromInclude; i < toExclude; i += step {
				if !yield(i) {
					return
				}
			}
		} else {
			for i := fromInclude; i > toExclude; i += step {
				if !yield(i) {
					return
				}
			}
		}

	}
}

// Repeat create an infinite Seq,
// which all elements is same as the input element.
// 使用输入的单个元素创建一个无限序列
func Repeat[T any](e T) Seq[T] {
	return func(yield func(T) bool) {
		for {
			if !yield(e) {
				return
			}
		}
	}
}

// Generate build a Seq,
// which each element is generate by the Supplier.
// 通过生成器生成一个序列
func Generate[T any](get types.Supplier[T]) Seq[T] {
	return func(yield func(T) bool) {
		for {
			if !yield(get()) {
				return
			}
		}
	}
}
