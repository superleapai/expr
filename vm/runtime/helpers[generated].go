package runtime

import (
	"fmt"
	"math"
	"reflect"
	"strconv"
	"time"
)

// Safe version of ToFloat64 that returns error instead of panicking
func ToFloat64Safe(a any) (float64, bool) {
	if IsNil(a) {
		return 0, true
	}
	switch x := a.(type) {
	case bool:
		if x {
			return 1, true
		}
		return 0, true
	case string:
		if x == "" {
			return 0.0, true // empty string converts to 0.0
		}
		if f, err := strconv.ParseFloat(x, 64); err == nil {
			return f, true
		}
		return 0, false
	case float32:
		return float64(x), true
	case float64:
		return x, true
	case int:
		return float64(x), true
	case int8:
		return float64(x), true
	case int16:
		return float64(x), true
	case int32:
		return float64(x), true
	case int64:
		return float64(x), true
	case uint:
		return float64(x), true
	case uint8:
		return float64(x), true
	case uint16:
		return float64(x), true
	case uint32:
		return float64(x), true
	case uint64:
		return float64(x), true
	default:
		return 0, false
	}
}

func Equal(a, b interface{}) bool {
	// Handle nil values first - nil only equals nil
	if IsNil(a) && IsNil(b) {
		return true
	}
	if IsNil(a) || IsNil(b) {
		return false
	}

	// Handle numeric operations
	switch x := a.(type) {
	case uint:
		switch y := b.(type) {
		case uint:
			return int(x) == int(y)
		case uint8:
			return int(x) == int(y)
		case uint16:
			return int(x) == int(y)
		case uint32:
			return int(x) == int(y)
		case uint64:
			return int(x) == int(y)
		case int:
			return int(x) == int(y)
		case int8:
			return int(x) == int(y)
		case int16:
			return int(x) == int(y)
		case int32:
			return int(x) == int(y)
		case int64:
			return int(x) == int(y)
		case float32:
			return float64(x) == float64(y)
		case float64:
			return float64(x) == float64(y)
		case bool:
			return ToFloat64(x) == ToFloat64(y)
		case string:
			return ToFloat64(x) == ToFloat64(y)
		}
	case uint8:
		switch y := b.(type) {
		case uint:
			return int(x) == int(y)
		case uint8:
			return int(x) == int(y)
		case uint16:
			return int(x) == int(y)
		case uint32:
			return int(x) == int(y)
		case uint64:
			return int(x) == int(y)
		case int:
			return int(x) == int(y)
		case int8:
			return int(x) == int(y)
		case int16:
			return int(x) == int(y)
		case int32:
			return int(x) == int(y)
		case int64:
			return int(x) == int(y)
		case float32:
			return float64(x) == float64(y)
		case float64:
			return float64(x) == float64(y)
		case bool:
			return ToFloat64(x) == ToFloat64(y)
		case string:
			return ToFloat64(x) == ToFloat64(y)
		}
	case uint16:
		switch y := b.(type) {
		case uint:
			return int(x) == int(y)
		case uint8:
			return int(x) == int(y)
		case uint16:
			return int(x) == int(y)
		case uint32:
			return int(x) == int(y)
		case uint64:
			return int(x) == int(y)
		case int:
			return int(x) == int(y)
		case int8:
			return int(x) == int(y)
		case int16:
			return int(x) == int(y)
		case int32:
			return int(x) == int(y)
		case int64:
			return int(x) == int(y)
		case float32:
			return float64(x) == float64(y)
		case float64:
			return float64(x) == float64(y)
		case bool:
			return ToFloat64(x) == ToFloat64(y)
		case string:
			return ToFloat64(x) == ToFloat64(y)
		}
	case uint32:
		switch y := b.(type) {
		case uint:
			return int(x) == int(y)
		case uint8:
			return int(x) == int(y)
		case uint16:
			return int(x) == int(y)
		case uint32:
			return int(x) == int(y)
		case uint64:
			return int(x) == int(y)
		case int:
			return int(x) == int(y)
		case int8:
			return int(x) == int(y)
		case int16:
			return int(x) == int(y)
		case int32:
			return int(x) == int(y)
		case int64:
			return int(x) == int(y)
		case float32:
			return float64(x) == float64(y)
		case float64:
			return float64(x) == float64(y)
		case bool:
			return ToFloat64(x) == ToFloat64(y)
		case string:
			return ToFloat64(x) == ToFloat64(y)
		}
	case uint64:
		switch y := b.(type) {
		case uint:
			return int(x) == int(y)
		case uint8:
			return int(x) == int(y)
		case uint16:
			return int(x) == int(y)
		case uint32:
			return int(x) == int(y)
		case uint64:
			return int(x) == int(y)
		case int:
			return int(x) == int(y)
		case int8:
			return int(x) == int(y)
		case int16:
			return int(x) == int(y)
		case int32:
			return int(x) == int(y)
		case int64:
			return int(x) == int(y)
		case float32:
			return float64(x) == float64(y)
		case float64:
			return float64(x) == float64(y)
		case bool:
			return ToFloat64(x) == ToFloat64(y)
		case string:
			return ToFloat64(x) == ToFloat64(y)
		}
	case int:
		switch y := b.(type) {
		case uint:
			return int(x) == int(y)
		case uint8:
			return int(x) == int(y)
		case uint16:
			return int(x) == int(y)
		case uint32:
			return int(x) == int(y)
		case uint64:
			return int(x) == int(y)
		case int:
			return int(x) == int(y)
		case int8:
			return int(x) == int(y)
		case int16:
			return int(x) == int(y)
		case int32:
			return int(x) == int(y)
		case int64:
			return int(x) == int(y)
		case float32:
			return float64(x) == float64(y)
		case float64:
			return float64(x) == float64(y)
		case bool:
			return ToFloat64(x) == ToFloat64(y)
		case string:
			return ToFloat64(x) == ToFloat64(y)
		}
	case int8:
		switch y := b.(type) {
		case uint:
			return int(x) == int(y)
		case uint8:
			return int(x) == int(y)
		case uint16:
			return int(x) == int(y)
		case uint32:
			return int(x) == int(y)
		case uint64:
			return int(x) == int(y)
		case int:
			return int(x) == int(y)
		case int8:
			return int(x) == int(y)
		case int16:
			return int(x) == int(y)
		case int32:
			return int(x) == int(y)
		case int64:
			return int(x) == int(y)
		case float32:
			return float64(x) == float64(y)
		case float64:
			return float64(x) == float64(y)
		case bool:
			return ToFloat64(x) == ToFloat64(y)
		case string:
			return ToFloat64(x) == ToFloat64(y)
		}
	case int16:
		switch y := b.(type) {
		case uint:
			return int(x) == int(y)
		case uint8:
			return int(x) == int(y)
		case uint16:
			return int(x) == int(y)
		case uint32:
			return int(x) == int(y)
		case uint64:
			return int(x) == int(y)
		case int:
			return int(x) == int(y)
		case int8:
			return int(x) == int(y)
		case int16:
			return int(x) == int(y)
		case int32:
			return int(x) == int(y)
		case int64:
			return int(x) == int(y)
		case float32:
			return float64(x) == float64(y)
		case float64:
			return float64(x) == float64(y)
		case bool:
			return ToFloat64(x) == ToFloat64(y)
		case string:
			return ToFloat64(x) == ToFloat64(y)
		}
	case int32:
		switch y := b.(type) {
		case uint:
			return int(x) == int(y)
		case uint8:
			return int(x) == int(y)
		case uint16:
			return int(x) == int(y)
		case uint32:
			return int(x) == int(y)
		case uint64:
			return int(x) == int(y)
		case int:
			return int(x) == int(y)
		case int8:
			return int(x) == int(y)
		case int16:
			return int(x) == int(y)
		case int32:
			return int(x) == int(y)
		case int64:
			return int(x) == int(y)
		case float32:
			return float64(x) == float64(y)
		case float64:
			return float64(x) == float64(y)
		case bool:
			return ToFloat64(x) == ToFloat64(y)
		case string:
			return ToFloat64(x) == ToFloat64(y)
		}
	case int64:
		switch y := b.(type) {
		case uint:
			return int(x) == int(y)
		case uint8:
			return int(x) == int(y)
		case uint16:
			return int(x) == int(y)
		case uint32:
			return int(x) == int(y)
		case uint64:
			return int(x) == int(y)
		case int:
			return int(x) == int(y)
		case int8:
			return int(x) == int(y)
		case int16:
			return int(x) == int(y)
		case int32:
			return int(x) == int(y)
		case int64:
			return int(x) == int(y)
		case float32:
			return float64(x) == float64(y)
		case float64:
			return float64(x) == float64(y)
		case bool:
			return ToFloat64(x) == ToFloat64(y)
		case string:
			return ToFloat64(x) == ToFloat64(y)
		}
	case float32:
		switch y := b.(type) {
		case uint:
			return float64(x) == float64(y)
		case uint8:
			return float64(x) == float64(y)
		case uint16:
			return float64(x) == float64(y)
		case uint32:
			return float64(x) == float64(y)
		case uint64:
			return float64(x) == float64(y)
		case int:
			return float64(x) == float64(y)
		case int8:
			return float64(x) == float64(y)
		case int16:
			return float64(x) == float64(y)
		case int32:
			return float64(x) == float64(y)
		case int64:
			return float64(x) == float64(y)
		case float32:
			return float64(x) == float64(y)
		case float64:
			return float64(x) == float64(y)
		case bool:
			return ToFloat64(x) == ToFloat64(y)
		case string:
			return ToFloat64(x) == ToFloat64(y)
		}
	case float64:
		switch y := b.(type) {
		case uint:
			return float64(x) == float64(y)
		case uint8:
			return float64(x) == float64(y)
		case uint16:
			return float64(x) == float64(y)
		case uint32:
			return float64(x) == float64(y)
		case uint64:
			return float64(x) == float64(y)
		case int:
			return float64(x) == float64(y)
		case int8:
			return float64(x) == float64(y)
		case int16:
			return float64(x) == float64(y)
		case int32:
			return float64(x) == float64(y)
		case int64:
			return float64(x) == float64(y)
		case float32:
			return float64(x) == float64(y)
		case float64:
			return float64(x) == float64(y)
		case bool:
			return ToFloat64(x) == ToFloat64(y)
		case string:
			return ToFloat64(x) == ToFloat64(y)
		}
	case []any:
		switch y := b.(type) {
		case []string:
			if len(x) != len(y) {
				return false
			}
			for i := range x {
				if !Equal(x[i], y[i]) {
					return false
				}
			}
			return true
		case []uint:
			if len(x) != len(y) {
				return false
			}
			for i := range x {
				if !Equal(x[i], y[i]) {
					return false
				}
			}
			return true
		case []uint8:
			if len(x) != len(y) {
				return false
			}
			for i := range x {
				if !Equal(x[i], y[i]) {
					return false
				}
			}
			return true
		case []uint16:
			if len(x) != len(y) {
				return false
			}
			for i := range x {
				if !Equal(x[i], y[i]) {
					return false
				}
			}
			return true
		case []uint32:
			if len(x) != len(y) {
				return false
			}
			for i := range x {
				if !Equal(x[i], y[i]) {
					return false
				}
			}
			return true
		case []uint64:
			if len(x) != len(y) {
				return false
			}
			for i := range x {
				if !Equal(x[i], y[i]) {
					return false
				}
			}
			return true
		case []int:
			if len(x) != len(y) {
				return false
			}
			for i := range x {
				if !Equal(x[i], y[i]) {
					return false
				}
			}
			return true
		case []int8:
			if len(x) != len(y) {
				return false
			}
			for i := range x {
				if !Equal(x[i], y[i]) {
					return false
				}
			}
			return true
		case []int16:
			if len(x) != len(y) {
				return false
			}
			for i := range x {
				if !Equal(x[i], y[i]) {
					return false
				}
			}
			return true
		case []int32:
			if len(x) != len(y) {
				return false
			}
			for i := range x {
				if !Equal(x[i], y[i]) {
					return false
				}
			}
			return true
		case []int64:
			if len(x) != len(y) {
				return false
			}
			for i := range x {
				if !Equal(x[i], y[i]) {
					return false
				}
			}
			return true
		case []float32:
			if len(x) != len(y) {
				return false
			}
			for i := range x {
				if !Equal(x[i], y[i]) {
					return false
				}
			}
			return true
		case []float64:
			if len(x) != len(y) {
				return false
			}
			for i := range x {
				if !Equal(x[i], y[i]) {
					return false
				}
			}
			return true
		case []any:
			if len(x) != len(y) {
				return false
			}
			for i := range x {
				if !Equal(x[i], y[i]) {
					return false
				}
			}
			return true
		}
	case []string:
		switch y := b.(type) {
		case []any:
			return Equal(y, x)
		case []string:
			if len(x) != len(y) {
				return false
			}
			for i := range x {
				if x[i] != y[i] {
					return false
				}
			}
			return true
		}
	case []uint:
		switch y := b.(type) {
		case []any:
			return Equal(y, x)
		case []uint:
			if len(x) != len(y) {
				return false
			}
			for i := range x {
				if x[i] != y[i] {
					return false
				}
			}
			return true
		}
	case []uint8:
		switch y := b.(type) {
		case []any:
			return Equal(y, x)
		case []uint8:
			if len(x) != len(y) {
				return false
			}
			for i := range x {
				if x[i] != y[i] {
					return false
				}
			}
			return true
		}
	case []uint16:
		switch y := b.(type) {
		case []any:
			return Equal(y, x)
		case []uint16:
			if len(x) != len(y) {
				return false
			}
			for i := range x {
				if x[i] != y[i] {
					return false
				}
			}
			return true
		}
	case []uint32:
		switch y := b.(type) {
		case []any:
			return Equal(y, x)
		case []uint32:
			if len(x) != len(y) {
				return false
			}
			for i := range x {
				if x[i] != y[i] {
					return false
				}
			}
			return true
		}
	case []uint64:
		switch y := b.(type) {
		case []any:
			return Equal(y, x)
		case []uint64:
			if len(x) != len(y) {
				return false
			}
			for i := range x {
				if x[i] != y[i] {
					return false
				}
			}
			return true
		}
	case []int:
		switch y := b.(type) {
		case []any:
			return Equal(y, x)
		case []int:
			if len(x) != len(y) {
				return false
			}
			for i := range x {
				if x[i] != y[i] {
					return false
				}
			}
			return true
		}
	case []int8:
		switch y := b.(type) {
		case []any:
			return Equal(y, x)
		case []int8:
			if len(x) != len(y) {
				return false
			}
			for i := range x {
				if x[i] != y[i] {
					return false
				}
			}
			return true
		}
	case []int16:
		switch y := b.(type) {
		case []any:
			return Equal(y, x)
		case []int16:
			if len(x) != len(y) {
				return false
			}
			for i := range x {
				if x[i] != y[i] {
					return false
				}
			}
			return true
		}
	case []int32:
		switch y := b.(type) {
		case []any:
			return Equal(y, x)
		case []int32:
			if len(x) != len(y) {
				return false
			}
			for i := range x {
				if x[i] != y[i] {
					return false
				}
			}
			return true
		}
	case []int64:
		switch y := b.(type) {
		case []any:
			return Equal(y, x)
		case []int64:
			if len(x) != len(y) {
				return false
			}
			for i := range x {
				if x[i] != y[i] {
					return false
				}
			}
			return true
		}
	case []float32:
		switch y := b.(type) {
		case []any:
			return Equal(y, x)
		case []float32:
			if len(x) != len(y) {
				return false
			}
			for i := range x {
				if x[i] != y[i] {
					return false
				}
			}
			return true
		}
	case []float64:
		switch y := b.(type) {
		case []any:
			return Equal(y, x)
		case []float64:
			if len(x) != len(y) {
				return false
			}
			for i := range x {
				if x[i] != y[i] {
					return false
				}
			}
			return true
		}
	case string:
		switch y := b.(type) {
		case string:
			return x == y
		case bool:
			// JavaScript-like coercion: empty string == false, non-empty string != true/false
			if x == "" {
				return y == false
			}
			return false
		case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64:
			// Try to convert string to number and compare
			if x == "" {
				return ToInt(y) == 0
			}
			return ToFloat64(x) == ToFloat64(y)
		}
	case time.Time:
		switch y := b.(type) {
		case time.Time:
			return x.Equal(y)
		}
	case time.Duration:
		switch y := b.(type) {
		case time.Duration:
			return x == y
		}
	case bool:
		switch y := b.(type) {
		case bool:
			return x == y
		case string:
			// JavaScript-like coercion: false == empty string
			if y == "" {
				return x == false
			}
			// Try to convert string to number
			return ToInt(x) == ToInt(y)
		case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64:
			// JavaScript-like coercion: true == 1, false == 0
			return ToInt(x) == ToInt(y)
		}
	}
	if IsNil(a) && IsNil(b) {
		return true
	}
	return reflect.DeepEqual(a, b)
}

