package repository_test

import (
	"context"
	"database/sql"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-sql-driver/mysql"
	"github.com/golang/mock/gomock"
	"github.com/imansohibul/otp-service/entity"
	"github.com/imansohibul/otp-service/internal/repository"
	"github.com/stretchr/testify/assert"
)

func TestOTPRepository_Create(t *testing.T) {
	type Input struct {
		ctx context.Context
		otp *entity.OTP
	}

	now := time.Now()
	expiresAt := now.Add(5 * time.Minute)

	dummyOTP := entity.OTP{
		UserID:    "user123",
		OTPCode:   "123456",
		Status:    entity.OTPStatusCreated,
		ExpiresAt: expiresAt,
	}

	expectedQuery := regexp.QuoteMeta("INSERT INTO otps (user_id, otp_code, status, expires_at) VALUES (?, ?, ?, ?)")

	tests := []struct {
		name           string
		mockDependency func(*repositoryDependency)
		assertFn       func(error)
		input          Input
	}{
		{
			name: "Should successfully create a new OTP",
			input: Input{
				ctx: context.TODO(),
				otp: &dummyOTP,
			},
			mockDependency: func(dependency *repositoryDependency) {
				dependency.mockedSQL.
					ExpectExec(expectedQuery).
					WithArgs("user123", "123456", entity.OTPStatusCreated, expiresAt).
					WillReturnResult(sqlmock.NewResult(1, 1)).
					WillReturnError(nil)
			},
			assertFn: func(err error) {
				assert.Nil(t, err)
			},
		},
		{
			name: "Should successfully create OTP with different values",
			input: Input{
				ctx: context.TODO(),
				otp: &entity.OTP{
					UserID:    "user456",
					OTPCode:   "654321",
					Status:    entity.OTPStatusCreated,
					ExpiresAt: expiresAt,
				},
			},
			mockDependency: func(dependency *repositoryDependency) {
				dependency.mockedSQL.
					ExpectExec(expectedQuery).
					WithArgs("user456", "654321", entity.OTPStatusCreated, expiresAt).
					WillReturnResult(sqlmock.NewResult(2, 1)).
					WillReturnError(nil)
			},
			assertFn: func(err error) {
				assert.Nil(t, err)
			},
		},
		{
			name: "Should return duplicate error when unique constraint violated",
			input: Input{
				ctx: context.TODO(),
				otp: &dummyOTP,
			},
			mockDependency: func(dependency *repositoryDependency) {
				dependency.mockedSQL.
					ExpectExec(expectedQuery).
					WithArgs("user123", "123456", entity.OTPStatusCreated, expiresAt).
					WillReturnError(&mysql.MySQLError{
						Number:  1062,
						Message: "Duplicate entry 'user123-123456' for key 'unique_user_otp'",
					})
			},
			assertFn: func(err error) {
				assert.NotNil(t, err)
				assert.Equal(t, entity.ErrOTPDuplicate, err)
			},
		},
		{
			name: "Should return error when exec context fails",
			input: Input{
				ctx: context.TODO(),
				otp: &dummyOTP,
			},
			mockDependency: func(dependency *repositoryDependency) {
				dependency.mockedSQL.
					ExpectExec(expectedQuery).
					WithArgs("user123", "123456", entity.OTPStatusCreated, expiresAt).
					WillReturnError(sqlmock.ErrCancelled)
			},
			assertFn: func(err error) {
				assert.NotNil(t, err)
				assert.Equal(t, sqlmock.ErrCancelled, err)
			},
		},
		{
			name: "Should return error when database connection fails",
			input: Input{
				ctx: context.TODO(),
				otp: &dummyOTP,
			},
			mockDependency: func(dependency *repositoryDependency) {
				dependency.mockedSQL.
					ExpectExec(expectedQuery).
					WithArgs("user123", "123456", entity.OTPStatusCreated, expiresAt).
					WillReturnError(sql.ErrConnDone)
			},
			assertFn: func(err error) {
				assert.NotNil(t, err)
				assert.Equal(t, sql.ErrConnDone, err)
			},
		},
		{
			name: "Should return error when transaction fails",
			input: Input{
				ctx: context.TODO(),
				otp: &dummyOTP,
			},
			mockDependency: func(dependency *repositoryDependency) {
				dependency.mockedSQL.
					ExpectExec(expectedQuery).
					WithArgs("user123", "123456", entity.OTPStatusCreated, expiresAt).
					WillReturnError(sql.ErrTxDone)
			},
			assertFn: func(err error) {
				assert.NotNil(t, err)
				assert.Equal(t, sql.ErrTxDone, err)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var (
				ctrl                 = gomock.NewController(t)
				repositoryDependency = newRepoDependency()
				repo                 = repository.NewOTPRepository(repositoryDependency.mockedDB)
			)

			defer ctrl.Finish()
			defer repositoryDependency.mockedDB.Close()

			tt.mockDependency(repositoryDependency)
			tt.assertFn(repo.Create(tt.input.ctx, tt.input.otp))

			// Verify all expectations were met
			assert.NoError(t, repositoryDependency.mockedSQL.ExpectationsWereMet())
		})
	}
}

