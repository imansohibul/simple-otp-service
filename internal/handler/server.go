package handler

import (
	"context"
	"net/http"

	"github.com/imansohibul/otp-service/generated"
	intmiddleware "github.com/imansohibul/otp-service/internal/handler/middleware"
	"github.com/labstack/echo-contrib/echoprometheus"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	oapimiddleware "github.com/oapi-codegen/echo-middleware"
	"github.com/rs/zerolog/log"
)

// RestServer encapsulates the Echo instance and usecases
type RestAPIServer struct {
	Echo       *echo.Echo
	OtpUsecase OTPUsecase
}

// NewRestAPIServer constructs the server with injected usecases
func NewRestAPIServer(otpUsecase OTPUsecase) *RestAPIServer {
	var (
		e      = echo.New()
		server = &RestAPIServer{
			Echo:       e,
			OtpUsecase: otpUsecase,
		}
	)

	// Set up middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())
	e.Use(middleware.RequestID())
	e.Use(echoprometheus.NewMiddleware("otp-service")) // adds middleware to gather metrics

	spec, err := generated.GetSwagger()
	if err != nil {
		log.Fatal().Err(err).Msg("REST API server stopped with error")
	}

	e.Use(oapimiddleware.OapiRequestValidatorWithOptions(spec, &oapimiddleware.Options{ErrorHandler: validationErrorHandler}))

	e.GET("/metrics", echoprometheus.NewHandler()) // adds route to serve gathered metrics
	e.HTTPErrorHandler = intmiddleware.ErrorHandler

	v1 := e.Group("/api/v1")
	generated.RegisterHandlers(v1, server)

	return server
}

// Start launches the Echo HTTP server
func (s *RestAPIServer) Start(address string) error {
	return s.Echo.Start(address)
}

// Shutdown gracefully shuts down the server
// It waits for all active connections to finish before closing
func (s *RestAPIServer) Shutdown(ctx context.Context) error {
	return s.Echo.Shutdown(ctx)
}

// validationErrorHandler handles OpenAPI validation errors and returns 400 Bad Request
func validationErrorHandler(c echo.Context, err *echo.HTTPError) error {
	// Log the validation error
	log.Warn().
		Err(err).
		Str("path", c.Request().URL.Path).
		Str("method", c.Request().Method).
		Msg("Request validation failed")

	// Return 400 Bad Request with error details
	return c.JSON(http.StatusBadRequest, map[string]interface{}{
		"error":             "bad_request",
		"error_description": err.Message,
	})
}
