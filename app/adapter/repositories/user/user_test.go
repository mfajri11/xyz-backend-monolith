package repository

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/mfajri11/xyz-backend-monolith/app/core/domain"
	uhttp "github.com/mfajri11/xyz-backend-monolith/util/http"
	"github.com/stretchr/testify/assert"
)

// func newIOReader(b interface{}) io.Reader {
// 	bb, _ := json.Marshal(b)
// 	return io.NopCloser(bytes.NewReader(bb))
// }

func newIOReaderCloser(b interface{}) io.ReadCloser {
	bb, _ := json.Marshal(b)
	return io.NopCloser(bytes.NewReader(bb))
}

func newByte(v interface{}) []byte {
	bb, _ := json.Marshal(v)
	return bb
}

func TestUserRepository_FindOneByNationalID(t *testing.T) {
	type args struct {
		ctx context.Context
		nid string
	}
	type mock struct {
		sqlmock.Sqlmock
	}
	tests := []struct {
		name        string
		repo        *UserRepository
		args        args
		prepareMock func(mock *mock)
		wantUser    *domain.UserEntity
		wantErr     bool
		wantNoRows  bool
	}{
		{
			name: "Given a valid national id, it should return a user",
			repo: &UserRepository{},
			args: args{
				ctx: context.Background(),
				nid: "12345678912345678",
			},
			prepareMock: func(mock *mock) {
				mock.
					ExpectQuery(regexp.QuoteMeta(getUserByNationalID)).
					WithArgs("12345678912345678").WillReturnRows(sqlmock.NewRows([]string{"id", "national_id", "full_name", "legal_name", "is_nid_valid", "is_photo_valid"}).AddRow(1, "12345678912345678", "John Doe", "John Doe", true, true))
			},
			wantUser: &domain.UserEntity{
				ID:                    1,
				NationalID:            "12345678912345678",
				FullName:              "John Doe",
				LegalName:             "John Doe",
				IsNationalIDValidated: true,
				IsPhotoValidated:      true,
			},
		},
		{
			name: "Given a invalid national id, it should return errors no rows",
			repo: &UserRepository{},
			args: args{
				ctx: context.Background(),
				nid: "0000000000000000",
			},
			prepareMock: func(mock *mock) {
				mock.
					ExpectQuery(regexp.QuoteMeta(getUserByNationalID)).
					WithArgs("0000000000000000").WillReturnError(sql.ErrNoRows)
			},
		},
		{
			name: "Given a invalid national id, it should return error",
			repo: &UserRepository{},
			args: args{
				ctx: context.Background(),
				nid: "12345678912345678",
			},
			prepareMock: func(mock *mock) {
				mock.
					ExpectQuery(regexp.QuoteMeta(getUserByNationalID)).
					WithArgs("12345678912345678").WillReturnError(sql.ErrConnDone)
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			conn, sqlMock, err := sqlmock.New()
			assert.NoError(t, err)
			defer conn.Close()
			tt.prepareMock(&mock{
				Sqlmock: sqlMock,
			})

			tt.repo.dbConn = conn

			gotUser, err := tt.repo.FindOneByNationalID(tt.args.ctx, tt.args.nid)
			assert.Equal(t, tt.wantErr, err != nil, err)
			assert.Equal(t, tt.wantUser, gotUser, gotUser)
		})
	}
}

