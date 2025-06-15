package port

import (
	"context"

	"github.com/mfajri11/xyz-backend-monolith/app/core/domain"
)

type LoanRepository interface {
	CreateLoan(ctx context.Context, loan domain.Loan) error
	GetLoanByUserIDAndContractNumber(ctx context.Context, uid int64, contractNumber string) (*domain.LoanAll, error)
	CreateLoanPayment(ctx context.Context, loanPayment domain.LoanPayment) error
	GetLoanPaymentsByLoanID(ctx context.Context, loanID int64) ([]domain.LoanPayment, error)
	GetLoanPaymentsByUserIDAndContractNumber(ctx context.Context, uid int64, contractNumber string) ([]domain.LoanPayment, error)
}

type LoanService interface {
	CreateLoan(ctx context.Context, loan domain.Loan) error
	CalculateOTRAmount(amount float64) float64
	GetLoanByUserIDAndContractNumber(ctx context.Context, uid int64, contractNumber string) (*domain.LoanAll, error)
	CreateLoanPayment(ctx context.Context, loanPayment domain.LoanPayment) error
	GetLoanPaymentsByUserIDAndContractNumber(ctx context.Context, uid int64, contractNumber string) ([]domain.LoanPayment, error)
}
