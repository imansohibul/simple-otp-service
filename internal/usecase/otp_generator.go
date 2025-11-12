package usecase

import (
	"crypto/rand"
	"fmt"
	"math/big"
)

type entityOTPGenerator struct{}

// GenerateOTP generates a secure 6-digit numeric OTP using crypto/rand.
// It produces values from 000000 to 999999 and guarantees leading zeros.
// Returns the generated OTP string or an error if random generation fails.
func (g *entityOTPGenerator) Generate() (string, error) {
	// rand.Int generates a cryptographically secure random number
	// in the range [0, 1,000,000).
	n, err := rand.Int(rand.Reader, big.NewInt(1000000))
	if err != nil {
		return "", err
	}

	// Format the number as 6 digits, padding with leading zeros if needed
	return fmt.Sprintf("%06d", n.Int64()), nil
}

func NewOTPGenerator() OTPGenerator {
	return &entityOTPGenerator{}
}