// EqualIn performs JavaScript-like equality comparison for "in" operations
// This is specifically designed for array/slice membership testing
func EqualIn(a, b interface{}) bool {
	// Handle nil values first - nil only equals nil
	if IsNil(a) && IsNil(b) {
		return true
	}
	if IsNil(a) || IsNil(b) {
		return false
	}

	// Handle numeric operations with JavaScript-like coercion
	switch x := a.(type) {
	case uint:
		switch y := b.(type) {
		case uint:
			return int(x) == int(y)
		case uint8:
			return int(x) == int(y)
		case uint16:
			return int(x) == int(y)
		case uint32:
			return int(x) == int(y)
		case uint64:
			return int(x) == int(y)
		case int:
			return int(x) == int(y)
		case int8:
			return int(x) == int(y)
		case int16:
			return int(x) == int(y)
		case int32:
			return int(x) == int(y)
		case int64:
			return int(x) == int(y)
		case float32:
			return float64(x) == float64(y)
		case float64:
			return float64(x) == float64(y)
		case bool:
			return ToFloat64(x) == ToFloat64(y)
		case string:
			// Try numeric conversion first
			if val, ok := ToFloat64Safe(y); ok {
				return ToFloat64(x) == val
			}
			// String comparison
			return false
		}
	case uint8:
		switch y := b.(type) {
		case uint:
			return int(x) == int(y)
		case uint8:
			return int(x) == int(y)
		case uint16:
			return int(x) == int(y)
		case uint32:
			return int(x) == int(y)
		case uint64:
			return int(x) == int(y)
		case int:
			return int(x) == int(y)
		case int8:
			return int(x) == int(y)
		case int16:
			return int(x) == int(y)
		case int32:
			return int(x) == int(y)
		case int64:
			return int(x) == int(y)
		case float32:
			return float64(x) == float64(y)
		case float64:
			return float64(x) == float64(y)
		case bool:
			return ToFloat64(x) == ToFloat64(y)
		case string:
			if val, ok := ToFloat64Safe(y); ok {
				return ToFloat64(x) == val
			}
			return false
		}
	case uint16:
		switch y := b.(type) {
		case uint:
			return int(x) == int(y)
		case uint8:
			return int(x) == int(y)
		case uint16:
			return int(x) == int(y)
		case uint32:
			return int(x) == int(y)
		case uint64:
			return int(x) == int(y)
		case int:
			return int(x) == int(y)
		case int8:
			return int(x) == int(y)
		case int16:
			return int(x) == int(y)
		case int32:
			return int(x) == int(y)
		case int64:
			return int(x) == int(y)
		case float32:
			return float64(x) == float64(y)
		case float64:
			return float64(x) == float64(y)
		case bool:
			return ToFloat64(x) == ToFloat64(y)
		case string:
			if val, ok := ToFloat64Safe(y); ok {
				return ToFloat64(x) == val
			}
			return false
		}
	case uint32:
		switch y := b.(type) {
		case uint:
			return int(x) == int(y)
		case uint8:
			return int(x) == int(y)
		case uint16:
			return int(x) == int(y)
		case uint32:
			return int(x) == int(y)
		case uint64:
			return int(x) == int(y)
		case int:
			return int(x) == int(y)
		case int8:
			return int(x) == int(y)
		case int16:
			return int(x) == int(y)
		case int32:
			return int(x) == int(y)
		case int64:
			return int(x) == int(y)
		case float32:
			return float64(x) == float64(y)
		case float64:
			return float64(x) == float64(y)
		case bool:
			return ToFloat64(x) == ToFloat64(y)
		case string:
			if val, ok := ToFloat64Safe(y); ok {
				return ToFloat64(x) == val
			}
			return false
		}
	case uint64:
		switch y := b.(type) {
		case uint:
			return int(x) == int(y)
		case uint8:
			return int(x) == int(y)
		case uint16:
			return int(x) == int(y)
		case uint32:
			return int(x) == int(y)
		case uint64:
			return int(x) == int(y)
		case int:
			return int(x) == int(y)
		case int8:
			return int(x) == int(y)
		case int16:
			return int(x) == int(y)
		case int32:
			return int(x) == int(y)
		case int64:
			return int(x) == int(y)
		case float32:
			return float64(x) == float64(y)
		case float64:
			return float64(x) == float64(y)
		case bool:
			return ToFloat64(x) == ToFloat64(y)
		case string:
			if val, ok := ToFloat64Safe(y); ok {
				return ToFloat64(x) == val
			}
			return false
		}
	case int:
		switch y := b.(type) {
		case uint:
			return int(x) == int(y)
		case uint8:
			return int(x) == int(y)
		case uint16:
			return int(x) == int(y)
		case uint32:
			return int(x) == int(y)
		case uint64:
			return int(x) == int(y)
		case int:
			return int(x) == int(y)
		case int8:
			return int(x) == int(y)
		case int16:
			return int(x) == int(y)
		case int32:
			return int(x) == int(y)
		case int64:
			return int(x) == int(y)
		case float32:
			return float64(x) == float64(y)
		case float64:
			return float64(x) == float64(y)
		case bool:
			return ToFloat64(x) == ToFloat64(y)
		case string:
			if val, ok := ToFloat64Safe(y); ok {
				return ToFloat64(x) == val
			}
			return false
		}
	case int8:
		switch y := b.(type) {
		case uint:
			return int(x) == int(y)
		case uint8:
			return int(x) == int(y)
		case uint16:
			return int(x) == int(y)
		case uint32:
			return int(x) == int(y)
		case uint64:
			return int(x) == int(y)
		case int:
			return int(x) == int(y)
		case int8:
			return int(x) == int(y)
		case int16:
			return int(x) == int(y)
		case int32:
			return int(x) == int(y)
		case int64:
			return int(x) == int(y)
		case float32:
			return float64(x) == float64(y)
		case float64:
			return float64(x) == float64(y)
		case bool:
			return ToFloat64(x) == ToFloat64(y)
		case string:
			if val, ok := ToFloat64Safe(y); ok {
				return ToFloat64(x) == val
			}
			return false
		}
	case int16:
		switch y := b.(type) {
		case uint:
			return int(x) == int(y)
		case uint8:
			return int(x) == int(y)
		case uint16:
			return int(x) == int(y)
		case uint32:
			return int(x) == int(y)
		case uint64:
			return int(x) == int(y)
		case int:
			return int(x) == int(y)
		case int8:
			return int(x) == int(y)
		case int16:
			return int(x) == int(y)
		case int32:
			return int(x) == int(y)
		case int64:
			return int(x) == int(y)
		case float32:
			return float64(x) == float64(y)
		case float64:
			return float64(x) == float64(y)
		case bool:
			return ToFloat64(x) == ToFloat64(y)
		case string:
			if val, ok := ToFloat64Safe(y); ok {
				return ToFloat64(x) == val
			}
			return false
		}
	case int32:
		switch y := b.(type) {
		case uint:
			return int(x) == int(y)
		case uint8:
			return int(x) == int(y)
		case uint16:
			return int(x) == int(y)
		case uint32:
			return int(x) == int(y)
		case uint64:
			return int(x) == int(y)
		case int:
			return int(x) == int(y)
		case int8:
			return int(x) == int(y)
		case int16:
			return int(x) == int(y)
		case int32:
			return int(x) == int(y)
		case int64:
			return int(x) == int(y)
		case float32:
			return float64(x) == float64(y)
		case float64:
			return float64(x) == float64(y)
		case bool:
			return ToFloat64(x) == ToFloat64(y)
		case string:
			if val, ok := ToFloat64Safe(y); ok {
				return ToFloat64(x) == val
			}
			return false
		}
	case int64:
		switch y := b.(type) {
		case uint:
			return int(x) == int(y)
		case uint8:
			return int(x) == int(y)
		case uint16:
			return int(x) == int(y)
		case uint32:
			return int(x) == int(y)
		case uint64:
			return int(x) == int(y)
		case int:
			return int(x) == int(y)
		case int8:
			return int(x) == int(y)
		case int16:
			return int(x) == int(y)
		case int32:
			return int(x) == int(y)
		case int64:
			return int(x) == int(y)
		case float32:
			return float64(x) == float64(y)
		case float64:
			return float64(x) == float64(y)
		case bool:
			return ToFloat64(x) == ToFloat64(y)
		case string:
			if val, ok := ToFloat64Safe(y); ok {
				return ToFloat64(x) == val
			}
			return false
		}
	case float32:
		switch y := b.(type) {
		case uint:
			return float64(x) == float64(y)
		case uint8:
			return float64(x) == float64(y)
		case uint16:
			return float64(x) == float64(y)
		case uint32:
			return float64(x) == float64(y)
		case uint64:
			return float64(x) == float64(y)
		case int:
			return float64(x) == float64(y)
		case int8:
			return float64(x) == float64(y)
		case int16:
			return float64(x) == float64(y)
		case int32:
			return float64(x) == float64(y)
		case int64:
			return float64(x) == float64(y)
		case float32:
			return float64(x) == float64(y)
		case float64:
			return float64(x) == float64(y)
		case bool:
			return ToFloat64(x) == ToFloat64(y)
		case string:
			if val, ok := ToFloat64Safe(y); ok {
				return ToFloat64(x) == val
			}
			return false
		}
	case float64:
		switch y := b.(type) {
		case uint:
			return float64(x) == float64(y)
		case uint8:
			return float64(x) == float64(y)
		case uint16:
			return float64(x) == float64(y)
		case uint32:
			return float64(x) == float64(y)
		case uint64:
			return float64(x) == float64(y)
		case int:
			return float64(x) == float64(y)
		case int8:
			return float64(x) == float64(y)
		case int16:
			return float64(x) == float64(y)
		case int32:
			return float64(x) == float64(y)
		case int64:
			return float64(x) == float64(y)
		case float32:
			return float64(x) == float64(y)
		case float64:
			return float64(x) == float64(y)
		case bool:
			return ToFloat64(x) == ToFloat64(y)
		case string:
			if val, ok := ToFloat64Safe(y); ok {
				return ToFloat64(x) == val
			}
			return false
		}
	case string:
		switch y := b.(type) {
		case string:
			return x == y
		case uint, uint8, uint16, uint32, uint64, int, int8, int16, int32, int64, float32, float64:
			// Try numeric conversion
			if val, ok := ToFloat64Safe(x); ok {
				return val == ToFloat64(y)
			}
			return false
		case bool:
			// JavaScript-like coercion
			if x == "" {
				return y == false
			}
			if val, ok := ToFloat64Safe(x); ok {
				return val == ToFloat64(y)
			}
			return false
		}
	case bool:
		switch y := b.(type) {
		case bool:
			return x == y
		default:
			return false
		}
	}
	// Fallback to reflect.DeepEqual for complex types
	return reflect.DeepEqual(a, b)
}

