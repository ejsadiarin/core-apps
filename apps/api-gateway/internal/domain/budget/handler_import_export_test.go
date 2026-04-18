package budget

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"core-gateway/internal/domain/auth"
	"core-gateway/internal/repository/sqlc"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type stubRow struct {
	values []any
	err    error
}

func (r *stubRow) Scan(dest ...any) error {
	if r.err != nil {
		return r.err
	}
	if len(dest) != len(r.values) {
		return errors.New("scan destination count mismatch")
	}
	for i := range dest {
		switch d := dest[i].(type) {
		case *uuid.UUID:
			*d = r.values[i].(uuid.UUID)
		case *string:
			*d = r.values[i].(string)
		case *pgtype.Numeric:
			*d = r.values[i].(pgtype.Numeric)
		case *pgtype.Text:
			*d = r.values[i].(pgtype.Text)
		case *pgtype.Date:
			*d = r.values[i].(pgtype.Date)
		case *pgtype.Timestamp:
			*d = r.values[i].(pgtype.Timestamp)
		case *pgtype.Bool:
			*d = r.values[i].(pgtype.Bool)
		case *pgtype.UUID:
			*d = r.values[i].(pgtype.UUID)
		default:
			return errors.New("unsupported scan destination type")
		}
	}
	return nil
}

type stubRows struct {
	rows  [][]any
	idx   int
	err   error
	close bool
}

func (r *stubRows) Close() { r.close = true }

func (r *stubRows) Err() error { return r.err }

func (r *stubRows) CommandTag() pgconn.CommandTag { return pgconn.CommandTag{} }

func (r *stubRows) FieldDescriptions() []pgconn.FieldDescription { return nil }

func (r *stubRows) Next() bool {
	if r.err != nil {
		return false
	}
	if r.idx >= len(r.rows) {
		return false
	}
	r.idx++
	return true
}

func (r *stubRows) Scan(dest ...any) error {
	if r.idx == 0 || r.idx > len(r.rows) {
		return errors.New("scan called without current row")
	}
	row := r.rows[r.idx-1]
	if len(dest) != len(row) {
		return errors.New("scan destination count mismatch")
	}
	for i := range dest {
		switch d := dest[i].(type) {
		case *uuid.UUID:
			*d = row[i].(uuid.UUID)
		case *string:
			*d = row[i].(string)
		case *pgtype.Numeric:
			*d = row[i].(pgtype.Numeric)
		case *pgtype.Text:
			*d = row[i].(pgtype.Text)
		case *pgtype.Date:
			*d = row[i].(pgtype.Date)
		case *pgtype.Timestamp:
			*d = row[i].(pgtype.Timestamp)
		case *pgtype.Bool:
			*d = row[i].(pgtype.Bool)
		case *pgtype.UUID:
			*d = row[i].(pgtype.UUID)
		default:
			return errors.New("unsupported scan destination type")
		}
	}
	return nil
}

func (r *stubRows) Values() ([]any, error) { return nil, errors.New("not implemented") }

func (r *stubRows) RawValues() [][]byte { return nil }

func (r *stubRows) Conn() *pgx.Conn { return nil }

type stubDB struct {
	queryRows map[string][][]any
	queryErr  map[string]error
	queryRow  map[string][]any

	execCalls []string
	lastExec  map[string][]any
	lastQuery map[string][]any
}

func newStubDB() *stubDB {
	return &stubDB{
		queryRows: make(map[string][][]any),
		queryErr:  make(map[string]error),
		queryRow:  make(map[string][]any),
		lastExec:  make(map[string][]any),
		lastQuery: make(map[string][]any),
	}
}

func (d *stubDB) Exec(_ context.Context, sql string, args ...interface{}) (pgconn.CommandTag, error) {
	d.execCalls = append(d.execCalls, sql)
	d.lastExec[sql] = args
	return pgconn.CommandTag{}, nil
}

