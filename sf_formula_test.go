package expr_test

import (
	"testing"
	"time"

	"github.com/expr-lang/expr"
	"github.com/expr-lang/expr/internal/testify/assert"
	"github.com/expr-lang/expr/internal/testify/require"
)

// sfOpts returns standard options for SF formula tests.
func sfOpts() []expr.Option {
	return []expr.Option{
		expr.AllowUndefinedVariables(),
		expr.AsAny(),
		expr.WithAllFormulaPacks(),
	}
}

// sfEnv returns a standard test environment with common SF fields.
func sfEnv() map[string]any {
	return map[string]any{
		// System variables (flattened from $User, $Profile, etc.)
		"Profile_Name":               "Direct Sales Rep",
		"User_Email":                 "test.user@razorpay.com",
		"User_Id":                    "005OX00000CdHtr",
		"User_Username":              "test.user@razorpay.com",
		"UserRole_RollupDescription": "Direct Sales",
		"Permission_LeadValidation":  false,

		// Record fields
		"Account_ParentId":         "001ABC",
		"OwnerId":                  "005XYZ",
		"Owner_Email":              "owner@razorpay.com",
		"Owner_FirstName":          "John",
		"Owner_Role":               "Direct Sales",
		"RecordType_Name":          "Sales",
		"RecordType_DeveloperName": "Ezetap",
		"RecordTypeId":             "012C50000004IXRIA2",

		// Custom fields
		"AM_manager":                  true,
		"Paid_AM_Plan":                "Premium",
		"AM_Prepaid_Amount":           1000.0,
		"Expiry_Date":                 time.Date(2024, 12, 31, 0, 0, 0, 0, time.UTC),
		"BU_Owner":                    "Sales",
		"StageName":                   "Working",
		"Status":                      "Open",
		"Type":                        "Payment Gateway",
		"Ranking":                     "5abc",
		"Screenshot_link":             "https://drive.google.com/file/123",
		"Lifetime_GMV":                5000000.0,
		"From_Amount":                 100.0,
		"To_Amount":                   200.0,
		"CloseDate":                   time.Date(2024, 6, 15, 0, 0, 0, 0, time.UTC),
		"CreatedDate":                 time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		"IsConverted":                 false,
		"IsWon":                       false,
		"FirstName":                   "Jane",
		"LastName":                    "Doe",
		"Description":                 "Test description",
		"Methods":                     "Credit Card;Debit Card",
		"Product":                     "Payment Gateway;Opfin",
		"Merchant_ID":                 "MID123",
		"No_of_Employees":             50.0,
		"Existing_Payroll_System":     "ADP",
		"Transacting":                 true,
		"Upsell":                      false,
		"Competitor_Name":             "Stripe",
		"Meeting_Notes":               "Discussed pricing",
		"Approval_Status":             "Pending",
		"Inactive_Record":             false,
		"Dropped_Record":              false,
		"From_LWC":                    false,
		"Campaign_Name":               "ICICI Base_Calling Activity",
		"Is_Prime_Opportunity":        true,
		"BANT_Qualification_Rec":      "BANT-001",
		"Referral_Details":            "",
		"Other_Reason":                "Custom reason",
		"What_is_required":            "https://razorpay.slack.com/channel",
		"Developer_Sign_off":          false,
		"Reject_Remarks":              "",
		"Expected_GMV":                "1234.56",
		"Subject":                     "Pricing Review",
		"LastActivityDate":            time.Date(2024, 3, 1, 0, 0, 0, 0, time.UTC),
		"Remarks":                     "Some remarks",
		"No_of_director":              "3",
		"Directors_With_Existing_DSC": 2.0,

		// Multi-select
		"Modules": "Payroll;Other",

		// Picklist fields (string values)
		"Settlement_Schedule":       "Other",
		"Designation_of_stake":      "Others",
		"Engagement":                "Follow Up Engagement",
		"User_Status":               "Active",
		"Stages":                    "New",
		"Request_Category":          "Integration Error / Debugging",
		"Backup_Terminal":           "Yes",
		"Backup_Terminal_Id":        "TERM001",
		"Product_Type":              "Payment Gateway",
		"Revivable":                 "Active",
		"Loss_Reason":               "Price",
		"Confidence_Level":          "High",
		"Business_Unit":             "Payments",
		"Ops_UW_Status":             "NEW",
		"LeadSource":                "Website",
		"Sub_Stage":                 "Connected",
		"Ticket_Status":             "New",
		"Lead_Status":               "Open",
		"CA_Pitched_Partner_Bank":   "ICICI",
		"Variant_Sold":              "Standard",
		"Payment_Collected":         "Yes",
		"Amount_Collected":          5000.0,
		"Corporate_Card_Interested": "Yes",
		"CA_Interested":             "Yes",
		"Company_Type":              "Private Limited Company (Pvt. Ltd.)",
		"Subscription":              "Monthly",
		"Service_Provider":          "Simpl",

		// Relationship fields (flattened)
		"Account_Owner_Role":                   "Direct Sales",
		"LastModifiedBy_FirstName":             "Admin",
		"LastModifiedBy_Profile_Name":          "System Administrator",
		"Primary_Merchant_POC":                 "POC001",
		"Primary_Merchant_POC_Owner_Role_Name": "Sales",
		"Account_Name":                         "Acme Corp",
		"Link_Previous_Engagement":             "ENG001",
		"Price_Book_IsActive":                  true,
		"Product_IsActive":                     true,

		// Context data for ISCHANGED/PRIORVALUE/ISNEW
		"_changed": map[string]bool{
			"OwnerId":    true,
			"StageName":  true,
			"Owner_Role": true,
			"AM_manager": true,
			"BU_Owner":   true,
			"Status":     false,
		},
		"_prior": map[string]any{
			"OwnerId":    "005OLD",
			"StageName":  "Open",
			"Owner_Role": "Banking",
			"Status":     "New",
		},
		"_isNew": false,
		"_tz":    time.UTC,
	}
}

// =============================================================================
// SF Formula Pattern Tests
// Each test adapts a real SF validation rule formula to our expr syntax.
// Key adaptations:
//   SF `=` (comparison) → `==`
//   SF `$User.Email` → `User_Email`
//   SF `$Profile.Name` → `Profile_Name`
//   SF `Owner.FirstName` → `Owner_FirstName`
//   SF `RecordType.Name` → `RecordType_Name`
//   SF `Account__r.ParentId` → `Account_ParentId`
//   SF `field__c` → `field` (stripped suffix for readability)
//   SF `null` → `NULL`
//   SF single-line /* comment */ → removed
// =============================================================================

