package repository

import (
	"context"
	"database/sql"

	"github.com/imansohibul/otp-service/entity"
	"github.com/jmoiron/sqlx"
)

// otpRepository implements the OTPRepository interface
type otpRepository struct {
	db *sqlx.DB
}

// NewOTPRepository creates a new instance of otpRepository
func NewOTPRepository(db *sqlx.DB) *otpRepository {
	return &otpRepository{
		db: db,
	}
}

// Create inserts a new OTP into the database
func (o *otpRepository) Create(ctx context.Context, otp *entity.OTP) error {
	const query = `
		INSERT INTO otps (user_id, otp_code, status, expires_at)
		VALUES (?, ?, ?, ?)
	`
	_, err := getExecutor(ctx, o.db).ExecContext(
		ctx,
		query,
		otp.UserID,
		otp.OTPCode,
		otp.Status,
		otp.ExpiresAt,
	)
	if err != nil {
		// Check if the error is a unique constraint violation
		if isUniqueConstraintViolation(err) {
			return entity.ErrOTPDuplicate
		}
		return err
	}

	return err
}

// FindByUserIDAndCode retrieves an OTP by user ID and OTP code from the database
func (o *otpRepository) FindByUserIDAndCode(ctx context.Context, userID string, otpCode string) (*entity.OTP, error) {
	const query = `
		SELECT id, user_id, otp_code, status, created_at, expires_at, validated_at
		FROM otps
		WHERE user_id = ? AND otp_code = ?
	`

	var otpRow otpRow
	if err := getExecutor(ctx, o.db).GetContext(ctx, &otpRow, query, userID, otpCode); err != nil {
		// Check if the error is sql.ErrNoRows to return entity.ErrNotFound
		if err == sql.ErrNoRows {
			return nil, entity.ErrOTPNotFound
		}
		return nil, err
	}

	return otpRow.ToEntity(), nil
}

// Update updates an OTP record in the database
func (o *otpRepository) Update(ctx context.Context, otp *entity.OTP) error {
	const query = `
		UPDATE otps
		SET status = ?, validated_at = ?
		WHERE id = ?
	`
	_, err := getExecutor(ctx, o.db).ExecContext(
		ctx,
		query,
		otp.Status,
		otp.ValidatedAt,
		otp.ID,
	)

	return err
}

// GetLastByUserID retrieves the most recent OTP for a specific user,
// ordered by creation timestamp descending. Returns entity.ErrOTPNotFound
// if no OTP exists for the user.
func (o *otpRepository) GetLastByUserID(ctx context.Context, userID string) (*entity.OTP, error) {
	const query = `
		SELECT id, user_id, otp_code, status, created_at, expires_at, validated_at
		FROM otps
		WHERE user_id = ?
		ORDER BY created_at DESC
		LIMIT 1
	`

	var otpRow otpRow
	if err := getExecutor(ctx, o.db).GetContext(ctx, &otpRow, query, userID); err != nil {
		// Check if the error is sql.ErrNoRows to return entity.ErrOTPNotFound
		if err == sql.ErrNoRows {
			return nil, entity.ErrOTPNotFound
		}
		return nil, err
	}

	return otpRow.ToEntity(), nil
}
