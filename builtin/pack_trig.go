package builtin

import (
	"fmt"
	"math"
)

// SIN returns the sine of n (in radians). nil input returns nil.
func SIN(args ...any) (any, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("SIN expects 1 argument, got %d", len(args))
	}
	n, wasNil := CoerceToFloat64(args[0])
	if wasNil {
		return nil, nil
	}
	r := math.Sin(n)
	if math.IsNaN(r) {
		return nil, nil
	}
	return r, nil
}

// COS returns the cosine of n (in radians). nil input returns nil.
func COS(args ...any) (any, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("COS expects 1 argument, got %d", len(args))
	}
	n, wasNil := CoerceToFloat64(args[0])
	if wasNil {
		return nil, nil
	}
	r := math.Cos(n)
	if math.IsNaN(r) {
		return nil, nil
	}
	return r, nil
}

// TAN returns the tangent of n (in radians). nil input returns nil.
func TAN(args ...any) (any, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("TAN expects 1 argument, got %d", len(args))
	}
	n, wasNil := CoerceToFloat64(args[0])
	if wasNil {
		return nil, nil
	}
	r := math.Tan(n)
	if math.IsNaN(r) {
		return nil, nil
	}
	return r, nil
}

// ASIN returns the arcsine of n. |n|>1 returns nil. nil input returns nil.
func ASIN(args ...any) (any, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("ASIN expects 1 argument, got %d", len(args))
	}
	n, wasNil := CoerceToFloat64(args[0])
	if wasNil {
		return nil, nil
	}
	if n < -1 || n > 1 {
		return nil, nil
	}
	r := math.Asin(n)
	if math.IsNaN(r) {
		return nil, nil
	}
	return r, nil
}

// ACOS returns the arccosine of n. |n|>1 returns nil. nil input returns nil.
func ACOS(args ...any) (any, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("ACOS expects 1 argument, got %d", len(args))
	}
	n, wasNil := CoerceToFloat64(args[0])
	if wasNil {
		return nil, nil
	}
	if n < -1 || n > 1 {
		return nil, nil
	}
	r := math.Acos(n)
	if math.IsNaN(r) {
		return nil, nil
	}
	return r, nil
}

// ATAN returns the arctangent of n. nil input returns nil.
func ATAN(args ...any) (any, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("ATAN expects 1 argument, got %d", len(args))
	}
	n, wasNil := CoerceToFloat64(args[0])
	if wasNil {
		return nil, nil
	}
	r := math.Atan(n)
	if math.IsNaN(r) {
		return nil, nil
	}
	return r, nil
}

// ATAN2 returns the arctangent of y/x using the signs to determine the quadrant.
// Either nil input returns nil.
func ATAN2(args ...any) (any, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("ATAN2 expects 2 arguments, got %d", len(args))
	}
	y, yNil := CoerceToFloat64(args[0])
	if yNil {
		return nil, nil
	}
	x, xNil := CoerceToFloat64(args[1])
	if xNil {
		return nil, nil
	}
	r := math.Atan2(y, x)
	if math.IsNaN(r) {
		return nil, nil
	}
	return r, nil
}

// PI returns the mathematical constant pi. No arguments.
func PI(args ...any) (any, error) {
	return math.Pi, nil
}

// TrigFunctions returns all trigonometric pack functions.
func TrigFunctions() []*PackFunction {
	variadicType := PackFuncTypes(new(func(...any) any))
	piType := PackFuncTypes(new(func() float64))
	return []*PackFunction{
		{Name: "SIN", Fn: SIN, Types: variadicType},
		{Name: "COS", Fn: COS, Types: variadicType},
		{Name: "TAN", Fn: TAN, Types: variadicType},
		{Name: "ASIN", Fn: ASIN, Types: variadicType},
		{Name: "ACOS", Fn: ACOS, Types: variadicType},
		{Name: "ATAN", Fn: ATAN, Types: variadicType},
		{Name: "ATAN2", Fn: ATAN2, Types: variadicType},
		{Name: "PI", Fn: PI, Types: piType},
	}
}
