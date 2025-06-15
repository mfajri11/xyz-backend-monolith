package domain

import (
	"database/sql"
	"time"
)

type UserEntity struct {
	ID                    int
	NationalID            string
	FullName              string
	LegalName             string
	BirthOfPlace          sql.NullString
	BirthOfDate           sql.NullTime
	Salary                sql.NullFloat64
	NationalIDPhoto       sql.NullByte
	UserPhoto             sql.NullByte
	IsNationalIDValidated bool
	IsPhotoValidated      bool
	ISSalaryValidated     bool
	CreatedAt             time.Time
	UpdatedAt             time.Time
	CreatedBy             string
	UpdatedBy             string
}
