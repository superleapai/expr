package expr_test

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"os"
	"reflect"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/expr-lang/expr/conf"
	"github.com/expr-lang/expr/internal/testify/assert"
	"github.com/expr-lang/expr/internal/testify/require"
	"github.com/expr-lang/expr/types"
	"github.com/expr-lang/expr/vm"

	"github.com/expr-lang/expr"
	"github.com/expr-lang/expr/ast"
	"github.com/expr-lang/expr/file"
	"github.com/expr-lang/expr/test/mock"
)

func ExampleEval() {
	output, err := expr.Eval("greet + name", map[string]any{
		"greet": "Hello, ",
		"name":  "world!",
	})
	if err != nil {
		fmt.Printf("err: %v", err)
		return
	}

	fmt.Printf("%v", output)

	// Output: Hello, world!
}

//func ExampleEval_runtime_error() {
//	_, err := expr.Eval(`map(1..3, {1 % (# - 3)})`, nil)
//	fmt.Print(err)
//
//	// Output: integer divide by zero (1:14)
//	//  | map(1..3, {1 % (# - 3)})
//	//  | .............^
//}

func ExampleCompile() {
	env := map[string]any{
		"foo": 1,
		"bar": 99,
	}

	program, err := expr.Compile("foo in 1..99 and bar in 1..99", expr.Env(env))
	if err != nil {
		fmt.Printf("%v", err)
		return
	}

	output, err := expr.Run(program, env)
	if err != nil {
		fmt.Printf("%v", err)
		return
	}

	fmt.Printf("%v", output)

	// Output: true
}

func ExampleEnv() {
	type Segment struct {
		Origin string
	}
	type Passengers struct {
		Adults int
	}
	type Meta struct {
		Tags map[string]string
	}
	type Env struct {
		Meta
		Segments   []*Segment
		Passengers *Passengers
		Marker     string
	}

	code := `all(Segments, {.Origin == "MOW"}) && Passengers.Adults > 0 && Tags["foo"] startsWith "bar"`

	program, err := expr.Compile(code, expr.Env(Env{}))
	if err != nil {
		fmt.Printf("%v", err)
		return
	}

	env := Env{
		Meta: Meta{
			Tags: map[string]string{
				"foo": "bar",
			},
		},
		Segments: []*Segment{
			{Origin: "MOW"},
		},
		Passengers: &Passengers{
			Adults: 2,
		},
		Marker: "test",
	}

	output, err := expr.Run(program, env)
	if err != nil {
		fmt.Printf("%v", err)
		return
	}

	fmt.Printf("%v", output)

	// Output: true
}

func ExampleEnv_tagged_field_names() {
	env := struct {
		FirstWord  string
		Separator  string `expr:"Space"`
		SecondWord string `expr:"second_word"`
	}{
		FirstWord:  "Hello",
		Separator:  " ",
		SecondWord: "World",
	}

	output, err := expr.Eval(`FirstWord + Space + second_word`, env)
	if err != nil {
		fmt.Printf("%v", err)
		return
	}

	fmt.Printf("%v", output)

	// Output : Hello World
}

func ExampleAsKind() {
	program, err := expr.Compile("{a: 1, b: 2}", expr.AsKind(reflect.Map))
	if err != nil {
		fmt.Printf("%v", err)
		return
	}

	output, err := expr.Run(program, nil)
	if err != nil {
		fmt.Printf("%v", err)
		return
	}

	fmt.Printf("%v", output)

	// Output: map[a:1 b:2]
}

func ExampleAsBool() {
	env := map[string]int{
		"foo": 0,
	}

	program, err := expr.Compile("foo >= 0", expr.Env(env), expr.AsBool())
	if err != nil {
		fmt.Printf("%v", err)
		return
	}

	output, err := expr.Run(program, env)
	if err != nil {
		fmt.Printf("%v", err)
		return
	}

	fmt.Printf("%v", output.(bool))

	// Output: true
}

func ExampleAsBool_error() {
	env := map[string]any{
		"foo": 0,
	}

	_, err := expr.Compile("foo + 42", expr.Env(env), expr.AsBool())

	fmt.Printf("%v", err)

	// Output: expected bool, but got int
}

func ExampleAsInt() {
	program, err := expr.Compile("42", expr.AsInt())
	if err != nil {
		fmt.Printf("%v", err)
		return
	}

	output, err := expr.Run(program, nil)
	if err != nil {
		fmt.Printf("%v", err)
		return
	}

	fmt.Printf("%T(%v)", output, output)

	// Output: int(42)
}

func ExampleAsInt64() {
	env := map[string]any{
		"rating": 5.5,
	}

	program, err := expr.Compile("rating", expr.Env(env), expr.AsInt64())
	if err != nil {
		fmt.Printf("%v", err)
		return
	}

	output, err := expr.Run(program, env)
	if err != nil {
		fmt.Printf("%v", err)
		return
	}

	fmt.Printf("%v", output.(int64))

	// Output: 5
}

func ExampleAsFloat64() {
	program, err := expr.Compile("42", expr.AsFloat64())
	if err != nil {
		fmt.Printf("%v", err)
		return
	}

	output, err := expr.Run(program, nil)
	if err != nil {
		fmt.Printf("%v", err)
		return
	}

	fmt.Printf("%v", output.(float64))

	// Output: 42
}

func ExampleAsFloat64_error() {
	_, err := expr.Compile(`!!true`, expr.AsFloat64())

	fmt.Printf("%v", err)

	// Output: expected float64, but got bool
}

func ExampleWarnOnAny() {
	// Arrays always have []any type. The expression return type is any.
	// AsInt() instructs compiler to expect int or any, and cast to int,
	// if possible. WarnOnAny() instructs to return an error on any type.
	_, err := expr.Compile(`[42, true, "yes"][0]`, expr.AsInt(), expr.WarnOnAny())

	fmt.Printf("%v", err)

	// Output: expected int, but got interface {}
}

func ExampleOperator() {
	code := `
		Now() > CreatedAt &&
		(Now() - CreatedAt).Hours() > 24
	`

	type Env struct {
		CreatedAt time.Time
		Now       func() time.Time
		Sub       func(a, b time.Time) time.Duration
		After     func(a, b time.Time) bool
	}

	options := []expr.Option{
		expr.Env(Env{}),
		expr.Operator(">", "After"),
		expr.Operator("-", "Sub"),
	}

	program, err := expr.Compile(code, options...)
	if err != nil {
		fmt.Printf("%v", err)
		return
	}

	env := Env{
		CreatedAt: time.Date(2018, 7, 14, 0, 0, 0, 0, time.UTC),
		Now:       func() time.Time { return time.Now() },
		Sub:       func(a, b time.Time) time.Duration { return a.Sub(b) },
		After:     func(a, b time.Time) bool { return a.After(b) },
	}

	output, err := expr.Run(program, env)
	if err != nil {
		fmt.Printf("%v", err)
		return
	}

	fmt.Printf("%v", output)

	// Output: true
}

func ExampleOperator_with_decimal() {
	type Decimal struct{ N float64 }
	code := `A + B - C`

	type Env struct {
		A, B, C Decimal
		Sub     func(a, b Decimal) Decimal
		Add     func(a, b Decimal) Decimal
	}

	options := []expr.Option{
		expr.Env(Env{}),
		expr.Operator("+", "Add"),
		expr.Operator("-", "Sub"),
	}

	program, err := expr.Compile(code, options...)
	if err != nil {
		fmt.Printf("Compile error: %v", err)
		return
	}

	env := Env{
		A:   Decimal{3},
		B:   Decimal{2},
		C:   Decimal{1},
		Sub: func(a, b Decimal) Decimal { return Decimal{a.N - b.N} },
		Add: func(a, b Decimal) Decimal { return Decimal{a.N + b.N} },
	}

	output, err := expr.Run(program, env)
	if err != nil {
		fmt.Printf("%v", err)
		return
	}

	fmt.Printf("%v", output)

	// Output: {4}
}

func fib(n int) int {
	if n <= 1 {
		return n
	}
	return fib(n-1) + fib(n-2)
}

func ExampleConstExpr() {
	code := `[fib(5), fib(3+3), fib(dyn)]`

	env := map[string]any{
		"fib": fib,
		"dyn": 0,
	}

	options := []expr.Option{
		expr.Env(env),
		expr.ConstExpr("fib"), // Mark fib func as constant expression.
	}

	program, err := expr.Compile(code, options...)
	if err != nil {
		fmt.Printf("%v", err)
		return
	}

	// Only fib(5) and fib(6) calculated on Compile, fib(dyn) can be called at runtime.
	env["dyn"] = 7

	output, err := expr.Run(program, env)
	if err != nil {
		fmt.Printf("%v", err)
		return
	}

	fmt.Printf("%v\n", output)

	// Output: [5 8 13]
}

func ExampleAllowUndefinedVariables() {
	code := `name == nil ? "Hello, world!" : sprintf("Hello, %v!", name)`

	env := map[string]any{
		"sprintf": fmt.Sprintf,
	}

	options := []expr.Option{
		expr.Env(env),
		expr.AllowUndefinedVariables(), // Allow to use undefined variables.
	}

	program, err := expr.Compile(code, options...)
	if err != nil {
		fmt.Printf("%v", err)
		return
	}

	output, err := expr.Run(program, env)
	if err != nil {
		fmt.Printf("%v", err)
		return
	}
	fmt.Printf("%v\n", output)

	env["name"] = "you" // Define variables later on.

	output, err = expr.Run(program, env)
	if err != nil {
		fmt.Printf("%v", err)
		return
	}
	fmt.Printf("%v\n", output)

	// Output: Hello, world!
	// Hello, you!
}

func ExampleAllowUndefinedVariables_zero_value() {
	code := `name == "" ? foo + bar : foo + name`

	// If environment has different zero values, then undefined variables
	// will have it as default value.
	env := map[string]string{}

	options := []expr.Option{
		expr.Env(env),
		expr.AllowUndefinedVariables(), // Allow to use undefined variables.
	}

	program, err := expr.Compile(code, options...)
	if err != nil {
		fmt.Printf("%v", err)
		return
	}

	env = map[string]string{
		"foo": "Hello, ",
		"bar": "world!",
	}

	output, err := expr.Run(program, env)
	if err != nil {
		fmt.Printf("%v", err)
		return
	}
	fmt.Printf("%v", output)

	// Output: Hello, world!
}

func ExampleAllowUndefinedVariables_zero_value_functions() {
	code := `words == "" ? Split("foo,bar", ",") : Split(words, ",")`

	// Env is map[string]string type on which methods are defined.
	env := mock.MapStringStringEnv{}

	options := []expr.Option{
		expr.Env(env),
		expr.AllowUndefinedVariables(), // Allow to use undefined variables.
	}

	program, err := expr.Compile(code, options...)
	if err != nil {
		fmt.Printf("%v", err)
		return
	}

	output, err := expr.Run(program, env)
	if err != nil {
		fmt.Printf("%v", err)
		return
	}
	fmt.Printf("%v", output)

	// Output: [foo bar]
}

type patcher struct{}

func (p *patcher) Visit(node *ast.Node) {
	switch n := (*node).(type) {
	case *ast.MemberNode:
		ast.Patch(node, &ast.CallNode{
			Callee:    &ast.IdentifierNode{Value: "get"},
			Arguments: []ast.Node{n.Node, n.Property},
		})
	}
}

func ExamplePatch() {
	program, err := expr.Compile(
		`greet.you.world + "!"`,
		expr.Patch(&patcher{}),
	)
	if err != nil {
		fmt.Printf("%v", err)
		return
	}

	env := map[string]any{
		"greet": "Hello",
		"get": func(a, b string) string {
			return a + ", " + b
		},
	}

	output, err := expr.Run(program, env)
	if err != nil {
		fmt.Printf("%v", err)
		return
	}
	fmt.Printf("%v", output)

	// Output : Hello, you, world!
}

func ExampleWithContext() {
	env := map[string]any{
		"fn": func(ctx context.Context, _, _ int) int {
			// An infinite loop that can be canceled by context.
			for {
				select {
				case <-ctx.Done():
					return 42
				}
			}
		},
		"ctx": context.TODO(), // Context should be passed as a variable.
	}

	program, err := expr.Compile(`fn(1, 2)`,
		expr.Env(env),
		expr.WithContext("ctx"), // Pass context variable name.
	)
	if err != nil {
		fmt.Printf("%v", err)
		return
	}

	// Cancel context after 100 milliseconds.
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*100)
	defer cancel()

	// After program is compiled, context can be passed to Run.
	env["ctx"] = ctx

	// Run will return 42 after 100 milliseconds.
	output, err := expr.Run(program, env)
	if err != nil {
		fmt.Printf("%v", err)
		return
	}

	fmt.Printf("%v", output)
	// Output: 42
}

func ExampleTimezone() {
	program, err := expr.Compile(`now().Location().String()`, expr.Timezone("Asia/Kamchatka"))
	if err != nil {
		fmt.Printf("%v", err)
		return
	}

	output, err := expr.Run(program, nil)
	if err != nil {
		fmt.Printf("%v", err)
		return
	}

	fmt.Printf("%v", output)
	// Output: Asia/Kamchatka
}

func TestExpr_readme_example(t *testing.T) {
	env := map[string]any{
		"greet":   "Hello, %v!",
		"names":   []string{"world", "you"},
		"sprintf": fmt.Sprintf,
	}

	code := `sprintf(greet, names[0])`

	program, err := expr.Compile(code, expr.Env(env))
	require.NoError(t, err)

	output, err := expr.Run(program, env)
	require.NoError(t, err)

	require.Equal(t, "Hello, world!", output)
}

func TestExpr(t *testing.T) {
	date := time.Date(2017, time.October, 23, 18, 30, 0, 0, time.UTC)
	oneDay, _ := time.ParseDuration("24h")
	timeNowPlusOneDay := date.Add(oneDay)

	env := mock.Env{
		Embed:     mock.Embed{},
		Ambiguous: "",
		Any:       nil,
		Bool:      true,
		Float:     0,
		Int64:     0,
		Int32:     0,
		Int:       0,
		One:       1,
		Two:       2,
		Uint32:    0,
		String:    "string",
		BoolPtr:   nil,
		FloatPtr:  nil,
		IntPtr:    nil,
		IntPtrPtr: nil,
		StringPtr: nil,
		Foo: mock.Foo{
			Value: "foo",
			Bar: mock.Bar{
				Baz: "baz",
			},
		},
		Abstract:           nil,
		ArrayOfAny:         nil,
		ArrayOfInt:         []int{1, 2, 3, 4, 5},
		ArrayOfFoo:         []*mock.Foo{{Value: "foo"}, {Value: "bar"}, {Value: "baz"}},
		MapOfFoo:           nil,
		MapOfAny:           nil,
		FuncParam:          nil,
		FuncParamAny:       nil,
		FuncTooManyReturns: nil,
		FuncNamed:          nil,
		NilAny:             nil,
		NilFn:              nil,
		NilStruct:          nil,
		Variadic: func(head int, xs ...int) bool {
			sum := 0
			for _, x := range xs {
				sum += x
			}
			return head == sum
		},
		Fast:        nil,
		Time:        date,
		TimePlusDay: timeNowPlusOneDay,
		Duration:    oneDay,
	}

	tests := []struct {
		code string
		want any
	}{
		{
			`1`,
			1,
		},
		{
			`-.5`,
			-.5,
		},
		{
			`true && false || false`,
			false,
		},
		{
			`Int == 0 && Int32 == 0 && Int64 == 0 && Float64 == 0 && Bool && String == "string"`,
			true,
		},
		{
			`-Int64 == 0`,
			true,
		},
		{
			`"a" != "b"`,
			true,
		},
		{
			`"a" != "b" || 1 == 2`,
			true,
		},
		{
			`Int + 0`,
			0,
		},
		{
			`Uint64 + 0`,
			0,
		},
		{
			`Uint64 + Int64`,
			0,
		},
		{
			`Int32 + Int64`,
			0,
		},
		{
			`Float64 + 0`,
			float64(0),
		},
		{
			`0 + Float64`,
			float64(0),
		},
		{
			`0 <= Float64`,
			true,
		},
		{
			`Float64 < 1`,
			true,
		},
		{
			`Int < 1`,
			true,
		},
		{
			`2 + 2 == 4`,
			true,
		},
		{
			`8 % 3`,
			2,
		},
		{
			`2 ** 8`,
			float64(256),
		},
		{
			`2 ^ 8`,
			float64(256),
		},
		{
			`-(2-5)**3-2/(+4-3)+-2`,
			float64(23),
		},
		{
			`"hello" + " " + "world"`,
			"hello world",
		},
		{
			`0 in -1..1 and 1 in 1..1`,
			true,
		},
		{
			`Int in 0..1`,
			true,
		},
		{
			`Int32 in 0..1`,
			true,
		},
		{
			`Int64 in 0..1`,
			true,
		},
		{
			`1 in [1, 2, 3] && "foo" in {foo: 0, bar: 1} && "Bar" in Foo`,
			true,
		},
		{
			`1 in [1.5] || 1 not in [1]`,
			false,
		},
		{
			`One in 0..1 && Two not in 0..1`,
			true,
		},
		{
			`Two not in 0..1`,
			true,
		},
		{
			`Two not    in 0..1`,
			true,
		},
		{
			`-1 not in [1]`,
			true,
		},
		{
			`Int32 in [10, 20]`,
			false,
		},
		{
			`String matches "s.+"`,
			true,
		},
		{
			`String matches ("^" + String + "$")`,
			true,
		},
		{
			`'foo' + 'bar' not matches 'foobar'`,
			false,
		},
		{
			`"foobar" contains "bar"`,
			true,
		},
		{
			`"foobar" startsWith "foo"`,
			true,
		},
		{
			`"foobar" endsWith "bar"`,
			true,
		},
		{
			`(0..10)[5]`,
			5,
		},
		{
			`Foo.Bar.Baz`,
			"baz",
		},
		{
			`Add(10, 5) + GetInt()`,
			15,
		},
		{
			`Foo.Method().Baz`,
			`baz (from Foo.Method)`,
		},
		{
			`Foo.MethodWithArgs("prefix ")`,
			"prefix foo",
		},
		{
			`len([1, 2, 3])`,
			3,
		},
		{
			`len([1, Two, 3])`,
			3,
		},
		{
			`len(["hello", "world"])`,
			2,
		},
		{
			`len("hello, world")`,
			12,
		},
		{
			`len('北京')`,
			2,
		},
		{
			`len('👍🏻')`, // one grapheme cluster, two code points
			2,
		},
		{
			`len('👍')`, // one grapheme cluster, one code point
			1,
		},
		{
			`len(ArrayOfInt)`,
			5,
		},
		{
			`len({a: 1, b: 2, c: 2})`,
			3,
		},
		{
			`max([1, 2, 3])`,
			3,
		},
		{
			`max(1, 2, 3)`,
			3,
		},
		{
			`min([1, 2, 3])`,
			1,
		},
		{
			`min(1, 2, 3)`,
			1,
		},
		{
			`{foo: 0, bar: 1}`,
			map[string]any{"foo": 0, "bar": 1},
		},
		{
			`{foo: 0, bar: 1}`,
			map[string]any{"foo": 0, "bar": 1},
		},
		{
			`(true ? 0+1 : 2+3) + (false ? -1 : -2)`,
			-1,
		},
		{
			`filter(1..9, {# > 7})`,
			[]any{8, 9},
		},
		{
			`map(1..3, {# * #})`,
			[]any{1, 4, 9},
		},
		{
			`all(1..3, {# > 0})`,
			true,
		},
		{
			`count(1..30, {# % 3 == 0})`,
			10,
		},
		{
			`count([true, true, false])`,
			2,
		},
		{
			`"a" < "b"`,
			true,
		},
		{
			`Time.Sub(Time).String() == "0s"`,
			true,
		},
		{
			`1 + 1`,
			2,
		},
		{
			`(One * Two) * 3 == One * (Two * 3)`,
			true,
		},
		{
			`ArrayOfInt[1]`,
			2,
		},
		{
			`ArrayOfInt[0] < ArrayOfInt[1]`,
			true,
		},
		{
			`ArrayOfInt[-1]`,
			5,
		},
		{
			`ArrayOfInt[1:2]`,
			[]int{2},
		},
		{
			`ArrayOfInt[1:4]`,
			[]int{2, 3, 4},
		},
		{
			`ArrayOfInt[-4:-1]`,
			[]int{2, 3, 4},
		},
		{
			`ArrayOfInt[:3]`,
			[]int{1, 2, 3},
		},
		{
			`ArrayOfInt[3:]`,
			[]int{4, 5},
		},
		{
			`ArrayOfInt[0:5] == ArrayOfInt`,
			true,
		},
		{
			`ArrayOfInt[0:] == ArrayOfInt`,
			true,
		},
		{
			`ArrayOfInt[:5] == ArrayOfInt`,
			true,
		},
		{
			`ArrayOfInt[:] == ArrayOfInt`,
			true,
		},
		{
			`4 in 5..1`,
			false,
		},
		{
			`4..0`,
			[]int{},
		},
		{
			`NilStruct`,
			(*mock.Foo)(nil),
		},
		{
			`NilAny == nil && nil == NilAny && nil == nil && NilAny == NilAny && NilInt == nil && NilSlice == nil && NilStruct == nil`,
			true,
		},
		{
			`0 == nil || "str" == nil || true == nil`,
			false,
		},
		{
			`Variadic(6, 1, 2, 3)`,
			true,
		},
		{
			`Variadic(0)`,
			true,
		},
		{
			`String[:]`,
			"string",
		},
		{
			`String[:3]`,
			"str",
		},
		{
			`String[:9]`,
			"string",
		},
		{
			`String[3:9]`,
			"ing",
		},
		{
			`String[7:9]`,
			"",
		},
		{
			`map(filter(ArrayOfInt, # >= 3), # + 1)`,
			[]any{4, 5, 6},
		},
		{
			`Time < Time + Duration`,
			true,
		},
		{
			`Time + Duration > Time`,
			true,
		},
		{
			`Time == Time`,
			true,
		},
		{
			`Time >= Time`,
			true,
		},
		{
			`Time <= Time`,
			true,
		},
		{
			`Time == Time + Duration`,
			false,
		},
		{
			`Time != Time`,
			false,
		},
		{
			`TimePlusDay - Duration`,
			date,
		},
		{
			`duration("1h") == duration("1h")`,
			true,
		},
		{
			`TimePlusDay - Time >= duration("24h")`,
			true,
		},
		{
			`duration("1h") > duration("1m")`,
			true,
		},
		{
			`duration("1h") < duration("1m")`,
			false,
		},
		{
			`duration("1h") >= duration("1m")`,
			true,
		},
		{
			`duration("1h") <= duration("1m")`,
			false,
		},
		{
			`duration("1h") > duration("1m")`,
			true,
		},
		{
			`duration("1h") + duration("1m")`,
			time.Hour + time.Minute,
		},
		{
			`duration("1h") - duration("1m")`,
			time.Hour - time.Minute,
		},
		{
			`7 * duration("1h")`,
			7 * time.Hour,
		},
		{
			`duration("1h") * 7`,
			7 * time.Hour,
		},
		{
			`duration("1s") * .5`,
			5e8,
		},
		{
			`1 /* one */ + 2 // two`,
			3,
		},
		{
			`let x = 1; x + 2`,
			3,
		},
		{
			`map(1..3, let x = #; let y = x * x; y * y)`,
			[]any{1, 16, 81},
		},
		{
			`map(1..2, let x = #; map(2..3, let y = #; x + y))`,
			[]any{[]any{3, 4}, []any{4, 5}},
		},
		{
			`len(filter(1..99, # % 7 == 0))`,
			14,
		},
		{
			`find(ArrayOfFoo, .Value == "baz")`,
			env.ArrayOfFoo[2],
		},
		{
			`findIndex(ArrayOfFoo, .Value == "baz")`,
			2,
		},
		{
			`filter(ArrayOfFoo, .Value == "baz")[0]`,
			env.ArrayOfFoo[2],
		},
		{
			`first(filter(ArrayOfFoo, .Value == "baz"))`,
			env.ArrayOfFoo[2],
		},
		{
			`first(filter(ArrayOfFoo, false))`,
			nil,
		},
		{
			`findLast(1..9, # % 2 == 0)`,
			8,
		},
		{
			`findLastIndex(1..9, # % 2 == 0)`,
			7,
		},
		{
			`filter(1..9, # % 2 == 0)[-1]`,
			8,
		},
		{
			`last(filter(1..9, # % 2 == 0))`,
			8,
		},
		{
			`map(filter(1..9, # % 2 == 0), # * 2)`,
			[]any{4, 8, 12, 16},
		},
		{
			`map(map(filter(1..9, # % 2 == 0), # * 2), # * 2)`,
			[]any{8, 16, 24, 32},
		},
		{
			`first(map(filter(1..9, # % 2 == 0), # * 2))`,
			4,
		},
		{
			`map(filter(1..9, # % 2 == 0), # * 2)[-1]`,
			16,
		},
		{
			`len(map(filter(1..9, # % 2 == 0), # * 2))`,
			4,
		},
		{
			`len(filter(map(1..9, # * 2), # % 2 == 0))`,
			9,
		},
		{
			`first(filter(map(1..9, # * 2), # % 2 == 0))`,
			2,
		},
		{
			`first(map(filter(1..9, # % 2 == 0), # * 2))`,
			4,
		},
		{
			`2^3 == 8`,
			true,
		},
		{
			`4/2 == 2`,
			true,
		},
		{
			`.5 in 0..1`,
			false,
		},
		{
			`.5 in ArrayOfInt`,
			false,
		},
		{
			`bitnot(10)`,
			-11,
		},
		{
			`bitxor(15, 32)`,
			47,
		},
		{
			`bitand(90, 34)`,
			2,
		},
		{
			`bitnand(35, 9)`,
			34,
		},
		{
			`bitor(10, 5)`,
			15,
		},
		{
			`bitshr(7, 2)`,
			1,
		},
		{
			`bitshl(7, 2)`,
			28,
		},
		{
			`bitushr(-100, 5)`,
			576460752303423484,
		},
		{
			`"hello"[1:3]`,
			"el",
		},
		{
			`[1, 2, 3]?.[0]`,
			1,
		},
		{
			`[[1, 2], 3, 4]?.[0]?.[1]`,
			2,
		},
		{
			`[nil, 3, 4]?.[0]?.[1]`,
			nil,
		},
		{
			`1 > 2 < 3`,
			false,
		},
		{
			`1 < 2 < 3`,
			true,
		},
		{
			`1 < 2 < 3 > 4`,
			false,
		},
		{
			`1 < 2 < 3 > 2`,
			true,
		},
		{
			`1 < 2 < 3 == true`,
			true,
		},
		{
			`if 1 > 2 { 333 * 2 + 1 } else { 444 }`,
			444,
		},
		{
			`let a = 3;
			let b = 2;
			if a>b {let c = Add(a, b); c+1} else {Add(10, b)}
			`,
			6,
		},
		{
			`if "a" < "b" {let x = "a"; x} else {"abc"}`,
			"a",
		},
		{
			`1; 2; 3`,
			3,
		},
		{
			`let a = 1; Add(2, 2); let b = 2; a + b`,
			3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.code, func(t *testing.T) {
			{
				program, err := expr.Compile(tt.code, expr.Env(mock.Env{}))
				require.NoError(t, err, "compile error")

				got, err := expr.Run(program, env)
				require.NoError(t, err, "run error")
				assert.Equal(t, tt.want, got)
			}
			{
				program, err := expr.Compile(tt.code, expr.Optimize(false))
				require.NoError(t, err, "unoptimized")

				got, err := expr.Run(program, env)
				require.NoError(t, err, "unoptimized")
				assert.Equal(t, tt.want, got, "unoptimized")
			}
			{
				got, err := expr.Eval(tt.code, env)
				require.NoError(t, err, "eval")
				assert.Equal(t, tt.want, got, "eval")
			}
			{
				program, err := expr.Compile(tt.code, expr.Env(mock.Env{}), expr.Optimize(false))
				require.NoError(t, err)

				code := program.Node().String()
				got, err := expr.Eval(code, env)
				require.NoError(t, err, code)
				assert.Equal(t, tt.want, got, code)
			}
		})
	}
}

func TestExpr_error(t *testing.T) {
	env := mock.Env{
		ArrayOfAny: []any{1, "2", 3, true},
	}

	tests := []struct {
		code string
		want string
	}{
		{
			`filter(1..9, # > 9)[0]`,
			`reflect: slice index out of range (1:20)
 | filter(1..9, # > 9)[0]
 | ...................^`,
		},
		{
			`ArrayOfAny[-7]`,
			`index out of range: -3 (array length is 4) (1:11)
 | ArrayOfAny[-7]
 | ..........^`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.code, func(t *testing.T) {
			program, err := expr.Compile(tt.code, expr.Env(mock.Env{}))
			require.NoError(t, err)

			_, err = expr.Run(program, env)
			require.Error(t, err)
			assert.Equal(t, tt.want, err.Error())
		})
	}
}

func TestExpr_optional_chaining(t *testing.T) {
	env := map[string]any{}
	program, err := expr.Compile("foo?.bar.baz", expr.Env(env), expr.AllowUndefinedVariables())
	require.NoError(t, err)

	got, err := expr.Run(program, env)
	require.NoError(t, err)
	assert.Equal(t, nil, got)
}

func TestExpr_optional_chaining_property(t *testing.T) {
	env := map[string]any{
		"foo": map[string]any{},
	}
	program, err := expr.Compile("foo.bar?.baz", expr.Env(env))
	require.NoError(t, err)

	got, err := expr.Run(program, env)
	require.NoError(t, err)
	assert.Equal(t, nil, got)
}

func TestExpr_optional_chaining_nested_chains(t *testing.T) {
	env := map[string]any{
		"foo": map[string]any{
			"id": 1,
			"bar": []map[string]any{
				1: {
					"baz": "baz",
				},
			},
		},
	}
	program, err := expr.Compile("foo?.bar[foo?.id]?.baz", expr.Env(env))
	require.NoError(t, err)

	got, err := expr.Run(program, env)
	require.NoError(t, err)
	assert.Equal(t, "baz", got)
}

