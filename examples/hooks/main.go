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

	fmt.Printf("[model hook] %s gathers mana before summoning.\n", w.Name)
	return nil
}

func (w *Wizard) AfterSummon() error {
	fmt.Printf("[model hook] %s has entered the dungeon.\n", w.Name)
	return nil
}

func (w *Wizard) AfterDivine() error {
	fmt.Printf("[model hook] %s was discovered by divination.\n", w.Name)
	return nil
}

func (w *Wizard) BeforeSacrifice() error {
	fmt.Printf("[model hook] %s prepares for sacrifice.\n", w.Name)
	return nil
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
		Name: "Eldrin",
		Mana: 1200,
	})
	if err != nil {
		if errors.Is(err, goblinorm.ErrHookFailed) {
			fmt.Printf("Hook failed: %v\n", err)
			return
		}

		panic(err)
	}

	var wizard Wizard

	err = dungeon.Divine(&wizard, "mana > ?", 100)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Hook example found wizard: %+v\n", wizard)

	err = dungeon.Sacrifice(&wizard)
	if err != nil {
		panic(err)
	}
}
