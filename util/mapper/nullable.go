package mapper

import (
	"database/sql"
	"time"
)

func NewSQLNullableFloat64(f float64) sql.NullFloat64 {
	return sql.NullFloat64{
		Float64: f,
		Valid:   true,
	}
}

func NewSQLNUllableTime(s string) (sql.NullTime, error) {
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		return sql.NullTime{}, err
	}
	return sql.NullTime{
		Time:  t,
		Valid: true,
	}, err
}
