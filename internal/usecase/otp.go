package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/imansohibul/otp-service/entity"
)

const (
	otpValidityDuration = 2 * time.Minute
	otpRateLimitWindow  = 2 * time.Minute
)

type otpUsecase struct {
	otpRepo      OTPRepository
	otpGenerator OTPGenerator
}

func NewOtpUsecase(otpRepo OTPRepository, otpGenerator OTPGenerator) *otpUsecase {
	return &otpUsecase{
		otpRepo:      otpRepo,
		otpGenerator: otpGenerator,
	}
}

func (o *otpUsecase) Create(ctx context.Context, userID string) (*entity.OTP, error) {
	// Check rate limiting
	lastOTP, _ := o.otpRepo.GetLastByUserID(ctx, userID)
	if lastOTP != nil && lastOTP.Status == entity.OTPStatusCreated {
		if time.Since(lastOTP.CreatedAt) < otpRateLimitWindow {
			return nil, entity.ErrOTPRateLimitExceeded
		}
	}

	otpCode, err := o.otpGenerator.Generate()
	if err != nil {
		return nil, fmt.Errorf("failed to generate OTP code: %w", err)
	}

	otp := &entity.OTP{
		UserID:    userID,
		OTPCode:   otpCode,
		Status:    entity.OTPStatusCreated,
		ExpiresAt: time.Now().Add(2 * time.Minute),
	}
	if err := o.otpRepo.Create(ctx, otp); err != nil {
		return nil, err
	}

	return otp, nil
}

func (o *otpUsecase) Validate(ctx context.Context, userID string, otpCode string) (*entity.OTP, error) {
	otp, err := o.otpRepo.FindByUserIDAndCode(ctx, userID, otpCode)
	if err != nil {
		return nil, err
	}

	// Validate OTP status and expiration
	if err := o.validateOTPStatus(ctx, otp); err != nil {
		return nil, err
	}

	// Mark OTP as used
	if err := o.markOTPAsValidated(ctx, otp); err != nil {
		return nil, fmt.Errorf("failed to update OTP status: %w", err)
	}

	return otp, nil
}

// validateOTPStatus checks if OTP is expired or already used
func (o *otpUsecase) validateOTPStatus(ctx context.Context, otp *entity.OTP) error {
	now := time.Now()

	// Check if already used (before expiration check for security)
	if otp.Status == entity.OTPStatusValidated {
		return entity.ErrOTPUsed
	}

	// Check expiration
	if now.After(otp.ExpiresAt) {
		if otp.Status != entity.OTPStatusExpired {
			otp.Status = entity.OTPStatusExpired
			if err := o.otpRepo.Update(ctx, otp); err != nil {
				return err
			}
		}
		return entity.ErrOTPExpired
	}

	return nil
}

// markOTPAsValidated updates OTP status to used
func (o *otpUsecase) markOTPAsValidated(ctx context.Context, otp *entity.OTP) error {
	now := time.Now()
	otp.Status = entity.OTPStatusValidated
	otp.ValidatedAt = &now
	return o.otpRepo.Update(ctx, otp)
}
