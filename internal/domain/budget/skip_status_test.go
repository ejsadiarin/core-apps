package budget

import (
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
	sharedvalidator "core-gateway/internal/shared/validator"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type skipStatusDB struct {
	lastSQL  string
	lastArgs []any

	incomeRuleID  uuid.UUID
	expenseRuleID uuid.UUID
}

type skipStatusRow struct {
	values []any
	err    error
}

func (r *skipStatusRow) Scan(dest ...any) error {
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
		case *bool:
			*d = r.values[i].(bool)
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

func (d *skipStatusDB) Exec(context.Context, string, ...interface{}) (pgconn.CommandTag, error) {
	return pgconn.CommandTag{}, nil
}

func (d *skipStatusDB) Query(context.Context, string, ...interface{}) (pgx.Rows, error) {
	return &stubRows{err: errors.New("query not configured")}, nil
}

func (d *skipStatusDB) QueryRow(_ context.Context, sql string, args ...interface{}) pgx.Row {
	d.lastSQL = sql
	d.lastArgs = args
	now := pgtype.Timestamp{Time: time.Now().UTC(), Valid: true}

	if strings.HasPrefix(sql, "-- name: UpsertSkippedIncome :one") {
		return &skipStatusRow{values: []any{
			uuid.New(),
			args[0].(uuid.UUID),
			float64ToNumeric(0),
			pgtype.Text{String: "USD", Valid: true},
			args[1].(pgtype.Date),
			pgtype.Text{String: "rule skip", Valid: true},
			pgtype.Text{String: "monthly", Valid: true},
			pgtype.Date{Valid: false},
			now,
			now,
			pgtype.Date{Valid: false},
			pgtype.Bool{Bool: false, Valid: true},
			"skipped",
			pgtype.UUID{Bytes: d.incomeRuleID, Valid: true},
		}}
	}

	if strings.HasPrefix(sql, "-- name: UpsertSkippedExpense :one") {
		return &skipStatusRow{values: []any{
			uuid.New(),
			"Skipped recurring expense",
			float64ToNumeric(0),
			pgtype.Text{String: "USD", Valid: true},
			pgtype.UUID{Valid: false},
			args[1].(pgtype.Date),
			now,
			now,
			pgtype.Text{Valid: false},
			args[0].(uuid.UUID),
			pgtype.Text{String: "monthly", Valid: true},
			pgtype.Date{Valid: false},
			pgtype.Date{Valid: false},
			pgtype.UUID{Valid: false},
			pgtype.Bool{Bool: false, Valid: true},
			"skipped",
			pgtype.UUID{Bytes: d.expenseRuleID, Valid: true},
		}}
	}

	if strings.HasPrefix(sql, "-- name: CheckSkippedExpense :one") {
		return &skipStatusRow{values: []any{true}}
	}

	return &skipStatusRow{err: errors.New("query row not configured")}
}

func TestSkipStatusIncomeHandlerUsesStatusRows(t *testing.T) {
	e := echo.New()
	e.Validator = sharedvalidator.NewValidator()

	ruleID := uuid.New()
	userID := uuid.New()
	db := &skipStatusDB{incomeRuleID: ruleID}
	queries := sqlc.New(db)
	logger := zerolog.Nop()
	h := NewHandler(queries, &logger)

	body := `{"date":"2026-06-15","source_rule_id":"` + ruleID.String() + `"}`
	req := httptest.NewRequest(http.MethodPost, "/api/budget/incomes/skip", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("user", &auth.UserContext{ID: userID, Role: auth.RoleUser, Email: "user@example.com"})

	err := h.SkipIncome(c)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, db.lastSQL, "status = 'skipped'")
	assert.NotContains(t, db.lastSQL, "Skipped:")

	var payload map[string]any
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &payload))
	assert.Equal(t, "skipped", payload["status"])
}

func TestSkipStatusExpenseHandlerUsesStatusRows(t *testing.T) {
	e := echo.New()
	e.Validator = sharedvalidator.NewValidator()

	ruleID := uuid.New()
	userID := uuid.New()
	db := &skipStatusDB{expenseRuleID: ruleID}
	queries := sqlc.New(db)
	logger := zerolog.Nop()
	h := NewHandler(queries, &logger)

	body := `{"expense_date":"2026-06-16","source_rule_id":"` + ruleID.String() + `"}`
	req := httptest.NewRequest(http.MethodPost, "/api/budget/expenses/skip", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("user", &auth.UserContext{ID: userID, Role: auth.RoleUser, Email: "user@example.com"})

	err := h.SkipExpense(c)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, db.lastSQL, "status = 'skipped'")
	assert.NotContains(t, db.lastSQL, "Skipped:")

	var payload map[string]any
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &payload))
	assert.Equal(t, "2026-06-16", payload["expense_date"])
}