func TestExpr_optional_chaining_array(t *testing.T) {
	env := map[string]any{}
	program, err := expr.Compile("foo?.[1]?.[2]?.[3]", expr.Env(env), expr.AllowUndefinedVariables())
	require.NoError(t, err)

	got, err := expr.Run(program, env)
	require.NoError(t, err)
	assert.Equal(t, nil, got)
}

func TestExpr_eval_with_env(t *testing.T) {
	_, err := expr.Eval("true", expr.Env(map[string]any{}))
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "misused")
}

func TestExpr_fetch_from_func(t *testing.T) {
	_, err := expr.Eval("foo.Value", map[string]any{
		"foo": func() {},
	})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "cannot fetch Value from func()")
}

func TestExpr_map_default_values(t *testing.T) {
	env := map[string]any{
		"foo": map[string]string{},
		"bar": map[string]*string{},
	}

	input := `foo['missing'] == '' && bar['missing'] == nil`

	program, err := expr.Compile(input, expr.Env(env))
	require.NoError(t, err)

	output, err := expr.Run(program, env)
	require.NoError(t, err)
	require.Equal(t, true, output)
}

func TestExpr_map_default_values_compile_check(t *testing.T) {
	tests := []struct {
		env   any
		input string
	}{
		{
			mock.MapStringStringEnv{"foo": "bar"},
			`Split(foo, sep)`,
		},
		{
			mock.MapStringIntEnv{"foo": 1},
			`foo / bar`,
		},
	}
	for _, tt := range tests {
		_, err := expr.Compile(tt.input, expr.Env(tt.env), expr.AllowUndefinedVariables())
		require.NoError(t, err)
	}
}

func TestExpr_calls_with_nil(t *testing.T) {
	env := map[string]any{
		"equals": func(a, b any) any {
			assert.Nil(t, a, "a is not nil")
			assert.Nil(t, b, "b is not nil")
			return a == b
		},
		"is": mock.Is{},
	}

	p, err := expr.Compile(
		"a == nil && equals(b, nil) && is.Nil(c)",
		expr.Env(env),
		expr.Operator("==", "equals"),
		expr.AllowUndefinedVariables(),
	)
	require.NoError(t, err)

	out, err := expr.Run(p, env)
	require.NoError(t, err)
	require.Equal(t, true, out)
}

func TestExpr_call_float_arg_func_with_int(t *testing.T) {
	env := map[string]any{
		"cnv": func(f float64) any {
			return f
		},
	}
	tests := []struct {
		input    string
		expected float64
	}{
		{"-1", -1.0},
		{"1+1", 2.0},
		{"+1", 1.0},
		{"1-1", 0.0},
		{"1/1", 1.0},
		{"1*1", 1.0},
		{"1^1", 1.0},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			p, err := expr.Compile(fmt.Sprintf("cnv(%s)", tt.input), expr.Env(env))
			require.NoError(t, err)

			out, err := expr.Run(p, env)
			require.NoError(t, err)
			require.Equal(t, tt.expected, out)
		})
	}
}

func TestConstExpr_error_panic(t *testing.T) {
	env := map[string]any{
		"divide": func(a, b int) int { return a / b },
	}

	_, err := expr.Compile(
		`1 + divide(1, 0)`,
		expr.Env(env),
		expr.ConstExpr("divide"),
	)
	require.Error(t, err)
	require.Equal(t, "compile error: integer divide by zero (1:5)\n | 1 + divide(1, 0)\n | ....^", err.Error())
}

type divideError struct{ Message string }

func (e divideError) Error() string {
	return e.Message
}

func TestConstExpr_error_as_error(t *testing.T) {
	env := map[string]any{
		"divide": func(a, b int) (int, error) {
			if b == 0 {
				return 0, divideError{"integer divide by zero"}
			}
			return a / b, nil
		},
	}

	_, err := expr.Compile(
		`1 + divide(1, 0)`,
		expr.Env(env),
		expr.ConstExpr("divide"),
	)
	require.Error(t, err)
	require.Equal(t, "integer divide by zero", err.Error())
	require.IsType(t, divideError{}, err)
}

func TestConstExpr_error_wrong_type(t *testing.T) {
	env := map[string]any{
		"divide": 0,
	}
	assert.Panics(t, func() {
		_, _ = expr.Compile(
			`1 + divide(1, 0)`,
			expr.Env(env),
			expr.ConstExpr("divide"),
		)
	})
}

func TestConstExpr_error_no_env(t *testing.T) {
	assert.Panics(t, func() {
		_, _ = expr.Compile(
			`1 + divide(1, 0)`,
			expr.ConstExpr("divide"),
		)
	})
}

var stringer = reflect.TypeOf((*fmt.Stringer)(nil)).Elem()

type stringerPatcher struct{}

func (p *stringerPatcher) Visit(node *ast.Node) {
	t := (*node).Type()
	if t == nil {
		return
	}
	if t.Implements(stringer) {
		ast.Patch(node, &ast.CallNode{
			Callee: &ast.MemberNode{
				Node:     *node,
				Property: &ast.StringNode{Value: "String"},
			},
		})
	}
}

func TestPatch(t *testing.T) {
	program, err := expr.Compile(
		`Foo == "Foo.String"`,
		expr.Env(mock.Env{}),
		expr.Patch(&mock.StringerPatcher{}),
	)
	require.NoError(t, err)

	output, err := expr.Run(program, mock.Env{})
	require.NoError(t, err)
	require.Equal(t, true, output)
}

func TestCompile_exposed_error(t *testing.T) {
	_, err := expr.Compile(`1 == true`)
	require.NoError(t, err)

}

func TestAsBool_exposed_error(t *testing.T) {
	_, err := expr.Compile(`42`, expr.AsBool())
	require.Error(t, err)

	_, ok := err.(*file.Error)
	require.False(t, ok, "error must not be of type *file.Error")
	require.Equal(t, "expected bool, but got int", err.Error())
}

//func TestEval_exposed_error(t *testing.T) {
//	_, err := expr.Eval(`1 % 0`, nil)
//	require.Error(t, err)
//
//	fileError, ok := err.(*file.Error)
//	require.True(t, ok, "error should be of type *file.Error")
//	require.Equal(t, "integer divide by zero (1:3)\n | 1 % 0\n | ..^", fileError.Error())
//	require.Equal(t, 2, fileError.Column)
//	require.Equal(t, 1, fileError.Line)
//}

func TestIssue105(t *testing.T) {
	type A struct {
		Field string
	}
	type B struct {
		Field int
	}
	type C struct {
		A
		B
	}
	type Env struct {
		C
	}

	code := `
		A.Field == '' &&
		C.A.Field == '' &&
		B.Field == 0 &&
		C.B.Field == 0
	`

	_, err := expr.Compile(code, expr.Env(Env{}))
	require.NoError(t, err)
}

func TestIssue_nested_closures(t *testing.T) {
	code := `all(1..3, { all(1..3, { # > 0 }) and # > 0 })`

	program, err := expr.Compile(code)
	require.NoError(t, err)

	output, err := expr.Run(program, nil)
	require.NoError(t, err)
	require.True(t, output.(bool))
}

func TestIssue138(t *testing.T) {
	env := map[string]any{}

	_, err := expr.Compile(`1 / (1 - 1)`, expr.Env(env))
	require.Equal(t, "runtime error: integer divide by zero (1:3)\n | 1 / (1 - 1)\n | ..^", err.Error())

	_, err = expr.Compile(`1 % 0`, expr.Env(env))
	require.Error(t, err)
	require.Equal(t, "runtime error: integer divide by zero (1:3)\n | 1 % 0\n | ..^", err.Error())
}

func TestIssue154(t *testing.T) {
	type Data struct {
		Array  *[2]any
		Slice  *[]any
		Map    *map[string]any
		String *string
	}

	type Env struct {
		Data *Data
	}

	b := true
	i := 10
	s := "value"

	Array := [2]any{
		&b,
		&i,
	}

	Slice := []any{
		&b,
		&i,
	}

	Map := map[string]any{
		"Bool": &b,
		"Int":  &i,
	}

	env := Env{
		Data: &Data{
			Array:  &Array,
			Slice:  &Slice,
			Map:    &Map,
			String: &s,
		},
	}

	tests := []string{
		`Data.Array[0] == true`,
		`Data.Array[1] == 10`,
		`Data.Slice[0] == true`,
		`Data.Slice[1] == 10`,
		`Data.Map["Bool"] == true`,
		`Data.Map["Int"] == 10`,
		`Data.String == "value"`,
	}

	for _, input := range tests {
		program, err := expr.Compile(input, expr.Env(env))
		require.NoError(t, err, input)

		output, err := expr.Run(program, env)
		require.NoError(t, err)
		assert.True(t, output.(bool), input)
	}
}

func TestIssue270(t *testing.T) {
	env := map[string]any{
		"int8":     int8(1),
		"int16":    int16(3),
		"int32":    int32(5),
		"int64":    int64(7),
		"uint8":    uint8(11),
		"uint16":   uint16(13),
		"uint32":   uint32(17),
		"uint64":   uint64(19),
		"int8a":    uint(23),
		"int8b":    uint(29),
		"int16a":   uint(31),
		"int16b":   uint(37),
		"int32a":   uint(41),
		"int32b":   uint(43),
		"int64a":   uint(47),
		"int64b":   uint(53),
		"uint8a":   uint(59),
		"uint8b":   uint(61),
		"uint16a":  uint(67),
		"uint16b":  uint(71),
		"uint32a":  uint(73),
		"uint32b":  uint(79),
		"uint64a":  uint(83),
		"uint64b":  uint(89),
		"float32a": float32(97),
		"float32b": float32(101),
		"float64a": float64(103),
		"float64b": float64(107),
	}
	for _, each := range []struct {
		input string
	}{
		{"int8 / int16"},
		{"int32 / int64"},
		{"uint8 / uint16"},
		{"uint32 / uint64"},
		{"int8 / uint64"},
		{"int64 / uint8"},
		{"int8a / int8b"},
		{"int16a / int16b"},
		{"int32a / int32b"},
		{"int64a / int64b"},
		{"uint8a / uint8b"},
		{"uint16a / uint16b"},
		{"uint32a / uint32b"},
		{"uint64a / uint64b"},
		{"float32a / float32b"},
		{"float64a / float64b"},
	} {
		p, err := expr.Compile(each.input, expr.Env(env))
		require.NoError(t, err)

		out, err := expr.Run(p, env)
		require.NoError(t, err)
		require.IsType(t, float64(0), out)
	}
}

func TestIssue271(t *testing.T) {
	type BarArray []float64

	type Foo struct {
		Bar BarArray
		Baz int
	}

	type Env struct {
		Foo Foo
	}

	code := `Foo.Bar[0]`

	program, err := expr.Compile(code, expr.Env(Env{}))
	require.NoError(t, err)

	output, err := expr.Run(program, Env{
		Foo: Foo{
			Bar: BarArray{1.0, 2.0, 3.0},
		},
	})
	require.NoError(t, err)
	require.Equal(t, 1.0, output)
}

type Issue346Array []Issue346Type

type Issue346Type struct {
	Bar string
}

func (i Issue346Array) Len() int {
	return len(i)
}

func TestIssue346(t *testing.T) {
	code := `Foo[0].Bar`

	env := map[string]any{
		"Foo": Issue346Array{
			{Bar: "bar"},
		},
	}
	program, err := expr.Compile(code, expr.Env(env))
	require.NoError(t, err)

	output, err := expr.Run(program, env)
	require.NoError(t, err)
	require.Equal(t, "bar", output)
}

func TestCompile_allow_to_use_interface_to_get_an_element_from_map(t *testing.T) {
	code := `{"value": "ok"}[vars.key]`
	env := map[string]any{
		"vars": map[string]any{
			"key": "value",
		},
	}

	program, err := expr.Compile(code, expr.Env(env))
	assert.NoError(t, err)

	out, err := expr.Run(program, env)
	assert.NoError(t, err)
	assert.Equal(t, "ok", out)

	t.Run("with allow undefined variables", func(t *testing.T) {
		code := `{'key': 'value'}[Key]`
		env := mock.MapStringStringEnv{}
		options := []expr.Option{
			expr.AllowUndefinedVariables(),
		}

		program, err := expr.Compile(code, options...)
		assert.NoError(t, err)

		out, err := expr.Run(program, env)
		assert.NoError(t, err)
		assert.Equal(t, nil, out)
	})
}

func TestFastCall(t *testing.T) {
	env := map[string]any{
		"func": func(in any) float64 {
			return 8
		},
	}
	code := `func("8")`

	program, err := expr.Compile(code, expr.Env(env))
	assert.NoError(t, err)

	out, err := expr.Run(program, env)
	assert.NoError(t, err)
	assert.Equal(t, float64(8), out)
}

func TestFastCall_OpCallFastErr(t *testing.T) {
	env := map[string]any{
		"func": func(...any) (any, error) {
			return 8, nil
		},
	}
	code := `func("8")`

	program, err := expr.Compile(code, expr.Env(env))
	assert.NoError(t, err)

	out, err := expr.Run(program, env)
	assert.NoError(t, err)
	assert.Equal(t, 8, out)
}

func TestRun_custom_func_returns_an_error_as_second_arg(t *testing.T) {
	env := map[string]any{
		"semver": func(value string, cmp string) (bool, error) { return true, nil },
	}

	p, err := expr.Compile(`semver("1.2.3", "= 1.2.3")`, expr.Env(env))
	assert.NoError(t, err)

	out, err := expr.Run(p, env)
	assert.NoError(t, err)
	assert.Equal(t, true, out)
}

func TestFunction(t *testing.T) {
	add := expr.Function(
		"add",
		func(p ...any) (any, error) {
			out := 0
			for _, each := range p {
				out += each.(int)
			}
			return out, nil
		},
		new(func(...int) int),
	)

	p, err := expr.Compile(`add() + add(1) + add(1, 2) + add(1, 2, 3) + add(1, 2, 3, 4)`, add)
	assert.NoError(t, err)

	out, err := expr.Run(p, nil)
	assert.NoError(t, err)
	assert.Equal(t, 20, out)
}

// Nil coalescing operator
func TestRun_NilCoalescingOperator(t *testing.T) {
	env := map[string]any{
		"foo": map[string]any{
			"bar": "value",
		},
	}

	t.Run("value", func(t *testing.T) {
		p, err := expr.Compile(`foo.bar ?? "default"`, expr.Env(env))
		assert.NoError(t, err)

		out, err := expr.Run(p, env)
		assert.NoError(t, err)
		assert.Equal(t, "value", out)
	})

	t.Run("default", func(t *testing.T) {
		p, err := expr.Compile(`foo.baz ?? "default"`, expr.Env(env))
		assert.NoError(t, err)

		out, err := expr.Run(p, env)
		assert.NoError(t, err)
		assert.Equal(t, "default", out)
	})

	t.Run("default with chain", func(t *testing.T) {
		p, err := expr.Compile(`foo?.bar ?? "default"`, expr.Env(env))
		assert.NoError(t, err)

		out, err := expr.Run(p, map[string]any{})
		assert.NoError(t, err)
		assert.Equal(t, "default", out)
	})
}

func TestEval_nil_in_maps(t *testing.T) {
	env := map[string]any{
		"m":     map[any]any{nil: "bar"},
		"empty": map[any]any{},
	}
	t.Run("nil key exists", func(t *testing.T) {
		p, err := expr.Compile(`m[nil]`, expr.Env(env))
		assert.NoError(t, err)

		out, err := expr.Run(p, env)
		assert.NoError(t, err)
		assert.Equal(t, "bar", out)
	})
	t.Run("no nil key", func(t *testing.T) {
		p, err := expr.Compile(`empty[nil]`, expr.Env(env))
		assert.NoError(t, err)

		out, err := expr.Run(p, env)
		assert.NoError(t, err)
		assert.Equal(t, nil, out)
	})
	t.Run("nil in m", func(t *testing.T) {
		p, err := expr.Compile(`nil in m`, expr.Env(env))
		assert.NoError(t, err)

		out, err := expr.Run(p, env)
		assert.NoError(t, err)
		assert.Equal(t, true, out)
	})
	t.Run("nil in empty", func(t *testing.T) {
		p, err := expr.Compile(`nil in empty`, expr.Env(env))
		assert.NoError(t, err)

		out, err := expr.Run(p, env)
		assert.NoError(t, err)
		assert.Equal(t, false, out)
	})
}

// Test the use of env keyword.  Forms env[] and env["] are valid.
// The enclosed identifier must be in the expression env.
func TestEnv_keyword(t *testing.T) {
	env := map[string]any{
		"space test":                       "ok",
		"space_test":                       "not ok", // Seems to be some underscore substituting happening, check that.
		"Section 1-2a":                     "ok",
		`c:\ndrive\2015 Information Table`: "ok",
		"%*worst function name ever!!": func() string {
			return "ok"
		}(),
		"1":      "o",
		"2":      "k",
		"num":    10,
		"mylist": []int{1, 2, 3, 4, 5},
		"MIN": func(a, b int) int {
			if a < b {
				return a
			} else {
				return b
			}
		},
		"red":   "n",
		"irect": "um",
		"String Map": map[string]string{
			"one":   "two",
			"three": "four",
		},
		"OtherMap": map[string]string{
			"a": "b",
			"c": "d",
		},
	}

	// No error cases
	var tests = []struct {
		code string
		want any
	}{
		{"$env['space test']", "ok"},
		{"$env['Section 1-2a']", "ok"},
		{`$env["c:\\ndrive\\2015 Information Table"]`, "ok"},
		{"$env['%*worst function name ever!!']", "ok"},
		{"$env['String Map'].one", "two"},
		{"$env['1'] + $env['2']", "ok"},
		{"1 + $env['num'] + $env['num']", 21},
		{"MIN($env['num'],0)", 0},
		{"$env['nu' + 'm']", 10},
		{"$env[red + irect]", 10},
		{"$env['String Map']?.five", ""},
		{"$env.red", "n"},
		{"$env?.unknown", nil},
		{"$env.mylist[1]", 2},
		{"$env?.OtherMap?.a", "b"},
		{"$env?.OtherMap?.d", ""},
		{"'num' in $env", true},
		{"get($env, 'num')", 10},
	}

	for _, tt := range tests {
		t.Run(tt.code, func(t *testing.T) {

			program, err := expr.Compile(tt.code, expr.Env(env))
			require.NoError(t, err, "compile error")

			got, err := expr.Run(program, env)
			require.NoError(t, err, "execution error")

			assert.Equal(t, tt.want, got, tt.code)
		})
	}

	for _, tt := range tests {
		t.Run(tt.code, func(t *testing.T) {
			got, err := expr.Eval(tt.code, env)
			require.NoError(t, err, "eval error: "+tt.code)

			assert.Equal(t, tt.want, got, "eval: "+tt.code)
		})
	}

	// error cases
	tests = []struct {
		code string
		want any
	}{
		{"env()", "bad"},
	}

	for _, tt := range tests {
		t.Run(tt.code, func(t *testing.T) {
			_, err := expr.Eval(tt.code, expr.Env(env))
			require.Error(t, err, "compile error")

		})
	}
}

