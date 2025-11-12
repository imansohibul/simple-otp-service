package usecase_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/imansohibul/otp-service/entity"
	"github.com/imansohibul/otp-service/internal/usecase"
	"github.com/imansohibul/otp-service/internal/usecase/mock"
	"github.com/stretchr/testify/assert"
)

func TestOtpUsecase_Create(t *testing.T) {
	type useCaseDependency struct {
		otpRepo      *mock.MockOTPRepository
		otpGenerator *mock.MockOTPGenerator
	}

	tests := []struct {
		name           string
		userID         string
		mockDependency func(dep *useCaseDependency)
		assertFn       func(*entity.OTP, error)
	}{
		{
			name:   "should return error when repository Create fails",
			userID: "user-1",
			mockDependency: func(dep *useCaseDependency) {
				dep.otpRepo.EXPECT().
					GetLastByUserID(gomock.Any(), "user-1").
					Return(nil, nil)
				dep.otpGenerator.EXPECT().
					Generate().
					Return("123456", nil)
				dep.otpRepo.EXPECT().
					Create(gomock.Any(), gomock.Any()).
					Return(errors.New("db error"))
			},
			assertFn: func(otp *entity.OTP, err error) {
				assert.Nil(t, otp)
				assert.NotNil(t, err)
				assert.EqualError(t, err, "db error")
			},
		},
		{
			name:   "should create otp successfully",
			userID: "user-1",
			mockDependency: func(dep *useCaseDependency) {
				dep.otpRepo.EXPECT().
					GetLastByUserID(gomock.Any(), "user-1").
					Return(nil, nil)
				dep.otpGenerator.EXPECT().
					Generate().
					Return("123456", nil)
				dep.otpRepo.EXPECT().
					Create(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, otp *entity.OTP) error {
						assert.Equal(t, "user-1", otp.UserID)
						assert.Equal(t, "123456", otp.OTPCode)
						assert.Equal(t, entity.OTPStatusCreated, otp.Status)
						assert.WithinDuration(t, time.Now().Add(2*time.Minute), otp.ExpiresAt, 2*time.Second)
						return nil
					})
			},
			assertFn: func(otp *entity.OTP, err error) {
				assert.NotNil(t, otp)
				assert.Nil(t, err)
				assert.Equal(t, entity.OTPStatusCreated, otp.Status)
				assert.Equal(t, "123456", otp.OTPCode)
			},
		},
		{
			name:   "should return rate limit error if OTP requested too soon",
			userID: "user-1",
			mockDependency: func(dep *useCaseDependency) {
				dep.otpRepo.EXPECT().
					GetLastByUserID(gomock.Any(), "user-1").
					Return(&entity.OTP{
						UserID:    "user-1",
						OTPCode:   "654321",
						Status:    entity.OTPStatusCreated,
						CreatedAt: time.Now().Add(-1 * time.Minute),
						ExpiresAt: time.Now().Add(1 * time.Minute),
					}, nil)
			},
			assertFn: func(otp *entity.OTP, err error) {
				assert.Nil(t, otp)
				assert.NotNil(t, err)
				assert.Equal(t, entity.ErrOTPRateLimitExceeded, err)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			dep := useCaseDependency{
				otpRepo:      mock.NewMockOTPRepository(ctrl),
				otpGenerator: mock.NewMockOTPGenerator(ctrl),
			}

			tt.mockDependency(&dep)

			usc := usecase.NewOtpUsecase(dep.otpRepo, dep.otpGenerator)

			otp, err := usc.Create(context.Background(), tt.userID)

			tt.assertFn(otp, err)
		})
	}
}

