package stream

import (
	"reflect"
	"sort"

	"github.com/youthlin/stream/optional"
	"github.com/youthlin/stream/types"
)

// stream is a node show as below. which source is a iterator. head stream has no prev node.
// terminal operate create a terminalStage,
// then this terminalStage will use a downStage of prev node and wrap a new stage,
// finally trigger source to iterate it's data to the wrappedStage
//
// stream 是数据流中的一个节点，见下图。头节点没有前驱节点。每个操作会创建一个新的节点，连接到原有节点后。
// 终止操作会创建一个 terminalStage, 这个 terminalStage 会作为最后一个节点的 downStage, 依次往前调用 wrap 方法，生成最终的 wrappedStage
// 最后触发数据源迭代每个元素给 wrappedStage
//
//            head filter   map   for-each
//            +--+    +---+    +--+
//     nil <- |  | <- |   | <- |  | <- terminalStage
//            +--+    +---+    +--+
//
//                +-filter----------------+
//  source -->    |                       |
//                |       +-map-----------+
//                |       |               |
//                |       |    +-for-each-+
//                |       |    |          | terminalStage
//                +-------+----+----------+
//
//               <----- wrapped stage ----->
type stream struct {
	source iterator // 数据源
	prev   *stream  // 前一个流
	wrap   func(stage) stage
}

// region help methods 帮助方法

// newHead 构造头节点
func newHead(source iterator) *stream {
	return &stream{source: source}
}

// newNode 构造中间节点
func newNode(prev *stream, wrap func(down stage) stage) *stream {
	return &stream{
		source: prev.source,
		prev:   prev,
		wrap:   wrap,
	}
}

// terminal 终止操作调用。触发包装各项操作，开始元素遍历
func (s *stream) terminal(ts *terminalStage) {
	stage := s.wrapStage(ts)
	source := s.source
	stage.Begin(source.GetSizeIfKnown())
	for source.HasNext() && !stage.CanFinish() {
		stage.Accept(source.Next())
	}
	stage.End()
}

// wrapStage 将所有操作"包装"为一个操作。从终止操作开始往前(因为 wrap 的参数是 downStage)包装
func (s *stream) wrapStage(terminalStage stage) stage {
	stage := terminalStage
	for i := s; i.prev != nil; i = i.prev {
		stage = i.wrap(stage)
	}
	return stage
}

// endregion 帮助方法

// region 无状态操作

// Filter 过滤操作
// test is a Predicate, return true then keep the element 返回 true 的将保留
func (s *stream) Filter(test types.Predicate) Stream {
	return newNode(s, func(down stage) stage {
		return newChainedStage(down, begin(func(int64) {
			down.Begin(unknownSize) // 过滤后个数不确定
		}), action(func(t types.T) {
			if test(t) {
				down.Accept(t)
			}
		}))
	})
}

// Map 转换操作
// apply is a Function, convert the element to another 转换元素
func (s *stream) Map(apply types.Function) Stream {
	return newNode(s, func(down stage) stage {
		return newChainedStage(down, action(func(t types.T) {
			down.Accept(apply(t))
		}))
	})
}

// FlatMap 打平集合为元素。[[1,2],[3,4]] -> [1,2,3,4]
func (s *stream) FlatMap(flatten func(t types.T) Stream) Stream {
	return newNode(s, func(down stage) stage {
		return newChainedStage(down, begin(func(int64) {
			down.Begin(unknownSize) // 最终个数不确定
		}), action(func(t types.T) {
			ss := flatten(t)        // 元素是集合，转为流
			ss.ForEach(down.Accept) // 消费流中的元素
		}))
	})
}

// Peek visit every element and leave them on stream so that they can be operated by next action  访问流中每个元素而不消费它，可用于 debug
func (s *stream) Peek(consumer types.Consumer) Stream {
	return newNode(s, func(down stage) stage {
		return newChainedStage(down, action(func(t types.T) {
			consumer(t)
			down.Accept(t)
		}))
	})
}

// endregion 无状态操作

// region 有状态操作

// Distinct remove duplicate 去重操作
// distincter is a IntFunction, which return a int hashcode to identity each element 返回元素的唯一标识用于区分每个元素
func (s *stream) Distinct(distincter types.IntFunction) Stream {
	return newNode(s, func(down stage) stage {
		var set map[int]bool
		return newChainedStage(down, begin(func(int64) {
			set = make(map[int]bool)
			down.Begin(unknownSize) // 去重后个数不确定
		}), action(func(t types.T) {
			hash := distincter(t)
			_, has := set[hash]
			if !has { // 唯一的元素才往下游发送
				down.Accept(t)
				set[hash] = true
			}
		}), end(func() {
			set = nil
			down.End()
		}))
	})
}

// Sorted sort by Comparator 排序
func (s *stream) Sorted(cmp types.Comparator) Stream {
	return newNode(s, func(down stage) stage {
		var list []types.T
		return newChainedStage(down, begin(func(size int64) {
			if size > 0 {
				list = make([]types.T, 0, size)
			} else {
				list = make([]types.T, 0)
			}
			down.Begin(size)
		}), action(func(t types.T) {
			list = append(list, t)
		}), end(func() {
			a := &Sortable{
				List: list,
				Cmp:  cmp,
			}
			sort.Sort(a)
			down.Begin(int64(len(a.List)))
			i := it(a.List...)
			for i.HasNext() && !down.CanFinish() {
				down.Accept(i.Next())
			}
			list = nil
			a = nil
			down.End()
		}))
	})
}