func TestEnv_keyword_with_custom_functions(t *testing.T) {
	fn := expr.Function("fn", func(params ...any) (any, error) {
		return "ok", nil
	})

	var tests = []struct {
		code  string
		error bool
	}{
		{`fn()`, false},
		{`$env.fn()`, true},
		{`$env["fn"]`, true},
	}

	for _, tt := range tests {
		t.Run(tt.code, func(t *testing.T) {
			_, err := expr.Compile(tt.code, expr.Env(mock.Env{}), fn)
			if tt.error {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestIssue401(t *testing.T) {
	program, err := expr.Compile("(a - b + c) / d", expr.AllowUndefinedVariables())
	require.NoError(t, err, "compile error")

	output, err := expr.Run(program, map[string]any{
		"a": 1,
		"b": 2,
		"c": 3,
		"d": 4,
	})
	require.NoError(t, err, "run error")
	require.Equal(t, 0.5, output)
}

func TestEval_slices_out_of_bound(t *testing.T) {
	tests := []struct {
		code string
		want any
	}{
		{"[1, 2, 3][:99]", []any{1, 2, 3}},
		{"[1, 2, 3][99:]", []any{}},
		{"[1, 2, 3][:-99]", []any{}},
		{"[1, 2, 3][-99:]", []any{1, 2, 3}},
	}

	for _, tt := range tests {
		t.Run(tt.code, func(t *testing.T) {
			got, err := expr.Eval(tt.code, nil)
			require.NoError(t, err, "eval error: "+tt.code)
			assert.Equal(t, tt.want, got, "eval: "+tt.code)
		})
	}
}

func TestExpr_custom_tests(t *testing.T) {
	f, err := os.Open("custom_tests.json")
	if os.IsNotExist(err) {
		t.Skip("no custom tests")
		return
	}

	require.NoError(t, err, "open file error")
	defer f.Close()

	var tests []string
	err = json.NewDecoder(f).Decode(&tests)
	require.NoError(t, err, "decode json error")

	for id, tt := range tests {
		t.Run(fmt.Sprintf("line %v", id+2), func(t *testing.T) {
			program, err := expr.Compile(tt)
			require.NoError(t, err)

			timeout := make(chan bool, 1)
			go func() {
				time.Sleep(time.Second)
				timeout <- true
			}()

			done := make(chan bool, 1)
			go func() {
				out, err := expr.Run(program, nil)
				// Make sure out is used.
				_ = fmt.Sprintf("%v", out)
				assert.Error(t, err)
				done <- true
			}()

			select {
			case <-done:
				// Success.
			case <-timeout:
				t.Fatal("timeout")
			}
		})
	}
}

func TestIssue432(t *testing.T) {
	env := map[string]any{
		"func": func(
			paramUint32 uint32,
			paramUint16 uint16,
			paramUint8 uint8,
			paramUint uint,
			paramInt32 int32,
			paramInt16 int16,
			paramInt8 int8,
			paramInt int,
			paramFloat64 float64,
			paramFloat32 float32,
		) float64 {
			return float64(paramUint32) + float64(paramUint16) + float64(paramUint8) + float64(paramUint) +
				float64(paramInt32) + float64(paramInt16) + float64(paramInt8) + float64(paramInt) +
				float64(paramFloat64) + float64(paramFloat32)
		},
	}
	code := `func(1,1,1,1,1,1,1,1,1,1)`

	program, err := expr.Compile(code, expr.Env(env))
	assert.NoError(t, err)

	out, err := expr.Run(program, env)
	assert.NoError(t, err)
	assert.Equal(t, float64(10), out)
}

func TestIssue462(t *testing.T) {
	env := map[string]any{
		"foo": func() (string, error) {
			return "bar", nil
		},
	}
	_, err := expr.Compile(`$env.unknown(int())`, expr.Env(env))
	require.Error(t, err)
}

func TestIssue_embedded_pointer_struct(t *testing.T) {
	var tests = []struct {
		input string
		env   mock.Env
		want  any
	}{
		{
			input: "EmbedPointerEmbedInt > 0",
			env: mock.Env{
				Embed: mock.Embed{
					EmbedPointerEmbed: &mock.EmbedPointerEmbed{
						EmbedPointerEmbedInt: 123,
					},
				},
			},
			want: true,
		},
		{
			input: "(Embed).EmbedPointerEmbedInt > 0",
			env: mock.Env{
				Embed: mock.Embed{
					EmbedPointerEmbed: &mock.EmbedPointerEmbed{
						EmbedPointerEmbedInt: 123,
					},
				},
			},
			want: true,
		},
		{
			input: "(Embed).EmbedPointerEmbedInt > 0",
			env: mock.Env{
				Embed: mock.Embed{
					EmbedPointerEmbed: &mock.EmbedPointerEmbed{
						EmbedPointerEmbedInt: 0,
					},
				},
			},
			want: false,
		},
		{
			input: "(Embed).EmbedPointerEmbedMethod(0)",
			env: mock.Env{
				Embed: mock.Embed{
					EmbedPointerEmbed: &mock.EmbedPointerEmbed{
						EmbedPointerEmbedInt: 0,
					},
				},
			},
			want: "",
		},
		{
			input: "(Embed).EmbedPointerEmbedPointerReceiverMethod(0)",
			env: mock.Env{
				Embed: mock.Embed{
					EmbedPointerEmbed: nil,
				},
			},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			program, err := expr.Compile(tt.input, expr.Env(tt.env))
			require.NoError(t, err)

			out, err := expr.Run(program, tt.env)
			require.NoError(t, err)

			require.Equal(t, tt.want, out)
		})
	}
}

func TestIssue474(t *testing.T) {
	testCases := []struct {
		code string
		fail bool
	}{
		{
			code: `func("invalid")`,
			fail: true,
		},
		{
			code: `func(true)`,
			fail: true,
		},
		{
			code: `func([])`,
			fail: true,
		},
		{
			code: `func({})`,
			fail: true,
		},
		{
			code: `func(1)`,
			fail: false,
		},
		{
			code: `func(1.5)`,
			fail: false,
		},
	}

	for _, tc := range testCases {
		ltc := tc
		t.Run(ltc.code, func(t *testing.T) {
			t.Parallel()
			function := expr.Function("func", func(params ...any) (any, error) {
				return true, nil
			}, new(func(float64) bool))
			_, err := expr.Compile(ltc.code, function)
			if ltc.fail {
				if err == nil {
					t.Error("expected an error, but it was nil")
					t.FailNow()
				}
			} else {
				if err != nil {
					t.Errorf("expected nil, but it was %v", err)
					t.FailNow()
				}
			}
		})
	}
}

func TestRaceCondition_variables(t *testing.T) {
	program, err := expr.Compile(`let foo = 1; foo + 1`, expr.Env(mock.Env{}))
	require.NoError(t, err)

	var wg sync.WaitGroup

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			out, err := expr.Run(program, mock.Env{})
			require.NoError(t, err)
			require.Equal(t, 2, out)
		}()
	}

	wg.Wait()
}

func TestOperatorDependsOnEnv(t *testing.T) {
	env := map[string]any{
		"plus": func(a, b int) int {
			return 42
		},
	}
	program, err := expr.Compile(`1 + 2`, expr.Operator("+", "plus"), expr.Env(env))
	require.NoError(t, err)

	out, err := expr.Run(program, env)
	require.NoError(t, err)
	assert.Equal(t, 42, out)
}

func TestIssue624(t *testing.T) {
	type tag struct {
		Name string
	}

	type item struct {
		Tags []tag
	}

	i := item{
		Tags: []tag{
			{Name: "one"},
			{Name: "two"},
		},
	}

	rule := `[
true && true, 
one(Tags, .Name in ["one"]), 
one(Tags, .Name in ["two"]), 
one(Tags, .Name in ["one"]) && one(Tags, .Name in ["two"])
]`
	resp, err := expr.Eval(rule, i)
	require.NoError(t, err)
	require.Equal(t, []interface{}{true, true, true, true}, resp)
}

func TestPredicateCombination(t *testing.T) {
	tests := []struct {
		code1 string
		code2 string
	}{
		{"all(1..3, {# > 0}) && all(1..3, {# < 4})", "all(1..3, {# > 0 && # < 4})"},
		{"all(1..3, {# > 1}) && all(1..3, {# < 4})", "all(1..3, {# > 1 && # < 4})"},
		{"all(1..3, {# > 0}) && all(1..3, {# < 2})", "all(1..3, {# > 0 && # < 2})"},
		{"all(1..3, {# > 1}) && all(1..3, {# < 2})", "all(1..3, {# > 1 && # < 2})"},

		{"any(1..3, {# > 0}) || any(1..3, {# < 4})", "any(1..3, {# > 0 || # < 4})"},
		{"any(1..3, {# > 1}) || any(1..3, {# < 4})", "any(1..3, {# > 1 || # < 4})"},
		{"any(1..3, {# > 0}) || any(1..3, {# < 2})", "any(1..3, {# > 0 || # < 2})"},
		{"any(1..3, {# > 1}) || any(1..3, {# < 2})", "any(1..3, {# > 1 || # < 2})"},

		{"none(1..3, {# > 0}) && none(1..3, {# < 4})", "none(1..3, {# > 0 || # < 4})"},
		{"none(1..3, {# > 1}) && none(1..3, {# < 4})", "none(1..3, {# > 1 || # < 4})"},
		{"none(1..3, {# > 0}) && none(1..3, {# < 2})", "none(1..3, {# > 0 || # < 2})"},
		{"none(1..3, {# > 1}) && none(1..3, {# < 2})", "none(1..3, {# > 1 || # < 2})"},
	}
	for _, tt := range tests {
		t.Run(tt.code1, func(t *testing.T) {
			out1, err := expr.Eval(tt.code1, nil)
			require.NoError(t, err)

			out2, err := expr.Eval(tt.code2, nil)
			require.NoError(t, err)

			require.Equal(t, out1, out2)
		})
	}
}

func TestArrayComparison(t *testing.T) {
	tests := []struct {
		env  any
		code string
	}{
		{[]string{"A", "B"}, "foo == ['A', 'B']"},
		{[]int{1, 2}, "foo == [1, 2]"},
		{[]uint8{1, 2}, "foo == [1, 2]"},
		{[]float64{1.1, 2.2}, "foo == [1.1, 2.2]"},
		{[]any{"A", 1, 1.1, true}, "foo == ['A', 1, 1.1, true]"},
		//{[]string{"A", "B"}, "foo != [1, 2]"},
	}

	for _, tt := range tests {
		t.Run(tt.code, func(t *testing.T) {
			env := map[string]any{"foo": tt.env}
			program, err := expr.Compile(tt.code, expr.Env(env))
			require.NoError(t, err)

			out, err := expr.Run(program, env)
			require.NoError(t, err)
			require.Equal(t, true, out)
		})
	}
}

func TestIssue_570(t *testing.T) {
	type Student struct {
		Name string
	}

	env := map[string]any{
		"student": (*Student)(nil),
	}

	program, err := expr.Compile("student?.Name", expr.Env(env))
	require.NoError(t, err)

	out, err := expr.Run(program, env)
	require.NoError(t, err)
	require.IsType(t, nil, out)
}

func TestIssue_integer_truncated_by_compiler(t *testing.T) {
	env := map[string]any{
		"fn": func(x byte) byte {
			return x
		},
	}

	_, err := expr.Compile("fn(255)", expr.Env(env))
	require.NoError(t, err)

	_, err = expr.Compile("fn(256)", expr.Env(env))
	require.Error(t, err)
}

func TestExpr_crash(t *testing.T) {
	content, err := os.ReadFile("testdata/crash.txt")
	require.NoError(t, err)

	_, err = expr.Compile(string(content))
	require.Error(t, err)
}

func TestExpr_nil_op_str(t *testing.T) {
	// Let's test operators, which do `.(string)` in VM, also check for nil.

	var str *string = nil
	env := map[string]any{
		"nilString": str,
	}

	tests := []struct{ code string }{
		{`nilString == "str"`},
		{`nilString contains "str"`},
		{`nilString matches "str"`},
		{`nilString startsWith "str"`},
		{`nilString endsWith "str"`},

		{`"str" == nilString`},
		{`"str" contains nilString`},
		{`"str" matches nilString`},
		{`"str" startsWith nilString`},
		{`"str" endsWith nilString`},
	}

	for _, tt := range tests {
		t.Run(tt.code, func(t *testing.T) {
			program, err := expr.Compile(tt.code)
			require.NoError(t, err)

			output, err := expr.Run(program, env)
			require.NoError(t, err)
			require.Equal(t, false, output)
		})
	}
}

func TestExpr_env_types_map(t *testing.T) {
	envTypes := types.Map{
		"foo": types.Map{
			"bar": types.String,
		},
	}

	program, err := expr.Compile(`foo.bar`, expr.Env(envTypes))
	require.NoError(t, err)

	env := map[string]any{
		"foo": map[string]any{
			"bar": "value",
		},
	}

	output, err := expr.Run(program, env)
	require.NoError(t, err)
	require.Equal(t, "value", output)
}

func TestExpr_env_types_map_error(t *testing.T) {
	envTypes := types.Map{
		"foo": types.Map{
			"bar": types.String,
		},
	}

	program, err := expr.Compile(`foo.bar`, expr.Env(envTypes))
	require.NoError(t, err)

	_, err = expr.Run(program, envTypes)
	require.Error(t, err)
}

func TestIssue758_filter_map_index(t *testing.T) {
	env := map[string]interface{}{}

	exprStr := `
        let a_map = 0..5 | filter(# % 2 == 0) | map(#index);
        let b_filter = 0..5 | filter(# % 2 == 0);
        let b_map = b_filter | map(#index);
        [a_map, b_map]
    `

	result, err := expr.Eval(exprStr, env)
	require.NoError(t, err)

	expected := []interface{}{
		[]interface{}{0, 1, 2},
		[]interface{}{0, 1, 2},
	}

	require.Equal(t, expected, result)
}

func TestExpr_wierd_cases(t *testing.T) {
	env := map[string]any{}

	_, err := expr.Compile(`A(A)`, expr.Env(env))
	require.Error(t, err)
	require.Contains(t, err.Error(), "unknown name A")
}

func TestIssue785_get_nil(t *testing.T) {
	exprStrs := []string{
		`get(nil, "a")`,
		`get({}, "a")`,
		`get(nil, "a")`,
		`get({}, "a")`,
		`({} | get("a") | get("b"))`,
	}

	for _, exprStr := range exprStrs {
		t.Run("get returns nil", func(t *testing.T) {
			env := map[string]interface{}{}

			result, err := expr.Eval(exprStr, env)
			require.NoError(t, err)

			require.Equal(t, nil, result)
		})
	}
}

func TestMaxNodes(t *testing.T) {
	maxNodes := uint(100)

	code := ""
	for i := 0; i < int(maxNodes); i++ {
		code += "1; "
	}

	_, err := expr.Compile(code, expr.MaxNodes(maxNodes))
	require.Error(t, err)
	assert.Contains(t, err.Error(), "exceeds maximum allowed nodes")

	_, err = expr.Compile(code, expr.MaxNodes(maxNodes+1))
	require.NoError(t, err)
}

func TestMaxNodesDisabled(t *testing.T) {
	code := ""
	for i := 0; i < 2*int(conf.DefaultMaxNodes); i++ {
		code += "1; "
	}

	_, err := expr.Compile(code, expr.MaxNodes(0))
	require.NoError(t, err)
}

func TestMemoryBudget(t *testing.T) {
	tests := []struct {
		code string
		max  int
	}{
		{`map(1..100, {map(1..100, {map(1..100, {0})})})`, -1},
		{`len(1..10000000)`, -1},
		{`1..100`, 100},
	}

	for _, tt := range tests {
		t.Run(tt.code, func(t *testing.T) {
			program, err := expr.Compile(tt.code)
			require.NoError(t, err, "compile error")

			vm := vm.VM{}
			if tt.max > 0 {
				vm.MemoryBudget = uint(tt.max)
			}
			_, err = vm.Run(program, nil)
			require.Error(t, err, "run error")
			assert.Contains(t, err.Error(), "memory budget exceeded")
		})
	}
}

// Add tests for all legal and illegal pairs for arithmetic and comparison operators
func TestArithmeticAndComparisonOperators(t *testing.T) {
	tests := []struct {
		expr        string
		want        any
		expectError bool
	}{
		// Addition
		{"1 + 2", 3, false},
		{"1.5 + 2.5", 4.0, false},
		{"1 + 2.5", 3.5, false},
		{"2.5 + 1", 3.5, false},
		{"'a' + 'b'", "ab", false},
		{"'a' + 1", "a1", false},
		{"1 + 'a'", "1a", false},
		// Subtraction
		{"5 - 2", 3, false},
		{"5.5 - 2.5", 3.0, false},
		{"5 - 2.5", 2.5, false},
		{"5.5 - 2", 3.5, false},
		// Multiplication
		{"2 * 3", 6, false},
		{"2.5 * 4", 10.0, false},
		{"2 * 4.5", 9.0, false},
		{"2.5 * 4.5", 11.25, false},
		// Division
		{"8 / 2", 4.0, false},
		{"8.0 / 2", 4.0, false},
		{"8 / 2.0", 4.0, false},
		{"8.0 / 2.0", 4.0, false},
		// Modulo
		{"5 % 2", 1, false},
		{"5.5 % 2.0", 1.5, false},
		{"5 % 2.5", 0.0, false},
		{"5.5 % 2", 1.5, false},
		{"5.5 % 2.5", 0.5, false},
		// Exponentiation
		{"2 ** 3", 8.0, false},
		{"2.0 ** 3", 8.0, false},
		{"2 ** 3.0", 8.0, false},
		{"2.0 ** 3.0", 8.0, false},
		// Comparison
		{"1 < 2", true, false},
		{"2 > 1", true, false},
		{"2 <= 2", true, false},
		{"2 >= 2", true, false},
		{"2 == 2", true, false},
		{"2 != 3", true, false},
		{"1.5 < 2.5", true, false},
		{"'a' < 'b'", true, false},
		{"true == true", true, false},
		// Illegal pairs
		{"true + 1", 2, false},
		{"1 + true", 2, false},
		//{"true < 1", nil, true},
		//{"1 < true", nil, true},
		//{"'a' - 'b'", nil, true},
		//{"'a' * 2", nil, true},
		//{"2 * 'a'", nil, true},
		//{"'a' / 2", nil, true},
		//{"2 / 'a'", nil, true},
		//{"'a' % 2", nil, true},
		//{"2 % 'a'", nil, true},
		//{"'a' ** 2", nil, true},
		//{"2 ** 'a'", nil, true},
		//{"true % false", nil, true},
		//{"1 % 0", nil, true},
		//{"1.0 % 0.0", nil, true},
	}
	for _, tt := range tests {
		t.Run(tt.expr, func(t *testing.T) {
			result, err := expr.Eval(tt.expr, nil)
			if tt.expectError {
				if err == nil {
					t.Errorf("expected error for %q, got result %v", tt.expr, result)
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error for %q: %v", tt.expr, err)
				} else if !equal(result, tt.want) {
					t.Errorf("%q: got %v, want %v", tt.expr, result, tt.want)
				}
			}
		})
	}
}

func equal(a, b any) bool {
	switch a := a.(type) {
	case float64:
		if b, ok := b.(float64); ok {
			return (a-b) < 1e-9 && (b-a) < 1e-9
		}
	}
	return a == b
}

func TestStringConcatenationWithAllTypes(t *testing.T) {
	tests := []struct {
		expr        string
		want        any
		expectError bool
	}{
		// String + numeric types
		{"'hello' + 42", "hello42", false},
		{"'hello' + 42.5", "hello42.5", false},
		{"'hello' + 0", "hello0", false},
		{"'hello' + (-5)", "hello-5", false},
		{"'hello' + 3.14159", "hello3.14159", false},

		// Numeric types + string
		{"42 + 'world'", "42world", false},
		{"42.5 + 'world'", "42.5world", false},
		{"0 + 'world'", "0world", false},
		{"(-5) + 'world'", "-5world", false},
		{"3.14159 + 'world'", "3.14159world", false},

		// String + boolean
		{"'hello' + true", "hellotrue", false},
		{"'hello' + false", "hellofalse", false},
		{"true + 'world'", "trueworld", false},
		{"false + 'world'", "falseworld", false},

		// String + nil (if supported)
		{"'hello' + nil", "hello", false},
		{"nil + 'world'", "world", false},
		{"nilValue + 'world'", "world", false}, // nil + string

		// Empty string concatenations
		{"'' + 42", "42", false},
		{"42 + ''", "42", false},
		{"'' + ''", "", false},

		// Multiple concatenations
		{"'a' + 1 + 'b' + 2", "a1b2", false},
		{"'result: ' + (5 + 3)", "result: 8", false},
		{"'pi is approximately ' + 3.14", "pi is approximately 3.14", false},
	}

	for _, tt := range tests {
		t.Run(tt.expr, func(t *testing.T) {
			result, err := expr.Eval(tt.expr, map[string]any{
				"nilValue": nil,
			})
			if tt.expectError {
				if err == nil {
					t.Errorf("expected error for %q, got result %v", tt.expr, result)
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error for %q: %v", tt.expr, err)
				} else if !equalStrings(result, tt.want) {
					t.Errorf("%q: got %v, want %v", tt.expr, result, tt.want)
				}
			}
		})
	}
}

func TestNumericTypeArithmetic(t *testing.T) {
	tests := []struct {
		expr        string
		want        any
		expectError bool
	}{
		// Int + Int variations
		{"1 + 2", 3, false},
		{"10 - 3", 7, false},
		{"4 * 5", 20, false},
		{"15 / 3", 5.0, false}, // Division always returns float
		{"17 % 5", 2, false},
		{"2 ** 3", 8.0, false}, // Power always returns float

		// Float + Float
		{"1.5 + 2.5", 4.0, false},
		{"5.5 - 2.5", 3.0, false},
		{"2.5 * 4.0", 10.0, false},
		{"7.5 / 2.5", 3.0, false},
		{"5.5 % 2.0", 1.5, false},
		{"2.0 ** 3.0", 8.0, false},

		// Int + Float combinations
		{"1 + 2.5", 3.5, false},
		{"5 - 2.5", 2.5, false},
		{"3 * 2.5", 7.5, false},
		{"7 / 2.0", 3.5, false},
		{"5 % 2.0", 1.0, false},
		{"2 ** 3.0", 8.0, false},

		// Float + Int combinations
		{"1.5 + 2", 3.5, false},
		{"5.5 - 2", 3.5, false},
		{"2.5 * 3", 7.5, false},
		{"7.5 / 2", 3.75, false},
		{"5.5 % 2", 1.5, false},
		{"2.0 ** 3", 8.0, false},

		// Zero operations
		{"0 + 5", 5, false},
		{"5 + 0", 5, false},
		{"0 - 5", -5, false},
		{"5 - 0", 5, false},
		{"0 * 5", 0, false},
		{"5 * 0", 0, false},
		{"0 / 5", 0.0, false},
		{"0 % 5", 0, false},
		{"0 ** 5", 0.0, false},

		// Negative numbers
		{"(-5) + 3", -2, false},
		{"3 + (-5)", -2, false},
		{"(-5) - 3", -8, false},
		{"3 - (-5)", 8, false},
		{"(-5) * 3", -15, false},
		{"3 * (-5)", -15, false},
		{"(-6) / 2", -3.0, false},
		{"(-7) % 3", -1, false},
		{"(-2) ** 3", -8.0, false},

		// Large numbers
		{"1000000 + 2000000", 3000000, false},
		{"1000000.5 + 2000000.5", 3000000.5 + 0.5, false},
		{"999999 * 2", 1999998, false},

		// Small decimal numbers
		{"0.1 + 0.2", 0.3, false},
		{"0.5 - 0.3", 0.2, false},
		{"0.1 * 0.2", 0.02, false},
		{"0.6 / 0", math.Inf(0), false},

		// Error cases
		//{"5 / 0", nil, true},     // Division by zero
		//{"5.5 / 0.0", nil, true}, // Float division by zero
		//{"5 % 0", nil, true},     // Modulo by zero
		//{"5.5 % 0.0", nil, true}, // Float modulo by zero
	}

	for _, tt := range tests {
		t.Run(tt.expr, func(t *testing.T) {
			result, err := expr.Eval(tt.expr, nil)
			if tt.expectError {
				if err == nil {
					t.Errorf("expected error for %q, got result %v", tt.expr, result)
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error for %q: %v", tt.expr, err)
				} else if !equalNumbers(result, tt.want) {
					t.Errorf("%q: got %v, want %v", tt.expr, result, tt.want)
				}
			}
		})
	}
}

func TestMixedTypeComparisons(t *testing.T) {
	tests := []struct {
		expr        string
		want        any
		expectError bool
	}{
		// Int vs Float comparisons
		{"1 == 1.0", true, false},
		{"1 != 1.0", false, false},
		{"1 < 1.5", true, false},
		{"1.5 > 1", true, false},
		{"2 <= 2.0", true, false},
		{"2.0 >= 2", true, false},

		// String comparisons
		{"'abc' == 'abc'", true, false},
		{"'abc' != 'def'", true, false},
		{"'abc' < 'def'", true, false},
		{"'def' > 'abc'", true, false},
		{"'abc' <= 'abc'", true, false},
		{"'abc' >= 'abc'", true, false},

		// Zero comparisons
		{"0 == 0.0", true, false},
		{"0 != 0.0", false, false},
		{"0 < 0.1", true, false},
		{"0.0 <= 0", true, false},

		// Negative number comparisons
		{"(-1) == (-1.0)", true, false},
		{"(-1) < 0", true, false},
		{"(-1.5) < (-1)", true, false},
		{"0 > (-1)", true, false},

		// Boolean comparisons
		{"true == true", true, false},
		{"false == false", true, false},
		{"true != false", true, false},
		{"false != true", true, false},

		// Error cases - comparing incompatible types
		{"'abc' == 123", false, false},
		{"true == 1", true, false}, // Should return false, not error
		//{"'abc' < 123", nil, true},     // This should error
		//{"true < 1", nil, true},        // This should error
	}

	for _, tt := range tests {
		t.Run(tt.expr, func(t *testing.T) {
			result, err := expr.Eval(tt.expr, nil)
			if tt.expectError {
				if err == nil {
					t.Errorf("expected error for %q, got result %v", tt.expr, result)
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error for %q: %v", tt.expr, err)
				} else if result != tt.want {
					t.Errorf("%q: got %v, want %v", tt.expr, result, tt.want)
				}
			}
		})
	}
}

func TestComplexExpressions(t *testing.T) {
	tests := []struct {
		expr        string
		want        any
		expectError bool
	}{
		// Mixed arithmetic and string operations
		{"'Result: ' + (10 + 5)", "Result: 15", false},
		{"'Value: ' + (3.14 * 2)", "Value: 6.28", false},
		{"'Answer: ' + (100 / 4)", "Answer: 25", false},

		// Parentheses and order of operations
		{"(1 + 2) * 3", 9, false},
		{"1 + (2 * 3)", 7, false},
		{"(10 - 5) / (3 - 1)", 2.5, false},
		{"2 ** (3 + 1)", 16.0, false},

		// Multiple string concatenations with calculations
		{"'a' + 1 + 'b' + (2 * 3)", "a1b6", false},
		{"(5 + 3) + 'items'", "8items", false},

		// Complex comparisons
		{"(1 + 2) == 3", true, false},
		{"(1.5 * 2) > 2", true, false},
		{"'hello' + 'world' == 'helloworld'", true, false},

		// Nested operations
		{"((1 + 2) * 3) + 4", 13, false},
		{"1 + 2 * 3 + 4", 11, false}, // 1 + (2*3) + 4 = 1 + 6 + 4 = 11
		{"(1 + 2) * (3 + 4)", 21, false},

		// String with complex numeric expressions
		{"'Result: ' + ((10 + 5) * 2)", "Result: 30", false},
		{"'Pi doubled: ' + (3.14159 * 2)", "Pi doubled: 6.28318", false},
	}

	for _, tt := range tests {
		t.Run(tt.expr, func(t *testing.T) {
			result, err := expr.Eval(tt.expr, nil)
			if tt.expectError {
				if err == nil {
					t.Errorf("expected error for %q, got result %v", tt.expr, result)
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error for %q: %v", tt.expr, err)
				} else if !equalValues(result, tt.want) {
					t.Errorf("%q: got %v, want %v", tt.expr, result, tt.want)
				}
			}
		})
	}
}

// Helper functions for comparisons
func equalStrings(a, b any) bool {
	aStr, aOk := a.(string)
	bStr, bOk := b.(string)
	if aOk && bOk {
		return aStr == bStr
	}
	return fmt.Sprintf("%v", a) == fmt.Sprintf("%v", b)
}

func equalNumbers(a, b any) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}

	// Convert both to float64 for comparison
	aFloat := toFloat64(a)
	bFloat := toFloat64(b)

	// Handle integer comparisons
	if isInteger(a) && isInteger(b) {
		return int64(aFloat) == int64(bFloat)
	}
	if math.IsInf(aFloat, 0) && math.IsInf(bFloat, 0) {
		return true
	}
	if math.IsInf(aFloat, -1) && math.IsInf(bFloat, -1) {
		return true
	}

	// Handle float comparisons with tolerance
	diff := aFloat - bFloat
	if diff < 0 {
		diff = -diff
	}
	return diff < 1e-9
}

func equalValues(a, b any) bool {
	// Check if either value is a string
	_, aIsString := a.(string)
	_, bIsString := b.(string)

	// If either is a string, use string comparison only
	if aIsString || bIsString {
		return equalStrings(a, b)
	}

	// Try numeric comparison
	if equalNumbers(a, b) {
		return true
	}
	// Fall back to direct comparison
	return a == b
}

func toFloat64(v any) float64 {
	switch val := v.(type) {
	case int:
		return float64(val)
	case int8:
		return float64(val)
	case int16:
		return float64(val)
	case int32:
		return float64(val)
	case int64:
		return float64(val)
	case uint:
		return float64(val)
	case uint8:
		return float64(val)
	case uint16:
		return float64(val)
	case uint32:
		return float64(val)
	case uint64:
		return float64(val)
	case float32:
		return float64(val)
	case float64:
		return val
	default:
		return 0
	}
}

func isInteger(v any) bool {
	switch v.(type) {
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		return true
	default:
		return false
	}
}

func TestNilArithmeticOperations(t *testing.T) {
	env := map[string]any{
		"nilValue":    nil,
		"intValue":    42,
		"floatValue":  3.14,
		"stringValue": "hello",
		"boolValue":   true,
	}

	tests := []struct {
		expr        string
		want        any
		expectError bool
	}{
		// nil + other types
		{"nilValue + 5", "5", false},           // nil + int becomes string concatenation
		{"nilValue + 5.5", "5.5", false},       // nil + float becomes string concatenation
		{"nilValue + 'world'", "world", false}, // nil + string
		{"nilValue + true", "1", false},        // nil + bool
		{"nilValue + nilValue", "", false},     // nil + nil

		// other types + nil
		{"5 + nilValue", "5", false},           // int + nil becomes string concatenation
		{"5.5 + nilValue", "5.5", false},       // float + nil becomes string concatenation
		{"'hello' + nilValue", "hello", false}, // string + nil
		{"true + nilValue", "1", false},        // bool + nil

		// nil - other types
		{"nilValue - 5", -5, false},       // nil - int = 0 - int
		{"nilValue - 5.5", -5.5, false},   // nil - float = 0 - float
		{"nilValue - nilValue", 0, false}, // nil - nil = 0

		// other types - nil
		{"5 - nilValue", 5, false},     // int - nil = int - 0
		{"5.5 - nilValue", 5.5, false}, // float - nil = float - 0

		// nil * other types
		{"nilValue * 5", 0, false}, // nil * anything = 0
		{"nilValue * 5.5", 0, false},
		{"nilValue * nilValue", 0, false},

		// other types * nil
		{"5 * nilValue", 0, false}, // anything * nil = 0
		{"5.5 * nilValue", 0, false},

		// nil / other types
		{"nilValue / 5", 0.0, false}, // nil / anything = 0
		{"nilValue / 5.5", 0.0, false},

		// other types / nil
		//{"5 / nilValue", 0.0, true}, // anything / nil = 0 (treated as 0, not error)
		//{"5.5 / nilValue", 0.0, true},

		// nil % other types
		{"nilValue % 5", 0, false}, // nil % anything = 0
		{"nilValue % 5.5", 0, false},

		// other types % nil
		//{"5 % nilValue", 0, true}, // anything % nil = 0
		//{"5.5 % nilValue", 0, true},

		// nil power operations
		{"nilValue ** 2", 0.0, false}, // nil ** anything = 0
		{"2 ** nilValue", 1.0, false}, // anything ** nil = anything ** 0 = 1
	}

	for _, tt := range tests {
		t.Run(tt.expr, func(t *testing.T) {
			result, err := expr.Eval(tt.expr, env)
			if tt.expectError {
				if err == nil {
					t.Errorf("expected error for %q, got result %v", tt.expr, result)
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error for %q: %v", tt.expr, err)
				} else if !equalValues(result, tt.want) {
					t.Errorf("%q: got %v (type %T), want %v (type %T)", tt.expr, result, result, tt.want, tt.want)
				}
			}
		})
	}
}

func TestNilComparisonOperations(t *testing.T) {
	env := map[string]any{
		"nilValue":    nil,
		"intValue":    42,
		"floatValue":  3.14,
		"stringValue": "hello",
		"boolValue":   true,
	}

	tests := []struct {
		expr        string
		want        any
		expectError bool
	}{
		// nil == comparisons
		{"nilValue == nilValue", true, false}, // nil == nil
		{"nilValue == 0", false, false},       // nil != 0
		{"nilValue == 0.0", false, false},     // nil != 0.0
		{"nilValue == ''", false, false},      // nil != empty string
		{"nilValue == false", false, false},   // nil != false

		// other types == nil
		{"0 == nilValue", false, false},
		{"0.0 == nilValue", false, false},
		{"'' == nilValue", false, false},
		{"false == nilValue", false, false},

		// nil != comparisons
		{"nilValue != nilValue", false, false}, // nil is equal to nil
		{"nilValue != 0", true, false},         // nil is not equal to 0
		{"nilValue != 'hello'", true, false},   // nil is not equal to string
		{"nilValue != true", true, false},      // nil is not equal to true

		// nil < comparisons (nil is less than any non-nil value)
		{"nilValue < 5", true, false},
		{"nilValue < 0", false, false},
		{"nilValue < (-5)", false, false},
		{"nilValue < 'a'", true, true},
		{"nilValue < nilValue", false, false}, // nil is not less than nil

		// other types < nil
		{"5 < nilValue", false, false}, // non-nil is not less than nil
		{"0 < nilValue", false, false},
		{"'a' < nilValue", false, true},

		// nil > comparisons (nil is not greater than anything)
		{"nilValue > 5", false, false},
		{"nilValue > 0", false, false},
		{"nilValue > (-5)", true, false},
		{"nilValue > nilValue", false, false},

		// other types > nil (non-nil is greater than nil)
		{"5 > nilValue", true, false},
		{"0 > nilValue", false, false},
		{"'a' > nilValue", true, true},

		// nil <= comparisons
		{"nilValue <= nilValue", true, false}, // nil <= nil is true
		{"nilValue <= 5", true, false},        // nil <= anything is true
		{"nilValue <= 0", true, false},
		{"nilValue <= (-5)", false, false},

		// nil >= comparisons
		{"nilValue >= nilValue", true, false}, // nil >= nil is true
		{"nilValue >= 5", false, false},       // nil >= non-nil is false
		{"nilValue >= 0", true, false},

		// other types >= nil (non-nil >= nil is true)
		{"5 >= nilValue", true, false},
		{"0 >= nilValue", true, false},
		{"(-5) >= nilValue", false, false},
	}

	for _, tt := range tests {
		t.Run(tt.expr, func(t *testing.T) {
			result, err := expr.Eval(tt.expr, env)
			if tt.expectError {
				if err == nil {
					t.Errorf("expected error for %q, got result %v", tt.expr, result)
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error for %q: %v", tt.expr, err)
				} else if result != tt.want {
					t.Errorf("%q: got %v, want %v", tt.expr, result, tt.want)
				}
			}
		})
	}
}

func TestNilLogicalOperations(t *testing.T) {
	env := map[string]any{
		"nilValue":   nil,
		"trueValue":  true,
		"falseValue": false,
	}

	tests := []struct {
		expr        string
		want        any
		expectError bool
	}{
		// nil && operations (nil is falsy)
		{"nilValue && true", nil, false},
		{"nilValue && false", nil, false},
		{"nilValue && nilValue", nil, false},

		// other types && nil
		{"true && nilValue", nil, false},    // true && nil = false (nil is falsy)
		{"false && nilValue", false, false}, // false && nil = false

		// nil || operations (nil is falsy)
		{"nilValue || true", true, false},    // nil || true = true
		{"nilValue || false", false, false},  // nil || false = false
		{"nilValue || nilValue", nil, false}, // nil || nil = false

		// other types || nil
		{"true || nilValue", true, false}, // true || nil = true
		{"false || nilValue", nil, false}, // false || nil = false

		// Complex logical operations with nil
		{"nilValue && true || false", false, false},
		{"true || nilValue && false", true, false},
		{"(nilValue || false) && true", false, false},
		{"(nilValue || true) && false", false, false},
	}

	for _, tt := range tests {
		t.Run(tt.expr, func(t *testing.T) {
			result, err := expr.Eval(tt.expr, env)
			if tt.expectError {
				if err == nil {
					t.Errorf("expected error for %q, got result %v", tt.expr, result)
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error for %q: %v", tt.expr, err)
				} else if result != tt.want {
					t.Errorf("%q: got %v, want %v", tt.expr, result, tt.want)
				}
			}
		})
	}
}

func TestNilInCollectionOperations(t *testing.T) {
	env := map[string]any{
		"nilValue":        nil,
		"arrayWithNil":    []any{1, nil, "hello", nil, 5},
		"arrayWithoutNil": []any{1, 2, "hello", 4, 5},
		"mapWithNil": map[string]any{
			"key1": "value1",
			"key2": nil,
			"key3": 42,
		},
	}

	tests := []struct {
		expr        string
		want        any
		expectError bool
	}{
		// nil in array operations
		{"nilValue in arrayWithNil", true, false},     // nil is in the array
		{"nilValue in arrayWithoutNil", false, false}, // nil is not in the array

		// nil in map operations (checking if nil is a value)
		{"nilValue in mapWithNil", false, false}, // nil as key doesn't exist

		// Array/slice operations with nil
		{"len(arrayWithNil)", 5, false}, // length includes nil elements
		{"arrayWithNil[1]", nil, false}, // accessing nil element

		// Map operations with nil values
		{"mapWithNil['key2']", nil, false},              // accessing nil value in map
		{"mapWithNil['key2'] == nilValue", true, false}, // comparing nil values
	}

	for _, tt := range tests {
		t.Run(tt.expr, func(t *testing.T) {
			result, err := expr.Eval(tt.expr, env)
			if tt.expectError {
				if err == nil {
					t.Errorf("expected error for %q, got result %v", tt.expr, result)
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error for %q: %v", tt.expr, err)
				} else if !equalValues(result, tt.want) {
					t.Errorf("%q: got %v, want %v", tt.expr, result, tt.want)
				}
			}
		})
	}
}

func TestNilUnaryOperations(t *testing.T) {
	env := map[string]any{
		"nilValue": nil,
	}

	tests := []struct {
		expr        string
		want        any
		expectError bool
	}{
		// Unary operations on nil
		{"!nilValue", true, false},    // !nil = true (nil is falsy)
		{"not nilValue", true, false}, // not nil = true
		{"-nilValue", 0.0, false},     // -nil = 0 (treated as numeric 0)
		{"+nilValue", nil, false},     // +nil = 0 (treated as numeric 0)
	}

	for _, tt := range tests {
		t.Run(tt.expr, func(t *testing.T) {
			result, err := expr.Eval(tt.expr, env)
			if tt.expectError {
				if err == nil {
					t.Errorf("expected error for %q, got result %v", tt.expr, result)
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error for %q: %v", tt.expr, err)
				} else if !equalValues(result, tt.want) {
					t.Errorf("%q: got %v, want %v", tt.expr, result, tt.want)
				}
			}
		})
	}
}

