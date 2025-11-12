// This file contains types that are used in the repository layer.
package repository

import (
	"time"

	"github.com/imansohibul/otp-service/entity"
)

// Note: Why not use the entity types directly?
// Defining separate types in the repository layer provides clear separation of concerns.
// The repository struct maps directly to the database schema and may include tags (e.g., db tags) specific to the database driver.
// This allows flexibility in how data is stored or queried, without coupling the repository to business-layer concerns.
// It also makes it easier to switch database drivers or ORMs, as changes in the data access layer won’t leak into the domain layer.

// otpRow represents the OTP table row structure for database operations
type otpRow struct {
	ID          uint64     `db:"id"`
	UserID      string     `db:"user_id"`
	OTPCode     string     `db:"otp_code"`
	Status      int        `db:"status"`
	CreatedAt   time.Time  `db:"created_at"`
	ExpiresAt   time.Time  `db:"expires_at"`
	ValidatedAt *time.Time `db:"validated_at"` // Nullable field
}

// ToEntity converts otpRow to entity.OTP
func (r *otpRow) ToEntity() *entity.OTP {
	return &entity.OTP{
		ID:          r.ID,
		UserID:      r.UserID,
		OTPCode:     r.OTPCode,
		Status:      entity.OTPStatus(r.Status),
		CreatedAt:   r.CreatedAt,
		ExpiresAt:   r.ExpiresAt,
		ValidatedAt: r.ValidatedAt,
	}
}

// QueryOption type to represent query modifiers
type QueryOption string

// Define available options
const (
	WithForUpdate QueryOption = "FOR UPDATE"
)
