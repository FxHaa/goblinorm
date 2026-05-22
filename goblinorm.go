package goblinorm

import (
	"database/sql"
	"fmt"
	"goblinorm/internal/goblinerr"
	"goblinorm/internal/logger"
	"io"

	_ "modernc.org/sqlite"

	"goblinorm/dialects/sqlite"
	"goblinorm/migration"
	"goblinorm/query"
)

type Dungeon struct {
	db      *sql.DB
	dialect Dialect
	logger  *logger.Logger
}

type Dialect interface {
	CreateTableSQL(model any) (string, error)
	InsertSQL(model any) (string, []any, error)
	SelectSQL(model any, condition string) (string, error)
	DeleteSQL(model any) (string, []any, error)
}

var (
	ErrInvalidModel         = goblinerr.ErrInvalidModel
	ErrNoRows               = goblinerr.ErrNoRows
	ErrMissingPrimaryKey    = goblinerr.ErrMissingPrimaryKey
	ErrUnsupportedFieldType = goblinerr.ErrUnsupportedFieldType
	ErrUnsupportedDriver    = goblinerr.ErrUnsupportedDriver
	ErrScanFailed           = goblinerr.ErrScanFailed
	ErrMigrationFailed      = goblinerr.ErrMigrationFailed
	ErrInsertFailed         = goblinerr.ErrInsertFailed
	ErrSelectFailed         = goblinerr.ErrSelectFailed
	ErrDeleteFailed         = goblinerr.ErrDeleteFailed
	ErrHookFailed           = goblinerr.ErrHookFailed
)

func Open(driverName, dataSourceName string) (*Dungeon, error) {
	db, err := sql.Open(driverName, dataSourceName)
	if err != nil {
		return nil, err
	}

	var dialect Dialect

	switch driverName {
	case "sqlite3", "sqlite":
		dialect = sqlite.Dialect{}
	default:
		return nil, fmt.Errorf("unsupported driver: %s", driverName)
	}

	return &Dungeon{
		db:      db,
		dialect: dialect,
		logger:  logger.New(),
	}, nil
}

func (d *Dungeon) Close() error {
	return d.db.Close()
}

func (d *Dungeon) EnableEnhancedLogging() {
	d.logger.Enable()
}

func (d *Dungeon) DisableEnhancedLogging() {
	d.logger.Disable()
}

func (d *Dungeon) Debug() {
	d.logger.EnableDebug()
}

func (d *Dungeon) DisableDebug() {
	d.logger.DisableDebug()
}

func (d *Dungeon) SetLogOutput(out io.Writer) {
	d.logger.SetOutput(out)
}

func (d *Dungeon) RaiseTheDead(model any) error {
	return migration.RaiseTheDead(d.db, d.dialect, d.logger, model)
}

func (d *Dungeon) Summon(model any) error {
	return query.Summon(d.db, d.dialect, d.logger, model)
}

func (d *Dungeon) Divine(dest any, condition string, args ...any) error {
	return query.Divine(d.db, d.dialect, d.logger, dest, condition, args...)
}

func (d *Dungeon) Sacrifice(model any) error {
	return query.Sacrifice(d.db, d.dialect, d.logger, model)
}
