package builtin

import (
	"fmt"
	"math"
)

// ABS returns the absolute value. nil input returns nil.
func ABS(args ...any) (any, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("ABS expects 1 argument, got %d", len(args))
	}
	n, wasNil := CoerceToFloat64(args[0])
	if wasNil {
		return nil, nil
	}
	r := math.Abs(n)
	if math.IsNaN(r) {
		return nil, nil
	}
	return r, nil
}

// CEILING returns the smallest integer >= n. nil input returns nil.
func CEILING(args ...any) (any, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("CEILING expects 1 argument, got %d", len(args))
	}
	n, wasNil := CoerceToFloat64(args[0])
	if wasNil {
		return nil, nil
	}
	r := math.Ceil(n)
	if math.IsNaN(r) {
		return nil, nil
	}
	return r, nil
}

// FLOOR returns the largest integer <= n. nil input returns nil.
func FLOOR(args ...any) (any, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("FLOOR expects 1 argument, got %d", len(args))
	}
	n, wasNil := CoerceToFloat64(args[0])
	if wasNil {
		return nil, nil
	}
	r := math.Floor(n)
	if math.IsNaN(r) {
		return nil, nil
	}
	return r, nil
}

// ROUND rounds n to the given number of decimal places using half-away-from-zero.
// nil input returns nil.
func ROUND(args ...any) (any, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("ROUND expects 2 arguments, got %d", len(args))
	}
	n, wasNil := CoerceToFloat64(args[0])
	if wasNil {
		return nil, nil
	}
	p, wasNil := CoerceToFloat64(args[1])
	if wasNil {
		return nil, nil
	}
	scale := math.Pow(10, p)
	var r float64
	if n >= 0 {
		r = math.Floor(n*scale+0.5) / scale
	} else {
		r = -math.Floor(-n*scale+0.5) / scale
	}
	if math.IsNaN(r) {
		return nil, nil
	}
	return r, nil
}

// MAX returns the larger of two values. Both nil returns nil; one nil returns the other.
func MAX(args ...any) (any, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("MAX expects 2 arguments, got %d", len(args))
	}
	a, aNil := CoerceToFloat64(args[0])
	b, bNil := CoerceToFloat64(args[1])
	if aNil && bNil {
		return nil, nil
	}
	if aNil {
		return b, nil
	}
	if bNil {
		return a, nil
	}
	if a >= b {
		return a, nil
	}
	return b, nil
}

// MIN returns the smaller of two values. Both nil returns nil; one nil returns the other.
func MIN(args ...any) (any, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("MIN expects 2 arguments, got %d", len(args))
	}
	a, aNil := CoerceToFloat64(args[0])
	b, bNil := CoerceToFloat64(args[1])
	if aNil && bNil {
		return nil, nil
	}
	if aNil {
		return b, nil
	}
	if bNil {
		return a, nil
	}
	if a <= b {
		return a, nil
	}
	return b, nil
}

// MOD returns the remainder of n / divisor. divisor=0 returns nil. nil input returns nil.
func MOD(args ...any) (any, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("MOD expects 2 arguments, got %d", len(args))
	}
	n, wasNil := CoerceToFloat64(args[0])
	if wasNil {
		return nil, nil
	}
	d, wasNil := CoerceToFloat64(args[1])
	if wasNil {
		return nil, nil
	}
	if d == 0 {
		return nil, nil
	}
	r := math.Mod(n, d)
	if math.IsNaN(r) {
		return nil, nil
	}
	return r, nil
}

// SQRT returns the square root of n. n<0 returns nil. nil input returns nil.
func SQRT(args ...any) (any, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("SQRT expects 1 argument, got %d", len(args))
	}
	n, wasNil := CoerceToFloat64(args[0])
	if wasNil {
		return nil, nil
	}
	if n < 0 {
		return nil, nil
	}
	r := math.Sqrt(n)
	if math.IsNaN(r) {
		return nil, nil
	}
	return r, nil
}

// LOG returns the base-10 logarithm. n<=0 returns nil. nil input returns nil.
func LOG(args ...any) (any, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("LOG expects 1 argument, got %d", len(args))
	}
	n, wasNil := CoerceToFloat64(args[0])
	if wasNil {
		return nil, nil
	}
	if n <= 0 {
		return nil, nil
	}
	r := math.Log10(n)
	if math.IsNaN(r) {
		return nil, nil
	}
	return r, nil
}