func TestSFFormula_Compile_And_Run(t *testing.T) {
	env := sfEnv()
	opts := append(sfOpts(), expr.Env(env))

	tests := []struct {
		name string
		expr string
		want any
	}{
		// Formula 1: Simple NOT(ISBLANK(...))
		{
			name: "F1_NOT_ISBLANK",
			expr: `NOT(ISBLANK(Account_ParentId))`,
			want: true,
		},
		// Formula 3: IF with OR of field comparisons
		{
			name: "F3_IF_OR_field_comparison",
			expr: `IF(
				User_Email == "pearl.sahni@razorpay.com" || User_Email == "pratik.pujara@razorpay.com",
				FALSE,
				TRUE
			)`,
			want: true,
		},
		// Formula 5: AND with ISCHANGED, PRIORVALUE, <>
		{
			name: "F5_AND_ISCHANGED_PRIORVALUE_diamond",
			expr: `AND(
				ISCHANGED("OwnerId"),
				PRIORVALUE("Owner_Role") == "Banking",
				Owner_Role <> "Banking"
			)`,
			want: true, // OwnerId changed, prior was Banking, current is Direct Sales
		},
		// Formula 6: AND with <> and ==
		{
			name: "F6_AND_diamond_eq",
			expr: `AND(Owner_Email <> "salesforce@razorpay.com", Profile_Name == "Direct Sales Rep")`,
			want: true,
		},
		// Formula 12: REGEX + NOT + ISBLANK + !=
		{
			name: "F12_REGEX_NOT_ISBLANK",
			expr: `!REGEX(Ranking, "(^[0-9])") && NOT(Permission_LeadValidation) && !ISBLANK(Ranking) && (Profile_Name != "Integration User" || Profile_Name != "System Administrator")`,
			want: false, // Ranking="5abc" starts with digit, REGEX returns true, !true = false → whole expr is false
		},
		// Formula 15: AND with <> and PRIORVALUE on OwnerId
		{
			name: "F15_AND_PRIORVALUE_diamond",
			expr: `AND(
				OwnerId <> PRIORVALUE("OwnerId"),
				PRIORVALUE("OwnerId") <> "005OX00000eAK3QYAW",
				User_Id == "005OX00000CdHtr"
			)`,
			want: true,
		},
		// Formula 19: IF with AND/OR, ISPICKVAL, ISBLANK
		{
			name: "F19_IF_AND_OR_ISPICKVAL_ISBLANK",
			expr: `IF(AND(OR(RecordType_Name == "Sales", RecordType_Name == "Service"), ISPICKVAL(Settlement_Schedule, "Other"), ISBLANK(NULL)), true, false)`,
			want: true,
		},
		// Formula 20: Simple comparison
		{
			name: "F20_simple_comparison_gte",
			expr: `From_Amount >= To_Amount`,
			want: false, // 100 >= 200 = false
		},
		// Formula 27: NOT/ISBLANK/BEGINS/RIGHT
		{
			name: "F27_BEGINS_RIGHT",
			expr: `NOT(ISBLANK(Screenshot_link)) && (
				NOT(BEGINS(Screenshot_link, "https://drive.google.com/")) &&
				NOT(BEGINS(Screenshot_link, "drive.google.com/")) &&
				NOT(BEGINS(Screenshot_link, "http://drive.google.com/")) || RIGHT(Screenshot_link, 4) == "com/")`,
			want: false, // Screenshot_link starts with https://drive.google.com/
		},
		// Formula 34: AND with ISPICKVAL + ISBLANK
		{
			name: "F34_ISPICKVAL_ISBLANK",
			expr: `AND(ISPICKVAL(Backup_Terminal, "Yes"), ISBLANK(NULL))`,
			want: true, // Backup_Terminal="Yes", NULL is blank
		},
		// Formula 36: AND with ISPICKVAL and <> True
		{
			name: "F36_ISPICKVAL_diamond_True",
			expr: `AND(ISPICKVAL(User_Status, "Inactive"), Inactive_Record <> TRUE)`,
			want: false, // User_Status is "Active", not "Inactive"
		},
		// Formula 42: CASE with < comparison (stage progression check)
		{
			name: "F42_CASE_stage_progression",
			expr: `ISCHANGED("StageName") &&
			CASE(StageName,
				"New", 1,
				"Under-Discussion", 2,
				"Requirement Discussion Closed", 3,
				"Integration In Progress", 4,
				"Testing & QA", 5,
				"Integration Audit", 6,
				"Integration Completed", 7,
				"MX Go-Live", 8,
				"Hypercare Period", 9,
				"Integration Closed", 10,
				"Request Rejected", 11,
				0
			)
			<
			CASE(PRIORVALUE("StageName"),
				"New", 1,
				"Under-Discussion", 2,
				"Requirement Discussion Closed", 3,
				"Integration In Progress", 4,
				"Testing & QA", 5,
				"Integration Audit", 6,
				"Integration Completed", 7,
				"MX Go-Live", 8,
				"Hypercare Period", 9,
				"Integration Closed", 10,
				"Request Rejected", 11,
				0
			)`,
			want: false, // StageName="Working"→0, Prior="Open"→0, 0 < 0 = false
		},
		// Formula 47: AND with != , INCLUDES, ISBLANK
		{
			name: "F47_AND_INCLUDES_ISBLANK",
			expr: `AND(Profile_Name != "Integration User", Profile_Name != "Direct Sales Rep", IsConverted, NOT(INCLUDES(Product, "Opfin")), OR(ISBLANK(Product), ISBLANK(Merchant_ID)))`,
			want: false, // Profile_Name is "Direct Sales Rep" so != fails
		},
		// Formula 48: Ispickval (case-insensitive) + ISBLANK
		{
			name: "F48_Ispickval_caseinsensitive",
			expr: `Ispickval(Loss_Reason, "Other") && ISBLANK(Other_Reason) && NOT(RecordType_DeveloperName == "Ezetap")`,
			want: false, // Loss_Reason is "Price", not "Other"
		},
		// Formula 50: IF with AND, OR, ISNEW, ISCHANGED, ISBLANK
		{
			name: "F50_IF_AND_OR_ISNEW_ISCHANGED",
			expr: `IF(AND(OR(User_Id == "005OX00000CdHtr", User_Id == "005C5000000K3s1"), NOT(ISNEW()), ISCHANGED("OwnerId"), ISBLANK(NULL)), true, false)`,
			want: true, // User_Id matches, not new, OwnerId changed, NULL is blank
		},
		// Formula 52: INCLUDES + ISNULL
		{
			name: "F52_INCLUDES_ISNULL",
			expr: `INCLUDES(Product, "Partnership") && ISNULL(Methods) && NOT(RecordType_DeveloperName == "Ezetap")`,
			want: false, // Product doesn't include "Partnership" as semicolon-separated
		},
		// Formula 54: ISPICKVAL with multiple booleans
		{
			name: "F54_ISPICKVAL_booleans",
			expr: `AND(ISPICKVAL(Status, "Qualified"), true, true, true, true, true)`,
			want: false, // Status is "Open", not "Qualified"
		},
		// Formula 59: IF with INCLUDES returning FALSE/TRUE
		{
			name: "F59_IF_RecordType_INCLUDES",
			expr: `IF(
				RecordType_DeveloperName == "Rize_Leads",
				IF(INCLUDES(Product, "Incorporation"), FALSE, TRUE),
				FALSE
			)`,
			want: false, // RecordType_DeveloperName is "Ezetap", not "Rize_Leads"
		},
		// Formula 61: OR with VALUE, TEXT, comparisons
		{
			name: "F61_OR_VALUE_TEXT",
			expr: `OR(Directors_With_Existing_DSC < 0, Directors_With_Existing_DSC > VALUE(No_of_director), VALUE(No_of_director) > 5)`,
			want: false, // 2 < 0 no, 2 > 3 no, 3 > 5 no
		},
		// Formula 62: VALUE, LEFT, FIND complex
		{
			name: "F62_VALUE_LEFT_FIND",
			expr: `OR(VALUE(LEFT(Expected_GMV, FIND(".", Expected_GMV) - 1)) > 9999, VALUE(Expected_GMV) > 9999)`,
			want: false, // LEFT("1234.56", 4) = "1234", VALUE("1234") = 1234, 1234 > 9999 no; 1234.56 > 9999 no
		},
		// Formula 66: AND ISPICKVAL empty string
		{
			name: "F66_ISPICKVAL_empty",
			expr: `AND(ISPICKVAL(Engagement, ""), NOT(From_LWC))`,
			want: false, // Engagement is "Follow Up Engagement", not ""
		},
		// Formula 67: field == null check
		{
			name: "F67_field_null_check",
			expr: `AND(
				ISPICKVAL(Engagement, "Follow Up Engagement"),
				Link_Previous_Engagement == NULL,
				NOT(From_LWC)
			)`,
			want: false, // Link_Previous_Engagement is "ENG001", not null
		},
		// Formula 69: field == null
		{
			name: "F69_field_null",
			expr: `AND(Account_Name == NULL, NOT(From_LWC))`,
			want: false, // Account_Name is "Acme Corp"
		},
		// Formula 76: Date comparison with DATEVALUE
		{
			name: "F76_date_comparison",
			expr: `Is_Prime_Opportunity == TRUE && (CloseDate < DATEVALUE(CreatedDate))`,
			want: false, // CloseDate (2024-06-15) < CreatedDate (2024-01-01) is false
		},
		// Formula 89: Complex multi-line with ISCHANGED, ISPICKVAL, ISNEW, !=
		{
			name: "F89_complex_multiline",
			expr: `(UserRole_RollupDescription != "")
				&&
				((ISCHANGED("Type") || ISCHANGED("CloseDate"))
				&&
				ISPICKVAL(StageName, "Closed Won"))
				&&
				Transacting == true
				&&
				!ISNEW()
				&&
				(LastModifiedBy_Profile_Name != "Integration User" || LastModifiedBy_Profile_Name != "System Administrator")
				&&
				RecordType_DeveloperName == "Ezetap"`,
			want: false, // StageName is "Working", not "Closed Won"
		},
		// Formula 95: CloseDate >= Today() with ISNEW
		{
			name: "F95_CloseDate_Today_isnew",
			expr: `AND(CloseDate >= TODAY(), UserRole_RollupDescription <> "Automation User", ISNEW(), TRUE)`,
			want: false, // ISNEW() is false
		},
		// Formula 99: Nested IF
		{
			name: "F99_nested_IF",
			expr: `IF(AND(
				Account_Owner_Role != "Automation User",
				OR(User_Id == "005OX00000CdHtr", User_Id == "005C5000000K3s1"),
				ISCHANGED("OwnerId"))
			, true,
			IF(AND(
				NOT(ISPICKVAL(Business_Unit, "Cross Border")), Account_Owner_Role == "Automation User", OR(User_Id == "005OX00000CdHtr", User_Id == "005C5000000K3s1")),
			 true, false))`,
			want: true, // First IF: Account_Owner_Role != "Automation User"=true, User_Id matches, OwnerId changed → true
		},
		// Formula 107: AND with ISPICKVAL, NOT, field!=
		{
			name: "F107_AND_ISPICKVAL_not_field",
			expr: `AND(ISPICKVAL(StageName, "Closed Won"), NOT(ISPICKVAL(Loss_Reason, "")), NOT(ISPICKVAL(Loss_Reason, "")), RecordTypeId == "012C50000004IXRIA2", NOT(Competitor_Name == ""), NOT(Existing_Payroll_System == ""), RecordType_DeveloperName == "Ezetap")`,
			want: false, // StageName is "Working"
		},
		// Formula 118: IF with CASE < CASE (stage regression)
		{
			name: "F118_IF_CASE_stage_regression",
			expr: `IF(Profile_Name != "Integration User" && Profile_Name != "System Administrator",
			CASE(StageName,
				"Open", 1,
				"Working", 2,
				"Qualify", 3,
				"Agreement Signed", 4,
				"Onboarding", 5,
				"Closed Won", 6,
				"Closed Lost", 7,
				0)
			<
			CASE(PRIORVALUE("StageName"),
				"Open", 1,
				"Working", 2,
				"Qualify", 3,
				"Agreement Signed", 4,
				"Onboarding", 5,
				"Closed Won", 6,
				"Closed Lost", 7,
				0)
			, false)`,
			want: false, // StageName="Working"→2, Prior="Open"→1, 2 < 1 = false
		},
		// Formula 122: IF ISPICKVAL with ISCHANGED
		{
			name: "F122_IF_ISPICKVAL_ISCHANGED",
			expr: `IF(ISPICKVAL(StageName, "Closed Lost"), ISCHANGED("OwnerId"), false)`,
			want: false, // StageName is "Working", not "Closed Lost"
		},
		// Formula 123: ISNEW() == FALSE
		{
			name: "F123_ISNEW_FALSE",
			expr: `ISNEW() == FALSE
				&&
				UserRole_RollupDescription != "Automation User"
				&&
				RecordType_DeveloperName == "Ezetap"
				&&
				ISBLANK(NULL)`,
			want: true, // ISNEW()=false==FALSE is true, rest matches
		},
		// Formula 150: AND with ISPICKVAL || ISPICKVAL, nested OR with VALUE
		{
			name: "F150_ISPICKVAL_OR_VALUE",
			expr: `AND(
				ISPICKVAL(Company_Type, "Private Limited Company (Pvt. Ltd.)")
				|| ISPICKVAL(Company_Type, "Limited Liability Partnership (LLP)"),
				OR(
					Directors_With_Existing_DSC < 0,
					Directors_With_Existing_DSC > VALUE(No_of_director),
					VALUE(No_of_director) > 5
				)
			)`,
			want: false, // Company_Type matches, but OR conditions all false
		},
		// Formula 155: AND with NOT and ISNEW
		{
			name: "F155_NOT_ISNEW",
			expr: `AND(NOT(Profile_Name == "System Administrator"), NOT(ISNEW()))`,
			want: true, // Profile is not SysAdmin, not new
		},
		// Formula 158: IF with TEXT(field) == and ISBLANK
		{
			name: "F158_IF_TEXT_ISBLANK",
			expr: `IF(AND(Status == "Reject", ISBLANK(Reject_Remarks)), TRUE, FALSE)`,
			want: false, // Status is "Open", not "Reject"
		},
		// Formula 159: IF ISNULL with PRIORVALUE
		{
			name: "F159_IF_ISNULL_PRIORVALUE",
			expr: `IF(ISNULL(PRIORVALUE("CloseDate")), FALSE, TRUE)`,
			want: false, // PRIORVALUE("CloseDate") returns nil (not in _prior map), so ISNULL=true → FALSE
		},
		// Formula 160: CONTAINS with ISNEW
		{
			name: "F160_CONTAINS_ISNEW",
			expr: `AND(CONTAINS(What_is_required, "https://razorpay.slack.com"), ISNEW())`,
			want: false, // ISNEW() is false
		},
		// Formula 170: IF with ISPICKVAL branching to ISBLANK
		{
			name: "F170_IF_ISPICKVAL_branch_ISBLANK",
			expr: `IF(ISPICKVAL(Loss_Reason, "Flat"), ISBLANK(NULL), ISBLANK(NULL))`,
			want: true, // Loss_Reason is "Price" not "Flat", goes to else, ISBLANK(NULL)=true
		},
		// Formula 176: ISNEW with ISBLANK, ==null
		{
			name: "F176_ISNEW_ISBLANK_null",
			expr: `IF(ISNEW() && (ISBLANK(NULL) || NULL == NULL || ISBLANK(NULL) || NULL == NULL || ISBLANK(NULL) || NULL == NULL || ISBLANK(NULL) || NULL == NULL), true, false)`,
			want: false, // ISNEW() is false, so && short-circuits
		},
		// Formula 179: AND ISCHANGED ISPICKVAL
		{
			name: "F179_AND_ISCHANGED_ISPICKVAL",
			expr: `AND(ISCHANGED("StageName"), ISPICKVAL(Status, "Completed"))`,
			want: false, // Status is "Open", not "Completed"
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			program, err := expr.Compile(tt.expr, opts...)
			require.NoError(t, err, "compile error for %s", tt.name)
			got, err := expr.Run(program, env)
			require.NoError(t, err, "run error for %s", tt.name)
			assert.Equal(t, tt.want, got, "wrong result for %s", tt.name)
		})
	}
}