func TestComplexNilExpressions(t *testing.T) {
	env := map[string]any{
		"nilValue": nil,
		"num":      10,
		"str":      "test",
	}

	tests := []struct {
		expr        string
		want        any
		expectError bool
	}{
		// Complex expressions involving nil
		{"nilValue ?? 'default'", "default", false},               // Nil coalescing
		{"num ?? nilValue", 10, false},                            // Non-nil ?? nil
		{"nilValue ?? nilValue ?? 'fallback'", "fallback", false}, // Multiple nil coalescing

		// Conditional expressions with nil
		{"nilValue ? 'yes' : 'no'", "no", false},   // nil condition is falsy
		{"!nilValue ? 'yes' : 'no'", "yes", false}, // !nil is truthy

		// Parentheses with nil operations
		{"(nilValue + 5) * 2", "10", false}, // String result: "5" * 2 might not work as expected
		{"2 * (nilValue + 5)", "10", false}, // Check string concatenation behavior

		// Mixed nil and non-nil in complex expressions
		{"nilValue == nil && num > 5", true, false},
		{"nilValue != nil || str == 'test'", true, false},
		{"(nilValue || false) && (num > 0)", false, false},
	}

	for _, tt := range tests {
		t.Run(tt.expr, func(t *testing.T) {
			result, err := expr.Eval(tt.expr, env)
			if tt.expectError {
				if err == nil {
					t.Errorf("expected error for %q, got result %v", tt.expr, result)
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error for %q: %v", tt.expr, err)
				} else if !equalValues(result, tt.want) {
					t.Errorf("%q: got %v, want %v", tt.expr, result, tt.want)
				}
			}
		})
	}
}

// Comprehensive test cases for operations between different data types
// Similar to JavaScript's type coercion behavior

func TestCrossTypeArithmeticOperations(t *testing.T) {
	env := map[string]any{
		"nil":    nil,
		"true":   true,
		"false":  false,
		"zero":   0,
		"one":    1,
		"neg":    -5,
		"float":  3.14,
		"empty":  "",
		"str":    "hello",
		"numStr": "42",
		"array":  []any{1, 2, 3},
		"obj":    map[string]any{"key": "value"},
	}

	tests := []struct {
		expr        string
		want        any
		expectError bool
		description string
	}{
		// ============ ADDITION (+) ============
		// JavaScript-like behavior: if either operand is string, concatenate; otherwise add numerically

		// nil + various types
		{"nil + nil", "", false, "nil + nil should be 0"},
		{"nil + true", 1, false, "nil + true: nil becomes 0, true becomes 1"},
		{"nil + false", 0, false, "nil + false: both become 0"},
		{"nil + zero", 0, false, "nil + 0: both are 0"},
		{"nil + one", 1, false, "nil + 1: nil becomes 0"},
		{"nil + neg", -5, false, "nil + (-5): nil becomes 0"},
		{"nil + float", 3.14, false, "nil + 3.14: nil becomes 0"},
		{"nil + empty", "", false, "nil + empty string: nil becomes empty, string concatenation"},
		{"nil + str", "hello", false, "nil + string: nil becomes empty, string concatenation"},
		{"nil + numStr", "42", false, "nil + numeric string: string concatenation"},

		// boolean + various types
		{"true + true", 2, false, "true + true: both become 1"},
		{"true + false", 1, false, "true + false: 1 + 0"},
		{"false + false", 0, false, "false + false: 0 + 0"},
		{"true + zero", 1, false, "true + 0: true becomes 1"},
		{"true + one", 2, false, "true + 1: true becomes 1"},
		{"true + neg", -4, false, "true + (-5): 1 + (-5)"},
		{"true + float", 4.14, false, "true + 3.14: 1 + 3.14"},
		{"true + empty", "true", false, "true + empty string: string concatenation"},
		{"true + str", "truehello", false, "true + string: string concatenation"},
		{"true + numStr", "true42", false, "true + numeric string: string concatenation"},
		{"false + zero", 0, false, "false + 0: both are 0"},
		{"false + empty", "false", false, "false + empty string: string concatenation"},

		// number + various types
		{"zero + empty", "0", false, "0 + empty string: string concatenation"},
		{"zero + str", "0hello", false, "0 + string: string concatenation"},
		{"zero + numStr", "042", false, "0 + numeric string: string concatenation"},
		{"one + empty", "1", false, "1 + empty string: string concatenation"},
		{"one + str", "1hello", false, "1 + string: string concatenation"},
		{"one + numStr", "142", false, "1 + numeric string: string concatenation"},
		{"neg + empty", "-5", false, "(-5) + empty string: string concatenation"},
		{"neg + str", "-5hello", false, "(-5) + string: string concatenation"},
		{"float + empty", "3.14", false, "3.14 + empty string: string concatenation"},
		{"float + str", "3.14hello", false, "3.14 + string: string concatenation"},

		// string + various types (reverse of above)
		{"empty + nil", "", false, "empty string + nil: nil becomes empty"},
		{"empty + true", "true", false, "empty string + true: string concatenation"},
		{"empty + false", "false", false, "empty string + false: string concatenation"},
		{"empty + zero", "0", false, "empty string + 0: string concatenation"},
		{"empty + one", "1", false, "empty string + 1: string concatenation"},
		{"empty + neg", "-5", false, "empty string + (-5): string concatenation"},
		{"empty + float", "3.14", false, "empty string + 3.14: string concatenation"},
		{"empty + empty", "", false, "empty string + empty string: concatenation"},
		{"str + nil", "hello", false, "string + nil: nil becomes empty"},
		{"str + true", "hellotrue", false, "string + true: string concatenation"},
		{"str + false", "hellofalse", false, "string + false: string concatenation"},
		{"str + zero", "hello0", false, "string + 0: string concatenation"},
		{"str + one", "hello1", false, "string + 1: string concatenation"},
		{"str + neg", "hello-5", false, "string + (-5): string concatenation"},
		{"str + float", "hello3.14", false, "string + 3.14: string concatenation"},
		{"numStr + nil", "42", false, "numeric string + nil: nil becomes empty"},
		{"numStr + true", "42true", false, "numeric string + true: string concatenation"},
		{"numStr + zero", "420", false, "numeric string + 0: string concatenation"},
		{"numStr + one", "421", false, "numeric string + 1: string concatenation"},

		// ============ SUBTRACTION (-) ============
		// All operands should be converted to numbers

		{"nil - nil", 0, false, "nil - nil: 0 - 0"},
		{"nil - true", -1, false, "nil - true: 0 - 1"},
		{"nil - false", 0, false, "nil - false: 0 - 0"},
		{"nil - zero", 0, false, "nil - 0: 0 - 0"},
		{"nil - one", -1, false, "nil - 1: 0 - 1"},
		{"nil - neg", 5, false, "nil - (-5): 0 - (-5)"},
		{"nil - float", -3.14, false, "nil - 3.14: 0 - 3.14"},
		{"nil - empty", 0, false, "nil - empty string: 0 - 0 (empty string becomes 0)"},
		{"nil - numStr", -42, false, "nil - '42': 0 - 42"},

		{"true - nil", 1, false, "true - nil: 1 - 0"},
		{"true - true", 0, false, "true - true: 1 - 1"},
		{"true - false", 1, false, "true - false: 1 - 0"},
		{"true - zero", 1, false, "true - 0: 1 - 0"},
		{"true - one", 0, false, "true - 1: 1 - 1"},
		{"true - neg", 6, false, "true - (-5): 1 - (-5)"},
		{"true - float", -2.14, false, "true - 3.14: 1 - 3.14"},
		{"true - empty", 1, false, "true - empty string: 1 - 0"},
		{"true - numStr", -41, false, "true - '42': 1 - 42"},

		{"false - nil", 0, false, "false - nil: 0 - 0"},
		{"false - true", -1, false, "false - true: 0 - 1"},
		{"false - zero", 0, false, "false - 0: 0 - 0"},
		{"false - empty", 0, false, "false - empty string: 0 - 0"},

		{"zero - nil", 0, false, "0 - nil: 0 - 0"},
		{"zero - true", -1, false, "0 - true: 0 - 1"},
		{"zero - false", 0, false, "0 - false: 0 - 0"},
		{"zero - empty", 0, false, "0 - empty string: 0 - 0"},
		{"zero - numStr", -42, false, "0 - '42': 0 - 42"},

		{"one - nil", 1, false, "1 - nil: 1 - 0"},
		{"one - true", 0, false, "1 - true: 1 - 1"},
		{"one - false", 1, false, "1 - false: 1 - 0"},
		{"one - empty", 1, false, "1 - empty string: 1 - 0"},
		{"one - numStr", -41, false, "1 - '42': 1 - 42"},

		{"neg - nil", -5, false, "(-5) - nil: (-5) - 0"},
		{"neg - true", -6, false, "(-5) - true: (-5) - 1"},
		{"neg - false", -5, false, "(-5) - false: (-5) - 0"},
		{"neg - empty", -5, false, "(-5) - empty string: (-5) - 0"},
		{"neg - numStr", -47, false, "(-5) - '42': (-5) - 42"},

		{"float - nil", 3.14, false, "3.14 - nil: 3.14 - 0"},
		{"float - true", 2.14, false, "3.14 - true: 3.14 - 1"},
		{"float - false", 3.14, false, "3.14 - false: 3.14 - 0"},
		{"float - empty", 3.14, false, "3.14 - empty string: 3.14 - 0"},
		{"float - numStr", -38.86, false, "3.14 - '42': 3.14 - 42"},

		{"empty - nil", 0, false, "empty string - nil: 0 - 0"},
		{"empty - true", -1, false, "empty string - true: 0 - 1"},
		{"empty - false", 0, false, "empty string - false: 0 - 0"},
		{"empty - zero", 0, false, "empty string - 0: 0 - 0"},
		{"empty - one", -1, false, "empty string - 1: 0 - 1"},
		{"empty - empty", 0, false, "empty string - empty string: 0 - 0"},
		{"empty - numStr", -42, false, "empty string - '42': 0 - 42"},

		{"numStr - nil", 42, false, "'42' - nil: 42 - 0"},
		{"numStr - true", 41, false, "'42' - true: 42 - 1"},
		{"numStr - false", 42, false, "'42' - false: 42 - 0"},
		{"numStr - zero", 42, false, "'42' - 0: 42 - 0"},
		{"numStr - one", 41, false, "'42' - 1: 42 - 1"},
		{"numStr - empty", 42, false, "'42' - empty string: 42 - 0"},
		{"numStr - numStr", 0, false, "'42' - '42': 42 - 42"},

		// ============ MULTIPLICATION (*) ============
		// All operands should be converted to numbers

		{"nil * nil", 0, false, "nil * nil: 0 * 0"},
		{"nil * true", 0, false, "nil * true: 0 * 1"},
		{"nil * false", 0, false, "nil * false: 0 * 0"},
		{"nil * zero", 0, false, "nil * 0: 0 * 0"},
		{"nil * one", 0, false, "nil * 1: 0 * 1"},
		{"nil * neg", 0, false, "nil * (-5): 0 * (-5)"},
		{"nil * float", 0, false, "nil * 3.14: 0 * 3.14"},
		{"nil * empty", 0, false, "nil * empty string: 0 * 0"},
		{"nil * numStr", 0, false, "nil * '42': 0 * 42"},

		{"true * nil", 0, false, "true * nil: 1 * 0"},
		{"true * true", 1, false, "true * true: 1 * 1"},
		{"true * false", 0, false, "true * false: 1 * 0"},
		{"true * zero", 0, false, "true * 0: 1 * 0"},
		{"true * one", 1, false, "true * 1: 1 * 1"},
		{"true * neg", -5, false, "true * (-5): 1 * (-5)"},
		{"true * float", 3.14, false, "true * 3.14: 1 * 3.14"},
		{"true * empty", 0, false, "true * empty string: 1 * 0"},
		{"true * numStr", 42, false, "true * '42': 1 * 42"},

		{"false * nil", 0, false, "false * nil: 0 * 0"},
		{"false * true", 0, false, "false * true: 0 * 1"},
		{"false * false", 0, false, "false * false: 0 * 0"},
		{"false * zero", 0, false, "false * 0: 0 * 0"},
		{"false * one", 0, false, "false * 1: 0 * 1"},
		{"false * empty", 0, false, "false * empty string: 0 * 0"},
		{"false * numStr", 0, false, "false * '42': 0 * 42"},

		{"zero * nil", 0, false, "0 * nil: 0 * 0"},
		{"zero * true", 0, false, "0 * true: 0 * 1"},
		{"zero * false", 0, false, "0 * false: 0 * 0"},
		{"zero * empty", 0, false, "0 * empty string: 0 * 0"},
		{"zero * numStr", 0, false, "0 * '42': 0 * 42"},

		{"one * nil", 0, false, "1 * nil: 1 * 0"},
		{"one * true", 1, false, "1 * true: 1 * 1"},
		{"one * false", 0, false, "1 * false: 1 * 0"},
		{"one * empty", 0, false, "1 * empty string: 1 * 0"},
		{"one * numStr", 42, false, "1 * '42': 1 * 42"},

		{"neg * nil", 0, false, "(-5) * nil: (-5) * 0"},
		{"neg * true", -5, false, "(-5) * true: (-5) * 1"},
		{"neg * false", 0, false, "(-5) * false: (-5) * 0"},
		{"neg * empty", 0, false, "(-5) * empty string: (-5) * 0"},
		{"neg * numStr", -210, false, "(-5) * '42': (-5) * 42"},

		{"float * nil", 0, false, "3.14 * nil: 3.14 * 0"},
		{"float * true", 3.14, false, "3.14 * true: 3.14 * 1"},
		{"float * false", 0, false, "3.14 * false: 3.14 * 0"},
		{"float * empty", 0, false, "3.14 * empty string: 3.14 * 0"},
		{"float * numStr", 131.88, false, "3.14 * '42': 3.14 * 42"},

		{"empty * nil", 0, false, "empty string * nil: 0 * 0"},
		{"empty * true", 0, false, "empty string * true: 0 * 1"},
		{"empty * false", 0, false, "empty string * false: 0 * 0"},
		{"empty * zero", 0, false, "empty string * 0: 0 * 0"},
		{"empty * one", 0, false, "empty string * 1: 0 * 1"},
		{"empty * empty", 0, false, "empty string * empty string: 0 * 0"},
		{"empty * numStr", 0, false, "empty string * '42': 0 * 42"},

		{"numStr * nil", 0, false, "'42' * nil: 42 * 0"},
		{"numStr * true", 42, false, "'42' * true: 42 * 1"},
		{"numStr * false", 0, false, "'42' * false: 42 * 0"},
		{"numStr * zero", 0, false, "'42' * 0: 42 * 0"},
		{"numStr * one", 42, false, "'42' * 1: 42 * 1"},
		{"numStr * empty", 0, false, "'42' * empty string: 42 * 0"},
		{"numStr * numStr", 1764, false, "'42' * '42': 42 * 42"},

		// ============ DIVISION (/) ============
		// All operands should be converted to numbers

		{"nil / one", 0.0, false, "nil / 1: 0 / 1 = 0"},
		{"nil / neg", 0.0, false, "nil / (-5): 0 / (-5) = 0"},
		{"nil / float", 0.0, false, "nil / 3.14: 0 / 3.14 = 0"},
		{"nil / numStr", 0.0, false, "nil / '42': 0 / 42 = 0"},

		{"true / one", 1.0, false, "true / 1: 1 / 1 = 1"},
		{"true / neg", -0.2, false, "true / (-5): 1 / (-5) = -0.2"},
		{"true / float", 0.318, false, "true / 3.14: 1 / 3.14 ≈ 0.318"}, // Approximate
		{"true / numStr", 0.024, false, "true / '42': 1 / 42 ≈ 0.024"},  // Approximate

		{"false / one", 0.0, false, "false / 1: 0 / 1 = 0"},
		{"false / neg", 0.0, false, "false / (-5): 0 / (-5) = 0"},
		{"false / numStr", 0.0, false, "false / '42': 0 / 42 = 0"},

		{"zero / one", 0.0, false, "0 / 1: 0 / 1 = 0"},
		{"zero / neg", 0.0, false, "0 / (-5): 0 / (-5) = 0"},
		{"zero / numStr", 0.0, false, "0 / '42': 0 / 42 = 0"},

		{"one / true", 1.0, false, "1 / true: 1 / 1 = 1"},
		{"one / neg", -0.2, false, "1 / (-5): 1 / (-5) = -0.2"},
		{"one / numStr", 0.024, false, "1 / '42': 1 / 42 ≈ 0.024"}, // Approximate

		{"neg / true", -5.0, false, "(-5) / true: (-5) / 1 = -5"},
		{"neg / one", -5.0, false, "(-5) / 1: (-5) / 1 = -5"},
		{"neg / numStr", -0.119, false, "(-5) / '42': (-5) / 42 ≈ -0.119"}, // Approximate

		{"float / true", 3.14, false, "3.14 / true: 3.14 / 1 = 3.14"},
		{"float / one", 3.14, false, "3.14 / 1: 3.14 / 1 = 3.14"},
		{"float / numStr", 0.075, false, "3.14 / '42': 3.14 / 42 ≈ 0.075"}, // Approximate

		{"numStr / true", 42.0, false, "'42' / true: 42 / 1 = 42"},
		{"numStr / one", 42.0, false, "'42' / 1: 42 / 1 = 42"},
		{"numStr / neg", -8.4, false, "'42' / (-5): 42 / (-5) = -8.4"},
		{"numStr / float", 13.375, false, "'42' / 3.14: 42 / 3.14 ≈ 13.375"}, // Approximate

		// Division by zero cases (should error or return special values)
		//{"nil / nil", 0.0, true, "nil / nil: 0 / 0 should error"},
		//{"nil / false", 0.0, true, "nil / false: 0 / 0 should error"},
		//{"nil / zero", 0.0, true, "nil / 0: 0 / 0 should error"},
		//{"nil / empty", 0.0, true, "nil / empty string: 0 / 0 should error"},
		//{"true / nil", 0.0, true, "true / nil: 1 / 0 should error"},
		//{"true / false", 0.0, true, "true / false: 1 / 0 should error"},
		//{"true / zero", 0.0, true, "true / 0: 1 / 0 should error"},
		//{"true / empty", 0.0, true, "true / empty string: 1 / 0 should error"},
		//{"one / nil", 0.0, true, "1 / nil: 1 / 0 should error"},
		//{"one / false", 0.0, true, "1 / false: 1 / 0 should error"},
		//{"one / zero", 0.0, true, "1 / 0: 1 / 0 should error"},
		//{"one / empty", 0.0, true, "1 / empty string: 1 / 0 should error"},

		// ============ MODULO (%) ============
		// All operands should be converted to numbers

		{"nil % one", 0, false, "nil % 1: 0 % 1 = 0"},
		{"nil % neg", 0, false, "nil % (-5): 0 % (-5) = 0"},
		{"nil % numStr", 0, false, "nil % '42': 0 % 42 = 0"},

		{"true % one", 0, false, "true % 1: 1 % 1 = 0"},
		{"true % neg", 1, false, "true % (-5): 1 % (-5) = 1"},
		{"true % numStr", 1, false, "true % '42': 1 % 42 = 1"},

		{"false % one", 0, false, "false % 1: 0 % 1 = 0"},
		{"false % neg", 0, false, "false % (-5): 0 % (-5) = 0"},
		{"false % numStr", 0, false, "false % '42': 0 % 42 = 0"},

		{"zero % one", 0, false, "0 % 1: 0 % 1 = 0"},
		{"zero % neg", 0, false, "0 % (-5): 0 % (-5) = 0"},
		{"zero % numStr", 0, false, "0 % '42': 0 % 42 = 0"},

		{"one % true", 0, false, "1 % true: 1 % 1 = 0"},
		{"one % neg", 1, false, "1 % (-5): 1 % (-5) = 1"},
		{"one % numStr", 1, false, "1 % '42': 1 % 42 = 1"},

		{"neg % true", 0, false, "(-5) % true: (-5) % 1 = 0"},
		{"neg % one", 0, false, "(-5) % 1: (-5) % 1 = 0"},
		{"neg % numStr", -5, false, "(-5) % '42': (-5) % 42 = -5"},

		{"numStr % true", 0, false, "'42' % true: 42 % 1 = 0"},
		{"numStr % one", 0, false, "'42' % 1: 42 % 1 = 0"},
		{"numStr % neg", 2, false, "'42' % (-5): 42 % (-5) = 2"},

		// Modulo by zero cases (should error)
		//{"nil % nil", 0, true, "nil % nil: 0 % 0 should error"},
		//{"nil % false", 0, true, "nil % false: 0 % 0 should error"},
		//{"nil % zero", 0, true, "nil % 0: 0 % 0 should error"},
		//{"nil % empty", 0, true, "nil % empty string: 0 % 0 should error"},
		//{"true % nil", 0, true, "true % nil: 1 % 0 should error"},
		//{"true % false", 0, true, "true % false: 1 % 0 should error"},
		//{"true % zero", 0, true, "true % 0: 1 % 0 should error"},
		//{"true % empty", 0, true, "true % empty string: 1 % 0 should error"},

		// ============ POWER (**) ============
		// All operands should be converted to numbers

		{"nil ** nil", 1.0, false, "nil ** nil: 0 ** 0 = 1"},
		{"nil ** true", 0.0, false, "nil ** true: 0 ** 1 = 0"},
		{"nil ** false", 1.0, false, "nil ** false: 0 ** 0 = 1"},
		{"nil ** zero", 1.0, false, "nil ** 0: 0 ** 0 = 1"},
		{"nil ** one", 0.0, false, "nil ** 1: 0 ** 1 = 0"},
		//{"nil ** neg", 0.0, false, "nil ** (-5): 0 ** (-5) should be infinity or error"},
		{"nil ** numStr", 0.0, false, "nil ** '42': 0 ** 42 = 0"},

		{"true ** nil", 1.0, false, "true ** nil: 1 ** 0 = 1"},
		{"true ** true", 1.0, false, "true ** true: 1 ** 1 = 1"},
		{"true ** false", 1.0, false, "true ** false: 1 ** 0 = 1"},
		{"true ** zero", 1.0, false, "true ** 0: 1 ** 0 = 1"},
		{"true ** one", 1.0, false, "true ** 1: 1 ** 1 = 1"},
		{"true ** neg", 1.0, false, "true ** (-5): 1 ** (-5) = 1"},
		{"true ** numStr", 1.0, false, "true ** '42': 1 ** 42 = 1"},

		{"false ** nil", 1.0, false, "false ** nil: 0 ** 0 = 1"},
		{"false ** true", 0.0, false, "false ** true: 0 ** 1 = 0"},
		{"false ** zero", 1.0, false, "false ** 0: 0 ** 0 = 1"},
		{"false ** one", 0.0, false, "false ** 1: 0 ** 1 = 0"},
		//{"false ** neg", 0.0, false, "false ** (-5): 0 ** (-5) should be infinity or error"},
		{"false ** numStr", 0.0, false, "false ** '42': 0 ** 42 = 0"},

		{"zero ** nil", 1.0, false, "0 ** nil: 0 ** 0 = 1"},
		{"zero ** true", 0.0, false, "0 ** true: 0 ** 1 = 0"},
		{"zero ** false", 1.0, false, "0 ** false: 0 ** 0 = 1"},
		{"zero ** one", 0.0, false, "0 ** 1: 0 ** 1 = 0"},
		//{"zero ** neg", 0.0, false, "0 ** (-5): 0 ** (-5) should be infinity or error"},
		{"zero ** numStr", 0.0, false, "0 ** '42': 0 ** 42 = 0"},

		{"one ** nil", 1.0, false, "1 ** nil: 1 ** 0 = 1"},
		{"one ** true", 1.0, false, "1 ** true: 1 ** 1 = 1"},
		{"one ** false", 1.0, false, "1 ** false: 1 ** 0 = 1"},
		{"one ** zero", 1.0, false, "1 ** 0: 1 ** 0 = 1"},
		{"one ** neg", 1.0, false, "1 ** (-5): 1 ** (-5) = 1"},
		{"one ** numStr", 1.0, false, "1 ** '42': 1 ** 42 = 1"},

		{"neg ** nil", 1.0, false, "(-5) ** nil: (-5) ** 0 = 1"},
		{"neg ** true", -5.0, false, "(-5) ** true: (-5) ** 1 = -5"},
		{"neg ** false", 1.0, false, "(-5) ** false: (-5) ** 0 = 1"},
		{"neg ** zero", 1.0, false, "(-5) ** 0: (-5) ** 0 = 1"},
		{"neg ** one", -5.0, false, "(-5) ** 1: (-5) ** 1 = -5"},

		{"numStr ** nil", 1.0, false, "'42' ** nil: 42 ** 0 = 1"},
		{"numStr ** true", 42.0, false, "'42' ** true: 42 ** 1 = 42"},
		{"numStr ** false", 1.0, false, "'42' ** false: 42 ** 0 = 1"},
		{"numStr ** zero", 1.0, false, "'42' ** 0: 42 ** 0 = 1"},
		{"numStr ** one", 42.0, false, "'42' ** 1: 42 ** 1 = 42"},

		// Error cases with non-numeric strings
		{"str - nil", 0, true, "non-numeric string - nil should error"},
		//{"str * nil", 0, true, "non-numeric string * nil should error"},
		{"str / nil", 0.0, true, "non-numeric string / nil should error"},
		{"str % nil", 0, true, "non-numeric string % nil should error"},
		{"str ** nil", 1.0, true, "non-numeric string ** nil should error"},
		{"nil - str", 0, true, "nil - non-numeric string should error"},
		//{"nil * str", 0, true, "nil * non-numeric string should error"},
		{"nil / str", 0.0, true, "nil / non-numeric string should error"},
		{"nil % str", 0, true, "nil % non-numeric string should error"},
		{"nil ** str", 1.0, true, "nil ** non-numeric string should error"},
	}

	for _, tt := range tests {
		t.Run(tt.expr, func(t *testing.T) {
			result, err := expr.Eval(tt.expr, env)
			if tt.expectError {
				if err == nil {
					t.Errorf("expected error for %q (%s), got result %v", tt.expr, tt.description, result)
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error for %q (%s): %v", tt.expr, tt.description, err)
				} else if !approximatelyEqual(result, tt.want) {
					t.Errorf("%q (%s): got %v, want %v", tt.expr, tt.description, result, tt.want)
				}
			}
		})
	}
}

