package expr_test

import (
	"testing"

	"github.com/expr-lang/expr"
	"github.com/expr-lang/expr/internal/testify/assert"
	"github.com/expr-lang/expr/internal/testify/require"
)

func TestExtractDeps_RecordFields(t *testing.T) {
	tests := []struct {
		name   string
		expr   string
		fields []string
	}{
		{
			name:   "single field",
			expr:   `Name`,
			fields: []string{"Name"},
		},
		{
			name:   "multiple fields in binary",
			expr:   `Amount > 100`,
			fields: []string{"Amount"},
		},
		{
			name:   "multiple distinct fields",
			expr:   `Name != "" && Amount > 0 && Status == "Active"`,
			fields: []string{"Amount", "Name", "Status"},
		},
		{
			name:   "field in function argument",
			expr:   `LEN(Name) > 10`,
			fields: []string{"Name"},
		},
		{
			name:   "deduplication",
			expr:   `Name == "" || Name == "test"`,
			fields: []string{"Name"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			deps, err := expr.ExtractDeps(tt.expr)
			require.NoError(t, err)
			assert.Equal(t, tt.fields, deps.RecordFields)
		})
	}
}

func TestExtractDeps_RelatedFields(t *testing.T) {
	tests := []struct {
		name    string
		expr    string
		related []string
		record  []string
	}{
		{
			name:    "single dotted path",
			expr:    `Account.Name == "Acme"`,
			related: []string{"Account.Name"},
		},
		{
			name:    "deep dotted path",
			expr:    `Account.Owner.ProfileId == "admin"`,
			related: []string{"Account.Owner.ProfileId"},
		},
		{
			name:    "mixed record and related",
			expr:    `Name != "" && Account.Name == "Acme"`,
			related: []string{"Account.Name"},
			record:  []string{"Name"},
		},
		{
			name:    "multiple related paths",
			expr:    `Account.Name == "Acme" && Contact.Email != ""`,
			related: []string{"Account.Name", "Contact.Email"},
		},
		{
			name:    "deduplication of related",
			expr:    `Account.Name == "" || Account.Name == "test"`,
			related: []string{"Account.Name"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			deps, err := expr.ExtractDeps(tt.expr)
			require.NoError(t, err)
			assert.Equal(t, tt.related, deps.RelatedFields)
			if tt.record != nil {
				assert.Equal(t, tt.record, deps.RecordFields)
			}
		})
	}
}

func TestExtractDeps_OldFields(t *testing.T) {
	tests := []struct {
		name string
		expr string
		old  []string
	}{
		{
			name: "ISCHANGED single field",
			expr: `ISCHANGED("Status")`,
			old:  []string{"Status"},
		},
		{
			name: "PRIORVALUE single field",
			expr: `PRIORVALUE("OwnerId") != OwnerId`,
			old:  []string{"OwnerId"},
		},
		{
			name: "both ISCHANGED and PRIORVALUE",
			expr: `ISCHANGED("Status") && PRIORVALUE("Amount") > 0`,
			old:  []string{"Amount", "Status"},
		},
		{
			name: "case insensitive function names",
			expr: `ischanged("Email")`,
			old:  []string{"Email"},
		},
		{
			name: "deduplication",
			expr: `ISCHANGED("Status") || PRIORVALUE("Status") == "Open"`,
			old:  []string{"Status"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			deps, err := expr.ExtractDeps(tt.expr)
			require.NoError(t, err)
			assert.Equal(t, tt.old, deps.OldFields)
		})
	}
}

func TestExtractDeps_NeedsIsNew(t *testing.T) {
	deps, err := expr.ExtractDeps(`ISNEW() && Status == "Draft"`)
	require.NoError(t, err)
	assert.True(t, deps.NeedsIsNew)
	assert.Equal(t, []string{"Status"}, deps.RecordFields)

	deps, err = expr.ExtractDeps(`Name != ""`)
	require.NoError(t, err)
	assert.False(t, deps.NeedsIsNew)
}

func TestExtractDeps_NeedsTimezone(t *testing.T) {
	tests := []struct {
		name string
		expr string
		want bool
	}{
		{"TODAY", `CloseDate < TODAY()`, true},
		{"NOW", `CreatedDate > NOW()`, true},
		{"TIMENOW", `TimeField == TIMENOW()`, true},
		{"no timezone func", `Amount > 0`, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			deps, err := expr.ExtractDeps(tt.expr)
			require.NoError(t, err)
			assert.Equal(t, tt.want, deps.NeedsTimezone)
		})
	}
}

func TestExtractDeps_FunctionNamesNotFields(t *testing.T) {
	// Function names like AND, OR, IF, ISBLANK should NOT appear as record fields.
	deps, err := expr.ExtractDeps(`AND(Name != "", Amount > 0)`)
	require.NoError(t, err)
	assert.Equal(t, []string{"Amount", "Name"}, deps.RecordFields)

	deps, err = expr.ExtractDeps(`IF(ISBLANK(Email), "none", Email)`)
	require.NoError(t, err)
	assert.Equal(t, []string{"Email"}, deps.RecordFields)

	deps, err = expr.ExtractDeps(`NOT(ISNULL(Phone))`)
	require.NoError(t, err)
	assert.Equal(t, []string{"Phone"}, deps.RecordFields)
}

