# 🧌 GoblinORM

> A chaotic fantasy-themed ORM for Go where databases are dungeons, structs are creatures, and SQL is forbidden magic.

GoblinORM is a tiny, cursed ORM for Go.

Summon structs. Divine records. Sacrifice rows. Raise tables from the dead.

## ✨ Features

- SQLite support
- auto table creation from structs
- insert, query, and delete records
- SQL debug mode
- lifecycle hooks
- typed errors
- fantasy-themed logs

## 🧙 Example

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

	err = dungeon.Summon(&Wizard{Name: "Grug", Mana: 999})
	if err != nil {
		panic(err)
	}

	var wizard Wizard
	err = dungeon.Divine(&wizard, "mana > ?", 100)
	if err != nil {
		panic(err)
	}

	fmt.Printf("The dungeon reveals: %+v\n", wizard)
}
```

## 🏰 Spells

| Spell | Meaning |
|---|---|
| `Open` | open a database |
| `RaiseTheDead` | create a table |
| `Summon` | insert a record |
| `Divine` | query a record |
| `Sacrifice` | delete a record |
| `Debug` | reveal generated SQL |

## 📜 Debug Output

```txt
[🧟] Raising table for Wizard...
[📜] Forbidden SQL revealed:
CREATE TABLE IF NOT EXISTS wizards (...);
[🧙] Summoning Wizard...
[🔥] Mana accepted.
[☠️] INSERT successful.
```

## 🪄 Hooks

```go
func (w *Wizard) BeforeSummon() error {
	if w.Mana <= 0 {
		return fmt.Errorf("%s has no mana", w.Name)
	}

	return nil
}
```

Supported hooks:

```go
BeforeSummon() error
AfterSummon() error
BeforeDivine() error
AfterDivine() error
BeforeSacrifice() error
AfterSacrifice() error
```

## 🧪 Examples

```bash
go run ./examples/basic
go run ./examples/debug
go run ./examples/hooks
```

## 🧹 Tests

```bash
go test ./...
```

## ⚠️ Warning

Do not use this to guard treasure.

Do not use this to run a bank.

Do not anger the goblins.
