package repository_test

import (
	"context"
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/imansohibul/otp-service/internal/repository"
	"github.com/stretchr/testify/assert"
)

func TestTransactionManager_WithTransaction(t *testing.T) {
	tests := []struct {
		name          string
		mockSetup     func(mock sqlmock.Sqlmock)
		txFunc        func(ctx context.Context) error
		expectedError error
	}{
		{
			name: "should successfully commit the database transaction",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectCommit()
			},
			txFunc: func(ctx context.Context) error {
				return nil
			},
			expectedError: nil,
		},
		{
			name: "should rollback transaction when error occurs",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectRollback()
			},
			txFunc: func(ctx context.Context) error {
				return sql.ErrTxDone
			},
			expectedError: sql.ErrTxDone,
		},
		{
			name: "should return error when transaction begins",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin().WillReturnError(sql.ErrConnDone)
			},
			txFunc: func(ctx context.Context) error {
				return nil
			},
			expectedError: sql.ErrConnDone,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			repoDependency := newRepoDependency()
			txManager := repository.NewTransactionManager(repoDependency.mockedDB)

			tt.mockSetup(repoDependency.mockedSQL)

			err := txManager.WithTransaction(ctx, tt.txFunc)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
