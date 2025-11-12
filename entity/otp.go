package entity

import (
	"time"
)

// OTPStatus represents the current status of an OTP (One-Time Password).
type OTPStatus int8

const (
	// OTPStatusCreated means the OTP has been generated but not yet validated.
	OTPStatusCreated OTPStatus = iota + 1
	// OTPStatusValidated means the OTP has been successfully validated.
	OTPStatusValidated
	// OTPStatusExpired means the OTP has expired and can no longer be used.
	OTPStatusExpired
)

// String returns the string representation of OTPStatus.
func (o OTPStatus) String() string {
	statusToStringMap := map[OTPStatus]string{
		OTPStatusCreated:   "created",
		OTPStatusValidated: "validated",
		OTPStatusExpired:   "expired",
	}

	str, _ := statusToStringMap[o]
	return str
}

// OTP represents a one-time password (OTP)
type OTP struct {
	ID          uint64
	UserID      string
	OTPCode     string
	Status      OTPStatus
	CreatedAt   time.Time
	ExpiresAt   time.Time
	ValidatedAt *time.Time
}