func TestUserRepository_UpdateByID(t *testing.T) {
	type args struct {
		ctx  context.Context
		user domain.UserEntity
	}
	type mock struct {
		sqlmock.Sqlmock
	}
	tests := []struct {
		name        string
		repo        *UserRepository
		args        args
		prepareMock func(mock *mock)
		wantErr     bool
	}{
		{
			name: "Given a valid id and user data, it should return no error",
			repo: &UserRepository{},
			args: args{
				ctx: context.Background(),
				user: domain.UserEntity{
					ID:                    1,
					NationalID:            "1122334455667788",
					FullName:              "John Doe",
					LegalName:             "John Doe",
					IsNationalIDValidated: true,
					IsPhotoValidated:      true,
				},
			},
			prepareMock: func(m *mock) {
				m.ExpectExec(regexp.QuoteMeta(queryUpdateUserById)).WillReturnResult(sqlmock.NewResult(1, 1))
			},
		},
		{
			name: "Given a valid id and user data, it should return no error",
			repo: &UserRepository{},
			args: args{
				ctx: context.Background(),
				user: domain.UserEntity{
					ID:                    -1,
					NationalID:            "1122334455667788",
					FullName:              "John Doe",
					LegalName:             "John Doe",
					IsNationalIDValidated: true,
					IsPhotoValidated:      true,
				},
			},
			prepareMock: func(m *mock) {
				m.ExpectExec(regexp.QuoteMeta(queryUpdateUserById)).WillReturnResult(nil).WillReturnError(sql.ErrTxDone)
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			conn, sqlMock, err := sqlmock.New()
			assert.NoError(t, err)
			defer conn.Close()
			tt.prepareMock(&mock{
				Sqlmock: sqlMock,
			})

			tt.repo.dbConn = conn

			err = tt.repo.UpdateByID(tt.args.ctx, tt.args.user)

			assert.Equal(t, tt.wantErr, err != nil, err)

		})
	}
}

func TestUserRepository_ValidateNationalID(t *testing.T) {
	type args struct {
		ctx context.Context
		req domain.KYCValidateNationalIDReq
	}
	type mock func(r *http.Request, w *httptest.ResponseRecorder) uhttp.DoFunc
	tests := []struct {
		name    string
		repo    *UserRepository
		args    args
		doMock  mock
		want    *domain.KYCValidateNationalIDResp
		wantErr bool
	}{
		{
			name: "Given a valid national id, it should return true",
			repo: &UserRepository{},
			args: args{
				ctx: context.Background(),
				req: domain.KYCValidateNationalIDReq{
					NationalID:  "1122334455667788",
					LegalName:   "John Doe",
					DateOfBirth: "2000-01-01",
					ReferenceID: "12345678912345678",
				},
			},
			doMock: func(r *http.Request, w *httptest.ResponseRecorder) uhttp.DoFunc {
				return func(r *http.Request) (*http.Response, error) {
					w.WriteHeader(http.StatusOK)
					w.Write(newByte(domain.KYCValidateNationalIDResp{
						Message: "Success",
						Data: domain.KYCData{
							NationalID:  true,
							LegalName:   true,
							DateOfBirth: true,
							ReferenceID: "12345678912345678",
						},
					}))
					return w.Result(), nil

				}
			},
			want: &domain.KYCValidateNationalIDResp{
				Message: "Success",
				Data: domain.KYCData{
					NationalID:  true,
					LegalName:   true,
					DateOfBirth: true,
					ReferenceID: "12345678912345678",
				},
			},
		},
		{
			name: "Given a invalid national id, it should return false",
			repo: &UserRepository{},
			args: args{
				ctx: context.Background(),
				req: domain.KYCValidateNationalIDReq{
					NationalID:  "1122334455667788",
					LegalName:   "John Doe",
					DateOfBirth: "2000-01-01",
					ReferenceID: "12345678912345678",
				},
			},
			doMock: func(r *http.Request, w *httptest.ResponseRecorder) uhttp.DoFunc {
				return func(r *http.Request) (*http.Response, error) {
					w.WriteHeader(http.StatusOK)
					w.Write(newByte(domain.KYCValidateNationalIDResp{
						Message: "Fail validate national id",
						Data: domain.KYCData{
							NationalID:  false,
							LegalName:   false,
							DateOfBirth: false,
							ReferenceID: "12345678912345678",
						},
					}))
					return w.Result(), nil

				}
			},
			want: &domain.KYCValidateNationalIDResp{
				Message: "Fail validate national id",
				Data: domain.KYCData{
					NationalID:  false,
					LegalName:   false,
					DateOfBirth: false,
					ReferenceID: "12345678912345678",
				},
			},
		},
		{
			name: "Given a valid national id, but it should return error",
			repo: &UserRepository{},
			args: args{
				ctx: context.Background(),
				req: domain.KYCValidateNationalIDReq{
					NationalID:  "1122334455667788",
					LegalName:   "John Doe",
					DateOfBirth: "2000-01-01",
					ReferenceID: "12345678912345678",
				},
			},
			doMock: func(r *http.Request, w *httptest.ResponseRecorder) uhttp.DoFunc {
				return func(r *http.Request) (*http.Response, error) {
					w.WriteHeader(http.StatusInternalServerError)
					return w.Result(), nil

				}
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httptest.NewRequest(http.MethodPost, "/verify/national-id", newIOReaderCloser(tt.args.req))
			w := httptest.NewRecorder()
			tt.repo.kycClient = uhttp.NewMock(tt.doMock(r, w))

			got, err := tt.repo.ValidateNationalID(tt.args.ctx, tt.args.req)

			assert.Equal(t, tt.wantErr, err != nil, err)
			assert.Equal(t, tt.want, got, got)
		})
	}
}

