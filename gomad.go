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
func (self *TypedErrMonad[T]) Bind(f interface{}, args ...interface{}) *TypedErrMonad[T] {
	var zero T
	// Propagate previous error if any
	if self.Err != nil {
		return self
	}
	// Check if f is a function
	fnVal := reflect.ValueOf(f)
	if fnVal.Kind() != reflect.Func {
		self.Err = fmt.Errorf("Bind expects a function")
		self.Val = zero
		return self
	}

	// Retrieve the output types of the function
	numOut := fnVal.Type().NumOut()
	if numOut > 2 || (numOut == 1 && fnVal.Type().Out(0) != reflect.TypeOf(zero) && fnVal.Type().Out(0) != reflect.TypeOf((*error)(nil)).Elem()) {
		self.Err = fmt.Errorf("Function must return T, error, or just error")
		self.Val = zero
		return self
	}

	// Check if len of args match to what function f expects
	if fnVal.Type().NumIn() != len(args) {
		self.Err = fmt.Errorf("Bound function expected %d arguments, got %d", fnVal.Type().NumIn(), len(args))
		self.Val = zero
		return self
	}

	// Convert args to reflect values and call f
	fnArgs := make([]reflect.Value, len(args))
	for i, arg := range args {
		fnArgs[i] = reflect.ValueOf(arg)
	}
	res := fnVal.Call(fnArgs)

	// Handle the return values based on the function's signature
	switch numOut {
	case 1:
		// Function returns just `T` or `error`
		switch fnVal.Type().Out(0) {
		case reflect.TypeOf(zero):
			// Output type is `T`
			if valRes, ok := res[0].Interface().(T); ok {
				self.Val = valRes
				self.Err = nil
			} else {
				self.Err = fmt.Errorf("Failed to convert function output to T")
				self.Val = zero
			}
		case reflect.TypeOf((*error)(nil)).Elem():
			// Output type `nil error`
			if res[0].Interface() == nil {
				self.Err = nil
			} else if errRes, ok := res[0].Interface().(error); ok {
				// Output type `error`
				self.Err = errRes
				self.Val = zero
			} else {
				self.Err = fmt.Errorf("Expected an error type return")
			}
		}
	case 2:
		// Function returns `T` and `error`
		if fnVal.Type().Out(0) == reflect.TypeOf(zero) && fnVal.Type().Out(1) == reflect.TypeOf((*error)(nil)).Elem() {
			if valRes, ok := res[0].Interface().(T); ok {
				self.Val = valRes
				// Output type `nil error`
				if res[1].Interface() == nil {
					self.Err = nil
				} else if errRes, ok := res[1].Interface().(error); ok {
					// Output type `error`
					self.Err = errRes
					self.Val = zero
				} else {
					self.Err = fmt.Errorf("Failed to convert function output to error")
					self.Val = zero
				}
			} else {
				self.Err = fmt.Errorf("Failed to convert function output to T")
				self.Val = zero
			}
		} else {
			self.Err = fmt.Errorf("Function return types do not match expected T and error")
			self.Val = zero
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
