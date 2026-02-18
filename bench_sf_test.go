package expr_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/expr-lang/expr"
	"github.com/expr-lang/expr/vm"
)

// buildRecords creates N simulated SF records, each with ~10 fields.
func buildRecords(n int) []map[string]any {
	records := make([]map[string]any, n)
	stages := []string{"Prospecting", "Qualification", "Proposal", "Negotiation", "Closed Won", "Closed Lost"}
	industries := []string{"Technology", "Finance", "Healthcare", "Retail", "Government"}

	for i := 0; i < n; i++ {
		records[i] = map[string]any{
			"Name":             fmt.Sprintf("Deal-%d", i),
			"Amount":           float64(10000 + i*100),
			"StageName":        stages[i%len(stages)],
			"CloseDate":        time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC).AddDate(0, 0, i),
			"Email":            fmt.Sprintf("user%d@example.com", i),
			"OwnerId":          fmt.Sprintf("005%06d", i),
			"Status":           stages[i%3],
			"Account.Name":     fmt.Sprintf("Account-%d", i%50),
			"Account.Industry": industries[i%len(industries)],
			"Contact.Email":    fmt.Sprintf("contact%d@example.com", i),
			"_changed":         map[string]bool{"Status": i%3 == 0},
			"_prior":           map[string]any{"Status": "Prospecting"},
			"_isNew":           i%5 == 0,
		}
	}
	return records
}

// Five realistic validation rules of varying complexity.
var validationRules = []string{
	// Rule 1: simple field check
	`AND(ISPICKVAL(StageName, "Closed Won"), ISBLANK(Amount))`,

	// Rule 2: ISCHANGED + PRIORVALUE
	`AND(ISCHANGED("Status"), PRIORVALUE("Status") == "Prospecting", Amount > 50000)`,

	// Rule 3: related field + text functions
	`AND(NOT(ISBLANK(Email)), NOT(CONTAINS(Email, "@example.com")), Account.Industry <> "Government")`,

	// Rule 4: ISNEW + nested logic
	`IF(ISNEW(), ISBLANK(Amount) || ISBLANK(Email), AND(ISCHANGED("Status"), ISBLANK(Account.Name)))`,

	// Rule 5: CASE + math
	`CASE(StageName, "Prospecting", 1, "Qualification", 2, "Proposal", 3, "Negotiation", 4, "Closed Won", 5, "Closed Lost", 6, 0) < 3 && Amount > 25000`,
}

// TestTiming_500Records is a regular test that prints wall-clock time
// so you can see the result without -bench flags.
func TestTiming_500Records(t *testing.T) {
	const numRecords = 500
	records := buildRecords(numRecords)

	// --- Phase 1: Compile all rules once ---
	compileStart := time.Now()
	programs := make([]*vm.Program, len(validationRules))
	for i, formula := range validationRules {
		// Use a sample env for type info.
		sampleEnv := expr.InflateEnv(records[0])
		p, err := expr.Compile(formula,
			expr.Env(sampleEnv),
			expr.AllowUndefinedVariables(),
			expr.WithAllFormulaPacks(),
		)
		if err != nil {
			t.Fatalf("compile rule %d: %v", i, err)
		}
		programs[i] = p
	}
	compileDur := time.Since(compileStart)

	// --- Phase 2: ExtractDeps for all rules ---
	depsStart := time.Now()
	for _, formula := range validationRules {
		_, err := expr.ExtractDeps(formula)
		if err != nil {
			t.Fatalf("deps: %v", err)
		}
	}
	depsDur := time.Since(depsStart)

	// --- Phase 3: Run all rules against all 500 records ---
	runStart := time.Now()
	totalEvals := 0
	fired := 0
	for _, rec := range records {
		env := expr.InflateEnv(rec)
		for _, program := range programs {
			result, err := expr.Run(program, env)
			if err != nil {
				t.Fatalf("run: %v", err)
			}
			totalEvals++
			if result == true {
				fired++
			}
		}
	}
	runDur := time.Since(runStart)

	t.Logf("=== Timing Results ===")
	t.Logf("Records:        %d", numRecords)
	t.Logf("Rules:          %d", len(validationRules))
	t.Logf("Total evals:    %d  (records x rules)", totalEvals)
	t.Logf("")
	t.Logf("Compile all:    %v  (one-time cost)", compileDur)
	t.Logf("ExtractDeps:    %v  (one-time cost)", depsDur)
	t.Logf("Run all:        %v", runDur)
	t.Logf("Per record:     %v  (%d rules)", runDur/time.Duration(numRecords), len(validationRules))
	t.Logf("Per eval:       %v", runDur/time.Duration(totalEvals))
	t.Logf("Validations fired: %d / %d", fired, totalEvals)
}

// BenchmarkSF_CompileRule benchmarks compiling a single rule.
func BenchmarkSF_CompileRule(b *testing.B) {
	sampleEnv := expr.InflateEnv(buildRecords(1)[0])
	formula := validationRules[1] // medium complexity

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := expr.Compile(formula,
			expr.Env(sampleEnv),
			expr.AllowUndefinedVariables(),
			expr.WithAllFormulaPacks(),
		)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkSF_ExtractDeps benchmarks dependency extraction.
func BenchmarkSF_ExtractDeps(b *testing.B) {
	formula := validationRules[1]
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := expr.ExtractDeps(formula)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkSF_RunSingleEval benchmarks running one pre-compiled rule against one record.
func BenchmarkSF_RunSingleEval(b *testing.B) {
	rec := buildRecords(1)[0]
	env := expr.InflateEnv(rec)
	program, _ := expr.Compile(validationRules[1],
		expr.Env(env),
		expr.AllowUndefinedVariables(),
		expr.WithAllFormulaPacks(),
	)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := expr.Run(program, env)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkSF_500Records_5Rules benchmarks the full 500-record x 5-rule workload.
func BenchmarkSF_500Records_5Rules(b *testing.B) {
	records := buildRecords(500)
	envs := make([]map[string]any, len(records))
	for i, rec := range records {
		envs[i] = expr.InflateEnv(rec)
	}

	programs := make([]*vm.Program, len(validationRules))
	for i, formula := range validationRules {
		p, _ := expr.Compile(formula,
			expr.Env(envs[0]),
			expr.AllowUndefinedVariables(),
			expr.WithAllFormulaPacks(),
		)
		programs[i] = p
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, env := range envs {
			for _, program := range programs {
				_, _ = expr.Run(program, env)
			}
		}
	}
}
