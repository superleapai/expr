package expr_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/expr-lang/expr"
	"github.com/expr-lang/expr/internal/testify/assert"
	"github.com/expr-lang/expr/internal/testify/require"
	"github.com/expr-lang/expr/vm"
)

// ---------------------------------------------------------------------------
// This file demonstrates the full end-to-end flow for Salesforce-style
// validation rules using the expr formula engine:
//
//   1. Store the formula string (adapted from SF syntax)
//   2. ExtractDeps  → know which fields/relations/old-values are needed
//   3. Fetch field values from your DB (simulated here)
//   4. InflateEnv   → convert flat "Account.Name" keys to nested maps
//   5. Inject runtime context (_changed, _prior, _isNew, timezone)
//   6. Compile + Run → get bool result (true = validation error triggered)
// ---------------------------------------------------------------------------

// simulateDBFetch pretends to query the database for the given field lists.
// In a real app this would be a SOQL-like query built from the deps.
func simulateDBFetch(recordFields, relatedFields []string) map[string]any {
	// Simulated record data — flat keys, including dotted relation paths.
	db := map[string]any{
		"Name":                "Acme Deal",
		"Amount":              float64(50000),
		"StageName":           "Proposal",
		"CloseDate":           time.Date(2026, 3, 15, 0, 0, 0, 0, time.UTC),
		"OwnerId":             "005ABC",
		"Status":              "Negotiation",
		"Email":               "sales@acme.com",
		"Account.Name":        "Acme Corp",
		"Account.Industry":    "Technology",
		"Account.Owner.Email": "owner@acme.com",
		"Contact.Email":       "john@acme.com",
	}

	result := make(map[string]any)
	for _, f := range recordFields {
		if v, ok := db[f]; ok {
			result[f] = v
		}
		// Fields not in DB stay absent (nil in expr)
	}
	for _, f := range relatedFields {
		if v, ok := db[f]; ok {
			result[f] = v
		}
	}
	return result
}

// simulateOldValues fetches prior values for ISCHANGED/PRIORVALUE fields.
func simulateOldValues(oldFields []string) (changed map[string]bool, prior map[string]any) {
	// Simulated: Status was changed from "Qualification" to "Negotiation"
	oldDB := map[string]any{
		"Status":  "Qualification",
		"OwnerId": "005ABC", // not changed
	}
	changed = make(map[string]bool)
	prior = make(map[string]any)
	currentDB := map[string]any{
		"Status":  "Negotiation",
		"OwnerId": "005ABC",
	}
	for _, f := range oldFields {
		oldVal := oldDB[f]
		curVal := currentDB[f]
		prior[f] = oldVal
		changed[f] = oldVal != curVal
	}
	return
}

// ---------------------------------------------------------------------------
// Full end-to-end test: one formula, full pipeline
// ---------------------------------------------------------------------------

