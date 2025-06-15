package app

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/mfajri11/xyz-backend-monolith/app/adapter/handler"
	loanRepository "github.com/mfajri11/xyz-backend-monolith/app/adapter/repositories/loan"
	userRepository "github.com/mfajri11/xyz-backend-monolith/app/adapter/repositories/user"
	loanService "github.com/mfajri11/xyz-backend-monolith/app/core/service/loan"
	userService "github.com/mfajri11/xyz-backend-monolith/app/core/service/user"
	"github.com/mfajri11/xyz-backend-monolith/infra/db/mysql"
	"github.com/mfajri11/xyz-backend-monolith/util/config"
	uhttp "github.com/mfajri11/xyz-backend-monolith/util/http"
)

func Run() error {

	cfg := config.Get()

	db := mysql.MustNew(mysql.DataSource(cfg.Database.Username, cfg.Database.Password, cfg.Database.Host, cfg.Database.Port, cfg.Database.DataBaseName), mysql.WithMaxIdleConns(cfg.Database.IdleConnection), mysql.WithMaxOpenConns(cfg.Database.OpenConnection), mysql.WithMaxLifetimeConn(cfg.Database.ConnectionMaxLifeTime))

	kycClient := uhttp.NewClient(cfg.KYCClient.BaseURL, cfg.KYCClient.APIKey, cfg.KYCClient.APPID)
	userRepo := userRepository.New(db, kycClient)
	loanRepo := loanRepository.New(db)

	userSvc := userService.New(userRepo)
	loanSerice := loanService.New(loanRepo)

	loanHandler := handler.New(loanSerice, userSvc)
	router := gin.Default()
	router.POST("/loan", loanHandler.CreateLoan)
	router.GET("/loans/:contractNumber", loanHandler.GetLoanByContractNumber)
	router.GET("/loan/:contractNumber/payments", loanHandler.GetLoanPaymentsByContractNumber)
	return router.Run(fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port))
}
