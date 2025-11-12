package usecase

//go:generate mockgen -destination=mock/usecase.go -package=mock -source=usecase.go

type OTPGenerator interface {
	// GenerateOTP generates a secure 6-digit numeric OTP using crypto/rand.
	// It produces values from 000000 to 999999 and guarantees leading zeros.
	// Returns the generated OTP string or an error if random generation fails.
	Generate() (string, error)
}
