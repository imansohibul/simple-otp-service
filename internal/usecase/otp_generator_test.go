package usecase_test

import (
	"crypto/rand"
	"errors"
	"testing"

	"github.com/imansohibul/otp-service/internal/usecase"
	"github.com/stretchr/testify/assert"
)

// MockReader is a reader that always returns an error
type MockReader struct{}

func (r *MockReader) Read(p []byte) (int, error) {
	return 0, errors.New("mock rand error")
}

func TestEntityOTPGenerator_Generate(t *testing.T) {
	t.Run("should generate 6-digit OTP successfully", func(t *testing.T) {
		gen := usecase.NewOTPGenerator()

		otp, err := gen.Generate()
		assert.NoError(t, err)
		assert.Len(t, otp, 6)
		// Ensure all characters are digits
		for _, ch := range otp {
			assert.GreaterOrEqual(t, ch, '0')
			assert.LessOrEqual(t, ch, '9')
		}
	})

	t.Run("should return error when rand.Reader fails", func(t *testing.T) {
		// Temporarily replace rand.Reader with a reader that always fails
		oldReader := rand.Reader
		rand.Reader = &MockReader{}
		defer func() { rand.Reader = oldReader }() // restore after test

		gen := usecase.NewOTPGenerator()

		otp, err := gen.Generate()
		assert.Empty(t, otp)
		assert.Error(t, err)
		assert.EqualError(t, err, "mock rand error")
	})
}
