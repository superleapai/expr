package expr

import (
	"fmt"
	"sort"
	"strings"

	"github.com/expr-lang/expr/ast"
	"github.com/expr-lang/expr/parser"
)

// FormulaDeps holds the dependencies extracted from a parsed formula expression.
type FormulaDeps struct {
	RecordFields  []string // flat fields: ["Name", "Amount", "Status"]
	RelatedFields []string // dotted paths: ["Account.Name", "Account.Owner.ProfileId"]
	GlobalVars    []string // $ paths: ["$User.Email", "$Profile.Name"]
	OldFields     []string // from ISCHANGED/PRIORVALUE string args: ["Status"]
	NeedsIsNew    bool     // formula calls ISNEW()
	NeedsTimezone bool     // formula calls TODAY()/NOW()/TIMENOW()
}

// ExtractDeps parses a formula string and extracts all field references,
// related-object paths, global variable references, and special function usage.
func ExtractDeps(formula string) (*FormulaDeps, error) {
	tree, err := parser.Parse(formula)
	if err != nil {
		return nil, err
	}

	w := &depsWalker{
		recordFields:  make(map[string]bool),
		relatedFields: make(map[string]bool),
		globalVars:    make(map[string]bool),
		oldFields:     make(map[string]bool),
		localVars:     make(map[string]bool),
	}

	if err := w.walk(tree.Node); err != nil {
		return nil, err
	}

	return &FormulaDeps{
		RecordFields:  sortedKeys(w.recordFields),
		RelatedFields: sortedKeys(w.relatedFields),
		GlobalVars:    sortedKeys(w.globalVars),
		OldFields:     sortedKeys(w.oldFields),
		NeedsIsNew:    w.needsIsNew,
		NeedsTimezone: w.needsTimezone,
	}, nil
}

type depsWalker struct {
	recordFields  map[string]bool
	relatedFields map[string]bool
	globalVars    map[string]bool
	oldFields     map[string]bool
	localVars     map[string]bool // variables declared via let
	needsIsNew    bool
	needsTimezone bool
}

// walk recursively visits AST nodes. isCallee indicates the node is
// the callee position of a CallNode (i.e. a function name, not a field).
func (w *depsWalker) walk(node ast.Node) error {
	return w.walkNode(node, false)
}

func (w *depsWalker) walkNode(node ast.Node, isCallee bool) error {
	if node == nil {
		return nil
	}
	switch n := node.(type) {

	case *ast.IdentifierNode:
		if isCallee {
			// Function name — check for special functions, but don't add as field.
			w.checkFuncName(n.Value)
			return nil
		}
		if w.localVars[n.Value] {
			// Variable declared by let — not a field reference.
		} else if strings.HasPrefix(n.Value, "$") {
			w.globalVars[n.Value] = true
		} else {
			w.recordFields[n.Value] = true
		}

	case *ast.MemberNode:
		path, ok := resolveDotPath(n)
		if !ok {
			// Dynamic property (e.g. array[idx]) — walk children normally.
			if err := w.walkNode(n.Node, false); err != nil {
				return err
			}
			return nil
		}
		if isCallee {
			// Method-style call on an object — the root object is still a dep.
			// Strip the last segment (method name) and treat the rest as a field path.
			if idx := strings.LastIndex(path, "."); idx > 0 {
				path = path[:idx]
			}
		}
		w.addPath(path)
		return nil // children already consumed by resolveDotPath

	case *ast.CallNode:
		// Walk callee as a function name.
		if err := w.walkNode(n.Callee, true); err != nil {
			return err
		}
		// Check for ISCHANGED / PRIORVALUE special handling.
		name := calleeName(n.Callee)
		upper := strings.ToUpper(name)
		if upper == "ISCHANGED" || upper == "PRIORVALUE" {
			if len(n.Arguments) < 1 {
				return fmt.Errorf("%s requires at least 1 argument", upper)
			}
			strNode, ok := n.Arguments[0].(*ast.StringNode)
			if !ok {
				return fmt.Errorf("%s requires a quoted field name as first argument", upper)
			}
			w.oldFields[strNode.Value] = true
		}
		// Walk arguments.
		for _, arg := range n.Arguments {
			if err := w.walkNode(arg, false); err != nil {
				return err
			}
		}

	case *ast.BuiltinNode:
		// BuiltinNodes are created by the checker, not parser.Parse,
		// but handle them for safety.
		w.checkFuncName(n.Name)
		for _, arg := range n.Arguments {
			if err := w.walkNode(arg, false); err != nil {
				return err
			}
		}

	case *ast.UnaryNode:
		return w.walkNode(n.Node, false)

	case *ast.BinaryNode:
		if err := w.walkNode(n.Left, false); err != nil {
			return err
		}
		return w.walkNode(n.Right, false)

	case *ast.ConditionalNode:
		if err := w.walkNode(n.Cond, false); err != nil {
			return err
		}
		if err := w.walkNode(n.Exp1, false); err != nil {
			return err
		}
		return w.walkNode(n.Exp2, false)

	case *ast.ChainNode:
		return w.walkNode(n.Node, false)

	case *ast.SliceNode:
		if err := w.walkNode(n.Node, false); err != nil {
			return err
		}
		if err := w.walkNode(n.From, false); err != nil {
			return err
		}
		return w.walkNode(n.To, false)

	case *ast.ArrayNode:
		for _, child := range n.Nodes {
			if err := w.walkNode(child, false); err != nil {
				return err
			}
		}

	case *ast.MapNode:
		for _, pair := range n.Pairs {
			if err := w.walkNode(pair, false); err != nil {
				return err
			}
		}

	case *ast.PairNode:
		if err := w.walkNode(n.Key, false); err != nil {
			return err
		}
		return w.walkNode(n.Value, false)

	case *ast.SequenceNode:
		for _, child := range n.Nodes {
			if err := w.walkNode(child, false); err != nil {
				return err
			}
		}

	case *ast.VariableDeclaratorNode:
		w.localVars[n.Name] = true
		if err := w.walkNode(n.Value, false); err != nil {
			return err
		}
		return w.walkNode(n.Expr, false)

	case *ast.PredicateNode:
		return w.walkNode(n.Node, false)

	case *ast.NilNode, *ast.IntegerNode, *ast.FloatNode,
		*ast.BoolNode, *ast.StringNode, *ast.ConstantNode, *ast.PointerNode:
		// Leaf nodes — no deps.
	}
	return nil
}

