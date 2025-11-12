package handler_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/imansohibul/otp-service/entity"
	"github.com/imansohibul/otp-service/generated"
	"github.com/imansohibul/otp-service/internal/handler"
	usecasemock "github.com/imansohibul/otp-service/internal/handler/mock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestPostOtpRequest(t *testing.T) {
	tests := []struct {
		name               string
		requestBody        interface{}
		mockSetup          func(*testing.T, *usecasemock.MockOTPUsecase)
		expectedStatusCode int
		expectedBody       string
	}{
		{
			name:        "Request OTP - Success",
			requestBody: &generated.PostOtpRequestJSONRequestBody{UserId: "user123"},
			mockSetup: func(t *testing.T, otpUsecase *usecasemock.MockOTPUsecase) {
				otpUsecase.EXPECT().
					Create(gomock.Any(), "user123").
					Return(&entity.OTP{UserID: "user123", OTPCode: "123456"}, nil)
			},
			expectedStatusCode: http.StatusOK,
			expectedBody:       `"user_id":"user123"`,
		},
		{
			name:        "Request OTP - Success with Different User",
			requestBody: &generated.PostOtpRequestJSONRequestBody{UserId: "user456"},
			mockSetup: func(t *testing.T, otpUsecase *usecasemock.MockOTPUsecase) {
				otpUsecase.EXPECT().
					Create(gomock.Any(), "user456").
					Return(&entity.OTP{UserID: "user456", OTPCode: "654321"}, nil)
			},
			expectedStatusCode: http.StatusOK,
			expectedBody:       `"otp":"654321"`,
		},
		{
			name:        "Request OTP - Invalid Request Body",
			requestBody: "invalid json",
			mockSetup: func(t *testing.T, otpUsecase *usecasemock.MockOTPUsecase) {
				// No mock needed for invalid request
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedBody:       "", // Will contain entity.ErrInvalidRequest
		},
		{
			name:        "Request OTP - Nil Request Body",
			requestBody: nil,
			mockSetup: func(t *testing.T, otpUsecase *usecasemock.MockOTPUsecase) {
				// No mock needed for nil request
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedBody:       "",
		},
		{
			name:        "Request OTP - Usecase Error",
			requestBody: &generated.PostOtpRequestJSONRequestBody{UserId: "user789"},
			mockSetup: func(t *testing.T, otpUsecase *usecasemock.MockOTPUsecase) {
				otpUsecase.EXPECT().
					Create(gomock.Any(), "user789").
					Return(nil, entity.ErrOTPDuplicate)
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedBody:       "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			e := echo.New()

			bodyBytes, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest(http.MethodPost, "/otp/request", bytes.NewReader(bodyBytes))
			if tt.requestBody != nil && tt.requestBody != "invalid json" {
				req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			}
			rec := httptest.NewRecorder()

			mockOTPUsecase := usecasemock.NewMockOTPUsecase(ctrl)
			tt.mockSetup(t, mockOTPUsecase)

			server := handler.RestAPIServer{
				Echo:       e,
				OtpUsecase: mockOTPUsecase,
			}

			c := e.NewContext(req, rec)
			err := server.PostOtpRequest(c)
			if err != nil {
				t.Errorf("Error: %v", err)
			}

			assert.Equal(t, tt.expectedStatusCode, rec.Code)
			if tt.expectedBody != "" {
				assert.Contains(t, rec.Body.String(), tt.expectedBody)
			}
		})
	}
}

func TestPostOtpValidate(t *testing.T) {
	tests := []struct {
		name               string
		requestBody        interface{}
		mockSetup          func(*testing.T, *usecasemock.MockOTPUsecase)
		expectedStatusCode int
		expectedBody       string
	}{
		{
			name:        "Validate OTP - Success",
			requestBody: &generated.PostOtpValidateJSONRequestBody{UserId: "user123", Otp: "123456"},
			mockSetup: func(t *testing.T, otpUsecase *usecasemock.MockOTPUsecase) {
				otpUsecase.EXPECT().
					Validate(gomock.Any(), "user123", "123456").
					Return(&entity.OTP{UserID: "user123", OTPCode: "123456"}, nil)
			},
			expectedStatusCode: http.StatusOK,
			expectedBody:       `"user_id":"user123"`,
		},
		{
			name:        "Validate OTP - Success with Different User",
			requestBody: &generated.PostOtpValidateJSONRequestBody{UserId: "user456", Otp: "654321"},
			mockSetup: func(t *testing.T, otpUsecase *usecasemock.MockOTPUsecase) {
				otpUsecase.EXPECT().
					Validate(gomock.Any(), "user456", "654321").
					Return(&entity.OTP{UserID: "user456", OTPCode: "654321"}, nil)
			},
			expectedStatusCode: http.StatusOK,
			expectedBody:       `"message":"OTP Validated successfully"`,
		},
		{
			name:        "Validate OTP - Invalid Request Body",
			requestBody: "invalid json",
			mockSetup: func(t *testing.T, otpUsecase *usecasemock.MockOTPUsecase) {
				// No mock needed for invalid request
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedBody:       "", // Will contain entity.ErrInvalidRequest
		},
		{
			name:        "Validate OTP - Nil Request Body",
			requestBody: nil,
			mockSetup: func(t *testing.T, otpUsecase *usecasemock.MockOTPUsecase) {
				// No mock needed for nil request
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedBody:       "",
		},
		{
			name:        "Validate OTP - OTP Expired",
			requestBody: &generated.PostOtpValidateJSONRequestBody{UserId: "user789", Otp: "789012"},
			mockSetup: func(t *testing.T, otpUsecase *usecasemock.MockOTPUsecase) {
				otpUsecase.EXPECT().
					Validate(gomock.Any(), "user789", "789012").
					Return(nil, entity.ErrOTPExpired)
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedBody:       "",
		},
		{
			name:        "Validate OTP - OTP Already Used",
			requestBody: &generated.PostOtpValidateJSONRequestBody{UserId: "user101", Otp: "101010"},
			mockSetup: func(t *testing.T, otpUsecase *usecasemock.MockOTPUsecase) {
				otpUsecase.EXPECT().
					Validate(gomock.Any(), "user101", "101010").
					Return(nil, entity.ErrOTPUsed)
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedBody:       "",
		},
		{
			name:        "Validate OTP - OTP Not Found",
			requestBody: &generated.PostOtpValidateJSONRequestBody{UserId: "user202", Otp: "999999"},
			mockSetup: func(t *testing.T, otpUsecase *usecasemock.MockOTPUsecase) {
				otpUsecase.EXPECT().
					Validate(gomock.Any(), "user202", "999999").
					Return(nil, entity.ErrOTPNotFound)
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedBody:       "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			e := echo.New()

			bodyBytes, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest(http.MethodPost, "/otp/validate", bytes.NewReader(bodyBytes))
			if tt.requestBody != nil && tt.requestBody != "invalid json" {
				req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			}
			rec := httptest.NewRecorder()

			mockOTPUsecase := usecasemock.NewMockOTPUsecase(ctrl)
			tt.mockSetup(t, mockOTPUsecase)

			server := handler.RestAPIServer{
				Echo:       e,
				OtpUsecase: mockOTPUsecase,
			}

			c := e.NewContext(req, rec)
			err := server.PostOtpValidate(c)
			if err != nil {
				t.Errorf("Error: %v", err)
			}

			assert.Equal(t, tt.expectedStatusCode, rec.Code)
			if tt.expectedBody != "" {
				assert.Contains(t, rec.Body.String(), tt.expectedBody)
			}
		})
	}
}
