package sqlite

import (
	"fmt"
	"goblinorm/internal/goblinerr"
	"reflect"
	"strings"

	"goblinorm/schema"
)

type Dialect struct{}

func (Dialect) CreateTableSQL(model any) (string, error) {
	parsed, err := schema.Parse(model)
	if err != nil {
		return "", err
	}

	columns := make([]string, 0, len(parsed.Fields))

	for _, field := range parsed.Fields {
		columnType, err := sqliteType(field.Type)
		if err != nil {
			return "", err
		}

		column := fmt.Sprintf("%s %s", field.ColumnName, columnType)

		if field.PrimaryKey {
			column += " PRIMARY KEY"
			if field.Type.Kind() == reflect.Int {
				column += " AUTOINCREMENT"
			}
		}

		columns = append(columns, column)
	}

	return fmt.Sprintf(
		"CREATE TABLE IF NOT EXISTS %s (%s);",
		parsed.TableName,
		strings.Join(columns, ", "),
	), nil
}

func (Dialect) InsertSQL(model any) (string, []any, error) {
	parsed, err := schema.Parse(model)
	if err != nil {
		return "", nil, err
	}

	var columns []string
	var placeholders []string
	var values []any

	for _, field := range parsed.Fields {
		if field.PrimaryKey && isZero(field.Value) {
			continue
		}

		columns = append(columns, field.ColumnName)
		placeholders = append(placeholders, "?")
		values = append(values, field.Value)
	}

	sql := fmt.Sprintf(
		"INSERT INTO %s (%s) VALUES (%s);",
		parsed.TableName,
		strings.Join(columns, ", "),
		strings.Join(placeholders, ", "),
	)

	return sql, values, nil
}

func (Dialect) SelectSQL(model any, condition string) (string, error) {
	parsed, err := schema.Parse(model)
	if err != nil {
		return "", err
	}

	sql := fmt.Sprintf("SELECT * FROM %s", parsed.TableName)

	if condition != "" {
		sql += " WHERE " + condition
	}

	sql += " LIMIT 1;"

	return sql, nil
}

func (Dialect) DeleteSQL(model any) (string, []any, error) {
	parsed, err := schema.Parse(model)
	if err != nil {
		return "", nil, err
	}

	for _, field := range parsed.Fields {
		if field.PrimaryKey {
			return fmt.Sprintf(
				"DELETE FROM %s WHERE %s = ?;",
				parsed.TableName,
				field.ColumnName,
			), []any{field.Value}, nil
		}
	}

	return "", nil, fmt.Errorf("%w: %s has no primary key field", goblinerr.ErrMissingPrimaryKey, parsed.Name)
}

func sqliteType(t reflect.Type) (string, error) {
	switch t.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return "INTEGER", nil
	case reflect.String:
		return "TEXT", nil
	case reflect.Bool:
		return "BOOLEAN", nil
	case reflect.Float32, reflect.Float64:
		return "REAL", nil
	default:
		return "", fmt.Errorf("%w: %s", goblinerr.ErrUnsupportedFieldType, t.String())
	}
}

func isZero(v any) bool {
	return reflect.ValueOf(v).IsZero()
}