func TestE2E_SF_Validation_FullPipeline(t *testing.T) {
	// --- Step 1: The formula (adapted from SF syntax) ---
	// "If stage changed AND the prior stage was Qualification AND close date
	//  is in the past, fire a validation error."
	formula := `AND(
		ISCHANGED("Status"),
		PRIORVALUE("Status") == "Qualification",
		CloseDate < TODAY()
	)`

	// --- Step 2: Extract deps ---
	deps, err := expr.ExtractDeps(formula)
	require.NoError(t, err)

	t.Logf("RecordFields:  %v", deps.RecordFields)
	t.Logf("RelatedFields: %v", deps.RelatedFields)
	t.Logf("OldFields:     %v", deps.OldFields)
	t.Logf("NeedsIsNew:    %v", deps.NeedsIsNew)
	t.Logf("NeedsTimezone: %v", deps.NeedsTimezone)

	assert.Equal(t, []string{"CloseDate"}, deps.RecordFields)
	assert.Equal(t, []string{"Status"}, deps.OldFields)
	assert.True(t, deps.NeedsTimezone) // TODAY() is used

	// --- Step 3: Fetch values from DB ---
	flat := simulateDBFetch(deps.RecordFields, deps.RelatedFields)

	// --- Step 4: Fetch old values for ISCHANGED/PRIORVALUE ---
	changed, prior := simulateOldValues(deps.OldFields)
	flat["_changed"] = changed
	flat["_prior"] = prior

	// --- Step 5: ISNEW context (this is an update, not a new record) ---
	if deps.NeedsIsNew {
		flat["_isNew"] = false
	}

	// --- Step 6: InflateEnv for dotted paths ---
	env := expr.InflateEnv(flat)

	// --- Step 7: Compile ---
	program, err := expr.Compile(formula,
		expr.Env(env),
		expr.AllowUndefinedVariables(),
		expr.WithAllFormulaPacks(),
	)
	require.NoError(t, err)

	// --- Step 8: Run ---
	result, err := expr.Run(program, env)
	require.NoError(t, err)

	t.Logf("Validation result: %v", result)

	// Status was changed (true), prior was "Qualification" (true),
	// CloseDate is 2026-03-15 which is in the future relative to TODAY.
	// So CloseDate < TODAY() is false → AND is false.
	assert.Equal(t, false, result)
}

// ---------------------------------------------------------------------------
// Multiple validation rules evaluated in sequence
// ---------------------------------------------------------------------------

func TestE2E_SF_MultipleRules(t *testing.T) {
	rules := []struct {
		name    string
		formula string
		wantErr bool // true = validation fires
	}{
		{
			name:    "Amount required for Closed Won",
			formula: `AND(ISPICKVAL(StageName, "Closed Won"), ISBLANK(Amount))`,
			wantErr: false, // StageName is "Proposal", not "Closed Won"
		},
		{
			name:    "Email must not be blank",
			formula: `ISBLANK(Email)`,
			wantErr: false, // Email is set
		},
		{
			name:    "Account name required",
			formula: `ISBLANK(Account.Name)`,
			wantErr: false, // Account.Name is "Acme Corp"
		},
		{
			name:    "Amount cannot exceed 100k for non-Enterprise",
			formula: `AND(Amount > 100000, Account.Industry <> "Enterprise")`,
			wantErr: false, // Amount is 50k
		},
		{
			name:    "Status must change for update",
			formula: `AND(NOT(ISNEW()), NOT(ISCHANGED("Status")))`,
			wantErr: false, // Status IS changed
		},
	}

	// Collect all deps across all rules first.
	allRecord := make(map[string]bool)
	allRelated := make(map[string]bool)
	allOld := make(map[string]bool)
	needsIsNew := false
	needsTZ := false

	for _, r := range rules {
		deps, err := expr.ExtractDeps(r.formula)
		require.NoError(t, err, r.name)
		for _, f := range deps.RecordFields {
			allRecord[f] = true
		}
		for _, f := range deps.RelatedFields {
			allRelated[f] = true
		}
		for _, f := range deps.OldFields {
			allOld[f] = true
		}
		if deps.NeedsIsNew {
			needsIsNew = true
		}
		if deps.NeedsTimezone {
			needsTZ = true
		}
	}

	t.Logf("All record fields: %v", allRecord)
	t.Logf("All related fields: %v", allRelated)
	t.Logf("All old fields:     %v", allOld)
	t.Logf("Needs ISNEW: %v, Needs TZ: %v", needsIsNew, needsTZ)

	// Build env once for all rules.
	recordList := make([]string, 0, len(allRecord))
	for k := range allRecord {
		recordList = append(recordList, k)
	}
	relatedList := make([]string, 0, len(allRelated))
	for k := range allRelated {
		relatedList = append(relatedList, k)
	}
	oldList := make([]string, 0, len(allOld))
	for k := range allOld {
		oldList = append(oldList, k)
	}

	flat := simulateDBFetch(recordList, relatedList)

	if len(oldList) > 0 {
		changed, prior := simulateOldValues(oldList)
		flat["_changed"] = changed
		flat["_prior"] = prior
	}
	if needsIsNew {
		flat["_isNew"] = false // it's an update
	}

	env := expr.InflateEnv(flat)

	// Run each rule.
	for _, r := range rules {
		t.Run(r.name, func(t *testing.T) {
			program, err := expr.Compile(r.formula,
				expr.Env(env),
				expr.AllowUndefinedVariables(),
				expr.WithAllFormulaPacks(),
			)
			require.NoError(t, err)

			result, err := expr.Run(program, env)
			require.NoError(t, err)

			fired := result == true
			assert.Equal(t, r.wantErr, fired,
				"rule %q: expected fired=%v got %v (result=%v)", r.name, r.wantErr, fired, result)
		})
	}
}