func Less(a, b interface{}) bool {
	// Convert nil values to 0 for numeric comparisons
	if IsNil(a) {
		a = 0
	}
	if IsNil(b) {
		b = 0
	}

	// Handle numeric operations
	switch x := a.(type) {
	case uint:
		switch y := b.(type) {
		case uint:
			return int(x) < int(y)
		case uint8:
			return int(x) < int(y)
		case uint16:
			return int(x) < int(y)
		case uint32:
			return int(x) < int(y)
		case uint64:
			return int(x) < int(y)
		case int:
			return int(x) < int(y)
		case int8:
			return int(x) < int(y)
		case int16:
			return int(x) < int(y)
		case int32:
			return int(x) < int(y)
		case int64:
			return int(x) < int(y)
		case float32:
			return float64(x) < float64(y)
		case float64:
			return float64(x) < float64(y)
		case bool:
			return ToFloat64(x) < ToFloat64(y)
		case string:
			// Only allow comparison with numeric strings
			if val, ok := ToFloat64Safe(y); ok {
				return ToFloat64(x) < val
			}
			panic(fmt.Sprintf("invalid operation: %T < %T", x, y))
		}
	case uint8:
		switch y := b.(type) {
		case uint:
			return int(x) < int(y)
		case uint8:
			return int(x) < int(y)
		case uint16:
			return int(x) < int(y)
		case uint32:
			return int(x) < int(y)
		case uint64:
			return int(x) < int(y)
		case int:
			return int(x) < int(y)
		case int8:
			return int(x) < int(y)
		case int16:
			return int(x) < int(y)
		case int32:
			return int(x) < int(y)
		case int64:
			return int(x) < int(y)
		case float32:
			return float64(x) < float64(y)
		case float64:
			return float64(x) < float64(y)
		case bool:
			return ToFloat64(x) < ToFloat64(y)
		case string:
			if val, ok := ToFloat64Safe(y); ok {
				return ToFloat64(x) < val
			}
			panic(fmt.Sprintf("invalid operation: %T < %T", x, y))
		}
	case uint16:
		switch y := b.(type) {
		case uint:
			return int(x) < int(y)
		case uint8:
			return int(x) < int(y)
		case uint16:
			return int(x) < int(y)
		case uint32:
			return int(x) < int(y)
		case uint64:
			return int(x) < int(y)
		case int:
			return int(x) < int(y)
		case int8:
			return int(x) < int(y)
		case int16:
			return int(x) < int(y)
		case int32:
			return int(x) < int(y)
		case int64:
			return int(x) < int(y)
		case float32:
			return float64(x) < float64(y)
		case float64:
			return float64(x) < float64(y)
		case bool:
			return ToFloat64(x) < ToFloat64(y)
		case string:
			if val, ok := ToFloat64Safe(y); ok {
				return ToFloat64(x) < val
			}
			panic(fmt.Sprintf("invalid operation: %T < %T", x, y))
		}
	case uint32:
		switch y := b.(type) {
		case uint:
			return int(x) < int(y)
		case uint8:
			return int(x) < int(y)
		case uint16:
			return int(x) < int(y)
		case uint32:
			return int(x) < int(y)
		case uint64:
			return int(x) < int(y)
		case int:
			return int(x) < int(y)
		case int8:
			return int(x) < int(y)
		case int16:
			return int(x) < int(y)
		case int32:
			return int(x) < int(y)
		case int64:
			return int(x) < int(y)
		case float32:
			return float64(x) < float64(y)
		case float64:
			return float64(x) < float64(y)
		case bool:
			return ToFloat64(x) < ToFloat64(y)
		case string:
			if val, ok := ToFloat64Safe(y); ok {
				return ToFloat64(x) < val
			}
			panic(fmt.Sprintf("invalid operation: %T < %T", x, y))
		}
	case uint64:
		switch y := b.(type) {
		case uint:
			return int(x) < int(y)
		case uint8:
			return int(x) < int(y)
		case uint16:
			return int(x) < int(y)
		case uint32:
			return int(x) < int(y)
		case uint64:
			return int(x) < int(y)
		case int:
			return int(x) < int(y)
		case int8:
			return int(x) < int(y)
		case int16:
			return int(x) < int(y)
		case int32:
			return int(x) < int(y)
		case int64:
			return int(x) < int(y)
		case float32:
			return float64(x) < float64(y)
		case float64:
			return float64(x) < float64(y)
		case bool:
			return ToFloat64(x) < ToFloat64(y)
		case string:
			if val, ok := ToFloat64Safe(y); ok {
				return ToFloat64(x) < val
			}
			panic(fmt.Sprintf("invalid operation: %T < %T", x, y))
		}
	case int:
		switch y := b.(type) {
		case uint:
			return int(x) < int(y)
		case uint8:
			return int(x) < int(y)
		case uint16:
			return int(x) < int(y)
		case uint32:
			return int(x) < int(y)
		case uint64:
			return int(x) < int(y)
		case int:
			return int(x) < int(y)
		case int8:
			return int(x) < int(y)
		case int16:
			return int(x) < int(y)
		case int32:
			return int(x) < int(y)
		case int64:
			return int(x) < int(y)
		case float32:
			return float64(x) < float64(y)
		case float64:
			return float64(x) < float64(y)
		case bool:
			return ToFloat64(x) < ToFloat64(y)
		case string:
			if val, ok := ToFloat64Safe(y); ok {
				return ToFloat64(x) < val
			}
			panic(fmt.Sprintf("invalid operation: %T < %T", x, y))
		}
	case int8:
		switch y := b.(type) {
		case uint:
			return int(x) < int(y)
		case uint8:
			return int(x) < int(y)
		case uint16:
			return int(x) < int(y)
		case uint32:
			return int(x) < int(y)
		case uint64:
			return int(x) < int(y)
		case int:
			return int(x) < int(y)
		case int8:
			return int(x) < int(y)
		case int16:
			return int(x) < int(y)
		case int32:
			return int(x) < int(y)
		case int64:
			return int(x) < int(y)
		case float32:
			return float64(x) < float64(y)
		case float64:
			return float64(x) < float64(y)
		case bool:
			return ToFloat64(x) < ToFloat64(y)
		case string:
			if val, ok := ToFloat64Safe(y); ok {
				return ToFloat64(x) < val
			}
			panic(fmt.Sprintf("invalid operation: %T < %T", x, y))
		}
	case int16:
		switch y := b.(type) {
		case uint:
			return int(x) < int(y)
		case uint8:
			return int(x) < int(y)
		case uint16:
			return int(x) < int(y)
		case uint32:
			return int(x) < int(y)
		case uint64:
			return int(x) < int(y)
		case int:
			return int(x) < int(y)
		case int8:
			return int(x) < int(y)
		case int16:
			return int(x) < int(y)
		case int32:
			return int(x) < int(y)
		case int64:
			return int(x) < int(y)
		case float32:
			return float64(x) < float64(y)
		case float64:
			return float64(x) < float64(y)
		case bool:
			return ToFloat64(x) < ToFloat64(y)
		case string:
			if val, ok := ToFloat64Safe(y); ok {
				return ToFloat64(x) < val
			}
			panic(fmt.Sprintf("invalid operation: %T < %T", x, y))
		}
	case int32:
		switch y := b.(type) {
		case uint:
			return int(x) < int(y)
		case uint8:
			return int(x) < int(y)
		case uint16:
			return int(x) < int(y)
		case uint32:
			return int(x) < int(y)
		case uint64:
			return int(x) < int(y)
		case int:
			return int(x) < int(y)
		case int8:
			return int(x) < int(y)
		case int16:
			return int(x) < int(y)
		case int32:
			return int(x) < int(y)
		case int64:
			return int(x) < int(y)
		case float32:
			return float64(x) < float64(y)
		case float64:
			return float64(x) < float64(y)
		case bool:
			return ToFloat64(x) < ToFloat64(y)
		case string:
			if val, ok := ToFloat64Safe(y); ok {
				return ToFloat64(x) < val
			}
			panic(fmt.Sprintf("invalid operation: %T < %T", x, y))
		}
	case int64:
		switch y := b.(type) {
		case uint:
			return int(x) < int(y)
		case uint8:
			return int(x) < int(y)
		case uint16:
			return int(x) < int(y)
		case uint32:
			return int(x) < int(y)
		case uint64:
			return int(x) < int(y)
		case int:
			return int(x) < int(y)
		case int8:
			return int(x) < int(y)
		case int16:
			return int(x) < int(y)
		case int32:
			return int(x) < int(y)
		case int64:
			return int(x) < int(y)
		case float32:
			return float64(x) < float64(y)
		case float64:
			return float64(x) < float64(y)
		case bool:
			return ToFloat64(x) < ToFloat64(y)
		case string:
			if val, ok := ToFloat64Safe(y); ok {
				return ToFloat64(x) < val
			}
			panic(fmt.Sprintf("invalid operation: %T < %T", x, y))
		}
	case float32:
		switch y := b.(type) {
		case uint:
			return float64(x) < float64(y)
		case uint8:
			return float64(x) < float64(y)
		case uint16:
			return float64(x) < float64(y)
		case uint32:
			return float64(x) < float64(y)
		case uint64:
			return float64(x) < float64(y)
		case int:
			return float64(x) < float64(y)
		case int8:
			return float64(x) < float64(y)
		case int16:
			return float64(x) < float64(y)
		case int32:
			return float64(x) < float64(y)
		case int64:
			return float64(x) < float64(y)
		case float32:
			return float64(x) < float64(y)
		case float64:
			return float64(x) < float64(y)
		case bool:
			return ToFloat64(x) < ToFloat64(y)
		case string:
			if val, ok := ToFloat64Safe(y); ok {
				return ToFloat64(x) < val
			}
			panic(fmt.Sprintf("invalid operation: %T < %T", x, y))
		}
	case float64:
		switch y := b.(type) {
		case uint:
			return float64(x) < float64(y)
		case uint8:
			return float64(x) < float64(y)
		case uint16:
			return float64(x) < float64(y)
		case uint32:
			return float64(x) < float64(y)
		case uint64:
			return float64(x) < float64(y)
		case int:
			return float64(x) < float64(y)
		case int8:
			return float64(x) < float64(y)
		case int16:
			return float64(x) < float64(y)
		case int32:
			return float64(x) < float64(y)
		case int64:
			return float64(x) < float64(y)
		case float32:
			return float64(x) < float64(y)
		case float64:
			return float64(x) < float64(y)
		case bool:
			return ToFloat64(x) < ToFloat64(y)
		case string:
			if val, ok := ToFloat64Safe(y); ok {
				return ToFloat64(x) < val
			}
			panic(fmt.Sprintf("invalid operation: %T < %T", x, y))
		}
	case bool:
		switch y := b.(type) {
		case uint:
			return ToFloat64(x) < ToFloat64(y)
		case uint8:
			return ToFloat64(x) < ToFloat64(y)
		case uint16:
			return ToFloat64(x) < ToFloat64(y)
		case uint32:
			return ToFloat64(x) < ToFloat64(y)
		case uint64:
			return ToFloat64(x) < ToFloat64(y)
		case int:
			return ToFloat64(x) < ToFloat64(y)
		case int8:
			return ToFloat64(x) < ToFloat64(y)
		case int16:
			return ToFloat64(x) < ToFloat64(y)
		case int32:
			return ToFloat64(x) < ToFloat64(y)
		case int64:
			return ToFloat64(x) < ToFloat64(y)
		case float32:
			return ToFloat64(x) < ToFloat64(y)
		case float64:
			return ToFloat64(x) < ToFloat64(y)
		case bool:
			return ToFloat64(x) < ToFloat64(y)
		case string:
			if val, ok := ToFloat64Safe(y); ok {
				return ToFloat64(x) < val
			}
			panic(fmt.Sprintf("invalid operation: %T < %T", x, y))
		}
	case string:
		switch y := b.(type) {
		case string:
			return x < y
		case uint, uint8, uint16, uint32, uint64, int, int8, int16, int32, int64, float32, float64:
			if val, ok := ToFloat64Safe(x); ok {
				return val < ToFloat64(y)
			}
			panic(fmt.Sprintf("invalid operation: %T < %T", x, y))
		case bool:
			if val, ok := ToFloat64Safe(x); ok {
				return val < ToFloat64(y)
			}
			panic(fmt.Sprintf("invalid operation: %T < %T", x, y))
		}
	case time.Time:
		switch y := b.(type) {
		case time.Time:
			return x.Before(y)
		}
	case time.Duration:
		switch y := b.(type) {
		case time.Duration:
			return x < y
		}
	}
	panic(fmt.Sprintf("invalid operation: %T < %T", a, b))
}

