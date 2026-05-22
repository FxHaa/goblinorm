package schema

import (
	"fmt"
	"goblinorm/internal/goblinerr"
	"reflect"
	"strings"
)

type Field struct {
	Name       string
	ColumnName string
	Type       reflect.Type
	Value      any
	PrimaryKey bool
}

type Model struct {
	Fields    []Field
	Name      string
	TableName string
}

func Parse(model any) (*Model, error) {
	t := reflect.TypeOf(model)
	v := reflect.ValueOf(model)

	if model == nil {
		return nil, fmt.Errorf("%w: model cannot be nil", goblinerr.ErrInvalidModel)
	}

	if t.Kind() != reflect.Pointer {
		return nil, fmt.Errorf("%w: model must be a pointer", goblinerr.ErrInvalidModel)
	}

	if v.IsNil() {
		return nil, fmt.Errorf("%w: model pointer cannot be nil", goblinerr.ErrInvalidModel)
	}

	t = t.Elem()
	v = v.Elem()

	if t.Kind() != reflect.Struct {
		return nil, fmt.Errorf("%w: model must point to a struct", goblinerr.ErrInvalidModel)
	}

	parsed := &Model{
		Name:      t.Name(),
		TableName: pluralize(toSnakeCase(t.Name())),
	}

	for i := 0; i < t.NumField(); i++ {
		structField := t.Field(i)

		if !structField.IsExported() {
			continue
		}
		fieldValue := v.Field(i)

		field := Field{
			Name:       structField.Name,
			ColumnName: toSnakeCase(structField.Name),
			Type:       structField.Type,
			PrimaryKey: structField.Name == "ID",
		}

		if fieldValue.IsValid() && fieldValue.CanInterface() {
			field.Value = fieldValue.Interface()
		}

		tag := structField.Tag.Get("goblin")
		if strings.Contains(tag, "primary_key") {
			field.PrimaryKey = true
		}

		parsed.Fields = append(parsed.Fields, field)
	}

	return parsed, nil
}

func toSnakeCase(s string) string {
	if s == "" {
		return s
	}

	runes := []rune(s)
	var out []rune

	for i, r := range runes {
		if i > 0 && shouldInsertUnderscore(runes, i) {
			out = append(out, '_')
		}

		out = append(out, []rune(strings.ToLower(string(r)))...)
	}

	return string(out)
}

func shouldInsertUnderscore(runes []rune, index int) bool {
	current := runes[index]
	previous := runes[index-1]

	if !isUpper(current) {
		return false
	}

	if isLower(previous) || isDigit(previous) {
		return true
	}

	if isUpper(previous) && index+1 < len(runes) && isLower(runes[index+1]) {
		return true
	}

	return false
}

func isUpper(r rune) bool {
	return r >= 'A' && r <= 'Z'
}

func isLower(r rune) bool {
	return r >= 'a' && r <= 'z'
}

func isDigit(r rune) bool {
	return r >= '0' && r <= '9'
}

func pluralize(s string) string {
	if strings.HasSuffix(s, "s") {
		return s
	}

	return s + "s"
}
