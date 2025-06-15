package port

import (
	"context"

	"github.com/mfajri11/xyz-backend-monolith/app/core/domain"
)

type UserRepository interface {
	FindOneByNationalID(ctx context.Context, nid string) (user *domain.UserEntity, err error)
	UpdateByID(ctx context.Context, user domain.UserEntity) error
	ValidateSalary(ctx context.Context, req domain.KYCValidateSalaryReq) (*domain.KYCValidateSalaryResp, error)
	ValidateNationalID(ctx context.Context, req domain.KYCValidateNationalIDReq) (*domain.KYCValidateNationalIDResp, error)
	ValidatePhoto(ctx context.Context, req domain.KYCValidatePhotoReq) (*domain.KYCValidatePhotoResp, error)
}

type UserService interface {
	ValidateData(ctx context.Context, req domain.ValidateUserReq) (bool, error)
}
