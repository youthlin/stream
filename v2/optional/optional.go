package optional

import (
	"fmt"
	"reflect"

	"github.com/youthlin/stream/v2/types"
)

var ErrAbsent = fmt.Errorf("value is absent")

func Of[T any](v T) types.Optional[T] {
	return &optional[T]{val: v, ok: true}
}

func Nil[T any]() types.Optional[T] {
	return &optional[T]{ok: false}
}

func OfPtr[T any](v *T) types.Optional[T] {
	if v == nil {
		return Nil[T]()
	}
	return Of(*v)
}

type optional[T any] struct {
	val T
	ok  bool
}

func (o *optional[T]) Get() (T, bool) { return o.val, o.ok }

func (o *optional[T]) Must() T {
	if o.ok {
		return o.val
	}
	panic(ErrAbsent)
}

func (o *optional[T]) IsPresent() bool { return o.ok }

func (o *optional[T]) IsAbsent() bool { return !o.ok }

func (o *optional[T]) IfPresent(accept types.Consumer[T]) {
	if o.ok {
		accept(o.val)
	}
}
func (o *optional[T]) IfAbsent(f func()) {
	if !o.ok {
		f()
	}
}

func (o *optional[T]) Value() T { return o.val }

func (o *optional[T]) Or(t T) T {
	if o.ok {
		return o.val
	}
	return t
}

func (o *optional[T]) OrZero() (t T) {
	if o.ok {
		return o.val
	}
	return
}

func (o *optional[T]) OrGet(get types.Supplier[T]) T {
	if o.ok {
		return o.val
	}
	return get()
}

func (o *optional[T]) Ptr() *T {
	if o.ok {
		return &o.val
	}
	return nil
}

func (o *optional[T]) Filter(test types.Predicate[T]) types.Optional[T] {
	if o.ok {
		if test(o.val) {
			return o
		}
	}
	return Nil[T]()
}

func (o *optional[T]) Map(f types.Function[T, T]) types.Optional[T] {
	if o.ok {
		v := f(o.val)
		if IsNil(v) {
			return Nil[T]()
		}
		return &optional[T]{val: v, ok: true}
	}
	return Nil[T]()
}

func (o *optional[T]) FlatMap(f types.Function[T, types.Optional[T]]) types.Optional[T] {
	if o.ok {
		return f(o.val)
	}
	return Nil[T]()
}

func IsNil(a any) bool {
	return a == nil || (Nillable(a) && reflect.ValueOf(a).IsNil())
}

func Nillable(a any) bool {
	rt := reflect.TypeOf(a)
	switch rt.Kind() {
	case reflect.Chan, reflect.Func, reflect.Map,
		reflect.Pointer, reflect.UnsafePointer,
		reflect.Slice, reflect.Interface:
		return true
	}
	return false
}