func (d *stubDB) Query(_ context.Context, sql string, args ...interface{}) (pgx.Rows, error) {
	d.lastQuery[sql] = args
	if err, ok := d.queryErr[sql]; ok {
		return &stubRows{err: err}, nil
	}
	rows := d.queryRows[sql]
	return &stubRows{rows: rows}, nil
}

func (d *stubDB) QueryRow(_ context.Context, sql string, args ...interface{}) pgx.Row {
	d.lastQuery[sql] = args
	if row, ok := d.queryRow[sql]; ok {
		return &stubRow{values: row}
	}

	now := pgtype.Timestamp{Time: time.Now().UTC(), Valid: true}

	if strings.HasPrefix(sql, "-- name: CreateIncome :one") {
		return &stubRow{values: []any{
			uuid.New(),          // id
			args[7].(uuid.UUID), // user_id
			args[0].(pgtype.Numeric),
			args[1].(pgtype.Text),
			args[2].(pgtype.Date),
			args[3].(pgtype.Text),
			args[4].(pgtype.Text),
			args[5].(pgtype.Date),
			now,
			now,
			args[6].(pgtype.Date),
			pgtype.Bool{Bool: false, Valid: true},
			"posted",
			pgtype.UUID{Valid: false},
		}}
	}

	if strings.HasPrefix(sql, "-- name: CreateExpense :one") {
		return &stubRow{values: []any{
			uuid.New(),       // id
			args[0].(string), // description
			args[1].(pgtype.Numeric),
			args[2].(pgtype.Text),
			args[3].(pgtype.UUID),
			args[4].(pgtype.Date),
			now,
			now,
			args[5].(pgtype.Text),
			args[6].(uuid.UUID),
			args[7].(pgtype.Text),
			args[8].(pgtype.Date),
			args[9].(pgtype.Date),
			args[10].(pgtype.UUID),
			"posted",
			pgtype.UUID{Valid: false},
		}}
	}

	return &stubRow{err: errors.New("no row configured")}
}

func setAuthenticatedUser(c echo.Context, userID uuid.UUID) {
	c.Set("user", &auth.UserContext{ID: userID, Role: auth.RoleUser, Email: "user@example.com"})
}