// TestSFFormula_CompileOnly tests that complex formulas compile without error,
// even if we don't verify the exact runtime result.
// These cover additional patterns not in the run tests above.
func TestSFFormula_CompileOnly(t *testing.T) {
	env := sfEnv()
	opts := append(sfOpts(), expr.Env(env))

	formulas := []struct {
		name string
		expr string
	}{
		// Formula 2: AND with NOT, ISCHANGED, ISBLANK, ISPICKVAL
		{
			name: "F2_complex_validation",
			expr: `NOT(Permission_LeadValidation) &&
				ISCHANGED("AM_manager") &&
				AM_manager == TRUE &&
				(ISBLANK(NULL) ||
				ISBLANK(NULL) ||
				ISBLANK(NULL))`,
		},
		// Formula 4: ISCHANGED with NOT(OR(...))
		{
			name: "F4_ISCHANGED_NOT_OR",
			expr: `ISCHANGED("BU_Owner") &&
				NOT(
					Profile_Name == "System Administrator" ||
					Profile_Name == "Integration User"
				)`,
		},
		// Formula 7: Complex multi-clause with PRIORVALUE, ISPICKVAL, !=
		{
			name: "F7_complex_multi_clause",
			expr: `(ISCHANGED("Owner_Role")
				&&
				((PRIORVALUE("Owner_Role") == "Direct Sales") || (PRIORVALUE("Owner_Role") == "Growth") || (PRIORVALUE("Owner_Role") == "KAM"))
				&&
				Owner_FirstName == "Integration")
				&&
				NOT(ISPICKVAL(Approval_Status, "Approved - ENTP to IUser"))
				&&
				(Profile_Name != "Integration User" || Profile_Name != "System Administrator" || Profile_Name != "Sales Manager")`,
		},
		// Formula 9: AND/OR with many PRIORVALUE comparisons
		{
			name: "F9_AND_OR_many_PRIORVALUE",
			expr: `AND(OR(
				PRIORVALUE("Owner_Role") == "IF - KAM", PRIORVALUE("Owner_Role") == "IF - SAM", PRIORVALUE("Owner_Role") == "IF - Startup", PRIORVALUE("Owner_Role") == "Hybrid - MM", PRIORVALUE("Owner_Role") == "Hybrid - SAM", PRIORVALUE("Owner_Role") == "Hybrid - KAM"
			),
			ISCHANGED("OwnerId"), User_Email == "adnan.noorie@razorpay.com")`,
		},
		// Formula 10: Deeply nested AND/OR with <>, RecordTypeId
		{
			name: "F10_deeply_nested",
			expr: `AND(
				ISCHANGED("OwnerId"),
				OR(
					AND(RecordTypeId == "012OX000000MgVF",
						Profile_Name <> "System Administrator"),
					AND(OwnerId <> "PARENT_OWNER",
						OR(Owner_Role == "IF - SAM",
							Owner_Role == "IF - KAM",
							Owner_Role == "IF - Startup"),
						OR(Owner_Role <> "Banking", Owner_Role <> "Partnership")
					)
				)
			)`,
		},
		// Formula 11: AND with NOT(OR(...)) with many email checks
		{
			name: "F11_NOT_OR_many_emails",
			expr: `AND(
				ISCHANGED("OwnerId"),
				NOT(
					OR(
						Profile_Name == "System Administrator",
						User_Email == "user1@razorpay.com",
						User_Email == "user2@razorpay.com",
						User_Email == "user3@razorpay.com"
					)
				)
			)`,
		},
		// Formula 17: !(OR(...)) with many conditions
		{
			name: "F17_bang_OR_ISPICKVAL",
			expr: `AND(!(OR((Profile_Name == "Integration User"), (Profile_Name == "System Administrator"), (User_Username == "user1@razorpay.com"))),
				ISCHANGED("OwnerId"),
				OR(
					ISPICKVAL(Ops_UW_Status, "PRE_VERIFICATION_COMPLETE"),
					ISPICKVAL(Ops_UW_Status, "OFFER_APPROVED"),
					ISPICKVAL(Ops_UW_Status, "OFFER_ACCEPTED")
				)
			)`,
		},
		// Formula 25: AND(AND(ISPICKVAL, ISBLANK(TEXT(...))), ISCHANGED)
		{
			name: "F25_nested_AND_ISPICKVAL_ISBLANK_TEXT",
			expr: `AND(AND(ISPICKVAL(Stages, "Solution Closure"), ISBLANK(Stages)), ISCHANGED("StageName"))`,
		},
		// Formula 28: Complex OR chain with ISPICKVAL, ISCHANGED, ISBLANK
		{
			name: "F28_complex_OR_chain",
			expr: `(ISPICKVAL(Status, "Issue Resolvable") && ISCHANGED("Status") && (ISBLANK(Description) || ISBLANK(NULL)) && RecordType_Name == "GMV-Wins" && NOT(Permission_LeadValidation))
			|| (ISPICKVAL(Status, "Qualified") && ISCHANGED("Status") && ISBLANK(NULL) && RecordType_Name == "GMV-Wins" && NOT(Permission_LeadValidation))
			|| (ISPICKVAL(Status, "Lead Won") && ISCHANGED("Status") && ISBLANK(Description) && RecordType_Name == "GMV-Wins" && NOT(Permission_LeadValidation))`,
		},
		// Formula 29: AND with multiple ISPICKVAL in OR
		{
			name: "F29_AND_multiple_ISPICKVAL_OR",
			expr: `AND(
				OR(
					ISPICKVAL(StageName, "Pricing"),
					ISPICKVAL(StageName, "Support"),
					ISPICKVAL(StageName, "Sales/Relationship"),
					ISPICKVAL(StageName, "Method Enablement")
				),
				ISBLANK(NULL),
				NOT(Permission_LeadValidation)
			)`,
		},
		// Formula 30: AND with ISPICKVAL + ISBLANK
		{
			name: "F30_AND_ISPICKVAL_ISBLANK",
			expr: `AND(ISPICKVAL(StageName, "Rank Upgrade"), ISBLANK(NULL))`,
		},
		// Formula 31: AND with OR of ISPICKVAL, OR of ISBLANK
		{
			name: "F31_AND_OR_ISPICKVAL_OR_ISBLANK",
			expr: `AND(
				OR(
					ISPICKVAL(Status, "Issue Resolvable"),
					ISPICKVAL(Status, "Issue Resolved"),
					ISPICKVAL(Status, "Closed Won"),
					ISPICKVAL(Status, "Closed Lost")
				),
				OR(
					ISBLANK(NULL),
					ISBLANK(NULL)
				)
			)`,
		},
		// Formula 35: AND with ISPICKVAL of PRIORVALUE
		{
			name: "F35_ISPICKVAL_PRIORVALUE",
			expr: `AND(
				ISPICKVAL(Request_Category, "Integration Error / Debugging"),
				ISPICKVAL(Stages, "Integration Closed"),
				OR(
					ISPICKVAL(PRIORVALUE("Stages"), "Integration In Progress"),
					ISPICKVAL(PRIORVALUE("Stages"), "New")
				)
			)`,
		},
		// Formula 39: AND with ISPICKVAL(PRIORVALUE(...))
		{
			name: "F39_ISPICKVAL_PRIORVALUE_complex",
			expr: `AND(ISPICKVAL(PRIORVALUE("User_Status"), "Inactive"),
				Inactive_Record == TRUE,
				NOT(ISPICKVAL(User_Status, "Dropped"))
			)`,
		},
		// Formula 43: Very complex nested validation (simplified)
		{
			name: "F43_complex_product_stage_validation",
			expr: `AND(ISPICKVAL(Product_Type, "Payment Gateway"),
				OR(
					AND(
						ISPICKVAL(Stages, "Integration Completed"), ISCHANGED("StageName"),
						NOT(OR(ISPICKVAL(Request_Category, "Integration Error / Debugging"), ISPICKVAL(Request_Category, "Integration Audit Post Go-live"))),
						OR(ISPICKVAL(PRIORVALUE("Stages"), "Testing & QA"), ISPICKVAL(PRIORVALUE("Stages"), "New")),
						OR(ISBLANK(NULL), LEN(Remarks) == 0, LEN(Description) == 0)
					),
					AND(
						NOT(OR(ISPICKVAL(Request_Category, "Integration Error / Debugging"), ISPICKVAL(Request_Category, "Integration Audit Post Go-live"))),
						OR(ISPICKVAL(Stages, "MX Go-Live"), ISPICKVAL(Stages, "Integration Closed")),
						OR(ISPICKVAL(PRIORVALUE("Stages"), "Testing & QA"), ISPICKVAL(PRIORVALUE("Stages"), "New")),
						OR(ISBLANK(NULL), LEN(Remarks) == 0, LEN(Description) == 0)
					)
				)
			)`,
		},
		// Formula 82: Complex with many ISPICKVAL in OR
		{
			name: "F82_many_ISPICKVAL_OR",
			expr: `AND(
				ISPICKVAL(Type, "Current_Account"),
				Owner_Role == "SME X",
				OR(
					ISPICKVAL(StageName, "Pitched"),
					ISPICKVAL(StageName, "Pending Confirmation"),
					ISPICKVAL(StageName, "Principal Approval"),
					ISPICKVAL(StageName, "Application Received")
				),
				OR(
					ISPICKVAL(CA_Pitched_Partner_Bank, ""),
					ISPICKVAL(Variant_Sold, ""),
					ISPICKVAL(Payment_Collected, ""),
					ISBLANK(Amount_Collected),
					ISBLANK(CloseDate)
				)
			)`,
		},
		// Formula 83: Complex with ! prefix, Ispickval case-insensitive, IsChanged case-insensitive
		{
			name: "F83_case_insensitive_functions",
			expr: `(
				(!Ispickval(StageName, "Open") && !Ispickval(StageName, "Working") && !Ispickval(StageName, "Qualify"))
				&&
				IsChanged("StageName")
				&&
				Ispickval(Type, "Current_Account")
				&&
				Ispickval(Corporate_Card_Interested, "")
				&&
				!ISNEW()
			) ||
			(
				(IsChanged("Corporate_Card_Interested"))
				&&
				(!Ispickval(StageName, "Open") && !Ispickval(StageName, "Working") && !Ispickval(StageName, "Qualify"))
				&&
				Ispickval(Corporate_Card_Interested, "")
				&&
				!ISNEW()
				&&
				Ispickval(Type, "Current_Account")
			)`,
		},
		// Formula 85: TEXT with <> comparison
		{
			name: "F85_TEXT_diamond",
			expr: `AND(
				ISPICKVAL(Type, "Incorporation"),
				PRIORVALUE("CA_Interested") == "Yes",
				CA_Interested <> "Yes"
			)`,
		},
		// Formula 92: ISNULL with && chain
		{
			name: "F92_ISNULL_chain",
			expr: `ISNULL(NULL) &&
				Upsell == TRUE &&
				!ISPICKVAL(Type, "Opfin") &&
				(LastModifiedBy_FirstName != "Arpitha" || LastModifiedBy_FirstName != "Integration") &&
				RecordType_DeveloperName == "Ezetap"`,
		},
		// Formula 93: ISBLANK + TEXT(StageName) == "Closed Won"
		{
			name: "F93_TEXT_equality",
			expr: `RecordType_Name == "Corporate Cards"
				&&
				ISBLANK(NULL)
				&&
				(StageName == "Closed Won")
				&& ISCHANGED("StageName")`,
		},
		// Formula 105: CONTAINS with Owner.UserRole
		{
			name: "F105_CONTAINS_nested_ref",
			expr: `CONTAINS(Owner_Role, "NIT") && (ISBLANK(NULL) || ISPICKVAL(LeadSource, "") || ISBLANK(NULL) || ISPICKVAL(CA_Pitched_Partner_Bank, ""))`,
		},
		// Formula 129: CASE with <> CASE + 1 (milestone progression)
		{
			name: "F129_CASE_diamond_CASE_plus1",
			expr: `AND(
				Profile_Name <> "Integration User",
				Profile_Name <> "rzprizeform Profile",
				RecordType_Name == "RIZE Opportunity",
				ISCHANGED("StageName"),
				StageName <> "Document Rejected",
				StageName <> "Document Link Shared",
				CASE(StageName,
					"Name Approval", 1,
					"Digital Signature", 2,
					"Incorporation Form Upload", 3,
					"Pending", 4,
					"Approved", 5,
					0
				)
				<>
				CASE(PRIORVALUE("StageName"),
					"Name Approval", 1,
					"Digital Signature", 2,
					"Incorporation Form Upload", 3,
					"Pending", 4,
					"Approved", 5,
					0
				)
				+ 1
			)`,
		},
		// Formula 132: CASE <> CASE + 1 with TEXT <> checks
		{
			name: "F132_CASE_plus1_TEXT_diamond",
			expr: `AND(
				Profile_Name <> "Integration User",
				RecordType_Name == "RIZE Opportunity",
				ISCHANGED("StageName"),
				StageName <> "Closed Lost",
				StageName <> "Open",
				StageName <> "Working",
				StageName <> "Closed",
				CASE(StageName,
					"Qualify", 1,
					"Agreement Signed", 2,
					"Onboarding", 3,
					"Closed Won", 4,
					0)
				<>
				CASE(PRIORVALUE("StageName"),
					"Qualify", 1,
					"Agreement Signed", 2,
					"Onboarding", 3,
					"Closed Won", 4,
					0
				)
				+ 1
			)`,
		},
		// Formula 139: INCLUDES with multiple checks
		{
			name: "F139_INCLUDES_multiple",
			expr: `AND(INCLUDES(Methods, "Cardless EMI"), Service_Provider != "Axio", Service_Provider != "Fibe", UserRole_RollupDescription <> "Automation User", ISCHANGED("Service_Provider"))`,
		},
		// Formula 161: IF with ISPICKVAL OR chain, RecordType
		{
			name: "F161_IF_ISPICKVAL_OR_RecordType",
			expr: `IF(Developer_Sign_off == FALSE &&
				(ISPICKVAL(Ticket_Status, "Work In Progress") ||
				ISPICKVAL(Ticket_Status, "Testing") ||
				ISPICKVAL(Ticket_Status, "Under Deployment") ||
				ISPICKVAL(Ticket_Status, "Closed")) &&
				RecordType_Name == "Sprint Ticket",
			TRUE, FALSE)`,
		},
		// Formula 172: Very deeply nested
		{
			name: "F172_deeply_nested_AND_OR",
			expr: `AND(AND(ISPICKVAL(Lead_Status, "Converted"), AND(ISPICKVAL(Type, "Partnership"), OR(ISBLANK(NULL),
				ISPICKVAL(Confidence_Level, ""), true, ISBLANK(NULL), ISBLANK(NULL)))), User_Email != "anuj.jain@razorpay.com")`,
		},
		// Formula 173: IF with PRIORVALUE equality
		{
			name: "F173_IF_PRIORVALUE",
			expr: `IF(PRIORVALUE("Status") == "Converted" || PRIORVALUE("Status") == "Qualified" || PRIORVALUE("Status") == "Unqualified", true, false)`,
		},
		// Formula 178: AND with ISCHANGED, Subject ==, <>
		{
			name: "F178_ISCHANGED_Subject",
			expr: `AND(ISCHANGED("OwnerId"), Subject == "Revamp Pricing", UserRole_RollupDescription <> "Automation User", Profile_Name <> "System Administrator", User_Email <> "exception@razorpay.com")`,
		},
	}

	for _, tt := range formulas {
		t.Run(tt.name, func(t *testing.T) {
			program, err := expr.Compile(tt.expr, opts...)
			require.NoError(t, err, "compile error for %s", tt.name)
			_, err = expr.Run(program, env)
			require.NoError(t, err, "run error for %s", tt.name)
		})
	}
}

