package entity

import "fmt"

// DomainError represents a custom error with a Code and Message
type DomainError struct {
	Code    string
	Message string
}

func (e DomainError) Error() string {
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

// DomainError creates a new AppError instance
func NewDomainError(code, message string) *DomainError {
	return &DomainError{Code: code, Message: message}
}

var (
	// GENERAL errors
	ErrInvalidRequest = NewDomainError("invalid_request", "Invalid request: Please check the request body and try again")

	// OTP specific errors
	ErrOTPExpired           = NewDomainError("otp_expired", "OTP has expired")
	ErrOTPUsed              = NewDomainError("otp_used", "OTP has already been used")
	ErrOTPNotFound          = NewDomainError("otp_not_found", "OTP Not Found")
	ErrOTPDuplicate         = NewDomainError("duplicate_otp_code", "OTP Code Already Exists")
	ErrOTPRateLimitExceeded = NewDomainError("otp_rete_limit_exceeded", "OTP requested too frequently, please wait before requesting again")
)
