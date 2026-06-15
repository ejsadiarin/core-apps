# TODO: Fix Budget RecurringType Type Mismatch

## Issue

sqlc generates `string` for BudgetExpense.RecurringType and BudgetIncome.RecurringType, but handler code still uses pgtype.Text conversion functions.

## Errors

```
handler.go:449: cannot use stringPtrToText(req.RecurringType) as string value in struct literal
handler.go:591: cannot use row.RecurringType (variable of type string) as pgtype.Text value
```

## Fix Required

1. For CREATE inputs (req.RecurringType) - use `stringPtrToString()` instead of `stringPtrToText()`
2. For READ outputs (row.RecurringType) - use `recurringTypeToStringPtr()` or remove `.String` from pgtype.Text patterns

## Files to Update

- handler.go
- handler_income.go
- handler_import_export.go
- handler_stats.go (replace .RecurringType.String with .RecurringType)

## After Fix

Run: `go test ./...`