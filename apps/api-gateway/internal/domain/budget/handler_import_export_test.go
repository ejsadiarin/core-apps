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

func assignScanValue(dest any, value any) error {
	switch d := dest.(type) {
	case *uuid.UUID:
		*d = value.(uuid.UUID)
	case *string:
		*d = value.(string)
	case *int64:
		*d = value.(int64)
	case *pgtype.Numeric:
		*d = value.(pgtype.Numeric)
	case *pgtype.Text:
		*d = value.(pgtype.Text)
	case *pgtype.Date:
		*d = value.(pgtype.Date)
	case *pgtype.Timestamp:
		*d = value.(pgtype.Timestamp)
	case *pgtype.Bool:
		*d = value.(pgtype.Bool)
	case *pgtype.UUID:
		*d = value.(pgtype.UUID)
	default:
		return errors.New("unsupported scan destination type")
	}
	return nil
}

func (r *stubRow) Scan(dest ...any) error {
	if r.err != nil {
		return r.err
	}
	if len(dest) != len(r.values) {
		return errors.New("scan destination count mismatch")
	}
	for i := range dest {
		if err := assignScanValue(dest[i], r.values[i]); err != nil {
			return err
		}
	}
	return nil
}

type stubRows struct {
	rows [][]any
	idx  int
	err  error
}

func (r *stubRows) Close() {}

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
		if err := assignScanValue(dest[i], row[i]); err != nil {
			return err
		}
	}
	return nil
}

func (r *stubRows) Values() ([]any, error) { return nil, errors.New("not implemented") }

func (r *stubRows) RawValues() [][]byte { return nil }

func (r *stubRows) Conn() *pgx.Conn { return nil }

type stubDB struct {
	exportIncomesRows  [][]any
	exportExpensesRows [][]any
}

func (d *stubDB) Exec(context.Context, string, ...interface{}) (pgconn.CommandTag, error) {
	return pgconn.CommandTag{}, nil
}

func (d *stubDB) Query(_ context.Context, sql string, _ ...interface{}) (pgx.Rows, error) {
	if strings.HasPrefix(sql, "-- name: ExportIncomes :many") {
		return &stubRows{rows: d.exportIncomesRows}, nil
	}
	if strings.HasPrefix(sql, "-- name: ExportExpenses :many") {
		return &stubRows{rows: d.exportExpensesRows}, nil
	}
	return &stubRows{err: errors.New("query not configured")}, nil
}

func (d *stubDB) QueryRow(context.Context, string, ...interface{}) pgx.Row {
	return &stubRow{err: errors.New("query row not configured")}
}

type stubTxBeginner struct {
	tx         *stubTx
	beginErr   error
	beginCalls int
}

func (b *stubTxBeginner) BeginTx(_ context.Context, _ pgx.TxOptions) (pgx.Tx, error) {
	b.beginCalls++
	if b.beginErr != nil {
		return nil, b.beginErr
	}
	return b.tx, nil
}

type stubTx struct {
	commitCalled   bool
	rollbackCalled bool

	createIncomeCalls  int
	createExpenseCalls int

	createExpenseErr error

	listCategoriesRows [][]any
}

func (t *stubTx) Begin(context.Context) (pgx.Tx, error) { return nil, errors.New("not implemented") }

func (t *stubTx) Commit(context.Context) error {
	t.commitCalled = true
	return nil
}

func (t *stubTx) Rollback(context.Context) error {
	t.rollbackCalled = true
	return nil
}

func (t *stubTx) CopyFrom(context.Context, pgx.Identifier, []string, pgx.CopyFromSource) (int64, error) {
	return 0, errors.New("not implemented")
}

func (t *stubTx) SendBatch(context.Context, *pgx.Batch) pgx.BatchResults { return nil }

func (t *stubTx) LargeObjects() pgx.LargeObjects { return pgx.LargeObjects{} }

func (t *stubTx) Prepare(context.Context, string, string) (*pgconn.StatementDescription, error) {
	return nil, errors.New("not implemented")
}

func (t *stubTx) Exec(context.Context, string, ...any) (pgconn.CommandTag, error) {
	return pgconn.CommandTag{}, nil
}

func (t *stubTx) Query(_ context.Context, sql string, _ ...any) (pgx.Rows, error) {
	if strings.HasPrefix(sql, "-- name: ListCategories :many") {
		return &stubRows{rows: t.listCategoriesRows}, nil
	}
	return &stubRows{err: errors.New("tx query not configured")}, nil
}

