package builtin

import (
	"fmt"
	"strings"
)

// ISPICKVAL checks if a picklist field equals a specific value.
// ISPICKVAL(field, value) → bool
// If field is nil and value is "", returns true (SF behavior for empty picklists).
// If field is nil and value is non-empty, returns false.
func sfISPICKVAL(args ...any) (any, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("ISPICKVAL expects 2 arguments, got %d", len(args))
	}
	field := args[0]
	value := CoerceToString(args[1])

	if IsNilVal(field) {
		return value == "", nil
	}
	return CoerceToString(field) == value, nil
}

// INCLUDES checks if a multi-select picklist contains a specific value.
// Multi-select values are semicolon-separated strings.
// INCLUDES(multiSelectField, value) → bool
func sfINCLUDES(args ...any) (any, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("INCLUDES expects 2 arguments, got %d", len(args))
	}
	if IsNilVal(args[0]) || IsNilVal(args[1]) {
		return false, nil
	}
	field := CoerceToString(args[0])
	value := CoerceToString(args[1])
	// SF multi-select values are semicolon-separated
	parts := strings.Split(field, ";")
	for _, p := range parts {
		if strings.TrimSpace(p) == value {
			return true, nil
		}
	}
	return false, nil
}

// sfISCHANGED is a context function that checks if a field has been changed.
// It expects the env to contain a map "_changed" with field names as keys.
// The first argument is the field value (ignored), but the function name
// and usage pattern will be transformed by the formula compiler.
// For now, it receives the current value and checks the _changed map by value identity.
func sfISCHANGED(env any, args ...any) (any, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("ISCHANGED expects 1 argument, got %d", len(args))
	}
	// In the real implementation, the formula compiler would pass the field name.
	// For testing, we look up the value in the _changed map from env.
	if envMap, ok := env.(map[string]any); ok {
		if changed, ok := envMap["_changed"]; ok {
			if changedMap, ok := changed.(map[string]bool); ok {
				// The argument is the current field value. We can't determine the field name
				// from the value alone. In real usage, the formula compiler rewrites this.
				// For testing, we pass a string field name as the argument.
				if fieldName, ok := args[0].(string); ok {
					return changedMap[fieldName], nil
				}
			}
		}
	}
	return false, nil
}

// sfPRIORVALUE is a context function that returns the prior value of a field.
// It expects the env to contain a map "_prior" with field names as keys.
func sfPRIORVALUE(env any, args ...any) (any, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("PRIORVALUE expects 1 argument, got %d", len(args))
	}
	if envMap, ok := env.(map[string]any); ok {
		if prior, ok := envMap["_prior"]; ok {
			if priorMap, ok := prior.(map[string]any); ok {
				if fieldName, ok := args[0].(string); ok {
					return priorMap[fieldName], nil
				}
			}
		}
	}
	return nil, nil
}

// sfISNEW is a context function that checks if the record is new.
// It expects the env to contain a bool "_isNew".
func sfISNEW(env any, args ...any) (any, error) {
	if envMap, ok := env.(map[string]any); ok {
		if isNew, ok := envMap["_isNew"]; ok {
			return ToBool(isNew), nil
		}
	}
	return false, nil
}

// BLANKVALUE returns the substitute value if the expression is blank
// (nil, empty string, or whitespace-only string), otherwise returns the expression value.
// BLANKVALUE(expr, substituteExpr) → any
func sfBLANKVALUE(args ...any) (any, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("BLANKVALUE expects 2 arguments, got %d", len(args))
	}
	if IsBlankVal(args[0]) {
		return args[1], nil
	}
	return args[0], nil
}

// NULLVALUE returns the substitute value if the expression is null,
// otherwise returns the expression value.
// NULLVALUE(expr, substituteExpr) → any
func sfNULLVALUE(args ...any) (any, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("NULLVALUE expects 2 arguments, got %d", len(args))
	}
	if IsNilVal(args[0]) {
		return args[1], nil
	}
	return args[0], nil
}

// SalesforceRuntimeFunctions returns the SF-specific runtime functions:
// ISPICKVAL, INCLUDES, ISCHANGED, PRIORVALUE, ISNEW, BLANKVALUE, NULLVALUE.
func SalesforceRuntimeFunctions() []*PackFunction {
	variadic := PackFuncTypes(new(func(...any) any))
	return []*PackFunction{
		{
			Name:  "ISPICKVAL",
			Fn:    sfISPICKVAL,
			Types: variadic,
		},
		{
			Name:  "INCLUDES",
			Fn:    sfINCLUDES,
			Types: variadic,
		},
		{
			Name:  "ISCHANGED",
			CtxFn: sfISCHANGED,
			Types: variadic,
			IsCtx: true,
		},
		{
			Name:  "PRIORVALUE",
			CtxFn: sfPRIORVALUE,
			Types: variadic,
			IsCtx: true,
		},
		{
			Name:  "ISNEW",
			CtxFn: sfISNEW,
			Types: PackFuncTypes(new(func() any)),
			IsCtx: true,
		},
		{
			Name:  "BLANKVALUE",
			Fn:    sfBLANKVALUE,
			Types: variadic,
		},
		{
			Name:  "NULLVALUE",
			Fn:    sfNULLVALUE,
			Types: variadic,
		},
	}
}
