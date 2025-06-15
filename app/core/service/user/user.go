package services

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/google/uuid"
	"github.com/mfajri11/xyz-backend-monolith/app/core/domain"
	"github.com/mfajri11/xyz-backend-monolith/app/core/port"
	"github.com/mfajri11/xyz-backend-monolith/util/apperror"
	"github.com/mfajri11/xyz-backend-monolith/util/mapper"
)

type UserService struct {
	repo port.UserRepository
}

func New(repo port.UserRepository) *UserService {
	return &UserService{
		repo: repo,
	}
}

func (svc *UserService) ValidateData(ctx context.Context, req domain.ValidateUserReq) (bool, error) {
	uid, ok := ctx.Value("uid").(int)
	if !ok {
		return false, apperror.WrapError(errors.New("invalid request"), apperror.ErrBadRequest)
	}

	refId, ok := ctx.Value("requestID").(string)
	if !ok {
		refId = uuid.New().String()
	}

	user, err := svc.repo.FindOneByNationalID(ctx, req.NationalID)
	if err != nil {
		return false, fmt.Errorf("ValidateData: error while find user: %w", err)
	}

	userToSave := domain.UserEntity{
		ID: uid,
	}
	if errors.Is(err, apperror.ErrNotFound) {

		if _, err := svc.validateNationalID(ctx, refId, req, &userToSave); err != nil {
			// already wrapped
			// TODO: might be wrapped on public function only?
			return false, err
		}

		if _, err := svc.validateSalary(ctx, refId, req, &userToSave); err != nil {
			return false, err
		}

		if _, err := svc.validatePhoto(ctx, refId, req, &userToSave); err != nil {
			return false, err
		}

		err = svc.repo.UpdateByID(ctx, userToSave)
		if err != nil {
			err = fmt.Errorf("ValidateData: error while update user: %w", err)
			return false, apperror.WrapError(err, apperror.ErrInternalServerError)
		}

		if _, err := svc.validatePhoto(ctx, refId, req, &userToSave); err != nil {
			return false, err
		}

		return true, nil
	}

	// check and re validate data
	if !user.IsNationalIDValidated {
		if _, err := svc.validateSalary(ctx, refId, req, user); err != nil {
			return false, err
		}
	}

	if !user.IsPhotoValidated {
		if _, err := svc.validatePhoto(ctx, refId, req, user); err != nil {
			return false, err
		}
	}

	return true, nil

}

func (svc *UserService) validateNationalID(ctx context.Context, refId string, req domain.ValidateUserReq, userToSave *domain.UserEntity) (bool, error) {
	validatedNID, err := svc.repo.ValidateNationalID(ctx, domain.KYCValidateNationalIDReq{
		NationalID:  req.NationalID,
		LegalName:   req.LegalName,
		DateOfBirth: req.BirthOfDate,
		ReferenceID: refId,
	})

	if err != nil {
		return false, fmt.Errorf("ValidateData: error while validate national id: %w", err)
	}
	if validatedNID.Data.NationalID {
		t, err := mapper.NewSQLNUllableTime(req.BirthOfDate)
		if err == nil {
			userToSave.BirthOfDate = t
		}
		userToSave.NationalID = req.NationalID
		userToSave.NationalID = req.NationalID
		userToSave.LegalName = req.LegalName
		userToSave.IsNationalIDValidated = true
	}

	return true, nil
}

func (svc *UserService) validateSalary(ctx context.Context, refId string, req domain.ValidateUserReq, userToSave *domain.UserEntity) (bool, error) {
	validatedSalary, err := svc.repo.ValidateSalary(ctx, domain.KYCValidateSalaryReq{
		NationalID: req.NationalID,
		LegalName:  req.LegalName,
		Salary:     req.Salary,
	})

	if err != nil {
		return false, fmt.Errorf("ValidateData: error while validate salary: %w", err)
	}

	userSalary, err := strconv.ParseFloat(req.Salary, 64)
	if err != nil {
		return false, apperror.WrapError(errors.New("salary isinvalid"), apperror.ErrBadRequest)
	}

	upperRangeSalary, err := strconv.ParseFloat(validatedSalary.Data.SalaryUper, 64)
	if err != nil {
		return false, apperror.WrapError(errors.New("error converting range upper salary from kyc response"), apperror.ErrInternalServerError)
	}

	lowerRangeSalary, err := strconv.ParseFloat(validatedSalary.Data.SalaryLower, 64)
	if err != nil {
		return false, apperror.WrapError(errors.New("error converting range lower salary from kyc response"), apperror.ErrInternalServerError)
	}

	if userSalary > lowerRangeSalary && userSalary < upperRangeSalary {
		userToSave.Salary = mapper.NewSQLNullableFloat64(userSalary)
		userToSave.ISSalaryValidated = true
	}

	return true, nil
}

func (svc *UserService) validatePhoto(ctx context.Context, refId string, req domain.ValidateUserReq, userToSave *domain.UserEntity) (bool, error) {
	validatedPhoto, err := svc.repo.ValidatePhoto(ctx, domain.KYCValidatePhotoReq{
		NationalID:  req.NationalID,
		ReferenceID: refId,
	})

	if err != nil {
		return false, fmt.Errorf("ValidateData: error while validate photo: %w", err)
	}

	if validatedPhoto.Data.Status == "valid" {
		userToSave.IsPhotoValidated = true
	}
	return true, nil
}
