package query

import (
	"database/sql"
	"errors"
	"fmt"
	"goblinorm/hooks"
	"goblinorm/internal/goblinerr"
	"reflect"

	"goblinorm/internal/logger"
	"goblinorm/schema"
)

type Dialect interface {
	InsertSQL(model any) (string, []any, error)
	SelectSQL(model any, condition string) (string, error)
	DeleteSQL(model any) (string, []any, error)
}

func Summon(db *sql.DB, dialect Dialect, log *logger.Logger, model any) error {
	parsed, err := schema.Parse(model)
	if err != nil {
		if log != nil {
			log.SpellFailed(err)
		}
		return err
	}

	if err := hooks.BeforeSummon(model, log); err != nil {
		return err
	}

	if log != nil {
		log.Summoning(parsed.Name)
		log.ManaAccepted()
	}

	sqlText, args, err := dialect.InsertSQL(model)
	if err != nil {
		if log != nil {
			log.SpellFailed(err)
		}
		return err
	}

	if log != nil {
		log.ForbiddenSQL(sqlText, args...)
	}

	_, err = db.Exec(sqlText, args...)
	if err != nil {
		err = fmt.Errorf("%w: %v", goblinerr.ErrInsertFailed, err)
		if log != nil {
			log.SpellFailed(err)
		}
		return err
	}

	if log != nil {
		log.InsertSuccessful()
	}

	if err := hooks.AfterSummon(model, log); err != nil {
		return err
	}

	return nil
}

func Divine(db *sql.DB, dialect Dialect, log *logger.Logger, dest any, condition string, args ...any) error {
	parsed, err := schema.Parse(dest)
	if err != nil {
		if log != nil {
			log.SpellFailed(err)
		}
		return err
	}

	if err := hooks.BeforeDivine(dest, log); err != nil {
		return err
	}

	if log != nil {
		log.Divining(parsed.Name)
	}

	sqlText, err := dialect.SelectSQL(dest, condition)
	if err != nil {
		if log != nil {
			log.SpellFailed(err)
		}
		return err
	}

	if log != nil {
		log.ForbiddenSQL(sqlText, args...)
	}

	row := db.QueryRow(sqlText, args...)
	err = scanRow(row, dest)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = fmt.Errorf("%w: %s", goblinerr.ErrNoRows, parsed.Name)
		} else {
			err = fmt.Errorf("%w: %v", goblinerr.ErrScanFailed, err)
		}

		if log != nil {
			log.SpellFailed(err)
		}
		return err
	}

	if log != nil {
		log.QuerySuccessful()
	}

	if err := hooks.AfterDivine(dest, log); err != nil {
		return err
	}

	return nil
}

func Sacrifice(db *sql.DB, dialect Dialect, log *logger.Logger, model any) error {
	parsed, err := schema.Parse(model)
	if err != nil {
		if log != nil {
			log.SpellFailed(err)
		}
		return err
	}

	if err := hooks.BeforeSacrifice(model, log); err != nil {
		return err
	}

	if log != nil {
		log.Sacrificing(parsed.Name)
	}

	sqlText, args, err := dialect.DeleteSQL(model)
	if err != nil {
		if log != nil {
			log.SpellFailed(err)
		}
		return err
	}

	if log != nil {
		log.ForbiddenSQL(sqlText, args...)
	}

	_, err = db.Exec(sqlText, args...)
	if err != nil {
		err = fmt.Errorf("%w: %v", goblinerr.ErrDeleteFailed, err)
		if log != nil {
			log.SpellFailed(err)
		}
		return err
	}

	if log != nil {
		log.DeleteSuccessful()
	}

	if err := hooks.AfterSacrifice(model, log); err != nil {
		return err
	}

	return nil
}

func scanRow(row *sql.Row, dest any) error {
	v := reflect.ValueOf(dest)
	if dest == nil {
		return fmt.Errorf("%w: destination cannot be nil", goblinerr.ErrInvalidModel)
	}

	if v.Kind() != reflect.Pointer || v.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("%w: destination must be pointer to struct", goblinerr.ErrInvalidModel)
	}

	structValue := v.Elem()

	values := make([]any, 0, structValue.NumField())

	for i := 0; i < structValue.NumField(); i++ {
		field := structValue.Field(i)

		if !field.CanSet() {
			continue
		}

		values = append(values, field.Addr().Interface())
	}

	return row.Scan(values...)
}
