package stream

import (
	"iter"

	"github.com/youthlin/stream/v2/types"
)

var _ types.Stream[int] = Of[int]()

type Seq[T any] iter.Seq[T]

func (it Seq[T]) Seq() iter.Seq[T] {
	return iter.Seq[T](it)
}

func (it Seq[T]) Filter(test types.Predicate[T]) types.Stream[T] {
	return Seq[T](Filter(iter.Seq[T](it), test))
}

func (it Seq[T]) Map(f types.Function[T, T]) types.Stream[T] {
	return Maps(it, f)
}

func (it Seq[T]) FlatMap(f types.Function[T, iter.Seq[T]]) types.Stream[T] {
	return FlatMaps(it, f)
}

func (it Seq[T]) Peek(accept types.Consumer[T]) types.Stream[T] {
	return Seq[T](Peek(iter.Seq[T](it), accept))
}

func (it Seq[T]) Distinct(f types.Function[T, int]) types.Stream[T] {
	return Seq[T](Distinct(iter.Seq[T](it), f))
}

func (it Seq[T]) Sorted(cmp types.Comparator[T]) types.Stream[T] {
	return Seq[T](Sorted(iter.Seq[T](it), cmp))
}

func (it Seq[T]) Limit(limit int64) types.Stream[T] {
	return Seq[T](Limit(iter.Seq[T](it), limit))
}

func (it Seq[T]) Skip(skip int64) types.Stream[T] {
	return Seq[T](Skip(iter.Seq[T](it), skip))
}

func (it Seq[T]) ForEach(accept types.Consumer[T]) {
	ForEach(iter.Seq[T](it), accept)
}

func (it Seq[T]) Collect() []T {
	return Collect(iter.Seq[T](it))
}

func (it Seq[T]) AllMatch(test types.Predicate[T]) bool {
	return AllMatch(iter.Seq[T](it), test)
}

func (it Seq[T]) NoneMatch(test types.Predicate[T]) bool {
	return NoneMatch(iter.Seq[T](it), test)
}

func (it Seq[T]) AnyMatch(test types.Predicate[T]) bool {
	return AnyMatch(iter.Seq[T](it), test)
}

func (it Seq[T]) Reduce(acc types.BinaryOperator[T]) types.Optional[T] {
	return Reduce(iter.Seq[T](it), acc)
}

func (it Seq[T]) ReduceFrom(initVal T, acc types.BinaryOperator[T]) T {
	return ReduceFrom(iter.Seq[T](it), initVal, acc)
}

func (it Seq[T]) FindFirst() types.Optional[T] {
	return FindFirst(iter.Seq[T](it))
}

func (it Seq[T]) Count() int64 {
	return Count(iter.Seq[T](it))
}
