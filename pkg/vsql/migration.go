package vsql

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/vertex-center/vertex/pkg/log"
	"github.com/vertex-center/vlog"
)

type Migration interface {
	Up(tx *sqlx.Tx) error
}

// Migrate executes all migrations that have not been executed yet.
// This function needs a migration table containing a version field.
func Migrate(migrations []Migration, db *sqlx.DB, current int) error {
	target := len(migrations)
	for i := current; i < target; i++ {
		log.Info("executing migration", vlog.Int("version", i+1))

		tx, err := db.Beginx()
		if err != nil {
			return fmt.Errorf("failed to start migration transaction: %w", err)
		}

		log.Info("migration transaction started", vlog.Int("version", i+1))

		err = migrations[i].Up(tx)
		if err != nil {
			_ = tx.Rollback()
			return fmt.Errorf("failed to execute migration: %w", err)
		}

		_, err = tx.Exec("UPDATE migrations SET version = $1 WHERE id = 1", i+1)
		if err != nil {
			_ = tx.Rollback()
			return fmt.Errorf("failed to update migration version: %w", err)
		}

		err = tx.Commit()
		if err != nil {
			return fmt.Errorf("failed to commit migration transaction: %w", err)
		}

		log.Info("migration transaction committed", vlog.Int("version", i+1))
	}
	return nil
}