func More(a, b interface{}) bool {
	// Convert nil values to 0 for numeric comparisons
	if IsNil(a) {
		a = 0
	}
	if IsNil(b) {
		b = 0
	}
	// Handle numeric operations
	switch x := a.(type) {
	case uint:
		switch y := b.(type) {
		case uint:
			return int(x) > int(y)
		case uint8:
			return int(x) > int(y)
		case uint16:
			return int(x) > int(y)
		case uint32:
			return int(x) > int(y)
		case uint64:
			return int(x) > int(y)
		case int:
			return int(x) > int(y)
		case int8:
			return int(x) > int(y)
		case int16:
			return int(x) > int(y)
		case int32:
			return int(x) > int(y)
		case int64:
			return int(x) > int(y)
		case float32:
			return float64(x) > float64(y)
		case float64:
			return float64(x) > float64(y)
		case bool:
			return ToFloat64(x) > ToFloat64(y)
		case string:
			return ToFloat64(x) > ToFloat64(y)
		}
	case uint8:
		switch y := b.(type) {
		case uint:
			return int(x) > int(y)
		case uint8:
			return int(x) > int(y)
		case uint16:
			return int(x) > int(y)
		case uint32:
			return int(x) > int(y)
		case uint64:
			return int(x) > int(y)
		case int:
			return int(x) > int(y)
		case int8:
			return int(x) > int(y)
		case int16:
			return int(x) > int(y)
		case int32:
			return int(x) > int(y)
		case int64:
			return int(x) > int(y)
		case float32:
			return float64(x) > float64(y)
		case float64:
			return float64(x) > float64(y)
		case bool:
			return ToFloat64(x) > ToFloat64(y)
		case string:
			return ToFloat64(x) > ToFloat64(y)
		}
	case uint16:
		switch y := b.(type) {
		case uint:
			return int(x) > int(y)
		case uint8:
			return int(x) > int(y)
		case uint16:
			return int(x) > int(y)
		case uint32:
			return int(x) > int(y)
		case uint64:
			return int(x) > int(y)
		case int:
			return int(x) > int(y)
		case int8:
			return int(x) > int(y)
		case int16:
			return int(x) > int(y)
		case int32:
			return int(x) > int(y)
		case int64:
			return int(x) > int(y)
		case float32:
			return float64(x) > float64(y)
		case float64:
			return float64(x) > float64(y)
		case bool:
			return ToFloat64(x) > ToFloat64(y)
		case string:
			return ToFloat64(x) > ToFloat64(y)
		}
	case uint32:
		switch y := b.(type) {
		case uint:
			return int(x) > int(y)
		case uint8:
			return int(x) > int(y)
		case uint16:
			return int(x) > int(y)
		case uint32:
			return int(x) > int(y)
		case uint64:
			return int(x) > int(y)
		case int:
			return int(x) > int(y)
		case int8:
			return int(x) > int(y)
		case int16:
			return int(x) > int(y)
		case int32:
			return int(x) > int(y)
		case int64:
			return int(x) > int(y)
		case float32:
			return float64(x) > float64(y)
		case float64:
			return float64(x) > float64(y)
		case bool:
			return ToFloat64(x) > ToFloat64(y)
		case string:
			return ToFloat64(x) > ToFloat64(y)
		}
	case uint64:
		switch y := b.(type) {
		case uint:
			return int(x) > int(y)
		case uint8:
			return int(x) > int(y)
		case uint16:
			return int(x) > int(y)
		case uint32:
			return int(x) > int(y)
		case uint64:
			return int(x) > int(y)
		case int:
			return int(x) > int(y)
		case int8:
			return int(x) > int(y)
		case int16:
			return int(x) > int(y)
		case int32:
			return int(x) > int(y)
		case int64:
			return int(x) > int(y)
		case float32:
			return float64(x) > float64(y)
		case float64:
			return float64(x) > float64(y)
		case bool:
			return ToFloat64(x) > ToFloat64(y)
		case string:
			return ToFloat64(x) > ToFloat64(y)
		}
	case int:
		switch y := b.(type) {
		case uint:
			return int(x) > int(y)
		case uint8:
			return int(x) > int(y)
		case uint16:
			return int(x) > int(y)
		case uint32:
			return int(x) > int(y)
		case uint64:
			return int(x) > int(y)
		case int:
			return int(x) > int(y)
		case int8:
			return int(x) > int(y)
		case int16:
			return int(x) > int(y)
		case int32:
			return int(x) > int(y)
		case int64:
			return int(x) > int(y)
		case float32:
			return float64(x) > float64(y)
		case float64:
			return float64(x) > float64(y)
		case bool:
			return ToFloat64(x) > ToFloat64(y)
		case string:
			return ToFloat64(x) > ToFloat64(y)
		}
	case int8:
		switch y := b.(type) {
		case uint:
			return int(x) > int(y)
		case uint8:
			return int(x) > int(y)
		case uint16:
			return int(x) > int(y)
		case uint32:
			return int(x) > int(y)
		case uint64:
			return int(x) > int(y)
		case int:
			return int(x) > int(y)
		case int8:
			return int(x) > int(y)
		case int16:
			return int(x) > int(y)
		case int32:
			return int(x) > int(y)
		case int64:
			return int(x) > int(y)
		case float32:
			return float64(x) > float64(y)
		case float64:
			return float64(x) > float64(y)
		case bool:
			return ToFloat64(x) > ToFloat64(y)
		case string:
			return ToFloat64(x) > ToFloat64(y)
		}
	case int16:
		switch y := b.(type) {
		case uint:
			return int(x) > int(y)
		case uint8:
			return int(x) > int(y)
		case uint16:
			return int(x) > int(y)
		case uint32:
			return int(x) > int(y)
		case uint64:
			return int(x) > int(y)
		case int:
			return int(x) > int(y)
		case int8:
			return int(x) > int(y)
		case int16:
			return int(x) > int(y)
		case int32:
			return int(x) > int(y)
		case int64:
			return int(x) > int(y)
		case float32:
			return float64(x) > float64(y)
		case float64:
			return float64(x) > float64(y)
		case bool:
			return ToFloat64(x) > ToFloat64(y)
		case string:
			return ToFloat64(x) > ToFloat64(y)
		}
	case int32:
		switch y := b.(type) {
		case uint:
			return int(x) > int(y)
		case uint8:
			return int(x) > int(y)
		case uint16:
			return int(x) > int(y)
		case uint32:
			return int(x) > int(y)
		case uint64:
			return int(x) > int(y)
		case int:
			return int(x) > int(y)
		case int8:
			return int(x) > int(y)
		case int16:
			return int(x) > int(y)
		case int32:
			return int(x) > int(y)
		case int64:
			return int(x) > int(y)
		case float32:
			return float64(x) > float64(y)
		case float64:
			return float64(x) > float64(y)
		case bool:
			return ToFloat64(x) > ToFloat64(y)
		case string:
			return ToFloat64(x) > ToFloat64(y)
		}
	case int64:
		switch y := b.(type) {
		case uint:
			return int(x) > int(y)
		case uint8:
			return int(x) > int(y)
		case uint16:
			return int(x) > int(y)
		case uint32:
			return int(x) > int(y)
		case uint64:
			return int(x) > int(y)
		case int:
			return int(x) > int(y)
		case int8:
			return int(x) > int(y)
		case int16:
			return int(x) > int(y)
		case int32:
			return int(x) > int(y)
		case int64:
			return int(x) > int(y)
		case float32:
			return float64(x) > float64(y)
		case float64:
			return float64(x) > float64(y)
		case bool:
			return ToFloat64(x) > ToFloat64(y)
		case string:
			return ToFloat64(x) > ToFloat64(y)
		}
	case float32:
		switch y := b.(type) {
		case uint:
			return float64(x) > float64(y)
		case uint8:
			return float64(x) > float64(y)
		case uint16:
			return float64(x) > float64(y)
		case uint32:
			return float64(x) > float64(y)
		case uint64:
			return float64(x) > float64(y)
		case int:
			return float64(x) > float64(y)
		case int8:
			return float64(x) > float64(y)
		case int16:
			return float64(x) > float64(y)
		case int32:
			return float64(x) > float64(y)
		case int64:
			return float64(x) > float64(y)
		case float32:
			return float64(x) > float64(y)
		case float64:
			return float64(x) > float64(y)
		case bool:
			return ToFloat64(x) > ToFloat64(y)
		case string:
			return ToFloat64(x) > ToFloat64(y)
		}
	case float64:
		switch y := b.(type) {
		case uint:
			return float64(x) > float64(y)
		case uint8:
			return float64(x) > float64(y)
		case uint16:
			return float64(x) > float64(y)
		case uint32:
			return float64(x) > float64(y)
		case uint64:
			return float64(x) > float64(y)
		case int:
			return float64(x) > float64(y)
		case int8:
			return float64(x) > float64(y)
		case int16:
			return float64(x) > float64(y)
		case int32:
			return float64(x) > float64(y)
		case int64:
			return float64(x) > float64(y)
		case float32:
			return float64(x) > float64(y)
		case float64:
			return float64(x) > float64(y)
		case bool:
			return ToFloat64(x) > ToFloat64(y)
		case string:
			return ToFloat64(x) > ToFloat64(y)
		}
	case bool:
		switch y := b.(type) {
		case uint:
			return ToFloat64(x) > ToFloat64(y)
		case uint8:
			return ToFloat64(x) > ToFloat64(y)
		case uint16:
			return ToFloat64(x) > ToFloat64(y)
		case uint32:
			return ToFloat64(x) > ToFloat64(y)
		case uint64:
			return ToFloat64(x) > ToFloat64(y)
		case int:
			return ToFloat64(x) > ToFloat64(y)
		case int8:
			return ToFloat64(x) > ToFloat64(y)
		case int16:
			return ToFloat64(x) > ToFloat64(y)
		case int32:
			return ToFloat64(x) > ToFloat64(y)
		case int64:
			return ToFloat64(x) > ToFloat64(y)
		case float32:
			return ToFloat64(x) > ToFloat64(y)
		case float64:
			return ToFloat64(x) > ToFloat64(y)
		case bool:
			return ToFloat64(x) > ToFloat64(y)
		case string:
			if val, ok := ToFloat64Safe(y); ok {
				return ToFloat64(x) > val
			}
			panic(fmt.Sprintf("invalid operation: %T > %T", x, y))
		}
	case string:
		switch y := b.(type) {
		case string:
			return x > y
		case uint, uint8, uint16, uint32, uint64, int, int8, int16, int32, int64, float32, float64:
			if val, ok := ToFloat64Safe(x); ok {
				return val > ToFloat64(y)
			}
			panic(fmt.Sprintf("invalid operation: %T > %T", x, y))
		case bool:
			if val, ok := ToFloat64Safe(x); ok {
				return val > ToFloat64(y)
			}
			panic(fmt.Sprintf("invalid operation: %T > %T", x, y))
		}
	case time.Time:
		switch y := b.(type) {
		case time.Time:
			return x.After(y)
		}
	case time.Duration:
		switch y := b.(type) {
		case time.Duration:
			return x > y
		}

	}
	panic(fmt.Sprintf("invalid operation: %T > %T", a, b))
}