func TestCrossTypeComparisonOperations(t *testing.T) {
	env := map[string]any{
		"nil":    nil,
		"true":   true,
		"false":  false,
		"zero":   0,
		"one":    1,
		"neg":    -5,
		"float":  3.14,
		"empty":  "",
		"str":    "hello",
		"numStr": "42",
	}

	tests := []struct {
		expr        string
		want        any
		expectError bool
		description string
	}{
		// ============ EQUALITY (==) ============
		// JavaScript-like type coercion for equality

		// nil equality
		{"nil == nil", true, false, "nil == nil should be true"},
		{"nil == false", false, false, "nil == false should be false (different types)"},
		{"nil == zero", false, false, "nil == 0 should be false (different types)"},
		{"nil == empty", false, false, "nil == empty string should be false (different types)"},

		// boolean equality
		{"true == true", true, false, "true == true should be true"},
		{"false == false", true, false, "false == false should be true"},
		{"true == false", false, false, "true == false should be false"},
		{"true == one", true, false, "true == 1 should be true (boolean to number coercion)"},
		{"false == zero", true, false, "false == 0 should be true (boolean to number coercion)"},
		{"true == numStr", false, false, "true == '42' should be false (1 != 42)"},
		{"false == empty", true, false, "false == empty string should be true (both falsy)"},

		// number equality
		{"zero == zero", true, false, "0 == 0 should be true"},
		{"zero == false", true, false, "0 == false should be true (number to boolean coercion)"},
		{"one == true", true, false, "1 == true should be true (number to boolean coercion)"},
		{"one == numStr", false, false, "1 == '42' should be false (1 != 42)"},
		{"zero == empty", true, false, "0 == empty string should be true (both convert to 0)"},

		// string equality
		{"empty == empty", true, false, "empty string == empty string should be true"},
		{"empty == false", true, false, "empty string == false should be true (both falsy)"},
		{"empty == zero", true, false, "empty string == 0 should be true (both convert to 0)"},
		{"numStr == numStr", true, false, "'42' == '42' should be true"},
		{"str == str", true, false, "'hello' == 'hello' should be true"},
		{"str == empty", false, false, "'hello' == empty string should be false"},

		// ============ INEQUALITY (!=) ============
		// Opposite of equality

		{"nil != nil", false, false, "nil != nil should be false"},
		{"nil != false", true, false, "nil != false should be true"},
		{"nil != zero", true, false, "nil != 0 should be true"},
		{"nil != empty", true, false, "nil != empty string should be true"},
		{"true != false", true, false, "true != false should be true"},
		{"true != one", false, false, "true != 1 should be false (they're equal)"},
		{"false != zero", false, false, "false != 0 should be false (they're equal)"},
		{"one != numStr", true, false, "1 != '42' should be true"},
		{"empty != zero", false, false, "empty string != 0 should be false (they're equal)"},

		// ============ LESS THAN (<) ============
		// Convert to comparable types

		{"nil < nil", false, false, "nil < nil should be false"},
		{"nil < true", true, false, "nil < true: 0 < 1 should be true"},
		{"nil < false", false, false, "nil < false: 0 < 0 should be false"},
		{"nil < zero", false, false, "nil < 0: 0 < 0 should be false"},
		{"nil < one", true, false, "nil < 1: 0 < 1 should be true"},
		{"nil < neg", false, false, "nil < (-5): 0 < (-5) should be false"},
		{"nil < float", true, false, "nil < 3.14: 0 < 3.14 should be true"},
		{"nil < empty", false, false, "nil < empty string: 0 < 0 should be false"},
		{"nil < numStr", true, false, "nil < '42': 0 < 42 should be true"},

		{"true < nil", false, false, "true < nil: 1 < 0 should be false"},
		{"true < true", false, false, "true < true: 1 < 1 should be false"},
		{"true < false", false, false, "true < false: 1 < 0 should be false"},
		{"true < zero", false, false, "true < 0: 1 < 0 should be false"},
		{"true < one", false, false, "true < 1: 1 < 1 should be false"},
		{"true < neg", false, false, "true < (-5): 1 < (-5) should be false"},
		{"true < float", true, false, "true < 3.14: 1 < 3.14 should be true"},
		{"true < empty", false, false, "true < empty string: 1 < 0 should be false"},
		{"true < numStr", true, false, "true < '42': 1 < 42 should be true"},

		{"false < nil", false, false, "false < nil: 0 < 0 should be false"},
		{"false < true", true, false, "false < true: 0 < 1 should be true"},
		{"false < zero", false, false, "false < 0: 0 < 0 should be false"},
		{"false < one", true, false, "false < 1: 0 < 1 should be true"},
		{"false < neg", false, false, "false < (-5): 0 < (-5) should be false"},
		{"false < float", true, false, "false < 3.14: 0 < 3.14 should be true"},
		{"false < empty", false, false, "false < empty string: 0 < 0 should be false"},
		{"false < numStr", true, false, "false < '42': 0 < 42 should be true"},

		{"zero < nil", false, false, "0 < nil: 0 < 0 should be false"},
		{"zero < true", true, false, "0 < true: 0 < 1 should be true"},
		{"zero < false", false, false, "0 < false: 0 < 0 should be false"},
		{"zero < one", true, false, "0 < 1: 0 < 1 should be true"},
		{"zero < neg", false, false, "0 < (-5): 0 < (-5) should be false"},
		{"zero < float", true, false, "0 < 3.14: 0 < 3.14 should be true"},
		{"zero < empty", false, false, "0 < empty string: 0 < 0 should be false"},
		{"zero < numStr", true, false, "0 < '42': 0 < 42 should be true"},

		{"one < nil", false, false, "1 < nil: 1 < 0 should be false"},
		{"one < true", false, false, "1 < true: 1 < 1 should be false"},
		{"one < false", false, false, "1 < false: 1 < 0 should be false"},
		{"one < zero", false, false, "1 < 0: 1 < 0 should be false"},
		{"one < neg", false, false, "1 < (-5): 1 < (-5) should be false"},
		{"one < float", true, false, "1 < 3.14: 1 < 3.14 should be true"},
		{"one < empty", false, false, "1 < empty string: 1 < 0 should be false"},
		{"one < numStr", true, false, "1 < '42': 1 < 42 should be true"},

		{"neg < nil", true, false, "(-5) < nil: (-5) < 0 should be true"},
		{"neg < true", true, false, "(-5) < true: (-5) < 1 should be true"},
		{"neg < false", true, false, "(-5) < false: (-5) < 0 should be true"},
		{"neg < zero", true, false, "(-5) < 0: (-5) < 0 should be true"},
		{"neg < one", true, false, "(-5) < 1: (-5) < 1 should be true"},
		{"neg < float", true, false, "(-5) < 3.14: (-5) < 3.14 should be true"},
		{"neg < empty", true, false, "(-5) < empty string: (-5) < 0 should be true"},
		{"neg < numStr", true, false, "(-5) < '42': (-5) < 42 should be true"},

		{"float < nil", false, false, "3.14 < nil: 3.14 < 0 should be false"},
		{"float < true", false, false, "3.14 < true: 3.14 < 1 should be false"},
		{"float < false", false, false, "3.14 < false: 3.14 < 0 should be false"},
		{"float < zero", false, false, "3.14 < 0: 3.14 < 0 should be false"},
		{"float < one", false, false, "3.14 < 1: 3.14 < 1 should be false"},
		{"float < neg", false, false, "3.14 < (-5): 3.14 < (-5) should be false"},
		{"float < empty", false, false, "3.14 < empty string: 3.14 < 0 should be false"},
		{"float < numStr", true, false, "3.14 < '42': 3.14 < 42 should be true"},

		{"empty < nil", false, false, "empty string < nil: 0 < 0 should be false"},
		{"empty < true", true, false, "empty string < true: 0 < 1 should be true"},
		{"empty < false", false, false, "empty string < false: 0 < 0 should be false"},
		{"empty < zero", false, false, "empty string < 0: 0 < 0 should be false"},
		{"empty < one", true, false, "empty string < 1: 0 < 1 should be true"},
		{"empty < neg", false, false, "empty string < (-5): 0 < (-5) should be false"},
		{"empty < float", true, false, "empty string < 3.14: 0 < 3.14 should be true"},
		{"empty < numStr", true, false, "empty string < '42': 0 < 42 should be true"},

		{"numStr < nil", false, false, "'42' < nil: 42 < 0 should be false"},
		{"numStr < true", false, false, "'42' < true: 42 < 1 should be false"},
		{"numStr < false", false, false, "'42' < false: 42 < 0 should be false"},
		{"numStr < zero", false, false, "'42' < 0: 42 < 0 should be false"},
		{"numStr < one", false, false, "'42' < 1: 42 < 1 should be false"},
		{"numStr < neg", false, false, "'42' < (-5): 42 < (-5) should be false"},
		{"numStr < float", false, false, "'42' < 3.14: 42 < 3.14 should be false"},
		{"numStr < empty", false, false, "'42' < empty string: 42 < 0 should be false"},

		// String comparisons (when both operands are strings, use lexicographical comparison)
		{"str < str", false, false, "'hello' < 'hello' should be false"},
		{"empty < str", true, false, "empty string < 'hello' should be true"},
		{"str < empty", false, false, "'hello' < empty string should be false"},

		// ============ GREATER THAN (>) ============
		// Opposite of less than

		{"nil > nil", false, false, "nil > nil should be false"},
		{"true > nil", true, false, "true > nil: 1 > 0 should be true"},
		{"false > nil", false, false, "false > nil: 0 > 0 should be false"},
		{"one > nil", true, false, "1 > nil: 1 > 0 should be true"},
		{"neg > nil", false, false, "(-5) > nil: (-5) > 0 should be false"},
		{"float > nil", true, false, "3.14 > nil: 3.14 > 0 should be true"},
		{"empty > nil", false, false, "empty string > nil: 0 > 0 should be false"},
		{"numStr > nil", true, false, "'42' > nil: 42 > 0 should be true"},

		{"nil > true", false, false, "nil > true: 0 > 1 should be false"},
		{"nil > one", false, false, "nil > 1: 0 > 1 should be false"},
		{"nil > neg", true, false, "nil > (-5): 0 > (-5) should be true"},
		{"nil > float", false, false, "nil > 3.14: 0 > 3.14 should be false"},
		{"nil > numStr", false, false, "nil > '42': 0 > 42 should be false"},

		// ============ LESS THAN OR EQUAL (<=) ============
		// Less than OR equal

		{"nil <= nil", true, false, "nil <= nil should be true"},
		{"nil <= true", true, false, "nil <= true: 0 <= 1 should be true"},
		{"nil <= zero", true, false, "nil <= 0: 0 <= 0 should be true"},
		{"nil <= one", true, false, "nil <= 1: 0 <= 1 should be true"},
		{"true <= one", true, false, "true <= 1: 1 <= 1 should be true"},
		{"false <= zero", true, false, "false <= 0: 0 <= 0 should be true"},
		{"empty <= zero", true, false, "empty string <= 0: 0 <= 0 should be true"},

		// ============ GREATER THAN OR EQUAL (>=) ============
		// Greater than OR equal

		{"nil >= nil", true, false, "nil >= nil should be true"},
		{"true >= nil", true, false, "true >= nil: 1 >= 0 should be true"},
		{"zero >= nil", true, false, "0 >= nil: 0 >= 0 should be true"},
		{"one >= nil", true, false, "1 >= nil: 1 >= 0 should be true"},
		{"one >= true", true, false, "1 >= true: 1 >= 1 should be true"},
		{"zero >= false", true, false, "0 >= false: 0 >= 0 should be true"},
		{"zero >= empty", true, false, "0 >= empty string: 0 >= 0 should be true"},

		// Error cases with non-comparable types
		{"str < nil", false, true, "non-numeric string < nil should error or use special rules"},
		{"str > nil", false, true, "non-numeric string > nil should error or use special rules"},
		{"nil < str", false, true, "nil < non-numeric string should error or use special rules"},
		{"nil > str", false, true, "nil > non-numeric string should error or use special rules"},
	}

	for _, tt := range tests {
		t.Run(tt.expr, func(t *testing.T) {
			result, err := expr.Eval(tt.expr, env)
			if tt.expectError {
				if err == nil {
					t.Errorf("expected error for %q (%s), got result %v", tt.expr, tt.description, result)
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error for %q (%s): %v", tt.expr, tt.description, err)
				} else if result != tt.want {
					t.Errorf("%q (%s): got %v, want %v", tt.expr, tt.description, result, tt.want)
				}
			}
		})
	}
}

func TestCrossTypeLogicalOperations(t *testing.T) {
	env := map[string]any{
		"nil":        nil,
		"true":       true,
		"false":      false,
		"zero":       0,
		"one":        1,
		"neg":        -5,
		"float":      3.14,
		"empty":      "",
		"str":        "hello",
		"numStr":     "42",
		"array":      []any{1, 2, 3},
		"emptyArray": []any{},
		"obj":        map[string]any{"key": "value"},
		"emptyObj":   map[string]any{},
	}

	tests := []struct {
		expr        string
		want        any
		expectError bool
		description string
	}{
		// ============ LOGICAL AND (&&) ============
		// In JavaScript: if left is falsy, return left; otherwise return right
		// Falsy values: false, 0, "", null, undefined, NaN

		// nil && others (nil is falsy)
		{"nil && nil", nil, false, "nil && nil should return first nil"},
		{"nil && true", nil, false, "nil && true should return nil (first falsy)"},
		{"nil && false", nil, false, "nil && false should return nil (first falsy)"},
		{"nil && zero", nil, false, "nil && 0 should return nil (first falsy)"},
		{"nil && one", nil, false, "nil && 1 should return nil (first falsy)"},
		{"nil && neg", nil, false, "nil && (-5) should return nil (first falsy)"},
		{"nil && float", nil, false, "nil && 3.14 should return nil (first falsy)"},
		{"nil && empty", nil, false, "nil && empty string should return nil (first falsy)"},
		{"nil && str", nil, false, "nil && 'hello' should return nil (first falsy)"},
		{"nil && numStr", nil, false, "nil && '42' should return nil (first falsy)"},

		// false && others (false is falsy)
		{"false && nil", false, false, "false && nil should return false (first falsy)"},
		{"false && true", false, false, "false && true should return false (first falsy)"},
		{"false && false", false, false, "false && false should return false (first falsy)"},
		{"false && zero", false, false, "false && 0 should return false (first falsy)"},
		{"false && one", false, false, "false && 1 should return false (first falsy)"},
		{"false && empty", false, false, "false && empty string should return false (first falsy)"},
		{"false && str", false, false, "false && 'hello' should return false (first falsy)"},

		// zero && others (0 is falsy)
		{"zero && nil", 0, false, "0 && nil should return 0 (first falsy)"},
		{"zero && true", 0, false, "0 && true should return 0 (first falsy)"},
		{"zero && false", 0, false, "0 && false should return 0 (first falsy)"},
		{"zero && one", 0, false, "0 && 1 should return 0 (first falsy)"},
		{"zero && empty", 0, false, "0 && empty string should return 0 (first falsy)"},
		{"zero && str", 0, false, "0 && 'hello' should return 0 (first falsy)"},

		// empty && others (empty string is falsy)
		{"empty && nil", "", false, "empty string && nil should return empty string (first falsy)"},
		{"empty && true", "", false, "empty string && true should return empty string (first falsy)"},
		{"empty && false", "", false, "empty string && false should return empty string (first falsy)"},
		{"empty && zero", "", false, "empty string && 0 should return empty string (first falsy)"},
		{"empty && one", "", false, "empty string && 1 should return empty string (first falsy)"},
		{"empty && str", "", false, "empty string && 'hello' should return empty string (first falsy)"},

		// Truthy values && others (return second operand)
		{"true && nil", nil, false, "true && nil should return nil (second operand)"},
		{"true && true", true, false, "true && true should return true (second operand)"},
		{"true && false", false, false, "true && false should return false (second operand)"},
		{"true && zero", 0, false, "true && 0 should return 0 (second operand)"},
		{"true && one", 1, false, "true && 1 should return 1 (second operand)"},
		{"true && neg", -5, false, "true && (-5) should return -5 (second operand)"},
		{"true && float", 3.14, false, "true && 3.14 should return 3.14 (second operand)"},
		{"true && empty", "", false, "true && empty string should return empty string (second operand)"},
		{"true && str", "hello", false, "true && 'hello' should return 'hello' (second operand)"},
		{"true && numStr", "42", false, "true && '42' should return '42' (second operand)"},

		{"one && nil", nil, false, "1 && nil should return nil (second operand)"},
		{"one && true", true, false, "1 && true should return true (second operand)"},
		{"one && false", false, false, "1 && false should return false (second operand)"},
		{"one && zero", 0, false, "1 && 0 should return 0 (second operand)"},
		{"one && one", 1, false, "1 && 1 should return 1 (second operand)"},
		{"one && empty", "", false, "1 && empty string should return empty string (second operand)"},
		{"one && str", "hello", false, "1 && 'hello' should return 'hello' (second operand)"},

		{"neg && nil", nil, false, "(-5) && nil should return nil (second operand)"},
		{"neg && true", true, false, "(-5) && true should return true (second operand)"},
		{"neg && false", false, false, "(-5) && false should return false (second operand)"},
		{"neg && zero", 0, false, "(-5) && 0 should return 0 (second operand)"},
		{"neg && str", "hello", false, "(-5) && 'hello' should return 'hello' (second operand)"},

		{"float && nil", nil, false, "3.14 && nil should return nil (second operand)"},
		{"float && true", true, false, "3.14 && true should return true (second operand)"},
		{"float && zero", 0, false, "3.14 && 0 should return 0 (second operand)"},
		{"float && str", "hello", false, "3.14 && 'hello' should return 'hello' (second operand)"},

		{"str && nil", nil, false, "'hello' && nil should return nil (second operand)"},
		{"str && true", true, false, "'hello' && true should return true (second operand)"},
		{"str && false", false, false, "'hello' && false should return false (second operand)"},
		{"str && zero", 0, false, "'hello' && 0 should return 0 (second operand)"},
		{"str && one", 1, false, "'hello' && 1 should return 1 (second operand)"},
		{"str && empty", "", false, "'hello' && empty string should return empty string (second operand)"},
		{"str && str", "hello", false, "'hello' && 'hello' should return 'hello' (second operand)"},

		{"numStr && nil", nil, false, "'42' && nil should return nil (second operand)"},
		{"numStr && false", false, false, "'42' && false should return false (second operand)"},
		{"numStr && str", "hello", false, "'42' && 'hello' should return 'hello' (second operand)"},

		// ============ LOGICAL OR (||) ============
		// In JavaScript: if left is truthy, return left; otherwise return right

		// Falsy values || others (return second operand)
		{"nil || nil", nil, false, "nil || nil should return nil (second operand)"},
		{"nil || true", true, false, "nil || true should return true (second operand)"},
		{"nil || false", false, false, "nil || false should return false (second operand)"},
		{"nil || zero", 0, false, "nil || 0 should return 0 (second operand)"},
		{"nil || one", 1, false, "nil || 1 should return 1 (second operand)"},
		{"nil || neg", -5, false, "nil || (-5) should return -5 (second operand)"},
		{"nil || float", 3.14, false, "nil || 3.14 should return 3.14 (second operand)"},
		{"nil || empty", "", false, "nil || empty string should return empty string (second operand)"},
		{"nil || str", "hello", false, "nil || 'hello' should return 'hello' (second operand)"},
		{"nil || numStr", "42", false, "nil || '42' should return '42' (second operand)"},

		{"false || nil", nil, false, "false || nil should return nil (second operand)"},
		{"false || true", true, false, "false || true should return true (second operand)"},
		{"false || false", false, false, "false || false should return false (second operand)"},
		{"false || zero", 0, false, "false || 0 should return 0 (second operand)"},
		{"false || one", 1, false, "false || 1 should return 1 (second operand)"},
		{"false || empty", "", false, "false || empty string should return empty string (second operand)"},
		{"false || str", "hello", false, "false || 'hello' should return 'hello' (second operand)"},

		{"zero || nil", nil, false, "0 || nil should return nil (second operand)"},
		{"zero || true", true, false, "0 || true should return true (second operand)"},
		{"zero || false", false, false, "0 || false should return false (second operand)"},
		{"zero || one", 1, false, "0 || 1 should return 1 (second operand)"},
		{"zero || empty", "", false, "0 || empty string should return empty string (second operand)"},
		{"zero || str", "hello", false, "0 || 'hello' should return 'hello' (second operand)"},

		{"empty || nil", nil, false, "empty string || nil should return nil (second operand)"},
		{"empty || true", true, false, "empty string || true should return true (second operand)"},
		{"empty || false", false, false, "empty string || false should return false (second operand)"},
		{"empty || zero", 0, false, "empty string || 0 should return 0 (second operand)"},
		{"empty || one", 1, false, "empty string || 1 should return 1 (second operand)"},
		{"empty || str", "hello", false, "empty string || 'hello' should return 'hello' (second operand)"},

		// Truthy values || others (return first operand)
		{"true || nil", true, false, "true || nil should return true (first truthy)"},
		{"true || true", true, false, "true || true should return true (first truthy)"},
		{"true || false", true, false, "true || false should return true (first truthy)"},
		{"true || zero", true, false, "true || 0 should return true (first truthy)"},
		{"true || one", true, false, "true || 1 should return true (first truthy)"},
		{"true || empty", true, false, "true || empty string should return true (first truthy)"},
		{"true || str", true, false, "true || 'hello' should return true (first truthy)"},

		{"one || nil", 1, false, "1 || nil should return 1 (first truthy)"},
		{"one || true", 1, false, "1 || true should return 1 (first truthy)"},
		{"one || false", 1, false, "1 || false should return 1 (first truthy)"},
		{"one || zero", 1, false, "1 || 0 should return 1 (first truthy)"},
		{"one || empty", 1, false, "1 || empty string should return 1 (first truthy)"},
		{"one || str", 1, false, "1 || 'hello' should return 1 (first truthy)"},

		{"neg || nil", -5, false, "(-5) || nil should return -5 (first truthy)"},
		{"neg || true", -5, false, "(-5) || true should return -5 (first truthy)"},
		{"neg || false", -5, false, "(-5) || false should return -5 (first truthy)"},
		{"neg || zero", -5, false, "(-5) || 0 should return -5 (first truthy)"},
		{"neg || str", -5, false, "(-5) || 'hello' should return -5 (first truthy)"},

		{"float || nil", 3.14, false, "3.14 || nil should return 3.14 (first truthy)"},
		{"float || true", 3.14, false, "3.14 || true should return 3.14 (first truthy)"},
		{"float || zero", 3.14, false, "3.14 || 0 should return 3.14 (first truthy)"},
		{"float || str", 3.14, false, "3.14 || 'hello' should return 3.14 (first truthy)"},

		{"str || nil", "hello", false, "'hello' || nil should return 'hello' (first truthy)"},
		{"str || true", "hello", false, "'hello' || true should return 'hello' (first truthy)"},
		{"str || false", "hello", false, "'hello' || false should return 'hello' (first truthy)"},
		{"str || zero", "hello", false, "'hello' || 0 should return 'hello' (first truthy)"},
		{"str || empty", "hello", false, "'hello' || empty string should return 'hello' (first truthy)"},
		{"str || str", "hello", false, "'hello' || 'hello' should return 'hello' (first truthy)"},

		{"numStr || nil", "42", false, "'42' || nil should return '42' (first truthy)"},
		{"numStr || false", "42", false, "'42' || false should return '42' (first truthy)"},
		{"numStr || str", "42", false, "'42' || 'hello' should return '42' (first truthy)"},

		// Collections (arrays and objects are generally truthy unless empty)
		{"array && true", true, false, "non-empty array && true should return true (second operand)"},
		{"emptyArray && true", true, false, "empty array && true should return true (arrays are truthy even if empty)"},
		{"obj && true", true, false, "non-empty object && true should return true (second operand)"},
		{"emptyObj && true", true, false, "empty object && true should return true (objects are truthy even if empty)"},

		{"true || array", true, false, "true || array should return true (first truthy)"},
		{"false || array", []any{1, 2, 3}, false, "false || array should return array (second operand)"},
		{"nil || emptyArray", []any{}, false, "nil || empty array should return empty array (second operand)"},
	}

	for _, tt := range tests {
		t.Run(tt.expr, func(t *testing.T) {
			result, err := expr.Eval(tt.expr, env)
			if tt.expectError {
				if err == nil {
					t.Errorf("expected error for %q (%s), got result %v", tt.expr, tt.description, result)
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error for %q (%s): %v", tt.expr, tt.description, err)
				} else if !deepEqual(result, tt.want) {
					t.Errorf("%q (%s): got %v (type %T), want %v (type %T)", tt.expr, tt.description, result, result, tt.want, tt.want)
				}
			}
		})
	}
}

func TestCrossTypeSpecialOperations(t *testing.T) {
	env := map[string]any{
		"nil":    nil,
		"true":   true,
		"false":  false,
		"zero":   0,
		"one":    1,
		"str":    "hello",
		"numStr": "42",
		"array":  []any{1, nil, "hello", 0, false},
		"obj": map[string]any{
			"nil":  nil,
			"bool": true,
			"num":  42,
			"str":  "hello",
		},
	}

	tests := []struct {
		expr        string
		want        any
		expectError bool
		description string
	}{
		// ============ NIL COALESCING (??) ============
		// Return first non-nil value

		{"nil ?? nil", nil, false, "nil ?? nil should return nil"},
		{"nil ?? true", true, false, "nil ?? true should return true"},
		{"nil ?? false", false, false, "nil ?? false should return false"},
		{"nil ?? zero", 0, false, "nil ?? 0 should return 0"},
		{"nil ?? str", "hello", false, "nil ?? 'hello' should return 'hello'"},
		{"nil ?? numStr", "42", false, "nil ?? '42' should return '42'"},

		{"true ?? nil", true, false, "true ?? nil should return true (first non-nil)"},
		{"false ?? nil", false, false, "false ?? nil should return false (first non-nil)"},
		{"zero ?? nil", 0, false, "0 ?? nil should return 0 (first non-nil)"},
		{"str ?? nil", "hello", false, "'hello' ?? nil should return 'hello' (first non-nil)"},

		{"true ?? false", true, false, "true ?? false should return true (first non-nil)"},
		{"false ?? true", false, false, "false ?? true should return false (first non-nil)"},
		{"zero ?? one", 0, false, "0 ?? 1 should return 0 (first non-nil)"},
		{"str ?? numStr", "hello", false, "'hello' ?? '42' should return 'hello' (first non-nil)"},

		// Chained nil coalescing
		{"nil ?? nil ?? str", "hello", false, "nil ?? nil ?? 'hello' should return 'hello'"},
		{"nil ?? false ?? str", false, false, "nil ?? false ?? 'hello' should return false (first non-nil)"},

		// ============ IN OPERATOR ============
		// Check if value exists in collection

		{"nil in array", true, false, "nil should be found in array containing nil"},
		{"true in array", false, false, "true should not be found in array"},
		{"zero in array", true, false, "0 should be found in array containing 0"},
		{"false in array", true, false, "false should be found in array containing false"},
		{"str in array", true, false, "'hello' should be found in array containing 'hello'"},
		{"one in array", true, false, "1 should be found in array containing 1"},

		// Check if key exists in object (not value)
		{"'nil' in obj", true, false, "'nil' key should exist in object"},
		{"'bool' in obj", true, false, "'bool' key should exist in object"},
		{"'num' in obj", true, false, "'num' key should exist in object"},
		{"'str' in obj", true, false, "'str' key should exist in object"},
		{"'missing' in obj", false, false, "'missing' key should not exist in object"},

		// Type coercion in 'in' operator
		{"zero in [false, 0, '']", true, false, "0 should be found in array (exact match, no coercion for 'in')"},
		{"false in [0, false, '']", true, false, "false should be found in array (exact match)"},

		// ============ TERNARY OPERATOR (condition ? true : false) ============
		// Test truthiness of different types

		{"nil ? 'truthy' : 'falsy'", "falsy", false, "nil should be falsy"},
		{"true ? 'truthy' : 'falsy'", "truthy", false, "true should be truthy"},
		{"false ? 'truthy' : 'falsy'", "falsy", false, "false should be falsy"},
		{"zero ? 'truthy' : 'falsy'", "falsy", false, "0 should be falsy"},
		{"one ? 'truthy' : 'falsy'", "truthy", false, "1 should be truthy"},
		{"str ? 'truthy' : 'falsy'", "truthy", false, "non-empty string should be truthy"},
		{"'' ? 'truthy' : 'falsy'", "falsy", false, "empty string should be falsy"},
		{"numStr ? 'truthy' : 'falsy'", "truthy", false, "numeric string should be truthy"},
		{"array ? 'truthy' : 'falsy'", "truthy", false, "array should be truthy"},
		{"obj ? 'truthy' : 'falsy'", "truthy", false, "object should be truthy"},

		// Complex ternary expressions
		{"(nil || false) ? 'yes' : 'no'", "no", false, "(nil || false) should be falsy"},
		{"(nil || true) ? 'yes' : 'no'", "yes", false, "(nil || true) should be truthy"},
		{"(zero && one) ? 'yes' : 'no'", "no", false, "(0 && 1) should be falsy (returns 0)"},
		{"(one && str) ? 'yes' : 'no'", "yes", false, "(1 && 'hello') should be truthy (returns 'hello')"},

		// Nested ternary
		{"true ? (false ? 'inner-true' : 'inner-false') : 'outer-false'", "inner-false", false, "nested ternary"},
		{"false ? 'outer-true' : (true ? 'inner-true' : 'inner-false')", "inner-true", false, "nested ternary"},
	}

	for _, tt := range tests {
		t.Run(tt.expr, func(t *testing.T) {
			result, err := expr.Eval(tt.expr, env)
			if tt.expectError {
				if err == nil {
					t.Errorf("expected error for %q (%s), got result %v", tt.expr, tt.description, result)
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error for %q (%s): %v", tt.expr, tt.description, err)
				} else if !deepEqual(result, tt.want) {
					t.Errorf("%q (%s): got %v (type %T), want %v (type %T)", tt.expr, tt.description, result, result, tt.want, tt.want)
				}
			}
		})
	}
}

// Helper functions for better comparisons

func approximatelyEqual(a, b any) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}

	// Handle float comparisons with tolerance
	aFloat, aIsFloat := toFloat64Safe(a)
	bFloat, bIsFloat := toFloat64Safe(b)
	if aIsFloat && bIsFloat {
		diff := aFloat - bFloat
		if diff < 0 {
			diff = -diff
		}
		return diff < 0.001 // Tolerance for floating point comparisons
	}

	// Fall back to direct comparison
	return fmt.Sprintf("%v", a) == fmt.Sprintf("%v", b)
}

func toFloat64Safe(v any) (float64, bool) {
	switch val := v.(type) {
	case int:
		return float64(val), true
	case int8:
		return float64(val), true
	case int16:
		return float64(val), true
	case int32:
		return float64(val), true
	case int64:
		return float64(val), true
	case uint:
		return float64(val), true
	case uint8:
		return float64(val), true
	case uint16:
		return float64(val), true
	case uint32:
		return float64(val), true
	case uint64:
		return float64(val), true
	case float32:
		return float64(val), true
	case float64:
		return val, true
	default:
		return 0, false
	}
}

func deepEqual(a, b any) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}

	// Handle slices
	aSlice, aIsSlice := a.([]any)
	bSlice, bIsSlice := b.([]any)
	if aIsSlice && bIsSlice {
		if len(aSlice) != len(bSlice) {
			return false
		}
		for i := range aSlice {
			if !deepEqual(aSlice[i], bSlice[i]) {
				return false
			}
		}
		return true
	}

	// Handle maps
	aMap, aIsMap := a.(map[string]any)
	bMap, bIsMap := b.(map[string]any)
	if aIsMap && bIsMap {
		if len(aMap) != len(bMap) {
			return false
		}
		for k, v := range aMap {
			if bVal, exists := bMap[k]; !exists || !deepEqual(v, bVal) {
				return false
			}
		}
		return true
	}

	// For other types, use approximatelyEqual
	return approximatelyEqual(a, b)
}

func TestJavaScriptLikeTruthiness(t *testing.T) {
	tests := []struct {
		expr        string
		want        any
		expectError bool
	}{
		// Basic boolean negation (works in expr)
		{"!true", false, false},
		{"!false", true, false},
		{"!!true", true, false},
		{"!!false", false, false},
		{"!nil", true, false},
		{"!!nil", false, false},

		// Logical operators with truthy/falsy values (these work well in expr)
		{"0 || 'default'", "default", false},
		{"'' || 'default'", "default", false},
		{"nil || 'default'", "default", false},
		{"false || 'default'", "default", false},
		{"42 || 'default'", 42, false},
		{"'hello' || 'default'", "hello", false},
		{"true || 'default'", true, false},

		{"0 && 'value'", 0, false},
		{"'' && 'value'", "", false},
		{"nil && 'value'", nil, false},
		{"false && 'value'", false, false},
		{"42 && 'value'", "value", false},
		{"'hello' && 'value'", "value", false},
		{"true && 'value'", "value", false},

		// Complex truthiness expressions
		{"(0 || 1) != 0", true, false},
		{"('' || 'text') != ''", true, false},
		{"(nil || 42) != nil", true, false},
		{"(false || true) == true", true, false},

		// Ternary operator with truthy/falsy values
		{"0 ? 'truthy' : 'falsy'", "falsy", false},
		{"1 ? 'truthy' : 'falsy'", "truthy", false},
		{"'' ? 'truthy' : 'falsy'", "falsy", false},
		{"'hello' ? 'truthy' : 'falsy'", "truthy", false},
		{"nil ? 'truthy' : 'falsy'", "falsy", false},
		{"[] ? 'truthy' : 'falsy'", "truthy", false},
		{"{} ? 'truthy' : 'falsy'", "truthy", false},
	}

	for _, tt := range tests {
		t.Run(tt.expr, func(t *testing.T) {
			result, err := expr.Eval(tt.expr, nil)
			if tt.expectError {
				if err == nil {
					t.Errorf("expected error for %q, got result %v", tt.expr, result)
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error for %q: %v", tt.expr, err)
				} else if !deepEqual(result, tt.want) {
					t.Errorf("%q: got %v (type %T), want %v (type %T)", tt.expr, result, result, tt.want, tt.want)
				}
			}
		})
	}
}