// ---------------------------------------------------------------------------
// Test InflateEnv directly
// ---------------------------------------------------------------------------

func TestInflateEnv(t *testing.T) {
	flat := map[string]any{
		"Name":                "Acme Deal",
		"Amount":              50000,
		"Account.Name":        "Acme Corp",
		"Account.Owner.Email": "owner@acme.com",
		"Contact.Email":       "john@acme.com",
		"_changed":            map[string]bool{"Status": true},
		"_isNew":              false,
	}

	env := expr.InflateEnv(flat)

	// Flat fields stay flat.
	assert.Equal(t, "Acme Deal", env["Name"])
	assert.Equal(t, 50000, env["Amount"])

	// Dotted paths become nested maps.
	account, ok := env["Account"].(map[string]any)
	require.True(t, ok, "Account should be a nested map")
	assert.Equal(t, "Acme Corp", account["Name"])

	owner, ok := account["Owner"].(map[string]any)
	require.True(t, ok, "Account.Owner should be a nested map")
	assert.Equal(t, "owner@acme.com", owner["Email"])

	contact, ok := env["Contact"].(map[string]any)
	require.True(t, ok, "Contact should be a nested map")
	assert.Equal(t, "john@acme.com", contact["Email"])

	// Underscore runtime keys stay flat.
	assert.Equal(t, map[string]bool{"Status": true}, env["_changed"])
	assert.Equal(t, false, env["_isNew"])
}

// ---------------------------------------------------------------------------
// Test InflateEnv + expr.Run for dotted field access
// ---------------------------------------------------------------------------

func TestInflateEnv_WithExprRun(t *testing.T) {
	flat := map[string]any{
		"Amount":              float64(75000),
		"StageName":           "Closed Won",
		"Account.Name":        "BigCo",
		"Account.Industry":    "Finance",
		"Account.Owner.Email": "boss@bigco.com",
	}
	env := expr.InflateEnv(flat)

	tests := []struct {
		name string
		expr string
		want any
	}{
		{"flat field", `Amount > 50000`, true},
		{"dotted 1 level", `Account.Name == "BigCo"`, true},
		{"dotted 2 levels", `Account.Owner.Email == "boss@bigco.com"`, true},
		{"dotted in function", `CONTAINS(Account.Owner.Email, "bigco")`, true},
		{"diamond on dotted", `Account.Industry <> "Technology"`, true},
		{"concat dotted", `Account.Name & " - " & Account.Industry`, "BigCo - Finance"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			program, err := expr.Compile(tt.expr,
				expr.Env(env),
				expr.AllowUndefinedVariables(),
				expr.WithAllFormulaPacks(),
			)
			require.NoError(t, err)
			result, err := expr.Run(program, env)
			require.NoError(t, err)
			assert.Equal(t, tt.want, result)
		})
	}
}

// ---------------------------------------------------------------------------
// Demonstrate how you'd use this in production code (pseudocode in comments)
// ---------------------------------------------------------------------------

// ---------------------------------------------------------------------------
// Multiple validation rules x multiple records (production pattern)
// ---------------------------------------------------------------------------

