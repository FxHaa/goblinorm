# GoblinORM

A chaotic fantasy-themed ORM for Go where databases are dungeons, structs are creatures, and SQL is forbidden magic.

## Features

- SQLite support
- Auto table creation from structs
- Insert, query, and delete records
- Reflection-based schema parsing
- Fantasy-themed logging
- SQL debug mode
- Lifecycle hooks
- Typed errors

## Example

```go
package main

import (
	"fmt"

	"goblinorm"
)

type Wizard struct {
	ID   int
	Name string
	Mana int
}

func main() {
	dungeon, err := goblinorm.Open("sqlite", ":memory:")
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := dungeon.Close(); err != nil {
			panic(err)
		}
	}()

	dungeon.Debug()

	err = dungeon.RaiseTheDead(&Wizard{})
	if err != nil {
		panic(err)
	}

	err = dungeon.Summon(&Wizard{
		Name: "Grug",
		Mana: 999,
	})
	if err != nil {
		panic(err)
	}

	var wizard Wizard

	err = dungeon.Divine(&wizard, "mana > ?", 100)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Found wizard: %+v\n", wizard)
}
```

## Hooks

```go
func (w *Wizard) BeforeSummon() error {
	if w.Mana <= 0 {
		return fmt.Errorf("%s has no mana", w.Name)
	}

	return nil
}
```

## Debug Output

```txt
[🧟] Raising table for Wizard...
[📜] Forbidden SQL revealed:
CREATE TABLE IF NOT EXISTS wizards (...);
[🧙] Summoning Wizard...
[🔥] Mana accepted.
[☠️] INSERT successful.
```

## Run Examples

```bash
go run ./examples/basic
go run ./examples/debug
go run ./examples/hooks
```

## Run Tests

```bash
go test ./...
```
