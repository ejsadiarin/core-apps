package budget

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

// type conversion helpers

func float64ToNumeric(f float64) pgtype.Numeric {
	n := pgtype.Numeric{}
	n.Scan(fmt.Sprintf("%f", f))
	return n
}

func numericToFloat64(n pgtype.Numeric) float64 {
	f, _ := n.Float64Value()
	return f.Float64
}

func interfaceToFloat64(i interface{}) float64 {
	if n, ok := i.(pgtype.Numeric); ok {
		return numericToFloat64(n)
	}
	if v, ok := i.(int64); ok {
		return float64(v)
	}
	return 0.0
}

func stringToDate(s string) pgtype.Date {
	d := pgtype.Date{}
	d.Scan(s)
	return d
}

func dateToString(d pgtype.Date) string {
	return d.Time.Format("2006-01-02")
}

func stringPtrToText(s *string) pgtype.Text {
	if s == nil {
		return pgtype.Text{Valid: false}
	}
	return pgtype.Text{String: *s, Valid: true}
}

func textToStringPtr(t pgtype.Text) *string {
	if !t.Valid {
		return nil
	}
	s := t.String
	return &s
}

func uuidPtrToNullUUID(u *uuid.UUID) pgtype.UUID {
	if u == nil {
		return pgtype.UUID{Valid: false}
	}
	return pgtype.UUID{Bytes: *u, Valid: true}
}

func stringPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

func stringPtrToDate(s *string) pgtype.Date {
	if s == nil {
		return pgtype.Date{Valid: false}
	}
	d := pgtype.Date{}
	d.Scan(*s)
	return d
}

func getCurrency(t pgtype.Text) string {
	if !t.Valid || t.String == "" {
		return "USD"
	}
	return t.String
}

func dateToNullableStringPtr(d pgtype.Date) *string {
	if !d.Valid {
		return nil
	}
	s := d.Time.Format("2006-01-02")
	return &s
}

func float64PtrToNumeric(f *float64) pgtype.Numeric {
	if f == nil {
		return pgtype.Numeric{Valid: false}
	}
	return float64ToNumeric(*f)
}
