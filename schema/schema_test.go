package schema

import (
	"errors"
	"reflect"
	"testing"

	"goblinorm/internal/goblinerr"
)

type testWizard struct {
	ID   int
	Name string
	Mana int
}

type testGoblin struct {
	GoblinID int    `goblin:"primary_key"`
	Name     string `goblin:"size:255"`
}

type testDragon struct {
	ID         int
	FireName   string
	IsAncient  bool
	PowerLevel float64
}

type testSecretCreature struct {
	ID     int
	Name   string
	secret string
}

func TestParseModel(t *testing.T) {
	wizard := &testWizard{
		ID:   7,
		Name: "Grug",
		Mana: 999,
	}

	model, err := Parse(wizard)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if model.Name != "testWizard" {
		t.Fatalf("expected model name testWizard, got %s", model.Name)
	}

	if model.TableName != "test_wizards" {
		t.Fatalf("expected table name test_wizards, got %s", model.TableName)
	}

	if len(model.Fields) != 3 {
		t.Fatalf("expected 3 fields, got %d", len(model.Fields))
	}

	idField := model.Fields[0]
	if idField.Name != "ID" {
		t.Fatalf("expected first field name ID, got %s", idField.Name)
	}

	if idField.ColumnName != "id" {
		t.Fatalf("expected ID column name id, got %s", idField.ColumnName)
	}

	if idField.Type.Kind() != reflect.Int {
		t.Fatalf("expected ID type int, got %s", idField.Type.Kind())
	}

	if idField.Value != 7 {
		t.Fatalf("expected ID value 7, got %v", idField.Value)
	}

	if !idField.PrimaryKey {
		t.Fatal("expected ID field to be primary key")
	}

	nameField := model.Fields[1]
	if nameField.Name != "Name" {
		t.Fatalf("expected second field name Name, got %s", nameField.Name)
	}

	if nameField.ColumnName != "name" {
		t.Fatalf("expected Name column name name, got %s", nameField.ColumnName)
	}

	if nameField.Value != "Grug" {
		t.Fatalf("expected Name value Grug, got %v", nameField.Value)
	}

	if nameField.PrimaryKey {
		t.Fatal("expected Name field not to be primary key")
	}

	manaField := model.Fields[2]
	if manaField.Name != "Mana" {
		t.Fatalf("expected third field name Mana, got %s", manaField.Name)
	}

	if manaField.ColumnName != "mana" {
		t.Fatalf("expected Mana column name mana, got %s", manaField.ColumnName)
	}

	if manaField.Value != 999 {
		t.Fatalf("expected Mana value 999, got %v", manaField.Value)
	}
}

func TestParsePrimaryKeyTag(t *testing.T) {
	model, err := Parse(&testGoblin{
		GoblinID: 42,
		Name:     "Snarg",
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(model.Fields) != 2 {
		t.Fatalf("expected 2 fields, got %d", len(model.Fields))
	}

	field := model.Fields[0]

	if field.Name != "GoblinID" {
		t.Fatalf("expected field GoblinID, got %s", field.Name)
	}

	if field.ColumnName != "goblin_id" {
		t.Fatalf("expected column name goblin_id, got %s", field.ColumnName)
	}

	if !field.PrimaryKey {
		t.Fatal("expected goblin primary_key tag to mark field as primary key")
	}
}

func TestParseConvertsFieldNamesToSnakeCase(t *testing.T) {
	model, err := Parse(&testDragon{})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	tests := map[string]string{
		"ID":         "id",
		"FireName":   "fire_name",
		"IsAncient":  "is_ancient",
		"PowerLevel": "power_level",
	}

	for _, field := range model.Fields {
		expectedColumn, ok := tests[field.Name]
		if !ok {
			t.Fatalf("unexpected field %s", field.Name)
		}

		if field.ColumnName != expectedColumn {
			t.Fatalf("expected column for %s to be %s, got %s", field.Name, expectedColumn, field.ColumnName)
		}
	}
}

func TestParseSkipsUnexportedFields(t *testing.T) {
	model, err := Parse(&testSecretCreature{
		ID:     1,
		Name:   "Hidden Goblin",
		secret: "forbidden",
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(model.Fields) != 2 {
		t.Fatalf("expected 2 exported fields, got %d", len(model.Fields))
	}

	for _, field := range model.Fields {
		if field.Name == "secret" {
			t.Fatal("expected unexported field secret to be skipped")
		}
	}
}

func TestParseRejectsNilModel(t *testing.T) {
	_, err := Parse(nil)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if !errors.Is(err, goblinerr.ErrInvalidModel) {
		t.Fatalf("expected ErrInvalidModel, got %v", err)
	}
}

func TestParseRejectsNilPointer(t *testing.T) {
	var wizard *testWizard

	_, err := Parse(wizard)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if !errors.Is(err, goblinerr.ErrInvalidModel) {
		t.Fatalf("expected ErrInvalidModel, got %v", err)
	}
}

func TestParseRejectsNonPointer(t *testing.T) {
	_, err := Parse(testWizard{})
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if !errors.Is(err, goblinerr.ErrInvalidModel) {
		t.Fatalf("expected ErrInvalidModel, got %v", err)
	}
}

func TestParseRejectsPointerToNonStruct(t *testing.T) {
	value := 10

	_, err := Parse(&value)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if !errors.Is(err, goblinerr.ErrInvalidModel) {
		t.Fatalf("expected ErrInvalidModel, got %v", err)
	}
}

func TestToSnakeCase(t *testing.T) {
	tests := map[string]string{
		"Wizard":         "wizard",
		"FireWizard":     "fire_wizard",
		"ManaPool":       "mana_pool",
		"ID":             "id",
		"GoblinID":       "goblin_id",
		"HTTPSpell":      "http_spell",
		"SpellURL":       "spell_url",
		"URLValue":       "url_value",
		"Mana2Spell":     "mana2_spell",
		"Spell2Mana":     "spell2_mana",
		"XMLHTTPRequest": "xmlhttp_request",
	}

	for input, expected := range tests {
		got := toSnakeCase(input)
		if got != expected {
			t.Fatalf("expected %s to become %s, got %s", input, expected, got)
		}
	}
}

func TestPluralize(t *testing.T) {
	tests := map[string]string{
		"wizard": "wizards",
		"boss":   "boss",
		"goblin": "goblins",
	}

	for input, expected := range tests {
		got := pluralize(input)
		if got != expected {
			t.Fatalf("expected %s to become %s, got %s", input, expected, got)
		}
	}
}