func (t *stubTx) QueryRow(_ context.Context, sql string, args ...any) pgx.Row {
	now := pgtype.Timestamp{Time: time.Now().UTC(), Valid: true}

	if strings.HasPrefix(sql, "-- name: CreateIncome :one") {
		t.createIncomeCalls++
		return &stubRow{values: []any{
			uuid.New(),
			args[7].(uuid.UUID),
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
		t.createExpenseCalls++
		if t.createExpenseErr != nil {
			return &stubRow{err: t.createExpenseErr}
		}
		return &stubRow{values: []any{
			uuid.New(),
			args[0].(string),
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

	return &stubRow{err: errors.New("tx query row not configured")}
}

func (t *stubTx) Conn() *pgx.Conn { return nil }

func setAuthenticatedUser(c echo.Context, userID uuid.UUID) {
	c.Set("user", &auth.UserContext{ID: userID, Role: auth.RoleUser, Email: "user@example.com"})
}

func TestExportBudgetJSONReturnsExpectedStructure(t *testing.T) {
	e := echo.New()
	db := &stubDB{}
	queries := sqlc.New(db)
	logger := zerolog.Nop()
	h := NewHandler(queries, &logger)

	userID := uuid.New()
	categoryID := uuid.New()
	priorityID := uuid.New()
	incomeID := uuid.New()
	expenseID := uuid.New()

	db.exportIncomesRows = [][]any{{
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

	db.exportExpensesRows = [][]any{{
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
	assert.Equal(t, expenseID.String(), payload.Expenses[0].ID)
	assert.Equal(t, priorityID.String(), *payload.Expenses[0].PriorityGroupID)
}

func TestImportBudgetJSONRollsBackOnMidImportFailure(t *testing.T) {
	e := echo.New()
	db := &stubDB{}
	queries := sqlc.New(db)
	logger := zerolog.Nop()

	userID := uuid.New()
	categoryID := uuid.New()

	tx := &stubTx{
		createExpenseErr: errors.New("insert failed"),
		listCategoriesRows: [][]any{{
			categoryID,
			"Food",
			pgtype.Text{Valid: false},
			pgtype.Text{Valid: false},
			pgtype.Timestamp{Valid: true},
			userID,
		}},
	}
	txBeginner := &stubTxBeginner{tx: tx}
	h := NewHandler(queries, &logger, txBeginner)

	// no existing records, so both incoming records are CREATE actions
	db.exportIncomesRows = [][]any{}
	db.exportExpensesRows = [][]any{}

	bodyPayload := BudgetExportPayload{
		Incomes: []ImportIncomeRecord{
			{Date: "2026-04-11", Description: "Side Gig", Amount: 300, Currency: "USD"},
		},
		Expenses: []ImportExpenseRecord{
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
	require.Equal(t, http.StatusInternalServerError, rec.Code)

	assert.Equal(t, 1, txBeginner.beginCalls)
	assert.Equal(t, 1, tx.createIncomeCalls)
	assert.Equal(t, 1, tx.createExpenseCalls)
	assert.False(t, tx.commitCalled)
	assert.True(t, tx.rollbackCalled)
}

func TestImportBudgetJSONInvalidPayloadReturns400(t *testing.T) {
	e := echo.New()
	db := &stubDB{}
	queries := sqlc.New(db)
	logger := zerolog.Nop()

	tx := &stubTx{}
	txBeginner := &stubTxBeginner{tx: tx}
	h := NewHandler(queries, &logger, txBeginner)

	bodyPayload := BudgetExportPayload{
		Incomes: []ImportIncomeRecord{{
			Date:     "invalid-date",
			Amount:   100,
			Currency: "USD",
		}},
	}
	bodyBytes, err := json.Marshal(bodyPayload)
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/api/budget/import", bytes.NewReader(bodyBytes))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	setAuthenticatedUser(c, uuid.New())

	err = h.ImportBudgetJSON(c)
	require.NoError(t, err)
	require.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Contains(t, rec.Body.String(), "incomes[0].date")

	assert.Equal(t, 0, txBeginner.beginCalls)
	assert.False(t, tx.commitCalled)
	assert.False(t, tx.rollbackCalled)
}