func TestE2E_SF_MultipleRules_MultipleRecords(t *testing.T) {
	// ---------------------------------------------------------------
	// Step 1: Define validation rules (loaded from DB in production)
	// ---------------------------------------------------------------
	type ValidationRule struct {
		Name       string
		Formula    string
		ErrMessage string
	}

	rules := []ValidationRule{
		{
			Name:       "Amount required for Closed Won",
			Formula:    `AND(ISPICKVAL(StageName, "Closed Won"), ISBLANK(Amount))`,
			ErrMessage: "Amount is required when Stage is Closed Won",
		},
		{
			Name:       "Email format check",
			Formula:    `AND(NOT(ISBLANK(Email)), NOT(CONTAINS(Email, "@")))`,
			ErrMessage: "Email must contain @",
		},
		{
			Name:       "High amount needs manager",
			Formula:    `AND(Amount > 100000, ISBLANK(Account.Owner.Email))`,
			ErrMessage: "Deals over 100k require an Account Owner with email",
		},
		{
			Name:       "Stage regression blocked",
			Formula:    `AND(ISCHANGED("StageName"), CASE(StageName, "Prospecting", 1, "Qualification", 2, "Proposal", 3, "Negotiation", 4, "Closed Won", 5, 0) < CASE(PRIORVALUE("StageName"), "Prospecting", 1, "Qualification", 2, "Proposal", 3, "Negotiation", 4, "Closed Won", 5, 0))`,
			ErrMessage: "Cannot move stage backwards",
		},
		{
			Name:       "New record must have name",
			Formula:    `AND(ISNEW(), ISBLANK(Name))`,
			ErrMessage: "Name is required for new records",
		},
	}

	// ---------------------------------------------------------------
	// Step 2: Extract deps from ALL rules (union of all deps)
	// ---------------------------------------------------------------
	allRecord := make(map[string]bool)
	allRelated := make(map[string]bool)
	allOld := make(map[string]bool)
	needsIsNew := false
	needsTZ := false

	for _, r := range rules {
		deps, err := expr.ExtractDeps(r.Formula)
		require.NoError(t, err, r.Name)
		for _, f := range deps.RecordFields {
			allRecord[f] = true
		}
		for _, f := range deps.RelatedFields {
			allRelated[f] = true
		}
		for _, f := range deps.OldFields {
			allOld[f] = true
		}
		if deps.NeedsIsNew {
			needsIsNew = true
		}
		if deps.NeedsTimezone {
			needsTZ = true
		}
	}

	t.Logf("Union deps — record: %v, related: %v, old: %v, isNew: %v, tz: %v",
		allRecord, allRelated, allOld, needsIsNew, needsTZ)

	// ---------------------------------------------------------------
	// Step 3: Pre-compile all rules once (reused across all records)
	// ---------------------------------------------------------------
	// We need a sample env for type-checking at compile time.
	sampleFlat := map[string]any{
		"Name":                "",
		"Amount":              float64(0),
		"StageName":           "",
		"Email":               "",
		"Account.Owner.Email": "",
		"_changed":            map[string]bool{},
		"_prior":              map[string]any{},
		"_isNew":              false,
	}
	sampleEnv := expr.InflateEnv(sampleFlat)

	type CompiledRule struct {
		ValidationRule
		Program *vm.Program
	}
	compiled := make([]CompiledRule, len(rules))
	for i, r := range rules {
		p, err := expr.Compile(r.Formula,
			expr.Env(sampleEnv),
			expr.AllowUndefinedVariables(),
			expr.WithAllFormulaPacks(),
		)
		require.NoError(t, err, "compile %s", r.Name)
		compiled[i] = CompiledRule{ValidationRule: r, Program: p}
	}

	// ---------------------------------------------------------------
	// Step 4: Simulate multiple records (from DB in production)
	// ---------------------------------------------------------------
	type Record struct {
		Label  string
		Fields map[string]any // flat current values
		Old    map[string]any // flat prior values (nil if new)
		IsNew  bool
	}

	records := []Record{
		{
			Label: "Existing deal, stage advanced",
			Fields: map[string]any{
				"Name":                "Big Deal",
				"Amount":              float64(200000),
				"StageName":           "Negotiation",
				"Email":               "rep@acme.com",
				"Account.Owner.Email": "boss@acme.com",
			},
			Old:   map[string]any{"StageName": "Proposal"},
			IsNew: false,
		},
		{
			Label: "New record with blank name",
			Fields: map[string]any{
				"Name":                nil,
				"Amount":              float64(5000),
				"StageName":           "Prospecting",
				"Email":               "new@example.com",
				"Account.Owner.Email": "mgr@example.com",
			},
			Old:   nil,
			IsNew: true,
		},
		{
			Label: "Closed Won missing amount",
			Fields: map[string]any{
				"Name":                "Won Deal",
				"Amount":              nil,
				"StageName":           "Closed Won",
				"Email":               "closer@acme.com",
				"Account.Owner.Email": "boss@acme.com",
			},
			Old:   map[string]any{"StageName": "Negotiation"},
			IsNew: false,
		},
		{
			Label: "Stage regression Negotiation→Qualification",
			Fields: map[string]any{
				"Name":                "Regressed Deal",
				"Amount":              float64(30000),
				"StageName":           "Qualification",
				"Email":               "rep@acme.com",
				"Account.Owner.Email": "boss@acme.com",
			},
			Old:   map[string]any{"StageName": "Negotiation"},
			IsNew: false,
		},
		{
			Label: "Bad email format, no @ sign",
			Fields: map[string]any{
				"Name":                "Email Deal",
				"Amount":              float64(20000),
				"StageName":           "Proposal",
				"Email":               "not-an-email",
				"Account.Owner.Email": "owner@acme.com",
			},
			Old:   map[string]any{"StageName": "Proposal"},
			IsNew: false,
		},
		{
			Label: "High amount, no account owner email",
			Fields: map[string]any{
				"Name":                "Whale Deal",
				"Amount":              float64(250000),
				"StageName":           "Negotiation",
				"Email":               "whale@bigcorp.com",
				"Account.Owner.Email": nil,
			},
			Old:   map[string]any{"StageName": "Proposal"},
			IsNew: false,
		},
		{
			Label: "Multiple violations: new, no name, bad email",
			Fields: map[string]any{
				"Name":                nil,
				"Amount":              float64(8000),
				"StageName":           "Prospecting",
				"Email":               "bademail",
				"Account.Owner.Email": "mgr@example.com",
			},
			Old:   nil,
			IsNew: true,
		},
		{
			Label: "Clean record, no errors",
			Fields: map[string]any{
				"Name":                "Good Deal",
				"Amount":              float64(50000),
				"StageName":           "Proposal",
				"Email":               "sales@acme.com",
				"Account.Owner.Email": "owner@acme.com",
			},
			Old:   map[string]any{"StageName": "Proposal"},
			IsNew: false,
		},
	}

	// ---------------------------------------------------------------
	// Step 5: Run all rules against all records, collect errors
	// ---------------------------------------------------------------
	type ValidationError struct {
		RecordLabel string
		RuleName    string
		Message     string
	}

	var errors []ValidationError

	for _, rec := range records {
		// Build the env for this record.
		flat := make(map[string]any)
		for k, v := range rec.Fields {
			flat[k] = v
		}

		// Build _changed and _prior from old values.
		if rec.Old != nil {
			changed := make(map[string]bool)
			prior := make(map[string]any)
			for _, f := range []string{"StageName"} { // only track fields in allOld
				oldVal, hasOld := rec.Old[f]
				curVal := rec.Fields[f]
				prior[f] = oldVal
				changed[f] = hasOld && oldVal != curVal
			}
			flat["_changed"] = changed
			flat["_prior"] = prior
		} else {
			flat["_changed"] = map[string]bool{}
			flat["_prior"] = map[string]any{}
		}

		if needsIsNew {
			flat["_isNew"] = rec.IsNew
		}

		env := expr.InflateEnv(flat)

		// Evaluate every rule.
		for _, cr := range compiled {
			result, err := expr.Run(cr.Program, env)
			require.NoError(t, err, "%s / %s", rec.Label, cr.Name)

			if result == true {
				errors = append(errors, ValidationError{
					RecordLabel: rec.Label,
					RuleName:    cr.Name,
					Message:     cr.ErrMessage,
				})
			}
		}
	}

	// ---------------------------------------------------------------
	// Step 6: Assert expected validation errors
	// ---------------------------------------------------------------
	t.Logf("Validation errors found: %d", len(errors))
	for _, e := range errors {
		t.Logf("  [%s] %s: %s", e.RecordLabel, e.RuleName, e.Message)
	}

	// Build a lookup for easier assertion: "RecordLabel|RuleName" → true
	errorSet := make(map[string]bool)
	for _, e := range errors {
		errorSet[e.RecordLabel+"|"+e.RuleName] = true
	}

	// Expect exactly 7 errors across 8 records:
	require.Len(t, errors, 7)

	// Record "New record with blank name" → fires: "New record must have name"
	assert.True(t, errorSet["New record with blank name|New record must have name"])

	// Record "Closed Won missing amount" → fires: "Amount required for Closed Won"
	assert.True(t, errorSet["Closed Won missing amount|Amount required for Closed Won"])

	// Record "Stage regression" → fires: "Stage regression blocked"
	assert.True(t, errorSet["Stage regression Negotiation→Qualification|Stage regression blocked"])

	// Record "Bad email format" → fires: "Email format check"
	assert.True(t, errorSet["Bad email format, no @ sign|Email format check"])

	// Record "High amount, no account owner email" → fires: "High amount needs manager"
	assert.True(t, errorSet["High amount, no account owner email|High amount needs manager"])

	// Record "Multiple violations" → fires BOTH: "New record must have name" AND "Email format check"
	assert.True(t, errorSet["Multiple violations: new, no name, bad email|New record must have name"])
	assert.True(t, errorSet["Multiple violations: new, no name, bad email|Email format check"])

	// Record "Existing deal, stage advanced" → no errors (stage went forward)
	assert.False(t, errorSet["Existing deal, stage advanced|Stage regression blocked"])

	// Record "Clean record, no errors" → no errors at all
	for _, e := range errors {
		assert.NotEqual(t, "Clean record, no errors", e.RecordLabel)
	}
}

