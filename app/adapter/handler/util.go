package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/mfajri11/xyz-backend-monolith/app/core/domain"
	"github.com/mfajri11/xyz-backend-monolith/util/apperror"
)

func writeError(c *gin.Context, err error) {
	serr := apperror.SentinelError(err)
	if serr != nil {
		sc, msg := serr.APIError()
		c.JSON(sc, domain.GeneralResponse{
			Success: false,
			Message: msg,
		})
	}

	c.JSON(500, domain.GeneralResponse{
		Success: false,
		Message: "opps! something went wrong",
	})
}

func writeSuccess(c *gin.Context, data interface{}) {
	c.JSON(200, domain.GeneralResponse{
		Success: true,
		Message: "success",
		Data:    data,
	})
}
