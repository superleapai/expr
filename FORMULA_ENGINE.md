# Formula Engine Documentation

A Salesforce-compatible formula engine built on expr-lang. Supports 77 built-in functions,
Salesforce operators, null-safe evaluation, and runtime context functions.

---

## Table of Contents

- [Converting Salesforce Formulas](#converting-salesforce-formulas)
- [Quick Start](#quick-start)
- [Full Validation Pipeline](#full-validation-pipeline)
- [Operators](#operators)
- [Keywords](#keywords)
- [Type Coercion](#type-coercion)
- [Function Reference](#function-reference)
  - [Logical Functions](#logical-functions)
  - [Math Functions](#math-functions)
  - [Trigonometric Functions](#trigonometric-functions)
  - [Text Functions](#text-functions)
  - [Date/Time Functions](#datetime-functions)
  - [Regex Functions](#regex-functions)
  - [Encoding Functions](#encoding-functions)
  - [Salesforce Runtime Functions](#salesforce-runtime-functions)
- [ExtractDeps](#extractdeps)
- [InflateEnv](#inflateenv)
- [Null Handling](#null-handling)
- [Date Arithmetic](#date-arithmetic)
- [Common Patterns](#common-patterns)
- [Performance](#performance)

---

## Converting Salesforce Formulas

When migrating validation rules from Salesforce, you need to apply these syntax changes.
Most functions (AND, OR, IF, ISBLANK, ISPICKVAL, etc.) work identically — only the
operators and variable references need adjustment.

### Required Changes

#### 1. Equality: `=` -> `==`

Salesforce uses single `=` for equality. This engine uses `==`.

```
// Salesforce:
Status = "Active"

// This engine:
Status == "Active"
```

#### 2. Not-Equal: both `<>` and `!=` work

No change needed — `<>` is supported natively. `!=` also works.

```
// Salesforce:
Status <> "Closed"

// This engine (both work):
Status <> "Closed"
Status != "Closed"
```

#### 3. String Concatenation: both `&` and `+` work

No change needed — `&` is supported natively.

```
// Salesforce:
FirstName & " " & LastName

// This engine (works as-is):
FirstName & " " & LastName
```

#### 4. Field References: remove `$` prefix from global variables

Salesforce uses `$User.Email`, `$Profile.Name`, etc. Convert these to regular
field paths or flat env keys.

```
// Salesforce:
$User.Email
$Profile.Name
$User.Id

// This engine (option A: dotted paths, use InflateEnv):
User.Email
Profile.Name
User.Id

// This engine (option B: flat keys, no InflateEnv needed):
User_Email
Profile_Name
User_Id
```

#### 5. ISCHANGED / PRIORVALUE: use quoted string field names

In Salesforce, these take field references: `ISCHANGED(Status)`.
In this engine, they take **string literals**: `ISCHANGED("Status")`.

```
// Salesforce:
ISCHANGED(Status)
PRIORVALUE(OwnerId)

// This engine:
ISCHANGED("Status")
PRIORVALUE("OwnerId")
```

This is because the engine looks up the field name in the `_changed` / `_prior`
maps at runtime, rather than having a compiler-level field reference.

#### 6. Null literal: both `null` and `NULL` work

```
// Salesforce:
Amount = null

// This engine:
Amount == NULL
Amount == nil     // also works
```

#### 7. Boolean literals: both cases work

```
// Salesforce:
TRUE, FALSE

// This engine (all work):
TRUE, FALSE
true, false
True, False
```

#### 8. Cross-object field references: keep the dot notation

```
// Salesforce:
Account.Name
Account.Owner.Email
Contact.Account.Industry

// This engine (same syntax, but env must use InflateEnv):
Account.Name
Account.Owner.Email
Contact.Account.Industry
```

Make sure to pass these as dotted keys to `InflateEnv`:
```go
flat := map[string]any{
    "Account.Name":        "Acme",
    "Account.Owner.Email": "owner@acme.com",
}
env := expr.InflateEnv(flat)
```

### No Changes Needed

These work identically to Salesforce:

| Feature | Example |
|---------|---------|
| AND / OR / NOT | `AND(cond1, cond2)` |
| IF / CASE | `IF(cond, val1, val2)` |
| ISBLANK / ISNULL | `ISBLANK(field)` |
| ISPICKVAL | `ISPICKVAL(Stage, "Won")` |
| INCLUDES | `INCLUDES(MultiSelect, "Val")` |
| ISNEW | `ISNEW()` |
| Math operators | `+`, `-`, `*`, `/` |
| Comparison operators | `<`, `>`, `<=`, `>=` |
| Parentheses | `(a + b) * c` |
| String literals | `"hello"` |
| Number literals | `42`, `3.14` |
| All functions | `LEN()`, `LEFT()`, `CONTAINS()`, etc. |

### Conversion Examples

#### Example 1: Simple validation

```
// Salesforce:
AND(
    ISPICKVAL(StageName, "Closed Won"),
    ISBLANK(Amount)
)

// This engine (no changes needed!):
AND(
    ISPICKVAL(StageName, "Closed Won"),
    ISBLANK(Amount)
)
```

#### Example 2: With equality and global variables

```
// Salesforce:
AND(
    Status = "Active",
    $User.Profile.Name <> "System Administrator",
    Amount > 10000
)

// This engine:
AND(
    Status == "Active",                          // = → ==
    User_Profile_Name <> "System Administrator", // $User.Profile.Name → flat key
    Amount > 10000
)
```

#### Example 3: ISCHANGED / PRIORVALUE

```
// Salesforce:
AND(
    ISCHANGED(OwnerId),
    PRIORVALUE(OwnerId) <> "005000000000001",
    $User.Id = "005000000000002"
)

// This engine:
AND(
    ISCHANGED("OwnerId"),                           // field ref → string
    PRIORVALUE("OwnerId") <> "005000000000001",     // field ref → string
    User_Id == "005000000000002"                     // $User.Id → flat, = → ==
)
```

#### Example 4: Complex nested formula

```
// Salesforce:
IF(
    OR(
        ISPICKVAL(Type, "New Customer"),
        ISPICKVAL(Type, "Existing Customer")
    ),
    AND(
        NOT(ISBLANK(Account.Industry)),
        Amount > 0,
        CloseDate > TODAY()
    ),
    TRUE
)

// This engine (only TODAY needs no change, Account.Industry needs InflateEnv):
IF(
    OR(
        ISPICKVAL(Type, "New Customer"),
        ISPICKVAL(Type, "Existing Customer")
    ),
    AND(
        NOT(ISBLANK(Account.Industry)),
        Amount > 0,
        CloseDate > TODAY()
    ),
    TRUE
)
```

### Quick Reference: Conversion Checklist

| # | Salesforce Syntax | This Engine | Change Required? |
|---|-------------------|-------------|-----------------|
| 1 | `=` (equality) | `==` | **Yes** |
| 2 | `<>` (not equal) | `<>` or `!=` | No |
| 3 | `&` (concat) | `&` | No |
| 4 | `$User.Email` | `User.Email` or `User_Email` | **Yes** (remove `$`) |
| 5 | `ISCHANGED(Field)` | `ISCHANGED("Field")` | **Yes** (add quotes) |
| 6 | `PRIORVALUE(Field)` | `PRIORVALUE("Field")` | **Yes** (add quotes) |
| 7 | `null` | `NULL` or `nil` | No |
| 8 | `TRUE` / `FALSE` | `TRUE` / `FALSE` | No |
| 9 | `Account.Name` | `Account.Name` + `InflateEnv` | No (but env needs nesting) |
| 10 | All functions | Same names | No |

### Automated Conversion

For bulk conversion, you can apply these regex replacements:

```
1. Single = to ==  (but not <= >= <> != ==):
   Replace:  (?<![<>!=])=(?!=)
   With:     ==

2. $Variable references to flat keys:
   Replace:  \$(\w+)\.(\w+)
   With:     $1_$2
   (repeat for deeper nesting like $User.Profile.Name)

3. ISCHANGED(FieldName) to ISCHANGED("FieldName"):
   Replace:  ISCHANGED\((\w+)\)
   With:     ISCHANGED("$1")

4. PRIORVALUE(FieldName) to PRIORVALUE("FieldName"):
   Replace:  PRIORVALUE\((\w+)\)
   With:     PRIORVALUE("$1")
```

---



```go
import "github.com/expr-lang/expr"

// 1. Compile the formula (do this once, reuse for many records)
program, err := expr.Compile(
    `AND(Amount > 10000, NOT(ISBLANK(Email)))`,
    expr.AllowUndefinedVariables(),
    expr.WithAllFormulaPacks(),
)

// 2. Run against a record
env := map[string]any{
    "Amount": float64(50000),
    "Email":  "user@example.com",
}
result, err := expr.Run(program, env)
// result = true
```

### Compile Options

| Option | Purpose |
|--------|---------|
| `expr.WithAllFormulaPacks()` | Register all 77 formula functions |
| `expr.AllowUndefinedVariables()` | Don't error on fields missing from env (treat as nil) |
| `expr.Env(env)` | Provide a sample env for compile-time type checking |
| `expr.WithLogicalFunctions()` | Register only logical functions |
| `expr.WithMathFunctions()` | Register only math functions |
| `expr.WithTrigFunctions()` | Register only trig functions |
| `expr.WithTextFunctions()` | Register only text functions |
| `expr.WithDateTimeFunctions()` | Register only date/time functions |
| `expr.WithRegexFunctions()` | Register only REGEX |
| `expr.WithEncodingFunctions()` | Register only encoding functions |
| `expr.WithSalesforceRuntimeFunctions()` | Register only SF runtime functions |

---

## Full Validation Pipeline

This is the recommended production pattern for evaluating Salesforce-style validation rules.

```go
// ================================================================
// STEP 1: Extract dependencies from all formulas
// ================================================================
// This tells you exactly which fields, relations, and old values
// you need to fetch from the database.

formula := `AND(
    ISCHANGED("StageName"),
    ISPICKVAL(StageName, "Closed Won"),
    ISBLANK(Amount),
    Account.Industry <> "Government"
)`

deps, err := expr.ExtractDeps(formula)
// deps.RecordFields  = ["Amount", "StageName"]
// deps.RelatedFields = ["Account.Industry"]
// deps.OldFields     = ["StageName"]
// deps.NeedsIsNew    = false
// deps.NeedsTimezone = false

// ================================================================
// STEP 2: Fetch field values from your database
// ================================================================
// Use deps to build your query — only fetch what's needed.
// Return values as a flat map. Dotted keys are OK.

flat := map[string]any{
    "Amount":           nil,
    "StageName":        "Closed Won",
    "Account.Industry": "Technology",
}

// ================================================================
// STEP 3: Add runtime context for ISCHANGED/PRIORVALUE/ISNEW
// ================================================================
if len(deps.OldFields) > 0 {
    flat["_changed"] = map[string]bool{"StageName": true}
    flat["_prior"]   = map[string]any{"StageName": "Negotiation"}
}
if deps.NeedsIsNew {
    flat["_isNew"] = false // true for insert, false for update
}

// ================================================================
// STEP 4: Inflate dotted paths into nested maps
// ================================================================
// "Account.Industry" becomes {"Account": {"Industry": "Technology"}}
env := expr.InflateEnv(flat)

// ================================================================
// STEP 5: Compile (once per formula, cache and reuse)
// ================================================================
program, err := expr.Compile(formula,
    expr.Env(env),
    expr.AllowUndefinedVariables(),
    expr.WithAllFormulaPacks(),
)

// ================================================================
// STEP 6: Run against each record
// ================================================================
result, err := expr.Run(program, env)
// result == true means validation error fires
```

### Multiple Rules x Multiple Records

```go
// Pre-compile all rules once
programs := make([]*vm.Program, len(rules))
for i, rule := range rules {
    programs[i], _ = expr.Compile(rule.Formula,
        expr.Env(sampleEnv),
        expr.AllowUndefinedVariables(),
        expr.WithAllFormulaPacks(),
    )
}

// Evaluate every rule against every record
for _, record := range records {
    env := expr.InflateEnv(buildFlatEnv(record))
    for _, program := range programs {
        result, _ := expr.Run(program, env)
        if result == true {
            // Validation error — collect error message
        }
    }
}
```

---

## Operators

### Comparison Operators

| Operator | Description | Example |
|----------|-------------|---------|
| `==` | Equal | `Status == "Active"` |
| `!=` | Not equal | `Status != "Closed"` |
| `<>` | Not equal (SF-style) | `Status <> "Closed"` |
| `<` | Less than | `Amount < 1000` |
| `>` | Greater than | `Amount > 1000` |
| `<=` | Less or equal | `Amount <= 1000` |
| `>=` | Greater or equal | `Amount >= 1000` |

### Logical Operators

| Operator | Description | Example |
|----------|-------------|---------|
| `&&` | Logical AND | `a > 0 && b > 0` |
| `\|\|` | Logical OR | `a > 0 \|\| b > 0` |
| `!` | Logical NOT | `!ISBLANK(Name)` |

### Arithmetic Operators

| Operator | Description | Example |
|----------|-------------|---------|
| `+` | Add (or string concat if left is string) | `Amount + Tax` |
| `-` | Subtract | `Total - Discount` |
| `*` | Multiply | `Price * Quantity` |
| `/` | Divide (returns nil for divide-by-zero) | `Total / Count` |
| `%` | Modulo | `Value % 2` |

### String Concatenation

| Operator | Description | Example |
|----------|-------------|---------|
| `&` | String concatenation (SF-style) | `FirstName & " " & LastName` |
| `+` | Also concatenates when left side is string | `"Count: " + Amount` |

The `&` operator always coerces both sides to strings. The `+` operator concatenates
only when the left side is already a string.

---

## Keywords

| Keyword | Aliases | Description |
|---------|---------|-------------|
| `true` | `TRUE` | Boolean true |
| `false` | `FALSE` | Boolean false |
| `nil` | `NULL` | Null value |

All keywords are case-insensitive: `true`, `TRUE`, `True` all work.

---

## Type Coercion

### Numeric Types

`int`, `int64`, `float64` are interchangeable in comparisons and arithmetic:

```
Amount > 100       // works with int, int64, or float64
Amount == 50000.0  // int(50000) == float64(50000.0) is true
Amount + 10        // preserves original type (int+int=int, float+int=float)
```

### String to Number

Strings containing numbers auto-coerce in comparisons:

```
"50000" > 100      // true (auto-coerces string to number)
"abc" > 100        // RUNTIME ERROR — non-numeric string can't compare
```

**Warning:** `string + number` does string concatenation, NOT addition:

```
"50000" + 10       // "5000010" (string concat, NOT 50010)
```

Use `VALUE()` for explicit conversion:

```
VALUE("50000") + 10  // 50010.0 (proper numeric addition)
VALUE(Amount) > 100  // safe even if Amount is a string
```

### Number to String

Use `TEXT()` for explicit conversion:

```
TEXT(50000)          // "50000"
TEXT(3.14)           // "3.14"
TEXT(Amount) & " USD"  // "50000 USD"
```

### Nil Coercion

| Context | Nil behavior |
|---------|-------------|
| Comparisons (`>`, `<`) | nil coerces to 0 |
| String concatenation (`&`) | nil coerces to "" |
| Boolean context (AND, OR, NOT) | nil is falsy |
| Division by zero | Returns nil (not error) |
| Function argument | Most functions return nil (null propagation) |

---

## Function Reference

All function names are **case-insensitive**: `ISBLANK`, `isblank`, `IsBlank` all work.

### Logical Functions

Enabled with `expr.WithLogicalFunctions()`.

#### AND(value1, value2, ...)

Returns `true` if ALL arguments are truthy.

```
AND(Amount > 0, NOT(ISBLANK(Name)))    // true if both conditions met
AND(true, true, false)                  // false
AND()                                   // true (no args = true)
```

#### OR(value1, value2, ...)

Returns `true` if ANY argument is truthy.

```
OR(Status == "Active", Status == "Pending")  // true if either
OR(false, false, true)                        // true
OR()                                          // false (no args = false)
```

#### NOT(value)

Returns the logical negation.

```
NOT(ISBLANK(Email))      // true if Email is not blank
NOT(true)                 // false
NOT(nil)                  // true (nil is falsy)
```

#### IF(condition, trueValue, falseValue)

Returns `trueValue` if condition is truthy, else `falseValue`.

```
IF(Amount > 1000, "High", "Low")           // "High" or "Low"
IF(ISBLANK(Name), "Unknown", Name)          // fallback for blank
IF(TRUE, 1, 2)                              // 1
```

#### CASE(expression, val1, result1, val2, result2, ..., defaultResult)

Compares expression against value/result pairs. Returns matching result or default.

```
CASE(StageName,
    "Prospecting", 1,
    "Qualification", 2,
    "Proposal", 3,
    "Negotiation", 4,
    "Closed Won", 5,
    0                    // default
)
```

#### ISBLANK(value)

Returns `true` if value is nil, empty string `""`, or whitespace-only.

```
ISBLANK(nil)             // true
ISBLANK("")              // true
ISBLANK("  ")            // true
ISBLANK("hello")         // false
ISBLANK(0)               // false (0 is not blank)
```

#### ISNULL(value)

Identical to ISBLANK. Returns `true` if value is nil, empty string, or whitespace-only.

```
ISNULL(nil)              // true
ISNULL("")               // true
```

#### ISNUMBER(value)

Returns `true` if value is numeric or a string parseable as a number.

```
ISNUMBER(42)             // true
ISNUMBER(3.14)           // true
ISNUMBER("100")          // true
ISNUMBER("abc")          // false
ISNUMBER(nil)            // false
```

---

### Math Functions

Enabled with `expr.WithMathFunctions()`.

All math functions return `nil` for nil input (null propagation).

#### ABS(number)

Returns absolute value.

```
ABS(-42)       // 42
ABS(42)        // 42
ABS(nil)       // nil
```

#### CEILING(number)

Returns smallest integer >= number.

```
CEILING(2.1)   // 3
CEILING(-2.9)  // -2
CEILING(nil)   // nil
```

#### FLOOR(number)

Returns largest integer <= number.

```
FLOOR(2.9)     // 2
FLOOR(-2.1)    // -3
FLOOR(nil)     // nil
```

#### ROUND(number, decimalPlaces)

Rounds to given decimal places (half-away-from-zero).

```
ROUND(1.5, 0)      // 2
ROUND(2.345, 2)    // 2.35
ROUND(-1.5, 0)     // -2
ROUND(nil, 2)      // nil
```

#### MAX(a, b)

Returns the larger of two values.

```
MAX(10, 20)    // 20
MAX(nil, 5)    // 5  (nil is ignored)
MAX(nil, nil)  // nil
```

#### MIN(a, b)

Returns the smaller of two values.

```
MIN(10, 20)    // 10
MIN(nil, 5)    // 5  (nil is ignored)
MIN(nil, nil)  // nil
```

#### MOD(number, divisor)

Returns remainder. Division by zero returns nil.

```
MOD(10, 3)     // 1
MOD(10, 0)     // nil
MOD(nil, 3)    // nil
```

#### SQRT(number)

Returns square root. Negative input returns nil.

```
SQRT(16)       // 4
SQRT(-1)       // nil
SQRT(nil)      // nil
```

#### LOG(number)

Returns base-10 logarithm. Non-positive input returns nil.

```
LOG(100)       // 2
LOG(0)         // nil
LOG(nil)       // nil
```

#### LN(number)

Returns natural logarithm. Non-positive input returns nil.

```
LN(2.71828)   // ~1.0
LN(0)          // nil
LN(nil)        // nil
```

#### EXP(number)

Returns e^number.

```
EXP(1)         // 2.718281828...
EXP(0)         // 1
EXP(nil)       // nil
```

#### MCEILING(number)

Ceiling away from zero. Positive: rounds up. Negative: rounds down (away from zero).

```
MCEILING(2.1)   // 3
MCEILING(-2.1)  // -3
```

#### MFLOOR(number)

Floor toward zero. Positive: rounds down. Negative: rounds up (toward zero).

```
MFLOOR(2.9)    // 2
MFLOOR(-2.9)   // -2
```

#### TRUNC(number) / TRUNC(number, places)

Truncates toward zero. Optional decimal places (default 0).

```
TRUNC(2.9)      // 2
TRUNC(-2.9)     // -2
TRUNC(3.456, 2) // 3.45
```

---

### Trigonometric Functions

Enabled with `expr.WithTrigFunctions()`.

All trig functions return `nil` for nil input.

| Function | Description | Example |
|----------|-------------|---------|
| `SIN(radians)` | Sine | `SIN(PI() / 2)` = 1 |
| `COS(radians)` | Cosine | `COS(0)` = 1 |
| `TAN(radians)` | Tangent | `TAN(PI() / 4)` = ~1 |
| `ASIN(number)` | Arcsine (-1 to 1) | `ASIN(1)` = ~1.5708 |
| `ACOS(number)` | Arccosine (-1 to 1) | `ACOS(0)` = ~1.5708 |
| `ATAN(number)` | Arctangent | `ATAN(1)` = ~0.7854 |
| `ATAN2(y, x)` | Arctangent of y/x | `ATAN2(1, 1)` = ~0.7854 |
| `PI()` | Pi constant | `PI()` = 3.14159... |

Out-of-range inputs (e.g., `ASIN(2)`) return nil.

---

### Text Functions

Enabled with `expr.WithTextFunctions()`.

#### BEGINS(text, prefix)

Returns `true` if text starts with prefix.

```
BEGINS("Hello World", "Hello")  // true
BEGINS("Hello", "World")        // false
BEGINS(nil, "test")             // false
```

#### CONTAINS(text, substring)

Returns `true` if text contains substring.

```
CONTAINS("Hello World", "World")  // true
CONTAINS("Hello", "xyz")          // false
CONTAINS(nil, "test")             // false
```

#### FIND(searchString, text) / FIND(searchString, text, startPos)

Returns 1-based position of first occurrence, or 0 if not found.

```
FIND("World", "Hello World")     // 7
FIND("xyz", "Hello World")       // 0
FIND("l", "Hello World", 4)      // 4 (search from position 4)
```

#### LEFT(text, numChars)

Returns leftmost N characters.

```
LEFT("Hello World", 5)   // "Hello"
LEFT("Hi", 10)            // "Hi" (clips to available)
LEFT(nil, 3)              // nil
```

#### RIGHT(text, numChars)

Returns rightmost N characters.

```
RIGHT("Hello World", 5)  // "World"
RIGHT("Hi", 10)           // "Hi"
RIGHT(nil, 3)             // nil
```

#### MID(text, startPos, length)

Returns substring from 1-based start position.

```
MID("Hello World", 7, 5)  // "World"
MID("Hello", 1, 3)        // "Hel"
MID(nil, 1, 3)            // nil
```

#### LEN(text)

Returns character count (Unicode-aware, not byte count).

```
LEN("Hello")     // 5
LEN("")           // 0
LEN(nil)          // 0
```

#### LOWER(text)

Converts to lowercase.

```
LOWER("HELLO")   // "hello"
LOWER(nil)        // nil
```

#### UPPER(text)

Converts to uppercase.

```
UPPER("hello")   // "HELLO"
UPPER(nil)        // nil
```

#### TRIM(text)

Removes leading and trailing whitespace.

```
TRIM("  hello  ")  // "hello"
TRIM(nil)           // nil
```

#### INITCAP(text)

Capitalizes first letter of each word.

```
INITCAP("hello world")    // "Hello World"
INITCAP("HELLO WORLD")    // "Hello World"
INITCAP(nil)               // nil
```

#### SUBSTITUTE(text, oldString, newString)

Replaces all occurrences.

```
SUBSTITUTE("Hello World", "World", "Go")  // "Hello Go"
SUBSTITUTE("aaa", "a", "bb")               // "bbbbbb"
SUBSTITUTE(nil, "a", "b")                  // nil
```

#### LPAD(text, length) / LPAD(text, length, padChar)

Left-pads to target length. Default pad character is space.

```
LPAD("42", 5)           // "   42"
LPAD("42", 5, "0")      // "00042"
LPAD("Hello", 3)        // "Hello" (already >= length)
```

#### RPAD(text, length) / RPAD(text, length, padChar)

Right-pads to target length. Default pad character is space.

```
RPAD("42", 5)           // "42   "
RPAD("42", 5, "0")      // "42000"
```

#### TEXT(value)

Converts any value to string.

```
TEXT(42)               // "42"
TEXT(3.14)             // "3.14"
TEXT(true)             // "true"
TEXT(nil)              // ""
TEXT(50000)            // "50000" (no trailing .0)
```

#### VALUE(text)

Parses string as number. Returns nil if not parseable.

```
VALUE("42")            // 42.0
VALUE("3.14")          // 3.14
VALUE("abc")           // nil
VALUE("")              // nil
VALUE(nil)             // nil
```

**Use VALUE() when your data might be strings but you need numeric operations:**

```
VALUE(Amount) > 100    // safe even if Amount is "50000"
VALUE(Amount) + 10     // numeric addition, not string concat
```

#### REVERSE(text)

Reverses character order (Unicode-aware).

```
REVERSE("Hello")       // "olleH"
REVERSE(nil)           // nil
```

#### BR()

Returns a newline character `"\n"`.

```
FirstName & BR() & LastName  // "John\nDoe"
```

#### ASCII(text)

Returns Unicode code point of first character.

```
ASCII("A")             // 65
ASCII("abc")           // 97
ASCII(nil)             // nil
```

#### CHR(codePoint)

Returns character for Unicode code point.

```
CHR(65)                // "A"
CHR(97)                // "a"
CHR(nil)               // nil
```

---

### Date/Time Functions

Enabled with `expr.WithDateTimeFunctions()`.

#### DATE(year, month, day)

Creates a date from components (UTC).

```
DATE(2026, 3, 15)      // 2026-03-15 00:00:00 UTC
DATE(nil, 1, 1)        // nil
```

#### DATEVALUE(value)

Extracts date portion (sets time to midnight). Parses strings automatically.

```
DATEVALUE("2026-03-15")                  // 2026-03-15 00:00:00 UTC
DATEVALUE("2026-03-15T14:30:00Z")        // 2026-03-15 00:00:00 UTC
DATEVALUE(nil)                            // nil
```

Supported string formats:
- `2006-01-02`
- `2006-01-02T15:04:05Z07:00` (RFC3339)
- `2006-01-02 15:04:05`
- `01/02/2006`
- `Jan 2, 2006`
- `January 2, 2006`
- `02-Jan-2006`

#### DATETIMEVALUE(value)

Converts to datetime (preserves time portion).

```
DATETIMEVALUE("2026-03-15T14:30:00Z")   // 2026-03-15 14:30:00 UTC
```

#### DAY(date)

Returns day of month (1-31).

```
DAY(DATE(2026, 3, 15))   // 15
DAY(nil)                   // nil
```

#### MONTH(date)

Returns month (1-12).

```
MONTH(DATE(2026, 3, 15))  // 3
```

#### YEAR(date)

Returns year.

```
YEAR(DATE(2026, 3, 15))   // 2026
```

#### WEEKDAY(date)

Returns day of week. Salesforce convention: 1=Sunday, 7=Saturday.

```
WEEKDAY(DATE(2026, 3, 15))  // 1 (Sunday)
```

#### ADDMONTHS(date, numMonths)

Adds months to a date, clamping to month-end if needed.

```
ADDMONTHS(DATE(2026, 1, 31), 1)  // 2026-02-28 (clamped)
ADDMONTHS(DATE(2026, 3, 15), -1) // 2026-02-15
ADDMONTHS(nil, 1)                 // nil
```

#### HOUR(datetime) / MINUTE(datetime) / SECOND(datetime) / MILLISECOND(datetime)

Extract time components.

```
HOUR(DATETIMEVALUE("2026-03-15T14:30:45Z"))    // 14
MINUTE(DATETIMEVALUE("2026-03-15T14:30:45Z"))  // 30
SECOND(DATETIMEVALUE("2026-03-15T14:30:45Z"))  // 45
```

#### TIMEVALUE(datetime)

Returns duration since midnight.

```
TIMEVALUE(DATETIMEVALUE("2026-03-15T14:30:00Z"))
// 14h30m0s (as time.Duration)
```

#### DAYOFYEAR(date)

Returns day of year (1-366).

```
DAYOFYEAR(DATE(2026, 2, 1))  // 32
```

#### ISOWEEK(date) / ISOYEAR(date)

Returns ISO 8601 week number and year.

```
ISOWEEK(DATE(2026, 1, 1))   // 1
ISOYEAR(DATE(2026, 1, 1))   // 2026
```

#### FROMUNIXTIME(timestamp)

Converts Unix timestamp to UTC time.

```
FROMUNIXTIME(1700000000)  // 2023-11-14 22:13:20 UTC
```

#### UNIXTIMESTAMP(datetime)

Converts datetime to Unix timestamp.

```
UNIXTIMESTAMP(DATE(2026, 1, 1))  // 1767225600
```

#### FORMATDURATION(startDateTime, endDateTime)

Formats duration between two dates as "DD:HH:MM:SS".

```
FORMATDURATION(
    DATETIMEVALUE("2026-01-01T00:00:00Z"),
    DATETIMEVALUE("2026-01-02T03:30:45Z")
)
// "1:03:30:45"
```

#### TODAY()

Returns today's date at midnight. **Context function** — reads timezone from env.

```
CloseDate < TODAY()       // is close date in the past?
TODAY() + 30              // 30 days from today (date arithmetic)
```

#### NOW()

Returns current datetime. **Context function** — reads timezone from env.

```
CreatedDate > NOW() - 7   // created within last 7 days
```

#### TIMENOW()

Returns duration since midnight today. **Context function**.

```
TIMENOW()                 // e.g., 14h30m0s
```

---

### Regex Functions

Enabled with `expr.WithRegexFunctions()`.

#### REGEX(text, pattern)

Returns `true` if text matches the regex pattern. Patterns are compiled and cached.

```
REGEX(Email, "^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}$")
REGEX(Phone, "^\\+?[0-9]{10,15}$")
REGEX("12345", "^[0-9]+$")      // true
REGEX(nil, ".*")                  // false
```

---

### Encoding Functions

Enabled with `expr.WithEncodingFunctions()`.

| Function | Description | Example |
|----------|-------------|---------|
| `HTMLENCODE(text)` | HTML-escape (`<` -> `&lt;`) | `HTMLENCODE("<b>hi</b>")` -> `"&lt;b&gt;hi&lt;/b&gt;"` |
| `URLENCODE(text)` | URL-encode (spaces -> `%20`) | `URLENCODE("hello world")` -> `"hello%20world"` |
| `JSENCODE(text)` | JavaScript-escape (quotes, newlines) | `JSENCODE("it's \"ok\"")` -> `"it\'s \\\"ok\\\""` |
| `JSINHTMLENCODE(text)` | JS-escape then HTML-escape | Combined encoding |

All return `nil` for nil input.

---

### Salesforce Runtime Functions

Enabled with `expr.WithSalesforceRuntimeFunctions()`.

These functions require runtime context in the env map.

#### ISPICKVAL(field, value)

Checks if a picklist field equals a value. Handles nil fields like Salesforce:
nil field matches empty string `""`.

```
ISPICKVAL(StageName, "Closed Won")  // true if StageName is "Closed Won"
ISPICKVAL(nil, "")                   // true (SF behavior)
ISPICKVAL(nil, "Active")             // false
```

#### INCLUDES(multiSelectField, value)

Checks if a multi-select picklist contains a value. Multi-select values are
semicolon-separated in Salesforce.

```
// If Products = "Widget;Gadget;Tool"
INCLUDES(Products, "Widget")    // true
INCLUDES(Products, "Other")     // false
INCLUDES(nil, "Widget")         // false
```

#### ISCHANGED(fieldName)

Returns `true` if the field was changed in this transaction. Requires `_changed`
map in env.

```
ISCHANGED("Status")    // true if Status was modified
ISCHANGED("OwnerId")  // true if owner changed
```

**Env setup:**
```go
env["_changed"] = map[string]bool{
    "Status":  true,   // was changed
    "OwnerId": false,  // was not changed
}
```

#### PRIORVALUE(fieldName)

Returns the previous value of a field before the update. Requires `_prior`
map in env.

```
PRIORVALUE("Status") == "Open"     // was it previously "Open"?
PRIORVALUE("Amount") > Amount      // did amount decrease?
```

**Env setup:**
```go
env["_prior"] = map[string]any{
    "Status": "Open",
    "Amount": float64(10000),
}
```

#### ISNEW()

Returns `true` if the record is being created (insert, not update). Requires
`_isNew` in env.

```
AND(ISNEW(), ISBLANK(Name))   // name required on new records
NOT(ISNEW())                    // only on updates
```

**Env setup:**
```go
env["_isNew"] = true  // or false for updates
```

---

## ExtractDeps

`ExtractDeps` parses a formula and returns all dependencies without compiling or running it.
Use this to know which fields to fetch from the database.

```go
deps, err := expr.ExtractDeps(formula)
```

### FormulaDeps struct

```go
type FormulaDeps struct {
    RecordFields  []string  // flat fields: ["Name", "Amount", "Status"]
    RelatedFields []string  // dotted paths: ["Account.Name", "Account.Owner.ProfileId"]
    GlobalVars    []string  // $ paths: ["$User.Email", "$Profile.Name"]
    OldFields     []string  // ISCHANGED/PRIORVALUE args: ["Status"]
    NeedsIsNew    bool      // formula uses ISNEW()
    NeedsTimezone bool      // formula uses TODAY()/NOW()/TIMENOW()
}
```

### How it works

- Parses the formula into an AST (no compilation needed)
- Walks all nodes to classify identifiers and member chains
- Function names (AND, ISBLANK, etc.) are NOT included as fields
- Dotted paths (Account.Name) go to RelatedFields
- ISCHANGED("X") / PRIORVALUE("X") string arguments go to OldFields
- All lists are deduplicated and sorted

### Merging deps across multiple rules

```go
allRecord := make(map[string]bool)
allRelated := make(map[string]bool)
allOld := make(map[string]bool)

for _, rule := range rules {
    deps, _ := expr.ExtractDeps(rule.Formula)
    for _, f := range deps.RecordFields  { allRecord[f] = true }
    for _, f := range deps.RelatedFields { allRelated[f] = true }
    for _, f := range deps.OldFields     { allOld[f] = true }
}
// Now fetch the union of all fields in one query
```

---

## InflateEnv

Converts a flat map with dotted keys into nested maps required by the expression engine.

```go
flat := map[string]any{
    "Name":                "Acme Deal",
    "Amount":              float64(50000),
    "Account.Name":        "Acme Corp",
    "Account.Owner.Email": "owner@acme.com",
    "_changed":            map[string]bool{"Status": true},
    "_isNew":              false,
}

env := expr.InflateEnv(flat)
// Result:
// {
//   "Name": "Acme Deal",
//   "Amount": 50000,
//   "Account": {
//     "Name": "Acme Corp",
//     "Owner": {
//       "Email": "owner@acme.com"
//     }
//   },
//   "_changed": {"Status": true},
//   "_isNew": false
// }
```

- Non-dotted keys stay as-is
- Dotted keys become nested maps: `"A.B.C": val` -> `{"A": {"B": {"C": val}}}`
- Underscore-prefixed keys (`_changed`, `_prior`, `_isNew`) stay flat (no dots in them)
- Multiple keys sharing a prefix merge correctly:
  `"Account.Name"` and `"Account.Industry"` both go under the same `"Account"` map

---

## Null Handling

The engine follows Salesforce null propagation semantics:

| Scenario | Result |
|----------|--------|
| `nil + 5` | 5 (nil coerces to 0 in arithmetic) |
| `nil > 0` | false (nil coerces to 0 in comparison) |
| `nil == nil` | true |
| `nil & "text"` | "text" (nil coerces to "" in concat) |
| `10 / 0` | nil (not error, not infinity) |
| `ABS(nil)` | nil (null propagation) |
| `LEN(nil)` | 0 |
| `ISBLANK(nil)` | true |
| `AND(nil, true)` | false (nil is falsy) |
| `OR(nil, true)` | true |
| Undefined variable | nil (with AllowUndefinedVariables) |

---

## Date Arithmetic

Dates support arithmetic with numbers (interpreted as days):

```
CloseDate + 30           // add 30 days
CloseDate - 7            // subtract 7 days
CloseDate - CreatedDate  // difference in days (float64)
TODAY() + 90             // 90 days from now

// Comparison
CloseDate < TODAY()      // is close date in the past?
CloseDate >= DATE(2026, 1, 1)  // on or after Jan 1 2026
```

---

## Common Patterns

### Required field on new records
```
AND(ISNEW(), ISBLANK(FieldName))
```

### Field changed to specific value
```
AND(ISCHANGED("Status"), ISPICKVAL(Status, "Closed"))
```

### Prevent stage regression
```
AND(
    ISCHANGED("StageName"),
    CASE(StageName,
        "Prospecting", 1,
        "Qualification", 2,
        "Proposal", 3,
        "Negotiation", 4,
        "Closed Won", 5,
        0
    ) < CASE(PRIORVALUE("StageName"),
        "Prospecting", 1,
        "Qualification", 2,
        "Proposal", 3,
        "Negotiation", 4,
        "Closed Won", 5,
        0
    )
)
```

### Email format validation
```
AND(NOT(ISBLANK(Email)), NOT(REGEX(Email, "^[^@]+@[^@]+\\.[^@]+$")))
```

### Close date must be in the future
```
AND(NOT(ISNEW()), CloseDate <= TODAY())
```

### Conditional requirement based on picklist
```
AND(
    ISPICKVAL(StageName, "Closed Won"),
    OR(ISBLANK(Amount), ISBLANK(CloseDate))
)
```

### Cross-object validation
```
AND(
    Amount > 100000,
    ISBLANK(Account.Owner.Email)
)
```

### Multi-select picklist check
```
AND(
    INCLUDES(Products, "Premium"),
    Amount < 50000
)
```

### String manipulation
```
AND(
    BEGINS(RIGHT(AccountNumber, 3), "00"),
    LEN(AccountNumber) == 10
)
```

### Numeric string handling
```
VALUE(LEFT(AccountNumber, 3)) > 100
```

---

## Performance

Benchmarked on Apple M1 Max:

| Operation | Time |
|-----------|------|
| Compile one rule | ~29 us |
| ExtractDeps one rule | ~4 us |
| Run one evaluation | ~240 ns |
| 500 records x 5 rules (2,500 evals) | ~1.1 ms |
| 10,000 records x 10 rules (100k evals) | ~24 ms |

**Key insight:** Compilation is a one-time cost. Always compile once and reuse the
program across all records. The per-evaluation cost is ~240 nanoseconds — negligible
compared to any database query.

```go
// DO THIS: compile once, run many
program, _ := expr.Compile(formula, opts...)
for _, record := range records {
    result, _ := expr.Run(program, buildEnv(record))
}

// DON'T DO THIS: compile per record
for _, record := range records {
    program, _ := expr.Compile(formula, opts...)  // wasteful!
    result, _ := expr.Run(program, buildEnv(record))
}
```
