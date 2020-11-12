package stream

import (
	"github.com/youthlin/stream/types"
)

// stage 记录一个操作
// Begin 用于操作开始，参数是元素的个数，如果个数不确定，则是 unknownSize
// Accept 接收每个元素
// CanFinish 用于判断是否可以提前结束
// End 是收尾动作
type stage interface {
	Begin(size int64)
	Accept(types.T)
	CanFinish() bool
	End()
}

// region baseStage

type baseStage struct {
	begin     func(int64) // begin(size)
	action    types.Consumer
	canFinish func() bool // canFinish() bool
	end       func()      // end()
}

func (b *baseStage) Begin(size int64) {
	b.begin(size)
}

func (b *baseStage) Accept(t types.T) {
	b.action(t)
}

func (b *baseStage) CanFinish() bool {
	return b.canFinish()
}

func (b *baseStage) End() {
	b.end()
}

type option func(b *baseStage)

func begin(onBegin func(int64)) option {
	return func(c *baseStage) {
		c.begin = onBegin
	}
}

func canFinish(judge func() bool) option {
	return func(c *baseStage) {
		c.canFinish = judge
	}
}

func action(onAction types.Consumer) option {
	return func(c *baseStage) {
		c.action = onAction
	}
}
func end(onEnd func()) option {
	return func(c *baseStage) {
		c.end = onEnd
	}
}

// endregion baseStage

// region chainedStage

// chainedStage 串起下一个操作
// down 是下一个操作
type chainedStage struct {
	*baseStage
}

func defaultChainedStage(down stage) *chainedStage {
	return &chainedStage{
		baseStage: &baseStage{
			begin:     down.Begin,
			action:    down.Accept,
			canFinish: down.CanFinish,
			end:       down.End,
		},
	}
}

func newChainedStage(down stage, opt ...option) *chainedStage {
	s := defaultChainedStage(down)
	for _, o := range opt {
		o(s.baseStage)
	}
	return s
}

//  endregion

// region terminalStage

type terminalStage struct {
	*baseStage
}

func defaultTerminal(action types.Consumer) *terminalStage {
	return &terminalStage{&baseStage{
		begin:     func(int64) {},
		action:    action,
		canFinish: func() bool { return false },
		end:       func() {},
	}}
}

func newTerminalStage(action types.Consumer, opt ...option) *terminalStage {
	s := defaultTerminal(action)
	for _, o := range opt {
		o(s.baseStage)
	}
	return s
}

// endregion terminalStage
