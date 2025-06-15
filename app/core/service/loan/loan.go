package loan

import (
	"context"
	"fmt"

	"github.com/mfajri11/xyz-backend-monolith/app/core/domain"
	"github.com/mfajri11/xyz-backend-monolith/app/core/port"
)

type LoanService struct {
	repo port.LoanRepository
}

func New(repo port.LoanRepository) *LoanService {
	return &LoanService{
		repo: repo,
	}
}

func (svc *LoanService) CreateLoan(ctx context.Context, loan domain.Loan) error {
	err := svc.repo.CreateLoan(ctx, domain.Loan{
		UserID:          loan.UserID,
		ContractNumber:  loan.ContractNumber,
		OTRAmount:       loan.OTRAmount,
		PrincipalAmount: loan.PrincipalAmount,
		AssetName:       loan.AssetName,
		LoanTypeID:      loan.LoanTypeID,
		LimitTypeID:     loan.LimitTypeID,
		Status:          loan.Status,
		StartDate:       loan.StartDate,
		InterestRate:    loan.InterestRate,
	})

	if err != nil {
		return fmt.Errorf("CreateLoan: error insert loan: %w", err)
	}

	return nil
}

func (svc *LoanService) CreateLoanPayment(ctx context.Context, loanPayment domain.LoanPayment) error {
	err := svc.repo.CreateLoanPayment(ctx, domain.LoanPayment{
		LoanID:  loanPayment.LoanID,
		Amount:  loanPayment.Amount,
		Date:    loanPayment.Date,
		Channel: loanPayment.Channel,
	})

	if err != nil {
		return fmt.Errorf("CreateLoanPayment: error insert loan payment: %w", err)
	}

	return nil
}

func (svc *LoanService) GetLoanByUserIDAndContractNumber(ctx context.Context, uid int64, contractNumber string) (*domain.LoanAll, error) {
	return svc.repo.GetLoanByUserIDAndContractNumber(ctx, uid, contractNumber)
}

func (svc *LoanService) GetLoanPaymentsByUserIDAndContractNumber(ctx context.Context, uid int64, contractNumber string) ([]domain.LoanPayment, error) {
	return svc.repo.GetLoanPaymentsByUserIDAndContractNumber(ctx, uid, contractNumber)
}

func (svc *LoanService) CalculateOTRAmount(amount float64) float64 {
	// not yet implemented
	return amount + amount*0.3
}