func LessOrEqual(a, b interface{}) bool {
	// Convert nil values to 0 for numeric comparisons
	if IsNil(a) {
		a = 0
	}
	if IsNil(b) {
		b = 0
	}

	// Handle numeric operations for non-nil values
	switch x := a.(type) {
	case uint:
		switch y := b.(type) {
		case uint:
			return int(x) <= int(y)
		case uint8:
			return int(x) <= int(y)
		case uint16:
			return int(x) <= int(y)
		case uint32:
			return int(x) <= int(y)
		case uint64:
			return int(x) <= int(y)
		case int:
			return int(x) <= int(y)
		case int8:
			return int(x) <= int(y)
		case int16:
			return int(x) <= int(y)
		case int32:
			return int(x) <= int(y)
		case int64:
			return int(x) <= int(y)
		case float32:
			return float64(x) <= float64(y)
		case float64:
			return float64(x) <= float64(y)
		case bool:
			return float64(x) <= ToFloat64(y)
		case string:
			return float64(x) <= ToFloat64(y)
		}
	case uint8:
		switch y := b.(type) {
		case uint:
			return int(x) <= int(y)
		case uint8:
			return int(x) <= int(y)
		case uint16:
			return int(x) <= int(y)
		case uint32:
			return int(x) <= int(y)
		case uint64:
			return int(x) <= int(y)
		case int:
			return int(x) <= int(y)
		case int8:
			return int(x) <= int(y)
		case int16:
			return int(x) <= int(y)
		case int32:
			return int(x) <= int(y)
		case int64:
			return int(x) <= int(y)
		case float32:
			return float64(x) <= float64(y)
		case float64:
			return float64(x) <= float64(y)
		case bool:
			return float64(x) <= ToFloat64(y)
		case string:
			return float64(x) <= ToFloat64(y)
		}
	case uint16:
		switch y := b.(type) {
		case uint:
			return int(x) <= int(y)
		case uint8:
			return int(x) <= int(y)
		case uint16:
			return int(x) <= int(y)
		case uint32:
			return int(x) <= int(y)
		case uint64:
			return int(x) <= int(y)
		case int:
			return int(x) <= int(y)
		case int8:
			return int(x) <= int(y)
		case int16:
			return int(x) <= int(y)
		case int32:
			return int(x) <= int(y)
		case int64:
			return int(x) <= int(y)
		case float32:
			return float64(x) <= float64(y)
		case float64:
			return float64(x) <= float64(y)
		case bool:
			return float64(x) <= ToFloat64(y)
		case string:
			return float64(x) <= ToFloat64(y)
		}
	case uint32:
		switch y := b.(type) {
		case uint:
			return int(x) <= int(y)
		case uint8:
			return int(x) <= int(y)
		case uint16:
			return int(x) <= int(y)
		case uint32:
			return int(x) <= int(y)
		case uint64:
			return int(x) <= int(y)
		case int:
			return int(x) <= int(y)
		case int8:
			return int(x) <= int(y)
		case int16:
			return int(x) <= int(y)
		case int32:
			return int(x) <= int(y)
		case int64:
			return int(x) <= int(y)
		case float32:
			return float64(x) <= float64(y)
		case float64:
			return float64(x) <= float64(y)
		case bool:
			return float64(x) <= ToFloat64(y)
		case string:
			return float64(x) <= ToFloat64(y)
		}
	case uint64:
		switch y := b.(type) {
		case uint:
			return int(x) <= int(y)
		case uint8:
			return int(x) <= int(y)
		case uint16:
			return int(x) <= int(y)
		case uint32:
			return int(x) <= int(y)
		case uint64:
			return int(x) <= int(y)
		case int:
			return int(x) <= int(y)
		case int8:
			return int(x) <= int(y)
		case int16:
			return int(x) <= int(y)
		case int32:
			return int(x) <= int(y)
		case int64:
			return int(x) <= int(y)
		case float32:
			return float64(x) <= float64(y)
		case float64:
			return float64(x) <= float64(y)
		case bool:
			return float64(x) <= ToFloat64(y)
		case string:
			return float64(x) <= ToFloat64(y)
		}
	case int:
		switch y := b.(type) {
		case uint:
			return int(x) <= int(y)
		case uint8:
			return int(x) <= int(y)
		case uint16:
			return int(x) <= int(y)
		case uint32:
			return int(x) <= int(y)
		case uint64:
			return int(x) <= int(y)
		case int:
			return int(x) <= int(y)
		case int8:
			return int(x) <= int(y)
		case int16:
			return int(x) <= int(y)
		case int32:
			return int(x) <= int(y)
		case int64:
			return int(x) <= int(y)
		case float32:
			return float64(x) <= float64(y)
		case float64:
			return float64(x) <= float64(y)
		case bool:
			return float64(x) <= ToFloat64(y)
		case string:
			return float64(x) <= ToFloat64(y)
		}
	case int8:
		switch y := b.(type) {
		case uint:
			return int(x) <= int(y)
		case uint8:
			return int(x) <= int(y)
		case uint16:
			return int(x) <= int(y)
		case uint32:
			return int(x) <= int(y)
		case uint64:
			return int(x) <= int(y)
		case int:
			return int(x) <= int(y)
		case int8:
			return int(x) <= int(y)
		case int16:
			return int(x) <= int(y)
		case int32:
			return int(x) <= int(y)
		case int64:
			return int(x) <= int(y)
		case float32:
			return float64(x) <= float64(y)
		case float64:
			return float64(x) <= float64(y)
		case bool:
			return float64(x) <= ToFloat64(y)
		case string:
			return float64(x) <= ToFloat64(y)
		}
	case int16:
		switch y := b.(type) {
		case uint:
			return int(x) <= int(y)
		case uint8:
			return int(x) <= int(y)
		case uint16:
			return int(x) <= int(y)
		case uint32:
			return int(x) <= int(y)
		case uint64:
			return int(x) <= int(y)
		case int:
			return int(x) <= int(y)
		case int8:
			return int(x) <= int(y)
		case int16:
			return int(x) <= int(y)
		case int32:
			return int(x) <= int(y)
		case int64:
			return int(x) <= int(y)
		case float32:
			return float64(x) <= float64(y)
		case float64:
			return float64(x) <= float64(y)
		case bool:
			return float64(x) <= ToFloat64(y)
		case string:
			return float64(x) <= ToFloat64(y)
		}
	case int32:
		switch y := b.(type) {
		case uint:
			return int(x) <= int(y)
		case uint8:
			return int(x) <= int(y)
		case uint16:
			return int(x) <= int(y)
		case uint32:
			return int(x) <= int(y)
		case uint64:
			return int(x) <= int(y)
		case int:
			return int(x) <= int(y)
		case int8:
			return int(x) <= int(y)
		case int16:
			return int(x) <= int(y)
		case int32:
			return int(x) <= int(y)
		case int64:
			return int(x) <= int(y)
		case float32:
			return float64(x) <= float64(y)
		case float64:
			return float64(x) <= float64(y)
		case bool:
			return float64(x) <= ToFloat64(y)
		case string:
			return float64(x) <= ToFloat64(y)
		}
	case int64:
		switch y := b.(type) {
		case uint:
			return int(x) <= int(y)
		case uint8:
			return int(x) <= int(y)
		case uint16:
			return int(x) <= int(y)
		case uint32:
			return int(x) <= int(y)
		case uint64:
			return int(x) <= int(y)
		case int:
			return int(x) <= int(y)
		case int8:
			return int(x) <= int(y)
		case int16:
			return int(x) <= int(y)
		case int32:
			return int(x) <= int(y)
		case int64:
			return int(x) <= int(y)
		case float32:
			return float64(x) <= float64(y)
		case float64:
			return float64(x) <= float64(y)
		case bool:
			return float64(x) <= ToFloat64(y)
		case string:
			return float64(x) <= ToFloat64(y)
		}
	case float32:
		switch y := b.(type) {
		case uint:
			return float64(x) <= float64(y)
		case uint8:
			return float64(x) <= float64(y)
		case uint16:
			return float64(x) <= float64(y)
		case uint32:
			return float64(x) <= float64(y)
		case uint64:
			return float64(x) <= float64(y)
		case int:
			return float64(x) <= float64(y)
		case int8:
			return float64(x) <= float64(y)
		case int16:
			return float64(x) <= float64(y)
		case int32:
			return float64(x) <= float64(y)
		case int64:
			return float64(x) <= float64(y)
		case float32:
			return float64(x) <= float64(y)
		case float64:
			return float64(x) <= float64(y)
		case bool:
			return float64(x) <= ToFloat64(y)
		case string:
			return float64(x) <= ToFloat64(y)
		}
	case float64:
		switch y := b.(type) {
		case uint:
			return float64(x) <= float64(y)
		case uint8:
			return float64(x) <= float64(y)
		case uint16:
			return float64(x) <= float64(y)
		case uint32:
			return float64(x) <= float64(y)
		case uint64:
			return float64(x) <= float64(y)
		case int:
			return float64(x) <= float64(y)
		case int8:
			return float64(x) <= float64(y)
		case int16:
			return float64(x) <= float64(y)
		case int32:
			return float64(x) <= float64(y)
		case int64:
			return float64(x) <= float64(y)
		case float32:
			return float64(x) <= float64(y)
		case float64:
			return float64(x) <= float64(y)
		case bool:
			return float64(x) <= ToFloat64(y)
		case string:
			return float64(x) <= ToFloat64(y)
		}
	case bool:
		switch y := b.(type) {
		case uint:
			return ToFloat64(x) <= float64(y)
		case uint8:
			return ToFloat64(x) <= float64(y)
		case uint16:
			return ToFloat64(x) <= float64(y)
		case uint32:
			return ToFloat64(x) <= float64(y)
		case uint64:
			return ToFloat64(x) <= float64(y)
		case int:
			return ToFloat64(x) <= float64(y)
		case int8:
			return ToFloat64(x) <= float64(y)
		case int16:
			return ToFloat64(x) <= float64(y)
		case int32:
			return ToFloat64(x) <= float64(y)
		case int64:
			return ToFloat64(x) <= float64(y)
		case float32:
			return ToFloat64(x) <= float64(y)
		case float64:
			return ToFloat64(x) <= float64(y)
		case bool:
			return ToFloat64(x) <= ToFloat64(y)
		case string:
			return ToFloat64(x) <= ToFloat64(y)
		}
	case string:
		switch y := b.(type) {
		case uint:
			return ToFloat64(x) <= float64(y)
		case uint8:
			return ToFloat64(x) <= float64(y)
		case uint16:
			return ToFloat64(x) <= float64(y)
		case uint32:
			return ToFloat64(x) <= float64(y)
		case uint64:
			return ToFloat64(x) <= float64(y)
		case int:
			return ToFloat64(x) <= float64(y)
		case int8:
			return ToFloat64(x) <= float64(y)
		case int16:
			return ToFloat64(x) <= float64(y)
		case int32:
			return ToFloat64(x) <= float64(y)
		case int64:
			return ToFloat64(x) <= float64(y)
		case float32:
			return ToFloat64(x) <= float64(y)
		case float64:
			return ToFloat64(x) <= float64(y)
		case bool:
			return ToFloat64(x) <= ToFloat64(y)
		case string:
			return x <= y
		}
	case time.Time:
		switch y := b.(type) {
		case time.Time:
			return x.Before(y) || x.Equal(y)
		}
	case time.Duration:
		switch y := b.(type) {
		case time.Duration:
			return x <= y
		}
	}
	panic(fmt.Sprintf("invalid operation: %T <= %T", a, b))
}