// TestSFFormula_SpecificResults verifies specific SF formula patterns
// with carefully crafted env to produce known results.
func TestSFFormula_SpecificResults(t *testing.T) {
	opts := sfOpts()

	t.Run("ISPICKVAL_match", func(t *testing.T) {
		env := map[string]any{"Status": "Closed Won"}
		program, err := expr.Compile(`ISPICKVAL(Status, "Closed Won")`, append(opts, expr.Env(env))...)
		require.NoError(t, err)
		got, err := expr.Run(program, env)
		require.NoError(t, err)
		assert.Equal(t, true, got)
	})

	t.Run("ISPICKVAL_empty_match_nil", func(t *testing.T) {
		env := map[string]any{"Status": nil}
		program, err := expr.Compile(`ISPICKVAL(Status, "")`, append(opts, expr.Env(env))...)
		require.NoError(t, err)
		got, err := expr.Run(program, env)
		require.NoError(t, err)
		assert.Equal(t, true, got) // SF: nil picklist == "" is true
	})

	t.Run("ISPICKVAL_empty_no_match", func(t *testing.T) {
		env := map[string]any{"Status": "Open"}
		program, err := expr.Compile(`ISPICKVAL(Status, "")`, append(opts, expr.Env(env))...)
		require.NoError(t, err)
		got, err := expr.Run(program, env)
		require.NoError(t, err)
		assert.Equal(t, false, got)
	})

	t.Run("INCLUDES_match", func(t *testing.T) {
		env := map[string]any{"Product": "PG;Opfin;CA"}
		program, err := expr.Compile(`INCLUDES(Product, "Opfin")`, append(opts, expr.Env(env))...)
		require.NoError(t, err)
		got, err := expr.Run(program, env)
		require.NoError(t, err)
		assert.Equal(t, true, got)
	})

	t.Run("INCLUDES_no_match", func(t *testing.T) {
		env := map[string]any{"Product": "PG;CA"}
		program, err := expr.Compile(`INCLUDES(Product, "Opfin")`, append(opts, expr.Env(env))...)
		require.NoError(t, err)
		got, err := expr.Run(program, env)
		require.NoError(t, err)
		assert.Equal(t, false, got)
	})

	t.Run("ISCHANGED_true", func(t *testing.T) {
		env := map[string]any{
			"_changed": map[string]bool{"OwnerId": true},
		}
		program, err := expr.Compile(`ISCHANGED("OwnerId")`, append(opts, expr.Env(env))...)
		require.NoError(t, err)
		got, err := expr.Run(program, env)
		require.NoError(t, err)
		assert.Equal(t, true, got)
	})

	t.Run("ISCHANGED_false", func(t *testing.T) {
		env := map[string]any{
			"_changed": map[string]bool{"OwnerId": false},
		}
		program, err := expr.Compile(`ISCHANGED("OwnerId")`, append(opts, expr.Env(env))...)
		require.NoError(t, err)
		got, err := expr.Run(program, env)
		require.NoError(t, err)
		assert.Equal(t, false, got)
	})

	t.Run("PRIORVALUE_returns_old", func(t *testing.T) {
		env := map[string]any{
			"_prior": map[string]any{"StageName": "Open"},
		}
		program, err := expr.Compile(`PRIORVALUE("StageName") == "Open"`, append(opts, expr.Env(env))...)
		require.NoError(t, err)
		got, err := expr.Run(program, env)
		require.NoError(t, err)
		assert.Equal(t, true, got)
	})

	t.Run("ISNEW_true", func(t *testing.T) {
		env := map[string]any{"_isNew": true}
		program, err := expr.Compile(`ISNEW()`, append(opts, expr.Env(env))...)
		require.NoError(t, err)
		got, err := expr.Run(program, env)
		require.NoError(t, err)
		assert.Equal(t, true, got)
	})

	t.Run("ISNEW_false", func(t *testing.T) {
		env := map[string]any{"_isNew": false}
		program, err := expr.Compile(`ISNEW()`, append(opts, expr.Env(env))...)
		require.NoError(t, err)
		got, err := expr.Run(program, env)
		require.NoError(t, err)
		assert.Equal(t, false, got)
	})

	t.Run("CASE_stage_regression_detected", func(t *testing.T) {
		env := map[string]any{
			"StageName": "Open",
			"_changed":  map[string]bool{"StageName": true},
			"_prior":    map[string]any{"StageName": "Working"},
		}
		// Open=1, Working=2 → 1 < 2 = true (regression detected)
		program, err := expr.Compile(`ISCHANGED("StageName") &&
			CASE(StageName, "Open", 1, "Working", 2, "Qualify", 3, 0)
			<
			CASE(PRIORVALUE("StageName"), "Open", 1, "Working", 2, "Qualify", 3, 0)`,
			append(opts, expr.Env(env))...)
		require.NoError(t, err)
		got, err := expr.Run(program, env)
		require.NoError(t, err)
		assert.Equal(t, true, got)
	})

	t.Run("CASE_stage_no_regression", func(t *testing.T) {
		env := map[string]any{
			"StageName": "Working",
			"_changed":  map[string]bool{"StageName": true},
			"_prior":    map[string]any{"StageName": "Open"},
		}
		// Working=2, Open=1 → 2 < 1 = false (no regression)
		program, err := expr.Compile(`ISCHANGED("StageName") &&
			CASE(StageName, "Open", 1, "Working", 2, "Qualify", 3, 0)
			<
			CASE(PRIORVALUE("StageName"), "Open", 1, "Working", 2, "Qualify", 3, 0)`,
			append(opts, expr.Env(env))...)
		require.NoError(t, err)
		got, err := expr.Run(program, env)
		require.NoError(t, err)
		assert.Equal(t, false, got)
	})

	t.Run("complex_validation_rule_true", func(t *testing.T) {
		env := map[string]any{
			"StageName":       "Closed Won",
			"Loss_Reason":     "Other",
			"Other_Reason":    "",
			"RecordType_Name": "Ezetap",
			"_changed":        map[string]bool{},
			"_prior":          map[string]any{},
			"_isNew":          false,
			"_tz":             time.UTC,
		}
		program, err := expr.Compile(
			`ISPICKVAL(StageName, "Closed Won") && ISPICKVAL(Loss_Reason, "Other") && ISBLANK(Other_Reason) && RecordType_Name == "Ezetap"`,
			append(opts, expr.Env(env))...)
		require.NoError(t, err)
		got, err := expr.Run(program, env)
		require.NoError(t, err)
		assert.Equal(t, true, got)
	})
}

