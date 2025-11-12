package middleware_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/imansohibul/otp-service/entity"
	"github.com/imansohibul/otp-service/generated"
	"github.com/imansohibul/otp-service/internal/handler/middleware"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestErrorHandler(t *testing.T) {
	tests := []struct {
		name       string
		err        error
		committed  bool
		wantStatus int
		wantBody   generated.ErrorResponse
	}{
		{
			name:       "response already committed",
			err:        errors.New("some error"),
			committed:  true,
			wantStatus: http.StatusOK,
			wantBody:   generated.ErrorResponse{}, // will not modify response
		},
		{
			name:       "DomainError - NotFound",
			err:        &entity.DomainError{Code: entity.ErrOTPNotFound.Code, Message: "Resource not found"},
			committed:  false,
			wantStatus: http.StatusNotFound,
			wantBody:   generated.ErrorResponse{Error: entity.ErrOTPNotFound.Code, ErrorDescription: "Resource not found"},
		},
		{
			name:       "DomainError - BadRequest",
			err:        &entity.DomainError{Code: entity.ErrInvalidRequest.Code, Message: "Invalid request"},
			committed:  false,
			wantStatus: http.StatusBadRequest,
			wantBody:   generated.ErrorResponse{Error: entity.ErrInvalidRequest.Code, ErrorDescription: "Invalid request"},
		},
		{
			name:       "Other error - InternalServerError",
			err:        errors.New("some internal error"),
			committed:  false,
			wantStatus: http.StatusInternalServerError,
			wantBody:   generated.ErrorResponse{Error: "INTERNAL_SERVER_ERROR", ErrorDescription: "Something went wrong! Please try again later"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			rec := httptest.NewRecorder()
			ctx := e.NewContext(req, rec)

			// Simulate committed response if needed
			if tt.committed {
				ctx.Response().Committed = true
				ctx.Response().WriteHeader(http.StatusOK)
				ctx.Response().Write([]byte("Already committed"))
			}

			middleware.ErrorHandler(tt.err, ctx)

			if tt.committed {
				// Response should not be modified
				assert.Equal(t, http.StatusOK, rec.Code)
				assert.Equal(t, "Already committed", rec.Body.String())
			} else {
				assert.Equal(t, tt.wantStatus, rec.Code)
				var resp generated.ErrorResponse
				err := json.Unmarshal(rec.Body.Bytes(), &resp)
				assert.NoError(t, err)
				assert.Equal(t, tt.wantBody.Error, resp.Error)
				assert.Equal(t, tt.wantBody.ErrorDescription, resp.ErrorDescription)
			}
		})
	}
}
