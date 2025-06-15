package loan

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/mfajri11/xyz-backend-monolith/app/core/domain"
	"github.com/mfajri11/xyz-backend-monolith/util/apperror"
)

type LoanRepositories struct {
	dbConn *sql.DB
}

func New(db *sql.DB) *LoanRepositories {
	return &LoanRepositories{
		dbConn: db,
	}
}

func (repo *LoanRepositories) CreateLoan(ctx context.Context, loan domain.Loan) error {
	_, err := repo.dbConn.ExecContext(ctx, createLoan, loan.UserID, loan.ContractNumber, loan.OTRAmount, loan.PrincipalAmount, loan.AssetName,
		loan.LoanTypeID, loan.LimitTypeID, loan.Status, loan.StartDate, loan.InterestRate)
	if err != nil {
		err = fmt.Errorf("CreateLoan: error insert loan: %w", err)
		return apperror.WrapError(err, apperror.ErrInternalServerError)
	}

	return nil
}

func (repo *LoanRepositories) GetLoanByUserIDAndContractNumber(ctx context.Context, uid int64, contractNumber string) (*domain.LoanAll, error) {
	var loan domain.LoanAll
	err := repo.dbConn.QueryRowContext(context.Background(), getLoanByContractNumber, uid, contractNumber).
		Scan(
			&loan.ID,
			&loan.UserID,
			&loan.ContractNumber,
			&loan.OTRAmount,
			&loan.PrincipalAmount,
			&loan.AssetName,
			&loan.LoanType.Name,
			&loan.LimitType.Amount,
			&loan.LimitType.Term,
			&loan.Status,
			&loan.StartDate,
			&loan.InterestRate,
		)

	if err == sql.ErrNoRows {
		err = fmt.Errorf("GetLoanByContractNumber: loan with contract number %s not found", contractNumber)
		return nil, apperror.WrapError(err, apperror.ErrNotFound)
	}

	if err != nil {
		err = fmt.Errorf("GetLoanByContractNumber: error select query: %w", err)
		return nil, apperror.WrapError(err, apperror.ErrInternalServerError)
	}

	return &loan, nil

}

func (repo *LoanRepositories) CreateLoanPayment(ctx context.Context, loanPayment domain.LoanPayment) error {
	_, err := repo.dbConn.ExecContext(ctx, createLoanPayment, loanPayment.LoanID, loanPayment.Amount, loanPayment.Date, loanPayment.Channel)
	if err != nil {
		err = fmt.Errorf("createLoanPayment: error insert loan payment: %w", err)
		return apperror.WrapError(err, apperror.ErrInternalServerError)
	}

	return nil
}

func (repo *LoanRepositories) GetLoanPaymentsByLoanID(ctx context.Context, loanID int64) ([]domain.LoanPayment, error) {
	var loanPayments []domain.LoanPayment
	rows, err := repo.dbConn.QueryContext(ctx, getLoanPaymentsByLoanID, loanID)
	if err != nil {
		err = fmt.Errorf("GetLoanPaymentsByLoanID: error select query: %w", err)
		return nil, apperror.WrapError(err, apperror.ErrInternalServerError)
	}

	for rows.Next() {
		var loanPayment domain.LoanPayment
		if err := rows.Scan(&loanPayment.Amount, &loanPayment.Date, &loanPayment.Channel); err != nil {
			err = fmt.Errorf("GetLoanPaymentsByLoanID: error scan query: %w", err)
			return nil, apperror.WrapError(err, apperror.ErrInternalServerError)
		}
		loanPayments = append(loanPayments, loanPayment)
	}

	return loanPayments, nil
}

func (repo *LoanRepositories) GetLoanPaymentsByUserIDAndContractNumber(ctx context.Context, uid int64, contractNumber string) ([]domain.LoanPayment, error) {
	var loanPayments []domain.LoanPayment
	rows, err := repo.dbConn.QueryContext(ctx, getLoanPaymentsByContractNumber, uid, contractNumber)
	if err != nil {
		err = fmt.Errorf("getLoanPaymentsByContractNumber: error select query: %w", err)
		return nil, apperror.WrapError(err, apperror.ErrInternalServerError)
	}

	for rows.Next() {
		var loanPayment domain.LoanPayment
		if err := rows.Scan(&loanPayment.Amount, &loanPayment.Date, &loanPayment.Channel); err != nil {
			err = fmt.Errorf("getLoanPaymentsByContractNumber: error scan query: %w", err)
			return nil, apperror.WrapError(err, apperror.ErrInternalServerError)
		}
		loanPayments = append(loanPayments, loanPayment)
	}

	return loanPayments, nil
}
