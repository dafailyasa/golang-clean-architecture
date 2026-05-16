package mysql

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type TransactionManagerTest struct {
	suite.Suite
	db   *sql.DB
	mock sqlmock.Sqlmock
	gdb  *gorm.DB
}

func (t *TransactionManagerTest) SetupTest() {
	db, mock, err := sqlmock.New()
	t.Require().NoError(err)

	dialector := mysql.New(mysql.Config{
		Conn:                      db,
		SkipInitializeWithVersion: true,
	})

	gdb, err := gorm.Open(dialector, &gorm.Config{})
	t.Require().NoError(err)

	t.db = db
	t.mock = mock
	t.gdb = gdb
}

func (t *TransactionManagerTest) TearDownTest() {
	t.db.Close()
}

func TestTransactionManager(t *testing.T) {
	suite.Run(t, new(TransactionManagerTest))
}

func (t *TransactionManagerTest) TestNewTransactionManagerNilDB() {
	tm, err := NewTransactionManager(nil)
	t.Error(err)
	t.Nil(tm)
}

func (t *TransactionManagerTest) TestNewTransactionManagerSuccess() {
	tm, err := NewTransactionManager(t.gdb)
	t.NoError(err)
	t.NotNil(tm)
}

func (t *TransactionManagerTest) TestWithTransactionCommit() {
	tm, err := NewTransactionManager(t.gdb)
	t.Require().NoError(err)

	t.mock.ExpectBegin()
	t.mock.ExpectExec("SELECT 1").WillReturnResult(sqlmock.NewResult(0, 1))
	t.mock.ExpectCommit()

	err = tm.WithTransaction(context.Background(), func(ctx context.Context) error {
		tx := GetDB(ctx, t.gdb)
		return tx.Exec("SELECT 1").Error
	})

	t.NoError(err)
	t.NoError(t.mock.ExpectationsWereMet())
}

func (t *TransactionManagerTest) TestWithTransactionRollbackOnError() {
	tm, err := NewTransactionManager(t.gdb)
	t.Require().NoError(err)

	t.mock.ExpectBegin()
	t.mock.ExpectRollback()

	err = tm.WithTransaction(context.Background(), func(ctx context.Context) error {
		return errors.New("fail")
	})

	t.Error(err)
	t.NoError(t.mock.ExpectationsWereMet())
}

func (t *TransactionManagerTest) TestGetDBReturnsBaseWithoutTx() {
	db := GetDB(context.Background(), t.gdb)
	t.Equal(t.gdb, db)
}
