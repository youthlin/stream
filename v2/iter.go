package stream

import (
	"iter"
	"sort"

	"github.com/youthlin/stream/v2/optional"
	"github.com/youthlin/stream/v2/types"
)

// Filter keep elements which satisfy the Predicate.
// 保留满足断言的元素
func Filter[T any](it iter.Seq[T], test types.Predicate[T]) iter.Seq[T] {
	return func(yield func(T) bool) {
		next, stop := iter.Pull(it)
		defer stop()
		for {
			v, ok := next()
			if !ok {
				return
			}
			if test(v) && !yield(v) {
				return
			}
		}
	}
}

// Map transform the element use Fuction.
// 使用输入函数对每个元素进行转换
func Map[T, R any](it iter.Seq[T], f types.Function[T, R]) iter.Seq[R] {
	return func(yield func(R) bool) {
		next, stop := iter.Pull(it)
		defer stop()
		for {
			v, ok := next()
			if !ok || !yield(f(v)) {
				return
			}
		}
	}
}

// Maps is alias of Map.
// 同 Map, 不过参数是 Seq 而不是 iter.Seq.
func Maps[T, R any](it Seq[T], f types.Function[T, R]) Seq[R] {
	return Seq[R](Map(iter.Seq[T](it), f))
}

// FlatMap transform each element in Seq[T] to a new Seq[R].
// 将原本序列中的每个元素都转换为一个新的序列，
// 并将所有转换后的序列依次连接起来生成一个新的序列
func FlatMap[T, R any](it iter.Seq[T], flatten types.Function[T, iter.Seq[R]]) iter.Seq[R] {
	return func(yield func(R) bool) {
		next, stop := iter.Pull(it)
		defer stop()
		for {
			v, ok := next()
			if !ok {
				return
			}
			for v := range flatten(v) {
				if !yield(v) {
					return
				}
			}
		}
	}
}

// FlatMaps is alias of FlatMap.
// 同 FlatMap, 不过参数是 Seq 而不是 iter.Seq.
func FlatMaps[T, R any](it Seq[T], flatten types.Function[T, iter.Seq[R]]) Seq[R] {
	return Seq[R](FlatMap(iter.Seq[T](it), func(input T) iter.Seq[R] {
		return flatten(input)
	}))
}

// Peek visit every element in the Seq and leave them on the Seq.
// 访问序列中的每个元素而不消费它
func Peek[T any](it iter.Seq[T], accept types.Consumer[T]) iter.Seq[T] {
	return func(yield func(T) bool) {
		next, stop := iter.Pull(it)
		defer stop()
		for {
			v, ok := next()
			if !ok {
				return
			}
			accept(v)
			if !yield(v) {
				return
			}
		}
	}
}

// Distinct remove duplicate elements.
// 对序列中的元素去重
func Distinct[T any, Cmp comparable](it iter.Seq[T], f types.Function[T, Cmp]) iter.Seq[T] {
	return func(yield func(T) bool) {
		next, stop := iter.Pull(it)
		defer stop()
		var set = make(map[Cmp]struct{})
		for {
			v, ok := next()
			if !ok {
				return
			}
			k := f(v)
			_, ok = set[k]
			set[k] = struct{}{}
			if !ok && !yield(v) {
				return
			}
		}
	}
}

// Sorted sort elements in the Seq by Comparator.
// 对序列中的元素排序
func Sorted[T any](it iter.Seq[T], cmp types.Comparator[T]) iter.Seq[T] {
	return func(yield func(T) bool) {
		vals := Collect(it)
		sort.SliceStable(vals, func(i, j int) bool {
			return cmp(vals[i], vals[j]) < 0
		})
		for _, v := range vals {
			if !yield(v) {
				return
			}
		}
	}
}

// Limit limits the number of elements in Seq.
// 限制元素个数
func Limit[T any, Number types.Int](it iter.Seq[T], limit Number) iter.Seq[T] {
	return func(yield func(T) bool) {
		next, stop := iter.Pull(it)
		defer stop()
		for {
			limit--
			if limit < 0 {
				return
			}
			v, ok := next()
			if !ok || !yield(v) {
				return
			}
		}
	}
}

// Skip drop some elements of the Seq.
// 跳过指定个数的元素
func Skip[T any, Number types.Int](it iter.Seq[T], skip Number) iter.Seq[T] {
	return func(yield func(T) bool) {
		next, stop := iter.Pull(it)
		defer stop()
		for {
			v, ok := next()
			skip--
			if skip < 0 {
				if !ok || !yield(v) {
					return
				}
			}
		}
	}
}

// ForEach consume every elements in the Seq.
// 消费序列中的每个元素
func ForEach[T any](it iter.Seq[T], accept types.Consumer[T]) {
	for v := range it {
		accept(v)
	}
}

// Collect return all elements as a slice.
// 将序列中所有元素收集为切片返回
func Collect[T any](it iter.Seq[T]) (result []T) {
	for v := range it {
		result = append(result, v)
	}
	return
}

// AllMatch test if every elements are all match the Predicate.
// 是否每个元素都满足条件
func AllMatch[T any](it iter.Seq[T], test types.Predicate[T]) bool {
	for v := range it {
		if !test(v) {
			return false
		}
	}
	return true
}

// NoneMatch test if none element matches the Predicate.
// 是否没有元素满足条件
func NoneMatch[T any](it iter.Seq[T], test types.Predicate[T]) bool {
	for v := range it {
		if test(v) {
			return false
		}
	}
	return true
}

// AnyMatch test if any element matches the Predicate.
// 是否有任意元素满足条件
func AnyMatch[T any](it iter.Seq[T], test types.Predicate[T]) bool {
	for v := range it {
		if test(v) {
			return true
		}
	}
	return false
}

// Reduce accumulate each element using the binary operation.
// 使用给定的累加函数, 累加序列中的每个元素.
// 序列中可能没有元素因此返回的是 Optional
func Reduce[T any](it iter.Seq[T], acc types.BinaryOperator[T]) types.Optional[T] {
	var result T
	var has bool
	for v := range it {
		if !has {
			result = v
			has = true
		} else {
			result = acc(result, v)
		}
	}
	if has {
		return optional.Of(result)
	}
	return optional.Nil[T]()
}

// ReduceFrom accumulate each element using the binary operation
// starting from the initial value.
// 从初始值开始, 通过 acc 二元运算累加每个元素
func ReduceFrom[T any](it iter.Seq[T], initVal T, acc types.BinaryOperator[T]) (result T) {
	result = initVal
	for v := range it {
		result = acc(result, v)
	}
	return
}

// ReduceWith accumulate each element using the BiFunction
// starting from the initial value.
// 从初始值开始, 通过 acc 函数累加每个元素
func ReduceWith[T, R any](it iter.Seq[T], initVal R, acc types.BiFunction[R, T, R]) (result R) {
	result = initVal
	for v := range it {
		result = acc(result, v)
	}
	return
}

// FindFirst find the first element in the Seq.
// 返回序列中的第一个元素(如有).
func FindFirst[T any](it iter.Seq[T]) types.Optional[T] {
	for v := range it {
		return optional.Of(v)
	}
	return optional.Nil[T]()
}

// Count return the count of elements in the Seq.
// 返回序列中的元素个数
func Count[T any](it iter.Seq[T]) (count int64) {
	for _ = range it {
		count++
	}
	return
}
