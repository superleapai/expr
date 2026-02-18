package builtin

import (
	"fmt"
	"math"
	"reflect"
	"strconv"
	"strings"
	"time"
)

// ToBool converts any value to bool using SF-style coercion.
func ToBool(v interface{}) bool {
	if v == nil {
		return false
	}
	if rv := reflect.ValueOf(v); rv.Kind() == reflect.Ptr || rv.Kind() == reflect.Interface {
		if rv.IsNil() {
			return false
		}
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
	case time.Time:
		return true
	default:
		return true
	}
}

// CoerceToFloat64 converts any value to float64. Returns (value, wasNil).
func CoerceToFloat64(v interface{}) (float64, bool) {
	if v == nil {
		return 0, true
	}
	if rv := reflect.ValueOf(v); rv.Kind() == reflect.Ptr || rv.Kind() == reflect.Interface {
		if rv.IsNil() {
			return 0, true
		}
	}
	switch x := v.(type) {
	case bool:
		if x {
			return 1, false
		}
		return 0, false
	case int:
		return float64(x), false
	case int8:
		return float64(x), false
	case int16:
		return float64(x), false
	case int32:
		return float64(x), false
	case int64:
		return float64(x), false
	case uint:
		return float64(x), false
	case uint8:
		return float64(x), false
	case uint16:
		return float64(x), false
	case uint32:
		return float64(x), false
	case uint64:
		return float64(x), false
	case float32:
		return float64(x), false
	case float64:
		return x, false
	case string:
		if x == "" {
			return 0, false
		}
		if f, err := strconv.ParseFloat(x, 64); err == nil {
			return f, false
		}
		return 0, true // parse failed, treat as nil
	default:
		return 0, true
	}
}

// CoerceToString converts any value to string representation.
func CoerceToString(v interface{}) string {
	if v == nil {
		return ""
	}
	if rv := reflect.ValueOf(v); rv.Kind() == reflect.Ptr || rv.Kind() == reflect.Interface {
		if rv.IsNil() {
			return ""
		}
	}
	switch x := v.(type) {
	case string:
		return x
	case bool:
		if x {
			return "true"
		}
		return "false"
	case int:
		return strconv.Itoa(x)
	case int8:
		return strconv.FormatInt(int64(x), 10)
	case int16:
		return strconv.FormatInt(int64(x), 10)
	case int32:
		return strconv.FormatInt(int64(x), 10)
	case int64:
		return strconv.FormatInt(x, 10)
	case uint:
		return strconv.FormatUint(uint64(x), 10)
	case uint8:
		return strconv.FormatUint(uint64(x), 10)
	case uint16:
		return strconv.FormatUint(uint64(x), 10)
	case uint32:
		return strconv.FormatUint(uint64(x), 10)
	case uint64:
		return strconv.FormatUint(x, 10)
	case float32:
		return FormatFloat(float64(x))
	case float64:
		return FormatFloat(x)
	case time.Time:
		if x.Hour() == 0 && x.Minute() == 0 && x.Second() == 0 && x.Nanosecond() == 0 {
			return x.Format("2006-01-02")
		}
		return x.Format("2006-01-02 15:04:05")
	default:
		return fmt.Sprint(v)
	}
}

// FormatFloat formats a float64, trimming trailing .0 for whole numbers.
func FormatFloat(f float64) string {
	if math.IsInf(f, 0) || math.IsNaN(f) {
		return fmt.Sprint(f)
	}
	if f == float64(int64(f)) {
		return strconv.FormatInt(int64(f), 10)
	}
	return strconv.FormatFloat(f, 'f', -1, 64)
}

// CoerceToTime tries to convert a value to time.Time.
func CoerceToTime(v interface{}) (time.Time, bool) {
	if v == nil {
		return time.Time{}, false
	}
	switch x := v.(type) {
	case time.Time:
		return x, true
	case string:
		formats := []string{
			"2006-01-02",
			"2006-01-02 15:04:05",
			"2006-01-02T15:04:05Z07:00",
			"2006-01-02T15:04:05Z",
			"01/02/2006",
			"01/02/2006 15:04:05",
			"January 2, 2006",
		}
		for _, layout := range formats {
			if t, err := time.Parse(layout, x); err == nil {
				return t, true
			}
		}
		return time.Time{}, false
	default:
		return time.Time{}, false
	}
}

// EqualsCoerce compares two values with type coercion.
func EqualsCoerce(a, b interface{}) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		// nil vs "" should be true for ISPICKVAL compatibility
		if a == nil {
			if s, ok := b.(string); ok && s == "" {
				return true
			}
		}
		if b == nil {
			if s, ok := a.(string); ok && s == "" {
				return true
			}
		}
		return false
	}
	// Same types: direct compare
	if reflect.TypeOf(a) == reflect.TypeOf(b) {
		return reflect.DeepEqual(a, b)
	}
	// time.Time comparison
	if ta, ok := a.(time.Time); ok {
		if tb, ok2 := b.(time.Time); ok2 {
			return ta.Equal(tb)
		}
	}
	// Cross-type numeric comparison
	af, aOk := CoerceToFloat64(a)
	bf, bOk := CoerceToFloat64(b)
	if !aOk && !bOk { // both non-nil
		return af == bf
	}
	// String comparison as fallback
	return CoerceToString(a) == CoerceToString(b)
}

// IsBlankVal checks if a value is blank (nil, empty string, whitespace-only).
func IsBlankVal(v interface{}) bool {
	if v == nil {
		return true
	}
	if rv := reflect.ValueOf(v); rv.Kind() == reflect.Ptr || rv.Kind() == reflect.Interface {
		if rv.IsNil() {
			return true
		}
	}
	if s, ok := v.(string); ok {
		return strings.TrimSpace(s) == ""
	}
	return false
}

// IsNilVal checks if a value is nil (handles reflect nil).
func IsNilVal(v interface{}) bool {
	if v == nil {
		return true
	}
	rv := reflect.ValueOf(v)
	switch rv.Kind() {
	case reflect.Ptr, reflect.Interface, reflect.Map, reflect.Slice, reflect.Chan, reflect.Func:
		return rv.IsNil()
	}
	return false
}

// IsNumericVal checks if a value is or can be converted to a number.
func IsNumericVal(v interface{}) bool {
	if v == nil {
		return false
	}
	switch v.(type) {
	case int, int8, int16, int32, int64:
		return true
	case uint, uint8, uint16, uint32, uint64:
		return true
	case float32, float64:
		return true
	case string:
		s := v.(string)
		_, err := strconv.ParseFloat(s, 64)
		return err == nil
	default:
		return false
	}
}
