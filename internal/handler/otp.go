package handler

import (
	"net/http"

	"github.com/imansohibul/otp-service/entity"
	"github.com/imansohibul/otp-service/generated"
	"github.com/labstack/echo/v4"
)

// Request a new OTP
// (POST /otp/request)
func (r *RestAPIServer) PostOtpRequest(eCtx echo.Context) error {
	var (
		ctx = eCtx.Request().Context()
		req = new(generated.PostOtpRequestJSONRequestBody)
	)

	if err := eCtx.Bind(req); err != nil {
		return eCtx.JSON(http.StatusBadRequest, entity.ErrInvalidRequest)
	}

	otp, err := r.OtpUsecase.Create(ctx, req.UserId)
	if err != nil {
		return eCtx.JSON(http.StatusBadRequest, err)
	}

	return eCtx.JSON(http.StatusOK, generated.RequestOtpResponseSuccess{
		UserId: otp.UserID,
		Otp:    otp.OTPCode,
	})
}

// Validate an OTP
// (POST /otp/validate)
func (r *RestAPIServer) PostOtpValidate(eCtx echo.Context) error {
	var (
		ctx = eCtx.Request().Context()
		req = new(generated.PostOtpValidateJSONRequestBody)
	)

	if err := eCtx.Bind(req); err != nil {
		return eCtx.JSON(http.StatusBadRequest, entity.ErrInvalidRequest)
	}

	otp, err := r.OtpUsecase.Validate(ctx, req.UserId, req.Otp)
	if err != nil {
		return eCtx.JSON(http.StatusBadRequest, err)
	}

	return eCtx.JSON(http.StatusOK, generated.ValidateOtpResponseSuccess{
		UserId:  otp.UserID,
		Message: "OTP Validated successfully",
	})
}