func MoreOrEqual(a, b interface{}) bool {
	// Convert nil values to 0 for numeric comparisons
	if IsNil(a) {
		a = 0
	}
	if IsNil(b) {
		b = 0
	}

	// Handle numeric operations for non-nil values
	switch x := a.(type) {
	case uint:
		switch y := b.(type) {
		case uint:
			return int(x) >= int(y)
		case uint8:
			return int(x) >= int(y)
		case uint16:
			return int(x) >= int(y)
		case uint32:
			return int(x) >= int(y)
		case uint64:
			return int(x) >= int(y)
		case int:
			return int(x) >= int(y)
		case int8:
			return int(x) >= int(y)
		case int16:
			return int(x) >= int(y)
		case int32:
			return int(x) >= int(y)
		case int64:
			return int(x) >= int(y)
		case float32:
			return float64(x) >= float64(y)
		case float64:
			return float64(x) >= float64(y)
		case bool:
			return float64(x) >= ToFloat64(y)
		case string:
			return float64(x) >= ToFloat64(y)
		}
	case uint8:
		switch y := b.(type) {
		case uint:
			return int(x) >= int(y)
		case uint8:
			return int(x) >= int(y)
		case uint16:
			return int(x) >= int(y)
		case uint32:
			return int(x) >= int(y)
		case uint64:
			return int(x) >= int(y)
		case int:
			return int(x) >= int(y)
		case int8:
			return int(x) >= int(y)
		case int16:
			return int(x) >= int(y)
		case int32:
			return int(x) >= int(y)
		case int64:
			return int(x) >= int(y)
		case float32:
			return float64(x) >= float64(y)
		case float64:
			return float64(x) >= float64(y)
		case bool:
			return float64(x) >= ToFloat64(y)
		case string:
			return float64(x) >= ToFloat64(y)
		}
	case uint16:
		switch y := b.(type) {
		case uint:
			return int(x) >= int(y)
		case uint8:
			return int(x) >= int(y)
		case uint16:
			return int(x) >= int(y)
		case uint32:
			return int(x) >= int(y)
		case uint64:
			return int(x) >= int(y)
		case int:
			return int(x) >= int(y)
		case int8:
			return int(x) >= int(y)
		case int16:
			return int(x) >= int(y)
		case int32:
			return int(x) >= int(y)
		case int64:
			return int(x) >= int(y)
		case float32:
			return float64(x) >= float64(y)
		case float64:
			return float64(x) >= float64(y)
		case bool:
			return float64(x) >= ToFloat64(y)
		case string:
			return float64(x) >= ToFloat64(y)
		}
	case uint32:
		switch y := b.(type) {
		case uint:
			return int(x) >= int(y)
		case uint8:
			return int(x) >= int(y)
		case uint16:
			return int(x) >= int(y)
		case uint32:
			return int(x) >= int(y)
		case uint64:
			return int(x) >= int(y)
		case int:
			return int(x) >= int(y)
		case int8:
			return int(x) >= int(y)
		case int16:
			return int(x) >= int(y)
		case int32:
			return int(x) >= int(y)
		case int64:
			return int(x) >= int(y)
		case float32:
			return float64(x) >= float64(y)
		case float64:
			return float64(x) >= float64(y)
		case bool:
			return float64(x) >= ToFloat64(y)
		case string:
			return float64(x) >= ToFloat64(y)
		}
	case uint64:
		switch y := b.(type) {
		case uint:
			return int(x) >= int(y)
		case uint8:
			return int(x) >= int(y)
		case uint16:
			return int(x) >= int(y)
		case uint32:
			return int(x) >= int(y)
		case uint64:
			return int(x) >= int(y)
		case int:
			return int(x) >= int(y)
		case int8:
			return int(x) >= int(y)
		case int16:
			return int(x) >= int(y)
		case int32:
			return int(x) >= int(y)
		case int64:
			return int(x) >= int(y)
		case float32:
			return float64(x) >= float64(y)
		case float64:
			return float64(x) >= float64(y)
		case bool:
			return float64(x) >= ToFloat64(y)
		case string:
			return float64(x) >= ToFloat64(y)
		}
	case int:
		switch y := b.(type) {
		case uint:
			return int(x) >= int(y)
		case uint8:
			return int(x) >= int(y)
		case uint16:
			return int(x) >= int(y)
		case uint32:
			return int(x) >= int(y)
		case uint64:
			return int(x) >= int(y)
		case int:
			return int(x) >= int(y)
		case int8:
			return int(x) >= int(y)
		case int16:
			return int(x) >= int(y)
		case int32:
			return int(x) >= int(y)
		case int64:
			return int(x) >= int(y)
		case float32:
			return float64(x) >= float64(y)
		case float64:
			return float64(x) >= float64(y)
		case bool:
			return float64(x) >= ToFloat64(y)
		case string:
			return float64(x) >= ToFloat64(y)
		}
	case int8:
		switch y := b.(type) {
		case uint:
			return int(x) >= int(y)
		case uint8:
			return int(x) >= int(y)
		case uint16:
			return int(x) >= int(y)
		case uint32:
			return int(x) >= int(y)
		case uint64:
			return int(x) >= int(y)
		case int:
			return int(x) >= int(y)
		case int8:
			return int(x) >= int(y)
		case int16:
			return int(x) >= int(y)
		case int32:
			return int(x) >= int(y)
		case int64:
			return int(x) >= int(y)
		case float32:
			return float64(x) >= float64(y)
		case float64:
			return float64(x) >= float64(y)
		case bool:
			return float64(x) >= ToFloat64(y)
		case string:
			return float64(x) >= ToFloat64(y)
		}
	case int16:
		switch y := b.(type) {
		case uint:
			return int(x) >= int(y)
		case uint8:
			return int(x) >= int(y)
		case uint16:
			return int(x) >= int(y)
		case uint32:
			return int(x) >= int(y)
		case uint64:
			return int(x) >= int(y)
		case int:
			return int(x) >= int(y)
		case int8:
			return int(x) >= int(y)
		case int16:
			return int(x) >= int(y)
		case int32:
			return int(x) >= int(y)
		case int64:
			return int(x) >= int(y)
		case float32:
			return float64(x) >= float64(y)
		case float64:
			return float64(x) >= float64(y)
		case bool:
			return float64(x) >= ToFloat64(y)
		case string:
			return float64(x) >= ToFloat64(y)
		}
	case int32:
		switch y := b.(type) {
		case uint:
			return int(x) >= int(y)
		case uint8:
			return int(x) >= int(y)
		case uint16:
			return int(x) >= int(y)
		case uint32:
			return int(x) >= int(y)
		case uint64:
			return int(x) >= int(y)
		case int:
			return int(x) >= int(y)
		case int8:
			return int(x) >= int(y)
		case int16:
			return int(x) >= int(y)
		case int32:
			return int(x) >= int(y)
		case int64:
			return int(x) >= int(y)
		case float32:
			return float64(x) >= float64(y)
		case float64:
			return float64(x) >= float64(y)
		case bool:
			return float64(x) >= ToFloat64(y)
		case string:
			return float64(x) >= ToFloat64(y)
		}
	case int64:
		switch y := b.(type) {
		case uint:
			return int(x) >= int(y)
		case uint8:
			return int(x) >= int(y)
		case uint16:
			return int(x) >= int(y)
		case uint32:
			return int(x) >= int(y)
		case uint64:
			return int(x) >= int(y)
		case int:
			return int(x) >= int(y)
		case int8:
			return int(x) >= int(y)
		case int16:
			return int(x) >= int(y)
		case int32:
			return int(x) >= int(y)
		case int64:
			return int(x) >= int(y)
		case float32:
			return float64(x) >= float64(y)
		case float64:
			return float64(x) >= float64(y)
		case bool:
			return float64(x) >= ToFloat64(y)
		case string:
			return float64(x) >= ToFloat64(y)
		}
	case float32:
		switch y := b.(type) {
		case uint:
			return float64(x) >= float64(y)
		case uint8:
			return float64(x) >= float64(y)
		case uint16:
			return float64(x) >= float64(y)
		case uint32:
			return float64(x) >= float64(y)
		case uint64:
			return float64(x) >= float64(y)
		case int:
			return float64(x) >= float64(y)
		case int8:
			return float64(x) >= float64(y)
		case int16:
			return float64(x) >= float64(y)
		case int32:
			return float64(x) >= float64(y)
		case int64:
			return float64(x) >= float64(y)
		case float32:
			return float64(x) >= float64(y)
		case float64:
			return float64(x) >= float64(y)
		case bool:
			return float64(x) >= ToFloat64(y)
		case string:
			return float64(x) >= ToFloat64(y)
		}
	case float64:
		switch y := b.(type) {
		case uint:
			return float64(x) >= float64(y)
		case uint8:
			return float64(x) >= float64(y)
		case uint16:
			return float64(x) >= float64(y)
		case uint32:
			return float64(x) >= float64(y)
		case uint64:
			return float64(x) >= float64(y)
		case int:
			return float64(x) >= float64(y)
		case int8:
			return float64(x) >= float64(y)
		case int16:
			return float64(x) >= float64(y)
		case int32:
			return float64(x) >= float64(y)
		case int64:
			return float64(x) >= float64(y)
		case float32:
			return float64(x) >= float64(y)
		case float64:
			return float64(x) >= float64(y)
		case bool:
			return float64(x) >= ToFloat64(y)
		case string:
			return float64(x) >= ToFloat64(y)
		}
	case bool:
		switch y := b.(type) {
		case uint:
			return ToFloat64(x) >= float64(y)
		case uint8:
			return ToFloat64(x) >= float64(y)
		case uint16:
			return ToFloat64(x) >= float64(y)
		case uint32:
			return ToFloat64(x) >= float64(y)
		case uint64:
			return ToFloat64(x) >= float64(y)
		case int:
			return ToFloat64(x) >= float64(y)
		case int8:
			return ToFloat64(x) >= float64(y)
		case int16:
			return ToFloat64(x) >= float64(y)
		case int32:
			return ToFloat64(x) >= float64(y)
		case int64:
			return ToFloat64(x) >= float64(y)
		case float32:
			return ToFloat64(x) >= float64(y)
		case float64:
			return ToFloat64(x) >= float64(y)
		case bool:
			return ToFloat64(x) >= ToFloat64(y)
		case string:
			return ToFloat64(x) >= ToFloat64(y)
		}
	case string:
		switch y := b.(type) {
		case uint:
			return ToFloat64(x) >= float64(y)
		case uint8:
			return ToFloat64(x) >= float64(y)
		case uint16:
			return ToFloat64(x) >= float64(y)
		case uint32:
			return ToFloat64(x) >= float64(y)
		case uint64:
			return ToFloat64(x) >= float64(y)
		case int:
			return ToFloat64(x) >= float64(y)
		case int8:
			return ToFloat64(x) >= float64(y)
		case int16:
			return ToFloat64(x) >= float64(y)
		case int32:
			return ToFloat64(x) >= float64(y)
		case int64:
			return ToFloat64(x) >= float64(y)
		case float32:
			return ToFloat64(x) >= float64(y)
		case float64:
			return ToFloat64(x) >= float64(y)
		case bool:
			return ToFloat64(x) >= ToFloat64(y)
		case string:
			return x >= y
		}
	case time.Time:
		switch y := b.(type) {
		case time.Time:
			return x.After(y) || x.Equal(y)
		}
	case time.Duration:
		switch y := b.(type) {
		case time.Duration:
			return x >= y
		}
	}
	panic(fmt.Sprintf("invalid operation: %T >= %T", a, b))
}

func Add(a, b interface{}) interface{} {
	// Handle nil values with proper string concatenation behavior
	if IsNil(a) && IsNil(b) {
		return 0 // nil + nil = 0
	}
	if IsNil(a) {
		switch y := b.(type) {
		case string:
			return "" + y // nil + string = "" + string (string concatenation)
		case float32, float64:
			return ToFloat64(a) + ToFloat64(y) // nil + float = 0.0 + float
		default:
			return ToInt(a) + ToInt(y) // nil + numeric = 0 + numeric
		}
	}
	if IsNil(b) {
		switch x := a.(type) {
		case string:
			return x + "" // string + nil = string + "" (string concatenation)
		case float32, float64:
			return ToFloat64(x) + ToFloat64(b) // float + nil = float + 0.0
		default:
			return ToInt(x) + ToInt(b) // numeric + nil = numeric + 0
		}
	}

	// Handle string concatenation - if either operand is string, concatenate
	if str, ok := a.(string); ok {
		return str + fmt.Sprint(b)
	}
	if str, ok := b.(string); ok {
		return fmt.Sprint(a) + str
	}

	// Handle boolean + boolean operations
	if aBool, aOk := a.(bool); aOk {
		if bBool, bOk := b.(bool); bOk {
			return ToInt(aBool) + ToInt(bBool) // true + true = 1 + 1 = 2
		}
		// bool + numeric - check if numeric is float
		switch b.(type) {
		case float32, float64:
			return ToFloat64(aBool) + ToFloat64(b)
		default:
			return ToInt(aBool) + ToInt(b)
		}
	}
	if bBool, bOk := b.(bool); bOk {
		// numeric + bool - check if numeric is float
		switch a.(type) {
		case float32, float64:
			return ToFloat64(a) + ToFloat64(bBool)
		default:
			return ToInt(a) + ToInt(bBool)
		}
	}

	// Handle numeric operations with type promotion
	switch x := a.(type) {
	case uint:
		switch y := b.(type) {
		case uint:
			return int(x) + int(y)
		case uint8:
			return int(x) + int(y)
		case uint16:
			return int(x) + int(y)
		case uint32:
			return int(x) + int(y)
		case uint64:
			return int(x) + int(y)
		case int:
			return int(x) + int(y)
		case int8:
			return int(x) + int(y)
		case int16:
			return int(x) + int(y)
		case int32:
			return int(x) + int(y)
		case int64:
			return int(x) + int(y)
		case float32:
			return float64(x) + float64(y)
		case float64:
			return float64(x) + float64(y)
		}
	case uint8:
		switch y := b.(type) {
		case uint:
			return int(x) + int(y)
		case uint8:
			return int(x) + int(y)
		case uint16:
			return int(x) + int(y)
		case uint32:
			return int(x) + int(y)
		case uint64:
			return int(x) + int(y)
		case int:
			return int(x) + int(y)
		case int8:
			return int(x) + int(y)
		case int16:
			return int(x) + int(y)
		case int32:
			return int(x) + int(y)
		case int64:
			return int(x) + int(y)
		case float32:
			return float64(x) + float64(y)
		case float64:
			return float64(x) + float64(y)
		}
	case uint16:
		switch y := b.(type) {
		case uint:
			return int(x) + int(y)
		case uint8:
			return int(x) + int(y)
		case uint16:
			return int(x) + int(y)
		case uint32:
			return int(x) + int(y)
		case uint64:
			return int(x) + int(y)
		case int:
			return int(x) + int(y)
		case int8:
			return int(x) + int(y)
		case int16:
			return int(x) + int(y)
		case int32:
			return int(x) + int(y)
		case int64:
			return int(x) + int(y)
		case float32:
			return float64(x) + float64(y)
		case float64:
			return float64(x) + float64(y)
		}
	case uint32:
		switch y := b.(type) {
		case uint:
			return int(x) + int(y)
		case uint8:
			return int(x) + int(y)
		case uint16:
			return int(x) + int(y)
		case uint32:
			return int(x) + int(y)
		case uint64:
			return int(x) + int(y)
		case int:
			return int(x) + int(y)
		case int8:
			return int(x) + int(y)
		case int16:
			return int(x) + int(y)
		case int32:
			return int(x) + int(y)
		case int64:
			return int(x) + int(y)
		case float32:
			return float64(x) + float64(y)
		case float64:
			return float64(x) + float64(y)
		}
	case uint64:
		switch y := b.(type) {
		case uint:
			return int(x) + int(y)
		case uint8:
			return int(x) + int(y)
		case uint16:
			return int(x) + int(y)
		case uint32:
			return int(x) + int(y)
		case uint64:
			return int(x) + int(y)
		case int:
			return int(x) + int(y)
		case int8:
			return int(x) + int(y)
		case int16:
			return int(x) + int(y)
		case int32:
			return int(x) + int(y)
		case int64:
			return int(x) + int(y)
		case float32:
			return float64(x) + float64(y)
		case float64:
			return float64(x) + float64(y)
		}
	case int:
		switch y := b.(type) {
		case uint:
			return int(x) + int(y)
		case uint8:
			return int(x) + int(y)
		case uint16:
			return int(x) + int(y)
		case uint32:
			return int(x) + int(y)
		case uint64:
			return int(x) + int(y)
		case int:
			return int(x) + int(y)
		case int8:
			return int(x) + int(y)
		case int16:
			return int(x) + int(y)
		case int32:
			return int(x) + int(y)
		case int64:
			return int(x) + int(y)
		case float32:
			return float64(x) + float64(y)
		case float64:
			return float64(x) + float64(y)
		}
	case int8:
		switch y := b.(type) {
		case uint:
			return int(x) + int(y)
		case uint8:
			return int(x) + int(y)
		case uint16:
			return int(x) + int(y)
		case uint32:
			return int(x) + int(y)
		case uint64:
			return int(x) + int(y)
		case int:
			return int(x) + int(y)
		case int8:
			return int(x) + int(y)
		case int16:
			return int(x) + int(y)
		case int32:
			return int(x) + int(y)
		case int64:
			return int(x) + int(y)
		case float32:
			return float64(x) + float64(y)
		case float64:
			return float64(x) + float64(y)
		}
	case int16:
		switch y := b.(type) {
		case uint:
			return int(x) + int(y)
		case uint8:
			return int(x) + int(y)
		case uint16:
			return int(x) + int(y)
		case uint32:
			return int(x) + int(y)
		case uint64:
			return int(x) + int(y)
		case int:
			return int(x) + int(y)
		case int8:
			return int(x) + int(y)
		case int16:
			return int(x) + int(y)
		case int32:
			return int(x) + int(y)
		case int64:
			return int(x) + int(y)
		case float32:
			return float64(x) + float64(y)
		case float64:
			return float64(x) + float64(y)
		}
	case int32:
		switch y := b.(type) {
		case uint:
			return int(x) + int(y)
		case uint8:
			return int(x) + int(y)
		case uint16:
			return int(x) + int(y)
		case uint32:
			return int(x) + int(y)
		case uint64:
			return int(x) + int(y)
		case int:
			return int(x) + int(y)
		case int8:
			return int(x) + int(y)
		case int16:
			return int(x) + int(y)
		case int32:
			return int(x) + int(y)
		case int64:
			return int(x) + int(y)
		case float32:
			return float64(x) + float64(y)
		case float64:
			return float64(x) + float64(y)
		}
	case int64:
		switch y := b.(type) {
		case uint:
			return int(x) + int(y)
		case uint8:
			return int(x) + int(y)
		case uint16:
			return int(x) + int(y)
		case uint32:
			return int(x) + int(y)
		case uint64:
			return int(x) + int(y)
		case int:
			return int(x) + int(y)
		case int8:
			return int(x) + int(y)
		case int16:
			return int(x) + int(y)
		case int32:
			return int(x) + int(y)
		case int64:
			return int(x) + int(y)
		case float32:
			return float64(x) + float64(y)
		case float64:
			return float64(x) + float64(y)
		}
	case float32:
		switch y := b.(type) {
		case uint:
			return float64(x) + float64(y)
		case uint8:
			return float64(x) + float64(y)
		case uint16:
			return float64(x) + float64(y)
		case uint32:
			return float64(x) + float64(y)
		case uint64:
			return float64(x) + float64(y)
		case int:
			return float64(x) + float64(y)
		case int8:
			return float64(x) + float64(y)
		case int16:
			return float64(x) + float64(y)
		case int32:
			return float64(x) + float64(y)
		case int64:
			return float64(x) + float64(y)
		case float32:
			return float64(x) + float64(y)
		case float64:
			return float64(x) + float64(y)
		}
	case float64:
		switch y := b.(type) {
		case uint:
			return float64(x) + float64(y)
		case uint8:
			return float64(x) + float64(y)
		case uint16:
			return float64(x) + float64(y)
		case uint32:
			return float64(x) + float64(y)
		case uint64:
			return float64(x) + float64(y)
		case int:
			return float64(x) + float64(y)
		case int8:
			return float64(x) + float64(y)
		case int16:
			return float64(x) + float64(y)
		case int32:
			return float64(x) + float64(y)
		case int64:
			return float64(x) + float64(y)
		case float32:
			return float64(x) + float64(y)
		case float64:
			return float64(x) + float64(y)
		}
	case time.Time:
		switch y := b.(type) {
		case time.Duration:
			return x.Add(y)
		}
	case time.Duration:
		switch y := b.(type) {
		case time.Time:
			return y.Add(x)
		case time.Duration:
			return x + y
		}
	}
	panic(fmt.Sprintf("invalid operation: %T + %T", a, b))
}

