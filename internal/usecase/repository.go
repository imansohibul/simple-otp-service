package usecase

// This file contains the interfaces for the repository layer.
// The repository layer is responsible for interacting with the database.
// For testing purpose we will generate mock implementations of these
// interfaces using mockgen. See the Makefile for more information.

import (
	"context"

	"github.com/imansohibul/otp-service/entity"
)

//go:generate mockgen -destination=mock/repository.go -package=mock -source=repository.go

// TransactionManager defines the interface for managing database transactions.
type TransactionManager interface {
	WithTransaction(ctx context.Context, fn func(ctx context.Context) error) error
}

// OTPRepository defines the interface for OTP data access operations
type OTPRepository interface {
	// Create inserts a new OTP into the database.
	// Returns entity.ErrDuplicateOTP if an OTP with the same user_id and otp_code already exists.
	Create(ctx context.Context, otp *entity.OTP) error

	// FindByUserIDAndCode retrieves an OTP by user ID and OTP code from the database.
	FindByUserIDAndCode(ctx context.Context, userID string, otpCode string) (*entity.OTP, error)

	// Update updates an existing OTP record in the database.
	// Typically used to update the status and validated_at fields.
	Update(ctx context.Context, otp *entity.OTP) error

	// GetLastByUserID retrieves the most recent OTP record for a given user,
	// ordered by creation timestamp descending.
	GetLastByUserID(ctx context.Context, userID string) (*entity.OTP, error)
}