func TestOTPRepository_FindByUserIDAndCode(t *testing.T) {
	type Input struct {
		ctx     context.Context
		userID  string
		otpCode string
	}

	now := time.Now()
	expectedQuery := regexp.QuoteMeta(`
		SELECT id, user_id, otp_code, status, created_at, expires_at, validated_at
		FROM otps
		WHERE user_id = ? AND otp_code = ?
	`)

	tests := []struct {
		name           string
		mockDependency func(*repositoryDependency)
		assertFn       func(*testing.T, *entity.OTP, error)
		input          Input
	}{
		{
			name: "Should return OTP successfully",
			input: Input{
				ctx:     context.TODO(),
				userID:  "user123",
				otpCode: "123456",
			},
			mockDependency: func(dependency *repositoryDependency) {
				dependency.mockedSQL.
					ExpectQuery(expectedQuery).
					WithArgs("user123", "123456").
					WillReturnRows(sqlmock.NewRows([]string{
						"id", "user_id", "otp_code", "status", "created_at", "expires_at", "validated_at",
					}).AddRow(
						1, "user123", "123456", entity.OTPStatusCreated, now, now.Add(5*time.Minute), nil,
					))
			},
			assertFn: func(t *testing.T, otp *entity.OTP, err error) {
				assert.Nil(t, err)
				assert.NotNil(t, otp)
				assert.Equal(t, "user123", otp.UserID)
				assert.Equal(t, "123456", otp.OTPCode)
			},
		},
		{
			name: "Should return ErrOTPNotFound when no row found",
			input: Input{
				ctx:     context.TODO(),
				userID:  "user123",
				otpCode: "000000",
			},
			mockDependency: func(dependency *repositoryDependency) {
				dependency.mockedSQL.
					ExpectQuery(expectedQuery).
					WithArgs("user123", "000000").
					WillReturnError(sql.ErrNoRows)
			},
			assertFn: func(t *testing.T, otp *entity.OTP, err error) {
				assert.Nil(t, otp)
				assert.Equal(t, entity.ErrOTPNotFound, err)
			},
		},
		{
			name: "Should return error when DB fails",
			input: Input{
				ctx:     context.TODO(),
				userID:  "user123",
				otpCode: "123456",
			},
			mockDependency: func(dependency *repositoryDependency) {
				dependency.mockedSQL.
					ExpectQuery(expectedQuery).
					WithArgs("user123", "123456").
					WillReturnError(sql.ErrConnDone)
			},
			assertFn: func(t *testing.T, otp *entity.OTP, err error) {
				assert.Nil(t, otp)
				assert.Equal(t, sql.ErrConnDone, err)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repositoryDependency := newRepoDependency()
			repo := repository.NewOTPRepository(repositoryDependency.mockedDB)

			defer repositoryDependency.mockedDB.Close()

			tt.mockDependency(repositoryDependency)
			otp, err := repo.FindByUserIDAndCode(tt.input.ctx, tt.input.userID, tt.input.otpCode)
			tt.assertFn(t, otp, err)

			assert.NoError(t, repositoryDependency.mockedSQL.ExpectationsWereMet())
		})
	}
}