func Subtract(a, b interface{}) interface{} {
	// Handle nil values first - convert to numeric operations
	if IsNil(a) && IsNil(b) {
		return 0 // nil - nil = 0 - 0 = 0
	}
	if IsNil(a) {
		// nil - X: convert based on X's type to preserve precision
		switch b.(type) {
		case float32, float64:
			return ToFloat64(a) - ToFloat64(b) // nil - float = 0.0 - float
		default:
			return ToInt(a) - ToInt(b) // nil - numeric = 0 - numeric
		}
	}
	if IsNil(b) {
		// X - nil: convert based on X's type to preserve precision
		switch a.(type) {
		case float32, float64:
			return ToFloat64(a) - ToFloat64(b) // float - nil = float - 0.0
		default:
			return ToInt(a) - ToInt(b) // numeric - nil = numeric - 0
		}
	}

	// Handle boolean operations
	if aBool, aOk := a.(bool); aOk {
		if bBool, bOk := b.(bool); bOk {
			return ToInt(aBool) - ToInt(bBool) // true - false = 1 - 0 = 1
		}
		// bool - numeric - check if numeric is float
		switch b.(type) {
		case float32, float64:
			return ToFloat64(aBool) - ToFloat64(b)
		default:
			return ToInt(aBool) - ToInt(b)
		}
	}
	if bBool, bOk := b.(bool); bOk {
		// numeric - bool - check if numeric is float
		switch a.(type) {
		case float32, float64:
			return ToFloat64(a) - ToFloat64(bBool)
		default:
			return ToInt(a) - ToInt(bBool)
		}
	}

	// Handle string operations - convert to numeric or error
	if aStr, aOk := a.(string); aOk {
		// Check if string can be converted to number
		if _, ok := ToFloat64Safe(aStr); !ok {
			panic(fmt.Sprintf("invalid operation: string(%q) - %T", aStr, b))
		}
		if bStr, bOk := b.(string); bOk {
			// Check if both strings can be converted to numbers
			if _, ok := ToFloat64Safe(bStr); !ok {
				panic(fmt.Sprintf("invalid operation: string(%q) - string(%q)", aStr, bStr))
			}
			// string - string: convert both to numbers
			return ToInt(aStr) - ToInt(bStr)
		}
		// string - numeric: convert string to number, match numeric type
		switch b.(type) {
		case float32, float64:
			return ToFloat64(aStr) - ToFloat64(b)
		default:
			return ToInt(aStr) - ToInt(b)
		}
	}
	if bStr, bOk := b.(string); bOk {
		// Check if string can be converted to number
		if _, ok := ToFloat64Safe(bStr); !ok {
			panic(fmt.Sprintf("invalid operation: %T - string(%q)", a, bStr))
		}
		// numeric - string: convert string to number, match numeric type
		switch a.(type) {
		case float32, float64:
			return ToFloat64(a) - ToFloat64(bStr)
		default:
			return ToInt(a) - ToInt(bStr)
		}
	}

	// Handle numeric operations
	switch x := a.(type) {
	case uint:
		switch y := b.(type) {
		case uint:
			return int(x) - int(y)
		case uint8:
			return int(x) - int(y)
		case uint16:
			return int(x) - int(y)
		case uint32:
			return int(x) - int(y)
		case uint64:
			return int(x) - int(y)
		case int:
			return int(x) - int(y)
		case int8:
			return int(x) - int(y)
		case int16:
			return int(x) - int(y)
		case int32:
			return int(x) - int(y)
		case int64:
			return int(x) - int(y)
		case float32:
			return float64(x) - float64(y)
		case float64:
			return float64(x) - float64(y)
		}
	case uint8:
		switch y := b.(type) {
		case uint:
			return int(x) - int(y)
		case uint8:
			return int(x) - int(y)
		case uint16:
			return int(x) - int(y)
		case uint32:
			return int(x) - int(y)
		case uint64:
			return int(x) - int(y)
		case int:
			return int(x) - int(y)
		case int8:
			return int(x) - int(y)
		case int16:
			return int(x) - int(y)
		case int32:
			return int(x) - int(y)
		case int64:
			return int(x) - int(y)
		case float32:
			return float64(x) - float64(y)
		case float64:
			return float64(x) - float64(y)
		}
	case uint16:
		switch y := b.(type) {
		case uint:
			return int(x) - int(y)
		case uint8:
			return int(x) - int(y)
		case uint16:
			return int(x) - int(y)
		case uint32:
			return int(x) - int(y)
		case uint64:
			return int(x) - int(y)
		case int:
			return int(x) - int(y)
		case int8:
			return int(x) - int(y)
		case int16:
			return int(x) - int(y)
		case int32:
			return int(x) - int(y)
		case int64:
			return int(x) - int(y)
		case float32:
			return float64(x) - float64(y)
		case float64:
			return float64(x) - float64(y)
		}
	case uint32:
		switch y := b.(type) {
		case uint:
			return int(x) - int(y)
		case uint8:
			return int(x) - int(y)
		case uint16:
			return int(x) - int(y)
		case uint32:
			return int(x) - int(y)
		case uint64:
			return int(x) - int(y)
		case int:
			return int(x) - int(y)
		case int8:
			return int(x) - int(y)
		case int16:
			return int(x) - int(y)
		case int32:
			return int(x) - int(y)
		case int64:
			return int(x) - int(y)
		case float32:
			return float64(x) - float64(y)
		case float64:
			return float64(x) - float64(y)
		}
	case uint64:
		switch y := b.(type) {
		case uint:
			return int(x) - int(y)
		case uint8:
			return int(x) - int(y)
		case uint16:
			return int(x) - int(y)
		case uint32:
			return int(x) - int(y)
		case uint64:
			return int(x) - int(y)
		case int:
			return int(x) - int(y)
		case int8:
			return int(x) - int(y)
		case int16:
			return int(x) - int(y)
		case int32:
			return int(x) - int(y)
		case int64:
			return int(x) - int(y)
		case float32:
			return float64(x) - float64(y)
		case float64:
			return float64(x) - float64(y)
		}
	case int:
		switch y := b.(type) {
		case uint:
			return int(x) - int(y)
		case uint8:
			return int(x) - int(y)
		case uint16:
			return int(x) - int(y)
		case uint32:
			return int(x) - int(y)
		case uint64:
			return int(x) - int(y)
		case int:
			return int(x) - int(y)
		case int8:
			return int(x) - int(y)
		case int16:
			return int(x) - int(y)
		case int32:
			return int(x) - int(y)
		case int64:
			return int(x) - int(y)
		case float32:
			return float64(x) - float64(y)
		case float64:
			return float64(x) - float64(y)
		}
	case int8:
		switch y := b.(type) {
		case uint:
			return int(x) - int(y)
		case uint8:
			return int(x) - int(y)
		case uint16:
			return int(x) - int(y)
		case uint32:
			return int(x) - int(y)
		case uint64:
			return int(x) - int(y)
		case int:
			return int(x) - int(y)
		case int8:
			return int(x) - int(y)
		case int16:
			return int(x) - int(y)
		case int32:
			return int(x) - int(y)
		case int64:
			return int(x) - int(y)
		case float32:
			return float64(x) - float64(y)
		case float64:
			return float64(x) - float64(y)
		}
	case int16:
		switch y := b.(type) {
		case uint:
			return int(x) - int(y)
		case uint8:
			return int(x) - int(y)
		case uint16:
			return int(x) - int(y)
		case uint32:
			return int(x) - int(y)
		case uint64:
			return int(x) - int(y)
		case int:
			return int(x) - int(y)
		case int8:
			return int(x) - int(y)
		case int16:
			return int(x) - int(y)
		case int32:
			return int(x) - int(y)
		case int64:
			return int(x) - int(y)
		case float32:
			return float64(x) - float64(y)
		case float64:
			return float64(x) - float64(y)
		}
	case int32:
		switch y := b.(type) {
		case uint:
			return int(x) - int(y)
		case uint8:
			return int(x) - int(y)
		case uint16:
			return int(x) - int(y)
		case uint32:
			return int(x) - int(y)
		case uint64:
			return int(x) - int(y)
		case int:
			return int(x) - int(y)
		case int8:
			return int(x) - int(y)
		case int16:
			return int(x) - int(y)
		case int32:
			return int(x) - int(y)
		case int64:
			return int(x) - int(y)
		case float32:
			return float64(x) - float64(y)
		case float64:
			return float64(x) - float64(y)
		}
	case int64:
		switch y := b.(type) {
		case uint:
			return int(x) - int(y)
		case uint8:
			return int(x) - int(y)
		case uint16:
			return int(x) - int(y)
		case uint32:
			return int(x) - int(y)
		case uint64:
			return int(x) - int(y)
		case int:
			return int(x) - int(y)
		case int8:
			return int(x) - int(y)
		case int16:
			return int(x) - int(y)
		case int32:
			return int(x) - int(y)
		case int64:
			return int(x) - int(y)
		case float32:
			return float64(x) - float64(y)
		case float64:
			return float64(x) - float64(y)
		}
	case float32:
		switch y := b.(type) {
		case uint:
			return float64(x) - float64(y)
		case uint8:
			return float64(x) - float64(y)
		case uint16:
			return float64(x) - float64(y)
		case uint32:
			return float64(x) - float64(y)
		case uint64:
			return float64(x) - float64(y)
		case int:
			return float64(x) - float64(y)
		case int8:
			return float64(x) - float64(y)
		case int16:
			return float64(x) - float64(y)
		case int32:
			return float64(x) - float64(y)
		case int64:
			return float64(x) - float64(y)
		case float32:
			return float64(x) - float64(y)
		case float64:
			return float64(x) - float64(y)
		}
	case float64:
		switch y := b.(type) {
		case uint:
			return float64(x) - float64(y)
		case uint8:
			return float64(x) - float64(y)
		case uint16:
			return float64(x) - float64(y)
		case uint32:
			return float64(x) - float64(y)
		case uint64:
			return float64(x) - float64(y)
		case int:
			return float64(x) - float64(y)
		case int8:
			return float64(x) - float64(y)
		case int16:
			return float64(x) - float64(y)
		case int32:
			return float64(x) - float64(y)
		case int64:
			return float64(x) - float64(y)
		case float32:
			return float64(x) - float64(y)
		case float64:
			return float64(x) - float64(y)
		}
	case time.Time:
		switch y := b.(type) {
		case time.Time:
			return x.Sub(y)
		case time.Duration:
			return x.Add(-y)
		}
	case time.Duration:
		switch y := b.(type) {
		case time.Duration:
			return x - y
		}
	}
	panic(fmt.Sprintf("invalid operation: %T - %T", a, b))
}

