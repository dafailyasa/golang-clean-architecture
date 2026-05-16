package mysql

import (
	"context"
	"errors"

	"gorm.io/gorm"
)

type transactionKey struct{}

type TransactionManager struct {
	db *gorm.DB
}

func NewTransactionManager(db *gorm.DB) (*TransactionManager, error) {
	if db == nil {
		return nil, errors.New("failed to initialize transaction manager is empty")
	}
	return &TransactionManager{db: db}, nil
}

func (tm *TransactionManager) WithTransaction(ctx context.Context, fn func(ctx context.Context) error) error {
	tx := tm.db.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	ctx = context.WithValue(ctx, transactionKey{}, tx)

	if err := fn(ctx); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func GetDB(ctx context.Context, db *gorm.DB) *gorm.DB {
	if tx, ok := ctx.Value(transactionKey{}).(*gorm.DB); ok {
		return tx
	}
	return db
}
