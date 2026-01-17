package apihandler

import (
	"math/rand"
	"fmt"
)

type Pokemon struct {
	ID             int    `json:"id"`
	Name           string `json:"name"`
	BaseExperience int    `json:"base_experience"`
	Height         int    `json:"height"`
	IsDefault      bool   `json:"is_default"`
	Order          int    `json:"order"`
	Weight         int    `json:"weight"`
	Abilities      []struct {
		IsHidden bool `json:"is_hidden"`
		Slot     int  `json:"slot"`
		Ability  struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"ability"`
	} `json:"abilities"`
	Stats []struct {
		BaseStat int `json:"base_stat"`
		Effort   int `json:"effort"`
		Stat     struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"stat"`
	} `json:"stats"`
	Types []struct {
		Slot int `json:"slot"`
		Type struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"type"`
	}
}

type Pokedex struct {
	Data map[string]Pokemon
}

func catchPokemon(target string) {
	fmt.Println("Throwing a Pokeball at " + target + "...")
	if targetData == nil {
		fmt.Println("Error finding target.")
		return
	}
	if ( rollCatchChance(targetData.BaseExperience) ) {
		pdex.Data[target] = *(targetData)
		fmt.Printf("%s was caught!\n", target)
		fmt.Println("You may now inspect it with the inspect command.")
	} else {
		fmt.Printf("%s escaped!\n", target)
	}

}

func inspectPokemon(target string) {
	if pmon, ok := (*pdex).Data[target]; !ok {
		fmt.Println("you have not caught that pokemon.")

	} else {
		fmt.Printf("Name: %s\n", pmon.Name)
		fmt.Printf("Weight: %d\n", pmon.Weight)
		fmt.Printf("Stats: \n")
		for _, stat := range pmon.Stats {
			fmt.Printf("-%s: %d\n", stat.Stat.Name, stat.BaseStat)
		}
		fmt.Printf("Types: \n")
		for _, types := range pmon.Types {
			fmt.Printf("- %s\n", types.Type.Name)
		}
		return
	}
}

func displayPokedex() {
	if len(pdex.Data) == 0 {
		fmt.Println("You have not caught any Pokemon.")
		return

	} else {
		fmt.Println("Your Pokedex:")
		for _, pmon := range pdex.Data {
			fmt.Printf("- %s\n", pmon.Name)
		}
	}
	return
}

func rollCatchChance(baseExp int) bool {

	var caught bool = false

	
	denom := (baseExp / 25)
	rollChance := []bool{}
	for i := 0; i < denom; i++ {
		rollChance = append(rollChance, false)
	}

	rollChance[rand.Intn(denom)] = true

	roll := rand.Intn(denom)

	fmt.Printf("Chance to catch is 1 in %d\n", len(rollChance))
	if rollChance[roll] {
		fmt.Println("Caught")
		caught = true
	} else {
		caught = false
	}
	
	return caught
}