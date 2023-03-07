package postgres

import (
	"context"
	"log"

	"github.com/phantranhieunhan/s3-assignment/common"
	"github.com/phantranhieunhan/s3-assignment/pkg/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Database struct {
	db *gorm.DB
}

func NewDatabase() Database {
	dsn := config.C.PostgresULR
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatalln(err)
	}

	db.Debug()

	return Database{
		db: db,
	}
}

func NewDatabaseWithDB(db *gorm.DB) Database {
	return Database{
		db: db,
	}
}

type txKey struct{}

// injectTx injects transaction to context
func injectTx(ctx context.Context, tx *gorm.DB) context.Context {
	return context.WithValue(ctx, txKey{}, tx)
}

// extractTx extracts transaction from context
func extractTx(ctx context.Context) *gorm.DB {
	if tx, ok := ctx.Value(txKey{}).(*gorm.DB); ok {
		return tx
	}
	return nil
}

// model returns query model with context with or without transaction extracted from context
func (db Database) Model(ctx context.Context) *gorm.DB {
	tx := extractTx(ctx)
	if tx != nil {
		return tx
	}

	return db.db
}

// WithinTransaction runs function within transaction
//
// The transaction commits when function were finished without error
func (db Database) WithinTransaction(ctx context.Context, tFunc func(ctx context.Context) error) error {
	// begin transaction
	tx := db.db.Begin()
	if tx.Error != nil {
		return common.ErrDB(tx.Error)
	}

	defer func() {
		// finalize transaction on panic, etc.
		if r := recover(); r != nil {
			if errTx := tx.Rollback().Error; errTx != nil {
				log.Printf("close transaction: %v", errTx)
			}
		}
	}()

	// run callback
	err := tFunc(injectTx(ctx, tx))
	if err != nil {
		// if error, rollback
		if errRollback := tx.Rollback().Error; errRollback != nil {
			log.Printf("rollback transaction: %v", errRollback)
			return common.ErrDB(errRollback)
		}
		return err
	}
	// if no error, commit
	if errCommit := tx.Commit().Error; errCommit != nil {
		log.Printf("commit transaction: %v", errCommit)
		return common.ErrDB(errCommit)
	}
	return nil
}
