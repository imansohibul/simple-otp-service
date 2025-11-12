package middleware

import (
	"errors"
	"net/http"

	"github.com/imansohibul/otp-service/entity"
	"github.com/imansohibul/otp-service/generated"
	"github.com/labstack/echo/v4"
)

func ErrorHandler(err error, ctx echo.Context) {
	var domainErr *entity.DomainError

	if ctx.Response().Committed {
		return
	}

	// Check if the error is a DomainError
	if errors.As(err, &domainErr) {
		httpStatus := http.StatusBadRequest
		if domainErr.Code == entity.ErrOTPNotFound.Code  {
			httpStatus = http.StatusNotFound
		}

		_ = ctx.JSON(httpStatus, generated.ErrorResponse{
			Error:            domainErr.Code,
			ErrorDescription: domainErr.Message,
		})
		return
	}

	// For other errors, return a 500 Internal Server Error
	_ = ctx.JSON(http.StatusInternalServerError, generated.ErrorResponse{
		Error:            "INTERNAL_SERVER_ERROR",
		ErrorDescription: "Something went wrong! Please try again later"},
	)
}
