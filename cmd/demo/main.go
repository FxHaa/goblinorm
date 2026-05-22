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

func (w *Wizard) BeforeSummon() error {
	if w.Mana <= 0 {
		return fmt.Errorf("%s has no mana", w.Name)
	}

	return nil
}

func (w *Wizard) AfterSummon() error {
	fmt.Printf("%s has entered the dungeon.\n", w.Name)
	return nil
}

func (w *Wizard) AfterDivine() error {
	fmt.Printf("%s has been divined from the database.\n", w.Name)
	return nil
}

func (w *Wizard) BeforeSacrifice() error {
	fmt.Printf("%s prepares for sacrifice.\n", w.Name)
	return nil
}

func main() {
	dungeon, err := goblinorm.Open("sqlite", "./goblin.db")
	if err != nil {
		panic(err)
	}
	defer func(dungeon *goblinorm.Dungeon) {
		err := dungeon.Close()
		if err != nil {
			panic(err)
		}
	}(dungeon)

	dungeon.EnableEnhancedLogging()
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
		if errors.Is(err, goblinorm.ErrNoRows) {
			fmt.Println("No wizard found in this dungeon.")
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