func TestSFFormula_BLANKVALUE(t *testing.T) {
	opts := sfOpts()
	tests := []struct {
		name string
		expr string
		env  map[string]any
		want any
	}{
		{
			name: "nil field returns substitute",
			expr: `BLANKVALUE(Phone, "N/A")`,
			env:  map[string]any{"Phone": nil},
			want: "N/A",
		},
		{
			name: "empty string returns substitute",
			expr: `BLANKVALUE(Phone, "N/A")`,
			env:  map[string]any{"Phone": ""},
			want: "N/A",
		},
		{
			name: "whitespace-only returns substitute",
			expr: `BLANKVALUE(Phone, "N/A")`,
			env:  map[string]any{"Phone": "   "},
			want: "N/A",
		},
		{
			name: "non-blank returns original",
			expr: `BLANKVALUE(Phone, "N/A")`,
			env:  map[string]any{"Phone": "555-1234"},
			want: "555-1234",
		},
		{
			name: "numeric zero is not blank",
			expr: `BLANKVALUE(Amount, 100)`,
			env:  map[string]any{"Amount": 0},
			want: 0,
		},
		{
			name: "numeric value is not blank",
			expr: `BLANKVALUE(Amount, 100)`,
			env:  map[string]any{"Amount": 42.5},
			want: 42.5,
		},
		{
			name: "nil numeric returns substitute",
			expr: `BLANKVALUE(Amount, 100)`,
			env:  map[string]any{"Amount": nil},
			want: 100,
		},
		{
			name: "false is not blank",
			expr: `BLANKVALUE(IsActive, true)`,
			env:  map[string]any{"IsActive": false},
			want: false,
		},
		{
			name: "nested: BLANKVALUE with expression substitute",
			expr: `BLANKVALUE(Phone, "Call " & Name)`,
			env:  map[string]any{"Phone": nil, "Name": "Acme"},
			want: "Call Acme",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			program, err := expr.Compile(tt.expr, append(opts, expr.Env(tt.env))...)
			require.NoError(t, err)
			got, err := expr.Run(program, tt.env)
			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestSFFormula_EmptyAndZeroDateComparison(t *testing.T) {
	opts := sfOpts()
	someDate := time.Date(2024, 6, 15, 0, 0, 0, 0, time.UTC)
	tests := []struct {
		name string
		expr string
		env  map[string]any
		want any
	}{
		// any comparison between "" or 0 and time.Time always returns false
		{name: "empty_gt_datevalue", expr: `"" > DATEVALUE(D)`, env: map[string]any{"D": someDate}, want: false},
		{name: "empty_lt_datevalue", expr: `"" < DATEVALUE(D)`, env: map[string]any{"D": someDate}, want: false},
		{name: "empty_gte_datevalue", expr: `"" >= DATEVALUE(D)`, env: map[string]any{"D": someDate}, want: false},
		{name: "empty_lte_datevalue", expr: `"" <= DATEVALUE(D)`, env: map[string]any{"D": someDate}, want: false},
		{name: "datevalue_gt_empty", expr: `DATEVALUE(D) > ""`, env: map[string]any{"D": someDate}, want: false},
		{name: "datevalue_lt_empty", expr: `DATEVALUE(D) < ""`, env: map[string]any{"D": someDate}, want: false},
		{name: "datevalue_gte_empty", expr: `DATEVALUE(D) >= ""`, env: map[string]any{"D": someDate}, want: false},
		{name: "datevalue_lte_empty", expr: `DATEVALUE(D) <= ""`, env: map[string]any{"D": someDate}, want: false},
		{name: "zero_gt_datevalue", expr: `0 > DATEVALUE(D)`, env: map[string]any{"D": someDate}, want: false},
		{name: "zero_lt_datevalue", expr: `0 < DATEVALUE(D)`, env: map[string]any{"D": someDate}, want: false},
		{name: "datevalue_gt_zero", expr: `DATEVALUE(D) > 0`, env: map[string]any{"D": someDate}, want: false},
		{name: "datevalue_lt_zero", expr: `DATEVALUE(D) < 0`, env: map[string]any{"D": someDate}, want: false},
		{name: "empty_gt_date_field", expr: `"" > D`, env: map[string]any{"D": someDate}, want: false},
		{name: "zero_gt_date_field", expr: `0 > D`, env: map[string]any{"D": someDate}, want: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			program, err := expr.Compile(tt.expr, append(opts, expr.Env(tt.env))...)
			require.NoError(t, err)
			got, err := expr.Run(program, tt.env)
			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestSFFormula_NULLVALUE(t *testing.T) {
	opts := sfOpts()
	tests := []struct {
		name string
		expr string
		env  map[string]any
		want any
	}{
		{
			name: "nil returns substitute",
			expr: `NULLVALUE(Phone, "N/A")`,
			env:  map[string]any{"Phone": nil},
			want: "N/A",
		},
		{
			name: "empty string is NOT null - returns original",
			expr: `NULLVALUE(Phone, "N/A")`,
			env:  map[string]any{"Phone": ""},
			want: "",
		},
		{
			name: "non-nil returns original",
			expr: `NULLVALUE(Phone, "N/A")`,
			env:  map[string]any{"Phone": "555-1234"},
			want: "555-1234",
		},
		{
			name: "zero is not null",
			expr: `NULLVALUE(Amount, 100)`,
			env:  map[string]any{"Amount": 0},
			want: 0,
		},
		{
			name: "nil numeric returns substitute",
			expr: `NULLVALUE(Amount, 100)`,
			env:  map[string]any{"Amount": nil},
			want: 100,
		},
		{
			name: "false is not null",
			expr: `NULLVALUE(IsActive, true)`,
			env:  map[string]any{"IsActive": false},
			want: false,
		},
		{
			name: "whitespace is NOT null - returns original",
			expr: `NULLVALUE(Name, "Unknown")`,
			env:  map[string]any{"Name": "   "},
			want: "   ",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			program, err := expr.Compile(tt.expr, append(opts, expr.Env(tt.env))...)
			require.NoError(t, err)
			got, err := expr.Run(program, tt.env)
			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
