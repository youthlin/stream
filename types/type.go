package types

type (
	// T is a empty interface, that is `any` type.
	// since Go is not support generics now(but will coming soon),
	// so we use T to represent any type
	T interface{}
	// R is another `any` type used to distinguish T
	R interface{}
	// U is another `any` type used to distinguish T and R
	U interface{}
	// Function represents a conversion ability, which accepts one argument and produces a result
	Function func(e T) R
	// IntFunction is a Function, which result type is int
	IntFunction func(e T) int
	// Predicate is a Function, which produces a bool value. usually used to test a value whether satisfied condition
	Predicate func(e T) bool
	// UnaryOperator is a Function, which argument and result are the same type
	UnaryOperator func(e T) T
	// Consumer accepts one argument and not produces any result
	Consumer func(e T)
	// Supplier returns a result. each time invoked it can returns a new or distinct result
	Supplier func() T
	// BiFunction like Function, but is accepts two arguments and produces a result
	BiFunction func(t T, u U) R
	// BinaryOperator is a BiFunction which input and result are the same type
	BinaryOperator func(e1 T, e2 T) T
	// Comparator is a BiFunction, which two input arguments are the type, and returns a int.
	// if left is greater then right, it returns a positive number;
	// if left is less then right, it returns a negative number; if the two input are equal, it returns 0
	Comparator func(left T, right T) int
	Pair       struct {
		First  T
		Second R
	}
)

var (
	// IntComparator is a Comparator for int
	IntComparator Comparator = func(left, right T) int {
		if left.(int) < right.(int) {
			return -1
		}
		if left.(int) > right.(int) {
			return 1
		}
		return 0
	}
	// Int64Comparator is a Comparator for int64
	Int64Comparator Comparator = func(left, right T) int {
		if left.(int64) > right.(int64) {
			return 1
		}
		if left.(int64) < right.(int64) {
			return -1
		}
		return 0
	}
)

// ReverseOrder returns the reverse order of the input Comparator
func ReverseOrder(cmp Comparator) Comparator {
	return func(left, right T) int {
		return cmp(right, left)
	}
}
