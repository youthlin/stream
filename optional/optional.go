package optional

import (
	"errors"
	"reflect"

	"github.com/youthlin/stream/types"
)

// Optional is a container which may or may not contain a non-nil value.
// If value is not nil, IsPresent will return true, and Get will return the value.
type Optional interface {
	Get() types.T                                  // Get the value or panic when no value
	IsPresent() bool                               // IsPresent return true is value exist
	IfPresent(types.Consumer)                      // IfPresent will invoke the Consumer if value exist
	Filter(types.Predicate) Optional               // Filter: if has value, and the value matches the given Predicate return a Present Optional, or else return Empty
	Map(types.Function) Optional                   // Map: if has value, apply the given Function to it, or else return Empty
	FlatMap(func(t types.T) Optional) Optional     // Flatmap: if has value, apply the given flatten-Function and return the result, or else return Empty
	OrElse(types.T) types.T                        // OrElse: if absent, return the given value. if value present, return it
	OrElseGet(types.Supplier) types.T              // OrElseGet: if absent, call Supplier and return it's result. if value present return it
	OrPanic(panicArg interface{}) types.T          // OrPanic:if absent, panic with `panicArg`, if value present, return it
	OrPanicGet(getPanicArg types.Supplier) types.T // OrPanicGet: if absent, panic with the given supplier's result. if value present, return it
}

var (
	// ErrAbsent used to panic when call Optional.Get on empty Optional
	ErrAbsent = errors.New("absent value")
	// ErrNil used to panic when input a nil value to the Of method
	ErrNil = errors.New("nil value")
)

// Empty return a absent optional
func Empty() Optional {
	return empty
}

// Of returns a present optional, the input value must not be nil, or it will panic
func Of(value types.T) Optional {
	if IsNil(value) {
		panic(ErrNil)
	}
	return &present{value}
}

// OfNullable return a absent or present optional depends on the input value is nil or false
func OfNullable(value types.T) Optional {
	if IsNil(value) {
		return Empty()
	}
	return Of(value)
}

// IsNil used to determine a value whether is nil or not
func IsNil(value types.T) bool {
	if value == nil {
		return true
	}
	reflectValue := reflect.ValueOf(value)
	switch reflectValue.Kind() {
	case reflect.Chan, reflect.Func, reflect.Map, reflect.Ptr, reflect.UnsafePointer,
		reflect.Interface, reflect.Slice:
		if reflectValue.IsNil() {
			return true
		}
	}
	return false
}
