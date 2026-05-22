package main

import (
	"errors"
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
		if errors.Is(err, goblinorm.ErrNoRows) {
			fmt.Println("No wizard found.")
			return
		}

		panic(err)
	}

	fmt.Printf("Found wizard: %+v\n", wizard)

	err = dungeon.Sacrifice(&wizard)
	if err != nil {
		panic(err)
	}

	fmt.Println("Wizard sacrificed.")
}
