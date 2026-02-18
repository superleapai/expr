package builtin

import (
	"fmt"
	"regexp"
	"sync"
)

var regexCache sync.Map

// getCompiledRegex returns a cached compiled regex, or compiles and caches it.
func getCompiledRegex(pattern string) (*regexp.Regexp, error) {
	if cached, ok := regexCache.Load(pattern); ok {
		return cached.(*regexp.Regexp), nil
	}
	re, err := regexp.Compile(pattern)
	if err != nil {
		return nil, err
	}
	regexCache.Store(pattern, re)
	return re, nil
}

// REGEX returns true if the text matches the given regex pattern.
func REGEX(args ...any) (any, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("REGEX expects 2 arguments, got %d", len(args))
	}
	if args[0] == nil || args[1] == nil {
		return false, nil
	}
	text := CoerceToString(args[0])
	pattern := CoerceToString(args[1])
	re, err := getCompiledRegex(pattern)
	if err != nil {
		return false, nil
	}
	return re.MatchString(text), nil
}

// RegexFunctions returns all regex pack functions.
func RegexFunctions() []*PackFunction {
	variadicType := PackFuncTypes(new(func(...any) any))
	return []*PackFunction{
		{Name: "REGEX", Fn: REGEX, Types: variadicType},
	}
}