func TestExportBudgetJSONReturnsExpectedStructure(t *testing.T) {
	e := echo.New()
	db := newStubDB()
	queries := sqlc.New(db)
	logger := zerolog.Nop()
	h := NewHandler(queries, &logger)

	userID := uuid.New()
	categoryID := uuid.New()
	priorityID := uuid.New()
	incomeID := uuid.New()
	expenseID := uuid.New()

	incomeQuery := `-- name: ExportIncomes :many
SELECT id, user_id, amount, currency, date, description, recurring_type, start_date, created_at, updated_at, end_date, exclude_from_calculations, status, source_rule_id FROM budget_incomes
WHERE user_id = $1
    AND status = 'posted'
ORDER BY date ASC, created_at ASC
`

	expenseQuery := `-- name: ExportExpenses :many
SELECT e.id, e.description, e.amount, e.currency, e.category_id, e.expense_date, e.created_at, e.updated_at, e.notes, e.user_id, e.recurring_type, e.start_date, e.end_date, e.priority_group_id, e.status, e.source_rule_id, c.name as category_name
FROM budget_expenses e
LEFT JOIN budget_categories c ON e.category_id = c.id
WHERE e.user_id = $1
    AND e.status = 'posted'
ORDER BY e.expense_date ASC, e.created_at ASC
`

	db.queryRows[incomeQuery] = [][]any{{
		incomeID,
		userID,
		float64ToNumeric(5000),
		pgtype.Text{String: "USD", Valid: true},
		stringToDate("2026-04-01"),
		pgtype.Text{String: "Salary", Valid: true},
		pgtype.Text{Valid: false},
		pgtype.Date{Valid: false},
		pgtype.Timestamp{Valid: true},
		pgtype.Timestamp{Valid: true},
		pgtype.Date{Valid: false},
		pgtype.Bool{Bool: false, Valid: true},
		"posted",
		pgtype.UUID{Valid: false},
	}}

	db.queryRows[expenseQuery] = [][]any{{
		expenseID,
		"Rent",
		float64ToNumeric(1200),
		pgtype.Text{String: "USD", Valid: true},
		pgtype.UUID{Bytes: categoryID, Valid: true},
		stringToDate("2026-04-02"),
		pgtype.Timestamp{Valid: true},
		pgtype.Timestamp{Valid: true},
		pgtype.Text{String: "April rent", Valid: true},
		userID,
		pgtype.Text{Valid: false},
		pgtype.Date{Valid: false},
		pgtype.Date{Valid: false},
		pgtype.UUID{Bytes: priorityID, Valid: true},
		"posted",
		pgtype.UUID{Valid: false},
		pgtype.Text{String: "Housing", Valid: true},
	}}

	req := httptest.NewRequest(http.MethodGet, "/api/budget/export", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	setAuthenticatedUser(c, userID)

	err := h.ExportBudgetJSON(c)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, rec.Code)

	var payload BudgetExportPayload
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &payload))

	assert.Equal(t, budgetImportExportSchemaVersion, payload.Metadata.SchemaVersion)
	assert.Equal(t, budgetImportExportSource, payload.Metadata.Source)
	require.Len(t, payload.Incomes, 1)
	require.Len(t, payload.Expenses, 1)

	assert.Equal(t, incomeID.String(), payload.Incomes[0].ID)
	assert.Equal(t, "Salary", payload.Incomes[0].Description)
	assert.Equal(t, 5000.0, payload.Incomes[0].Amount)
	assert.Equal(t, "2026-04-01", payload.Incomes[0].Date)

	assert.Equal(t, expenseID.String(), payload.Expenses[0].ID)
	assert.Equal(t, "Rent", payload.Expenses[0].Description)
	assert.Equal(t, "Housing", payload.Expenses[0].CategoryName)
	assert.Equal(t, priorityID.String(), *payload.Expenses[0].PriorityGroupID)
}

