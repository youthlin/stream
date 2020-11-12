package optional

import "github.com/youthlin/stream/types"

type present struct {
	value types.T
}

func (p *present) Get() types.T {
	return p.value
}

func (p *present) IsPresent() bool {
	return true
}

func (p *present) IfPresent(action types.Consumer) {
	action(p.value)
}

func (p *present) Filter(test types.Predicate) Optional {
	if test(p.value) {
		return p
	}
	return empty
}

func (p *present) Map(mapper types.Function) Optional {
	return OfNullable(mapper(p.value))
}

func (p *present) FlatMap(flatMapper func(t types.T) Optional) Optional {
	return flatMapper(p.value)
}

func (p *present) OrElse(types.T) types.T {
	return p.value
}

func (p *present) OrElseGet(types.Supplier) types.T {
	return p.value
}

func (p *present) OrPanic(interface{}) types.T {
	return p.value
}

func (p *present) OrPanicGet(types.Supplier) types.T {
	return p.value
}
