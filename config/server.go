package config

import (
	"github.com/imansohibul/otp-service/internal/handler"
	"github.com/imansohibul/otp-service/internal/repository"
	"github.com/imansohibul/otp-service/internal/usecase"
)

func NewRestAPI() (*handler.RestAPIServer, error) {
	// Load configuration
	serviceConfig, err := LoadConfig()
	if err != nil {
		return nil, err
	}

	// Initialize database connection
	db := initDatabase(serviceConfig)

	// Initialize repositories
	var (
		otpRepository = repository.NewOTPRepository(db)
	)

	// Create usecases
	var (
		otpGenerator = usecase.NewOTPGenerator()
		otpUsecase   = usecase.NewOtpUsecase(
			otpRepository,
			otpGenerator,
		)
	)

	// Initialize Rest API server
	return handler.NewRestAPIServer(otpUsecase), nil
}
