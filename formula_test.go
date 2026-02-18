package expr_test

import (
	"math"
	"testing"
	"time"

	"github.com/expr-lang/expr"
	"github.com/expr-lang/expr/internal/testify/assert"
	"github.com/expr-lang/expr/internal/testify/require"
)

// =============================================================================
// TRUE / FALSE / NULL keywords
// =============================================================================

func TestKeywords_TRUE_FALSE_NULL(t *testing.T) {
	tests := []struct {
		expr string
		want any
	}{
		{"TRUE", true},
		{"FALSE", false},
		{"NULL", nil},
		{"TRUE && true", true},
		{"FALSE || TRUE", true},
		{"TRUE && FALSE", false},
		{"NULL == nil", true},
		{"NULL == NULL", true},
		{"TRUE != FALSE", true},
	}
	for _, tt := range tests {
		t.Run(tt.expr, func(t *testing.T) {
			program, err := expr.Compile(tt.expr, expr.AllowUndefinedVariables(), expr.AsAny())
			require.NoError(t, err)
			got, err := expr.Run(program, nil)
			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

// =============================================================================
// <> not-equal operator
// =============================================================================

func TestOperator_NotEqual_Diamond(t *testing.T) {
	tests := []struct {
		expr string
		want bool
	}{
		{"1 <> 2", true},
		{"1 <> 1", false},
		{`"a" <> "b"`, true},
		{`"a" <> "a"`, false},
		{"true <> false", true},
		{"true <> true", false},
	}
	for _, tt := range tests {
		t.Run(tt.expr, func(t *testing.T) {
			program, err := expr.Compile(tt.expr, expr.AllowUndefinedVariables(), expr.AsAny())
			require.NoError(t, err)
			got, err := expr.Run(program, nil)
			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

// =============================================================================
// & string concatenation operator
// =============================================================================

func TestOperator_Ampersand_Concat(t *testing.T) {
	tests := []struct {
		expr string
		want string
	}{
		{`"hello" & " world"`, "hello world"},
		{`"count: " & 5`, "count: 5"},
		{`"pi: " & 3.14`, "pi: 3.14"},
		{`"bool: " & true`, "bool: true"},
		{`"a" & "b" & "c"`, "abc"},
	}
	for _, tt := range tests {
		t.Run(tt.expr, func(t *testing.T) {
			program, err := expr.Compile(tt.expr, expr.AllowUndefinedVariables(), expr.AsAny())
			require.NoError(t, err)
			got, err := expr.Run(program, nil)
			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestOperator_Ampersand_Concat_Nil(t *testing.T) {
	program, err := expr.Compile(`"hello" & NULL`, expr.AllowUndefinedVariables(), expr.AsAny())
	require.NoError(t, err)
	got, err := expr.Run(program, nil)
	require.NoError(t, err)
	assert.Equal(t, "hello", got)
}

// =============================================================================
// Date comparisons
// =============================================================================

func TestDateComparison(t *testing.T) {
	d1 := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
	d2 := time.Date(2023, 6, 15, 0, 0, 0, 0, time.UTC)
	env := map[string]any{"d1": d1, "d2": d2}

	tests := []struct {
		expr string
		want bool
	}{
		{"d1 < d2", true},
		{"d2 < d1", false},
		{"d1 > d2", false},
		{"d2 > d1", true},
		{"d1 <= d2", true},
		{"d1 <= d1", true},
		{"d1 >= d2", false},
		{"d1 >= d1", true},
		{"d1 == d1", true},
		{"d1 == d2", false},
	}
	for _, tt := range tests {
		t.Run(tt.expr, func(t *testing.T) {
			program, err := expr.Compile(tt.expr, expr.Env(env), expr.AsAny())
			require.NoError(t, err)
			got, err := expr.Run(program, env)
			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

// =============================================================================
// ContextFunction support
// =============================================================================

func TestContextFunction(t *testing.T) {
	env := map[string]any{
		"_tz": time.FixedZone("UTC+5", 5*3600),
	}

	program, err := expr.Compile("MY_FUNC()",
		expr.AllowUndefinedVariables(),
		expr.AsAny(),
		expr.ContextFunction("MY_FUNC", func(env any, params ...any) (any, error) {
			return "got_context", nil
		}, new(func() string)),
	)
	require.NoError(t, err)
	got, err := expr.Run(program, env)
	require.NoError(t, err)
	assert.Equal(t, "got_context", got)
}

func TestContextFunction_WithArgs(t *testing.T) {
	program, err := expr.Compile(`MY_FUNC("hello")`,
		expr.AllowUndefinedVariables(),
		expr.AsAny(),
		expr.ContextFunction("MY_FUNC", func(env any, params ...any) (any, error) {
			return params[0].(string) + "_modified", nil
		}, new(func(string) string)),
	)
	require.NoError(t, err)
	got, err := expr.Run(program, map[string]any{})
	require.NoError(t, err)
	assert.Equal(t, "hello_modified", got)
}

// =============================================================================
// Case-insensitive function lookup (for user-registered Functions only)
// =============================================================================

func TestCaseInsensitiveFunction(t *testing.T) {
	tests := []struct {
		name string
		expr string
	}{
		{"exact case", "MYFUNC(1)"},
		{"lower case", "myfunc(1)"},
		{"mixed case", "MyFunc(1)"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			program, err := expr.Compile(tt.expr,
				expr.AllowUndefinedVariables(),
				expr.AsAny(),
				expr.Function("MYFUNC", func(params ...any) (any, error) {
					return params[0], nil
				}, new(func(any) any)),
			)
			require.NoError(t, err)
			got, err := expr.Run(program, nil)
			require.NoError(t, err)
			assert.Equal(t, 1, got)
		})
	}
}

func TestCaseInsensitive_DoesNotAffectBuiltins(t *testing.T) {
	// "String" is a variable, not the builtin "string" function
	env := map[string]any{
		"String": "hello world",
	}
	program, err := expr.Compile("String[:5]", expr.Env(env), expr.AsAny())
	require.NoError(t, err)
	got, err := expr.Run(program, env)
	require.NoError(t, err)
	assert.Equal(t, "hello", got)
}

// =============================================================================
// Logical function pack
// =============================================================================

func TestPack_Logical(t *testing.T) {
	opts := []expr.Option{
		expr.AllowUndefinedVariables(),
		expr.AsAny(),
		expr.WithLogicalFunctions(),
	}

	tests := []struct {
		expr string
		want any
	}{
		{"AND(true, true)", true},
		{"AND(true, false)", false},
		{"AND(true, true, true)", true},
		{"OR(false, false)", false},
		{"OR(false, true)", true},
		{"OR(false, false, true)", true},
		{"NOT(true)", false},
		{"NOT(false)", true},
		{"IF(true, 1, 2)", 1},
		{"IF(false, 1, 2)", 2},
		{"ISBLANK(NULL)", true},
		{`ISBLANK("")`, true},
		{"ISBLANK(0)", false},
		{`ISBLANK("hello")`, false},
		{"ISNULL(NULL)", true},
		{"ISNULL(0)", false},
		{"ISNUMBER(42)", true},
		{"ISNUMBER(3.14)", true},
		{`ISNUMBER("hello")`, false},
		{"ISNUMBER(NULL)", false},
	}
	for _, tt := range tests {
		t.Run(tt.expr, func(t *testing.T) {
			program, err := expr.Compile(tt.expr, opts...)
			require.NoError(t, err)
			got, err := expr.Run(program, nil)
			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestPack_Logical_CASE(t *testing.T) {
	opts := []expr.Option{
		expr.AllowUndefinedVariables(),
		expr.AsAny(),
		expr.WithLogicalFunctions(),
	}
	// CASE(expr, val1, result1, val2, result2, ..., default)
	program, err := expr.Compile(`CASE("B", "A", 1, "B", 2, "C", 3, 0)`, opts...)
	require.NoError(t, err)
	got, err := expr.Run(program, nil)
	require.NoError(t, err)
	assert.Equal(t, 2, got)
}

func TestPack_Logical_CASE_Default(t *testing.T) {
	opts := []expr.Option{
		expr.AllowUndefinedVariables(),
		expr.AsAny(),
		expr.WithLogicalFunctions(),
	}
	program, err := expr.Compile(`CASE("Z", "A", 1, "B", 2, 99)`, opts...)
	require.NoError(t, err)
	got, err := expr.Run(program, nil)
	require.NoError(t, err)
	assert.Equal(t, 99, got)
}

// =============================================================================
// Math function pack
// =============================================================================

func TestPack_Math(t *testing.T) {
	opts := []expr.Option{
		expr.AllowUndefinedVariables(),
		expr.AsAny(),
		expr.WithMathFunctions(),
	}

	tests := []struct {
		expr string
		want any
	}{
		{"ABS(-5)", 5.0},
		{"ABS(5)", 5.0},
		{"ABS(NULL)", nil},
		{"CEILING(4.2)", 5.0},
		{"CEILING(-4.2)", -4.0},
		{"FLOOR(4.8)", 4.0},
		{"FLOOR(-4.8)", -5.0},
		{"SQRT(25)", 5.0},
		{"SQRT(NULL)", nil},
		{"SQRT(-1)", nil},
		{"LOG(100)", 2.0},
		{"LOG(NULL)", nil},
		{"LN(1)", 0.0},
		{"EXP(0)", 1.0},
		{"MCEILING(4.2)", 5.0},
		{"MCEILING(-4.2)", -5.0},
		{"MFLOOR(4.8)", 4.0},
		{"MFLOOR(-4.8)", -4.0},
		{"TRUNC(4.9)", 4.0},
		{"TRUNC(-4.9)", -4.0},
	}
	for _, tt := range tests {
		t.Run(tt.expr, func(t *testing.T) {
			program, err := expr.Compile(tt.expr, opts...)
			require.NoError(t, err)
			got, err := expr.Run(program, nil)
			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestPack_Math_ROUND(t *testing.T) {
	opts := []expr.Option{
		expr.AllowUndefinedVariables(),
		expr.AsAny(),
		expr.WithMathFunctions(),
	}
	tests := []struct {
		expr string
		want float64
	}{
		{"ROUND(2.5, 0)", 3.0},   // half-away-from-zero
		{"ROUND(-2.5, 0)", -3.0}, // half-away-from-zero
		{"ROUND(3.1415, 2)", 3.14},
		{"ROUND(123.456, 1)", 123.5},
	}
	for _, tt := range tests {
		t.Run(tt.expr, func(t *testing.T) {
			program, err := expr.Compile(tt.expr, opts...)
			require.NoError(t, err)
			got, err := expr.Run(program, nil)
			require.NoError(t, err)
			assert.InDelta(t, tt.want, got.(float64), 1e-9)
		})
	}
}

func TestPack_Math_MAX_MIN(t *testing.T) {
	opts := []expr.Option{
		expr.AllowUndefinedVariables(),
		expr.AsAny(),
		expr.WithMathFunctions(),
	}
	tests := []struct {
		expr string
		want any
	}{
		{"MAX(3, 7)", 7.0},
		{"MAX(-1, -5)", -1.0},
		{"MAX(NULL, 5)", 5.0},
		{"MAX(5, NULL)", 5.0},
		{"MAX(NULL, NULL)", nil},
		{"MIN(3, 7)", 3.0},
		{"MIN(-1, -5)", -5.0},
		{"MIN(NULL, 5)", 5.0},
		{"MIN(5, NULL)", 5.0},
		{"MIN(NULL, NULL)", nil},
	}
	for _, tt := range tests {
		t.Run(tt.expr, func(t *testing.T) {
			program, err := expr.Compile(tt.expr, opts...)
			require.NoError(t, err)
			got, err := expr.Run(program, nil)
			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestPack_Math_MOD(t *testing.T) {
	opts := []expr.Option{
		expr.AllowUndefinedVariables(),
		expr.AsAny(),
		expr.WithMathFunctions(),
	}
	tests := []struct {
		expr string
		want any
	}{
		{"MOD(10, 3)", 1.0},
		{"MOD(10, 0)", nil},
		{"MOD(NULL, 3)", nil},
	}
	for _, tt := range tests {
		t.Run(tt.expr, func(t *testing.T) {
			program, err := expr.Compile(tt.expr, opts...)
			require.NoError(t, err)
			got, err := expr.Run(program, nil)
			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

// =============================================================================
// Trig function pack
// =============================================================================

func TestPack_Trig(t *testing.T) {
	opts := []expr.Option{
		expr.AllowUndefinedVariables(),
		expr.AsAny(),
		expr.WithTrigFunctions(),
	}

	tests := []struct {
		expr string
		want float64
	}{
		{"SIN(0)", 0.0},
		{"COS(0)", 1.0},
		{"TAN(0)", 0.0},
		{"ASIN(0)", 0.0},
		{"ACOS(1)", 0.0},
		{"ATAN(0)", 0.0},
		{"ATAN2(1, 1)", math.Pi / 4},
	}
	for _, tt := range tests {
		t.Run(tt.expr, func(t *testing.T) {
			program, err := expr.Compile(tt.expr, opts...)
			require.NoError(t, err)
			got, err := expr.Run(program, nil)
			require.NoError(t, err)
			assert.InDelta(t, tt.want, got.(float64), 1e-9)
		})
	}
}

func TestPack_Trig_PI(t *testing.T) {
	program, err := expr.Compile("PI()",
		expr.AllowUndefinedVariables(),
		expr.AsAny(),
		expr.WithTrigFunctions(),
	)
	require.NoError(t, err)
	got, err := expr.Run(program, nil)
	require.NoError(t, err)
	assert.Equal(t, math.Pi, got)
}

func TestPack_Trig_NullPropagation(t *testing.T) {
	opts := []expr.Option{
		expr.AllowUndefinedVariables(),
		expr.AsAny(),
		expr.WithTrigFunctions(),
	}
	tests := []struct {
		expr string
	}{
		{"SIN(NULL)"},
		{"COS(NULL)"},
		{"TAN(NULL)"},
		{"ASIN(NULL)"},
		{"ACOS(NULL)"},
		{"ATAN(NULL)"},
	}
	for _, tt := range tests {
		t.Run(tt.expr, func(t *testing.T) {
			program, err := expr.Compile(tt.expr, opts...)
			require.NoError(t, err)
			got, err := expr.Run(program, nil)
			require.NoError(t, err)
			assert.Nil(t, got)
		})
	}
}

func TestPack_Trig_OutOfRange(t *testing.T) {
	opts := []expr.Option{
		expr.AllowUndefinedVariables(),
		expr.AsAny(),
		expr.WithTrigFunctions(),
	}
	// ASIN/ACOS with |x| > 1 should return nil
	tests := []struct {
		expr string
	}{
		{"ASIN(2)"},
		{"ASIN(-2)"},
		{"ACOS(2)"},
		{"ACOS(-2)"},
	}
	for _, tt := range tests {
		t.Run(tt.expr, func(t *testing.T) {
			program, err := expr.Compile(tt.expr, opts...)
			require.NoError(t, err)
			got, err := expr.Run(program, nil)
			require.NoError(t, err)
			assert.Nil(t, got)
		})
	}
}

// =============================================================================
// Text function pack
// =============================================================================

func TestPack_Text(t *testing.T) {
	opts := []expr.Option{
		expr.AllowUndefinedVariables(),
		expr.AsAny(),
		expr.WithTextFunctions(),
	}

	tests := []struct {
		expr string
		want any
	}{
		{`BEGINS("Hello World", "Hello")`, true},
		{`BEGINS("Hello World", "World")`, false},
		{`BEGINS(NULL, "test")`, false},
		{`CONTAINS("Hello World", "World")`, true},
		{`CONTAINS("Hello World", "xyz")`, false},
		{`CONTAINS(NULL, "test")`, false},
		{`FIND("World", "Hello World")`, 7},
		{`FIND("xyz", "Hello World")`, 0},
		{`LEFT("Hello", 3)`, "Hel"},
		{`LEFT("Hi", 10)`, "Hi"},
		{`LEFT(NULL, 3)`, nil},
		{`RIGHT("Hello", 3)`, "llo"},
		{`RIGHT("Hi", 10)`, "Hi"},
		{`MID("Hello World", 7, 5)`, "World"},
		{`LEN("Hello")`, 5},
		{"LEN(NULL)", 0},
		{`LOWER("HELLO")`, "hello"},
		{"LOWER(NULL)", nil},
		{`UPPER("hello")`, "HELLO"},
		{"UPPER(NULL)", nil},
		{`TRIM("  hello  ")`, "hello"},
		{"TRIM(NULL)", nil},
		{`INITCAP("hello world")`, "Hello World"},
		{`SUBSTITUTE("Hello World", "World", "Go")`, "Hello Go"},
		{"SUBSTITUTE(NULL, \"a\", \"b\")", nil},
		{`LPAD("hi", 5)`, "   hi"},
		{`LPAD("hi", 5, "0")`, "000hi"},
		{`RPAD("hi", 5)`, "hi   "},
		{`RPAD("hi", 5, "0")`, "hi000"},
		{`TEXT(42)`, "42"},
		{`TEXT(3.14)`, "3.14"},
		{"TEXT(NULL)", ""},
		{`VALUE("42")`, 42.0},
		{`VALUE("3.14")`, 3.14},
		{"VALUE(NULL)", nil},
		{`VALUE("abc")`, nil},
		{`REVERSE("Hello")`, "olleH"},
		{"REVERSE(NULL)", nil},
		{"BR()", "\n"},
		{`ASCII("A")`, 65},
		{"ASCII(NULL)", nil},
		{`CHR(65)`, "A"},
		{"CHR(NULL)", nil},
	}
	for _, tt := range tests {
		t.Run(tt.expr, func(t *testing.T) {
			program, err := expr.Compile(tt.expr, opts...)
			require.NoError(t, err)
			got, err := expr.Run(program, nil)
			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestPack_Text_CaseInsensitive(t *testing.T) {
	opts := []expr.Option{
		expr.AllowUndefinedVariables(),
		expr.AsAny(),
		expr.WithTextFunctions(),
	}
	// Test that function lookup is case-insensitive
	tests := []struct {
		expr string
		want any
	}{
		{`lower("HELLO")`, "hello"},
		{`Upper("hello")`, "HELLO"},
		{`Len("abc")`, 3},
		{`trim("  hi  ")`, "hi"},
	}
	for _, tt := range tests {
		t.Run(tt.expr, func(t *testing.T) {
			program, err := expr.Compile(tt.expr, opts...)
			require.NoError(t, err)
			got, err := expr.Run(program, nil)
			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

// =============================================================================
// DateTime function pack
// =============================================================================

func TestPack_DateTime(t *testing.T) {
	opts := []expr.Option{
		expr.AllowUndefinedVariables(),
		expr.AsAny(),
		expr.WithDateTimeFunctions(),
	}

	d := time.Date(2023, 6, 15, 14, 30, 45, 123000000, time.UTC)
	env := map[string]any{"dt": d}

	tests := []struct {
		expr string
		want any
	}{
		{"DAY(dt)", 15},
		{"MONTH(dt)", 6},
		{"YEAR(dt)", 2023},
		{"WEEKDAY(dt)", 5}, // Thursday = 5 in SF convention (1=Sunday)
		{"HOUR(dt)", 14},
		{"MINUTE(dt)", 30},
		{"SECOND(dt)", 45},
		{"MILLISECOND(dt)", 123},
		{"DAYOFYEAR(dt)", 166},
		{"DAY(NULL)", nil},
		{"MONTH(NULL)", nil},
		{"YEAR(NULL)", nil},
	}
	for _, tt := range tests {
		t.Run(tt.expr, func(t *testing.T) {
			program, err := expr.Compile(tt.expr, append(opts, expr.Env(env))...)
			require.NoError(t, err)
			got, err := expr.Run(program, env)
			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestPack_DateTime_DATE(t *testing.T) {
	opts := []expr.Option{
		expr.AllowUndefinedVariables(),
		expr.AsAny(),
		expr.WithDateTimeFunctions(),
	}
	program, err := expr.Compile("DATE(2023, 6, 15)", opts...)
	require.NoError(t, err)
	got, err := expr.Run(program, nil)
	require.NoError(t, err)
	expected := time.Date(2023, 6, 15, 0, 0, 0, 0, time.UTC)
	assert.Equal(t, expected, got)
}

func TestPack_DateTime_ADDMONTHS(t *testing.T) {
	opts := []expr.Option{
		expr.AllowUndefinedVariables(),
		expr.AsAny(),
		expr.WithDateTimeFunctions(),
	}
	d := time.Date(2023, 1, 31, 0, 0, 0, 0, time.UTC)
	env := map[string]any{"d": d}
	program, err := expr.Compile("ADDMONTHS(d, 1)", append(opts, expr.Env(env))...)
	require.NoError(t, err)
	got, err := expr.Run(program, env)
	require.NoError(t, err)
	// Jan 31 + 1 month = Feb 28 (2023 is not a leap year)
	expected := time.Date(2023, 2, 28, 0, 0, 0, 0, time.UTC)
	assert.Equal(t, expected, got)
}

func TestPack_DateTime_DATEVALUE(t *testing.T) {
	opts := []expr.Option{
		expr.AllowUndefinedVariables(),
		expr.AsAny(),
		expr.WithDateTimeFunctions(),
	}
	program, err := expr.Compile(`DATEVALUE("2023-06-15")`, opts...)
	require.NoError(t, err)
	got, err := expr.Run(program, nil)
	require.NoError(t, err)
	expected := time.Date(2023, 6, 15, 0, 0, 0, 0, time.UTC)
	assert.Equal(t, expected, got)
}

func TestPack_DateTime_FROMUNIXTIME(t *testing.T) {
	opts := []expr.Option{
		expr.AllowUndefinedVariables(),
		expr.AsAny(),
		expr.WithDateTimeFunctions(),
	}
	program, err := expr.Compile("FROMUNIXTIME(0)", opts...)
	require.NoError(t, err)
	got, err := expr.Run(program, nil)
	require.NoError(t, err)
	expected := time.Unix(0, 0).UTC()
	assert.Equal(t, expected, got)
}

func TestPack_DateTime_UNIXTIMESTAMP(t *testing.T) {
	opts := []expr.Option{
		expr.AllowUndefinedVariables(),
		expr.AsAny(),
		expr.WithDateTimeFunctions(),
	}
	d := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
	env := map[string]any{"d": d}
	program, err := expr.Compile("UNIXTIMESTAMP(d)", append(opts, expr.Env(env))...)
	require.NoError(t, err)
	got, err := expr.Run(program, env)
	require.NoError(t, err)
	assert.Equal(t, d.Unix(), got)
}

func TestPack_DateTime_ContextFunctions(t *testing.T) {
	tz := time.FixedZone("UTC+5", 5*3600)
	env := map[string]any{"_tz": tz}

	opts := []expr.Option{
		expr.AllowUndefinedVariables(),
		expr.AsAny(),
		expr.WithDateTimeFunctions(),
	}

	// TODAY() should return a date (midnight) in the configured timezone
	program, err := expr.Compile("TODAY()", opts...)
	require.NoError(t, err)
	got, err := expr.Run(program, env)
	require.NoError(t, err)
	today := got.(time.Time)
	assert.Equal(t, 0, today.Hour())
	assert.Equal(t, 0, today.Minute())
	assert.Equal(t, 0, today.Second())

	// NOW() should return a datetime
	program2, err := expr.Compile("NOW()", opts...)
	require.NoError(t, err)
	got2, err := expr.Run(program2, env)
	require.NoError(t, err)
	now := got2.(time.Time)
	assert.False(t, now.IsZero())
}

func TestPack_DateTime_FORMATDURATION(t *testing.T) {
	opts := []expr.Option{
		expr.AllowUndefinedVariables(),
		expr.AsAny(),
		expr.WithDateTimeFunctions(),
	}
	d1 := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
	d2 := time.Date(2023, 1, 2, 3, 30, 45, 0, time.UTC)
	env := map[string]any{"d1": d1, "d2": d2}
	program, err := expr.Compile("FORMATDURATION(d1, d2)", append(opts, expr.Env(env))...)
	require.NoError(t, err)
	got, err := expr.Run(program, env)
	require.NoError(t, err)
	assert.Equal(t, "1:03:30:45", got)
}

func TestPack_DateTime_ISOWEEK_ISOYEAR(t *testing.T) {
	opts := []expr.Option{
		expr.AllowUndefinedVariables(),
		expr.AsAny(),
		expr.WithDateTimeFunctions(),
	}
	// Jan 1, 2023 is a Sunday — ISO week 52 of 2022
	d := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
	env := map[string]any{"d": d}

	program, err := expr.Compile("ISOWEEK(d)", append(opts, expr.Env(env))...)
	require.NoError(t, err)
	got, err := expr.Run(program, env)
	require.NoError(t, err)
	isoYear, isoWeek := d.ISOWeek()
	assert.Equal(t, isoWeek, got)

	program2, err := expr.Compile("ISOYEAR(d)", append(opts, expr.Env(env))...)
	require.NoError(t, err)
	got2, err := expr.Run(program2, env)
	require.NoError(t, err)
	assert.Equal(t, isoYear, got2)
}

// =============================================================================
// Regex function pack
// =============================================================================

func TestPack_Regex(t *testing.T) {
	opts := []expr.Option{
		expr.AllowUndefinedVariables(),
		expr.AsAny(),
		expr.WithRegexFunctions(),
	}

	tests := []struct {
		expr string
		want any
	}{
		{`REGEX("Hello123", "[0-9]+")`, true},
		{`REGEX("Hello", "[0-9]+")`, false},
		{`REGEX("abc@example.com", "^[a-z]+@[a-z]+\\.[a-z]+$")`, true},
		{`REGEX(NULL, "[0-9]+")`, false},
		{`REGEX("test", NULL)`, false},
	}
	for _, tt := range tests {
		t.Run(tt.expr, func(t *testing.T) {
			program, err := expr.Compile(tt.expr, opts...)
			require.NoError(t, err)
			got, err := expr.Run(program, nil)
			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

// =============================================================================
// Encoding function pack
// =============================================================================

func TestPack_Encoding(t *testing.T) {
	opts := []expr.Option{
		expr.AllowUndefinedVariables(),
		expr.AsAny(),
		expr.WithEncodingFunctions(),
	}

	tests := []struct {
		expr string
		want any
	}{
		{`HTMLENCODE("<b>bold</b>")`, "&lt;b&gt;bold&lt;/b&gt;"},
		{"HTMLENCODE(NULL)", nil},
		{`URLENCODE("hello world")`, "hello+world"},
		{"URLENCODE(NULL)", nil},
		{`JSENCODE("He said 'hi'")`, `He said \'hi\'`},
		{"JSENCODE(NULL)", nil},
	}
	for _, tt := range tests {
		t.Run(tt.expr, func(t *testing.T) {
			program, err := expr.Compile(tt.expr, opts...)
			require.NoError(t, err)
			got, err := expr.Run(program, nil)
			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

// =============================================================================
// WithAllFormulaPacks
// =============================================================================

func TestWithAllFormulaPacks(t *testing.T) {
	opts := []expr.Option{
		expr.AllowUndefinedVariables(),
		expr.AsAny(),
		expr.WithAllFormulaPacks(),
	}

	// Test that functions from multiple packs work together
	tests := []struct {
		expr string
		want any
	}{
		{"ABS(-5)", 5.0},             // math
		{"AND(true, true)", true},    // logical
		{"SIN(0)", 0.0},              // trig
		{`UPPER("hello")`, "HELLO"},  // text
		{`REGEX("abc", "^a")`, true}, // regex
	}
	for _, tt := range tests {
		t.Run(tt.expr, func(t *testing.T) {
			program, err := expr.Compile(tt.expr, opts...)
			require.NoError(t, err)
			got, err := expr.Run(program, nil)
			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

// =============================================================================
// Combined keyword + operator + function tests
// =============================================================================

func TestFormula_Combined(t *testing.T) {
	opts := []expr.Option{
		expr.AllowUndefinedVariables(),
		expr.AsAny(),
		expr.WithAllFormulaPacks(),
	}

	tests := []struct {
		expr string
		want any
	}{
		// TRUE/FALSE with logical functions
		{"IF(TRUE, 1, 2)", 1},
		{"IF(FALSE, 1, 2)", 2},
		{"AND(TRUE, TRUE)", true},
		{"OR(FALSE, TRUE)", true},
		{"NOT(FALSE)", true},
		// NULL with ISBLANK
		{"ISBLANK(NULL)", true},
		// <> with IF
		{"IF(1 <> 2, \"diff\", \"same\")", "diff"},
		{"IF(1 <> 1, \"diff\", \"same\")", "same"},
		// & concatenation with functions
		{`UPPER("hello") & " " & UPPER("world")`, "HELLO WORLD"},
		{`"Result: " & ABS(-42)`, "Result: 42"},
	}
	for _, tt := range tests {
		t.Run(tt.expr, func(t *testing.T) {
			program, err := expr.Compile(tt.expr, opts...)
			require.NoError(t, err)
			got, err := expr.Run(program, nil)
			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

// =============================================================================
// Edge cases
// =============================================================================

func TestFormula_NullPropagation(t *testing.T) {
	opts := []expr.Option{
		expr.AllowUndefinedVariables(),
		expr.AsAny(),
		expr.WithAllFormulaPacks(),
	}

	// Most functions should propagate null
	tests := []struct {
		expr string
	}{
		{"ABS(NULL)"},
		{"CEILING(NULL)"},
		{"FLOOR(NULL)"},
		{"SQRT(NULL)"},
		{"LOWER(NULL)"},
		{"UPPER(NULL)"},
		{"TRIM(NULL)"},
		{"REVERSE(NULL)"},
		{"LEFT(NULL, 3)"},
		{"RIGHT(NULL, 3)"},
		{"DAY(NULL)"},
		{"MONTH(NULL)"},
		{"YEAR(NULL)"},
	}
	for _, tt := range tests {
		t.Run(tt.expr, func(t *testing.T) {
			program, err := expr.Compile(tt.expr, opts...)
			require.NoError(t, err)
			got, err := expr.Run(program, nil)
			require.NoError(t, err)
			assert.Nil(t, got, "expected nil for %s", tt.expr)
		})
	}
}
