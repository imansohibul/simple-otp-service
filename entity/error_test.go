package entity_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/imansohibul/otp-service/entity"
	"github.com/imansohibul/otp-service/internal/handler/middleware"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestErrorHandler(t *testing.T) {
	tests := []struct {
		name       string
		err        error
		wantStatus int
		wantBody   string
	}{
		{
			name:       "DomainError - BadRequest",
			err:        &entity.DomainError{Code: "INVALID_INPUT", Message: "invalid input"},
			wantStatus: http.StatusBadRequest,
			wantBody:   `{"error":"INVALID_INPUT","error_description":"invalid input"}`,
		},
		{
			name:       "DomainError - NotFound",
			err:        &entity.DomainError{Code: entity.ErrOTPNotFound.Code, Message: "OTP not found"},
			wantStatus: http.StatusNotFound,
			wantBody:   `{"error":"otp_not_found","error_description":"OTP not found"}`,
		},
		{
			name:       "Other error - InternalServerError",
			err:        errors.New("random error"),
			wantStatus: http.StatusInternalServerError,
			wantBody:   `{"error":"INTERNAL_SERVER_ERROR","error_description":"Something went wrong! Please try again later"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			rec := httptest.NewRecorder()
			ctx := e.NewContext(req, rec)

			// Call ErrorHandler
			middleware.ErrorHandler(tt.err, ctx)

			assert.Equal(t, tt.wantStatus, rec.Code)
			assert.JSONEq(t, tt.wantBody, rec.Body.String())
		})
	}
}
