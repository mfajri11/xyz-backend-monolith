package loan

import (
	"context"
	"database/sql"
	"errors"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/mfajri11/xyz-backend-monolith/app/core/domain"
	"github.com/mfajri11/xyz-backend-monolith/util/mapper"
	"github.com/stretchr/testify/assert"
)

func TestLoanRepositories_CreateLoan(t *testing.T) {
	type args struct {
		ctx  context.Context
		loan domain.Loan
	}
	type mock struct {
		sqlmock.Sqlmock
	}
	tests := []struct {
		name        string
		repo        *LoanRepositories
		args        args
		prepareMock func(mock *mock)
		wantErr     bool
	}{
		{
			name: "Given a valid loan, it should return no error",
			repo: &LoanRepositories{},
			args: args{
				ctx: context.Background(),
				loan: domain.Loan{
					UserID:          1,
					ContractNumber:  "12345678912345678",
					OTRAmount:       1000,
					PrincipalAmount: 1000,
					AssetName:       "test",
					LoanTypeID:      mapper.NewSQLNUllableInt16(1),
					LimitTypeID:     mapper.NewSQLNUllableInt16(1),
					Status:          mapper.NewSQLNUllableString("active"),
					StartDate:       mapper.MustNewSQLNUllableTime("2021-01-01"),
					InterestRate:    mapper.NewSQLNullableFloat64(0.01),
				},
			},
			prepareMock: func(mock *mock) {
				mock.ExpectExec(regexp.QuoteMeta(createLoan)).WithArgs(1, "12345678912345678", float64(1000), float64(1000), "test", 1, 1, "active", time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC), 0.01).WillReturnResult(sqlmock.NewResult(1, 1))
			},
		},
		{
			name: "Given a valid loan, but insert it return error",
			repo: &LoanRepositories{},
			args: args{
				ctx: context.Background(),
				loan: domain.Loan{
					UserID:          1,
					ContractNumber:  "12345678912345678",
					OTRAmount:       1000,
					PrincipalAmount: 1000,
					AssetName:       "test",
					LoanTypeID:      mapper.NewSQLNUllableInt16(1),
					LimitTypeID:     mapper.NewSQLNUllableInt16(1),
					Status:          mapper.NewSQLNUllableString("active"),
					StartDate:       mapper.MustNewSQLNUllableTime("2021-01-01"),
					InterestRate:    mapper.NewSQLNullableFloat64(0.01),
				},
			},
			prepareMock: func(mock *mock) {
				mock.ExpectExec(regexp.QuoteMeta(createLoan)).WithArgs(1, "12345678912345678", float64(1000), float64(1000), "test", 1, 1, "active", time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC), 0.01).WillReturnResult(nil).WillReturnError(sql.ErrConnDone)
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			conn, sqlMock, err := sqlmock.New()
			assert.NoError(t, err)
			defer conn.Close()
			tt.repo.dbConn = conn
			tt.prepareMock(&mock{
				Sqlmock: sqlMock,
			})

			err = tt.repo.CreateLoan(tt.args.ctx, tt.args.loan)
			assert.Equal(t, tt.wantErr, err != nil, err)
		})
	}
}

func TestLoanRepositories_CreateLoanPayment(t *testing.T) {
	type args struct {
		ctx         context.Context
		loanPayment domain.LoanPayment
	}
	type mock struct {
		sqlmock.Sqlmock
	}
	tests := []struct {
		name        string
		repo        *LoanRepositories
		args        args
		prepareMock func(mock *mock)
		wantErr     bool
	}{
		{
			name: "Given a valid loan payment, it should return no error",
			repo: &LoanRepositories{},
			args: args{
				ctx: context.Background(),
				loanPayment: domain.LoanPayment{
					Amount:  float64(1000),
					LoanID:  1,
					Date:    time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
					Channel: "test",
				},
			},
			prepareMock: func(mock *mock) {
				mock.ExpectExec(regexp.QuoteMeta(createLoanPayment)).WithArgs(1, float64(1000), time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC), "test").WillReturnResult(sqlmock.NewResult(1, 1))
			},
		},
		{
			name: "Given a valid loan payment, but insert it return error",
			repo: &LoanRepositories{},
			args: args{
				ctx: context.Background(),
				loanPayment: domain.LoanPayment{
					Amount:  float64(1000),
					LoanID:  1,
					Date:    time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
					Channel: "test",
				},
			},
			prepareMock: func(mock *mock) {
				mock.ExpectExec(regexp.QuoteMeta(createLoanPayment)).WithArgs(1, float64(1000), time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC), "test").WillReturnResult(nil).WillReturnError(sql.ErrConnDone)
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			conn, sqlMock, err := sqlmock.New()
			assert.NoError(t, err)
			defer conn.Close()
			tt.repo.dbConn = conn
			tt.prepareMock(&mock{
				Sqlmock: sqlMock,
			})

			err = tt.repo.CreateLoanPayment(tt.args.ctx, tt.args.loanPayment)
			assert.Equal(t, tt.wantErr, err != nil, err)
		})
	}
}

