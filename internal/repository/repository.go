package repository

import (
	"errors"

	"github.com/go-sql-driver/mysql"
)

// Note: Repositories are separated by domain to follow the Single Responsibility Principle.
// This improves maintainability, testability, and scalability by isolating each domain’s data access logic.
// Each repository can evolve independently, e.g., changing queries or database structure without affecting others.
// It also simplifies mocking in unit tests by allowing targeted mocks per domain, leading to more focused and reliable test cases.

// isUniqueConstraintViolation checks if the error is a unique constraint violation
func isUniqueConstraintViolation(err error) bool {
 var mysqlErr *mysql.MySQLError
    if errors.As(err, &mysqlErr) {
        return mysqlErr.Number == 1062 // ER_DUP_ENTRY
    }
    return false
}
