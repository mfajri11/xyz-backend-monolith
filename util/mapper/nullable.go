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
	t, err := time.Parse(time.DateOnly, s)
	if err != nil {
		return sql.NullTime{}, err
	}
	return sql.NullTime{
		Time:  t,
		Valid: true,
	}, err
}

// for testing purpose
func MustNewSQLNUllableTime(s string) sql.NullTime {
	t, err := time.Parse(time.DateOnly, s)
	if err != nil {
		panic(err)
	}
	return sql.NullTime{
		Time:  t,
		Valid: true,
	}
}

func NewSQLNUllableInt16(v int16) sql.NullInt16 {
	return sql.NullInt16{
		Int16: v,
		Valid: true,
	}
}

func NewSQLNUllableString(s string) sql.NullString {
	return sql.NullString{
		String: s,
		Valid:  true,
	}
}
