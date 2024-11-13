package gomad

import (
	"fmt"
	"reflect"
)

type TypedErrMonad[T any] struct {
	Val T
	Err error
}

func NewTypedMonad[T any](val T) TypedErrMonad[T] {
	return TypedErrMonad[T]{Val: val, Err: nil}
}
func newErrTypedErrMonad[T any](err error) TypedErrMonad[T] {
	var zero T
	return TypedErrMonad[T]{Val: zero, Err: err}
}
func (self TypedErrMonad[T]) Bind(f interface{}, args ...interface{}) TypedErrMonad[T] {
	// Propagate previous error if any
	if self.Err != nil {
		return self
	}
	// Check if f is a function
	fnVal := reflect.ValueOf(f)
	if fnVal.Kind() != reflect.Func {
		return newErrTypedErrMonad[T](fmt.Errorf("Bind expects a function"))
	}
	// Ensure the function returns either `T` and `error`, or just `T`, or just `error`.
	if fnVal.Type().NumOut() > 2 || (fnVal.Type().NumOut() == 1 && fnVal.Type().Out(0) != reflect.TypeOf((*error)(nil)).Elem() && fnVal.Type().Out(0) != reflect.TypeOf(self.Val)) {
		return newErrTypedErrMonad[T](fmt.Errorf("Function must return T, error, or just error"))
	}
	// Check if len of args match to what function f expects
	if fnVal.Type().NumIn() != len(args) {
		return newErrTypedErrMonad[T](fmt.Errorf("Bound function expected %d arguments, got %d", fnVal.Type().NumIn(), len(args)))
	}
	// Convert args to reflect values and call f
	fnArgs := make([]reflect.Value, len(args))
	for i, arg := range args {
		fnArgs[i] = reflect.ValueOf(arg)
	}
	res := fnVal.Call(fnArgs)

	// Handle the return values
	switch len(res) {
	case 1:
		// If f only returned just `error` or just `T`
		if valRes, ok := res[0].Interface().(T); ok {
			return NewTypedMonad(valRes)
		} else if errRes, ok := res[0].Interface().(error); ok && errRes != nil {
			return newErrTypedErrMonad[T](errRes)
		}
		return newErrTypedErrMonad[T](fmt.Errorf("Failed to convert function output to T"))
	case 2:
		// If f returns `T` and `error`
		if errRes, ok := res[1].Interface().(error); ok && errRes != nil {
			return newErrTypedErrMonad[T](errRes)
		}
	}
	return self
}
func (self TypedErrMonad[T]) ValueOr(val T) T {
	if self.Err == nil {
		return self.Val
	}
	return val
}
