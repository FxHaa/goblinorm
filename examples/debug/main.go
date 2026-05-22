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
		Name: "Morg",
		Mana: 777,
	})
	if err != nil {
		panic(err)
	}

	var wizard Wizard

	err = dungeon.Divine(&wizard, "mana > ?", 500)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Debug mode found wizard: %+v\n", wizard)
}
