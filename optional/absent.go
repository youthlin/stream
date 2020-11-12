package optional

import (
	"github.com/youthlin/stream/types"
)

type absent int

const empty = absent(0)

func (a absent) Get() types.T {
	panic(ErrAbsent)
}

func (a absent) IsPresent() bool {
	return false
}

func (a absent) IfPresent(types.Consumer) {}

func (a absent) Filter(types.Predicate) Optional {
	return a
}

func (a absent) Map(types.Function) Optional {
	return a
}

func (a absent) FlatMap(func(t types.T) Optional) Optional {
	return a
}

func (a absent) OrElse(t types.T) types.T {
	return t
}

func (a absent) OrElseGet(get types.Supplier) types.T {
	return get()
}

func (a absent) OrPanic(panicArg interface{}) types.T {
	panic(panicArg)
}

func (a absent) OrPanicGet(getPanicArg types.Supplier) types.T {
	panic(getPanicArg())
}
