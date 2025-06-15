package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
	"github.com/mfajri11/xyz-backend-monolith/app/core/domain"
	"github.com/mfajri11/xyz-backend-monolith/app/core/port"
	"github.com/mfajri11/xyz-backend-monolith/util/apperror"
	"github.com/rs/zerolog/log"
)

type LoanHandler struct {
	loanService port.LoanService
	userService port.UserService
}

func New(loanService port.LoanService, userService port.UserService) *LoanHandler {
	return &LoanHandler{
		loanService: loanService,
		userService: userService,
	}
}

func (handler *LoanHandler) CreateLoan(c *gin.Context) {
	var req domain.CreateLoanReq
	rid := requestid.Get(c)
	logger := log.With().Str("requestID", rid).Logger()
	uid := c.GetInt64("uid")
	if uid == 0 {
		logger.Error().Err(fmt.Errorf("CreateLoan: invalid user id")).Msg("")
		writeError(c, apperror.ErrBadRequest)
		return
	}

	err := c.ShouldBindJSON(&req)
	if err != nil {
	}

	isValid, err := handler.userService.ValidateData(c, domain.ValidateUserReq{
		NationalID:      req.NationalID,
		LegalName:       req.LegalName,
		BirthOfDate:     req.BirthOfDate,
		Salary:          req.Salary,
		NationalIDPhoto: req.NationalIDPhoto,
		UserPhoto:       req.UserPhoto,
	})

	if err != nil {
		logger.Error().Err(err).Msg("error while validate data")
		writeError(c, err)
		return
	}

	if !isValid {
		logger.Err(fmt.Errorf("CreateLoan: invalid user data")).Msg("")
		writeError(c, apperror.ErrBadRequest)
		return
	}

	otrAmoutn := handler.loanService.CalculateOTRAmount(req.Amount)

	handler.loanService.CreateLoan(c, domain.Loan{
		UserID:          uid,
		ContractNumber:  req.ContractNumber,
		OTRAmount:       otrAmoutn,
		PrincipalAmount: otrAmoutn - req.DownPayment,
	})

	c.JSON(http.StatusOK, nil)
}

func (handler *LoanHandler) GetLoanByContractNumber(c *gin.Context) {
	contractNumber := c.Param("contractNumber")
	rid := requestid.Get(c)
	logger := log.With().Str("requestID", rid).Logger()
	uid := c.GetInt64("uid")
	if uid == 0 {
		logger.Error().Err(fmt.Errorf("CreateLoan: invalid user id")).Msg("")
		writeError(c, apperror.ErrBadRequest)
	}
	loan, err := handler.loanService.GetLoanByUserIDAndContractNumber(c, uid, contractNumber)
	if err != nil {
		writeError(c, err)
		return
	}

	writeSuccess(c, loan)
}

func (handler *LoanHandler) GetLoanPaymentsByContractNumber(c *gin.Context) {
	contractNumber := c.Param("contractNumber")
	rid := requestid.Get(c)
	logger := log.With().Str("requestID", rid).Logger()

	uid := c.GetInt64("uid")
	if uid == 0 {
		logger.Error().Err(fmt.Errorf("CreateLoan: invalid user id")).Msg("")
		writeError(c, apperror.ErrBadRequest)
	}
	loanPayments, err := handler.loanService.GetLoanPaymentsByUserIDAndContractNumber(c, uid, contractNumber)
	if err != nil {
		writeError(c, err)
		return
	}

	writeSuccess(c, loanPayments)
}
