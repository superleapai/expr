package builtin

import (
	"strconv"
	"strings"
	"unicode/utf8"
)

// TextFunctions returns the pack functions for text manipulation.
func TextFunctions() []*PackFunction {
	return []*PackFunction{
		{
			Name:  "BEGINS",
			Types: PackFuncTypes(new(func(...any) any)),
			Fn: func(args ...any) (any, error) {
				if len(args) < 2 {
					return false, nil
				}
				if args[0] == nil || args[1] == nil {
					return false, nil
				}
				text := CoerceToString(args[0])
				prefix := CoerceToString(args[1])
				return strings.HasPrefix(text, prefix), nil
			},
		},
		{
			Name:  "CONTAINS",
			Types: PackFuncTypes(new(func(...any) any)),
			Fn: func(args ...any) (any, error) {
				if len(args) < 2 {
					return false, nil
				}
				if args[0] == nil || args[1] == nil {
					return false, nil
				}
				text := CoerceToString(args[0])
				substr := CoerceToString(args[1])
				return strings.Contains(text, substr), nil
			},
		},
		{
			Name:  "FIND",
			Types: PackFuncTypes(new(func(...any) any)),
			Fn: func(args ...any) (any, error) {
				if len(args) < 2 {
					return 0, nil
				}
				if args[0] == nil || args[1] == nil {
					return 0, nil
				}
				search := CoerceToString(args[0])
				text := CoerceToString(args[1])

				start := 1
				if len(args) >= 3 && args[2] != nil {
					f, wasNil := CoerceToFloat64(args[2])
					if !wasNil {
						start = int(f)
					}
				}

				runes := []rune(text)
				searchRunes := []rune(search)

				// Convert 1-based start to 0-based
				startIdx := start - 1
				if startIdx < 0 {
					startIdx = 0
				}
				if startIdx >= len(runes) {
					return 0, nil
				}

				// Search within the substring starting at startIdx
				sub := string(runes[startIdx:])
				byteIdx := strings.Index(sub, string(searchRunes))
				if byteIdx < 0 {
					return 0, nil
				}

				// Count runes up to the found position and add back the start offset
				runePos := utf8.RuneCountInString(sub[:byteIdx])
				return runePos + startIdx + 1, nil // Return 1-based index
			},
		},
		{
			Name:  "LEFT",
			Types: PackFuncTypes(new(func(...any) any)),
			Fn: func(args ...any) (any, error) {
				if len(args) < 2 {
					return nil, nil
				}
				if args[0] == nil {
					return nil, nil
				}
				text := CoerceToString(args[0])
				f, _ := CoerceToFloat64(args[1])
				n := int(f)
				if n <= 0 {
					return "", nil
				}
				runes := []rune(text)
				if n > len(runes) {
					return text, nil
				}
				return string(runes[:n]), nil
			},
		},
		{
			Name:  "RIGHT",
			Types: PackFuncTypes(new(func(...any) any)),
			Fn: func(args ...any) (any, error) {
				if len(args) < 2 {
					return nil, nil
				}
				if args[0] == nil {
					return nil, nil
				}
				text := CoerceToString(args[0])
				f, _ := CoerceToFloat64(args[1])
				n := int(f)
				if n <= 0 {
					return "", nil
				}
				runes := []rune(text)
				if n > len(runes) {
					return text, nil
				}
				return string(runes[len(runes)-n:]), nil
			},
		},
		{
			Name:  "MID",
			Types: PackFuncTypes(new(func(...any) any)),
			Fn: func(args ...any) (any, error) {
				if len(args) < 3 {
					return nil, nil
				}
				if args[0] == nil {
					return nil, nil
				}
				text := CoerceToString(args[0])
				runes := []rune(text)

				sf, _ := CoerceToFloat64(args[1])
				start := int(sf) - 1 // Convert 1-based to 0-based
				if start < 0 {
					start = 0
				}

				lf, _ := CoerceToFloat64(args[2])
				length := int(lf)
				if length <= 0 {
					return "", nil
				}

				if start >= len(runes) {
					return "", nil
				}

				end := start + length
				if end > len(runes) {
					end = len(runes)
				}
				return string(runes[start:end]), nil
			},
		},
		{
			Name:  "LEN",
			Types: PackFuncTypes(new(func(...any) any)),
			Fn: func(args ...any) (any, error) {
				if len(args) < 1 {
					return 0, nil
				}
				// LEN(nil) → 0
				s := CoerceToString(args[0])
				return utf8.RuneCountInString(s), nil
			},
		},
		{
			Name:  "LOWER",
			Types: PackFuncTypes(new(func(...any) any)),
			Fn: func(args ...any) (any, error) {
				if len(args) < 1 {
					return nil, nil
				}
				if args[0] == nil {
					return nil, nil
				}
				text := CoerceToString(args[0])
				return strings.ToLower(text), nil
			},
		},
		{
			Name:  "UPPER",
			Types: PackFuncTypes(new(func(...any) any)),
			Fn: func(args ...any) (any, error) {
				if len(args) < 1 {
					return nil, nil
				}
				if args[0] == nil {
					return nil, nil
				}
				text := CoerceToString(args[0])
				return strings.ToUpper(text), nil
			},
		},
		{
			Name:  "TRIM",
			Types: PackFuncTypes(new(func(...any) any)),
			Fn: func(args ...any) (any, error) {
				if len(args) < 1 {
					return nil, nil
				}
				if args[0] == nil {
					return nil, nil
				}
				text := CoerceToString(args[0])
				return strings.TrimSpace(text), nil
			},
		},
		{
			Name:  "INITCAP",
			Types: PackFuncTypes(new(func(...any) any)),
			Fn: func(args ...any) (any, error) {
				if len(args) < 1 {
					return nil, nil
				}
				if args[0] == nil {
					return nil, nil
				}
				text := CoerceToString(args[0])
				words := strings.Fields(text)
				for i, w := range words {
					if len(w) == 0 {
						continue
					}
					runes := []rune(w)
					runes[0] = []rune(strings.ToUpper(string(runes[0])))[0]
					for j := 1; j < len(runes); j++ {
						runes[j] = []rune(strings.ToLower(string(runes[j])))[0]
					}
					words[i] = string(runes)
				}
				return strings.Join(words, " "), nil
			},
		},
		{
			Name:  "SUBSTITUTE",
			Types: PackFuncTypes(new(func(...any) any)),
			Fn: func(args ...any) (any, error) {
				if len(args) < 3 {
					return nil, nil
				}
				if args[0] == nil || args[1] == nil || args[2] == nil {
					return nil, nil
				}
				text := CoerceToString(args[0])
				old := CoerceToString(args[1])
				newStr := CoerceToString(args[2])
				if old == "" {
					return text, nil
				}
				return strings.ReplaceAll(text, old, newStr), nil
			},
		},
		{
			Name:  "LPAD",
			Types: PackFuncTypes(new(func(...any) any)),
			Fn: func(args ...any) (any, error) {
				if len(args) < 2 {
					return nil, nil
				}
				if args[0] == nil {
					return nil, nil
				}
				text := CoerceToString(args[0])
				f, _ := CoerceToFloat64(args[1])
				totalLen := int(f)

				padChar := ' '
				if len(args) >= 3 && args[2] != nil {
					ps := CoerceToString(args[2])
					if len(ps) > 0 {
						padChar = []rune(ps)[0]
					}
				}

				runes := []rune(text)
				if len(runes) >= totalLen {
					return text, nil
				}
				padding := strings.Repeat(string(padChar), totalLen-len(runes))
				return padding + text, nil
			},
		},
		{
			Name:  "RPAD",
			Types: PackFuncTypes(new(func(...any) any)),
			Fn: func(args ...any) (any, error) {
				if len(args) < 2 {
					return nil, nil
				}
				if args[0] == nil {
					return nil, nil
				}
				text := CoerceToString(args[0])
				f, _ := CoerceToFloat64(args[1])
				totalLen := int(f)

				padChar := ' '
				if len(args) >= 3 && args[2] != nil {
					ps := CoerceToString(args[2])
					if len(ps) > 0 {
						padChar = []rune(ps)[0]
					}
				}

				runes := []rune(text)
				if len(runes) >= totalLen {
					return text, nil
				}
				padding := strings.Repeat(string(padChar), totalLen-len(runes))
				return text + padding, nil
			},
		},
		{
			Name:  "TEXT",
			Types: PackFuncTypes(new(func(...any) any)),
			Fn: func(args ...any) (any, error) {
				if len(args) < 1 {
					return "", nil
				}
				// TEXT(nil) → "" (NOT nil)
				return CoerceToString(args[0]), nil
			},
		},
		{
			Name:  "VALUE",
			Types: PackFuncTypes(new(func(...any) any)),
			Fn: func(args ...any) (any, error) {
				if len(args) < 1 {
					return nil, nil
				}
				if args[0] == nil {
					return nil, nil
				}
				s := CoerceToString(args[0])
				if s == "" {
					return nil, nil
				}
				f, err := strconv.ParseFloat(s, 64)
				if err != nil {
					return nil, nil
				}
				return f, nil
			},
		},
		{
			Name:  "REVERSE",
			Types: PackFuncTypes(new(func(...any) any)),
			Fn: func(args ...any) (any, error) {
				if len(args) < 1 {
					return nil, nil
				}
				if args[0] == nil {
					return nil, nil
				}
				text := CoerceToString(args[0])
				runes := []rune(text)
				for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
					runes[i], runes[j] = runes[j], runes[i]
				}
				return string(runes), nil
			},
		},
		{
			Name:  "BR",
			Types: PackFuncTypes(new(func() string)),
			Fn: func(args ...any) (any, error) {
				return "\n", nil
			},
		},
		{
			Name:  "ASCII",
			Types: PackFuncTypes(new(func(...any) any)),
			Fn: func(args ...any) (any, error) {
				if len(args) < 1 {
					return nil, nil
				}
				if args[0] == nil {
					return nil, nil
				}
				text := CoerceToString(args[0])
				if text == "" {
					return nil, nil
				}
				runes := []rune(text)
				return int(runes[0]), nil
			},
		},
		{
			Name:  "CHR",
			Types: PackFuncTypes(new(func(...any) any)),
			Fn: func(args ...any) (any, error) {
				if len(args) < 1 {
					return nil, nil
				}
				if args[0] == nil {
					return nil, nil
				}
				f, wasNil := CoerceToFloat64(args[0])
				if wasNil {
					return nil, nil
				}
				code := int(f)
				if code == 0 {
					return nil, nil
				}
				return string(rune(code)), nil
			},
		},
	}
}
