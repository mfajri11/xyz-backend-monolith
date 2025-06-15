package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"

	"github.com/mfajri11/xyz-backend-monolith/app/core/domain"
	"github.com/mfajri11/xyz-backend-monolith/util/apperror"
	uhttp "github.com/mfajri11/xyz-backend-monolith/util/http"
)

type UserRepository struct {
	dbConn    *sql.DB
	kycClient *uhttp.HTTPClient
}

func New(db *sql.DB, kycClient *uhttp.HTTPClient) *UserRepository {
	return &UserRepository{
		dbConn:    db,
		kycClient: kycClient,
	}
}

func (repo *UserRepository) FindOneByNationalID(ctx context.Context, nid string) (user *domain.UserEntity, err error) {
	user = new(domain.UserEntity)
	err = repo.dbConn.QueryRowContext(ctx, getUserByNationalID, nid).Scan(
		&user.ID,
		&user.NationalID,
		&user.FullName,
		&user.LegalName,
		&user.IsNationalIDValidated,
		&user.IsPhotoValidated,
	)
	if err == sql.ErrNoRows {
		err = fmt.Errorf("FindOneByNationalID: user with national id %s not found", nid)
		return nil, apperror.WrapError(err, apperror.ErrNotFound)
	}
	if err != nil {
		err = fmt.Errorf("FindOneByNationalID: error select query: %w", err)
		return nil, apperror.WrapError(err, apperror.ErrInternalServerError)
	}

	return user, nil
}

func (repo *UserRepository) UpdateByID(ctx context.Context, user domain.UserEntity) error {
	_, err := repo.dbConn.ExecContext(ctx, queryUpdateUserById)
	if err != nil {
		err = fmt.Errorf("UpdateByID: error update user: %w", err)
		return apperror.WrapError(err, apperror.ErrInternalServerError)
	}

	return nil
}

func (repo *UserRepository) ValidateSalary(ctx context.Context, req domain.KYCValidateSalaryReq) (*domain.KYCValidateSalaryResp, error) {
	resp, err := repo.kycClient.Post(repo.kycClient.BaseURL+"/veryfi/national-id", req)
	if err != nil {
		err = fmt.Errorf("ValidateSalary: error kyc request: %w", err)
		return nil, apperror.WrapError(err, apperror.ErrInternalServerError)
	}

	if resp.StatusCode != 200 {
		err = fmt.Errorf("ValidateSalary: unexpected status code: %d", resp.StatusCode)
		return nil, apperror.WrapError(err, apperror.ErrBadRequest)
	}

	defer resp.Body.Close()
	var kycData domain.KYCValidateSalaryResp
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		err = fmt.Errorf("ValidateSalary: error read kyc body response: %w", err)
		return nil, err
	}

	if err = json.Unmarshal(b, &kycData); err != nil {
		err = fmt.Errorf("ValidateSalary: error unmarshalling kyc data: %w", err)
		return nil, apperror.WrapError(err, apperror.ErrInternalServerError)
	}

	return &kycData, nil
}

func (repo *UserRepository) ValidateNationalID(ctx context.Context, req domain.KYCValidateNationalIDReq) (*domain.KYCValidateNationalIDResp, error) {
	resp, err := repo.kycClient.Post(repo.kycClient.BaseURL+"/veryfi/national-id", req)
	if err != nil {
		err = fmt.Errorf("VerifyNationalID: error kyc request: %w", err)
		return nil, apperror.WrapError(err, apperror.ErrInternalServerError)
	}

	if resp.StatusCode != 200 {
		err = fmt.Errorf("VerifyNationalID: unexpected status code: %d", resp.StatusCode)
		return nil, apperror.WrapError(err, apperror.ErrBadRequest)
	}

	defer resp.Body.Close()
	var kycData domain.KYCValidateNationalIDResp
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		err = fmt.Errorf("VerifyNationalID: error read kyc body response: %w", err)
		return nil, apperror.WrapError(err, apperror.ErrInternalServerError)
	}

	if err = json.Unmarshal(b, &kycData); err != nil {
		err = fmt.Errorf("VerifyNationalID: error unmarshalling kyc data: %w", err)
		return nil, apperror.WrapError(err, apperror.ErrInternalServerError)
	}

	return &kycData, nil
}

func (repo *UserRepository) ValidatePhoto(ctx context.Context, req domain.KYCValidatePhotoReq) (*domain.KYCValidatePhotoResp, error) {
	resp, err := repo.kycClient.Post(repo.kycClient.BaseURL+"/veryfi/national-id", req)
	if err != nil {
		return nil, fmt.Errorf("VerifyPhoto: error kyc request: %w", err)
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("VerifyPhoto: unexpected status code: %d", resp.StatusCode)
	}

	defer resp.Body.Close()
	var kycData domain.KYCValidatePhotoResp
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("VerifyPhoto: error read kyc body response: %w", err)
	}

	if err = json.Unmarshal(b, &kycData); err != nil {
		return nil, fmt.Errorf("VerifyPhoto: error unmarshalling kyc data: %w", err)
	}

	return &kycData, nil
}