func TestLoanRepositories_GetLoanPaymentsByLoanID(t *testing.T) {
	type args struct {
		ctx    context.Context
		loanID int64
	}
	type mock struct {
		sqlmock.Sqlmock
	}
	tests := []struct {
		name        string
		repo        *LoanRepositories
		args        args
		prepareMock func(m *mock)
		want        []domain.LoanPayment
		wantErr     bool
	}{
		{
			name: "Given a valid loan id, it should return a list of loan payments",
			repo: &LoanRepositories{},
			args: args{
				ctx:    context.Background(),
				loanID: 1,
			},
			prepareMock: func(m *mock) {
				m.ExpectQuery(regexp.QuoteMeta(getLoanPaymentsByLoanID)).WithArgs(1).WillReturnRows(sqlmock.NewRows([]string{"amount", "date", "channel"}).AddRow(float64(1000), time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC), "test").AddRow(float64(4000), time.Date(2021, 1, 2, 0, 0, 0, 0, time.UTC), "test-2"))
			},
			want: []domain.LoanPayment{
				{
					Amount:  float64(1000),
					Date:    time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
					Channel: "test",
				},
				{
					Amount:  float64(4000),
					Date:    time.Date(2021, 1, 2, 0, 0, 0, 0, time.UTC),
					Channel: "test-2",
				},
			},
		},
		{
			name: "Given a invalid loan id, it should return no list",
			repo: &LoanRepositories{},
			args: args{
				ctx:    context.Background(),
				loanID: -1,
			},
			prepareMock: func(m *mock) {
				m.ExpectQuery(regexp.QuoteMeta(getLoanPaymentsByLoanID)).WithArgs(-1).WillReturnRows(sqlmock.NewRows([]string{"amount", "date", "channel"}))
			},
		},
		{
			name: "Given a invalid loan id, it should return error",
			repo: &LoanRepositories{},
			args: args{
				ctx:    context.Background(),
				loanID: 0,
			},
			prepareMock: func(m *mock) {
				m.ExpectQuery(regexp.QuoteMeta(getLoanPaymentsByLoanID)).WithArgs(0).WillReturnRows(sqlmock.NewRows([]string{"amount", "date", "channel"})).WillReturnError(errors.New("oops!"))
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			conn, sqlMock, err := sqlmock.New()
			assert.NoError(t, err)
			defer conn.Close()
			tt.repo.dbConn = conn
			tt.prepareMock(&mock{
				Sqlmock: sqlMock,
			})

			got, err := tt.repo.GetLoanPaymentsByLoanID(tt.args.ctx, tt.args.loanID)

			assert.Equal(t, tt.wantErr, err != nil, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestLoanRepositories_GetLoanPaymentsByUserIDAndContractNumber(t *testing.T) {
	type args struct {
		ctx            context.Context
		userID         int64
		contractNumber string
	}
	type mock struct {
		sqlmock.Sqlmock
	}
	tests := []struct {
		name        string
		repo        *LoanRepositories
		args        args
		prepareMock func(m *mock)
		want        []domain.LoanPayment
		wantErr     bool
	}{
		{
			name: "Given a valid loan id, it should return a list of loan payments",
			repo: &LoanRepositories{},
			args: args{
				ctx:            context.Background(),
				contractNumber: "XYZ-LAI-01",
			},
			prepareMock: func(m *mock) {
				m.ExpectQuery(regexp.QuoteMeta(getLoanByContractNumber)).WithArgs(1).WillReturnRows(sqlmock.NewRows([]string{"amount", "date", "channel"}).AddRow(float64(1000), time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC), "test").AddRow(float64(4000), time.Date(2021, 1, 2, 0, 0, 0, 0, time.UTC), "test-2"))
			},
			want: []domain.LoanPayment{
				{
					Amount:  float64(1000),
					Date:    time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
					Channel: "test",
				},
				{
					Amount:  float64(4000),
					Date:    time.Date(2021, 1, 2, 0, 0, 0, 0, time.UTC),
					Channel: "test-2",
				},
			},
		},
		{
			name: "Given a invalid loan id, it should return no list",
			repo: &LoanRepositories{},
			args: args{
				ctx:            context.Background(),
				contractNumber: "XYZ-LAI--1",
			},
			prepareMock: func(m *mock) {
				m.ExpectQuery(regexp.QuoteMeta(getLoanByContractNumber)).WithArgs(-1).WillReturnRows(sqlmock.NewRows([]string{"amount", "date", "channel"}))
			},
		},
		{
			name: "Given a invalid loan id, it should return error",
			repo: &LoanRepositories{},
			args: args{
				ctx:            context.Background(),
				contractNumber: "XYZ-LAI-0",
			},
			prepareMock: func(m *mock) {
				m.ExpectQuery(regexp.QuoteMeta(getLoanByContractNumber)).WithArgs(0).WillReturnRows(sqlmock.NewRows([]string{"amount", "date", "channel"})).WillReturnError(errors.New("oops!"))
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			conn, sqlMock, err := sqlmock.New()
			assert.NoError(t, err)
			defer conn.Close()
			tt.repo.dbConn = conn
			tt.prepareMock(&mock{
				Sqlmock: sqlMock,
			})

			got, err := tt.repo.GetLoanByUserIDAndContractNumber(tt.args.ctx, tt.args.userID, tt.args.contractNumber)

			assert.Equal(t, tt.wantErr, err != nil, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
