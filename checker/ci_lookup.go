package checker

import (
	"strings"

	"github.com/expr-lang/expr/builtin"
)

// resolveFunctionCI performs a case-insensitive lookup in a FunctionsTable.
// It returns nil if no match is found.
func resolveFunctionCI(name string, table map[string]*builtin.Function) *builtin.Function {
	upper := strings.ToUpper(name)
	for k, fn := range table {
		if strings.ToUpper(k) == upper {
			return fn
		}
	}
	return nil
}