// Limit 限制元素个数
func (s *stream) Limit(maxSize int64) Stream {
	return newNode(s, func(down stage) stage {
		count := int64(0)
		return newChainedStage(down, begin(func(size int64) {
			if size > 0 {
				if size > maxSize {
					size = maxSize
				}
			}
			down.Begin(size)
		}), action(func(t types.T) {
			if count < maxSize {
				down.Accept(t)
			}
			count++
		}), canFinish(func() bool {
			return count == maxSize // 已经到了限制数量，就可以提前结束了
		}))
	})
}

// SKip 跳过指定个数的元素
func (s *stream) Skip(n int64) Stream {
	return newNode(s, func(down stage) stage {
		count := int64(0)
		return newChainedStage(down, begin(func(size int64) {
			if size > 0 {
				size -= n
				if size < 0 {
					size = 0
				}
			}
			down.Begin(size)
		}), action(func(t types.T) {
			if count >= n {
				down.Accept(t)
			}
			count++
		}))
	})
}

// endregion 有状态操作

// region 终止操作

// ForEach 消费流中每个元素
func (s *stream) ForEach(consumer types.Consumer) {
	s.terminal(newTerminalStage(consumer))
}

// ToSlice 转为切片
func (s *stream) ToSlice() []types.T {
	return s.ReduceBy(func(count int64) types.R {
		if count >= 0 {
			return make([]types.T, 0, count)
		}
		return make([]types.T, 0)
	}, func(acc types.R, t types.T) types.R {
		slice := acc.([]types.T)
		slice = append(slice, t)
		return slice
	}).([]types.T)
}

// ToElementSlice needs a argument cause the stream may be empty
func (s *stream) ToElementSlice(some types.T) types.R {
	return s.ToSliceOf(reflect.TypeOf(some))
}

// ToRealSlice
func (s *stream) ToSliceOf(typ reflect.Type) types.R {
	sliceType := reflect.SliceOf(typ)
	return s.ReduceBy(func(size int64) types.R {
		if size >= 0 {
			return reflect.MakeSlice(sliceType, 0, int(size))
		}
		return reflect.MakeSlice(sliceType, 0, 16)
	}, func(acc types.R, t types.T) types.R {
		sliceValue := acc.(reflect.Value)
		sliceValue = reflect.Append(sliceValue, reflect.ValueOf(t))
		return sliceValue
	}).(reflect.Value).
		Interface()
}

// AllMatch 测试是否所有元素符合断言
func (s *stream) AllMatch(test types.Predicate) bool {
	var result = true
	s.terminal(newTerminalStage(func(t types.T) {
		if !test(t) {
			result = false // 有任意一个不符合
		}
	}, canFinish(func() bool {
		return !result // 有一个不符合就可以结束了
	})))
	return result
}

// NoneMatch 测试是否没有元素符合断言
func (s *stream) NoneMatch(test types.Predicate) bool {
	var result = true
	s.terminal(newTerminalStage(func(t types.T) {
		if test(t) {
			result = false // 有任意一个符合
		}
	}, canFinish(func() bool {
		return !result // 有一个符合就可以结束了
	})))
	return result
}

// AnyMatch 测试是否有任意一个元素符合断言
func (s *stream) AnyMatch(test types.Predicate) bool {
	var result = false
	s.terminal(newTerminalStage(func(t types.T) {
		if test(t) {
			result = true // 有任意一个符合
		}
	}, canFinish(func() bool {
		return result // 有任意一个符合就可以结束
	})))
	return result
}

func (s *stream) Reduce(accumulator types.BinaryOperator) optional.Optional {
	var result types.T = nil
	var hasElement = false
	s.terminal(newTerminalStage(func(t types.T) {
		if !hasElement {
			result = t
			hasElement = true
		} else {
			result = accumulator(result, t)
		}
	}))
	if hasElement {
		return optional.Of(result)
	}
	return optional.Empty()
}

// ReduceFrom 从给定的初始值 initValue(类型和元素类型相同) 开始迭代 使用 accumulator(2个入参类型和返回类型相同) 累计结果
func (s *stream) ReduceFrom(initValue types.T, accumulator types.BinaryOperator) types.T {
	var result = initValue
	s.terminal(newTerminalStage(func(t types.T) {
		result = accumulator(result, t)
	}))
	return result
}

// ReduceWith 使用给定的初始值 initValue(类型和元素类型不同) 开始迭代 使用 accumulator( R + T -> R) 累计结果
func (s *stream) ReduceWith(initValue types.R, accumulator func(types.R, types.T) types.R) types.R {
	var result = initValue
	s.terminal(newTerminalStage(func(t types.T) {
		result = accumulator(result, t)
	}))
	return result
}

// ReduceBy 使用给定的初始化方法(参数是元素个数，或-1)生成 initValue, 然后使用 accumulator 累计结果
// ReduceBy use `buildInitValue` to build the initValue, which parameter is a int64 means element size, or -1 if unknown size.
// Then use `accumulator` to add each element to previous result
func (s *stream) ReduceBy(buildInitValue func(int64) types.R, accumulator func(types.R, types.T) types.R) types.R {
	var result types.R
	s.terminal(newTerminalStage(func(e types.T) {
		result = accumulator(result, e)
	}, begin(func(count int64) {
		result = buildInitValue(count)
	})))
	return result
}

func (s *stream) FindFirst() optional.Optional {
	var result types.T
	var find = false
	s.terminal(newTerminalStage(func(t types.T) {
		if !find {
			result = t
			find = true
		}
	}, canFinish(func() bool {
		return find
	})))
	return optional.OfNullable(result)
}

// Count 计算元素个数
func (s *stream) Count() int64 {
	return s.ReduceWith(int64(0), func(count types.R, t types.T) types.R {
		return count.(int64) + 1
	}).(int64)
}

// endregion 终止操作
