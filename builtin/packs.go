package builtin

import (
	"reflect"
)

// PackFunction represents a function to be registered in a pack.
type PackFunction struct {
	Name  string
	Fn    func(args ...any) (any, error)
	CtxFn func(env any, args ...any) (any, error)
	Types []reflect.Type
	IsCtx bool
}

// PackFuncTypes is a helper to create []reflect.Type from function pointer types.
func PackFuncTypes(types ...any) []reflect.Type {
	ts := make([]reflect.Type, len(types))
	for i, t := range types {
		t := reflect.TypeOf(t)
		if t.Kind() == reflect.Ptr {
			t = t.Elem()
		}
		ts[i] = t
	}
	return ts
}
