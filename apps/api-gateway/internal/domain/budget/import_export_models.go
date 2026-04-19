package budget

// ImportExportMetadata describes high-level payload metadata for JSON import/export.
type ImportExportMetadata struct {
	SchemaVersion string `json:"schema_version"`
	ExportedAt    string `json:"exported_at"`
	Source        string `json:"source"`
}

// BudgetExportPayload is the JSON contract returned by export endpoints.
type BudgetExportPayload struct {
	Metadata ImportExportMetadata  `json:"metadata"`
	Incomes  []ImportIncomeRecord  `json:"incomes"`
	Expenses []ImportExpenseRecord `json:"expenses"`
}

// ImportIncomeRecord represents one income item in import/export payloads.
type ImportIncomeRecord struct {
	ID            string  `json:"id,omitempty"`
	Amount        float64 `json:"amount"`
	Currency      string  `json:"currency"`
	Date          string  `json:"date"`
	Description   string  `json:"description,omitempty"`
	RecurringType *string `json:"recurring_type,omitempty"`
	StartDate     *string `json:"start_date,omitempty"`
	EndDate       *string `json:"end_date,omitempty"`
}

// ImportExpenseRecord represents one expense item in import/export payloads.
type ImportExpenseRecord struct {
	ID              string  `json:"id,omitempty"`
	Description     string  `json:"description"`
	Amount          float64 `json:"amount"`
	Currency        string  `json:"currency"`
	ExpenseDate     string  `json:"expense_date"`
	CategoryName    string  `json:"category_name,omitempty"`
	Notes           *string `json:"notes,omitempty"`
	RecurringType   *string `json:"recurring_type,omitempty"`
	StartDate       *string `json:"start_date,omitempty"`
	EndDate         *string `json:"end_date,omitempty"`
	PriorityGroupID *string `json:"priority_group_id,omitempty"`
	IsDebt          bool    `json:"is_debt"`
}

type ImportMergeAction string

const (
	ImportMergeActionCreate       ImportMergeAction = "CREATE"
	ImportMergeActionSkipExisting ImportMergeAction = "SKIP_EXISTING"
	ImportMergeActionConflict     ImportMergeAction = "CONFLICT"
)

// ImportFieldDifference captures a mismatched field between incoming and existing items.
type ImportFieldDifference struct {
	Field    string `json:"field"`
	Incoming any    `json:"incoming"`
	Existing any    `json:"existing"`
}

// ImportConflictDetail provides conflict diagnostics for import responses.
type ImportConflictDetail struct {
	Entity      string                  `json:"entity"`
	Date        string                  `json:"date"`
	ExistingID  string                  `json:"existing_id,omitempty"`
	Incoming    any                     `json:"incoming"`
	Existing    any                     `json:"existing"`
	Differences []ImportFieldDifference `json:"differences"`
}

// ImportMatchResult is returned by pure matcher functions.
type ImportMatchResult struct {
	Action            ImportMergeAction     `json:"action"`
	MatchedExistingID *string               `json:"matched_existing_id,omitempty"`
	Conflict          *ImportConflictDetail `json:"conflict,omitempty"`
}

// BudgetImportResultMetadata describes metadata returned from an import operation.
type BudgetImportResultMetadata struct {
	ImportedAt string `json:"imported_at"`
	DryRun     bool   `json:"dry_run"`
}

type BudgetImportSummary struct {
	IncomesCreated    int `json:"incomes_created"`
	IncomesSkipped    int `json:"incomes_skipped"`
	IncomesConflicts  int `json:"incomes_conflicts"`
	ExpensesCreated   int `json:"expenses_created"`
	ExpensesSkipped   int `json:"expenses_skipped"`
	ExpensesConflicts int `json:"expenses_conflicts"`
}

type ImportIncomeResult struct {
	InputIndex int                   `json:"input_index"`
	Action     ImportMergeAction     `json:"action"`
	ExistingID *string               `json:"existing_id,omitempty"`
	Conflict   *ImportConflictDetail `json:"conflict,omitempty"`
}

type ImportExpenseResult struct {
	InputIndex int                   `json:"input_index"`
	Action     ImportMergeAction     `json:"action"`
	ExistingID *string               `json:"existing_id,omitempty"`
	Conflict   *ImportConflictDetail `json:"conflict,omitempty"`
}

// BudgetImportResult is the JSON contract returned by import endpoints.
type BudgetImportResult struct {
	Metadata  BudgetImportResultMetadata `json:"metadata"`
	Summary   BudgetImportSummary        `json:"summary"`
	Incomes   []ImportIncomeResult       `json:"incomes"`
	Expenses  []ImportExpenseResult      `json:"expenses"`
	Conflicts []ImportConflictDetail     `json:"conflicts,omitempty"`
}
