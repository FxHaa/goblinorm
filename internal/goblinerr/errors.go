package goblinerr

import "errors"

var (
	ErrInvalidModel         = errors.New("invalid creature model")
	ErrNoRows               = errors.New("no creatures found in this dungeon")
	ErrMissingPrimaryKey    = errors.New("creature has no primary rune")
	ErrUnsupportedFieldType = errors.New("forbidden field type")
	ErrUnsupportedDriver    = errors.New("unsupported dungeon driver")
	ErrScanFailed           = errors.New("failed to bind dungeon loot into creature")
	ErrMigrationFailed      = errors.New("failed to raise table from the dead")
	ErrInsertFailed         = errors.New("failed to summon creature")
	ErrSelectFailed         = errors.New("failed to divine creature")
	ErrDeleteFailed         = errors.New("failed to sacrifice creature")
	ErrHookFailed           = errors.New("creature hook failed")
)