func TestUserRepository_ValidateSalary(t *testing.T) {
	type args struct {
		ctx context.Context
		req domain.KYCValidateSalaryReq
	}
	type mock func(r *http.Request, w *httptest.ResponseRecorder) uhttp.DoFunc
	tests := []struct {
		name    string
		repo    *UserRepository
		args    args
		doMock  mock
		want    *domain.KYCValidateSalaryResp
		wantErr bool
	}{
		{
			name: "Given a valid salary, it should return true",
			repo: &UserRepository{},
			args: args{
				ctx: context.Background(),
				req: domain.KYCValidateSalaryReq{
					NationalID:  "1122334455667788",
					LegalName:   "John Doe",
					Salary:      "1000",
					ReferenceID: "12345678912345678",
				},
			},
			doMock: func(r *http.Request, w *httptest.ResponseRecorder) uhttp.DoFunc {
				return func(r *http.Request) (*http.Response, error) {
					w.WriteHeader(http.StatusOK)
					w.Write(newByte(domain.KYCValidateSalaryResp{
						Message: "Success",
						Data: domain.KYCData{
							NationalID:  true,
							LegalName:   true,
							ReferenceID: "12345678912345678",
							SalaryUper:  "2000",
							SalaryLower: "500",
						},
					}))
					return w.Result(), nil

				}
			},
			want: &domain.KYCValidateSalaryResp{
				Message: "Success",
				Data: domain.KYCData{
					NationalID:  true,
					LegalName:   true,
					ReferenceID: "12345678912345678",
					SalaryUper:  "2000",
					SalaryLower: "500",
				},
			},
		},
		{
			name: "Given a valid salary, but it should return error",
			repo: &UserRepository{},
			args: args{
				ctx: context.Background(),
				req: domain.KYCValidateSalaryReq{
					NationalID:  "1122334455667788",
					LegalName:   "John Doe",
					Salary:      "100000",
					ReferenceID: "12345678912345678",
				},
			},
			doMock: func(r *http.Request, w *httptest.ResponseRecorder) uhttp.DoFunc {
				return func(r *http.Request) (*http.Response, error) {
					w.WriteHeader(http.StatusInternalServerError)
					return w.Result(), nil

				}
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httptest.NewRequest(http.MethodPost, "/verify/income", newIOReaderCloser(tt.args.req))
			w := httptest.NewRecorder()
			tt.repo.kycClient = uhttp.NewMock(tt.doMock(r, w))

			got, err := tt.repo.ValidateSalary(tt.args.ctx, tt.args.req)

			assert.Equal(t, tt.wantErr, err != nil, err)
			assert.Equal(t, tt.want, got, got)
		})
	}
}