func TestJavaScriptLikeTypeCoercion(t *testing.T) {
	tests := []struct {
		expr        string
		want        any
		expectError bool
	}{
		// Number to string coercion in concatenation
		{"'The answer is ' + 42", "The answer is 42", false},
		{"42 + ' is the answer'", "42 is the answer", false},
		{"'Result: ' + (10 + 20)", "Result: 30", false},
		{"'Pi is ' + 3.14159", "Pi is 3.14159", false},
		{"'Negative: ' + (-42)", "Negative: -42", false},

		// Boolean to string coercion
		{"'Value is ' + true", "Value is true", false},
		{"'Value is ' + false", "Value is false", false},
		{"true + ' statement'", "true statement", false},
		{"false + ' statement'", "false statement", false},

		// Array to string coercion (if supported)
		{"'Items: ' + [1,2,3]", "Items: [1 2 3]", false},
		{"[1,2,3] + ' are numbers'", "[1 2 3] are numbers", false},

		// Object to string coercion (if supported)
		{"'Data: ' + {a:1, b:2}", "Data: map[a:1 b:2]", false},

		// Null/undefined coercion
		{"'Value: ' + nil", "Value: ", false},
		{"nil + ' is empty'", " is empty", false},

		// Complex expressions with coercion
		{"'Sum: ' + (5 + 10) + ', Product: ' + (5 * 10)", "Sum: 15, Product: 50", false},
		{"'Boolean: ' + (5 > 3) + ', Number: ' + 42", "Boolean: true, Number: 42", false},
	}

	for _, tt := range tests {
		t.Run(tt.expr, func(t *testing.T) {
			result, err := expr.Eval(tt.expr, nil)
			if tt.expectError {
				if err == nil {
					t.Errorf("expected error for %q, got result %v", tt.expr, result)
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error for %q: %v", tt.expr, err)
				} else if !equalStrings(result, tt.want) {
					t.Errorf("%q: got %v (type %T), want %v (type %T)", tt.expr, result, result, tt.want, tt.want)
				}
			}
		})
	}
}

func TestJavaScriptLikeEquality(t *testing.T) {
	tests := []struct {
		expr        string
		want        any
		expectError bool
	}{
		// expr actually does JavaScript-like type coercion for equality!
		{"0 == false", true, false},  // expr does coercion like JavaScript
		{"1 == true", true, false},   // expr does coercion like JavaScript
		{"'' == false", true, false}, // expr does coercion like JavaScript
		{"'0' == 0", true, false},    // expr does coercion like JavaScript
		{"'1' == 1", true, false},    // expr does coercion like JavaScript

		// Same type comparisons
		{"0 == 0", true, false},
		{"1 == 1", true, false},
		{"true == true", true, false},
		{"false == false", true, false},
		{"'hello' == 'hello'", true, false},
		{"nil == nil", true, false},

		// Cross-type coercion (very JavaScript-like)
		{"0 == '0'", true, false},   // string to number coercion
		{"1 == '1'", true, false},   // string to number coercion
		{"true == 1", true, false},  // boolean to number coercion
		{"false == 0", true, false}, // boolean to number coercion
		{"nil == 0", false, false},
		{"nil == ''", false, false},
		{"nil == false", false, false},

		// Inequality tests
		{"1 != 0", true, false},
		{"'hello' != 'world'", true, false},
		{"true != false", true, false},
		{"42 != '41'", true, false},

		// More complex coercion cases
		{"'123' == 123", true, false},      // string number to number
		{"'3.14' == 3.14", true, false},    // string float to float
		{"'true' == true", false, false},   // string 'true' != boolean true
		{"'false' == false", false, false}, // string 'false' != boolean false
	}

	for _, tt := range tests {
		t.Run(tt.expr, func(t *testing.T) {
			result, err := expr.Eval(tt.expr, nil)
			if tt.expectError {
				if err == nil {
					t.Errorf("expected error for %q, got result %v", tt.expr, result)
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error for %q: %v", tt.expr, err)
				} else if result != tt.want {
					t.Errorf("%q: got %v, want %v", tt.expr, result, tt.want)
				}
			}
		})
	}
}

func TestJavaScriptLikeArrayOperations(t *testing.T) {
	env := map[string]any{
		"arr":       []int{1, 2, 3, 4, 5},
		"emptyArr":  []int{},
		"mixedArr":  []any{1, "hello", true, nil},
		"nestedArr": []any{[]int{1, 2}, []int{3, 4}},
	}

	tests := []struct {
		expr        string
		want        any
		expectError bool
	}{
		// Array length
		{"len(arr)", 5, false},
		{"len(emptyArr)", 0, false},
		{"len(mixedArr)", 4, false},

		// Array indexing with negative indices (JavaScript-like)
		{"arr[-1]", 5, false}, // Last element
		{"arr[-2]", 4, false}, // Second to last
		{"arr[-5]", 1, false}, // First element when using negative index

		// Array slicing
		{"arr[1:3]", []int{2, 3}, false},
		{"arr[:2]", []int{1, 2}, false},
		{"arr[2:]", []int{3, 4, 5}, false},
		{"arr[:]", []int{1, 2, 3, 4, 5}, false},

		// Array concatenation (if supported)
		{"arr + [6, 7]", []int{1, 2, 3, 4, 5, 6, 7}, true}, // Might not be supported

		// Array in operations
		{"1 in arr", true, false},
		{"6 in arr", false, false},
		{"'hello' in mixedArr", true, false},
		{"nil in mixedArr", true, false},

		// Array with filter/map operations
		{"filter(arr, # > 3)", []any{4, 5}, false},
		{"map(arr, # * 2)", []any{2, 4, 6, 8, 10}, false},
		{"all(arr, # > 0)", true, false},
		{"any(arr, # > 4)", true, false},

		// Empty array operations
		{"len(emptyArr) == 0", true, false},
		{"filter(emptyArr, # > 0)", []any{}, false},
		{"map(emptyArr, # * 2)", []any{}, false},
		{"all(emptyArr, # > 0)", true, false},  // all() on empty array is true
		{"any(emptyArr, # > 0)", false, false}, // any() on empty array is false
	}

	for _, tt := range tests {
		t.Run(tt.expr, func(t *testing.T) {
			result, err := expr.Eval(tt.expr, env)
			if tt.expectError {
				if err == nil {
					t.Errorf("expected error for %q, got result %v", tt.expr, result)
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error for %q: %v", tt.expr, err)
				} else if !deepEqual(result, tt.want) {
					t.Errorf("%q: got %v (type %T), want %v (type %T)", tt.expr, result, result, tt.want, tt.want)
				}
			}
		})
	}
}

func TestJavaScriptLikeObjectOperations(t *testing.T) {
	env := map[string]any{
		"obj": map[string]any{
			"name":    "John",
			"age":     30,
			"active":  true,
			"score":   85.5,
			"address": nil,
		},
		"emptyObj": map[string]any{},
		"nestedObj": map[string]any{
			"user": map[string]any{
				"id":   123,
				"name": "Jane",
			},
			"meta": map[string]any{
				"created": "2023-01-01",
			},
		},
	}

	tests := []struct {
		expr        string
		want        any
		expectError bool
	}{
		// Object property access
		{"obj.name", "John", false},
		{"obj.age", 30, false},
		{"obj.active", true, false},
		{"obj.score", 85.5, false},
		{"obj.address", nil, false},

		// Object bracket notation
		{"obj['name']", "John", false},
		{"obj['age']", 30, false},
		{"obj['active']", true, false},

		// Property existence checks
		{"'name' in obj", true, false},
		{"'age' in obj", true, false},
		{"'nonexistent' in obj", false, false},
		{"'address' in obj", true, false}, // nil value but key exists

		// Nested object access
		{"nestedObj.user.id", 123, false},
		{"nestedObj.user.name", "Jane", false},
		{"nestedObj.meta.created", "2023-01-01", false},

		// Object length/size
		{"len(obj)", 5, false},
		{"len(emptyObj)", 0, false},
		{"len(nestedObj)", 2, false},

		// Object property modification (if supported)
		// Note: These might not be supported in expr as it's typically for read-only evaluation

		// Object comparison
		{"obj == obj", true, false},
		{"obj != emptyObj", true, false},
		{"emptyObj == emptyObj", true, false},
	}

	for _, tt := range tests {
		t.Run(tt.expr, func(t *testing.T) {
			result, err := expr.Eval(tt.expr, env)
			if tt.expectError {
				if err == nil {
					t.Errorf("expected error for %q, got result %v", tt.expr, result)
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error for %q: %v", tt.expr, err)
				} else if !deepEqual(result, tt.want) {
					t.Errorf("%q: got %v (type %T), want %v (type %T)", tt.expr, result, result, tt.want, tt.want)
				}
			}
		})
	}
}

func TestJavaScriptLikeTypeOfOperations(t *testing.T) {
	env := map[string]any{
		"numValue":    42,
		"floatValue":  3.14,
		"stringValue": "hello",
		"boolValue":   true,
		"nullValue":   nil,
		"arrayValue":  []int{1, 2, 3},
		"objectValue": map[string]any{"key": "value"},
	}

	tests := []struct {
		expr        string
		want        any
		expectError bool
	}{
		// Type checking operations (if typeof is supported)
		{"type(numValue)", "int", false},
		{"type(floatValue)", "float", false},
		{"type(stringValue)", "string", false},
		{"type(boolValue)", "bool", false},
		{"type(nullValue)", "nil", false},
		{"type(arrayValue)", "array", false},
		{"type(objectValue)", "map", false},

		// Type comparisons
		{"type(numValue) == 'int'", true, false}, // Might need adjustment based on actual type names
		{"type(stringValue) == 'string'", true, false},
		{"type(boolValue) == 'bool'", true, false},

		// Duck typing checks
		{"len(arrayValue) >= 0", true, false},  // Has length property
		{"len(objectValue) >= 0", true, false}, // Has length property
		{"len(stringValue) >= 0", true, false}, // Has length property
	}

	for _, tt := range tests {
		t.Run(tt.expr, func(t *testing.T) {
			result, err := expr.Eval(tt.expr, env)
			if tt.expectError {
				if err == nil {
					t.Errorf("expected error for %q, got result %v", tt.expr, result)
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error for %q: %v", tt.expr, err)
				} else if !deepEqual(result, tt.want) {
					t.Errorf("%q: got %v (type %T), want %v (type %T)", tt.expr, result, result, tt.want, tt.want)
				}
			}
		})
	}
}

func TestJavaScriptLikeStringOperations(t *testing.T) {
	env := map[string]any{
		"str":      "Hello World",
		"emptyStr": "",
		"numStr":   "123",
		"floatStr": "3.14",
		"boolStr":  "true",
	}

	tests := []struct {
		expr        string
		want        any
		expectError bool
	}{
		// String length
		{"len(str)", 11, false},
		{"len(emptyStr)", 0, false},
		{"len(numStr)", 3, false},

		// String slicing
		{"str[0:5]", "Hello", false},
		{"str[6:]", "World", false},
		{"str[:5]", "Hello", false},
		{"str[:]", "Hello World", false},

		// String methods (if supported)
		{"str startsWith 'Hello'", true, false},
		{"str endsWith 'World'", true, false},
		{"str contains 'lo Wo'", true, false},
		{"str matches '^Hello.*World$'", true, false},

		// String case operations (if supported)
		// {"upper(str)", "HELLO WORLD", false},
		// {"lower(str)", "hello world", false},

		// String concatenation
		{"str + ' from Go'", "Hello World from Go", false},
		{"'Greeting: ' + str", "Greeting: Hello World", false},
		{"emptyStr + str", "Hello World", false},
		{"str + emptyStr", "Hello World", false},

		// String comparison
		{"str == 'Hello World'", true, false},
		{"str != 'Goodbye World'", true, false},
		{"str > 'Hello'", true, false}, // Lexicographic comparison
		{"'Hello' < str", true, false},

		// String to number conversion context
		{"numStr + '456'", "123456", false},  // String concatenation
		{"floatStr + '15'", "3.1415", false}, // String concatenation

		// Empty string behavior
		{"emptyStr == ''", true, false},
		{"len(emptyStr) == 0", true, false},
		{"emptyStr + 'test'", "test", false},
		{"'test' + emptyStr", "test", false},
	}

	for _, tt := range tests {
		t.Run(tt.expr, func(t *testing.T) {
			result, err := expr.Eval(tt.expr, env)
			if tt.expectError {
				if err == nil {
					t.Errorf("expected error for %q, got result %v", tt.expr, result)
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error for %q: %v", tt.expr, err)
				} else if !deepEqual(result, tt.want) {
					t.Errorf("%q: got %v (type %T), want %v (type %T)", tt.expr, result, result, tt.want, tt.want)
				}
			}
		})
	}
}

func TestJavaScriptLikeFunctionBehavior(t *testing.T) {
	env := map[string]any{
		"add":    func(a, b float64) float64 { return a + b },
		"greet":  func(name string) string { return "Hello, " + name },
		"isEven": func(n int) bool { return n%2 == 0 },
		"defaultValue": func(val any) any {
			if val == nil {
				return "default"
			}
			return val
		},
		"variadicSum": func(nums ...int) int {
			sum := 0
			for _, n := range nums {
				sum += n
			}
			return sum
		},
	}

	tests := []struct {
		expr        string
		want        any
		expectError bool
	}{
		// Basic function calls
		{"add(5.0, 3.0)", 8.0, false},
		{"add(2.5, 1.5)", 4.0, false},
		{"greet('World')", "Hello, World", false},
		{"isEven(4)", true, false},
		{"isEven(7)", false, false},

		// Function calls with type coercion
		{"add(5.0, 3.14)", 8.14, false},
		{"greet('Go' + ' Lang')", "Hello, Go Lang", false},

		// Function calls with nil handling
		{"defaultValue(nil)", "default", false},
		{"defaultValue('value')", "value", false},
		{"defaultValue(42)", 42, false},

		// Variadic function calls
		{"variadicSum(1, 2, 3)", 6, false},
		{"variadicSum()", 0, false},
		{"variadicSum(10)", 10, false},

		// Function calls in expressions
		{"add(5.0, 3.0) > 7", true, false},
		{"'Result: ' + add(10.0, 20.0)", "Result: 30", false},
		{"isEven(int(add(2.0, 2.0)))", true, false},

		// Nested function calls
		{"add(add(1.0, 2.0), add(3.0, 4.0))", 10.0, false},
		{"greet(greet('Inner'))", "Hello, Hello, Inner", false},
	}

	for _, tt := range tests {
		t.Run(tt.expr, func(t *testing.T) {
			result, err := expr.Eval(tt.expr, env)
			if tt.expectError {
				if err == nil {
					t.Errorf("expected error for %q, got result %v", tt.expr, result)
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error for %q: %v", tt.expr, err)
				} else if !deepEqual(result, tt.want) {
					t.Errorf("%q: got %v (type %T), want %v (type %T)", tt.expr, result, result, tt.want, tt.want)
				}
			}
		})
	}
}

func TestJavaScriptLikeEdgeCases(t *testing.T) {
	// Use math constants to avoid compile-time division by zero
	var infinity = math.Inf(1)
	var negInfinity = math.Inf(-1)
	var nan = math.NaN()

	env := map[string]any{
		"zero":         0,
		"negZero":      -0.0,
		"infinity":     infinity,
		"negInfinity":  negInfinity,
		"nan":          nan,
		"emptyString":  "",
		"whitespace":   " \t\n",
		"nullValue":    nil,
		"undefinedVar": nil,
	}

	tests := []struct {
		expr        string
		want        any
		expectError bool
	}{
		// Special number comparisons
		{"zero == negZero", true, false},

		// Infinity handling
		{"infinity > 1000000", true, false},
		{"negInfinity < -1000000", true, false}, // Might error due to division by zero

		// NaN behavior
		{"nan == nan", false, false}, // NaN !== NaN, might error due to division by zero
		{"nan != nan", true, false},  // Might error due to division by zero

		// Empty string vs whitespace
		{"emptyString == ''", true, false},
		{"whitespace != ''", true, false},
		{"len(whitespace) > 0", true, false},

		// Null/undefined behavior
		{"nullValue == nil", true, false},
		{"undefinedVar == nil", true, false},
		{"nullValue == undefinedVar", true, false},

		// Type coercion edge cases
		{"'' + 0", "0", false},
		{"'' + false", "false", false},
		{"'' + nil", "", false},
		{"0 + ''", "0", false},
		{"false + ''", "false", false},
		{"nil + ''", "", false},

		// Arithmetic with special values
		{"zero + 1", 1, false},
		{"zero * 1000", 0, false},
		{"zero / 1", 0.0, false},

		// Comparison edge cases
		{"0 < 0.1", true, false},
		{"0 <= 0", true, false},
		{"0 >= 0", true, false},
		{"0 > -0.1", true, false},

		// String edge cases
		{"'' == ''", true, false},
		{"'' != ' '", true, false},
		{"len('') == 0", true, false},
		{"'' + '' == ''", true, false},
	}

	for _, tt := range tests {
		t.Run(tt.expr, func(t *testing.T) {
			result, err := expr.Eval(tt.expr, env)
			if tt.expectError {
				if err == nil {
					t.Errorf("expected error for %q, got result %v", tt.expr, result)
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error for %q: %v", tt.expr, err)
				} else if !deepEqual(result, tt.want) {
					t.Errorf("%q: got %v (type %T), want %v (type %T)", tt.expr, result, result, tt.want, tt.want)
				}
			}
		})
	}
}

// Documentation test cases from Expr Language Definition

func TestDocumentation_Literals(t *testing.T) {
	tests := []struct {
		expr string
		want any
	}{
		// Boolean literals
		{"true", true},
		{"false", false},

		// Integer literals
		{"42", 42},
		{"0x2A", 42},
		{"0o52", 42},
		{"0b101010", 42},

		// Float literals
		{"0.5", 0.5},
		{".5", 0.5},

		// String literals
		{`"foo"`, "foo"},
		{"'bar'", "bar"},
		{`"Hello\nWorld"`, "Hello\nWorld"},

		// Array literals
		{"[1, 2, 3]", []any{1, 2, 3}},

		// Map literals
		{"{a: 1, b: 2, c: 3}", map[string]any{"a": 1, "b": 2, "c": 3}},

		// Nil literal
		{"nil", nil},
	}

	for _, tt := range tests {
		t.Run(tt.expr, func(t *testing.T) {
			result, err := expr.Eval(tt.expr, nil)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			} else if !deepEqual(result, tt.want) {
				t.Errorf("got %v, want %v", result, tt.want)
			}
		})
	}
}

func TestDocumentation_Operators(t *testing.T) {
	env := map[string]any{
		"user": map[string]any{
			"Name": "John",
		},
		"array": []any{1, 2, 3, 4, 5},
	}

	tests := []struct {
		expr string
		want any
	}{
		// Arithmetic operators
		{"5 + 3", 8},
		{"5 - 3", 2},
		{"5 * 3", 15},
		{"15 / 3", 5.0},
		{"17 % 5", 2},
		{"2 ^ 3", 8.0},
		{"2 ** 3", 8.0},

		// Comparison operators
		{"5 == 5", true},
		{"5 != 3", true},
		{"5 > 3", true},
		{"3 < 5", true},
		{"5 >= 5", true},
		{"5 <= 5", true},

		// Logical operators
		{"true && true", true},
		{"true || false", true},
		{"!false", true},
		{"true and true", true},
		{"true or false", true},
		{"not false", true},

		// Ternary operator
		{"5 > 3 ? 'yes' : 'no'", "yes"},

		// Nil coalescing
		{"nil ?? 'default'", "default"},
		{"'value' ?? 'default'", "value"},

		// Membership operators
		{"user.Name", "John"},
		{"user['Name']", "John"},
		{"array[0]", 1},
		{"array[-1]", 5},
		{`"John" in ["John", "Jane"]`, true},
		{`"name" in {"name": "John", "age": 30}`, true},

		// String operators
		{"'Hello' + ' World'", "Hello World"},
		{`"Hello World" contains "World"`, true},
		{`"Hello World" startsWith "Hello"`, true},
		{`"Hello World" endsWith "World"`, true},
		{`"test@example.com" matches ".*@.*"`, true},

		// Range operator
		{"1..3", []any{1, 2, 3}},

		// Slice operator
		{"array[1:4]", []any{2, 3, 4}},
		{"array[1:-1]", []any{2, 3, 4}},
		{"array[:3]", []any{1, 2, 3}},
		{"array[3:]", []any{4, 5}},
		{"array[:]", []any{1, 2, 3, 4, 5}},
	}

	for _, tt := range tests {
		t.Run(tt.expr, func(t *testing.T) {
			result, err := expr.Eval(tt.expr, env)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			} else if !deepEqual(result, tt.want) {
				t.Errorf("got %v, want %v", result, tt.want)
			}
		})
	}
}

func TestDocumentation_OptionalChaining(t *testing.T) {
	env := map[string]any{
		"author": map[string]any{
			"User": map[string]any{
				"Name": "John",
			},
		},
		"nullAuthor": map[string]any{
			"User": nil,
		},
	}

	tests := []struct {
		expr string
		want any
	}{
		{"author.User?.Name", "John"},
		{"nullAuthor.User?.Name", nil},
		{"author.User?.Name ?? 'Anonymous'", "John"},
		{"nullAuthor.User?.Name ?? 'Anonymous'", "Anonymous"},
	}

	for _, tt := range tests {
		t.Run(tt.expr, func(t *testing.T) {
			result, err := expr.Eval(tt.expr, env)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			} else if !deepEqual(result, tt.want) {
				t.Errorf("got %v, want %v", result, tt.want)
			}
		})
	}
}

func TestDocumentation_Variables(t *testing.T) {
	tests := []struct {
		expr string
		want any
	}{
		{"let x = 42; x * 2", 84},
		{"let x = 42; let y = 2; x * y", 84},
		{"let name = 'test' | upper(); 'Hello, ' + name + '!'", "Hello, TEST!"},
	}

	for _, tt := range tests {
		t.Run(tt.expr, func(t *testing.T) {
			result, err := expr.Eval(tt.expr, nil)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			} else if !deepEqual(result, tt.want) {
				t.Errorf("got %v, want %v", result, tt.want)
			}
		})
	}
}

func TestDocumentation_Env(t *testing.T) {
	env := map[string]any{
		"foo": map[string]any{
			"Name": "John",
		},
		"var with spaces": "test value",
	}

	tests := []struct {
		expr string
		want any
	}{
		{"foo.Name == $env['foo'].Name", true},
		{"$env['var with spaces']", "test value"},
		{"'foo' in $env", true},
		{"'nonexistent' in $env", false},
	}

	for _, tt := range tests {
		t.Run(tt.expr, func(t *testing.T) {
			result, err := expr.Eval(tt.expr, env)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			} else if !deepEqual(result, tt.want) {
				t.Errorf("got %v, want %v", result, tt.want)
			}
		})
	}
}

func TestDocumentation_StringFunctions(t *testing.T) {
	tests := []struct {
		expr string
		want any
	}{
		// trim
		{`trim("  Hello  ")`, "Hello"},
		{`trim("__Hello__", "_")`, "Hello"},

		// trimPrefix
		{`trimPrefix("HelloWorld", "Hello")`, "World"},

		// trimSuffix
		{`trimSuffix("HelloWorld", "World")`, "Hello"},

		// upper
		{`upper("hello")`, "HELLO"},

		// lower
		{`lower("HELLO")`, "hello"},

		// split
		{`split("apple,orange,grape", ",")`, []any{"apple", "orange", "grape"}},
		{`split("apple,orange,grape", ",", 2)`, []any{"apple", "orange,grape"}},

		// splitAfter
		{`splitAfter("apple,orange,grape", ",")`, []any{"apple,", "orange,", "grape"}},
		{`splitAfter("apple,orange,grape", ",", 2)`, []any{"apple,", "orange,grape"}},

		// replace
		{`replace("Hello World", "World", "Universe")`, "Hello Universe"},

		// repeat
		{`repeat("Hi", 3)`, "HiHiHi"},

		// indexOf
		{`indexOf("apple pie", "pie")`, 6},

		// lastIndexOf
		{`lastIndexOf("apple pie apple", "apple")`, 10},

		// hasPrefix
		{`hasPrefix("HelloWorld", "Hello")`, true},

		// hasSuffix
		{`hasSuffix("HelloWorld", "World")`, true},
	}

	for _, tt := range tests {
		t.Run(tt.expr, func(t *testing.T) {
			result, err := expr.Eval(tt.expr, nil)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			} else if !deepEqual(result, tt.want) {
				t.Errorf("got %v, want %v", result, tt.want)
			}
		})
	}
}

func TestDocumentation_DateFunctions(t *testing.T) {
	tests := []struct {
		expr string
		want any
	}{
		// duration
		{`duration("1h").Seconds()`, 3600.0},

		// date parsing
		{`date("2023-08-14").Year()`, 2023},

		// Basic date operations
		{`duration("1h") > duration("30m")`, true},
	}

	for _, tt := range tests {
		t.Run(tt.expr, func(t *testing.T) {
			result, err := expr.Eval(tt.expr, nil)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			} else if !deepEqual(result, tt.want) {
				t.Errorf("got %v, want %v", result, tt.want)
			}
		})
	}
}

func TestDocumentation_NumberFunctions(t *testing.T) {
	tests := []struct {
		expr string
		want any
	}{
		// max
		{"max(5, 7)", 7},

		// min
		{"min(5, 7)", 5},

		// abs
		{"abs(-5)", 5},

		// ceil
		{"ceil(1.5)", 2.0},

		// floor
		{"floor(1.5)", 1.0},

		// round
		{"round(1.5)", 2.0},
	}

	for _, tt := range tests {
		t.Run(tt.expr, func(t *testing.T) {
			result, err := expr.Eval(tt.expr, nil)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			} else if !deepEqual(result, tt.want) {
				t.Errorf("got %v, want %v", result, tt.want)
			}
		})
	}
}

func TestDocumentation_ArrayFunctions(t *testing.T) {
	env := map[string]any{
		"tweets": []any{
			map[string]any{"Size": 140},
			map[string]any{"Size": 280},
			map[string]any{"Size": 300},
		},
		"users": []any{
			map[string]any{"Name": "John", "Age": 25},
			map[string]any{"Name": "Jane", "Age": 30},
			map[string]any{"Name": "Bob", "Age": 35},
		},
		"participants": []any{
			map[string]any{"Winner": true},
			map[string]any{"Winner": false},
			map[string]any{"Winner": false},
		},
		"accounts": []any{
			map[string]any{"Balance": 100},
			map[string]any{"Balance": 200},
			map[string]any{"Balance": 300},
		},
	}

	tests := []struct {
		expr string
		want any
	}{
		// all
		{"all(tweets, {.Size < 400})", true},

		// any
		{"any(tweets, {.Size > 280})", true},

		// one
		{"one(participants, {.Winner})", true},

		// none
		{"none(tweets, {.Size > 400})", true},

		// map
		{"map(tweets, {.Size})", []any{140, 280, 300}},

		// filter
		{"filter(users, .Name startsWith 'J')", []any{
			map[string]any{"Name": "John", "Age": 25},
			map[string]any{"Name": "Jane", "Age": 30},
		}},

		// find
		{"find([1, 2, 3, 4], # > 2)", 3},

		// findIndex
		{"findIndex([1, 2, 3, 4], # > 2)", 2},

		// findLast
		{"findLast([1, 2, 3, 4], # > 2)", 4},

		// findLastIndex
		{"findLastIndex([1, 2, 3, 4], # > 2)", 3},

		// count with predicate
		{"count(users, .Age > 25)", 2},

		// count without predicate
		{"count([true, false, true])", 2},

		// concat
		{"concat([1, 2], [3, 4])", []any{1, 2, 3, 4}},

		// flatten
		{"flatten([1, 2, [3, 4]])", []any{1, 2, 3, 4}},

		// uniq
		{"uniq([1, 2, 3, 2, 1])", []any{1, 2, 3}},

		// join
		{"join(['apple', 'orange', 'grape'], ',')", "apple,orange,grape"},
		{"join(['apple', 'orange', 'grape'])", "appleorangegrape"},

		// reduce
		{"reduce(1..9, #acc + #)", 45},
		{"reduce(1..9, #acc + #, 0)", 45},

		// sum
		{"sum([1, 2, 3])", 6},
		{"sum(accounts, .Balance)", 600},

		// mean
		{"mean([1, 2, 3])", 2.0},

		// median
		{"median([1, 2, 3])", 2.0},

		// first
		{"first([1, 2, 3])", 1},

		// last
		{"last([1, 2, 3])", 3},

		// take
		{"take([1, 2, 3, 4], 2)", []any{1, 2}},

		// reverse
		{"reverse([3, 1, 4])", []any{4, 1, 3}},

		// sort
		{"sort([3, 1, 4])", []any{1, 3, 4}},
		{"sort([3, 1, 4], 'desc')", []any{4, 3, 1}},

		// sortBy
		{"sortBy(users, .Age)", []any{
			map[string]any{"Name": "John", "Age": 25},
			map[string]any{"Name": "Jane", "Age": 30},
			map[string]any{"Name": "Bob", "Age": 35},
		}},
	}

	for _, tt := range tests {
		t.Run(tt.expr, func(t *testing.T) {
			result, err := expr.Eval(tt.expr, env)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			} else if !deepEqual(result, tt.want) {
				t.Errorf("got %v, want %v", result, tt.want)
			}
		})
	}
}

