package repository_test

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

type repositoryDependency struct {
	mockedDB  *sqlx.DB
	mockedSQL sqlmock.Sqlmock
}

func newRepoDependency() *repositoryDependency {
	mockDB, sqlMock, _ := sqlmock.New()
	sqlxDB := sqlx.NewDb(mockDB, "sqlmock")
	return &repositoryDependency{
		mockedDB:  sqlxDB,
		mockedSQL: sqlMock,
	}
}

func TestRepository(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Repository Suite")
}