func Multiply(a, b interface{}) interface{} {
	// Handle nil values first - any multiplication with nil results in 0
	if IsNil(a) || IsNil(b) {
		return 0
	}

	// Handle boolean operations
	if aBool, aOk := a.(bool); aOk {
		if bBool, bOk := b.(bool); bOk {
			return ToInt(aBool) * ToInt(bBool) // true * false = 1 * 0 = 0
		}
		// bool * numeric - check if numeric is float
		switch b.(type) {
		case float32, float64:
			return ToFloat64(aBool) * ToFloat64(b)
		default:
			return ToInt(aBool) * ToInt(b)
		}
	}
	if bBool, bOk := b.(bool); bOk {
		// numeric * bool - check if numeric is float
		switch a.(type) {
		case float32, float64:
			return ToFloat64(a) * ToFloat64(bBool)
		default:
			return ToInt(a) * ToInt(bBool)
		}
	}

	// Handle string operations - convert to numeric or error
	if aStr, aOk := a.(string); aOk {
		// Check if string can be converted to number
		if _, ok := ToFloat64Safe(aStr); !ok {
			panic(fmt.Sprintf("invalid operation: string(%q) * %T", aStr, b))
		}
		if bStr, bOk := b.(string); bOk {
			// Check if both strings can be converted to numbers
			if _, ok := ToFloat64Safe(bStr); !ok {
				panic(fmt.Sprintf("invalid operation: string(%q) * string(%q)", aStr, bStr))
			}
			// string * string: convert both to numbers
			return ToInt(aStr) * ToInt(bStr)
		}
		// string * numeric: convert string to number, match numeric type
		switch b.(type) {
		case float32, float64:
			return ToFloat64(aStr) * ToFloat64(b)
		default:
			return ToInt(aStr) * ToInt(b)
		}
	}
	if bStr, bOk := b.(string); bOk {
		// Check if string can be converted to number
		if _, ok := ToFloat64Safe(bStr); !ok {
			panic(fmt.Sprintf("invalid operation: %T * string(%q)", a, bStr))
		}
		// numeric * string: convert string to number, match numeric type
		switch a.(type) {
		case float32, float64:
			return ToFloat64(a) * ToFloat64(bStr)
		default:
			return ToInt(a) * ToInt(bStr)
		}
	}

	// Handle numeric operations
	switch x := a.(type) {
	case uint:
		switch y := b.(type) {
		case uint:
			return int(x) * int(y)
		case uint8:
			return int(x) * int(y)
		case uint16:
			return int(x) * int(y)
		case uint32:
			return int(x) * int(y)
		case uint64:
			return int(x) * int(y)
		case int:
			return int(x) * int(y)
		case int8:
			return int(x) * int(y)
		case int16:
			return int(x) * int(y)
		case int32:
			return int(x) * int(y)
		case int64:
			return int(x) * int(y)
		case float32:
			return float64(x) * float64(y)
		case float64:
			return float64(x) * float64(y)
		case time.Duration:
			return time.Duration(x) * time.Duration(y)
		}
	case uint8:
		switch y := b.(type) {
		case uint:
			return int(x) * int(y)
		case uint8:
			return int(x) * int(y)
		case uint16:
			return int(x) * int(y)
		case uint32:
			return int(x) * int(y)
		case uint64:
			return int(x) * int(y)
		case int:
			return int(x) * int(y)
		case int8:
			return int(x) * int(y)
		case int16:
			return int(x) * int(y)
		case int32:
			return int(x) * int(y)
		case int64:
			return int(x) * int(y)
		case float32:
			return float64(x) * float64(y)
		case float64:
			return float64(x) * float64(y)
		case time.Duration:
			return time.Duration(x) * time.Duration(y)
		}
	case uint16:
		switch y := b.(type) {
		case uint:
			return int(x) * int(y)
		case uint8:
			return int(x) * int(y)
		case uint16:
			return int(x) * int(y)
		case uint32:
			return int(x) * int(y)
		case uint64:
			return int(x) * int(y)
		case int:
			return int(x) * int(y)
		case int8:
			return int(x) * int(y)
		case int16:
			return int(x) * int(y)
		case int32:
			return int(x) * int(y)
		case int64:
			return int(x) * int(y)
		case float32:
			return float64(x) * float64(y)
		case float64:
			return float64(x) * float64(y)
		case time.Duration:
			return time.Duration(x) * time.Duration(y)
		}
	case uint32:
		switch y := b.(type) {
		case uint:
			return int(x) * int(y)
		case uint8:
			return int(x) * int(y)
		case uint16:
			return int(x) * int(y)
		case uint32:
			return int(x) * int(y)
		case uint64:
			return int(x) * int(y)
		case int:
			return int(x) * int(y)
		case int8:
			return int(x) * int(y)
		case int16:
			return int(x) * int(y)
		case int32:
			return int(x) * int(y)
		case int64:
			return int(x) * int(y)
		case float32:
			return float64(x) * float64(y)
		case float64:
			return float64(x) * float64(y)
		case time.Duration:
			return time.Duration(x) * time.Duration(y)
		}
	case uint64:
		switch y := b.(type) {
		case uint:
			return int(x) * int(y)
		case uint8:
			return int(x) * int(y)
		case uint16:
			return int(x) * int(y)
		case uint32:
			return int(x) * int(y)
		case uint64:
			return int(x) * int(y)
		case int:
			return int(x) * int(y)
		case int8:
			return int(x) * int(y)
		case int16:
			return int(x) * int(y)
		case int32:
			return int(x) * int(y)
		case int64:
			return int(x) * int(y)
		case float32:
			return float64(x) * float64(y)
		case float64:
			return float64(x) * float64(y)
		case time.Duration:
			return time.Duration(x) * time.Duration(y)
		}
	case int:
		switch y := b.(type) {
		case uint:
			return int(x) * int(y)
		case uint8:
			return int(x) * int(y)
		case uint16:
			return int(x) * int(y)
		case uint32:
			return int(x) * int(y)
		case uint64:
			return int(x) * int(y)
		case int:
			return int(x) * int(y)
		case int8:
			return int(x) * int(y)
		case int16:
			return int(x) * int(y)
		case int32:
			return int(x) * int(y)
		case int64:
			return int(x) * int(y)
		case float32:
			return float64(x) * float64(y)
		case float64:
			return float64(x) * float64(y)
		case time.Duration:
			return time.Duration(x) * time.Duration(y)
		}
	case int8:
		switch y := b.(type) {
		case uint:
			return int(x) * int(y)
		case uint8:
			return int(x) * int(y)
		case uint16:
			return int(x) * int(y)
		case uint32:
			return int(x) * int(y)
		case uint64:
			return int(x) * int(y)
		case int:
			return int(x) * int(y)
		case int8:
			return int(x) * int(y)
		case int16:
			return int(x) * int(y)
		case int32:
			return int(x) * int(y)
		case int64:
			return int(x) * int(y)
		case float32:
			return float64(x) * float64(y)
		case float64:
			return float64(x) * float64(y)
		case time.Duration:
			return time.Duration(x) * time.Duration(y)
		}
	case int16:
		switch y := b.(type) {
		case uint:
			return int(x) * int(y)
		case uint8:
			return int(x) * int(y)
		case uint16:
			return int(x) * int(y)
		case uint32:
			return int(x) * int(y)
		case uint64:
			return int(x) * int(y)
		case int:
			return int(x) * int(y)
		case int8:
			return int(x) * int(y)
		case int16:
			return int(x) * int(y)
		case int32:
			return int(x) * int(y)
		case int64:
			return int(x) * int(y)
		case float32:
			return float64(x) * float64(y)
		case float64:
			return float64(x) * float64(y)
		case time.Duration:
			return time.Duration(x) * time.Duration(y)
		}
	case int32:
		switch y := b.(type) {
		case uint:
			return int(x) * int(y)
		case uint8:
			return int(x) * int(y)
		case uint16:
			return int(x) * int(y)
		case uint32:
			return int(x) * int(y)
		case uint64:
			return int(x) * int(y)
		case int:
			return int(x) * int(y)
		case int8:
			return int(x) * int(y)
		case int16:
			return int(x) * int(y)
		case int32:
			return int(x) * int(y)
		case int64:
			return int(x) * int(y)
		case float32:
			return float64(x) * float64(y)
		case float64:
			return float64(x) * float64(y)
		case time.Duration:
			return time.Duration(x) * time.Duration(y)
		}
	case int64:
		switch y := b.(type) {
		case uint:
			return int(x) * int(y)
		case uint8:
			return int(x) * int(y)
		case uint16:
			return int(x) * int(y)
		case uint32:
			return int(x) * int(y)
		case uint64:
			return int(x) * int(y)
		case int:
			return int(x) * int(y)
		case int8:
			return int(x) * int(y)
		case int16:
			return int(x) * int(y)
		case int32:
			return int(x) * int(y)
		case int64:
			return int(x) * int(y)
		case float32:
			return float64(x) * float64(y)
		case float64:
			return float64(x) * float64(y)
		case time.Duration:
			return time.Duration(x) * time.Duration(y)
		}
	case float32:
		switch y := b.(type) {
		case uint:
			return float64(x) * float64(y)
		case uint8:
			return float64(x) * float64(y)
		case uint16:
			return float64(x) * float64(y)
		case uint32:
			return float64(x) * float64(y)
		case uint64:
			return float64(x) * float64(y)
		case int:
			return float64(x) * float64(y)
		case int8:
			return float64(x) * float64(y)
		case int16:
			return float64(x) * float64(y)
		case int32:
			return float64(x) * float64(y)
		case int64:
			return float64(x) * float64(y)
		case float32:
			return float64(x) * float64(y)
		case float64:
			return float64(x) * float64(y)
		case time.Duration:
			return float64(x) * float64(y)
		}
	case float64:
		switch y := b.(type) {
		case uint:
			return float64(x) * float64(y)
		case uint8:
			return float64(x) * float64(y)
		case uint16:
			return float64(x) * float64(y)
		case uint32:
			return float64(x) * float64(y)
		case uint64:
			return float64(x) * float64(y)
		case int:
			return float64(x) * float64(y)
		case int8:
			return float64(x) * float64(y)
		case int16:
			return float64(x) * float64(y)
		case int32:
			return float64(x) * float64(y)
		case int64:
			return float64(x) * float64(y)
		case float32:
			return float64(x) * float64(y)
		case float64:
			return float64(x) * float64(y)
		case time.Duration:
			return float64(x) * float64(y)
		}
	case time.Duration:
		switch y := b.(type) {
		case uint:
			return time.Duration(x) * time.Duration(y)
		case uint8:
			return time.Duration(x) * time.Duration(y)
		case uint16:
			return time.Duration(x) * time.Duration(y)
		case uint32:
			return time.Duration(x) * time.Duration(y)
		case uint64:
			return time.Duration(x) * time.Duration(y)
		case int:
			return time.Duration(x) * time.Duration(y)
		case int8:
			return time.Duration(x) * time.Duration(y)
		case int16:
			return time.Duration(x) * time.Duration(y)
		case int32:
			return time.Duration(x) * time.Duration(y)
		case int64:
			return time.Duration(x) * time.Duration(y)
		case float32:
			return float64(x) * float64(y)
		case float64:
			return float64(x) * float64(y)
		case time.Duration:
			return time.Duration(x) * time.Duration(y)
		}
	}
	panic(fmt.Sprintf("invalid operation: %T * %T", a, b))
}

func Divide(a, b interface{}) float64 {
	// Handle nil values first
	if IsNil(a) {
		a = 0
	}
	if IsNil(b) {
		b = 0
	}

	// Handle cross-type operations - convert booleans and strings to numbers
	switch x := a.(type) {
	case bool:
		a = ToInt(x)
	case string:
		// Check if string can be converted to number
		if _, ok := ToFloat64Safe(x); !ok {
			panic(fmt.Sprintf("invalid operation: string(%q) / %T", x, b))
		}
		a = ToFloat64(x) // Convert to float64 since division always returns float64
	}
	switch x := b.(type) {
	case bool:
		b = ToInt(x)
	case string:
		// Check if string can be converted to number
		if _, ok := ToFloat64Safe(x); !ok {
			panic(fmt.Sprintf("invalid operation: %T / string(%q)", a, x))
		}
		b = ToFloat64(x) // Convert to float64 since division always returns float64
	}

	// Check for division by zero after type conversion
	bVal := ToFloat64(b)
	if bVal == 0.0 {
		panic("integer divide by zero")
	}

	// Return the division result
	return ToFloat64(a) / bVal
}

func Modulo(a, b interface{}) interface{} {
	// Handle nil values first
	if IsNil(a) {
		a = 0
	}
	if IsNil(b) {
		b = 0
	}

	// Handle cross-type operations - convert booleans and strings to numbers
	switch x := a.(type) {
	case bool:
		a = ToInt(x)
	case string:
		// Check if string can be converted to number
		if _, ok := ToFloat64Safe(x); !ok {
			panic(fmt.Sprintf("invalid operation: string(%q) %% %T", x, b))
		}
		a = ToFloat64(x)
	}
	switch x := b.(type) {
	case bool:
		b = ToInt(x)
	case string:
		// Check if string can be converted to number
		if _, ok := ToFloat64Safe(x); !ok {
			panic(fmt.Sprintf("invalid operation: %T %% string(%q)", a, x))
		}
		b = ToFloat64(x)
	}

	// Check for modulo by zero after type conversion
	bVal := ToFloat64(b)
	if bVal == 0.0 {
		panic("integer divide by zero")
	}

	// For integer operations, return integer result
	aVal := ToFloat64(a)
	if isInteger(aVal) && isInteger(bVal) {
		return int(aVal) % int(bVal)
	}

	// For float operations, use math.Mod
	return math.Mod(aVal, bVal)
}

func isInteger(f float64) bool {
	return f == float64(int(f))
}

// IsTruthy determines if a value is truthy in JavaScript-like logic
func IsTruthy(v interface{}) bool {
	if IsNil(v) {
		return false
	}
	switch x := v.(type) {
	case bool:
		return x
	case int:
		return x != 0
	case int8:
		return x != 0
	case int16:
		return x != 0
	case int32:
		return x != 0
	case int64:
		return x != 0
	case uint:
		return x != 0
	case uint8:
		return x != 0
	case uint16:
		return x != 0
	case uint32:
		return x != 0
	case uint64:
		return x != 0
	case float32:
		return x != 0
	case float64:
		return x != 0
	case string:
		return x != ""
	case []interface{}:
		return true // Arrays are always truthy (even empty ones)
	case map[string]interface{}:
		return true // Objects are always truthy (even empty ones)
	default:
		return true // Other types are truthy
	}
}

// And performs JavaScript-like logical AND
func And(a, b interface{}) interface{} {
	if !IsTruthy(a) {
		return a
	}
	return b
}

// Or performs JavaScript-like logical OR
func Or(a, b interface{}) interface{} {
	if IsTruthy(a) {
		return a
	}
	return b
}