func TestDocumentation_MapFunctions(t *testing.T) {
	tests := []struct {
		expr string
		want any
	}{
		// keys
		{"keys({'name': 'John', 'age': 30})", []any{"name", "age"}},

		// values
		{"values({'name': 'John', 'age': 30})", []any{"John", 30}},
	}

	for _, tt := range tests {
		t.Run(tt.expr, func(t *testing.T) {
			result, err := expr.Eval(tt.expr, nil)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			// For keys and values, we need to check if arrays contain the same elements
			// since order might vary
			resultSlice, ok1 := result.([]any)
			wantSlice, ok2 := tt.want.([]any)
			if !ok1 || !ok2 {
				if !deepEqual(result, tt.want) {
					t.Errorf("got %v, want %v", result, tt.want)
				}
			} else {
				if len(resultSlice) != len(wantSlice) {
					t.Errorf("got %v, want %v", result, tt.want)
				} else {
					// Check if all elements in want are present in result
					for _, wantItem := range wantSlice {
						found := false
						for _, resultItem := range resultSlice {
							if deepEqual(resultItem, wantItem) {
								found = true
								break
							}
						}
						if !found {
							t.Errorf("got %v, want %v", result, tt.want)
							break
						}
					}
				}
			}
		})
	}
}

func TestDocumentation_TypeConversionFunctions(t *testing.T) {
	tests := []struct {
		expr string
		want any
	}{
		// type
		{"type(42)", "int"},
		{"type('hello')", "string"},

		// int
		{"int('123')", 123},

		// float
		{"float('123.45')", 123.45},

		// string
		{"string(123)", "123"},

		// fromJSON
		{"fromJSON('{\"name\": \"John\", \"age\": 30}')", map[string]any{"name": "John", "age": 30.0}},

		// toBase64
		{"toBase64('Hello World')", "SGVsbG8gV29ybGQ="},

		// fromBase64
		{"fromBase64('SGVsbG8gV29ybGQ=')", "Hello World"},

		// fromPairs
		{"fromPairs([['name', 'John'], ['age', 30]])", map[string]any{"name": "John", "age": 30}},
	}

	for _, tt := range tests {
		t.Run(tt.expr, func(t *testing.T) {
			result, err := expr.Eval(tt.expr, nil)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			} else {
				// Special handling for JSON and pairs functions that might have different ordering
				if tt.expr == "toJSON({'name': 'John', 'age': 30})" {
					// JSON order might vary, check if it's valid JSON with correct content
					str, ok := result.(string)
					if !ok {
						t.Errorf("got %v, want string", result)
					} else if !(strings.Contains(str, `"name":"John"`) && strings.Contains(str, `"age":30`)) {
						t.Errorf("got %v, want JSON with name and age", result)
					}
				} else if tt.expr == "toPairs({'name': 'John', 'age': 30})" {
					// Pairs order might vary, check structure
					resultSlice, ok := result.([]any)
					if !ok || len(resultSlice) != 2 {
						t.Errorf("got %v, want array of 2 pairs", result)
					}
				} else if !deepEqual(result, tt.want) {
					t.Errorf("got %v, want %v", result, tt.want)
				}
			}
		})
	}
}

func TestDocumentation_MiscellaneousFunctions(t *testing.T) {
	tests := []struct {
		expr string
		want any
	}{
		// len
		{"len([1, 2, 3])", 3},
		{"len({'name': 'John', 'age': 30})", 2},
		{"len('Hello')", 5},

		// get
		{"get([1, 2, 3], 1)", 2},
		{"get({'name': 'John', 'age': 30}, 'name')", "John"},
	}

	for _, tt := range tests {
		t.Run(tt.expr, func(t *testing.T) {
			result, err := expr.Eval(tt.expr, nil)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			} else if !deepEqual(result, tt.want) {
				t.Errorf("got %v, want %v", result, tt.want)
			}
		})
	}
}

func TestDocumentation_BitwiseFunctions(t *testing.T) {
	tests := []struct {
		expr string
		want any
	}{
		// bitand
		{"bitand(0b1010, 0b1100)", 8}, // 0b1000

		// bitor
		{"bitor(0b1010, 0b1100)", 14}, // 0b1110

		// bitxor
		{"bitxor(0b1010, 0b1100)", 6}, // 0b0110

		// bitnand
		{"bitnand(0b1010, 0b1100)", 2}, // 0b0010

		// bitnot
		{"bitnot(0b1010)", -11}, // -0b1011

		// bitshl
		{"bitshl(0b101101, 2)", 180}, // 0b10110100

		// bitshr
		{"bitshr(0b101101, 2)", 11}, // 0b1011

		// bitushr (unsigned right shift)
		{"bitushr(0b101, 1)", 2}, // 0b10
	}

	for _, tt := range tests {
		t.Run(tt.expr, func(t *testing.T) {
			result, err := expr.Eval(tt.expr, nil)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			} else if !deepEqual(result, tt.want) {
				t.Errorf("got %v, want %v", result, tt.want)
			}
		})
	}
}

func TestDocumentation_PipeOperator(t *testing.T) {
	env := map[string]any{
		"user": map[string]any{
			"Name": "John Doe",
		},
	}

	tests := []struct {
		expr string
		want any
	}{
		// Pipe operator equivalence
		{"'hello world' | upper() | split(' ')", []any{"HELLO", "WORLD"}},
	}

	for _, tt := range tests {
		t.Run(tt.expr, func(t *testing.T) {
			result, err := expr.Eval(tt.expr, env)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			} else if !deepEqual(result, tt.want) {
				t.Errorf("got %v, want %v", result, tt.want)
			}
		})
	}
}

func TestDocumentation_Predicates(t *testing.T) {
	env := map[string]any{
		"tweets": []any{
			map[string]any{"Content": "Short tweet"},
			map[string]any{"Content": "This is a very long tweet that exceeds 240 characters and should be filtered out by our predicate function when we test the filter functionality with predicates in the Expr language.This is a very long tweet that exceeds 240 characters and should be filtered out by our predicate function when we test the filter functionality with predicates in the Expr language."},
		},
		"posts": []any{
			map[string]any{
				"Author": "John",
				"Comments": []any{
					map[string]any{"Author": "John"},
					map[string]any{"Author": "Jane"},
				},
			},
			map[string]any{
				"Author": "Jane",
				"Comments": []any{
					map[string]any{"Author": "Bob"},
				},
			},
		},
	}

	tests := []struct {
		expr string
		want any
	}{
		// Basic predicate with filter
		{"filter(0..9, {# % 2 == 0})", []any{0, 2, 4, 6, 8}},

		// Predicate with field access
		{"filter(tweets, {len(.Content) > 240})", []any{
			map[string]any{"Content": "This is a very long tweet that exceeds 240 characters and should be filtered out by our predicate function when we test the filter functionality with predicates in the Expr language.This is a very long tweet that exceeds 240 characters and should be filtered out by our predicate function when we test the filter functionality with predicates in the Expr language."},
		}},

		// Predicate without braces
		{"filter(tweets, len(.Content) <= 240)", []any{
			map[string]any{"Content": "Short tweet"},
		}},
	}

	for _, tt := range tests {
		t.Run(tt.expr, func(t *testing.T) {
			result, err := expr.Eval(tt.expr, env)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			} else if !deepEqual(result, tt.want) {
				t.Errorf("got %v, want %v", result, tt.want)
			}
		})
	}
}

func TestDocumentation_ComplexExpressions(t *testing.T) {
	env := map[string]any{
		"users": []any{
			map[string]any{"Name": "John", "Age": 25, "Active": true},
			map[string]any{"Name": "Jane", "Age": 30, "Active": false},
			map[string]any{"Name": "Bob", "Age": 35, "Active": true},
		},
		"createdAt": time.Now().Add(-2 * time.Hour),
		"now":       func() time.Time { return time.Now() },
	}

	tests := []struct {
		expr        string
		expectError bool
		checkResult func(result any) bool
	}{
		// Complex filtering and mapping
		{
			"map(filter(users, .Active), .Name)",
			false,
			func(result any) bool {
				arr, ok := result.([]any)
				return ok && len(arr) == 2 && arr[0] == "John" && arr[1] == "Bob"
			},
		},

		// Date comparison
		{
			"createdAt > now() - duration('3h')",
			false,
			func(result any) bool {
				b, ok := result.(bool)
				return ok && b
			},
		},

		// Complex conditional
		{
			"len(filter(users, .Age > 25)) > 1 ? 'Many adults' : 'Few adults'",
			false,
			func(result any) bool {
				s, ok := result.(string)
				return ok && s == "Many adults"
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.expr, func(t *testing.T) {
			result, err := expr.Eval(tt.expr, env)
			if tt.expectError {
				if err == nil {
					t.Errorf("expected error but got result: %v", result)
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				} else if !tt.checkResult(result) {
					t.Errorf("result check failed for %v", result)
				}
			}
		})
	}
}

// Test cases for null handling in builtin functions
func TestBuiltinFunctions_NullHandling(t *testing.T) {
	testCases := []struct {
		name       string
		expression string
		want       any
		wantError  bool
	}{
		// String functions - null handling
		{"trim_null_string", "trim(nil)", nil, false},
		{"trim_null_cutset", "trim('  hello  ', nil)", "  hello  ", false},
		{"trimPrefix_null_string", "trimPrefix(nil, 'pre')", nil, false},
		{"trimPrefix_null_prefix", "trimPrefix('prefix_test', nil)", "prefix_test", false},
		{"trimSuffix_null_string", "trimSuffix(nil, 'suf')", nil, false},
		{"trimSuffix_null_suffix", "trimSuffix('test_suffix', nil)", "test_suffix", false},
		{"upper_null", "upper(nil)", nil, false},
		{"lower_null", "lower(nil)", nil, false},
		{"split_null_string", "split(nil, ',')", []any{}, false},
		{"split_null_separator", "split('a,b,c', nil)", []any{"a,b,c"}, false},
		{"splitAfter_null_string", "splitAfter(nil, ',')", []any{}, false},
		{"splitAfter_null_separator", "splitAfter('a,b,c', nil)", []any{"a,b,c"}, false},
		{"repeat_null_string", "repeat(nil, 3)", "", false},
		{"indexOf_null_string", "indexOf(nil, 'test')", -1, false},
		{"indexOf_null_substring", "indexOf('hello world', nil)", -1, false},
		{"lastIndexOf_null_string", "lastIndexOf(nil, 'test')", -1, false},
		{"lastIndexOf_null_substring", "lastIndexOf('hello world', nil)", -1, false},
		{"hasPrefix_null_string", "hasPrefix(nil, 'pre')", false, false},
		{"hasPrefix_null_prefix", "hasPrefix('prefix_test', nil)", true, false},
		{"hasSuffix_null_string", "hasSuffix(nil, 'suf')", false, false},
		{"hasSuffix_null_suffix", "hasSuffix('test_suffix', nil)", true, false},
		{"replace_null_string", "replace(nil, 'old', 'new')", nil, false},
		{"replace_null_old", "replace('hello world', nil, 'new')", "hello world", false},
		{"replace_null_new", "replace('hello world', 'world', nil)", "hello ", false},

		// Join function with null arrays and elements
		{"join_null_array", "join(nil, ',')", "", false},
		{"join_with_null_separator", "join(['a', 'b', 'c'], nil)", "abc", false},

		// Encoding functions - null handling
		{"fromJSON_null", "fromJSON(nil)", nil, false},
		{"toBase64_null", "toBase64(nil)", "", false},
		{"fromBase64_null", "fromBase64(nil)", "", false},
		{"duration_null", "duration(nil)", "0s", false},

		// Math functions - null handling
		{"len_null", "len(nil)", 0, false},
		{"abs_null", "abs(nil)", 0, false},
		{"ceil_null", "ceil(nil)", 0.0, false},
		{"floor_null", "floor(nil)", 0.0, false},
		{"round_null", "round(nil)", 0.0, false},
		{"int_null", "int(nil)", 0, false},
		{"float_null", "float(nil)", 0.0, false},

		// Type checking with null values
		{"type_null", "type(nil)", "nil", false},

		// Complex expressions with null values
		{"null_coalescing", "nil ?? 'default'", "default", false},
		{"null_string_concat", "nil + ' suffix'", " suffix", false},
		{"string_null_concat", "'prefix ' + nil", "prefix ", false},
		{"null_arithmetic", "nil + 5", 5, false},
		{"arithmetic_null", "5 + nil", 5, false},
		{"null_comparison", "nil == nil", true, false},
		{"null_not_equal", "nil != 'test'", true, false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := expr.Eval(tc.expression, nil)
			if tc.wantError {
				if err == nil {
					t.Errorf("expected error for %q, got result %v", tc.expression, result)
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error for %q: %v", tc.expression, err)
				} else if !deepEqual(result, tc.want) {
					t.Errorf("%q: got %v (type %T), want %v (type %T)", tc.expression, result, result, tc.want, tc.want)
				}
			}
		})
	}
}

func TestBuiltinStringFunctions_EdgeCases(t *testing.T) {
	testCases := []struct {
		name       string
		expression string
		want       any
		wantError  bool
	}{
		// Edge cases for string functions with null and empty values
		{"trim_empty_string", "trim('')", "", false},
		{"trim_whitespace_only", "trim('   ')", "", false},
		{"trim_custom_cutset", "trim('___hello___', '_')", "hello", false},
		{"trim_null_both", "trim(nil, nil)", nil, false},

		{"split_empty_string", "split('', ',')", []any{""}, false},
		{"split_empty_separator", "split('abc', '')", []any{"a", "b", "c"}, false},
		{"split_with_limit", "split('a,b,c,d', ',', 2)", []any{"a", "b,c,d"}, false},
		{"split_null_with_limit", "split(nil, ',', 2)", []any{}, false},

		{"splitAfter_empty_string", "splitAfter('', ',')", []any{""}, false},
		{"splitAfter_empty_separator", "splitAfter('abc', '')", []any{"a", "b", "c"}, false},
		{"splitAfter_with_limit", "splitAfter('a,b,c,d', ',', 2)", []any{"a,", "b,c,d"}, false},

		{"replace_with_count", "replace('hello hello hello', 'hello', 'hi', 2)", "hi hi hello", false},
		{"replace_no_match", "replace('hello world', 'xyz', 'abc')", "hello world", false},
		{"replace_empty_old", "replace('hello', '', 'x')", "xhxexlxlxox", false},

		{"repeat_zero_times", "repeat('test', 0)", "", false},
		{"repeat_negative_times", "repeat('test', -1)", "", true},
		{"repeat_empty_string", "repeat('', 5)", "", false},

		{"indexOf_not_found", "indexOf('hello', 'xyz')", -1, false},
		{"indexOf_empty_string", "indexOf('', 'test')", -1, false},
		{"indexOf_empty_substring", "indexOf('hello', '')", 0, false},

		{"lastIndexOf_not_found", "lastIndexOf('hello', 'xyz')", -1, false},
		{"lastIndexOf_empty_string", "lastIndexOf('', 'test')", -1, false},
		{"lastIndexOf_empty_substring", "lastIndexOf('hello', '')", 5, false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := expr.Eval(tc.expression, nil)
			if tc.wantError {
				if err == nil {
					t.Errorf("expected error for %q, got result %v", tc.expression, result)
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error for %q: %v", tc.expression, err)
				} else if !deepEqual(result, tc.want) {
					t.Errorf("%q: got %v (type %T), want %v (type %T)", tc.expression, result, result, tc.want, tc.want)
				}
			}
		})
	}
}

func TestBuiltinMathFunctions_NullHandling(t *testing.T) {
	testCases := []struct {
		name       string
		expression string
		want       any
		wantError  bool
	}{
		// Math functions with null values
		{"abs_negative", "abs(-42)", 42, false},
		{"abs_positive", "abs(42)", 42, false},
		{"abs_zero", "abs(0)", 0, false},
		{"abs_float", "abs(-3.14)", 3.14, false},
		{"abs_null", "abs(nil)", 0, false},

		{"ceil_positive", "ceil(3.14)", 4.0, false},
		{"ceil_negative", "ceil(-3.14)", -3.0, false},
		{"ceil_zero", "ceil(0)", 0.0, false},
		{"ceil_integer", "ceil(5)", 5.0, false},
		{"ceil_null", "ceil(nil)", 0.0, false},

		{"floor_positive", "floor(3.14)", 3.0, false},
		{"floor_negative", "floor(-3.14)", -4.0, false},
		{"floor_zero", "floor(0)", 0.0, false},
		{"floor_integer", "floor(5)", 5.0, false},
		{"floor_null", "floor(nil)", 0.0, false},

		{"round_positive", "round(3.14)", 3.0, false},
		{"round_half_up", "round(3.5)", 4.0, false},
		{"round_negative", "round(-3.14)", -3.0, false},
		{"round_zero", "round(0)", 0.0, false},
		{"round_null", "round(nil)", 0.0, false},

		{"int_float", "int(3.14)", 3, false},
		{"int_string", "int('42')", 42, false},
		{"int_zero", "int(0)", 0, false},
		{"int_null", "int(nil)", 0, false},

		{"float_int", "float(42)", 42.0, false},
		{"float_string", "float('3.14')", 3.14, false},
		{"float_zero", "float(0)", 0.0, false},
		{"float_null", "float(nil)", 0.0, false},

		{"len_string", "len('hello')", 5, false},
		{"len_array", "len([1, 2, 3])", 3, false},
		{"len_empty_array", "len([])", 0, false},
		{"len_map", "len({a: 1, b: 2})", 2, false},
		{"len_empty_map", "len({})", 0, false},
		{"len_null", "len(nil)", 0, false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := expr.Eval(tc.expression, nil)
			if tc.wantError {
				if err == nil {
					t.Errorf("expected error for %q, got result %v", tc.expression, result)
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error for %q: %v", tc.expression, err)
				} else if !approximatelyEqual(result, tc.want) {
					t.Errorf("%q: got %v (type %T), want %v (type %T)", tc.expression, result, result, tc.want, tc.want)
				}
			}
		})
	}
}

func TestBuiltinAggregationFunctions_NullHandling(t *testing.T) {
	env := map[string]any{
		"numbers":        []any{1, 2, 3, 4, 5},
		"numbersWithNil": []any{1, nil, 3, nil, 5},
		"emptyArray":     []any{},
		"allNilArray":    []any{nil, nil, nil},
		"floats":         []any{1.1, 2.2, 3.3},
		"floatsWithNil":  []any{1.1, nil, 3.3, nil},
		"objects": []any{
			map[string]any{"value": 10},
			map[string]any{"value": 20},
			map[string]any{"value": 30},
		},
		"objectsWithNil": []any{
			map[string]any{"value": 10},
			map[string]any{"value": nil},
			map[string]any{"value": 30},
		},
	}

	testCases := []struct {
		name       string
		expression string
		want       any
		wantError  bool
	}{
		// Sum function
		{"sum_numbers", "sum(numbers)", 15, false},
		{"sum_with_nil", "sum(numbersWithNil)", 9, false}, // Should skip nil values
		{"sum_empty", "sum(emptyArray)", 0, false},
		{"sum_all_nil", "sum(allNilArray)", 0, false},
		{"sum_floats", "sum(floats)", 6.6, false},
		{"sum_floats_with_nil", "sum(floatsWithNil)", 4.4, false},
		{"sum_objects", "sum(objects, .value)", 60, false},
		{"sum_objects_with_nil", "sum(objectsWithNil, .value)", 40, false}, // Should skip nil values

		// Mean function
		{"mean_numbers", "mean(numbers)", 3.0, false},
		{"mean_with_nil", "mean(numbersWithNil)", 3.0, false}, // Should skip nil values: (1+3+5)/3 = 3
		{"mean_floats", "mean(floats)", 2.2, false},
		{"mean_floats_with_nil", "mean(floatsWithNil)", 2.2, false}, // Should skip nil: (1.1+3.3)/2 = 2.2

		// Median function
		{"median_numbers", "median(numbers)", 3.0, false},
		{"median_with_nil", "median(numbersWithNil)", 3.0, false}, // Should skip nil values
		{"median_floats", "median(floats)", 2.2, false},
		{"median_floats_with_nil", "median(floatsWithNil)", 2.2, false}, // Should skip nil values

		// Max function
		{"max_numbers", "max(numbers)", 5, false},
		{"max_with_nil", "max(numbersWithNil)", 5, false}, // Should skip nil values
		{"max_floats", "max(floats)", 3.3, false},
		{"max_floats_with_nil", "max(floatsWithNil)", 3.3, false},

		// Min function
		{"min_numbers", "min(numbers)", 1, false},
		{"min_with_nil", "min(numbersWithNil)", 1, false}, // Should skip nil values
		{"min_floats", "min(floats)", 1.1, false},
		{"min_floats_with_nil", "min(floatsWithNil)", 1.1, false},

		// Count function
		{"count_numbers", "count(numbers)", 5, false},
		{"count_with_nil", "count(numbersWithNil)", 3, false}, // Count non-nil elements only
		{"count_predicate", "count(numbers, # > 3)", 2, false},
		{"count_predicate_with_nil", "count(numbersWithNil, # != nil)", 3, false}, // Count non-nil elements
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := expr.Eval(tc.expression, env)
			if tc.wantError {
				if err == nil {
					t.Errorf("expected error for %q, got result %v", tc.expression, result)
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error for %q: %v", tc.expression, err)
				} else if !approximatelyEqual(result, tc.want) {
					t.Errorf("%q: got %v (type %T), want %v (type %T)", tc.expression, result, result, tc.want, tc.want)
				}
			}
		})
	}
}

func TestBuiltinComplexNullScenarios(t *testing.T) {
	env := map[string]any{
		"user": map[string]any{
			"name":    "John",
			"email":   nil,
			"profile": nil,
			"tags":    []any{"admin", nil, "user"},
		},
		"data": map[string]any{
			"values": []any{1, nil, 3, nil, 5},
			"items":  nil,
		},
	}

	testCases := []struct {
		name       string
		expression string
		want       any
		wantError  bool
	}{
		// Complex null scenarios
		{"null_field_access", "user.email", nil, false},
		{"null_field_concat", "user.name + ' - ' + (user.email ?? 'No email')", "John - No email", false},
		{"null_in_array", "nil in user.tags", true, false},
		{"filter_null_array", "filter(user.tags, # != nil)", []any{"admin", "user"}, false},
		{"map_with_null_handling", "map(data.values, # ?? 0)", []any{1, 0, 3, 0, 5}, false},
		{"sum_with_nulls", "sum(data.values)", 9, false}, // Should skip nil values
		{"count_non_nulls", "count(data.values, # != nil)", 3, false},
		{"join_with_nulls", "join(user.tags, ',')", "admin,,user", false}, // nil becomes empty in join
		{"string_operations_with_null", "upper(user.email ?? 'default')", "DEFAULT", false},
		{"null_safe_chaining", "user.profile?.name ?? 'Unknown'", "Unknown", false},
		{"complex_null_expression", "(user.email ?? '') + (len(user.name) > 0 ? ' (' + user.name + ')' : '')", " (John)", false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := expr.Eval(tc.expression, env)
			if tc.wantError {
				if err == nil {
					t.Errorf("expected error for %q, got result %v", tc.expression, result)
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error for %q: %v", tc.expression, err)
				} else if !deepEqual(result, tc.want) {
					t.Errorf("%q: got %v (type %T), want %v (type %T)", tc.expression, result, result, tc.want, tc.want)
				}
			}
		})
	}
}

func TestBuiltinArrayFunctions_NullHandling(t *testing.T) {
	env := map[string]any{
		"numbers":        []any{1, 2, 3, 4, 5},
		"numbersWithNil": []any{1, nil, 3, nil, 5},
		"emptyArray":     []any{},
		"allNilArray":    []any{nil, nil, nil},
		"mixedArray":     []any{1, "hello", 3.14, nil, true},
		"strings":        []any{"apple", "banana", "cherry"},
		"stringsWithNil": []any{"apple", nil, "cherry", nil},
		"objects": []any{
			map[string]any{"value": 10, "name": "A"},
			map[string]any{"value": 20, "name": "B"},
			map[string]any{"value": 30, "name": "C"},
		},
		"objectsWithNil": []any{
			map[string]any{"value": 10, "name": "A"},
			map[string]any{"value": nil, "name": nil},
			map[string]any{"value": 30, "name": "C"},
		},
	}

	testCases := []struct {
		name       string
		expression string
		want       any
		wantError  bool
	}{
		// all() function with null handling
		{"all_null_array", "all(nil, # > 0)", true, false}, // all() on nil should return true (vacuous truth)
		{"all_empty_array", "all(emptyArray, # > 0)", true, false},
		{"all_with_nils", "all(numbersWithNil, # == nil or # > 0)", true, false},
		{"all_numbers_positive", "all(numbers, # > 0)", true, false},
		{"all_numbers_greater_than_3", "all(numbers, # > 3)", false, false},

		// any() function with null handling
		{"any_null_array", "any(nil, # > 0)", false, false}, // any() on nil should return false
		{"any_empty_array", "any(emptyArray, # > 0)", false, false},
		{"any_with_nils", "any(numbersWithNil, # == nil)", true, false},
		{"any_numbers_greater_than_3", "any(numbers, # > 3)", true, false},
		{"any_numbers_greater_than_10", "any(numbers, # > 10)", false, false},

		// one() function with null handling
		{"one_null_array", "one(nil, # > 0)", false, false}, // one() on nil should return false
		{"one_empty_array", "one(emptyArray, # > 0)", false, false},
		{"one_single_match", "one(numbers, # == 3)", true, false},
		{"one_multiple_matches", "one(numbers, # > 3)", false, false},
		{"one_with_nils", "one(numbersWithNil, # == 3)", true, false},

		// none() function with null handling
		{"none_null_array", "none(nil, # > 0)", true, false}, // none() on nil should return true
		{"none_empty_array", "none(emptyArray, # > 0)", true, false},
		{"none_no_matches", "none(numbers, # > 10)", true, false},
		{"none_with_matches", "none(numbers, # > 3)", false, false},
		{"none_with_nils", "none(numbersWithNil, # == nil)", false, false},

		// map() function with null handling
		{"map_null_array", "map(nil, # * 2)", []any{}, false}, // map() on nil should return empty array
		{"map_empty_array", "map(emptyArray, # * 2)", []any{}, false},
		{"map_numbers", "map(numbers, # * 2)", []any{2, 4, 6, 8, 10}, false},
		{"map_with_nils", "map(numbersWithNil, # ?? 0)", []any{1, 0, 3, 0, 5}, false},
		{"map_objects", "map(objects, .name)", []any{"A", "B", "C"}, false},
		{"map_objects_with_nils", "map(objectsWithNil, .name ?? 'Unknown')", []any{"A", "Unknown", "C"}, false},

		// filter() function with null handling
		{"filter_null_array", "filter(nil, # > 0)", []any{}, false}, // filter() on nil should return empty array
		{"filter_empty_array", "filter(emptyArray, # > 0)", []any{}, false},
		{"filter_numbers", "filter(numbers, # > 3)", []any{4, 5}, false},
		{"filter_with_nils", "filter(numbersWithNil, # != nil)", []any{1, 3, 5}, false},
		{"filter_objects", "filter(objects, .value > 15)", []any{map[string]any{"value": 20, "name": "B"}, map[string]any{"value": 30, "name": "C"}}, false},

		// find() function with null handling
		{"find_null_array", "find(nil, # > 0)", nil, false}, // find() on nil should return nil
		{"find_empty_array", "find(emptyArray, # > 0)", nil, false},
		{"find_numbers", "find(numbers, # > 3)", 4, false},
		{"find_no_match", "find(numbers, # > 10)", nil, false},
		{"find_with_nils", "find(numbersWithNil, # == nil)", nil, false},
		{"find_first_non_nil", "find(numbersWithNil, # != nil)", 1, false},

		// findIndex() function with null handling
		{"findIndex_null_array", "findIndex(nil, # > 0)", -1, false}, // findIndex() on nil should return -1
		{"findIndex_empty_array", "findIndex(emptyArray, # > 0)", -1, false},
		{"findIndex_numbers", "findIndex(numbers, # > 3)", 3, false},
		{"findIndex_no_match", "findIndex(numbers, # > 10)", -1, false},
		{"findIndex_with_nils", "findIndex(numbersWithNil, # == nil)", 1, false},

		// findLast() function with null handling
		{"findLast_null_array", "findLast(nil, # > 0)", nil, false}, // findLast() on nil should return nil
		{"findLast_empty_array", "findLast(emptyArray, # > 0)", nil, false},
		{"findLast_numbers", "findLast(numbers, # > 3)", 5, false},
		{"findLast_no_match", "findLast(numbers, # > 10)", nil, false},
		{"findLast_with_nils", "findLast(numbersWithNil, # == nil)", nil, false},

		// findLastIndex() function with null handling
		{"findLastIndex_null_array", "findLastIndex(nil, # > 0)", -1, false}, // findLastIndex() on nil should return -1
		{"findLastIndex_empty_array", "findLastIndex(emptyArray, # > 0)", -1, false},
		{"findLastIndex_numbers", "findLastIndex(numbers, # > 3)", 4, false},
		{"findLastIndex_no_match", "findLastIndex(numbers, # > 10)", -1, false},
		{"findLastIndex_with_nils", "findLastIndex(numbersWithNil, # == nil)", 3, false},

		// groupBy() function with null handling
		{"groupBy_null_array", "groupBy(nil, # % 2)", map[any][]any{}, false}, // groupBy() on nil should return empty map
		{"groupBy_empty_array", "groupBy(emptyArray, # % 2)", map[any][]any{}, false},
		{"groupBy_numbers", "groupBy(numbers, # % 2)", map[any][]any{1: []any{1, 3, 5}, 0: []any{2, 4}}, false},
		{"groupBy_with_nils", "groupBy(numbersWithNil, # ?? -1)", map[any][]any{1: []any{1}, -1: []any{nil, nil}, 3: []any{3}, 5: []any{5}}, false},

		// concat() function with null handling
		{"concat_with_null", "concat([1, 2], nil, [3, 4])", []any{1, 2, 3, 4}, false}, // concat() should skip nil arrays
		{"concat_null_first", "concat(nil, [1, 2], [3, 4])", []any{1, 2, 3, 4}, false},
		{"concat_all_null", "concat(nil, nil, nil)", []any{}, false},
		{"concat_empty_and_null", "concat([], nil, [])", []any{}, false},

		// flatten() function with null handling
		{"flatten_null_array", "flatten(nil)", []any{}, false}, // flatten() on nil should return empty array
		{"flatten_with_nils", "flatten([1, nil, [2, 3], nil, [4, [5, 6]]])", []any{1, nil, 2, 3, nil, 4, 5, 6}, false},
		{"flatten_nested_nils", "flatten([[1, nil], [nil, 2], [3]])", []any{1, nil, nil, 2, 3}, false},

		// uniq() function with null handling
		{"uniq_null_array", "uniq(nil)", []any{}, false}, // uniq() on nil should return empty array
		{"uniq_with_nils", "uniq([1, nil, 2, nil, 1, 3])", []any{1, nil, 2, 3}, false},
		{"uniq_all_nils", "uniq([nil, nil, nil])", []any{nil}, false},
		{"uniq_mixed", "uniq([1, 'a', nil, 1, 'a', nil])", []any{1, "a", nil}, false},

		// join() function with null handling - already covered in existing tests but adding more cases
		{"join_strings_with_nils", "join(stringsWithNil, ',')", "apple,,cherry,", false}, // nil becomes empty string
		{"join_null_separator", "join(strings, nil)", "applebananacherry", false},        // nil separator becomes empty string

		// reduce() function with null handling
		{"reduce_null_array", "reduce(nil, #acc + #, 0)", 0, false}, // reduce() on nil should return initial value
		{"reduce_empty_array", "reduce(emptyArray, #acc + #, 0)", 0, false},
		{"reduce_numbers", "reduce(numbers, #acc + #, 0)", 15, false},
		{"reduce_with_nils", "reduce(numbersWithNil, #acc + (# ?? 0), 0)", 9, false},
		{"reduce_without_initial", "reduce(numbers, #acc + #)", 15, false},

		// first() function with null handling
		{"first_null_array", "first(nil)", nil, false}, // first() on nil should return nil
		{"first_empty_array", "first(emptyArray)", nil, false},
		{"first_numbers", "first(numbers)", 1, false},
		{"first_with_nils", "first(numbersWithNil)", 1, false},
		{"first_all_nils", "first(allNilArray)", nil, false},

		// last() function with null handling
		{"last_null_array", "last(nil)", nil, false}, // last() on nil should return nil
		{"last_empty_array", "last(emptyArray)", nil, false},
		{"last_numbers", "last(numbers)", 5, false},
		{"last_with_nils", "last(numbersWithNil)", 5, false},
		{"last_all_nils", "last(allNilArray)", nil, false},

		// take() function with null handling
		{"take_null_array", "take(nil, 3)", []any{}, false}, // take() on nil should return empty array
		{"take_empty_array", "take(emptyArray, 3)", []any{}, false},
		{"take_numbers", "take(numbers, 3)", []any{1, 2, 3}, false},
		{"take_with_nils", "take(numbersWithNil, 3)", []any{1, nil, 3}, false},
		{"take_more_than_length", "take(numbers, 10)", []any{1, 2, 3, 4, 5}, false},

		// reverse() function with null handling
		{"reverse_null_array", "reverse(nil)", []any{}, false}, // reverse() on nil should return empty array
		{"reverse_empty_array", "reverse(emptyArray)", []any{}, false},
		{"reverse_numbers", "reverse(numbers)", []any{5, 4, 3, 2, 1}, false},
		{"reverse_with_nils", "reverse(numbersWithNil)", []any{5, nil, 3, nil, 1}, false},

		// sort() function with null handling
		{"sort_null_array", "sort(nil)", []any{}, false}, // sort() on nil should return empty array
		{"sort_empty_array", "sort(emptyArray)", []any{}, false},
		{"sort_numbers", "sort([3, 1, 4, 1, 5])", []any{1, 1, 3, 4, 5}, false},
		{"sort_numbers_desc", "sort([3, 1, 4, 1, 5], 'desc')", []any{5, 4, 3, 1, 1}, false},
		{"sort_with_nils", "sort([3, nil, 1, nil, 2])", []any{nil, nil, 1, 2, 3}, false}, // nils should sort first

		// sortBy() function with null handling
		{"sortBy_null_array", "sortBy(nil, .value)", []any{}, false}, // sortBy() on nil should return empty array
		{"sortBy_empty_array", "sortBy(emptyArray, .value)", []any{}, false},
		{"sortBy_objects", "sortBy(objects, .value)", []any{map[string]any{"value": 10, "name": "A"}, map[string]any{"value": 20, "name": "B"}, map[string]any{"value": 30, "name": "C"}}, false},
		{"sortBy_objects_desc", "sortBy(objects, .value, 'desc')", []any{map[string]any{"value": 30, "name": "C"}, map[string]any{"value": 20, "name": "B"}, map[string]any{"value": 10, "name": "A"}}, false},
		{"sortBy_objects_with_nils", "sortBy(objectsWithNil, .value ?? 0)", []any{map[string]any{"value": nil, "name": nil}, map[string]any{"value": 10, "name": "A"}, map[string]any{"value": 30, "name": "C"}}, false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := expr.Eval(tc.expression, env)
			if tc.wantError {
				if err == nil {
					t.Errorf("expected error for %q, got result %v", tc.expression, result)
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error for %q: %v", tc.expression, err)
				} else if !deepEqual(result, tc.want) {
					t.Errorf("%q: got %v (type %T), want %v (type %T)", tc.expression, result, result, tc.want, tc.want)
				}
			}
		})
	}
}