func TestOtpUsecase_Validate(t *testing.T) {
	type useCaseDependency struct {
		otpRepo *mock.MockOTPRepository
	}

	userID := "user-1"

	tests := []struct {
		name           string
		otpCode        string
		mockDependency func(dep *useCaseDependency)
		assertFn       func(*entity.OTP, error)
	}{
		{
			name:    "should return error if OTP not found",
			otpCode: "000000",
			mockDependency: func(dep *useCaseDependency) {
				dep.otpRepo.EXPECT().
					FindByUserIDAndCode(gomock.Any(), userID, "000000").
					Return(nil, entity.ErrOTPNotFound)
			},
			assertFn: func(otp *entity.OTP, err error) {
				assert.Nil(t, otp)
				assert.Equal(t, entity.ErrOTPNotFound, err)
			},
		},
		{
			name:    "should return error if OTP already validated",
			otpCode: "111111",
			mockDependency: func(dep *useCaseDependency) {
				dep.otpRepo.EXPECT().
					FindByUserIDAndCode(gomock.Any(), userID, "111111").
					Return(&entity.OTP{
						UserID:  userID,
						OTPCode: "111111",
						Status:  entity.OTPStatusValidated,
					}, nil)
			},
			assertFn: func(otp *entity.OTP, err error) {
				assert.Nil(t, otp)
				assert.Equal(t, entity.ErrOTPUsed, err)
			},
		},
		{
			name:    "should return error if OTP expired",
			otpCode: "222222",
			mockDependency: func(dep *useCaseDependency) {
				otp := &entity.OTP{
					UserID:    userID,
					OTPCode:   "222222",
					Status:    entity.OTPStatusCreated,
					ExpiresAt: time.Now().Add(-1 * time.Minute),
				}
				dep.otpRepo.EXPECT().
					FindByUserIDAndCode(gomock.Any(), userID, "222222").
					Return(otp, nil)
				dep.otpRepo.EXPECT().
					Update(gomock.Any(), otp).
					Return(nil) // update status to expired
			},
			assertFn: func(otp *entity.OTP, err error) {
				assert.Nil(t, otp)
				assert.Equal(t, entity.ErrOTPExpired, err)
			},
		},
		{
			name:    "should validate OTP successfully",
			otpCode: "333333",
			mockDependency: func(dep *useCaseDependency) {
				otp := &entity.OTP{
					UserID:    userID,
					OTPCode:   "333333",
					Status:    entity.OTPStatusCreated,
					ExpiresAt: time.Now().Add(1 * time.Minute),
				}
				dep.otpRepo.EXPECT().
					FindByUserIDAndCode(gomock.Any(), userID, "333333").
					Return(otp, nil)
				dep.otpRepo.EXPECT().
					Update(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, updatedOTP *entity.OTP) error {
						assert.Equal(t, entity.OTPStatusValidated, updatedOTP.Status)
						assert.NotNil(t, updatedOTP.ValidatedAt)
						return nil
					})
			},
			assertFn: func(otp *entity.OTP, err error) {
				assert.NotNil(t, otp)
				assert.Nil(t, err)
				assert.Equal(t, entity.OTPStatusValidated, otp.Status)
				assert.NotNil(t, otp.ValidatedAt)
			},
		},
		{
			name:    "should return error if update fails when validating",
			otpCode: "444444",
			mockDependency: func(dep *useCaseDependency) {
				otp := &entity.OTP{
					UserID:    userID,
					OTPCode:   "444444",
					Status:    entity.OTPStatusCreated,
					ExpiresAt: time.Now().Add(1 * time.Minute),
				}
				dep.otpRepo.EXPECT().
					FindByUserIDAndCode(gomock.Any(), userID, "444444").
					Return(otp, nil)
				dep.otpRepo.EXPECT().
					Update(gomock.Any(), gomock.Any()).
					Return(errors.New("db update failed"))
			},
			assertFn: func(otp *entity.OTP, err error) {
				assert.Nil(t, otp)
				assert.NotNil(t, err)
				assert.Contains(t, err.Error(), "failed to update OTP status")
			},
		},
		{
			name:    "should return error if update fails when marking expired OTP",
			otpCode: "555555",
			mockDependency: func(dep *useCaseDependency) {
				// Simulate an OTP that is expired and still in Created status
				otp := &entity.OTP{
					UserID:    userID,
					OTPCode:   "555555",
					Status:    entity.OTPStatusCreated,
					ExpiresAt: time.Now().Add(-1 * time.Minute), // expired
				}
				dep.otpRepo.EXPECT().
					FindByUserIDAndCode(gomock.Any(), userID, "555555").
					Return(otp, nil)
				// Simulate update failure when trying to mark it as expired
				dep.otpRepo.EXPECT().
					Update(gomock.Any(), otp).
					Return(errors.New("db update failed"))
			},
			assertFn: func(otp *entity.OTP, err error) {
				assert.Nil(t, otp)
				assert.NotNil(t, err)
				assert.EqualError(t, err, "db update failed") // error comes directly from repository
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			dep := useCaseDependency{
				otpRepo: mock.NewMockOTPRepository(ctrl),
			}

			tt.mockDependency(&dep)

			usc := usecase.NewOtpUsecase(dep.otpRepo, nil)

			otp, err := usc.Validate(context.Background(), userID, tt.otpCode)

			tt.assertFn(otp, err)
		})
	}
}
