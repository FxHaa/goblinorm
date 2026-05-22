package migration

import (
	"database/sql"
	"fmt"
	"goblinorm/internal/goblinerr"

	"goblinorm/internal/logger"
	"goblinorm/schema"
)

type Dialect interface {
	CreateTableSQL(model any) (string, error)
}

func RaiseTheDead(db *sql.DB, dialect Dialect, log *logger.Logger, model any) error {
	parsed, err := schema.Parse(model)
	if err != nil {
		if log != nil {
			log.SpellFailed(err)
		}
		return err
	}

	if log != nil {
		log.Raising(parsed.Name)
	}

	sqlText, err := dialect.CreateTableSQL(model)
	if err != nil {
		if log != nil {
			log.SpellFailed(err)
		}
		return err
	}

	if log != nil {
		log.ForbiddenSQL(sqlText)
	}

	_, err = db.Exec(sqlText)
	if err != nil {
		err = fmt.Errorf("%w: %v", goblinerr.ErrMigrationFailed, err)
		if log != nil {
			log.SpellFailed(err)
		}
		return err
	}

	if log != nil {
		log.MigrationSuccessful()
	}

	return nil
}
