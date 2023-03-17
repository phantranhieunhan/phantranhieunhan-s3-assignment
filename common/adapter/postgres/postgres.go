package postgres

import (
	"context"
	"database/sql"
	"log"

	_ "github.com/lib/pq"
	"github.com/phantranhieunhan/s3-assignment/common"
	"github.com/phantranhieunhan/s3-assignment/pkg/config"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

type Database struct {
	DB *sql.DB
}

func NewDatabase() Database {
	dsn := config.C.PostgresULR
	db, err := sql.Open("postgres", dsn)

	if err != nil {
		log.Fatalln(err)
	}

	// boil.SetDB(db)
	boil.DebugMode = true

	return Database{DB: db}
}

func NewDatabaseWithDB(db *sql.DB) Database {
	return Database{
		DB: db,
	}
}

type txKey struct{}

// injectTx injects transaction to context
func injectTx(ctx context.Context, tx *sql.Tx) context.Context {
	return context.WithValue(ctx, txKey{}, tx)
}

// extractTx extracts transaction from context
func extractTx(ctx context.Context) *sql.Tx {
	if tx, ok := ctx.Value(txKey{}).(*sql.Tx); ok {
		return tx
	}
	return nil
}

// model returns query model with context with or without transaction extracted from context
// func (db Database) Model(ctx context.Context) boil {
// 	tx := extractTx(ctx)
// 	if tx != nil {
// 		return tx
// 	}

// 	return db.db
// }

// WithinTransaction runs function within transaction
//
// The transaction commits when function were finished without error
func (db Database) WithinTransaction(ctx context.Context, tFunc func(ctx context.Context) error) error {
	// begin transaction
	tx, err := db.DB.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		return common.ErrDB(err)
	}

	defer func() {
		// finalize transaction on panic, etc.
		if r := recover(); r != nil {
			if errTx := tx.Rollback(); errTx != nil {
				log.Printf("close transaction: %v", errTx)
			}
		}
	}()

	// run callback
	err = tFunc(injectTx(ctx, tx))
	if err != nil {
		// if error, rollback
		if errRollback := tx.Rollback(); errRollback != nil {
			log.Printf("rollback transaction: %v", errRollback)
			return common.ErrDB(errRollback)
		}
		return err
	}
	// if no error, commit
	if errCommit := tx.Commit(); errCommit != nil {
		log.Printf("commit transaction: %v", errCommit)
		return common.ErrDB(errCommit)
	}
	return nil
}