func TestImportBudgetJSONSkipsAndConflictsOnExactDate(t *testing.T) {
	e := echo.New()
	db := newStubDB()
	queries := sqlc.New(db)
	logger := zerolog.Nop()
	h := NewHandler(queries, &logger)

	userID := uuid.New()
	incomeExistingID := uuid.New()
	expenseExistingID := uuid.New()
	categoryID := uuid.New()

	incomeQuery := `-- name: ExportIncomes :many
SELECT id, user_id, amount, currency, date, description, recurring_type, start_date, created_at, updated_at, end_date, exclude_from_calculations, status, source_rule_id FROM budget_incomes
WHERE user_id = $1
    AND status = 'posted'
ORDER BY date ASC, created_at ASC
`

	expenseQuery := `-- name: ExportExpenses :many
SELECT e.id, e.description, e.amount, e.currency, e.category_id, e.expense_date, e.created_at, e.updated_at, e.notes, e.user_id, e.recurring_type, e.start_date, e.end_date, e.priority_group_id, e.status, e.source_rule_id, c.name as category_name
FROM budget_expenses e
LEFT JOIN budget_categories c ON e.category_id = c.id
WHERE e.user_id = $1
    AND e.status = 'posted'
ORDER BY e.expense_date ASC, e.created_at ASC
`

	listCategoriesQuery := `-- name: ListCategories :many
SELECT id, name, color, icon, created_at, user_id FROM budget_categories
WHERE user_id = $1
ORDER BY name
`

	db.queryRows[incomeQuery] = [][]any{{
		incomeExistingID,
		userID,
		float64ToNumeric(2000),
		pgtype.Text{String: "USD", Valid: true},
		stringToDate("2026-04-10"),
		pgtype.Text{String: "Paycheck", Valid: true},
		pgtype.Text{Valid: false},
		pgtype.Date{Valid: false},
		pgtype.Timestamp{Valid: true},
		pgtype.Timestamp{Valid: true},
		pgtype.Date{Valid: false},
		pgtype.Bool{Bool: false, Valid: true},
		"posted",
		pgtype.UUID{Valid: false},
	}}

	db.queryRows[expenseQuery] = [][]any{{
		expenseExistingID,
		"Groceries",
		float64ToNumeric(40),
		pgtype.Text{String: "USD", Valid: true},
		pgtype.UUID{Bytes: categoryID, Valid: true},
		stringToDate("2026-04-12"),
		pgtype.Timestamp{Valid: true},
		pgtype.Timestamp{Valid: true},
		pgtype.Text{Valid: false},
		userID,
		pgtype.Text{Valid: false},
		pgtype.Date{Valid: false},
		pgtype.Date{Valid: false},
		pgtype.UUID{Valid: false},
		"posted",
		pgtype.UUID{Valid: false},
		pgtype.Text{String: "Food", Valid: true},
	}}

	db.queryRows[listCategoriesQuery] = [][]any{{
		categoryID,
		"Food",
		pgtype.Text{Valid: false},
		pgtype.Text{Valid: false},
		pgtype.Timestamp{Valid: true},
		userID,
	}}

	bodyPayload := BudgetExportPayload{
		Incomes: []ImportIncomeRecord{
			{Date: "2026-04-10", Description: "Paycheck", Amount: 2000, Currency: "USD"},
			{Date: "2026-04-10", Description: "Paycheck", Amount: 2100, Currency: "USD"},
			{Date: "2026-04-11", Description: "Side Gig", Amount: 300, Currency: "USD"},
		},
		Expenses: []ImportExpenseRecord{
			{ExpenseDate: "2026-04-12", Description: "Groceries", Amount: 40, Currency: "USD", CategoryName: "Food"},
			{ExpenseDate: "2026-04-12", Description: "Groceries", Amount: 45, Currency: "USD", CategoryName: "Food"},
			{ExpenseDate: "2026-04-13", Description: "Lunch", Amount: 12, Currency: "USD", CategoryName: "Food"},
		},
	}
	bodyBytes, err := json.Marshal(bodyPayload)
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/api/budget/import", bytes.NewReader(bodyBytes))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	setAuthenticatedUser(c, userID)

	err = h.ImportBudgetJSON(c)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, rec.Code)

	var result BudgetImportResult
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &result))

	assert.Equal(t, 1, result.Summary.IncomesCreated)
	assert.Equal(t, 1, result.Summary.IncomesSkipped)
	assert.Equal(t, 1, result.Summary.IncomesConflicts)
	assert.Equal(t, 1, result.Summary.ExpensesCreated)
	assert.Equal(t, 1, result.Summary.ExpensesSkipped)
	assert.Equal(t, 1, result.Summary.ExpensesConflicts)
	require.Len(t, result.Conflicts, 2)

	assert.Equal(t, ImportMergeActionSkipExisting, result.Incomes[0].Action)
	assert.Equal(t, ImportMergeActionConflict, result.Incomes[1].Action)
	assert.Equal(t, ImportMergeActionCreate, result.Incomes[2].Action)

	assert.Equal(t, ImportMergeActionSkipExisting, result.Expenses[0].Action)
	assert.Equal(t, ImportMergeActionConflict, result.Expenses[1].Action)
	assert.Equal(t, ImportMergeActionCreate, result.Expenses[2].Action)

	conflictByEntity := map[string]ImportConflictDetail{}
	for _, c := range result.Conflicts {
		conflictByEntity[c.Entity] = c
	}

	assert.True(t, hasDifferenceField(conflictByEntity["income"].Differences, "amount"))
	assert.True(t, hasDifferenceField(conflictByEntity["expense"].Differences, "amount"))

}
