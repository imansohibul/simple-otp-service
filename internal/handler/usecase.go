package handler

import (
	"context"

	"github.com/imansohibul/otp-service/entity"
)

//go:generate mockgen -destination=mock/usecase.go -package=mock -source=usecase.go

// OTPUsecase defines the business logic interface for OTP (One-Time Password) operations.
// It handles the creation and validation of OTPs.
type OTPUsecase interface {
	// Create generates a new OTP for the specified user and stores it in the system.
	// The OTP will have an expiration time and can only be used once.
	Create(ctx context.Context, userID string) (*entity.OTP, error)

	// Validate verifies that the provided OTP code is valid for the specified user.
	// This checks if the code matches, hasn't expired, and hasn't been used before.
	// Upon successful validation, the OTP should be marked as validated.
	Validate(ctx context.Context, userID string, otpCode string) (*entity.OTP, error)
}