// LN returns the natural logarithm. n<=0 returns nil. nil input returns nil.
func LN(args ...any) (any, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("LN expects 1 argument, got %d", len(args))
	}
	n, wasNil := CoerceToFloat64(args[0])
	if wasNil {
		return nil, nil
	}
	if n <= 0 {
		return nil, nil
	}
	r := math.Log(n)
	if math.IsNaN(r) {
		return nil, nil
	}
	return r, nil
}

// EXP returns e raised to the power of n. nil input returns nil.
func EXP(args ...any) (any, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("EXP expects 1 argument, got %d", len(args))
	}
	n, wasNil := CoerceToFloat64(args[0])
	if wasNil {
		return nil, nil
	}
	r := math.Exp(n)
	if math.IsNaN(r) {
		return nil, nil
	}
	return r, nil
}

// MCEILING returns ceiling away from zero. n>=0: Ceil, n<0: Floor. nil input returns nil.
func MCEILING(args ...any) (any, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("MCEILING expects 1 argument, got %d", len(args))
	}
	n, wasNil := CoerceToFloat64(args[0])
	if wasNil {
		return nil, nil
	}
	var r float64
	if n >= 0 {
		r = math.Ceil(n)
	} else {
		r = math.Floor(n)
	}
	if math.IsNaN(r) {
		return nil, nil
	}
	return r, nil
}

// MFLOOR returns floor toward zero. n>=0: Floor, n<0: Ceil. nil input returns nil.
func MFLOOR(args ...any) (any, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("MFLOOR expects 1 argument, got %d", len(args))
	}
	n, wasNil := CoerceToFloat64(args[0])
	if wasNil {
		return nil, nil
	}
	var r float64
	if n >= 0 {
		r = math.Floor(n)
	} else {
		r = math.Ceil(n)
	}
	if math.IsNaN(r) {
		return nil, nil
	}
	return r, nil
}

// TRUNC truncates n toward zero to the given number of decimal places.
// places defaults to 0 if only 1 argument is provided. nil input returns nil.
func TRUNC(args ...any) (any, error) {
	if len(args) < 1 || len(args) > 2 {
		return nil, fmt.Errorf("TRUNC expects 1 or 2 arguments, got %d", len(args))
	}
	n, wasNil := CoerceToFloat64(args[0])
	if wasNil {
		return nil, nil
	}
	var places float64
	if len(args) == 2 {
		p, pNil := CoerceToFloat64(args[1])
		if pNil {
			return nil, nil
		}
		places = p
	}
	scale := math.Pow(10, places)
	var r float64
	if n >= 0 {
		r = math.Floor(n*scale) / scale
	} else {
		r = math.Ceil(n*scale) / scale
	}
	if math.IsNaN(r) {
		return nil, nil
	}
	return r, nil
}

// MathFunctions returns all math pack functions.
func MathFunctions() []*PackFunction {
	variadicType := PackFuncTypes(new(func(...any) any))
	return []*PackFunction{
		{Name: "ABS", Fn: ABS, Types: variadicType},
		{Name: "CEILING", Fn: CEILING, Types: variadicType},
		{Name: "FLOOR", Fn: FLOOR, Types: variadicType},
		{Name: "ROUND", Fn: ROUND, Types: variadicType},
		{Name: "MAX", Fn: MAX, Types: variadicType},
		{Name: "MIN", Fn: MIN, Types: variadicType},
		{Name: "MOD", Fn: MOD, Types: variadicType},
		{Name: "SQRT", Fn: SQRT, Types: variadicType},
		{Name: "LOG", Fn: LOG, Types: variadicType},
		{Name: "LN", Fn: LN, Types: variadicType},
		{Name: "EXP", Fn: EXP, Types: variadicType},
		{Name: "MCEILING", Fn: MCEILING, Types: variadicType},
		{Name: "MFLOOR", Fn: MFLOOR, Types: variadicType},
		{Name: "TRUNC", Fn: TRUNC, Types: variadicType},
	}
}
