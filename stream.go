package stream

import (
	"reflect"

	"github.com/youthlin/stream/optional"
	"github.com/youthlin/stream/types"
)

// Stream is a interface which holds all supported operates.
// It has stateless operates(Filter, Map, FlatMap, Peek),
// stateful operates(Distinct, Sorted, Limit, Skip),
// and the left methods are terminal operates.
type Stream interface {
	// stateless operate 无状态操作

	Filter(types.Predicate) Stream         // 过滤
	Map(types.Function) Stream             // 转换
	FlatMap(func(t types.T) Stream) Stream // 打平
	Peek(types.Consumer) Stream            // peek 每个元素

	// stateful operate 有状态操作

	Distinct(types.IntFunction) Stream // 去重
	Sorted(types.Comparator) Stream    // 排序
	Limit(int64) Stream                // 限制个数
	Skip(int64) Stream                 // 跳过个数

	// terminal operate 终止操作

	// 遍历
	ForEach(types.Consumer)
	// return []T 转为切片
	ToSlice() []types.T
	// return []X which X is the type of some
	ToElementSlice(some types.T) types.R
	// return []X which X is same as the `typ` representation
	ToSliceOf(typ reflect.Type) types.R
	// 测试是否所有元素满足条件
	AllMatch(types.Predicate) bool
	// 测试是否没有元素满足条件
	NoneMatch(types.Predicate) bool
	// 测试是否有任意元素满足条件
	AnyMatch(types.Predicate) bool
	// Reduce return optional.Empty if no element. calculate result by (T, T) -> T from first element, panic if reduction is nil
	Reduce(accumulator types.BinaryOperator) optional.Optional
	// type of initValue is same as element.  (T, T) -> T
	ReduceFrom(initValue types.T, accumulator types.BinaryOperator) types.T
	// type of initValue is different from element. (R, T) -> R
	ReduceWith(initValue types.R, accumulator func(types.R, types.T) types.R) types.R
	FindFirst() optional.Optional
	// 返回元素个数
	Count() int64
}
