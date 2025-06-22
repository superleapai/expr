package runtime

import (
	"fmt"
	"math"
	"reflect"
	"time"
)

func Equal(a, b interface{}) bool {
	// Handle nil values first
	if IsNil(a) && IsNil(b) {
		return true
	}
	if a == nil || b == nil {
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
		}
	}
	if IsNil(a) && IsNil(b) {
		return true
	}
	return reflect.DeepEqual(a, b)
}

func Less(a, b interface{}) bool {
	// Handle nil values first
	if IsNil(a) && IsNil(b) {
		return false // nil is not less than nil
	}
	if a == nil {
		return true // nil is less than any non-nil value
	}
	if b == nil {
		return false // non-nil is not less than nil
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
		}
	case string:
		switch y := b.(type) {
		case string:
			return x < y
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
	// Handle nil values first
	if IsNil(a) && IsNil(b) {
		return false // nil is not more than nil
	}
	if a == nil {
		return false // nil is not more than any non-nil value
	}
	if b == nil {
		return true // non-nil is more than nil
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
		}
	case string:
		switch y := b.(type) {
		case string:
			return x > y
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
	// Handle nil values first
	if IsNil(a) && IsNil(b) {
		return true // nil is equal to nil
	}
	if a == nil {
		return true // nil is less than or equal to any non-nil value
	}
	if b == nil {
		return false // non-nil is not less than or equal to nil
	}

	// Handle numeric operations
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
		}
	case string:
		switch y := b.(type) {
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
	// Handle nil values first
	if IsNil(a) && IsNil(b) {
		return true // nil is equal to nil
	}
	if a == nil {
		return false // nil is not more than or equal to any non-nil value
	}
	if b == nil {
		return true // non-nil is more than or equal to nil
	}

	// Handle numeric operations
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
		}
	case string:
		switch y := b.(type) {
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
		if bStr, bOk := b.(string); bOk {
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
		if bStr, bOk := b.(string); bOk {
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
	if a == nil {
		return 0.0 // 0 divided by anything is 0
	}
	if b == nil {
		return 0.0 // Division by nil is treated as division by 0, which results in 0
	}
	// Handle cross-type operations - convert booleans and strings to numbers
	switch x := a.(type) {
	case bool:
		a = ToInt(x)
	case string:
		a = ToFloat64(x) // Convert to float64 since division always returns float64
	}
	switch x := b.(type) {
	case bool:
		b = ToInt(x)
	case string:
		b = ToFloat64(x) // Convert to float64 since division always returns float64
	}

	// Check for division by zero after type conversion
	bVal := ToFloat64(b)
	if bVal == 0.0 {
		panic("integer divide by zero")
	}

	// Return the division result
	return ToFloat64(a) / bVal
	switch x := a.(type) {
	case uint:
		switch y := b.(type) {
		case uint:
			if int(y) == 0 {
				panic("integer divide by zero")
			}
			return float64(x) / float64(y)
		case uint8:
			if int(y) == 0 {
				panic("integer divide by zero")
			}
			return float64(x) / float64(y)
		case uint16:
			if int(y) == 0 {
				panic("integer divide by zero")
			}
			return float64(x) / float64(y)
		case uint32:
			if int(y) == 0 {
				panic("integer divide by zero")
			}
			return float64(x) / float64(y)
		case uint64:
			if int(y) == 0 {
				panic("integer divide by zero")
			}
			return float64(x) / float64(y)
		case int:
			if int(y) == 0 {
				panic("integer divide by zero")
			}
			return float64(x) / float64(y)
		case int8:
			if int(y) == 0 {
				panic("integer divide by zero")
			}
			return float64(x) / float64(y)
		case int16:
			if int(y) == 0 {
				panic("integer divide by zero")
			}
			return float64(x) / float64(y)
		case int32:
			if int(y) == 0 {
				panic("integer divide by zero")
			}
			return float64(x) / float64(y)
		case int64:
			if int(y) == 0 {
				panic("integer divide by zero")
			}
			return float64(x) / float64(y)
		case float32:
			if float64(y) == 0.0 {
				panic("integer divide by zero")
			}
			return float64(x) / float64(y)
		case float64:
			if float64(y) == 0.0 {
				panic("integer divide by zero")
			}
			return float64(x) / float64(y)
		}
	case uint8:
		switch y := b.(type) {
		case uint:
			if int(y) == 0 {
				panic("integer divide by zero")
			}
			return float64(x) / float64(y)
		case uint8:
			if int(y) == 0 {
				panic("integer divide by zero")
			}
			return float64(x) / float64(y)
		case uint16:
			if int(y) == 0 {
				panic("integer divide by zero")
			}
			return float64(x) / float64(y)
		case uint32:
			if int(y) == 0 {
				panic("integer divide by zero")
			}
			return float64(x) / float64(y)
		case uint64:
			if int(y) == 0 {
				panic("integer divide by zero")
			}
			return float64(x) / float64(y)
		case int:
			if int(y) == 0 {
				panic("integer divide by zero")
			}
			return float64(x) / float64(y)
		case int8:
			if int(y) == 0 {
				panic("integer divide by zero")
			}
			return float64(x) / float64(y)
		case int16:
			if int(y) == 0 {
				panic("integer divide by zero")
			}
			return float64(x) / float64(y)
		case int32:
			if int(y) == 0 {
				panic("integer divide by zero")
			}
			return float64(x) / float64(y)
		case int64:
			if int(y) == 0 {
				panic("integer divide by zero")
			}
			return float64(x) / float64(y)
		case float32:
			if float64(y) == 0.0 {
				panic("integer divide by zero")
			}
			return float64(x) / float64(y)
		case float64:
			if float64(y) == 0.0 {
				panic("integer divide by zero")
			}
			return float64(x) / float64(y)
		}
	// Continue similar pattern for all other numeric types...
	case uint16:
		switch y := b.(type) {
		case uint:
			if int(y) == 0 {
				panic("integer divide by zero")
			}
			return float64(x) / float64(y)
		case uint8:
			if int(y) == 0 {
				panic("integer divide by zero")
			}
			return float64(x) / float64(y)
		case uint16:
			if int(y) == 0 {
				panic("integer divide by zero")
			}
			return float64(x) / float64(y)
		case uint32:
			if int(y) == 0 {
				panic("integer divide by zero")
			}
			return float64(x) / float64(y)
		case uint64:
			if int(y) == 0 {
				panic("integer divide by zero")
			}
			return float64(x) / float64(y)
		case int:
			if int(y) == 0 {
				panic("integer divide by zero")
			}
			return float64(x) / float64(y)
		case int8:
			if int(y) == 0 {
				panic("integer divide by zero")
			}
			return float64(x) / float64(y)
		case int16:
			if int(y) == 0 {
				panic("integer divide by zero")
			}
			return float64(x) / float64(y)
		case int32:
			if int(y) == 0 {
				panic("integer divide by zero")
			}
			return float64(x) / float64(y)
		case int64:
			if int(y) == 0 {
				panic("integer divide by zero")
			}
			return float64(x) / float64(y)
		case float32:
			if float64(y) == 0.0 {
				panic("integer divide by zero")
			}
			return float64(x) / float64(y)
		case float64:
			if float64(y) == 0.0 {
				panic("integer divide by zero")
			}
			return float64(x) / float64(y)
		}
	case uint32:
		switch y := b.(type) {
		case uint:
			if int(y) == 0 {
				panic("integer divide by zero")
			}
			return float64(x) / float64(y)
		case uint8:
			if int(y) == 0 {
				panic("integer divide by zero")
			}
			return float64(x) / float64(y)
		case uint16:
			if int(y) == 0 {
				panic("integer divide by zero")
			}
			return float64(x) / float64(y)
		case uint32:
			if int(y) == 0 {
				panic("integer divide by zero")
			}
			return float64(x) / float64(y)
		case uint64:
			if int(y) == 0 {
				panic("integer divide by zero")
			}
			return float64(x) / float64(y)
		case int:
			if int(y) == 0 {
				panic("integer divide by zero")
			}
			return float64(x) / float64(y)
		case int8:
			if int(y) == 0 {
				panic("integer divide by zero")
			}
			return float64(x) / float64(y)
		case int16:
			if int(y) == 0 {
				panic("integer divide by zero")
			}
			return float64(x) / float64(y)
		case int32:
			if int(y) == 0 {
				panic("integer divide by zero")
			}
			return float64(x) / float64(y)
		case int64:
			if int(y) == 0 {
				panic("integer divide by zero")
			}
			return float64(x) / float64(y)
		case float32:
			if float64(y) == 0.0 {
				panic("integer divide by zero")
			}
			return float64(x) / float64(y)
		case float64:
			if float64(y) == 0.0 {
				panic("integer divide by zero")
			}
			return float64(x) / float64(y)
		}
	case uint64:
		switch y := b.(type) {
		case uint:
			if int(y) == 0 {
				panic("integer divide by zero")
			}
			return float64(x) / float64(y)
		case uint8:
			if int(y) == 0 {
				panic("integer divide by zero")
			}
			return float64(x) / float64(y)
		case uint16:
			if int(y) == 0 {
				panic("integer divide by zero")
			}
			return float64(x) / float64(y)
		case uint32:
			if int(y) == 0 {
				panic("integer divide by zero")
			}
			return float64(x) / float64(y)
		case uint64:
			if int(y) == 0 {
				panic("integer divide by zero")
			}
			return float64(x) / float64(y)
		case int:
			if int(y) == 0 {
				panic("integer divide by zero")
			}
			return float64(x) / float64(y)
		case int8:
			if int(y) == 0 {
				panic("integer divide by zero")
			}
			return float64(x) / float64(y)
		case int16:
			if int(y) == 0 {
				panic("integer divide by zero")
			}
			return float64(x) / float64(y)
		case int32:
			if int(y) == 0 {
				panic("integer divide by zero")
			}
			return float64(x) / float64(y)
		case int64:
			if int(y) == 0 {
				panic("integer divide by zero")
			}
			return float64(x) / float64(y)
		case float32:
			if float64(y) == 0.0 {
				panic("integer divide by zero")
			}
			return float64(x) / float64(y)
		case float64:
			if float64(y) == 0.0 {
				panic("integer divide by zero")
			}
			return float64(x) / float64(y)
		}
	case int:
		switch y := b.(type) {
		case uint:
			if int(y) == 0 {
				panic("integer divide by zero")
			}
			return float64(x) / float64(y)
		case uint8:
			if int(y) == 0 {
				panic("integer divide by zero")
			}
			return float64(x) / float64(y)
		case uint16:
			if int(y) == 0 {
				panic("integer divide by zero")
			}
			return float64(x) / float64(y)
		case uint32:
			if int(y) == 0 {
				panic("integer divide by zero")
			}
			return float64(x) / float64(y)
		case uint64:
			if int(y) == 0 {
				panic("integer divide by zero")
			}
			return float64(x) / float64(y)
		case int:
			if int(y) == 0 {
				panic("integer divide by zero")
			}
			return float64(x) / float64(y)
		case int8:
			if int(y) == 0 {
				panic("integer divide by zero")
			}
			return float64(x) / float64(y)
		case int16:
			if int(y) == 0 {
				panic("integer divide by zero")
			}
			return float64(x) / float64(y)
		case int32:
			if int(y) == 0 {
				panic("integer divide by zero")
			}
			return float64(x) / float64(y)
		case int64:
			if int(y) == 0 {
				panic("integer divide by zero")
			}
			return float64(x) / float64(y)
		case float32:
			if float64(y) == 0.0 {
				panic("integer divide by zero")
			}
			return float64(x) / float64(y)
		case float64:
			if float64(y) == 0.0 {
				panic("integer divide by zero")
			}
			return float64(x) / float64(y)
		}
	case int8:
		switch y := b.(type) {
		case uint:
			if int(y) == 0 {
				panic("integer divide by zero")
			}
			return float64(x) / float64(y)
		case uint8:
			if int(y) == 0 {
				panic("integer divide by zero")
			}
			return float64(x) / float64(y)
		case uint16:
			if int(y) == 0 {
				panic("integer divide by zero")
			}
			return float64(x) / float64(y)
		case uint32:
			if int(y) == 0 {
				panic("integer divide by zero")
			}
			return float64(x) / float64(y)
		case uint64:
			if int(y) == 0 {
				panic("integer divide by zero")
			}
			return float64(x) / float64(y)
		case int:
			if int(y) == 0 {
				panic("integer divide by zero")
			}
			return float64(x) / float64(y)
		case int8:
			if int(y) == 0 {
				panic("integer divide by zero")
			}
			return float64(x) / float64(y)
		case int16:
			if int(y) == 0 {
				panic("integer divide by zero")
			}
			return float64(x) / float64(y)
		case int32:
			if int(y) == 0 {
				panic("integer divide by zero")
			}
			return float64(x) / float64(y)
		case int64:
			if int(y) == 0 {
				panic("integer divide by zero")
			}
			return float64(x) / float64(y)
		case float32:
			if float64(y) == 0.0 {
				panic("integer divide by zero")
			}
			return float64(x) / float64(y)
		case float64:
			if float64(y) == 0.0 {
				panic("integer divide by zero")
			}
			return float64(x) / float64(y)
		}
	case int16:
		switch y := b.(type) {
		case uint:
			if int(y) == 0 {
				panic("integer divide by zero")
			}
			return float64(x) / float64(y)
		case uint8:
			if int(y) == 0 {
				panic("integer divide by zero")
			}
			return float64(x) / float64(y)
		case uint16:
			if int(y) == 0 {
				panic("integer divide by zero")
			}
			return float64(x) / float64(y)
		case uint32:
			if int(y) == 0 {
				panic("integer divide by zero")
			}
			return float64(x) / float64(y)
		case uint64:
			if int(y) == 0 {
				panic("integer divide by zero")
			}
			return float64(x) / float64(y)
		case int:
			if int(y) == 0 {
				panic("integer divide by zero")
			}
			return float64(x) / float64(y)
		case int8:
			if int(y) == 0 {
				panic("integer divide by zero")
			}
			return float64(x) / float64(y)
		case int16:
			if int(y) == 0 {
				panic("integer divide by zero")
			}
			return float64(x) / float64(y)
		case int32:
			if int(y) == 0 {
				panic("integer divide by zero")
			}
			return float64(x) / float64(y)
		case int64:
			if int(y) == 0 {
				panic("integer divide by zero")
			}
			return float64(x) / float64(y)
		case float32:
			if float64(y) == 0.0 {
				panic("integer divide by zero")
			}
			return float64(x) / float64(y)
		case float64:
			if float64(y) == 0.0 {
				panic("integer divide by zero")
			}
			return float64(x) / float64(y)
		}
	case int32:
		switch y := b.(type) {
		case uint:
			if int(y) == 0 {
				panic("integer divide by zero")
			}
			return float64(x) / float64(y)
		case uint8:
			if int(y) == 0 {
				panic("integer divide by zero")
			}
			return float64(x) / float64(y)
		case uint16:
			if int(y) == 0 {
				panic("integer divide by zero")
			}
			return float64(x) / float64(y)
		case uint32:
			if int(y) == 0 {
				panic("integer divide by zero")
			}
			return float64(x) / float64(y)
		case uint64:
			if int(y) == 0 {
				panic("integer divide by zero")
			}
			return float64(x) / float64(y)
		case int:
			if int(y) == 0 {
				panic("integer divide by zero")
			}
			return float64(x) / float64(y)
		case int8:
			if int(y) == 0 {
				panic("integer divide by zero")
			}
			return float64(x) / float64(y)
		case int16:
			if int(y) == 0 {
				panic("integer divide by zero")
			}
			return float64(x) / float64(y)
		case int32:
			if int(y) == 0 {
				panic("integer divide by zero")
			}
			return float64(x) / float64(y)
		case int64:
			if int(y) == 0 {
				panic("integer divide by zero")
			}
			return float64(x) / float64(y)
		case float32:
			if float64(y) == 0.0 {
				panic("integer divide by zero")
			}
			return float64(x) / float64(y)
		case float64:
			if float64(y) == 0.0 {
				panic("integer divide by zero")
			}
			return float64(x) / float64(y)
		}
	case int64:
		switch y := b.(type) {
		case uint:
			if int(y) == 0 {
				panic("integer divide by zero")
			}
			return float64(x) / float64(y)
		case uint8:
			if int(y) == 0 {
				panic("integer divide by zero")
			}
			return float64(x) / float64(y)
		case uint16:
			if int(y) == 0 {
				panic("integer divide by zero")
			}
			return float64(x) / float64(y)
		case uint32:
			if int(y) == 0 {
				panic("integer divide by zero")
			}
			return float64(x) / float64(y)
		case uint64:
			if int(y) == 0 {
				panic("integer divide by zero")
			}
			return float64(x) / float64(y)
		case int:
			if int(y) == 0 {
				panic("integer divide by zero")
			}
			return float64(x) / float64(y)
		case int8:
			if int(y) == 0 {
				panic("integer divide by zero")
			}
			return float64(x) / float64(y)
		case int16:
			if int(y) == 0 {
				panic("integer divide by zero")
			}
			return float64(x) / float64(y)
		case int32:
			if int(y) == 0 {
				panic("integer divide by zero")
			}
			return float64(x) / float64(y)
		case int64:
			if int(y) == 0 {
				panic("integer divide by zero")
			}
			return float64(x) / float64(y)
		case float32:
			if float64(y) == 0.0 {
				panic("integer divide by zero")
			}
			return float64(x) / float64(y)
		case float64:
			if float64(y) == 0.0 {
				panic("integer divide by zero")
			}
			return float64(x) / float64(y)
		}
	case float32:
		switch y := b.(type) {
		case uint:
			if int(y) == 0 {
				panic("integer divide by zero")
			}
			return float64(x) / float64(y)
		case uint8:
			if int(y) == 0 {
				panic("integer divide by zero")
			}
			return float64(x) / float64(y)
		case uint16:
			if int(y) == 0 {
				panic("integer divide by zero")
			}
			return float64(x) / float64(y)
		case uint32:
			if int(y) == 0 {
				panic("integer divide by zero")
			}
			return float64(x) / float64(y)
		case uint64:
			if int(y) == 0 {
				panic("integer divide by zero")
			}
			return float64(x) / float64(y)
		case int:
			if int(y) == 0 {
				panic("integer divide by zero")
			}
			return float64(x) / float64(y)
		case int8:
			if int(y) == 0 {
				panic("integer divide by zero")
			}
			return float64(x) / float64(y)
		case int16:
			if int(y) == 0 {
				panic("integer divide by zero")
			}
			return float64(x) / float64(y)
		case int32:
			if int(y) == 0 {
				panic("integer divide by zero")
			}
			return float64(x) / float64(y)
		case int64:
			if int(y) == 0 {
				panic("integer divide by zero")
			}
			return float64(x) / float64(y)
		case float32:
			if float64(y) == 0.0 {
				panic("integer divide by zero")
			}
			return float64(x) / float64(y)
		case float64:
			if float64(y) == 0.0 {
				panic("integer divide by zero")
			}
			return float64(x) / float64(y)
		}
	case float64:
		switch y := b.(type) {
		case uint:
			if int(y) == 0 {
				panic("integer divide by zero")
			}
			return float64(x) / float64(y)
		case uint8:
			if int(y) == 0 {
				panic("integer divide by zero")
			}
			return float64(x) / float64(y)
		case uint16:
			if int(y) == 0 {
				panic("integer divide by zero")
			}
			return float64(x) / float64(y)
		case uint32:
			if int(y) == 0 {
				panic("integer divide by zero")
			}
			return float64(x) / float64(y)
		case uint64:
			if int(y) == 0 {
				panic("integer divide by zero")
			}
			return float64(x) / float64(y)
		case int:
			if int(y) == 0 {
				panic("integer divide by zero")
			}
			return float64(x) / float64(y)
		case int8:
			if int(y) == 0 {
				panic("integer divide by zero")
			}
			return float64(x) / float64(y)
		case int16:
			if int(y) == 0 {
				panic("integer divide by zero")
			}
			return float64(x) / float64(y)
		case int32:
			if int(y) == 0 {
				panic("integer divide by zero")
			}
			return float64(x) / float64(y)
		case int64:
			if int(y) == 0 {
				panic("integer divide by zero")
			}
			return float64(x) / float64(y)
		case float32:
			if float64(y) == 0.0 {
				panic("integer divide by zero")
			}
			return float64(x) / float64(y)
		case float64:
			if float64(y) == 0.0 {
				panic("integer divide by zero")
			}
			return float64(x) / float64(y)
		}
	}
	panic(fmt.Sprintf("invalid operation: %T / %T", a, b))
}

func Modulo(a, b interface{}) interface{} {
	// Handle nil values first
	if a == nil {
		return 0 // 0 modulo anything is 0
	}
	if b == nil {
		return 0 // Modulo by nil is treated as modulo by 0, which results in 0
	}

	// Handle numeric operations
	switch x := a.(type) {
	case int:
		switch y := b.(type) {
		case int:
			if int(y) == 0 {
				panic("integer divide by zero")
			}
			return int(x) % int(y)
		case int8:
			if int(y) == 0 {
				panic("integer divide by zero")
			}
			return int(x) % int(y)
		case int16:
			if int(y) == 0 {
				panic("integer divide by zero")
			}
			return int(x) % int(y)
		case int32:
			if int(y) == 0 {
				panic("integer divide by zero")
			}
			return int(x) % int(y)
		case int64:
			if int(y) == 0 {
				panic("integer divide by zero")
			}
			return int(x) % int(y)
		case float32:
			if float64(y) == 0.0 {
				panic("float modulo by zero")
			}
			return math.Mod(float64(x), float64(y))
		case float64:
			if float64(y) == 0.0 {
				panic("float modulo by zero")
			}
			return math.Mod(float64(x), float64(y))
		}
	case int8:
		switch y := b.(type) {
		case int:
			if int(y) == 0 {
				panic("integer divide by zero")
			}
			return int(x) % int(y)
		case int8:
			if int(y) == 0 {
				panic("integer divide by zero")
			}
			return int(x) % int(y)
		case int16:
			if int(y) == 0 {
				panic("integer divide by zero")
			}
			return int(x) % int(y)
		case int32:
			if int(y) == 0 {
				panic("integer divide by zero")
			}
			return int(x) % int(y)
		case int64:
			if int(y) == 0 {
				panic("integer divide by zero")
			}
			return int(x) % int(y)
		case float32:
			if float64(y) == 0.0 {
				panic("float modulo by zero")
			}
			return math.Mod(float64(x), float64(y))
		case float64:
			if float64(y) == 0.0 {
				panic("float modulo by zero")
			}
			return math.Mod(float64(x), float64(y))
		}
	case int16:
		switch y := b.(type) {
		case int:
			if int(y) == 0 {
				panic("integer divide by zero")
			}
			return int(x) % int(y)
		case int8:
			if int(y) == 0 {
				panic("integer divide by zero")
			}
			return int(x) % int(y)
		case int16:
			if int(y) == 0 {
				panic("integer divide by zero")
			}
			return int(x) % int(y)
		case int32:
			if int(y) == 0 {
				panic("integer divide by zero")
			}
			return int(x) % int(y)
		case int64:
			if int(y) == 0 {
				panic("integer divide by zero")
			}
			return int(x) % int(y)
		case float32:
			if float64(y) == 0.0 {
				panic("float modulo by zero")
			}
			return math.Mod(float64(x), float64(y))
		case float64:
			if float64(y) == 0.0 {
				panic("float modulo by zero")
			}
			return math.Mod(float64(x), float64(y))
		}
	case int32:
		switch y := b.(type) {
		case int:
			if int(y) == 0 {
				panic("integer divide by zero")
			}
			return int(x) % int(y)
		case int8:
			if int(y) == 0 {
				panic("integer divide by zero")
			}
			return int(x) % int(y)
		case int16:
			if int(y) == 0 {
				panic("integer divide by zero")
			}
			return int(x) % int(y)
		case int32:
			if int(y) == 0 {
				panic("integer divide by zero")
			}
			return int(x) % int(y)
		case int64:
			if int(y) == 0 {
				panic("integer divide by zero")
			}
			return int(x) % int(y)
		case float32:
			if float64(y) == 0.0 {
				panic("float modulo by zero")
			}
			return math.Mod(float64(x), float64(y))
		case float64:
			if float64(y) == 0.0 {
				panic("float modulo by zero")
			}
			return math.Mod(float64(x), float64(y))
		}
	case int64:
		switch y := b.(type) {
		case int:
			if int(y) == 0 {
				panic("integer divide by zero")
			}
			return int(x) % int(y)
		case int8:
			if int(y) == 0 {
				panic("integer divide by zero")
			}
			return int(x) % int(y)
		case int16:
			if int(y) == 0 {
				panic("integer divide by zero")
			}
			return int(x) % int(y)
		case int32:
			if int(y) == 0 {
				panic("integer divide by zero")
			}
			return int(x) % int(y)
		case int64:
			if int(y) == 0 {
				panic("integer divide by zero")
			}
			return int(x) % int(y)
		case float32:
			if float64(y) == 0.0 {
				panic("float modulo by zero")
			}
			return math.Mod(float64(x), float64(y))
		case float64:
			if float64(y) == 0.0 {
				panic("float modulo by zero")
			}
			return math.Mod(float64(x), float64(y))
		}
	case float32:
		switch y := b.(type) {
		case int:
			if float64(y) == 0.0 {
				panic("float modulo by zero")
			}
			return math.Mod(float64(x), float64(y))
		case int8:
			if float64(y) == 0.0 {
				panic("float modulo by zero")
			}
			return math.Mod(float64(x), float64(y))
		case int16:
			if float64(y) == 0.0 {
				panic("float modulo by zero")
			}
			return math.Mod(float64(x), float64(y))
		case int32:
			if float64(y) == 0.0 {
				panic("float modulo by zero")
			}
			return math.Mod(float64(x), float64(y))
		case int64:
			if float64(y) == 0.0 {
				panic("float modulo by zero")
			}
			return math.Mod(float64(x), float64(y))
		case float32:
			if float64(y) == 0.0 {
				panic("float modulo by zero")
			}
			return math.Mod(float64(x), float64(y))
		case float64:
			if float64(y) == 0.0 {
				panic("float modulo by zero")
			}
			return math.Mod(float64(x), float64(y))
		}
	case float64:
		switch y := b.(type) {
		case int:
			if float64(y) == 0.0 {
				panic("float modulo by zero")
			}
			return math.Mod(float64(x), float64(y))
		case int8:
			if float64(y) == 0.0 {
				panic("float modulo by zero")
			}
			return math.Mod(float64(x), float64(y))
		case int16:
			if float64(y) == 0.0 {
				panic("float modulo by zero")
			}
			return math.Mod(float64(x), float64(y))
		case int32:
			if float64(y) == 0.0 {
				panic("float modulo by zero")
			}
			return math.Mod(float64(x), float64(y))
		case int64:
			if float64(y) == 0.0 {
				panic("float modulo by zero")
			}
			return math.Mod(float64(x), float64(y))
		case float32:
			if float64(y) == 0.0 {
				panic("float modulo by zero")
			}
			return math.Mod(float64(x), float64(y))
		case float64:
			if float64(y) == 0.0 {
				panic("float modulo by zero")
			}
			return math.Mod(float64(x), float64(y))
		}

	}
	panic(fmt.Sprintf("invalid operation: %T %% %T", a, b))
}