func TestExtractDeps_ComplexFormulas(t *testing.T) {
	tests := []struct {
		name    string
		expr    string
		record  []string
		related []string
		old     []string
		isNew   bool
		tz      bool
	}{
		{
			name:   "SF validation: ISPICKVAL + ISBLANK",
			expr:   `AND(ISPICKVAL(StageName, "Closed Won"), ISBLANK(CloseDate))`,
			record: []string{"CloseDate", "StageName"},
		},
		{
			name:   "SF validation: ISCHANGED + PRIORVALUE + diamond",
			expr:   `AND(ISCHANGED("Status"), PRIORVALUE("Status") <> "Closed")`,
			record: nil,
			old:    []string{"Status"},
		},
		{
			name:   "SF validation: ISNEW + ISBLANK + TODAY",
			expr:   `ISNEW() && ISBLANK(CloseDate) && CreatedDate < TODAY()`,
			record: []string{"CloseDate", "CreatedDate"},
			isNew:  true,
			tz:     true,
		},
		{
			name:    "SF validation: related fields + record fields",
			expr:    `Account.Industry == "Tech" && Amount > 10000 && Contact.Email != ""`,
			record:  []string{"Amount"},
			related: []string{"Account.Industry", "Contact.Email"},
		},
		{
			name:   "nested IF with CASE",
			expr:   `IF(ISPICKVAL(Status, "New"), CASE(Priority, "High", 1, "Low", 2, 0), 0) > 0`,
			record: []string{"Priority", "Status"},
		},
		{
			name:   "INCLUDES with multi-select",
			expr:   `INCLUDES(Products, "Widget") && Amount > 100`,
			record: []string{"Amount", "Products"},
		},
		{
			name:   "string concat with &",
			expr:   `"Hello " & FirstName & " " & LastName`,
			record: []string{"FirstName", "LastName"},
		},
		{
			name:   "TEXT and VALUE functions",
			expr:   `VALUE(LEFT(AccountNumber, 3)) > 100`,
			record: []string{"AccountNumber"},
		},
		{
			name:   "REGEX with field",
			expr:   `REGEX(Email, "^[a-z]+@example\\.com$")`,
			record: []string{"Email"},
		},
		{
			name:   "ISCHANGED and PRIORVALUE record fields tracked",
			expr:   `ISCHANGED("OwnerId") && OwnerId != PRIORVALUE("OwnerId")`,
			record: []string{"OwnerId"},
			old:    []string{"OwnerId"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			deps, err := expr.ExtractDeps(tt.expr)
			require.NoError(t, err)
			assert.Equal(t, tt.record, deps.RecordFields, "RecordFields")
			if tt.related != nil {
				assert.Equal(t, tt.related, deps.RelatedFields, "RelatedFields")
			}
			if tt.old != nil {
				assert.Equal(t, tt.old, deps.OldFields, "OldFields")
			}
			assert.Equal(t, tt.isNew, deps.NeedsIsNew, "NeedsIsNew")
			assert.Equal(t, tt.tz, deps.NeedsTimezone, "NeedsTimezone")
		})
	}
}

func TestExtractDeps_ErrorCases(t *testing.T) {
	// Parse error
	_, err := expr.ExtractDeps(`AND(Name != ""`)
	assert.Error(t, err)

	// ISCHANGED without string arg
	_, err = expr.ExtractDeps(`ISCHANGED(Status)`)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "ISCHANGED requires a quoted field name")

	// PRIORVALUE without string arg
	_, err = expr.ExtractDeps(`PRIORVALUE(Amount)`)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "PRIORVALUE requires a quoted field name")
}

func TestExtractDeps_ConditionalAndTernary(t *testing.T) {
	deps, err := expr.ExtractDeps(`Status == "Active" ? Amount : DefaultAmount`)
	require.NoError(t, err)
	assert.Equal(t, []string{"Amount", "DefaultAmount", "Status"}, deps.RecordFields)
}

func TestExtractDeps_EmptyResult(t *testing.T) {
	// Expression with only literals
	deps, err := expr.ExtractDeps(`1 + 2 > 3`)
	require.NoError(t, err)
	assert.Nil(t, deps.RecordFields)
	assert.Nil(t, deps.RelatedFields)
	assert.Nil(t, deps.GlobalVars)
	assert.Nil(t, deps.OldFields)
	assert.False(t, deps.NeedsIsNew)
	assert.False(t, deps.NeedsTimezone)
}

func TestExtractDeps_LetExpression(t *testing.T) {
	deps, err := expr.ExtractDeps(`let x = Amount * 2; x > 100`)
	require.NoError(t, err)
	assert.Equal(t, []string{"Amount"}, deps.RecordFields)
}
