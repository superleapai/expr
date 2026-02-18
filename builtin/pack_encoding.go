package builtin

import (
	"html"
	"net/url"
	"strings"
)

// EncodingFunctions returns the pack functions for encoding/escaping.
func EncodingFunctions() []*PackFunction {
	return []*PackFunction{
		{
			Name:  "HTMLENCODE",
			Types: PackFuncTypes(new(func(...any) any)),
			Fn: func(args ...any) (any, error) {
				if len(args) < 1 {
					return nil, nil
				}
				if args[0] == nil {
					return nil, nil
				}
				text := CoerceToString(args[0])
				return html.EscapeString(text), nil
			},
		},
		{
			Name:  "URLENCODE",
			Types: PackFuncTypes(new(func(...any) any)),
			Fn: func(args ...any) (any, error) {
				if len(args) < 1 {
					return nil, nil
				}
				if args[0] == nil {
					return nil, nil
				}
				text := CoerceToString(args[0])
				return url.QueryEscape(text), nil
			},
		},
		{
			Name:  "JSENCODE",
			Types: PackFuncTypes(new(func(...any) any)),
			Fn: func(args ...any) (any, error) {
				if len(args) < 1 {
					return nil, nil
				}
				if args[0] == nil {
					return nil, nil
				}
				text := CoerceToString(args[0])
				return jsEncode(text), nil
			},
		},
		{
			Name:  "JSINHTMLENCODE",
			Types: PackFuncTypes(new(func(...any) any)),
			Fn: func(args ...any) (any, error) {
				if len(args) < 1 {
					return nil, nil
				}
				if args[0] == nil {
					return nil, nil
				}
				text := CoerceToString(args[0])
				return html.EscapeString(jsEncode(text)), nil
			},
		},
	}
}

// jsEncode performs JavaScript string escaping.
func jsEncode(s string) string {
	r := strings.NewReplacer(
		`\`, `\\`,
		`"`, `\"`,
		`'`, `\'`,
		"\n", `\n`,
		"\r", `\r`,
		"\t", `\t`,
		"/", `\/`,
	)
	return r.Replace(s)
}
