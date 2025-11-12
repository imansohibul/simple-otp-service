package repository

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
)

// ExecerContext defines the interface for executing SQL commands.
type ExecerContext interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	GetContext(ctx context.Context, dest any, query string, args ...any) error
	SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
}

// txKey is a key used to store the transaction in the context
type txKey struct{}

// transactionManager is implementation of TransactionManager interface
type transactionManager struct {
	db *sqlx.DB
}

// NewTransactionManger creates a new instance of transactionManager
func NewTransactionManager(db *sqlx.DB) *transactionManager {
	return &transactionManager{db: db}
}

// WithTransaction starts a new transaction and executes the provided function within that transaction.
func (t *transactionManager) WithTransaction(ctx context.Context, fn func(ctx context.Context) error) error {
	tx, err := t.db.BeginTxx(ctx, &sql.TxOptions{})
	if err != nil {
		return err
	}

	defer tx.Rollback() // Rollback the transaction if not committed

	ctx = context.WithValue(ctx, txKey{}, tx)
	if err := fn(ctx); err != nil {
		return err
	}

	return tx.Commit()
}

// transactionFromContext retrieves the transaction from the context
func transactionFromContext(ctx context.Context) *sqlx.Tx {
	tx, _ := ctx.Value(txKey{}).(*sqlx.Tx)
	return tx
}

// getExecutor retrieves the executor from the context
// If a transaction is present in the context, it returns the transaction executor.
// Otherwise, it returns the database executor.
func getExecutor(ctx context.Context, db *sqlx.DB) ExecerContext {
	if tx := transactionFromContext(ctx); tx != nil {
		return tx
	}

	return db
}