func ExampleExtractDeps() {
	formula := `AND(
		ISCHANGED("Status"),
		PRIORVALUE("Status") == "Open",
		NOT(ISBLANK(Amount)),
		Account.Industry <> "Government"
	)`

	deps, err := expr.ExtractDeps(formula)
	if err != nil {
		panic(err)
	}

	fmt.Println("RecordFields:", deps.RecordFields)
	fmt.Println("RelatedFields:", deps.RelatedFields)
	fmt.Println("OldFields:", deps.OldFields)
	fmt.Println("NeedsIsNew:", deps.NeedsIsNew)
	fmt.Println("NeedsTimezone:", deps.NeedsTimezone)

	// Output:
	// RecordFields: [Amount]
	// RelatedFields: [Account.Industry]
	// OldFields: [Status]
	// NeedsIsNew: false
	// NeedsTimezone: false
}

func ExampleInflateEnv() {
	flat := map[string]any{
		"Amount":           50000,
		"Account.Name":     "Acme",
		"Account.Industry": "Tech",
	}
	env := expr.InflateEnv(flat)

	program, err := expr.Compile(
		`Account.Name == "Acme" && Amount > 10000`,
		expr.Env(env),
		expr.AllowUndefinedVariables(),
		expr.WithAllFormulaPacks(),
	)
	if err != nil {
		panic(err)
	}
	result, _ := expr.Run(program, env)
	fmt.Println(result)

	// Output:
	// true
}