// checkFuncName records special function usage flags.
func (w *depsWalker) checkFuncName(name string) {
	switch strings.ToUpper(name) {
	case "ISNEW":
		w.needsIsNew = true
	case "TODAY", "NOW", "TIMENOW":
		w.needsTimezone = true
	}
}

// addPath classifies a resolved dotted path into the right bucket.
func (w *depsWalker) addPath(path string) {
	switch {
	case strings.HasPrefix(path, "$"):
		w.globalVars[path] = true
	case strings.Contains(path, "."):
		w.relatedFields[path] = true
	default:
		w.recordFields[path] = true
	}
}

// resolveDotPath walks a MemberNode chain and returns the full dotted path.
// Returns ("", false) if any property is not a StringNode (dynamic access).
func resolveDotPath(n *ast.MemberNode) (string, bool) {
	prop, ok := n.Property.(*ast.StringNode)
	if !ok {
		return "", false
	}
	switch inner := n.Node.(type) {
	case *ast.IdentifierNode:
		return inner.Value + "." + prop.Value, true
	case *ast.MemberNode:
		parent, ok := resolveDotPath(inner)
		if !ok {
			return "", false
		}
		return parent + "." + prop.Value, true
	default:
		return "", false
	}
}

// calleeName extracts the function name from a CallNode's Callee.
func calleeName(node ast.Node) string {
	if id, ok := node.(*ast.IdentifierNode); ok {
		return id.Value
	}
	return ""
}

// InflateEnv converts a flat map with dotted keys into a nested map structure
// suitable for expr evaluation. For example:
//
//	{"Account.Name": "Acme", "Account.Owner.Email": "a@b.com", "Amount": 100}
//
// becomes:
//
//	{"Account": {"Name": "Acme", "Owner": {"Email": "a@b.com"}}, "Amount": 100}
//
// Non-dotted keys are kept as-is. Keys starting with "_" (like "_changed",
// "_prior", "_isNew") are kept at the top level for runtime context.
func InflateEnv(flat map[string]any) map[string]any {
	root := make(map[string]any, len(flat))
	for key, val := range flat {
		if !strings.Contains(key, ".") {
			root[key] = val
			continue
		}
		parts := strings.Split(key, ".")
		m := root
		for _, p := range parts[:len(parts)-1] {
			if existing, ok := m[p]; ok {
				if sub, ok := existing.(map[string]any); ok {
					m = sub
				} else {
					// Conflict: a leaf value already exists at this path segment.
					// Overwrite with a nested map (last writer wins).
					sub := make(map[string]any)
					m[p] = sub
					m = sub
				}
			} else {
				sub := make(map[string]any)
				m[p] = sub
				m = sub
			}
		}
		m[parts[len(parts)-1]] = val
	}
	return root
}

// sortedKeys returns a sorted slice of map keys, or nil if the map is empty.
func sortedKeys(m map[string]bool) []string {
	if len(m) == 0 {
		return nil
	}
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	sort.Strings(out)
	return out
}
