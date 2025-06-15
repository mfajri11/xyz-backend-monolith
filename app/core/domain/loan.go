package domain

import (
	"database/sql"
	"database/sql/driver"
	"time"
)

type LoanStatus string
type LoanTypeName string

const (
	ACIIVE   LoanStatus = "ACTIVE"
	INCATIVE LoanStatus = "INACTIVE"
	REJECTED LoanStatus = "REJECTED"
)

const (
	CAR        LoanTypeName = "CAR"
	BIKE       LoanTypeName = "BIKE"
	WHITEGOODS LoanTypeName = "WHITE_GOODS"
)

type Loan struct {
	ID              int64
	UserID          int64
	ContractNumber  string
	OTRAmount       float64
	PrincipalAmount float64
	AssetName       string
	LoanTypeID      sql.NullInt16
	LimitTypeID     sql.NullInt16
	Status          sql.NullString
	StartDate       sql.NullTime
	InterestRate    sql.NullFloat64
}

type LoanAll struct {
	ID              int64
	UserID          int64
	ContractNumber  string
	OTRAmount       float64
	PrincipalAmount float64
	AssetName       string
	LoanType        LoanType
	LimitType       LimitType
	Status          sql.NullString
	StartDate       sql.NullTime
	InterestRate    sql.NullFloat64
}

type LoanPayment struct {
	ID      int64
	LoanID  int64
	Amount  float64
	Date    time.Time
	Channel string
}

type LimitType struct {
	ID     int8
	Amount sql.NullFloat64
	Term   int8
}

type LoanType struct {
	ID   int8
	Name LoanTypeName
}

// Scan for LoanStatus (implements sql.Scanner interface)
func (s *LoanStatus) Scan(value interface{}) error {
	if value == nil {
		*s = "" // Null value
		return nil
	}
	*s = LoanStatus(value.(string))
	return nil
}

// Value for LoanStatus (implements driver.Valuer interface)
func (s LoanStatus) Value() (driver.Value, error) {
	return string(s), nil
}

// Scan for LoanTypeName (implements sql.Scanner interface)
func (t *LoanTypeName) Scan(value interface{}) error {
	if value == nil {
		*t = "" // Null value
		return nil
	}
	*t = LoanTypeName(value.(string))
	return nil
}

// Value for LoanTypeName (implements driver.Valuer interface)
func (t LoanTypeName) Value() (driver.Value, error) {
	return string(t), nil
}