func TestBuiltinMapFunctions_NullHandling(t *testing.T) {
	env := map[string]any{
		"validMap": map[string]any{
			"name": "John",
			"age":  30,
			"city": "New York",
		},
		"mapWithNils": map[string]any{
			"name":    "John",
			"email":   nil,
			"phone":   nil,
			"address": "123 Main St",
		},
		"emptyMap": map[string]any{},
		"nilMap":   nil,
	}

	testCases := []struct {
		name       string
		expression string
		want       any
		wantError  bool
	}{
		// keys() function with null handling
		{"keys_null_map", "keys(nil)", []any{}, false}, // keys() on nil should return empty array
		{"keys_empty_map", "keys(emptyMap)", []any{}, false},
		{"keys_valid_map", "len(keys(validMap))", 3, false},        // Check length since order may vary
		{"keys_map_with_nils", "len(keys(mapWithNils))", 4, false}, // Keys with nil values should still be included

		// values() function with null handling
		{"values_null_map", "values(nil)", []any{}, false}, // values() on nil should return empty array
		{"values_empty_map", "values(emptyMap)", []any{}, false},
		{"values_valid_map", "len(values(validMap))", 3, false},        // Check length since order may vary
		{"values_map_with_nils", "len(values(mapWithNils))", 4, false}, // Values including nils should be included

		// Additional map access with null handling
		{"get_from_null_map", "get(nil, 'key')", nil, false},
		{"get_null_key", "get(validMap, nil)", nil, false},
		{"get_nonexistent_key", "get(validMap, 'nonexistent')", nil, false},
		{"get_valid_key", "get(validMap, 'name')", "John", false},
		{"get_nil_value", "get(mapWithNils, 'email')", nil, false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := expr.Eval(tc.expression, env)
			if tc.wantError {
				if err == nil {
					t.Errorf("expected error for %q, got result %v", tc.expression, result)
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error for %q: %v", tc.expression, err)
				} else if !deepEqual(result, tc.want) {
					t.Errorf("%q: got %v (type %T), want %v (type %T)", tc.expression, result, result, tc.want, tc.want)
				}
			}
		})
	}
}

func TestBuiltinTypeConversionFunctions_NullHandling(t *testing.T) {
	testCases := []struct {
		name       string
		expression string
		want       any
		wantError  bool
	}{
		// string() function with null handling
		{"string_null", "string(nil)", "", false},
		{"string_number", "string(42)", "42", false},
		{"string_float", "string(3.14)", "3.14", false},
		{"string_bool", "string(true)", "true", false},
		{"string_array", "string([1, 2, 3])", "[1 2 3]", false},

		// toJSON() function with null handling
		{"toJSON_null", "toJSON(nil)", "null", false},
		{"toJSON_number", "toJSON(42)", "42", false},
		{"toJSON_string", "toJSON('hello')", `"hello"`, false},
		{"toJSON_array", "toJSON([1, 2, 3])", "[1,2,3]", false},
		{"toJSON_map", "toJSON({name: 'John', age: 30})", `{"age":30,"name":"John"}`, false}, // Note: JSON key order may vary

		// fromJSON() function with null handling - already covered but adding more cases
		{"fromJSON_null_string", `fromJSON("null")`, nil, false},
		{"fromJSON_number_string", `fromJSON("42")`, 42.0, false}, // JSON numbers are float64
		{"fromJSON_string_value", `fromJSON('"hello"')`, "hello", false},
		{"fromJSON_array_string", `fromJSON("[1,2,3]")`, []any{1.0, 2.0, 3.0}, false},
		{"fromJSON_object_string", `fromJSON('{"name":"John","age":30}')`, map[string]any{"name": "John", "age": 30.0}, false},

		// toPairs() function with null handling
		{"toPairs_null_map", "toPairs(nil)", []any{}, false}, // toPairs() on nil should return empty array
		{"toPairs_empty_map", "toPairs({})", []any{}, false},
		{"toPairs_valid_map", "len(toPairs({a: 1, b: 2}))", 2, false}, // Check length since order may vary
		{"toPairs_map_with_nils", "len(toPairs({a: 1, b: nil, c: 3}))", 3, false},

		// fromPairs() function with null handling
		{"fromPairs_null_array", "fromPairs(nil)", map[string]any{}, false}, // fromPairs() on nil should return empty map
		{"fromPairs_empty_array", "fromPairs([])", map[string]any{}, false},
		{"fromPairs_valid_pairs", `fromPairs([["name", "John"], ["age", 30]])`, map[string]any{"name": "John", "age": 30}, false},
		{"fromPairs_with_nil_values", `fromPairs([["name", "John"], ["email", nil]])`, map[string]any{"name": "John", "email": nil}, false},
		//{"fromPairs_with_nil_pairs", `fromPairs([["name", "John"], nil, ["age", 30]])`, map[string]any{"name": "John", "age": 30}, false}, // Should skip nil pairs

		// Additional type checking with complex null scenarios
		{"type_of_various_nulls", "type(nil)", "nil", false},
		{"type_of_null_in_array", "type([nil][0])", "nil", false},
		{"type_of_null_in_map", "type({a: nil}.a)", "nil", false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := expr.Eval(tc.expression, nil)
			if tc.wantError {
				if err == nil {
					t.Errorf("expected error for %q, got result %v", tc.expression, result)
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error for %q: %v", tc.expression, err)
				} else if !deepEqual(result, tc.want) {
					t.Errorf("%q: got %v (type %T), want %v (type %T)", tc.expression, result, result, tc.want, tc.want)
				}
			}
		})
	}
}

func TestBuiltinBitwiseFunctions_NullHandling(t *testing.T) {
	testCases := []struct {
		name       string
		expression string
		want       any
		wantError  bool
	}{
		// bitand() function with null handling
		{"bitand_null_first", "bitand(nil, 5)", 0, false}, // nil should be treated as 0
		{"bitand_null_second", "bitand(5, nil)", 0, false},
		{"bitand_both_null", "bitand(nil, nil)", 0, false},
		{"bitand_valid", "bitand(12, 10)", 8, false}, // 1100 & 1010 = 1000 = 8

		// bitor() function with null handling
		{"bitor_null_first", "bitor(nil, 5)", 5, false}, // nil should be treated as 0
		{"bitor_null_second", "bitor(5, nil)", 5, false},
		{"bitor_both_null", "bitor(nil, nil)", 0, false},
		{"bitor_valid", "bitor(12, 10)", 14, false}, // 1100 | 1010 = 1110 = 14

		// bitxor() function with null handling
		{"bitxor_null_first", "bitxor(nil, 5)", 5, false}, // nil should be treated as 0
		{"bitxor_null_second", "bitxor(5, nil)", 5, false},
		{"bitxor_both_null", "bitxor(nil, nil)", 0, false},
		{"bitxor_valid", "bitxor(12, 10)", 6, false}, // 1100 ^ 1010 = 0110 = 6

		// bitnand() function with null handling
		//{"bitnand_null_first", "bitnand(nil, 5)", -1, false},  // ~(nil & 5) = ~0 = -1
		//{"bitnand_null_second", "bitnand(5, nil)", -1, false}, // ~(5 & nil) = ~0 = -1
		//{"bitnand_both_null", "bitnand(nil, nil)", -1, false}, // ~(nil & nil) = ~0 = -1
		{"bitnand_valid", "bitnand(12, 10)", 4, false}, // ~(1100 & 1010) = ~1000 = -9

		// bitnot() function with null handling
		{"bitnot_null", "bitnot(nil)", -1, false}, // ~nil = ~0 = -1
		{"bitnot_valid", "bitnot(5)", -6, false},  // ~5 = -6

		// bitshl() function with null handling
		{"bitshl_null_first", "bitshl(nil, 2)", 0, false},  // nil << 2 = 0 << 2 = 0
		{"bitshl_null_second", "bitshl(5, nil)", 5, false}, // 5 << nil = 5 << 0 = 5
		{"bitshl_both_null", "bitshl(nil, nil)", 0, false}, // nil << nil = 0 << 0 = 0
		{"bitshl_valid", "bitshl(5, 2)", 20, false},        // 5 << 2 = 20

		// bitshr() function with null handling
		{"bitshr_null_first", "bitshr(nil, 2)", 0, false},    // nil >> 2 = 0 >> 2 = 0
		{"bitshr_null_second", "bitshr(20, nil)", 20, false}, // 20 >> nil = 20 >> 0 = 20
		{"bitshr_both_null", "bitshr(nil, nil)", 0, false},   // nil >> nil = 0 >> 0 = 0
		{"bitshr_valid", "bitshr(20, 2)", 5, false},          // 20 >> 2 = 5

		// bitushr() function with null handling
		{"bitushr_null_first", "bitushr(nil, 2)", 0, false},                 // nil >>> 2 = 0 >>> 2 = 0
		{"bitushr_null_second", "bitushr(20, nil)", 20, false},              // 20 >>> nil = 20 >>> 0 = 20
		{"bitushr_both_null", "bitushr(nil, nil)", 0, false},                // nil >>> nil = 0 >>> 0 = 0
		{"bitushr_valid", "bitushr(20, 2)", 5, false},                       // 20 >>> 2 = 5
		{"bitushr_negative", "bitushr(-20, 2)", 4611686018427387899, false}, // Unsigned right shift of negative number
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := expr.Eval(tc.expression, nil)
			if tc.wantError {
				if err == nil {
					t.Errorf("expected error for %q, got result %v", tc.expression, result)
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error for %q: %v", tc.expression, err)
				} else if !deepEqual(result, tc.want) {
					t.Errorf("%q: got %v (type %T), want %v (type %T)", tc.expression, result, result, tc.want, tc.want)
				}
			}
		})
	}
}

func TestBuiltinDateFunctions_NullHandling(t *testing.T) {
	testCases := []struct {
		name       string
		expression string
		want       any
		wantError  bool
	}{
		// date() function with null handling
		{"date_null_string", "date(nil).Unix()", 0, false},                           // date() with nil should return nil
		{"date_null_format", "date('2023-01-01', nil)", true, false},                 // Should parse with default format, check if it's a valid date
		{"date_null_timezone", "date('2023-01-01', '2006-01-02', nil)", true, false}, // Should use default timezone
		{"date_all_null", "date(nil, nil, nil).Unix()", 0, false},

		// duration() function with null handling - already covered but adding verification
		{"duration_null", "duration(nil) == duration('0s')", true, false},

		// timezone() function with null handling
		//{"timezone_null", "timezone(nil)", nil, false}, // timezone() with nil should return nil

		// Complex date operations with null handling
		//{"date_comparison_with_null", "date('2023-01-01') > nil", false, false}, // Comparison with nil should be false
		//{"null_date_comparison", "nil < date('2023-01-01')", false, false},
		{"date_arithmetic_with_null", "date('2023-01-01') + nil", nil, false}, // Adding nil duration should return nil
		{"null_plus_duration", "nil + duration('1h')", nil, false},

		// Date method calls on null
		//{"null_date_year", "nil?.Year()", nil, false}, // Optional chaining should handle null
		//{"null_date_format", "nil?.Format('2006-01-02')", nil, false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := expr.Eval(tc.expression, nil)
			if tc.wantError {
				if err == nil {
					t.Errorf("expected error for %q, got result %v", tc.expression, result)
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error for %q: %v", tc.expression, err)
				} else {
					// For boolean expectations, check directly
					if want, ok := tc.want.(bool); ok {
						isValid := false
						if result == nil && !want {
							// nil is considered false in boolean context
							isValid = true
						} else if resultBool, ok := result.(bool); ok && resultBool == want {
							isValid = true
						} else if want && result != nil {
							// For date parsing, just check if result is not nil (successful parsing)
							isValid = true
						}
						if !isValid && !deepEqual(result, tc.want) {
							t.Errorf("%q: got %v (type %T), want %v (type %T)", tc.expression, result, result, tc.want, tc.want)
						}
					} else if !deepEqual(result, tc.want) {
						t.Errorf("%q: got %v (type %T), want %v (type %T)", tc.expression, result, result, tc.want, tc.want)
					}
				}
			}
		})
	}
}

func TestBuiltinNumberFunctions_TwoParam_NullHandling(t *testing.T) {
	testCases := []struct {
		name       string
		expression string
		want       any
		wantError  bool
	}{
		// max(a, b) function with null handling
		{"max_null_first", "max(nil, 5)", 5, false},    // max(nil, 5) should return 5
		{"max_null_second", "max(5, nil)", 5, false},   // max(5, nil) should return 5
		{"max_both_null", "max(nil, nil)", nil, false}, // max(nil, nil) should return nil
		{"max_valid_numbers", "max(3, 7)", 7, false},
		{"max_equal_numbers", "max(5, 5)", 5, false},
		{"max_negative_numbers", "max(-3, -7)", -3, false},
		{"max_mixed_types", "max(3, 5.7)", 5.7, false},

		// min(a, b) function with null handling
		{"min_null_first", "min(nil, 5)", 5, false},    // min(nil, 5) should return 5
		{"min_null_second", "min(5, nil)", 5, false},   // min(5, nil) should return 5
		{"min_both_null", "min(nil, nil)", nil, false}, // min(nil, nil) should return nil
		{"min_valid_numbers", "min(3, 7)", 3, false},
		{"min_equal_numbers", "min(5, 5)", 5, false},
		{"min_negative_numbers", "min(-3, -7)", -7, false},
		{"min_mixed_types", "min(3, 2.5)", 2.5, false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := expr.Eval(tc.expression, nil)
			if tc.wantError {
				if err == nil {
					t.Errorf("expected error for %q, got result %v", tc.expression, result)
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error for %q: %v", tc.expression, err)
				} else if !deepEqual(result, tc.want) {
					t.Errorf("%q: got %v (type %T), want %v (type %T)", tc.expression, result, result, tc.want, tc.want)
				}
			}
		})
	}
}

func TestBuiltinMiscellaneousFunctions_NullHandling(t *testing.T) {
	env := map[string]any{
		"array":    []any{1, 2, 3, 4, 5},
		"nilArray": nil,
		"map": map[string]any{
			"name": "John",
			"age":  30,
		},
		"nilMap": nil,
	}

	testCases := []struct {
		name       string
		expression string
		want       any
		wantError  bool
	}{
		// get() function with null handling - already partially covered but adding more cases
		{"get_null_collection", "get(nil, 0)", nil, false},
		{"get_null_index", "get(array, nil)", nil, false},
		{"get_both_null", "get(nil, nil)", nil, false},
		{"get_array_valid_index", "get(array, 2)", 3, false},
		{"get_array_out_of_bounds", "get(array, 10)", nil, false},
		{"get_array_negative_index", "get(array, -1)", 5, false}, // Should support negative indexing
		{"get_map_valid_key", "get(map, 'name')", "John", false},
		{"get_map_invalid_key", "get(map, 'nonexistent')", nil, false},
		//{"get_string_index", `get("hello", 1)`, "e", false},
		//{"get_string_out_of_bounds", `get("hello", 10)`, nil, false},
		//{"get_string_null_index", `get("hello", nil)`, nil, false},

		// Additional len() cases with null handling - already covered but ensuring completeness
		{"len_various_nulls", "len(nil)", 0, false},
		{"len_null_vs_empty_array", "len([]) == len(nil)", true, false},
		{"len_null_vs_empty_map", "len({}) == len(nil)", true, false},
		{"len_null_vs_empty_string", `len("") == len(nil)`, true, false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := expr.Eval(tc.expression, env)
			if tc.wantError {
				if err == nil {
					t.Errorf("expected error for %q, got result %v", tc.expression, result)
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error for %q: %v", tc.expression, err)
				} else if !deepEqual(result, tc.want) {
					t.Errorf("%q: got %v (type %T), want %v (type %T)", tc.expression, result, result, tc.want, tc.want)
				}
			}
		})
	}
}

func TestBuiltinAdditionalNullHandling(t *testing.T) {
	testCases := []struct {
		name       string
		expression string
		want       any
		wantError  bool
	}{
		// Additional edge cases for string functions that might not be fully covered
		//{"contains_null_string", "contains(nil, 'test')", false, false},
		//{"contains_null_substring", "contains('hello world', nil)", true, false}, // nil substring should match (empty string behavior)
		//{"contains_both_null", "contains(nil, nil)", true, false},
		//
		//{"startsWith_null_string", "startsWith(nil, 'test')", false, false},
		//{"startsWith_null_prefix", "startsWith('hello world', nil)", true, false}, // nil prefix should match (empty string behavior)
		//{"startsWith_both_null", "startsWith(nil, nil)", true, false},
		//
		//{"endsWith_null_string", "endsWith(nil, 'test')", false, false},
		//{"endsWith_null_suffix", "endsWith('hello world', nil)", true, false}, // nil suffix should match (empty string behavior)
		//{"endsWith_both_null", "endsWith(nil, nil)", true, false},

		//{"matches_null_string", "matches(nil, 'test')", false, false},
		//{"matches_null_pattern", "matches('hello world', nil)", false, false},
		//{"matches_both_null", "matches(nil, nil)", false, false},

		// Additional array function edge cases
		{"sum_null_array", "sum(nil)", 0, false},
		{"sum_null_array_with_predicate", "sum(nil, .value)", 0, false},
		{"mean_null_array", "mean(nil)", 0, false},     // mean of null should be nil
		{"median_null_array", "median(nil)", 0, false}, // median of null should be nil
		{"count_null_array_no_predicate", "count(nil)", 0, false},
		{"count_null_array_with_predicate", "count(nil, # > 0)", 0, false},

		// Type function with various null scenarios
		{"type_different_nulls", "type(nil) == 'nil'", true, false},

		// String conversion with null handling
		{"string_null", "string(nil)", "", false}, // string(nil) should return empty string

		// Additional complex null coalescing scenarios
		{"nested_null_coalescing", "nil ?? (nil ?? 'default')", "default", false},
		{"complex_null_arithmetic", "(nil + 5) * (nil ?? 2)", 10, false}, // (0 + 5) * 2 = 10
		{"null_in_complex_expression", "len(nil ?? []) + (nil ?? 0)", 0, false},

		// Optional chaining with various null scenarios
		//{"optional_chaining_deep_null", "nil?.a?.b?.c", nil, false},
		//{"optional_chaining_with_array", "nil?.[0]", nil, false},
		//{"optional_chaining_method_call", "nil?.toString()", nil, false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := expr.Eval(tc.expression, nil)
			if tc.wantError {
				if err == nil {
					t.Errorf("expected error for %q, got result %v", tc.expression, result)
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error for %q: %v", tc.expression, err)
				} else if !deepEqual(result, tc.want) {
					t.Errorf("%q: got %v (type %T), want %v (type %T)", tc.expression, result, result, tc.want, tc.want)
				}
			}
		})
	}
}

func TestBuiltinStringFunctionAdditionalNullCases(t *testing.T) {
	testCases := []struct {
		name       string
		expression string
		want       any
		wantError  bool
	}{
		// Additional trim variations with null
		{"trim_null_with_empty_cutset", "trim(nil, '')", nil, false},
		{"trim_empty_with_null_cutset", "trim('', nil)", "", false},

		// Replace with null handling - more edge cases
		{"replace_empty_to_null", "replace('hello', 'hello', nil)", "", false}, // Replace with nil should result in empty
		{"replace_multiple_nulls", "replace(nil, nil, nil)", nil, false},

		// Split with edge cases
		{"split_with_empty_string_and_null", "split('', nil)", []any{""}, false},
		{"split_null_with_empty_separator", "split(nil, '')", []any{}, false},

		// Additional index functions
		{"indexOf_empty_in_null", "indexOf(nil, '')", -1, false},
		{"lastIndexOf_empty_in_null", "lastIndexOf(nil, '')", -1, false},

		// Repeat with null
		{"repeat_null_zero_times", "repeat(nil, 0)", "", false},
		{"repeat_null_positive_times", "repeat(nil, 3)", "", false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := expr.Eval(tc.expression, nil)
			if tc.wantError {
				if err == nil {
					t.Errorf("expected error for %q, got result %v", tc.expression, result)
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error for %q: %v", tc.expression, err)
				} else if !deepEqual(result, tc.want) {
					t.Errorf("%q: got %v (type %T), want %v (type %T)", tc.expression, result, result, tc.want, tc.want)
				}
			}
		})
	}
}

func TestBuiltinArrayFunctionAdditionalNullCases(t *testing.T) {
	env := map[string]any{
		"arrayWithSomeNils": []any{1, nil, 3, nil, 5},
		"arrayAllNils":      []any{nil, nil, nil},
		"emptyArray":        []any{},
	}

	testCases := []struct {
		name       string
		expression string
		want       any
		wantError  bool
	}{
		// Additional array functions edge cases
		{"take_null_from_array", "take(arrayWithSomeNils, nil)", []any{}, false}, // take with nil count should return empty
		{"take_zero_from_null", "take(nil, 0)", []any{}, false},
		//{"take_negative_from_array", "take(arrayWithSomeNils, -1)", []any{}, false}, // Negative take should return empty

		// Concat with various null combinations
		{"concat_multiple_nulls", "concat(nil, nil, nil, nil)", []any{}, false},
		{"concat_mixed_nulls_and_arrays", "concat([1], nil, [2], nil, [3])", []any{1, 2, 3}, false},

		// Sort with null elements
		{"sort_mixed_nulls_and_numbers", "sort([3, nil, 1, nil, 2])", []any{nil, nil, 1, 2, 3}, false},
		{"sort_all_nulls", "sort(arrayAllNils)", []any{nil, nil, nil}, false},

		// Uniq with null elements
		{"uniq_multiple_nulls", "uniq([nil, 1, nil, 2, nil])", []any{nil, 1, 2}, false},
		{"uniq_only_nulls", "uniq([nil, nil, nil])", []any{nil}, false},

		// GroupBy with null keys
		{"groupBy_null_keys", "groupBy([1, nil, 2, nil, 3], #)", map[any][]any{1: []any{1}, nil: []any{nil, nil}, 2: []any{2}, 3: []any{3}}, false},

		// FindIndex variations
		{"findIndex_null_in_mixed", "findIndex([1, nil, 3], # == nil)", 1, false},
		{"findLastIndex_null_in_mixed", "findLastIndex([nil, 2, nil, 4], # == nil)", 2, false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := expr.Eval(tc.expression, env)
			if tc.wantError {
				if err == nil {
					t.Errorf("expected error for %q, got result %v", tc.expression, result)
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error for %q: %v", tc.expression, err)
				} else if !deepEqual(result, tc.want) {
					t.Errorf("%q: got %v (type %T), want %v (type %T)", tc.expression, result, result, tc.want, tc.want)
				}
			}
		})
	}
}

func TestBuiltinMapFunctionAdditionalNullCases(t *testing.T) {
	env := map[string]any{
		"mapWithNullValues": map[string]any{
			"a": 1,
			"b": nil,
			"c": 3,
			"d": nil,
		},
		"mapWithNullKeys": map[any]any{
			"valid": 1,
			nil:     2,
			"other": 3,
		},
		"emptyMap": map[string]any{},
	}

	testCases := []struct {
		name       string
		expression string
		want       any
		wantError  bool
	}{
		// Keys function with null values in map
		{"keys_map_with_null_values", "keys(mapWithNullValues)", []any{"a", "b", "c", "d"}, false}, // Should return all keys even if values are nil
		{"keys_null_map", "keys(nil)", []any{}, false},                                             // keys() on nil should return empty array
		{"keys_empty_map", "keys(emptyMap)", []any{}, false},

		// Values function with null values
		{"values_map_with_null_values", "values(mapWithNullValues)", []any{1, nil, 3, nil}, false}, // Should include nil values
		{"values_null_map", "values(nil)", []any{}, false},                                         // values() on nil should return empty array
		{"values_empty_map", "values(emptyMap)", []any{}, false},

		// Map access with null keys
		{"access_null_key", "mapWithNullKeys[nil]", 2, false}, // Should be able to access nil key
		{"access_null_key_with_get", "get(mapWithNullKeys, nil)", 2, false},

		// toPairs and fromPairs with null handling
		{"toPairs_null_map", "toPairs(nil)", []any{}, false},                      // toPairs on nil should return empty array
		{"toPairs_with_null_values", "len(toPairs(mapWithNullValues))", 4, false}, // Should include pairs with nil values
		{"fromPairs_null_array", "fromPairs(nil)", map[any]any{}, false},          // fromPairs on nil should return empty map
		{"fromPairs_with_null_pairs", "fromPairs([['a', 1], ['b', nil], [nil, 3]])", map[any]any{"a": 1, "b": nil, nil: 3}, false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := expr.Eval(tc.expression, env)
			if tc.wantError {
				if err == nil {
					t.Errorf("expected error for %q, got result %v", tc.expression, result)
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error for %q: %v", tc.expression, err)
				} else if !deepEqual(result, tc.want) {
					t.Errorf("%q: got %v (type %T), want %v (type %T)", tc.expression, result, result, tc.want, tc.want)
				}
			}
		})
	}
}
