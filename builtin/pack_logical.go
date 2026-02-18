package builtin

import "fmt"

// AND returns true if ALL arguments are truthy.
// 0 args returns true.
func AND(args ...any) (any, error) {
	for _, a := range args {
		if !ToBool(a) {
			return false, nil
		}
	}
	return true, nil
}

// OR returns true if ANY argument is truthy.
// 0 args returns false.
func OR(args ...any) (any, error) {
	for _, a := range args {
		if ToBool(a) {
			return true, nil
		}
	}
	return false, nil
}

// NOT returns the logical negation of a single argument.
func NOT(args ...any) (any, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("NOT expects 1 argument, got %d", len(args))
	}
	return !ToBool(args[0]), nil
}

// IF evaluates a condition and returns one of two values.
// Args: cond, trueVal, falseVal.
func IF(args ...any) (any, error) {
	if len(args) != 3 {
		return nil, fmt.Errorf("IF expects 3 arguments, got %d", len(args))
	}
	if ToBool(args[0]) {
		return args[1], nil
	}
	return args[2], nil
}

// CASE compares an expression against value/result pairs with a default.
// Args: expr, val1, result1, val2, result2, ..., defaultResult.
func CASE(args ...any) (any, error) {
	if len(args) < 2 {
		return nil, nil
	}
	expr := args[0]
	rest := args[1:]

	// rest must have an odd number of elements: pairs + default
	if len(rest)%2 == 0 {
		// Malformed: even number means no default (odd remaining after expr)
		return nil, nil
	}

	// Process value/result pairs
	pairs := len(rest) - 1
	for i := 0; i < pairs; i += 2 {
		if EqualsCoerce(expr, rest[i]) {
			return rest[i+1], nil
		}
	}

	// Return default (last element)
	return rest[len(rest)-1], nil
}

// ISBLANK returns true if the value is nil, empty string, or whitespace-only string.
func ISBLANK(args ...any) (any, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("ISBLANK expects 1 argument, got %d", len(args))
	}
	return IsBlankVal(args[0]), nil
}

// ISNULL returns true if the value is nil, empty string, or whitespace-only string.
// Identical to ISBLANK.
func ISNULL(args ...any) (any, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("ISNULL expects 1 argument, got %d", len(args))
	}
	return IsBlankVal(args[0]), nil
}

// ISNUMBER returns true if the value is numeric or can be parsed as a number.
func ISNUMBER(args ...any) (any, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("ISNUMBER expects 1 argument, got %d", len(args))
	}
	return IsNumericVal(args[0]), nil
}

// LogicalFunctions returns all logical pack functions.
func LogicalFunctions() []*PackFunction {
	variadicType := PackFuncTypes(new(func(...any) any))
	return []*PackFunction{
		{Name: "AND", Fn: AND, Types: variadicType},
		{Name: "OR", Fn: OR, Types: variadicType},
		{Name: "NOT", Fn: NOT, Types: variadicType},
		{Name: "IF", Fn: IF, Types: variadicType},
		{Name: "CASE", Fn: CASE, Types: variadicType},
		{Name: "ISBLANK", Fn: ISBLANK, Types: variadicType},
		{Name: "ISNULL", Fn: ISNULL, Types: variadicType},
		{Name: "ISNUMBER", Fn: ISNUMBER, Types: variadicType},
	}
}