func TestOTPRepository_Update(t *testing.T) {
	type Input struct {
		ctx context.Context
		otp *entity.OTP
	}

	now := time.Now()
	dummyOTP := &entity.OTP{
		ID:          1,
		UserID:      "user123",
		OTPCode:     "123456",
		Status:      entity.OTPStatusValidated,
		ValidatedAt: &now,
	}

	expectedQuery := regexp.QuoteMeta(`
		UPDATE otps
		SET status = ?, validated_at = ?
		WHERE id = ?
	`)

	tests := []struct {
		name           string
		mockDependency func(*repositoryDependency)
		assertFn       func(error)
		input          Input
	}{
		{
			name: "Should update OTP successfully",
			input: Input{
				ctx: context.TODO(),
				otp: dummyOTP,
			},
			mockDependency: func(dependency *repositoryDependency) {
				dependency.mockedSQL.
					ExpectExec(expectedQuery).
					WithArgs(entity.OTPStatusValidated, now, 1).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			assertFn: func(err error) {
				assert.Nil(t, err)
			},
		},
		{
			name: "Should return error when update fails",
			input: Input{
				ctx: context.TODO(),
				otp: dummyOTP,
			},
			mockDependency: func(dependency *repositoryDependency) {
				dependency.mockedSQL.
					ExpectExec(expectedQuery).
					WithArgs(entity.OTPStatusValidated, now, 1).
					WillReturnError(sql.ErrConnDone)
			},
			assertFn: func(err error) {
				assert.NotNil(t, err)
				assert.Equal(t, sql.ErrConnDone, err)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repositoryDependency := newRepoDependency()
			repo := repository.NewOTPRepository(repositoryDependency.mockedDB)

			defer repositoryDependency.mockedDB.Close()

			tt.mockDependency(repositoryDependency)
			err := repo.Update(tt.input.ctx, tt.input.otp)
			tt.assertFn(err)

			assert.NoError(t, repositoryDependency.mockedSQL.ExpectationsWereMet())
		})
	}
}

func TestOTPRepository_GetLastByUserID(t *testing.T) {
	type Input struct {
		ctx    context.Context
		userID string
	}

	now := time.Now()
	expectedQuery := regexp.QuoteMeta(`
		SELECT id, user_id, otp_code, status, created_at, expires_at, validated_at
		FROM otps
		WHERE user_id = ?
		ORDER BY created_at DESC
		LIMIT 1
	`)

	tests := []struct {
		name           string
		mockDependency func(*repositoryDependency)
		assertFn       func(*testing.T, *entity.OTP, error)
		input          Input
	}{
		{
			name: "Should return the most recent OTP successfully",
			input: Input{
				ctx:    context.TODO(),
				userID: "user123",
			},
			mockDependency: func(dependency *repositoryDependency) {
				dependency.mockedSQL.
					ExpectQuery(expectedQuery).
					WithArgs("user123").
					WillReturnRows(sqlmock.NewRows([]string{
						"id", "user_id", "otp_code", "status", "created_at", "expires_at", "validated_at",
					}).AddRow(
						1, "user123", "123456", entity.OTPStatusCreated, now, now.Add(2*time.Minute), nil,
					))
			},
			assertFn: func(t *testing.T, otp *entity.OTP, err error) {
				assert.Nil(t, err)
				assert.NotNil(t, otp)
				assert.Equal(t, "user123", otp.UserID)
				assert.Equal(t, "123456", otp.OTPCode)
				assert.Equal(t, entity.OTPStatusCreated, otp.Status)
			},
		},
		{
			name: "Should return ErrOTPNotFound when no OTP exists for user",
			input: Input{
				ctx:    context.TODO(),
				userID: "user999",
			},
			mockDependency: func(dependency *repositoryDependency) {
				dependency.mockedSQL.
					ExpectQuery(expectedQuery).
					WithArgs("user999").
					WillReturnError(sql.ErrNoRows)
			},
			assertFn: func(t *testing.T, otp *entity.OTP, err error) {
				assert.Nil(t, otp)
				assert.Equal(t, entity.ErrOTPNotFound, err)
			},
		},
		{
			name: "Should return error when DB fails",
			input: Input{
				ctx:    context.TODO(),
				userID: "user123",
			},
			mockDependency: func(dependency *repositoryDependency) {
				dependency.mockedSQL.
					ExpectQuery(expectedQuery).
					WithArgs("user123").
					WillReturnError(sql.ErrConnDone)
			},
			assertFn: func(t *testing.T, otp *entity.OTP, err error) {
				assert.Nil(t, otp)
				assert.Equal(t, sql.ErrConnDone, err)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repositoryDependency := newRepoDependency()
			repo := repository.NewOTPRepository(repositoryDependency.mockedDB)

			defer repositoryDependency.mockedDB.Close()

			tt.mockDependency(repositoryDependency)
			otp, err := repo.GetLastByUserID(tt.input.ctx, tt.input.userID)
			tt.assertFn(t, otp, err)

			assert.NoError(t, repositoryDependency.mockedSQL.ExpectationsWereMet())
		})
	}
}
